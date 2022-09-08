package client

import (
	"testing"

	"github.com/gotidy/ptr"
)

func TestRoute(t *testing.T) {
	input := &RouteInput{
		Locations: []*RouteLocation{},
		Costing:   ptr.String("auto"),
		Units:     ptr.String("km"),
	}

	input.Locations = append(input.Locations, &RouteLocation{Lat: ptr.Float64(48.390394), Lon: ptr.Float64(-4.486076)})
	input.Locations = append(input.Locations, &RouteLocation{Lat: ptr.Float64(48.45252), Lon: ptr.Float64(-4.25252)})

	clt := getTestClient()

	output, err := clt.Route(input)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}
