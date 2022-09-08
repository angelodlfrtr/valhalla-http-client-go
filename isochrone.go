package client

// @TODO: WIP

type IsochroneInputLocation struct {
	Lat *float64 `json:"lat,omitempty"`
	Lon *float64 `json:"lon,omitempty"`
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
}
