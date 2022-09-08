package client

import (
	"testing"

	"github.com/gotidy/ptr"
)

func TestElevation(t *testing.T) {
	input := &ElevationInput{
		HeightPrecision: ptr.Int(2),
		Shape:           []*ElevationPoint{},
	}

	input.Shape = append(input.Shape, &ElevationPoint{Lat: 42.913581, Lon: 0.137267})
	input.Shape = append(input.Shape, &ElevationPoint{Lat: 42.913612, Lon: 0.137234})

	clt := getTestClient()

	output, err := clt.Elevation(input)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}
