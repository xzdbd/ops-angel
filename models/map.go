package models

import (
	"net/url"
	"time"
)

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
