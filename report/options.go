package report

// All reports share this same options struct. Some options apply to all reports, some
// are interpreted creatively by others, and some only apply to one kind of report.
// They are all parsed of the http.request, including the report name.

import(
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	
	"github.com/skypies/geo"
	"github.com/skypies/geo/sfo"
	"github.com/skypies/util/widget"

	fdb "github.com/skypies/flightdb2"
)

type Options struct {
	Name               string
	Start, End         time.Time
	Tags             []string

	Procedures       []string  // Obsolete
	Waypoints        []string  // ARGH
	HackWaypoints    []string  // ARGH
	
	// Geo restriction 1: Box.
	Center             geo.Latlong
	RadiusKM           float64  // For defining areas of interest around waypoints
	SideKM             float64  // For defining areas of interest around waypoints

	// Geo-restriction 2: Window. Represent this a bit better.
	WindowTo, WindowFrom geo.Latlong	
	WindowMin, WindowMax float64

	// Data specification
	CanSeeFOIA         bool    // This is locked down to a few users. Upgrade to full ACL model?
	TrackDataSource    string  // TODO: replace with a cleverer track-specification thing
	
	// Options applicable to various reports
	ReferencePoint     geo.Latlong // Some reports do things in relation to a fixed point
	RefDistanceKM      float64     // ... and maybe within {dist} of {refpoint}
	AltitudeTolerance  float64  // Some reports care about this
	time.Duration      // embedded; a time tolerance
	
	// Formatting / output options
	ResultsFormat      string // csv, html
	MapLineOpacity     float64

	ReportLogLevel // For debugging
}

// If fields absent or blank, returns {0.0, 0.0}
// Expects two string fields that start with stem, with 'lat' or 'long' appended.
func FormValueLatlong(r *http.Request, stem string) geo.Latlong {
	lat  := widget.FormValueFloat64EatErrs(r, stem+"lat")
	long := widget.FormValueFloat64EatErrs(r, stem+"long")
	return geo.Latlong{lat,long}
}

func FormValueReportOptions(r *http.Request) (Options, error) {
	if r.FormValue("rep") == "" {
		return Options{}, fmt.Errorf("url arg 'rep' missing (no report specified")
	}

	s,e,err := widget.FormValueDateRange(r)	
	if err != nil { return Options{}, err }

//	cutoff,_ := time.Parse("2006.01.02", "2015.08.01")
//	if s.Before(cutoff) || e.Before(cutoff) {
//		return Options{}, fmt.Errorf("earliest date for reports is '%s'", cutoff)
//	}
	
	opt := Options{
		Name: r.FormValue("rep"),
		Start: s,
		End: e,
		Tags: widget.FormValueCommaSepStrings(r,"tags"),
		Procedures: widget.FormValueCommaSepStrings(r,"procedures"),
		Waypoints: []string{},

		Center: FormValueLatlong(r, "center"),
		RadiusKM: widget.FormValueFloat64EatErrs(r, "radiuskm"),
		SideKM: widget.FormValueFloat64EatErrs(r, "sidekm"),
		
		// Hmm.
		WindowTo: FormValueLatlong(r, "winto"),
		WindowFrom: FormValueLatlong(r, "winfrom"),
		WindowMin: widget.FormValueFloat64EatErrs(r, "winmin"),
		WindowMax: widget.FormValueFloat64EatErrs(r, "winmax"),

		AltitudeTolerance: widget.FormValueFloat64EatErrs(r, "altitudetolerance"),
		Duration: widget.FormValueDuration(r, "duration"),
		ReferencePoint: FormValueLatlong(r, "refpt"),
		RefDistanceKM: widget.FormValueFloat64EatErrs(r, "refdistancekm"),

		MapLineOpacity: widget.FormValueFloat64EatErrs(r, "maplineopacity"),
		ResultsFormat: r.FormValue("resultformat"),
	}
	
	switch r.FormValue("datasource") {
	case "ADSB": opt.TrackDataSource = "ADSB"
	case "fr24": opt.TrackDataSource = "fr24"
		// default means let the report pick
	}

	switch r.FormValue("log") {
	case "debug": opt.ReportLogLevel = DEBUG
	default:      opt.ReportLogLevel = INFO
	}

	for _,name := range []string{"waypoint1", "waypoint2", "waypoint3"} {
		if r.FormValue(name) != "" {
			waypoint := strings.ToUpper(r.FormValue(name))
			opt.HackWaypoints = append(opt.HackWaypoints, waypoint)
		}
	}
	
	// The form only sends one value (via dropdown), but keep it as a list just in case
	for _, waypoint := range widget.FormValueCommaSepStrings(r,"waypoint") {
		waypoint = strings.ToUpper(waypoint)
		if _,exists := sfo.KFixes[waypoint]; !exists {
			return Options{}, fmt.Errorf("Waypoint '%s' not known", waypoint)
		}
		opt.Waypoints = append(opt.Waypoints, waypoint)
	}

	if r.FormValue("refpt_waypoint") != "" {
		waypoint := strings.ToUpper(r.FormValue("refpt_waypoint"))
		if _,exists := sfo.KFixes[waypoint]; !exists {
			return Options{}, fmt.Errorf("Waypoint '%s' not known (ref pt)", waypoint)
		}
		opt.ReferencePoint = sfo.KFixes[waypoint]
	}
	
	if opt.Center.Lat>0.1 && opt.RadiusKM<=0.1 && opt.SideKM<=0.1 {
		return Options{}, fmt.Errorf("Must define a radius or boxside for the region")
	}
	
	return opt, nil
}
 
func (o Options)ListGeoRestrictors() []geo.GeoRestrictor {
	ret := []geo.GeoRestrictor{}
	if !o.WindowFrom.IsNil() && !o.WindowTo.IsNil() {
		window := geo.Window{
			LatlongLine: o.WindowFrom.BuildLine(o.WindowTo),
			MinAltitude: o.WindowMin,
			MaxAltitude: o.WindowMax,
		}
		ret = append(ret, window)
	}

	addCenteredRestriction := func(pos geo.Latlong) {
		if o.RadiusKM > 0 {
			ret = append(ret, pos.Circle(o.RadiusKM))
		} else if o.SideKM > 0 {
			ret = append(ret, pos.Box(o.SideKM,o.SideKM))
		}
	}

	if math.Abs(o.Center.Lat)>0.1 { addCenteredRestriction(o.Center) }
	for _,waypoint := range o.Waypoints { addCenteredRestriction(sfo.KFixes[waypoint]) }	

	return ret
}

func (o Options)GetRefpointRestrictor() geo.GeoRestrictor {
	return o.ReferencePoint.Box(o.RefDistanceKM, o.RefDistanceKM)
}

func (o Options)ListMapRenderers() []geo.MapRenderer {
	ret := []geo.MapRenderer{}

	for _,gr  := range o.ListGeoRestrictors() {
		ret = append(ret, gr)
	}

	if !o.ReferencePoint.IsNil() {
		dist := o.RefDistanceKM
		if dist == 0.0 { dist = 2 }
		ret = append(ret, o.ReferencePoint.Box(dist,dist))
	}

	for _,wp := range o.HackWaypoints {
		ret = append(ret, sfo.KFixes[wp].Box(fdb.KWaypointSnapKM,fdb.KWaypointSnapKM))
	}
	
	return ret
}

func (o Options)String() string {
	str := fmt.Sprintf("%#v\n--\n", o)
	str += "GeoRestrictors:-\n"
	for _,gr := range o.ListGeoRestrictors() {
		str += fmt.Sprintf("  %s\n", gr)
	}
	str += "--\n"
	
	return str
}

func LatlongToCGIArgs(stem string, pos geo.Latlong) string {
	return fmt.Sprintf("%slat=%.5f&%slong=%.5f", stem, pos.Lat, stem, pos.Long)
}

// A bare minimum of args, to embed in track links, so tracks can render with report tooltips
// and maps can see the geometry used
func (r *Report)ToCGIArgs() string {
	str := fmt.Sprintf("rep=%s&%s", r.Options.Name, widget.DateRangeToCGIArgs(r.Start, r.End))
	if r.TrackDataSource != "" { str += "&datasource="+r.TrackDataSource }

	if len(r.Waypoints) > 0 { str += "&waypoint=" + strings.Join(r.Options.Waypoints,",") }

	if !r.Center.IsNil() { str += "&" + LatlongToCGIArgs("center", r.Center) }
	if !r.ReferencePoint.IsNil() { str += "&" + LatlongToCGIArgs("refpt", r.ReferencePoint) }

	if r.RefDistanceKM > 0.0 { str += fmt.Sprintf("&refdistancekm=%.2f", r.RefDistanceKM) }
	if r.RadiusKM > 0.0 { str += fmt.Sprintf("&radiuskm=%.2f", r.RadiusKM) }
	if r.SideKM > 0.0 { str += fmt.Sprintf("&sidekm=%.2f", r.SideKM) }

	if !r.WindowFrom.IsNil() {
		str += "&" + LatlongToCGIArgs("winfrom", r.WindowFrom)
		str += "&" + LatlongToCGIArgs("winto", r.WindowTo)
		if r.WindowMin > 0.0 { str += fmt.Sprintf("&winmin=%.0f", r.WindowMin) }
		if r.WindowMax > 0.0 { str += fmt.Sprintf("&winmax=%.0f", r.WindowMax) }
	}
	
	if r.MapLineOpacity > 0.0 { str += fmt.Sprintf("&maplineopacity=%.2f", r.MapLineOpacity) }

	for i,wp := range r.HackWaypoints {
		str += fmt.Sprintf("&waypoint%d=%s", i+1, wp)
	}

	if len(r.Tags) > 0 { str += fmt.Sprintf("&tags=%s", strings.Join(r.Tags,",")) }
	
	return str
}
