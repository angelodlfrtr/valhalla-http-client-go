package client

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

// Elevationinput is the input for elevation service
type ElevationInput struct {
	// Range if true, both the height and cumulative distance are returned for each point.
	Range *bool `json:"range,omitempty"`

	// ResampleDistance parameter is a numeric value specifying the distance at which
	// the input polyline is sampled in order to provide uniform distances between
	// samples along the polyline.
	ResampleDistance *int `json:"resample_distance,omitempty"`

	// HeightPrecision specifies the precision (number of decimal places) of all
	// returned height values. Values of 0, 1, or 2 are admissable.
	// Defaults to 0 (integer precision).
	// Any other value will result in integer precision (0 decimal places).
	HeightPrecision *int `json:"height_precision,omitempty"`

	// Shape must include a latitude and longitude in decimal degrees,
	// and the locations are visited in the order specified.
	// The input coordinates can come from many input sources, such as a GPS location,
	// a point or a click on a map, a geocoding service, and so on.
	Shape []*Point `json:"shape,omitempty"`

	// ShapeFormat specifies whether the polyline is encoded with 6 digit precision (polyline6)
	// or 5 digit precision (polyline5).
	// If shape_format is not specified, the encoded polyline is expected to be 6 digit precision.
	ShapeFormat *string `json:"shape_format,omitempty"`

	// EncodedPolyline is a string of a polyline-encoded, with the specified precision,
	// shape and has the following parameters. Details on polyline encoding and decoding can
	// be found here: https://valhalla.readthedocs.io/en/latest/decoding/
	EncodedPolyline *string `json:"encoded_polyline,omitempty"`

	// ID name your elevation request. If id is specified,
	// the naming will be sent thru to the response.
	ID *string `json:"id,omitempty"`
}

// ElevationOutput is the output for elevation service
type ElevationOutput struct {
	// Shape contain the specified shape coordinates from the input request.
	Shape []*Point `json:"shape,omitempty"`

	// EncodedPolyline contain the specified encoded polyline,
	// with six degrees of precision, coordinates from the input request.
	EncodedPolyline *string `json:"encoded_polyline,omitempty"`

	// RangeHeight contain the 2D array of range (x) and height (y) per input latitude,
	// longitude coordinate.
	RangeHeight [][]float32 `json:"range_height,omitempty"`

	// Height contain an array of height for the associated latitude, longitude coordinates.
	Height []float32 `json:"height,omitempty"`
}

// Elevation returns the elevation for the given input
func (client *Client) Elevation(input *ElevationInput) (*ElevationOutput, error) {
	req, err := client.buildBaseRequest(fasthttp.MethodPost, "/height", input)
	if err != nil {
		return nil, fmt.Errorf("failed to build request for elevation: %w", err)
	}
	defer fasthttp.ReleaseRequest(req)

	// Acquire response
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.httpClient.Do(req, resp); err != nil {
		return nil, fmt.Errorf("error while calling http elevation service: %w", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		errRes := &ErrorResponse{}
		if err := json.Unmarshal(resp.Body(), errRes); err != nil {
			errRes.StatusCode = resp.StatusCode()
			errRes.ErrorMessage = string(resp.Body())
		}

		return nil, errRes
	}

	// Extract response
	output := &ElevationOutput{}
	if err := json.Unmarshal(resp.Body(), output); err != nil {
		return nil, fmt.Errorf("error while decoding http elevation json response data: %w", err)
	}

	return output, nil
}
