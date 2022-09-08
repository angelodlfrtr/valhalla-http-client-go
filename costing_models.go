package client

const (
	// CostingModelAuto standard costing for driving routes by car, motorcycle,
	// truck, and so on that obeys automobile driving rules, such as access and turn restrictions.
	// Auto provides a short time path (though not guaranteed to be shortest time)
	// and uses intersection costing to minimize turns and maneuvers or road name changes.
	// Routes also tend to favor highways and higher classification roads,
	// such as motorways and trunks.
	CostingModelAuto string = "auto"

	// CostingModelBicycle standard costing for travel by bicycle, with a slight preference
	// for using cycleways or roads with bicycle lanes. Bicycle routes follow regular roads
	// when needed, but avoid roads without bicycle access.
	CostingModelBicycle string = "bicycle"

	// CostingModelBus standard costing for bus routes. Bus costing inherits the auto
	// costing behaviors, but checks for bus access on the roads.
	CostingModelBus string = "bus"

	// CostingModelBikeshare BETA a combination of pedestrian and bicycle.
	// Use bike share station(amenity:bicycle_rental) to change the travel mode.
	CostingModelBikeshare string = "bikeshare"

	// CostingModelTruck standard costing for trucks. Truck costing inherits the auto costing
	// behaviors, but checks for truck access, width and height restrictions,
	// and weight limits on the roads.
	CostingModelTruck string = "truck"

	// CostingModelTaxi standard costing for taxi routes. Taxi costing inherits the
	// auto costing behaviors, but checks for taxi lane access on the roads and favors those roads.
	CostingModelTaxi string = "taxi"

	// CostingModelMotorScooter BETA standard costing for travel by motor scooter or moped.
	// By default, motor_scooter costing will avoid higher class roads unless the country
	// overrides allows motor scooters on these roads. Motor scooter routes follow regular
	// roads when needed, but avoid roads without motor_scooter, moped, or mofa access.
	CostingModelMotorScooter string = "motor_scooter"

	// CostingModelMultimodal Currently supports pedestrian and transit.
	// In the future, multimodal will support a combination of all of the above.
	CostingModelMultimodal string = "multimodal"

	// CostingModelPedestrian standard walking route that excludes roads without
	// pedestrian access. In general, pedestrian routes are shortest distance with the
	// following exceptions: walkways and footpaths are slightly favored,
	// while steps or stairs and alleys are slightly avoided.
	CostingModelPedestrian string = "pedestrian"
)

type (
	CostingModelOptionsBicycleBase struct {
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

	CostingModelOptionsMotorsBase struct {
		CostingModelOptionsBicycleBase

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

	CostingModelOptionsAuto CostingModelOptionsMotorsBase
	CostingModelOptionsTaxi CostingModelOptionsMotorsBase
	CostingModelOptionsBus  CostingModelOptionsMotorsBase

	CostingModelOptionsTruck struct {
		CostingModelOptionsMotorsBase

		// Length of the truck (in meters). Default 21.64.
		Length *float32 `json:"length,omitempty"`

		// Weight of the truck (in metric tons). Default 21.77.
		Weight *float32 `json:"weight,omitempty"`

		// AxleLoad of the truck (in metric tons). Default 9.07.
		AxleLoad *float32 `json:"axle_load,omitempty"`

		// Hazmat value indicating if the truck is carrying hazardous materials. Default false.
		Hazmat *bool `json:"hazmat,omitempty"`
	}

	CostingModelOptionsBicycle struct {
		CostingModelOptionsBicycleBase

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

	CostingModelOptionsMotorScooter struct {
		CostingModelOptionsMotorsBase

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

	CostingModelOptionsMotorcycle struct {
		CostingModelOptionsMotorsBase

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

	CostingModelOptionsPedestrian struct {
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

	CostingModelOptionsTransitFilter struct {
		Ids    []string `json:"ids,omitempty"`
		Action *string  `json:"action,omitempty"`
	}

	CostingModelOptionsTransit struct {
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
		Filters map[string]*CostingModelOptionsTransitFilter `json:"filters,omitempty"`
	}
)

type CostingModelOptions struct {
	Auto         *CostingModelOptionsAuto         `json:"auto,omitempty"`
	Taxi         *CostingModelOptionsTaxi         `json:"taxi,omitempty"`
	Bus          *CostingModelOptionsBus          `json:"bus,omitempty"`
	Truck        *CostingModelOptionsTruck        `json:"truck,omitempty"`
	Bicycle      *CostingModelOptionsBicycle      `json:"bicycle,omitempty"`
	MotorScooter *CostingModelOptionsMotorScooter `json:"motor_scooter,omitempty"`
	Motorcycle   *CostingModelOptionsMotorcycle   `json:"motorcycle,omitempty"`
	Pedestrian   *CostingModelOptionsPedestrian   `json:"pedestrian,omitempty"`

	Transit *CostingModelOptionsTransit `json:"transit,omitempty"`
}
