package models

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/httplib"
	"googlemaps.github.io/maps"
)

const (
	// Map Tool
	MapToolName     = "map"
	MapToolAlias    = "m"
	MapToolEndpoint = "/map"
	MapHelpMsg      = `map is a direction tool.

Usage:

1. 规划交通路线
	map direct PlaceA to PlaceB
		
Example 
	map direct 杭州火车东站 to 武林广场

2. 设置Home地址
	map set home Place	

3. 查询Home地址
	map get home

4. 规划回家交通路线	
	map go home Place
	or
	直接发送位置信息 
	
This tool is powered by Google Maps.`
)

type MapTool struct {
	toolBase
	Origin      string
	Destination string
	UserID      string
	HomeAddress string
	Latlng      string
	HelpMsg     string
}

type Routes struct {
	//Routes []maps.Route `json:"routes"`
	Routes []Route
}

type Route struct {
	// Summary contains a short textual description for the route, suitable for
	// naming and disambiguating the route from alternatives.
	Summary string `json:"summary"`

	// Legs contains information about a leg of the route, between two locations within the
	// given route. A separate leg will be present for each waypoint or destination specified.
	// (A route with no waypoints will contain exactly one leg within the legs array.)
	Legs []*Leg `json:"legs"`

	// Copyrights contains the copyrights text to be displayed for this route. You must handle
	// and display this information yourself.
	Copyrights string `json:"copyrights"`

	// Warnings contains an array of warnings to be displayed when showing these directions.
	// You must handle and display these warnings yourself.
	Warnings []string `json:"warnings"`
}

// Leg represents a single leg of a route.
type Leg struct {
	// Steps contains an array of steps denoting information about each separate step of the
	// leg of the journey.
	Steps []*Step `json:"steps"`

	// Distance indicates the total distance covered by this leg.
	Distance `json:"distance"`

	// Duration indicates total time required for this leg.
	Duration `json:"duration"`

	// StartAddress contains the human-readable address (typically a street address)
	// reflecting the start location of this leg.
	StartAddress string `json:"start_address"`

	// EndAddress contains the human-readable address (typically a street address)
	// reflecting the end location of this leg.
	EndAddress string `json:"end_address"`
}

// Step represents a single step of a leg.
type Step struct {
	// HTMLInstructions contains formatted instructions for this step, presented as an HTML text string.
	HTMLInstructions string `json:"html_instructions"`

	// Distance contains the distance covered by this step until the next step.
	Distance `json:"distance"`

	// Duration contains the typical time required to perform the step, until the next step.
	Duration `json:"duration"`

	// Steps contains detailed directions for walking or driving steps in transit directions. Substeps
	// are only available when travel_mode is set to "transit". The inner steps array is of the same
	// type as steps.
	Steps []*Step `json:"steps"`

	// TransitDetails contains transit specific information. This field is only returned with travel
	// mode is set to "transit".
	TransitDetails *TransitDetails `json:"transit_details"`

	// TravelMode indicates the travel mode of this step.
	TravelMode string `json:"travel_mode"`
}

// Distance is the API representation for a distance between two points.
type Distance struct {
	// HumanReadable is the human friendly distance. This is rounded and in an appropriate unit for the
	// request. The units can be overriden with a request parameter.
	HumanReadable string `json:"text"`
	// Meters is the numeric distance, always in meters. This is intended to be used only in
	// algorithmic situations, e.g. sorting results by some user specified metric.
	Meters int `json:"value"`
}

// TransitDetails contains additional information about the transit stop, transit line and transit agency.
type TransitDetails struct {
	// ArrivalStop contains information about the stop/station for this part of the trip.
	ArrivalStop TransitStop `json:"arrival_stop"`
	// DepartureStop contains information about the stop/station for this part of the trip.
	DepartureStop TransitStop `json:"departure_stop"`
	// Headsign specifies the direction in which to travel on this line, as it is marked on the vehicle or at the departure stop.
	Headsign string `json:"headsign"`
	// Headway specifies the expected number of seconds between departures from the same stop at this time
	Headway time.Duration `json:"headway"`
	// NumStops contains the number of stops in this step, counting the arrival stop, but not the departure stop
	NumStops uint `json:"num_stops"`
	// Line contains information about the transit line used in this step
	Line TransitLine `json:"line"`
}

// TransitLine contains information about the transit line used in this step
type TransitLine struct {
	// Name contains the full name of this transit line. eg. "7 Avenue Express".
	Name string `json:"name"`
	// ShortName contains the short name of this transit line.
	ShortName string `json:"short_name"`
	// Color contains the color commonly used in signage for this transit line.
	Color string `json:"color"`

	// URL contains the URL for this transit line as provided by the transit agency
	URL *url.URL `json:"url"`
	// Icon contains the URL for the icon associated with this line
	Icon *url.URL `json:"icon"`
	// TextColor contains the color of text commonly used for signage of this line
	TextColor string `json:"text_color"`
}

// TransitStop contains information about the stop/station for this part of the trip.
type TransitStop struct {
	// Name of the transit station/stop. eg. "Union Square".
	Name string `json:"name"`
}

type Duration struct {
	Value int64  `json:"value"`
	Text  string `json:"text"`
}

var (
	userHomeConfig config.Configer
)

func init() {
	var err error
	userHomeConfig, err = config.NewConfig("ini", "conf/userhome.conf")
	if err != nil {
		beego.Error("Failed to load userhome.conf file.")
	}
}

func (m *MapTool) NewTool(req Request) {
	m.name = MapToolName
	m.alias = MapToolAlias
	m.endpoint = MapToolAlias
	m.HelpMsg = MapHelpMsg
	m.UserID = req.FromUserName
}

func (m *MapTool) Directions() TextResponse {
	var textResp TextResponse
	var originPlaceID, destinationPlaceID string
	textResp.MsgType = MsgTypeText

	if m.Origin != "" && m.Destination != "" {
		var err error
		originPlaceID, _, err = getPlaceID(m.Origin)
		if err != nil {
			textResp.Content = fmt.Sprintf("查找地点失败，请尝试其他地点关键词。")
			return textResp
		}
		destinationPlaceID, _, err = getPlaceID(m.Destination)
		if err != nil {
			textResp.Content = fmt.Sprintf("查找地点失败，请尝试其他地点关键词。")
			return textResp
		}
	} else {
		textResp.Content = fmt.Sprintf("地点关键词不能为空。")
		return textResp
	}
	//originPlaceID = "ChIJWYij7kicTDQRCp51F2RCKfM"
	//destinationPlaceID = "ChIJwWnPHVdiSzQRN7O4WYYFC14"
	directionsStr, err := getDirections(originPlaceID, destinationPlaceID)
	if err != nil {
		textResp.Content = err.Error()
		return textResp
	}
	beego.Info("Directions Info:", directionsStr)
	textResp.Content = directionsStr
	return textResp
}

func (m *MapTool) SetHome() TextResponse {
	var textResp TextResponse
	textResp.MsgType = MsgTypeText
	if userHomeConfig.String(m.UserID+"::id") != "" {
		beego.Info("User ", m.UserID, "address is existed. Overwriting..")
	}

	homePlaceID, address, err := getPlaceID(m.HomeAddress)
	if err != nil {
		textResp.Content = fmt.Sprintf("设置Home地址失败，请尝试其他地址关键词。")
		return textResp
	}
	//homePlaceID := "idididdid"
	//address := "addressaddress"
	if err := userHomeConfig.Set(m.UserID+"::id", homePlaceID); err != nil {
		textResp.Content = fmt.Sprintf("设置Home地址失败。")
		return textResp
	}

	if err := userHomeConfig.Set(m.UserID+"::address", address); err != nil {
		textResp.Content = fmt.Sprintf("设置Home地址失败。")
		return textResp
	}

	if err := userHomeConfig.SaveConfigFile("conf/userhome.conf"); err != nil {
		textResp.Content = fmt.Sprintf("设置Home地址失败。")
		return textResp
	}
	textResp.Content = fmt.Sprintf("设置Home地址成功：%s", address)
	beego.Info("Set user home address:", m.UserID, homePlaceID, address)
	return textResp
}

func (m *MapTool) GetHome() TextResponse {
	var textResp TextResponse
	textResp.MsgType = MsgTypeText
	homePlaceID := userHomeConfig.String(m.UserID + "::id")
	address := userHomeConfig.String(m.UserID + "::address")

	if homePlaceID == "" || address == "" {
		textResp.Content = fmt.Sprintf("用户还未设置Home地址，使用map set home来设置Home地址。")
		return textResp
	}

	textResp.Content = fmt.Sprintf("Home地址：%s", address)
	beego.Info("Get user home address:", m.UserID, homePlaceID, address)
	return textResp
}

func (m *MapTool) GoHome() TextResponse {
	var textResp TextResponse
	var originPlaceID string
	var err error
	textResp.MsgType = MsgTypeText
	homePlaceID := userHomeConfig.String(m.UserID + "::id")
	address := userHomeConfig.String(m.UserID + "::address")

	if homePlaceID == "" || address == "" {
		textResp.Content = fmt.Sprintf("用户还未设置Home地址，使用map set home来设置Home地址。")
		return textResp
	}

	// Use nearby search first if Latlng exists
	if m.Latlng != "" {
		originPlaceID, err = getPleaceNearby(m.Origin, m.Latlng)
		if err != nil {
			originPlaceID, _, err = getPlaceID(m.Origin)
		}
	} else {
		originPlaceID, _, err = getPlaceID(m.Origin)
	}
	if err != nil {
		textResp.Content = fmt.Sprintf("查找地点失败，请尝试其他地点关键词。")
		return textResp
	}

	directionsStr, err := getDirections(originPlaceID, homePlaceID)
	if err != nil {
		textResp.Content = err.Error()
		return textResp
	}
	textResp.Content = directionsStr
	return textResp
}

func getPlaceID(keyword string) (placeID string, address string, err error) {
	var placeSearchResult *maps.PlacesSearchResponse
	req := httplib.Get(APIADDRESS + APIVERSION + MapToolEndpoint + "/place/search")
	req.Param("keyword", keyword)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err = req.ToJSON(&placeSearchResult)
	if err != nil {
		return "", "", err
	}
	if len(placeSearchResult.Results) < 1 {
		return "", "", errors.New("place not found")
	}
	beego.Info("Place Search Result: name:", placeSearchResult.Results[0].Name, "address:", placeSearchResult.Results[0].FormattedAddress, "PlaceID:", placeSearchResult.Results[0].PlaceID)
	placeID = placeSearchResult.Results[0].PlaceID
	address = placeSearchResult.Results[0].FormattedAddress
	return placeID, address, nil
}

func getPleaceNearby(keyword, latlng string) (placeID string, err error) {
	var placeSearchResult *maps.PlacesSearchResponse
	req := httplib.Get(APIADDRESS + APIVERSION + MapToolEndpoint + "/place/nearby")
	req.Param("keyword", keyword)
	req.Param("latlng", latlng)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err = req.ToJSON(&placeSearchResult)
	if err != nil {
		return "", err
	}
	if len(placeSearchResult.Results) < 1 {
		return "", errors.New("place not found")
	}
	beego.Info("Nearby Place Search Result: name:", placeSearchResult.Results[0].Name, "address:", placeSearchResult.Results[0].FormattedAddress, "PlaceID:", placeSearchResult.Results[0].PlaceID)
	placeID = placeSearchResult.Results[0].PlaceID
	return placeID, nil
}

func getDirections(originID string, destinationID string) (string, error) {
	var response *Routes
	var resultStr string
	req := httplib.Get(APIADDRESS + APIVERSION + MapToolEndpoint + "/direct/transit")
	req.Param("origin", originID)
	req.Param("destination", destinationID)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&response)
	if err != nil {
		beego.Trace("error:", err.Error())
		return "", err
	}

	if len(response.Routes) < 1 {
		return "", errors.New("查询线路失败，无可用线路。")
	}
	directionsResult := response.Routes
	if len(directionsResult[0].Legs) >= 1 {
		resultStr = fmt.Sprintf("路线总长%s，预计用时%s\n◇ %s\n", directionsResult[0].Legs[0].Distance.HumanReadable,
			directionsResult[0].Legs[0].Duration.Text,
			directionsResult[0].Legs[0].StartAddress)
		for i := 0; i < len(directionsResult[0].Legs[0].Steps); i++ {
			if directionsResult[0].Legs[0].Steps[i].TravelMode == "WALKING" {
				resultStr += fmt.Sprintf("    %s\n", directionsResult[0].Legs[0].Steps[i].HTMLInstructions)
				resultStr += fmt.Sprintf("    %s %s\n", directionsResult[0].Legs[0].Steps[i].Distance.HumanReadable, directionsResult[0].Legs[0].Steps[i].Duration.Text)
			} else if directionsResult[0].Legs[0].Steps[i].TravelMode == "TRANSIT" {
				resultStr += fmt.Sprintf("◇ %s\n", directionsResult[0].Legs[0].Steps[i].TransitDetails.DepartureStop.Name)
				resultStr += fmt.Sprintf("    %s %s %d站\n", directionsResult[0].Legs[0].Steps[i].HTMLInstructions,
					directionsResult[0].Legs[0].Steps[i].TransitDetails.Line.ShortName,
					directionsResult[0].Legs[0].Steps[i].TransitDetails.NumStops)
				resultStr += fmt.Sprintf("    %s %s\n", directionsResult[0].Legs[0].Steps[i].Distance.HumanReadable, directionsResult[0].Legs[0].Steps[i].Duration.Text)
				resultStr += fmt.Sprintf("◇ %s\n", directionsResult[0].Legs[0].Steps[i].TransitDetails.ArrivalStop.Name)
			}
		}
		resultStr += fmt.Sprintf("◇ %s\n", directionsResult[0].Legs[0].EndAddress)
	} else {
		return "", errors.New("查询线路失败，无可用线路。")
	}
	return resultStr, nil
}
