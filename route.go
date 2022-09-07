package client

const (
	RouteInputLocationTypeBreak        string = "break"
	RouteInputLocationTypeThrough      string = "through"
	RouteInputLocationTypeVia          string = "via"
	RouteInputLocationTypeBreakThrough string = "break_through"
)

const (
	RouteInputLocationPreferredSideSame     string = "same"
	RouteInputLocationPreferredSideOpposite string = "opposite"
	RouteInputLocationPreferredSideEither   string = "either"
)

const (
	// RouteInputCostingAuto standard costing for driving routes by car, motorcycle,
	// truck, and so on that obeys automobile driving rules, such as access and turn restrictions.
	// Auto provides a short time path (though not guaranteed to be shortest time)
	// and uses intersection costing to minimize turns and maneuvers or road name changes.
	// Routes also tend to favor highways and higher classification roads,
	// such as motorways and trunks.
	RouteInputCostingAuto string = "auto"

	// RouteInputCostingBicycle standard costing for travel by bicycle, with a slight preference
	// for using cycleways or roads with bicycle lanes. Bicycle routes follow regular roads
	// when needed, but avoid roads without bicycle access.
	RouteInputCostingBicycle string = "bicycle"

	// RouteInputCostingBus standard costing for bus routes. Bus costing inherits the auto
	// costing behaviors, but checks for bus access on the roads.
	RouteInputCostingBus string = "bus"

	// RouteInputCostingBikeshare BETA a combination of pedestrian and bicycle.
	// Use bike share station(amenity:bicycle_rental) to change the travel mode.
	RouteInputCostingBikeshare string = "bikeshare"

	// RouteInputCostingTruck standard costing for trucks. Truck costing inherits the auto costing
	// behaviors, but checks for truck access, width and height restrictions,
	// and weight limits on the roads.
	RouteInputCostingTruck string = "truck"

	// RouteInputCostingTaxi standard costing for taxi routes. Taxi costing inherits the
	// auto costing behaviors, but checks for taxi lane access on the roads and favors those roads.
	RouteInputCostingTaxi string = "taxi"

	// RouteInputCostingMotorScooter BETA standard costing for travel by motor scooter or moped.
	// By default, motor_scooter costing will avoid higher class roads unless the country
	// overrides allows motor scooters on these roads. Motor scooter routes follow regular
	// roads when needed, but avoid roads without motor_scooter, moped, or mofa access.
	RouteInputCostingMotorScooter string = "motor_scooter"

	// RouteInputCostingMultimodal Currently supports pedestrian and transit.
	// In the future, multimodal will support a combination of all of the above.
	RouteInputCostingMultimodal string = "multimodal"

	// RouteInputCostingPedestrian standard walking route that excludes roads without
	// pedestrian access. In general, pedestrian routes are shortest distance with the
	// following exceptions: walkways and footpaths are slightly favored,
	// while steps or stairs and alleys are slightly avoided.
	RouteInputCostingPedestrian string = "pedestrian"
)

// RouteInputLocationSearchFilter search filter for route input location.
type RouteInputLocationSearchFilter struct {
	// ExcludeTunnel whether to exclude roads marked as tunnels.
	ExcludeTunnel *bool `json:"exclude_tunnel,omitempty"`

	// ExcludeBridge whether to exclude roads marked as bridges.
	ExcludeBridge *bool `json:"exclude_bridge,omitempty"`

	// ExcludeRamp whether to exclude link roads marked as ramps,
	// note that some turn channels are also marked as ramps.
	ExcludeRamp *bool `json:"exclude_ramp,omitempty"`

	// ExcludeClosures whether to exclude roads considered closed due to live traffic closure.
	// Note: This option cannot be set if costing_options.<costing>.ignore_closures is also
	// specified. An error is returned if both options are specified.
	// Note 2: Ignoring closures at destination and source locations does NOT work
	// for date_time type 0/1 & 2 respectively.
	ExcludeClosures *bool `json:"exclude_closures,omitempty"`

	// MinRoadClass (defaults to "service_other"): lowest road class allowed.
	MinRoadClass *string `json:"min_road_class,omitempty"`

	// MaxRoadClass (defaults to "motorway"): highest road class allowed
	MaxRoadClass *string `json:"max_road_class,omitempty"`
}

// RouteInputLocation location for route input.
type RouteInputLocation struct {
	// Lon longitude of the location in degrees.
	// This is assumed to be both the routing location and the display location if
	// no display_lat and display_lon are provided.
	Lon *float64 `json:"lon"`

	// Lat latitude of the location in degrees.
	// This is assumed to be both the routing location and the display location
	// if no display_lat and display_lon are provided.
	Lat *float64 `json:"lat"`

	// Type of location, either break, through, via or break_through. Each type controls
	// two characteristics: whether or not to allow a u-turn at the location and whether
	// or not to generate guidance/legs at the location. A break is a location at which we
	// allows u-turns and generate legs and arrival/departure maneuvers.
	// A through location is a location at which we neither allow u-turns nor generate legs
	// or arrival/departure maneuvers. A via location is a location at which we allow u-turns
	// but do not generate legs or arrival/departure maneuvers.
	// A break_through location is a location at which we do not allow u-turns but do generate
	// legs and arrival/departure maneuvers. If no type is provided, the type is assumed
	// to be a break. The types of the first and last locations are ignored and are treated as breaks.
	Type *string `json:"type,omitempty"`

	// Heading (optional) preferred direction of travel for the start from the location.
	// This can be useful for mobile routing where a vehicle is traveling in a specific direction
	// along a road, and the route should start in that direction.
	// The heading is indicated in degrees from north in a clockwise direction,
	// where north is 0°, east is 90°, south is 180°, and west is 270°.
	Heading *float32 `json:"heading,omitempty"`

	// HeadingTolerance (optional) How close in degrees a given street's angle must be in order
	// for it to be considered as in the same direction of the heading parameter.
	// The default value is 60 degrees.
	HeadingTolerance *float32 `json:"heading_tolerance,omitempty"`

	// Street (optional) name. The street name may be used to assist finding the correct routing
	// location at the specified latitude, longitude. This is not currently implemented.
	Street *string `json:"street,omitempty"`

	// WayID (optional) OpenStreetMap identification number for a polyline way.
	// The way ID may be used to assist finding the correct routing location at the specified
	// latitude, longitude. This is not currently implemented.
	WayID *string `json:"way_id,omitempty"`

	// MinimumReachability minimum number of nodes (intersections) reachable for a given edge
	// (road between intersections) to consider that edge as belonging to a connected region.
	// When correlating this location to the route network, try to find candidates who
	// are reachable from this many or more nodes (intersections).
	// If a given candidate edge reaches less than this number of nodes its considered to be
	// a disconnected island and we'll search for more candidates until we find at least one that
	// isn't considered a disconnected island. If this value is larger than the configured service
	// limit it will be clamped to that limit. The default is a minimum of 50 reachable nodes.
	MinimumReachability int `json:"minimum_reachability,omitempty"`

	// Radius The number of meters about this input location within which edges
	// (roads between intersections) will be considered as candidates for said location.
	// When correlating this location to the route network, try to only return results within
	// this distance (meters) from this location. If there are no candidates within this distance
	// it will return the closest candidate within reason. If this value is larger than
	// the configured service limit it will be clamped to that limit. The default is 0 meters.
	Radius int `json:"radius,omitempty"`

	// RankCandidates whether or not to rank the edge candidates for this location.
	// The ranking is used as a penalty within the routing algorithm so that some edges
	// will be penalized more heavily than others. If true candidates will be ranked according
	// to their distance from the input and various other attributes.
	// If false the candidates will all be treated as equal which should lead to routes that
	// are just the most optimal path with emphasis about which edges were selected.
	RankCandidates *bool `json:"rank_candidates,omitempty"`

	// PreferredSide If the location is not offset from the road centerline or is closest
	// to an intersection this option has no effect. Otherwise the determined side of street
	// is used to determine whether or not the location should be visited from the same,
	// opposite or either side of the road with respect to the side of the road the given
	// locale drives on. In Germany (driving on the right side of the road),
	// passing a value of same will only allow you to leave from or arrive at a location such
	// that the location will be on your right. In Australia (driving on the left side of the road),
	// passing a value of same will force the location to be on your left.
	// A value of opposite will enforce arriving/departing from a location on the opposite side
	// of the road from that which you would be driving on while a value of either will make
	// no attempt limit the side of street that is available for the route.
	PreferredSide *string `json:"preferred_side,omitempty"`

	// DisplayLat latitude of the map location in degrees. If provided the lat and lon parameters
	// will be treated as the routing location and the display_lat and display_lon will
	// be used to determine the side of street. Both display_lat and display_lon must be
	// provided and valid to achieve the desired effect.
	DisplayLat *float64 `json:"display_lat,omitempty"`

	// DisplayLon longitude of the map location in degrees. If provided the lat and lon parameters
	// will be treated as the routing location and the display_lat and display_lon will
	// be used to determine the side of street. Both display_lat and display_lon must be
	// provided and valid to achieve the desired effect.
	DisplayLon *float64 `json:"display_lon,omitempty"`

	// @TODO: which type is it ?
	// SearchCutoff the cutoff at which we will assume the input is too far away
	// from civilisation to be worth correlating to the nearest graph elements.
	SearchCutoff *string `json:"search_cutoff,omitempty"`

	// NodeSnapTolerance during edge correlation this is the tolerance used to determine
	// whether or not to snap to the intersection rather than along the street,
	// if the snap location is within this distance from the intersection the intersection
	// is used instead. The default is 5 meters.
	NodeSnapTolerance *float64 `json:"node_snap_tolerance,omitempty"`

	// StreetSideTolerance if your input coordinate is less than this tolerance away
	// from the edge centerline then we set your side of street to none otherwise your side
	// of street will be left or right depending on direction of travel.
	StreetSideTolerance *float64 `json:"street_side_tolerance,omitempty"`

	// StreetSideMaxDistance the max distance in meters that the input coordinates or
	// display ll can be from the edge centerline for them to be used for determining the
	// side of street. Beyond this distance the side of street is set to none.
	StreetSideMaxDistance *float64 `json:"street_side_max_distance,omitempty"`

	// SearchFilter a set of optional filters to exclude candidate edges based on their attribution.
	SearchFilter *RouteInputLocationSearchFilter `json:"search_filter,omitempty"`

	// Next parameters has no effect pn routing and are returned as a convenience.

	// Name Location or business name.
	// The name may be used in the route narration directions,
	// such as "You have arrived at <business name>.")
	Name *string `json:"name,omitempty"`

	// City name.
	City *string `json:"city,omitempty"`

	// State name.
	State *string `json:"state,omitempty"`

	// PostalCode postal code.
	PostalCode *string `json:"postal_code,omitempty"`

	// Country name.
	Country *string `json:"country,omitempty"`

	// Phone phone number.
	Phone *string `json:"phone,omitempty"`

	// URL URL.
	URL *string `json:"url,omitempty"`

	// SideOfStreet (response only) The side of street of a break location that is
	// determined based on the actual route when the location is offset from the street.
	// The possible values are left and right.
	SideOfStreet *string `json:"side_of_street,omitempty"`

	// DateTime (response only) Expected date/time for the user to be at the location
	// using the ISO 8601 format (YYYY-MM-DDThh:mm) in the local time zone of departure or arrival.
	// For example "2015-12-29T08:00".
	DateTime *string `json:"date_time,omitempty"`
}

type (
	RouteInputCostingOptionsAuto struct {
		// ManeuverPenalty penalty applied when transitioning between roads that do not have
		// consistent naming–in other words, no road names in common.
		// This penalty can be used to create simpler routes that tend to have fewer maneuvers
		// or narrative guidance instructions. The default maneuver penalty is five seconds.
		ManeuverPenalty *int `json:"maneuver_penalty,omitempty"`

		// GateCost cost applied when a gate with undefined or private access is encountered.
		// This cost is added to the estimated time / elapsed time.
		// The default gate cost is 30 seconds.
		GateCost *int `json:"gate_cost,omitempty"`

		// GatePenalty penalty applied when a gate with no access information is on the road.
		// The default gate penalty is 300 seconds.
		GatePenalty *int `json:"gate_penalty,omitempty"`

		// TODO: continue
		// https://valhalla.readthedocs.io/en/latest/api/turn-by-turn/api-reference/
	}
	RouteInputCostingOptionsBus   RouteInputCostingOptionsAuto
	RouteInputCostingOptionsTruck RouteInputCostingOptionsAuto
)

type RouteInputCostingOptions struct {
	Auto  *RouteInputCostingOptionsAuto  `json:"auto,omitempty"`
	Bus   *RouteInputCostingOptionsBus   `json:"bus,omitempty"`
	Truck *RouteInputCostingOptionsTruck `json:"truck,omitempty"`
	// Bicycle      *RouteInputCostingOptionsBicycle      `json:"bicycle,omitempty"`
	// MotorScooter *RouteInputCostingOptionsMotorScooter `json:"motor_scooter,omitempty"`
	// Motorcycle   *RouteInputCostingOptionsMotorcycle   `json:"motorcycle,omitempty"`
	// Pedestrian   *RouteInputCostingOptionsPedestrian   `json:"pedestrian,omitempty"`
}

// RouteInput is the input for turn by turn routing service
type RouteInput struct {
	// Locations specify locations as an ordered list of two or more locations within a JSON array.
	// Locations are visited in the order specified.
	// A location must include a latitude and longitude in decimal degrees.
	// The coordinates can come from many input sources, such as a GPS location,
	// a point or a click on a map, a geocoding service, and so on.
	// Note that the Valhalla cannot search for names or addresses or perform geocoding or
	// reverse geocoding. External search services, such as Mapbox Geocoding, can be used
	// to find places and geocode addresses, which must be converted to coordinates for input.
	// To build a route, you need to specify two break locations. In addition,
	// you can include through, via or break_through locations to influence the route path.
	Locations []*RouteInputLocation `json:"locations,omitempty"`

	// Costing Valhalla's routing service uses dynamic, run-time costing to generate the route path.
	// The route request must include the name of the costing model and can include optional
	// parameters available for the chosen costing model.
	Costing *string `json:"costing,omitempty"`

	// CostingOptions (optional) Costing options for the specified costing model.
	CostingOptions *RouteInputCostingOptions `json:"costing_options,omitempty"`
}
