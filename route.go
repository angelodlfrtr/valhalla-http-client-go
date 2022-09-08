package client

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"
)

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
	RouteInputCostingOptionsBicycleTypeRoad     string = "Road"
	RouteInputCostingOptionsBicycleTypeHybrid   string = "Hybrid"
	RouteInputCostingOptionsBicycleTypeCity     string = "City"
	RouteInputCostingOptionsBicycleTypeCross    string = "Cross"
	RouteInputCostingOptionsBicycleTypeMountain string = "Mountain"
)

const (
	DirectionsTypeNone         string = "none"
	DirectionsTypeManeuvers    string = "maneuvers"
	DirectionsTypeInstructions string = "instructions"
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

// RouteLocation location for route input.
// Used both in input and output
type RouteLocation struct {
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
	//
	// Auto set when output
	SideOfStreet *string `json:"side_of_street,omitempty"`

	// DateTime (response only) Expected date/time for the user to be at the location
	// using the ISO 8601 format (YYYY-MM-DDThh:mm) in the local time zone of departure or arrival.
	// For example "2015-12-29T08:00".
	DateTime *string `json:"date_time,omitempty"`

	// Output only fields

	// OriginalIndex returned in output
	OriginalIndex *int `json:"original_index,omitempty"`
}

type (
	RouteInputCostingOptionsBicycleBase struct {
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

		// CountryCrossingCost cost applied when encountering an international border.
		// This cost is added to the estimated and elapsed times. The default cost is 600 seconds.
		CountryCrossingCost *int `json:"country_crossing_cost,omitempty"`

		// CountryCrossingPenalty penalty applied for a country crossing.
		// This penalty can be used to create paths that avoid spanning country boundaries.
		// The default penalty is 0.
		CountryCrossingPenalty *int `json:"country_crossing_penalty,omitempty"`

		// ServicePenalty penalty applied for transition to generic service road.
		// The default penalty is 0 trucks and 15 for cars, buses, motor scooters and motorcycles.
		ServicePenalty *int `json:"service_penalty,omitempty"`

		// Shortest changes the metric to quasi-shortest, i.e. purely distance-based costing.
		// Note, this will disable all other costings & penalties.
		// Also note, shortest will not disable hierarchy pruning, leading to potentially
		// sub-optimal routes for some costing models. The default is false.
		Shortest *bool `json:"shortest,omitempty"`
	}

	RouteInputCostingOptionsMotorsBase struct {
		RouteInputCostingOptionsBicycleBase

		// PrivateAccessPenalty penalty applied when a gate or bollard with access=private
		// is encountered. The default private access penalty is 450 seconds.
		PrivateAccessPenalty *int `json:"private_access_penalty,omitempty"`

		// TollBoothCost cost applied when a toll booth is encountered.
		// This cost is added to the estimated and elapsed times. The default cost is 15 seconds.
		TollBoothCost *int `json:"toll_booth_cost,omitempty"`

		// TollBoothPenalty penalty applied to the cost when a toll booth is encountered.
		// This penalty can be used to create paths that avoid toll roads.
		// The default toll booth penalty is 0.
		TollBoothPenalty *int `json:"toll_booth_penalty,omitempty"`

		// FerryCost cost applied when entering a ferry. This cost is added to the
		// estimated and elapsed times. The default cost is 300 seconds (5 minutes).
		FerryCost *int `json:"ferry_cost,omitempty"`

		// UseFerry value indicates the willingness to take ferries.
		// This is a range of values between 0 and 1.
		// Values near 0 attempt to avoid ferries and values near 1 will favor ferries.
		// The default value is 0.5. Note that sometimes ferries are required to complete
		// a route so values of 0 are not guaranteed to avoid ferries entirely.
		UseFerry *float32 `json:"use_ferry,omitempty"`

		// UseHighways value indicates the willingness to take highways.
		// This is a range of values between 0 and 1. Values near 0 attempt to avoid
		// highways and values near 1 will favor highways. The default value is 1.0.
		// Note that sometimes highways are required to complete a route so values
		// of 0 are not guaranteed to avoid highways entirely.
		UseHighways *float32 `json:"use_highways,omitempty"`

		// UseTolls value indicates the willingness to take roads with tolls.
		// This is a range of values between 0 and 1.
		// Values near 0 attempt to avoid tolls and values near 1 will not attempt to avoid them.
		// The default value is 0.5. Note that sometimes roads with tolls are required
		// to complete a route so values of 0 are not guaranteed to avoid them entirely.
		UseTolls *float32 `json:"use_tolls,omitempty"`

		// UseLivingStreets value indicates the willingness to take living streets.
		// This is a range of values between 0 and 1.
		// Values near 0 attempt to avoid living streets and values near 1 will favor living streets.
		// The default value is 0 for trucks, 0.1 for cars, buses, motor scooters and motorcycles.
		// Note that sometimes living streets are required to complete a route so values
		// of 0 are not guaranteed to avoid living streets entirely.
		UseLivingStreets *float32 `json:"use_living_streets,omitempty"`

		// UseTracks value indicates the willingness to take track roads.
		// This is a range of values between 0 and 1. Values near 0 attempt to avoid tracks
		// and values near 1 will favor tracks a little bit. The default value is 0 for autos,
		// 0.5 for motor scooters and motorcycles. Note that sometimes tracks are required
		// to complete a route so values of 0 are not guaranteed to avoid tracks entirely.
		UseTracks *float32 `json:"use_tracks,omitempty"`

		// ServiceFactor factor that modifies (multiplies) the cost when generic service
		// roads are encountered. The default service_factor is 1.
		ServiceFactor *float32 `json:"service_factor,omitempty"`

		// TopSpeed top speed the vehicle can go. Also used to avoid roads with higher speeds
		// than this value. top_speed must be between 10 and 252 KPH. The default value is 140 KPH.
		TopSpeed *int `json:"top_speed,omitempty"`

		// IgnoreClosures if set to true, ignores all closures, marked due to live traffic closures,
		// during routing. Note: This option cannot be set if location.search_filter.exclude_closures
		// is also specified in the request and will return an error if it is.
		IgnoreClosures *bool `json:"ignore_closures,omitempty"`

		// ClosureFactor factor that penalizes the cost when traversing a closed edge
		// (eg: if search_filter.exclude_closures is false for origin and/or destination location
		// and the route starts/ends on closed edges). Its value can range from 1.0 - don't
		// penalize closed edges, to 10.0 - apply high cost penalty to closed edges.
		// Default value is 9.0. Note: This factor is applicable only for motorized modes
		// of transport, i.e auto, motorcycle, motor_scooter, bus, truck & taxi.
		ClosureFactor *float32 `json:"closure_factor,omitempty"`

		// Height of the vehicle (in meters). Default 1.9 for car, bus, taxi and 2.6 for truck.
		Height *float32 `json:"height,omitempty"`

		// Width of the vehicle (in meters). Default 1.6 for car, bus, taxi and 4.11 for truck.
		Width *float32 `json:"width,omitempty"`

		// ExcludeUnpaved value indicates whether or not the path may include unpaved roads.
		// If exclude_unpaved is set to 1 it is allowed to start and end with unpaved roads,
		// but is not allowed to have them in the middle of the route path,
		// otherwise they are allowed. Default false.
		ExcludeUnpaved *bool `json:"exclude_unpaved,omitempty"`

		// ExcludeCashOnlyTolls value which indicates the desire to avoid routes
		// with cash-only tolls. Default false.
		ExcludeCashOnlyTolls *bool `json:"exclude_cash_only_tolls,omitempty"`

		// IncludeHov2 value which indicates the desire to include HOV roads with
		// a 2-occupant requirement in the route when advantageous. Default false.
		IncludeHov2 *bool `json:"include_hov2,omitempty"`

		// IncludeHov3 value which indicates the desire to include HOV roads with a 3-occupant
		// requirement in the route when advantageous. Default false.
		IncludeHov3 *bool `json:"include_hov3,omitempty"`

		// IncludeHot value which indicates the desire to include tolled HOV
		// roads which require the driver to pay a toll if the occupant requirement
		// isn't met. Default false.
		IncludeHot *bool `json:"include_hot,omitempty"`
	}

	RouteInputCostingOptionsAuto RouteInputCostingOptionsMotorsBase
	RouteInputCostingOptionsTaxi RouteInputCostingOptionsMotorsBase
	RouteInputCostingOptionsBus  RouteInputCostingOptionsMotorsBase

	RouteInputCostingOptionsTruck struct {
		RouteInputCostingOptionsMotorsBase

		// Length of the truck (in meters). Default 21.64.
		Length *float32 `json:"length,omitempty"`

		// Weight of the truck (in metric tons). Default 21.77.
		Weight *float32 `json:"weight,omitempty"`

		// AxleLoad of the truck (in metric tons). Default 9.07.
		AxleLoad *float32 `json:"axle_load,omitempty"`

		// Hazmat value indicating if the truck is carrying hazardous materials. Default false.
		Hazmat *bool `json:"hazmat,omitempty"`
	}

	RouteInputCostingOptionsBicycle struct {
		RouteInputCostingOptionsBicycleBase

		// BicycleType type of bicycle. The default type is Hybrid.
		//
		// Road: a road-style bicycle with narrow tires that is generally
		// lightweight and designed for speed on paved surfaces.
		//
		// Hybrid or City: a bicycle made mostly for city riding or casual riding on roads
		// and paths with good surfaces.
		//
		// Cross: a cyclo-cross bicycle, which is similar to a road bicycle
		// but with wider tires suitable to rougher surfaces.
		//
		// Mountain: a mountain bicycle suitable for most surfaces but generally
		// heavier and slower on paved surfaces.
		BicycleType *string `json:"bicycle_type,omitempty"`

		// CyclingSpeed is the average travel speed along smooth, flat roads.
		// This is meant to be the speed a rider can comfortably maintain over the desired distance
		// of the route. It can be modified (in the costing method) by surface type in conjunction
		// with bicycle type and (coming soon) by hilliness of the road section.
		// When no speed is specifically provided, the default speed is determined by the bicycle
		// type and are as follows: Road = 25 KPH (15.5 MPH), Cross = 20 KPH (13 MPH),
		// Hybrid/City = 18 KPH (11.5 MPH), and Mountain = 16 KPH (10 MPH).
		CyclingSpeed *float32 `json:"cycling_speed,omitempty"`

		// UseRoads a cyclist's propensity to use roads alongside other vehicles.
		// This is a range of values from 0 to 1, where 0 attempts to avoid roads and stay
		// on cycleways and paths, and 1 indicates the rider is more comfortable riding on roads.
		// Based on the use_roads factor, roads with certain classifications and higher speeds
		// are penalized in an attempt to avoid them when finding the best path.
		// The default value is 0.5.
		UseRoads *float32 `json:"use_roads,omitempty"`

		// UseHills a cyclist's desire to tackle hills in their routes.
		// This is a range of values from 0 to 1, where 0 attempts to avoid hills and steep grades
		// even if it means a longer (time and distance) path, while 1 indicates the rider
		// does not fear hills and steeper grades. Based on the use_hills factor, penalties
		// are applied to roads based on elevation change and grade. These penalties help the
		// path avoid hilly roads in favor of flatter roads or less steep grades where available.
		// Note that it is not always possible to find alternate paths to avoid hills
		// (for example when route locations are in mountainous areas). The default value is 0.5.
		UseHills *float32 `json:"use_hills,omitempty"`

		// UseFerry value indicates the willingness to take ferries.
		// This is a range of values between 0 and 1. Values near 0 attempt to avoid ferries
		// and values near 1 will favor ferries. Note that sometimes ferries are required
		// to complete a route so values of 0 are not guaranteed to avoid ferries entirely.
		// The default value is 0.5.
		UseFerry *float32 `json:"use_ferry,omitempty"`

		// UseLivingStreets value indicates the willingness to take living streets.
		// This is a range of values between 0 and 1. Values near 0 attempt to avoid living
		// streets and values from 0.5 to 1 will currently have no effect on route selection.
		// The default value is 0.5. Note that sometimes living streets are required to complete
		// a route so values of 0 are not guaranteed to avoid living streets entirely.
		UseLivingStreets *float32 `json:"use_living_streets,omitempty"`

		// AvoidBadSurfaces value is meant to represent how much a cyclist wants to avoid roads
		// with poor surfaces relative to the bicycle type being used.
		// This is a range of values between 0 and 1. When the value is 0, there is no penalization
		// of roads with different surface types; only bicycle speed on each surface is taken
		// into account. As the value approaches 1, roads with poor surfaces for the bike are
		// penalized heavier so that they are only taken if they significantly improve travel time.
		// When the value is equal to 1, all bad surfaces are completely disallowed from routing,
		// including start and end points. The default value is 0.25.
		AvoidBadSurfaces *float32 `json:"avoid_bad_surfaces,omitempty"`

		// BssReturnCost value is useful when bikeshare is chosen as travel mode.
		// It is meant to give the time will be used to return a rental bike.
		// This value will be displayed in the final directions and used to calculate the whole
		// duation. The default value is 120 seconds.
		BssReturnCost *int `json:"bss_return_cost,omitempty"`

		// BssReturnPenalty value is useful when bikeshare is chosen as travel mode.
		// It is meant to describe the potential effort to return a rental bike.
		// This value won't be displayed and used only inside of the algorithm.
		BssReturnPenalty *int `json:"bss_return_penalty,omitempty"`
	}

	RouteInputCostingOptionsMotorScooter struct {
		RouteInputCostingOptionsMotorsBase

		// TopSpeed the motorized scooter can go. Used to avoid roads with higher speeds
		// than this value. For motor_scooter this value must be between 20 and 120 KPH.
		// The default value is 45 KPH (~28 MPH)
		TopSpeed *float32 `json:"top_speed,omitempty"`

		// UsePrimary a riders's propensity to use primary roads.
		// This is a range of values from 0 to 1, where 0 attempts to avoid primary roads,
		// and 1 indicates the rider is more comfortable riding on primary roads.
		// Based on the use_primary factor, roads with certain classifications and higher speeds
		// are penalized in an attempt to avoid them when finding the best path.
		// The default value is 0.5.
		UsePrimary *float32 `json:"use_primary,omitempty"`

		// UseHills a riders's desire to tackle hills in their routes.
		// This is a range of values from 0 to 1, where 0 attempts to avoid hills and steep
		// grades even if it means a longer (time and distance) path, while 1 indicates the rider
		// does not fear hills and steeper grades. Based on the use_hills factor, penalties
		// are applied to roads based on elevation change and grade. These penalties help the
		// path avoid hilly roads in favor of flatter roads or less steep grades where available.
		// Note that it is not always possible to find alternate paths to avoid hills
		// (for example when route locations are in mountainous areas).
		// The default value is 0.5.
		UseHills *float32 `json:"use_hills,omitempty"`

		// Shortest changes the metric to quasi-shortest, i.e. purely distance-based costing.
		// Note, this will disable all other costings & penalties. Also note, shortest will not
		// disable hierarchy pruning, leading to potentially sub-optimal routes for some costing
		// models. The default is false.
		Shortest *bool `json:"shortest,omitempty"`
	}

	RouteInputCostingOptionsMotorcycle struct {
		RouteInputCostingOptionsMotorsBase

		// UseHighways a riders's propensity to prefer the use of highways.
		// This is a range of values from 0 to 1, where 0 attempts to avoid highways,
		// and values toward 1 indicates the rider prefers highways.
		// The default value is 1.0.
		UseHighways *float32 `json:"use_highways,omitempty"`

		// UseTrails a riders's desire for adventure in their routes.
		// This is a range of values from 0 to 1, where 0 will avoid trails, tracks,
		// unclassified or bad surfaces and values towards 1 will tend to avoid major
		// roads and route on secondary roads. The default value is 0.0.
		UseTrails *float32 `json:"use_trails,omitempty"`

		// Shortest changes the metric to quasi-shortest, i.e. purely distance-based costing.
		// Note, this will disable all other costings & penalties. Also note, shortest will not
		// disable hierarchy pruning, leading to potentially sub-optimal routes for some costing
		// models. The default is false.
		Shortest *bool `json:"shortest,omitempty"`
	}

	RouteInputCostingOptionsPedestrian struct {
		// Walking speed in kilometers per hour. Must be between 0.5 and 25 km/hr. Defaults to 5.1 km/hr (3.1 miles/hour).
		WalkingSpeed *float32 `json:"walking_speed,omitempty"`

		// A factor that modifies the cost when encountering roads classified as footway
		// (no motorized vehicles allowed), which may be designated footpaths or designated
		// sidewalks along residential roads. Pedestrian routes generally attempt to favor using
		// these walkways and sidewalks. The default walkway_factor is 1.0.
		WalkwayFactor *float32 `json:"walkway_factor,omitempty"`

		// A factor that modifies the cost when encountering roads with dedicated sidewalks.
		// Pedestrian routes generally attempt to favor using sidewalks.
		// The default sidewalk_factor is 1.0.
		SidewalkFactor *float32 `json:"sidewalk_factor,omitempty"`

		// A factor that modifies (multiplies) the cost when alleys are encountered.
		// Pedestrian routes generally want to avoid alleys or narrow service roads between buildings.
		// The default alley_factor is 2.0.
		AlleyFactor *float32 `json:"alley_factor,omitempty"`

		// A factor that modifies (multiplies) the cost when encountering a driveway,
		// which is often a private, service road. Pedestrian routes generally want to avoid
		// driveways (private). The default driveway factor is 5.0.
		DrivewayFactor *float32 `json:"driveway_factor,omitempty"`

		// A penalty in seconds added to each transition onto a path with steps or stairs.
		// Higher values apply larger cost penalties to avoid paths that contain flights of steps.
		StepPenalty *int `json:"step_penalty,omitempty"`

		// This value indicates the willingness to take ferries.
		// This is range of values between 0 and 1. Values near 0 attempt to avoid ferries and
		// values near 1 will favor ferries. The default value is 0.5.
		// Note that sometimes ferries are required to complete a route so values of 0
		// are not guaranteed to avoid ferries entirely.
		UseFerry *float32 `json:"use_ferry,omitempty"`

		// This value indicates the willingness to take living streets.
		// This is a range of values between 0 and 1.
		// Values near 0 attempt to avoid living streets and values near 1 will favor living streets.
		// The default value is 0.6. Note that sometimes living streets are required to complete
		// a route so values of 0 are not guaranteed to avoid living streets entirely.
		UseLivingStreets *float32 `json:"use_living_streets,omitempty"`

		// This value indicates the willingness to take track roads.
		// This is a range of values between 0 and 1. Values near 0 attempt to avoid tracks
		// and values near 1 will favor tracks a little bit.
		// The default value is 0.5. Note that sometimes tracks are required to complete
		// a route so values of 0 are not guaranteed to avoid tracks entirely.
		UseTracks *float32 `json:"use_tracks,omitempty"`

		// This is a range of values from 0 to 1, where 0 attempts to avoid hills and steep grades
		// even if it means a longer (time and distance) path, while 1 indicates the pedestrian
		// does not fear hills and steeper grades. Based on the use_hills factor,
		// penalties are applied to roads based on elevation change and grade.
		// These penalties help the path avoid hilly roads in favor of flatter roads
		// or less steep grades where available. Note that it is not always possible to find
		// alternate paths to avoid hills (for example when route locations are in
		// mountainous areas). The default value is 0.5.
		UseHills *float32 `json:"use_hills,omitempty"`

		// A penalty applied for transition to generic service road. The default penalty is 0.
		ServicePenalty *int `json:"service_penalty,omitempty"`

		// A factor that modifies (multiplies) the cost when generic service roads are encountered.
		// The default service_factor is 1.
		ServiceFactor *int `json:"service_factor,omitempty"`

		// This value indicates the maximum difficulty of hiking trails that is allowed.
		// Values between 0 and 6 are allowed. The values correspond to sac_scale values
		// within OpenStreetMap, see reference here. The default value is 1 which means that
		// well cleared trails that are mostly flat or slightly sloped are allowed.
		// Higher difficulty trails can be allowed by specifying a higher
		// value for max_hiking_difficulty.
		MaxHikingDifficulty *int `json:"max_hiking_difficulty,omitempty"`

		// This value is useful when bikeshare is chosen as travel mode.
		// It is meant to give the time will be used to rent a bike from a bike share station.
		// This value will be displayed in the final directions and used to calculate
		// the whole duation. The default value is 120 seconds.
		BssRentCost *int `json:"bss_rent_cost,omitempty"`

		// This value is useful when bikeshare is chosen as travel mode.
		// It is meant to describe the potential effort to rent a bike from a bike share station.
		// This value won't be displayed and used only inside of the algorithm.
		BssRentPenalty *int `json:"bss_rent_penalty,omitempty"`

		// Changes the metric to quasi-shortest, i.e. purely distance-based costing.
		// Note, this will disable all other costings & penalties.
		// Also note, shortest will not disable hierarchy pruning, leading to potentially
		// sub-optimal routes for some costing models. The default is false.
		Shortest *bool `json:"shortest,omitempty"`
	}

	RouteInputCostingOptionsTransitFilter struct {
		Ids    []string `json:"ids,omitempty"`
		Action *string  `json:"action,omitempty"`
	}

	RouteInputCostingOptionsTransit struct {
		// UseBus user's desire to use buses.
		// Range of values from 0 (try to avoid buses) to 1 (strong preference for riding buses).
		UseBus *float32 `json:"use_bus,omitempty"`

		// UseRail User's desire to use rail/subway/metro.
		// Range of values from 0 (try to avoid rail) to 1 (strong preference for riding rail).
		UseRail *float32 `json:"use_rail,omitempty"`

		// UseTransfers user's desire to favor transfers.
		// Range of values from 0 (try to avoid transfers) to 1 (totally comfortable with transfers).
		UseTransfers *float32 `json:"use_transfers,omitempty"`

		// TransitStartEndMaxDistance a pedestrian option that can be added to the request
		// to extend the defaults (2145 meters or approximately 1.5 miles).
		// This is the maximum walking distance at the beginning or end of a route.
		TransitStartEndMaxDistance *int `json:"transit_start_end_max_distance,omitempty"`

		// TransitTransferMaxDistance a pedestrian option that can be added to the request
		// to extend the defaults (800 meters or 0.5 miles).
		// This is the maximum walking distance between transfers.
		TransitTransferMaxDistance *int `json:"transit_transfer_max_distance,omitempty"`

		// Filters a way to filter for one or more stops, routes, or operators.
		// Filters must contain a list of Onestop IDs, which is a unique identifier
		// for Transitland data, and an action.
		//
		// ids: any number of Onestop IDs (such as o-9q9-bart)
		//
		// action: either exclude to exclude all of the ids listed in the filter
		// or include to include only the ids listed in the filter
		Filters map[string]*RouteInputCostingOptionsTransitFilter `json:"filters,omitempty"`
	}
)

type RouteInputCostingOptions struct {
	Auto         *RouteInputCostingOptionsAuto         `json:"auto,omitempty"`
	Taxi         *RouteInputCostingOptionsTaxi         `json:"taxi,omitempty"`
	Bus          *RouteInputCostingOptionsBus          `json:"bus,omitempty"`
	Truck        *RouteInputCostingOptionsTruck        `json:"truck,omitempty"`
	Bicycle      *RouteInputCostingOptionsBicycle      `json:"bicycle,omitempty"`
	MotorScooter *RouteInputCostingOptionsMotorScooter `json:"motor_scooter,omitempty"`
	Motorcycle   *RouteInputCostingOptionsMotorcycle   `json:"motorcycle,omitempty"`
	Pedestrian   *RouteInputCostingOptionsPedestrian   `json:"pedestrian,omitempty"`

	Transit *RouteInputCostingOptionsTransit `json:"transit,omitempty"`
}

type RouteInputDateTime struct {
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
	Locations []*RouteLocation `json:"locations,omitempty"`

	// Costing Valhalla's routing service uses dynamic, run-time costing to generate the route path.
	// The route request must include the name of the costing model and can include optional
	// parameters available for the chosen costing model.
	Costing *string `json:"costing,omitempty"`

	// CostingOptions (optional) Costing options for the specified costing model.
	CostingOptions *RouteInputCostingOptions `json:"costing_options,omitempty"`

	// Units distance units for output.
	// Allowable unit types are miles (or mi) and kilometers (or km).
	// If no unit type is specified, the units default to kilometers.
	Units *string `json:"units,omitempty"`

	// Language of the narration instructions based on the IETF BCP 47 language tag string.
	// If no language is specified or the specified language is unsupported,
	// United States-based English (en-US) is used.
	// See https://valhalla.readthedocs.io/en/latest/api/turn-by-turn/api-reference/#supported-language-tags
	Language *string `json:"language,omitempty"`

	// DirectionsType an enum with 3 values.
	//
	// none: indicating no maneuvers or instructions should be returned.
	//
	// maneuvers: indicating that only maneuvers be returned.
	//
	// instructions: indicating that maneuvers with instructions
	// should be returned (this is the default if not specified).
	DirectionsType *string `json:"directions_type,omitempty"`

	// Alternates a number denoting how many alternate routes should be provided.
	// There may be no alternates or less alternates than the user specifies.
	// Alternates are not yet supported on multipoint routes
	// (that is, routes with more than 2 locations).
	// They are also not supported on time dependent routes.
	Alternates *int `json:"alternates,omitempty"`

	// ExcludeLocations a set of locations to exclude or avoid within a route can be specified
	// using a JSON array of avoid_locations.
	// The avoid_locations have the same format as the locations list.
	// At a minimum each avoid location must include latitude and longitude.
	// The avoid_locations are mapped to the closest road or roads and these
	// roads are excluded from the route path computation.
	ExcludeLocations []*RouteLocation `json:"exclude_locations,omitempty"`

	// ExcludePolygons one or multiple exterior rings of polygons in the form of nested JSON arrays,
	// e.g. [[[lon1, lat1], [lon2,lat2]],[[lon1,lat1],[lon2,lat2]]].
	// Roads intersecting these rings will be avoided during path finding.
	// If you only need to avoid a few specific roads, it's much more efficient to use
	// exclude_locations. Valhalla will close open rings
	// (i.e. copy the first coordingate to the last position).
	ExcludePolygons [][][]float64 `json:"exclude_polygons,omitempty"`

	// DateTime this is the local date and time at the location.
	DateTime *RouteInputDateTime `json:"date_time,omitempty"`

	// OutFormat if no out_format is specified, JSON is returned.
	// Future work includes PBF (protocol buffer) support.
	OutFormat *string `json:"out_format,omitempty"`

	// ID name your route request. If id is specified, the naming will be sent thru to the response.
	ID *string `json:"id,omitempty"`

	// LinearReferences when present and true, the successful route response will include a key
	// linear_references. Its value is an array of base64-encoded OpenLR location references,
	// one for each graph edge of the road network matched by the input trace.
	LinearReferences *bool `json:"linear_references,omitempty"`

	// PrioritizeBidirectional prioritize bidirectional a* when date_time.type = depart_at/current.
	// By default time_dependent_forward a* is used in these cases,
	// but bidirectional a* is much faster. Currently it does not update the time (and speeds)
	// when searching for the route path, but the ETA on that route is recalculated
	// based on the time-dependent speeds
	PrioritizeBidirectional *bool `json:"prioritize_bidirectional,omitempty"`
}

type RouteOutputTripSummary struct {
	// Time Estimated elapsed time to complete the trip.
	Time *float64 `json:"time,omitempty"`

	// HasTimeRestrictions indicates if the trip has time restrictions.
	HasTimeRestrictions *bool `json:"has_time_restrictions,omitempty"`

	// MinLat minimum latitude of a bounding box containing the route.
	MinLat *float64 `json:"min_lat,omitempty"`

	// MinLon minimum longitude of a bounding box containing the route.
	MinLon *float64 `json:"min_lon,omitempty"`

	// MaxLat maximum latitude of a bounding box containing the route.
	MaxLat *float64 `json:"max_lat,omitempty"`

	// MaxLon maximum longitude of a bounding box containing the route.
	MaxLon *float64 `json:"max_lon,omitempty"`

	// Length length of the route in meters.
	Length *float64 `json:"length,omitempty"`

	// Cost cost of the route in the units specified by costing ?
	Cost *float64 `json:"cost,omitempty"`
}

type RouteOutputManeuverSign struct {
	// Text interchange sign text.
	Text *string `json:"text,omitempty"`

	// ConsecutiveCount the frequency of this sign element within a set a consecutive signs.
	// This item is optional.
	// TODO: int instead of float32 ?
	ConsecutiveCount *float32 `json:"consecutive_count,omitempty"`
}

type RouteOutputManeuverTransitInfoTransitStop struct {
	// Type of stop (simple stop=0; station=1).
	Type *int `json:"type,omitempty"`

	// Global transit stop identifier from Transitland.
	OnestopId *string `json:"onestop_id,omitempty"`

	// Name of the stop or station. For example "14 St - Union Sq".
	Name *string `json:"name,omitempty"`

	// Arrival date and time using the ISO 8601 format (YYYY-MM-DDThh:mm). For example, "2015-12-29T08:06".
	ArrivalDateTime *string `json:"arrival_date_time,omitempty"`

	// Departure date and time using the ISO 8601 format (YYYY-MM-DDThh:mm). For example, "2015-12-29T08:06".
	DepartureDateTime *string `json:"departure_date_time,omitempty"`

	// True if this stop is a marked as a parent stop.
	IsParentStop *bool `json:"is_parent_stop,omitempty"`

	// True if the times are based on an assumed schedule because the actual schedule is not known.
	AssumedSchedule *bool `json:"assumed_schedule,omitempty"`

	// Latitude of the transit stop in degrees.
	Lat *float64 `json:"lat,omitempty"`

	// Longitude of the transit stop in degrees.
	Lon *float64 `json:"lon,omitempty"`
}

type RouteOutputManeuverTransitInfo struct {
	// Global transit route identifier from Transitland.
	OnestopId *string `json:"OnestopId,omitempty"`

	// Short name describing the transit route. For example "N".
	ShortName *string `json:"short_name,omitempty"`

	// Long name describing the transit route. For example "Broadway Express".
	LongName *string `json:"long_name,omitempty"`

	// The sign on a public transport vehicle that identifies the route destination to passengers.
	// For example "ASTORIA - DITMARS BLVD".
	Headsign *string `json:"headsign,omitempty"`

	// The numeric color value associated with a transit route.
	// The value for yellow would be "16567306".
	Color *string `json:"color,omitempty"`

	// The numeric text color value associated with a transit route. The value for black would be "0".
	TextColor *string `json:"text_color,omitempty"`

	// The description of the the transit route. For example "Trains operate from Ditmars Boulevard,
	// Queens, to Stillwell Avenue, Brooklyn, at all times. N trains in Manhattan operate along
	// Broadway and across the Manhattan Bridge to and from Brooklyn.
	// Trains in Brooklyn operate along 4th Avenue, then through Borough Park to Gravesend.
	// Trains typically operate local in Queens, and either express or local in Manhattan
	// and Brooklyn, depending on the time. Late night trains operate via Whitehall Street,
	// Manhattan. Late night service is local".
	Description *string `json:"description,omitempty"`

	// Global operator/agency identifier from Transitland.
	OperatorOnestopId *string `json:"operator_onestop_id,omitempty"`

	// Operator/agency name. For example, "BART", "King County Marine Division", and so on. Short name is used over long name.
	OperatorName *string `json:"operator_name,omitempty"`

	// Operator/agency URL. For example, "http://web.mta.info/".
	OperatorUrl *string `json:"operator_url,omitempty"`

	// A list of the stops/stations associated with a specific transit route. See below for details.
	TransitStops []*RouteOutputManeuverTransitInfoTransitStop `json:"transit_stops,omitempty"`
}

type RouteOutputManeuver struct {
	// Type of maneuver. See doc for a list.
	// https://valhalla.readthedocs.io/en/latest/api/turn-by-turn/api-reference/
	Type *int `json:"type,omitempty"`

	// Instruction written maneuver instruction.
	// Describes the maneuver, such as "Turn right onto Main Street".
	Instruction *string `json:"instruction,omitempty"`

	// VerbalTransitionAlertInstruction text suitable for use as a verbal alert in a navigation
	// application. The transition alert instruction will prepare the user for the forthcoming
	// transition. For example: "Turn right onto North Prince Street".
	VerbalTransitionAlertInstruction *string `json:"verbal_transition_alert_instruction,omitempty"`

	// VerbalSuccinctTransitionInstruction TODO ? no doc
	VerbalSuccinctTransitionInstruction *string `json:"verbal_succinct_transition_instruction,omitempty"`

	// VerbalPreTransitionInstruction text suitable for use as a verbal message immediately
	// prior to the maneuver transition. For example "Turn right onto North Prince Street, U.S. 2 22".
	VerbalPreTransitionInstruction *string `json:"verbal_pre_transition_instruction,omitempty"`

	// VerbalPostTransitionInstruction text suitable for use as a verbal message immediately
	// after the maneuver transition. For example "Continue on U.S. 2 22 for 3.9 miles".
	VerbalPostTransitionInstruction *string `json:"verbal_post_transition_instruction,omitempty"`

	// StreetNames list of street names that are consistent along the entire nonobvious maneuver.
	StreetNames []string `json:"street_names,omitempty"`

	// BeginStreetNames when present, these are the street names at the beginning
	// (transition point) of the nonobvious maneuver (if they are different than the names
	// that are consistent along the entire nonobvious maneuver).
	BeginStreetNames []string `json:"begin_street_names,omitempty"`

	// Time estimated time along the maneuver in seconds.
	Time *float64 `json:"time,omitempty"`

	// Length maneuver length in the units specified.
	Length *float64 `json:"length,omitempty"`

	// Cost TODO
	Cost *float64 `json:"cost,omitempty"`

	// BeginShapeIndex index into the list of shape points for the start of the maneuver.
	BeginShapeIndex *int `json:"begin_shape_index,omitempty"`

	// EndShapeIndex index into the list of shape points for the end of the maneuver.
	EndShapeIndex *int `json:"end_shape_index,omitempty"`

	// Toll True if the maneuver has any toll, or portions of the maneuver are subject to a toll.
	Toll *bool `json:"toll,omitempty"`

	// Rough true if the maneuver is unpaved or rough pavement,
	// or has any portions that have rough pavement.
	Rough *bool `json:"rough,omitempty"`

	// Gate true if a gate is encountered on this maneuver.
	Gate *bool `json:"gate,omitempty"`

	// Ferry true if a ferry is encountered on this maneuver.
	Ferry *bool `json:"ferry,omitempty"`

	// contains the interchange guide information at a road junction associated with this maneuver.
	Sign map[string][]*RouteOutputManeuverSign `json:"sign,omitempty"`

	// RoundaboutExitCount the spoke to exit roundabout after entering.
	RoundaboutExitCount *int `json:"roundabout_exit_count,omitempty"`

	// DepartInstruction written depart time instruction.
	// Typically used with a transit maneuver, such as "Depart: 8:04 AM from 8 St - NYU".
	DepartInstruction *string `json:"depart_instruction,omitempty"`

	// VerbalDepartInstruction text suitable for use as a verbal depart time instruction.
	// Typically used with a transit maneuver, such as "Depart at 8:04 AM from 8 St - NYU".
	VerbalDepartInstruction *string `json:"verbal_depart_instruction,omitempty"`

	// ArriveInstruction written arrive time instruction.
	// Typically used with a transit maneuver, such as "Arrive: 8:10 AM at 34 St - Herald Sq".
	ArriveInstruction *string `json:"arrive_instruction,omitempty"`

	// VerbalArriveInstruction text suitable for use as a verbal arrive time instruction.
	// Typically used with a transit maneuver, such as "Arrive at 8:10 AM at 34 St - Herald Sq".
	VerbalArriveInstruction *string `json:"verbal_arrive_instruction,omitempty"`

	// TODO
	TransitInfo interface{} `json:"transit_info,omitempty"`

	// VerbalMultiCue true if the verbal_pre_transition_instruction has been appended
	// with the verbal instruction of the next maneuver.
	VerbalMultiCue *bool `json:"verbal_multi_cue,omitempty"`

	// TravelMode travel mode of the maneuver: drive, pedestrian, bicycle, transit
	TravelMode *string `json:"travel_mode,omitempty"`

	// TravelType travel type
	// Possibles: car, foot, road, tram, metro, rail, bus, ferry, cable_car, gondola, funicular
	TravelType *string `json:"travel_type,omitempty"`
}

type RouteOutputLeg struct {
	// Summary summary of the loeg.
	Summary *RouteOutputTripSummary `json:"summary,omitempty"`

	// Maneuvers a list of maneuvers.
	Maneuvers []*RouteOutputManeuver `json:"maneuvers,omitempty"`
}

type RouteOutputTrip struct {
	// Locations the locations used to generate the route.
	Locations []*RouteLocation `json:"locations,omitempty"`

	// Legs the legs of the route.
	Legs []*RouteOutputLeg `json:"legs,omitempty"`

	// Summary summary of the trip.
	Summary *RouteOutputTripSummary `json:"summary,omitempty"`

	// Shape an encoded polyline of the route path (with 6 digits decimal precision).
	Shape *string `json:"shape,omitempty"`
}

type RouteOutput struct {
	// ID from the id in request
	ID *string `json:"id,omitempty"`

	// Trip response
	Trip *RouteOutputTrip `json:"trip,omitempty"`
}

// Route returns the route between the given locations.
func (client *Client) Route(input *RouteInput) (*RouteOutput, error) {
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

	// Extract response
	output := &RouteOutput{}
	if err := json.Unmarshal(resp.Body(), output); err != nil {
		return nil, fmt.Errorf("error while decoding http route json response data: %w", err)
	}

	return output, nil
}
