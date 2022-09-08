package client

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/paulmach/go.geojson"
	"github.com/valyala/fasthttp"
)

type IsochroneInputLocation struct {
	Lat *float64 `json:"lat,omitempty"`
	Lon *float64 `json:"lon,omitempty"`
}

type IsochroneInputDateTime struct {
	// Type of date time.
	//
	// 0 - Current departure time.
	//
	// 1 - Specified departure time
	//
	// 2 - Specified arrival time. Not yet implemented for multimodal costing method.
	//
	// 3 - Invariant specified time. Time does not vary over the course of the path.
	// Not implemented for multimodal or bike share routing
	Type *int `json:"type,omitempty"`

	// Value the date and time is specified in ISO 8601 format (YYYY-MM-DDThh:mm)
	// in the local time zone of departure or arrival.
	// For example "2016-07-03T08:06"
	Value *string `json:"value,omitempty"`
}

type IsochroneInputContour struct {
	// Time a floating point value specifying the time in minutes for the contour.
	Time *float64 `json:"time,omitempty"`

	// Distance a floating point value specifying the distance in kilometers for the contour.
	Distance *float64 `json:"distance,omitempty"`

	// Color or the output of the contour. Specify it as a Hex value, but without the #,
	// such as "color":"ff0000" for red. If no color is specified,
	// the isochrone service will assign a default color to the output.
	Color *string `json:"color,omitempty"`
}

type IsochroneInput struct {
	// Locations must include a latitude and longitude in decimal degrees.
	// The coordinates can come from many input sources, such as a GPS location,
	// a point or a click on a map, a geocoding service, and so on.
	// External search services, such as Mapbox Geocoding can be used to find places
	// and geocode addresses, whose coordinates can be used as input to the service.
	Locations []*IsochroneInputLocation `json:"locations,omitempty"`

	// Costing the isochrone service uses the auto, bicycle, pedestrian, and multimodal costing
	// models available in the Valhalla Turn-by-Turn service.
	Costing *string `json:"costing,omitempty"`

	// CostingOptions (optional) Costing options for the specified costing model.
	CostingOptions *CostingModelOptions `json:"costing_options,omitempty"`

	// DateTime 	The local date and time at the location.
	// These parameters apply only for multimodal requests and are not used with other
	// costing methods.
	DateTime *IsochroneInputDateTime `json:"date_time,omitempty"`

	// ID name of the isochrone request.
	// If id is specified, the name is returned with the response.
	ID *string `json:"id,omitempty"`

	// Contours to use for each isochrone contour.
	Contours []*IsochroneInputContour `json:"contours,omitempty"`

	// Polygons determine whether to return geojson polygons or linestrings as the contours.
	// The default is false, which returns lines; when true, polygons are returned.
	// Note: When polygons is true, any contour that forms a ring is returned as a polygon.
	Polygons *bool `json:"polygons,omitempty"`

	// Denoise a floating point value from 0 to 1 (default of 1) which can be used to remove
	// smaller contours. A value of 1 will only return the largest contour for a given time value.
	// A value of 0.5 drops any contours that are less than half the area of the largest contour
	// in the set of contours for that same time value.
	Denoise *float64 `json:"denoise,omitempty"`

	// Generalize a floating point value in meters used as the tolerance for Douglas-Peucker
	// generalization. Note: Generalization of contours can lead to self-intersections,
	// as well as intersections of adjacent contours.
	Generalize *float64 `json:"generalize,omitempty"`

	// ShowLocations a boolean indicating whether the input locations should be returned
	// as MultiPoint features: one feature for the exact input coordinates and one feature
	// for the coordinates of the network node it snapped to. Default false.
	ShowLocations *bool `json:"show_locations,omitempty"`
}

// Isochrone returns the isochrone for the specified locations.
func (client *Client) Isochrone(input *IsochroneInput) (*geojson.FeatureCollection, error) {
	req, err := client.buildBaseRequest(fasthttp.MethodPost, "/route", input)
	if err != nil {
		return nil, fmt.Errorf("failed to build request for route: %w", err)
	}
	defer fasthttp.ReleaseRequest(req)

	// Acquire response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.httpClient.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error while calling http route service: %w", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		errRes := &ErrorResponse{}
		if err := json.Unmarshal(resp.Body(), errRes); err != nil {
			errRes.StatusCode = resp.StatusCode()
			errRes.ErrorMessage = string(resp.Body())
		}

		return nil, errRes
	}

	// Build geojson feature collection
	fc := geojson.NewFeatureCollection()
	if err := json.Unmarshal(resp.Body(), fc); err != nil {
		return nil, fmt.Errorf("error while decoding http isochrone json response data: %w", err)
	}

	return fc, nil
}
