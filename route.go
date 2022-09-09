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
	CostingOptions *CostingModelOptions `json:"costing_options,omitempty"`

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

	// Shape an encoded polyline of the route path (with 6 digits decimal precision).
	Shape *string `json:"shape,omitempty"`
}

type RouteOutputTrip struct {
	// Locations the locations used to generate the route.
	Locations []*RouteLocation `json:"locations,omitempty"`

	// Legs the legs of the route.
	Legs []*RouteOutputLeg `json:"legs,omitempty"`

	// Summary summary of the trip.
	Summary *RouteOutputTripSummary `json:"summary,omitempty"`
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
