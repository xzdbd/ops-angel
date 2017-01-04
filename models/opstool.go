package models

import (
	"crypto/tls"

	"strconv"

	"strings"

	"fmt"

	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/docker/go-dockercloud/dockercloud"
	"googlemaps.github.io/maps"
)

type Tool interface {
	NewTool()
}

type toolBase struct {
	name     string
	alias    string
	endpoint string
}

type GoogleTool struct {
	toolBase
	Key     string
	N       int
	HelpMsg string
}

type GoogleResult struct {
	Abstract  string
	Title     string
	URL       string
	Sitelinks []Sitelinks
}

type Sitelinks struct {
	Index    string
	Abstract string
	Title    string
	URL      string
}

type DockerCloudTool struct {
	toolBase
	ServiceName string
	Action      string
	Privileged  bool
	HelpMsg     string
}

type MapTool struct {
	toolBase
	Origin      string
	Destination string
	HomeAddress string
	HelpMsg     string
}

const (
	APIADDRESS = "https://api.xzdbd.com/"
	//APIADDRESS = "http://11.11.1.6:8098/"
	APIVERSION = "v1"

	// Google Tool
	GoogleToolName     = "google"
	GoogleToolAlias    = "g"
	GoogleToolEndpoint = "/search/google"
	GoogleHelpMsg      = `google is a google search tool. Enjoy!

Usage:
	google KEY 
or
	g KEY 

KEY:
	search key words.

Example:
	google happy day	
	`

	// Docker Cloud Tool
	DockerCloudToolName     = "dockercloud"
	DockerCloudToolAlias    = "dc"
	DockerCloudToolEndpoint = "/dockercloud"
	DockerCloudHelpMsg      = `dockercloud is an operations tool.

Usage:
	dockercloud service NAME [status]|start|stop
or
	dc service NAME [status]|start|stop
		
Example 
	dc service test status`

	// Map Tool
	MapToolName     = "map"
	MapToolAlias    = "m"
	MapToolEndpoint = "/map"
	MapHelpMsg      = `map is a direction tool.

Usage:
	dockercloud service NAME [status]|start|stop
or
	dc service NAME [status]|start|stop
		
Example 
	dc service test status`
)

var (
	apiuser     = beego.AppConfig.String("apiuser")
	apipassword = beego.AppConfig.String("apipassword")
)

func (g *GoogleTool) NewTool() {
	g.name = GoogleToolName
	g.alias = GoogleToolAlias
	g.endpoint = GoogleToolEndpoint
	g.HelpMsg = GoogleHelpMsg
}

func (g *GoogleTool) Run() (NewsResponse, error) {
	var GoogleResultList []*GoogleResult
	GoogleResultList = make([]*GoogleResult, g.N)
	var newsResp NewsResponse
	newsResp.MsgType = MsgTypeNews

	req := httplib.Get(APIADDRESS + APIVERSION + g.endpoint)
	req.Param("key", g.Key)
	req.Param("n", strconv.Itoa(g.N))
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	err := req.ToJSON(&GoogleResultList)
	if err != nil {
		return newsResp, err
	}

	newsResp.ArticleCount = len(GoogleResultList)
	for i := 0; i < newsResp.ArticleCount; i++ {
		//picUrl := getFavicons(GoogleResultList[i].URL)
		picUrl := "https://upload.wikimedia.org/wikipedia/commons/thumb/5/53/Google_%22G%22_Logo.svg/200px-Google_%22G%22_Logo.svg.png"
		item := Item{Title: GoogleResultList[i].Title, Description: GoogleResultList[i].Abstract, Url: GoogleResultList[i].URL, PicUrl: picUrl}
		newsResp.Articles = append(newsResp.Articles, &item)
	}

	return newsResp, nil
}

func (dc *DockerCloudTool) NewTool() {
	dc.name = DockerCloudToolName
	dc.alias = DockerCloudToolAlias
	dc.endpoint = DockerCloudToolAlias
	dc.HelpMsg = DockerCloudHelpMsg
}

func (dc *DockerCloudTool) Run() (TextResponse, error) {
	var dcList dockercloud.SListResponse
	var textResp TextResponse
	textResp.MsgType = MsgTypeText

	switch dc.Action {
	case "status":
		if dc.ServiceName != "" {
			if strings.ToLower(dc.ServiceName) == "all" {
				var err error
				dcList, err = getAllDockerCloudService()
				if err != nil {
					return textResp, err
				}
				textResp.Content = fmt.Sprintf("共有%d个服务。\n", dcList.Meta.TotalCount)
				for i := 0; i < dcList.Meta.TotalCount; i++ {
					textResp.Content += fmt.Sprintf("%d. %s: %s\n", i+1, dcList.Objects[i].Name, dcList.Objects[i].State)
				}
			} else {
				var err error
				dcList, err = getDockerCloudServiceByName(dc.ServiceName)
				if err != nil {
					return textResp, err
				}
				if dcList.Meta.TotalCount < 1 {
					textResp.Content = fmt.Sprintf("没有找到名称为%s的服务。", dc.ServiceName)
				} else {
					textResp.Content = fmt.Sprintf("%s: %s\n", dcList.Objects[0].Name, dcList.Objects[0].State)
				}
			}
		}
	case "start":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "start")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务启动错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务启动成功，请稍后查看该服务状态。")
			}
		}
	case "stop":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "stop")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务停止成功，请稍后查看该服务状态。")
			}
		}
	case "redeploy":
		if dc.ServiceName != "" {
			_, err := actionDockerCloudService(dc.ServiceName, "redeploy")
			if err != nil {
				textResp.Content = fmt.Sprintf("服务重新部署错误，错误信息：%s\n", err.Error())
			} else {
				textResp.Content = fmt.Sprintf("服务重新部署成功，请稍后查看该服务状态。")
			}
		}
	default:
		return textResp, errors.New("Invalid Action. Valid actions are 'start', 'stop' and 'status'")
	}
	return textResp, nil
}

func (m *MapTool) NewTool() {
	m.name = MapToolName
	m.alias = MapToolAlias
	m.endpoint = MapToolAlias
	m.HelpMsg = MapHelpMsg
}

func (m *MapTool) Directions() (TextResponse, error) {
	var textResp TextResponse
	var originPlaceID, destinationPlaceID string
	textResp.MsgType = MsgTypeText

	if m.Origin != "" && m.Destination != "" {
		var err error
		originPlaceID, err = getPlaceID(m.Origin)
		if err != nil {
			textResp.Content = fmt.Sprintf("查找地点失败，请尝试其他地点关键词。")
			return textResp, nil
		}
		destinationPlaceID, err = getPlaceID(m.Destination)
		if err != nil {
			textResp.Content = fmt.Sprintf("查找地点失败，请尝试其他地点关键词。")
			return textResp, nil
		}
	} else {
		textResp.Content = fmt.Sprintf("地点关键词不能为空。")
		return textResp, nil
	}
	//originPlaceID = "ChIJWYij7kicTDQRCp51F2RCKfM"
	//destinationPlaceID = "ChIJwWnPHVdiSzQRN7O4WYYFC14"
	directionsStr, err := getDirections(originPlaceID, destinationPlaceID)
	if err != nil {
		textResp.Content = err.Error()
		return textResp, nil
	}
	beego.Info("Directions Info:", directionsStr)
	textResp.Content = directionsStr
	return textResp, nil
}

func getPlaceID(keyword string) (string, error) {
	var placeSearchResult *maps.PlacesSearchResponse
	req := httplib.Get(APIADDRESS + APIVERSION + MapToolEndpoint + "/place/search")
	req.Param("keyword", keyword)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&placeSearchResult)
	if err != nil {
		return "", err
	}
	if len(placeSearchResult.Results) < 1 {
		return "", errors.New("place not found")
	}
	beego.Info("Place Search Result: name:", placeSearchResult.Results[0].Name, "address:", placeSearchResult.Results[0].FormattedAddress, "PlaceID:", placeSearchResult.Results[0].PlaceID)
	placeID := placeSearchResult.Results[0].PlaceID
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

func getAllDockerCloudService() (dockercloud.SListResponse, error) {
	var dcList dockercloud.SListResponse
	req := httplib.Get(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service")
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&dcList)
	if err != nil {
		return dcList, err
	}
	return dcList, nil
}

func getDockerCloudServiceByName(name string) (dockercloud.SListResponse, error) {
	var dcList dockercloud.SListResponse
	req := httplib.Get(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + name)
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err := req.ToJSON(&dcList)
	if err != nil {
		return dcList, err
	}
	return dcList, nil
}

func getDockerCloudServiceUuid(name string) (string, error) {
	dcList, err := getDockerCloudServiceByName(name)
	if err != nil {
		return "", err
	}
	if dcList.Meta.TotalCount < 1 {
		return "", fmt.Errorf("没有找到名称为%s的服务。", name)
	}
	return dcList.Objects[0].Uuid, nil
}

// start, stop and redeploy service
func actionDockerCloudService(name string, action string) (dockercloud.Service, error) {
	var service dockercloud.Service
	var req *httplib.BeegoHTTPRequest
	uuid, err := getDockerCloudServiceUuid(name)
	if err != nil {
		return service, err
	}
	switch action {
	case "start":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/start")
	case "stop":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/stop")
	case "redeploy":
		req = httplib.Post(APIADDRESS + APIVERSION + DockerCloudToolEndpoint + "/service/" + uuid + "/redeploy")
	}
	req.SetBasicAuth(apiuser, apipassword)
	req.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	err = req.ToJSON(&service)
	if err != nil {
		return service, err
	}
	return service, nil
}

func getFavicons(domain string) string {
	//return "https://www.google.com/s2/favicons?domain=" + domain
	return "https://api.byi.pw/favicon/?url=" + domain
}
