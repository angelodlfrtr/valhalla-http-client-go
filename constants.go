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
