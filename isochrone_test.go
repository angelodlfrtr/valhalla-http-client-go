package client

import (
	"testing"

	"github.com/gotidy/ptr"
)

func TestIsochrone(t *testing.T) {
	input := &IsochroneInput{
		Costing: ptr.String(CostingModelPedestrian),
	}

	input.Locations = append(input.Locations, &IsochroneInputLocation{Lat: ptr.Float64(42.913581), Lon: ptr.Float64(0.137267)})
	input.Contours = append(input.Contours, &IsochroneInputContour{Time: ptr.Float64(10)})

	clt := getTestClient()

	output, err := clt.Isochrone(input)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}
