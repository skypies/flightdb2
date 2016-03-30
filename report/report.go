package report

import(
	"fmt"
	"html/template"
	"sort"

	"github.com/skypies/util/histogram"
	fdb "github.com/skypies/flightdb2"
)

type FlightReportOutcome int
const(
	RejectedByGeoRestriction FlightReportOutcome = iota
	RejectedByReport
	Accepted
	Undefined
)
type ReportFunc func(*Report, *fdb.Flight, []fdb.TrackIntersection)(FlightReportOutcome,error)
type SummarizeFunc func(*Report)

type ReportLogLevel int
const(
	DEBUG = iota
	INFO
)

type Report struct {
	Name              string
	ReportingContext  // embedded
	Options           // embedded
	Func              ReportFunc
	SummarizeFunc     // embedded, but just to avoid a more confusing name
	TrackSpec       []string

	// Private state a report might accumulate (be careful about RAM though!)
	Blobs map[string]interface{}
	
	// Output state
	RowsHTML  [][]template.HTML
	RowsText  [][]string

	HeadersText []string
	
	I         map[string]int
	F         map[string]float64
	S         map[string]string
	H         histogram.Histogram
	
	Log string
}

func BlankReport() Report {
	return Report{
		I: map[string]int{},
		F: map[string]float64{},
		S: map[string]string{},
		RowsHTML: [][]template.HTML{},
		RowsText: [][]string{},
		HeadersText: []string{},
		Blobs: map[string]interface{}{},
	}
}

func (r *Report)Logger(level ReportLogLevel, s string) {
	if level < r.Options.ReportLogLevel { return }
	r.Log += s
}
func (r *Report)Info(s string) { r.Logger(INFO, s) }
func (r *Report)Debug(s string) { r.Logger(DEBUG, s) }

func (r *Report)SetHeaders(headers []string) {
	if len(r.HeadersText) == 0 { r.HeadersText = headers }
}
func (r *Report)AddRow(html *[]string, text *[]string) {
	htmlRow := []template.HTML{}
	for _,s  := range *html { htmlRow = append(htmlRow, template.HTML(s)) }
	if html != nil { r.RowsHTML = append(r.RowsHTML, htmlRow) }
	if text != nil { r.RowsText = append(r.RowsText, *text) }
}

func (r *Report)ListPreferredDataSources() []string {
	// Dumb logic for now ...
	if r.Options.TrackDataSource != "" {
		return []string{r.Options.TrackDataSource}
	}
	return r.TrackSpec
}

// Ensure the flight matches all the search restrictions
func (r *Report)PreProcess(f *fdb.Flight) (bool, []fdb.TrackIntersection) {
	r.I["[A] PreProcessed"]++
	
	// If restrictions were specified, only match flights that satisfy them
	failed := false
	intersections := []fdb.TrackIntersection{}
	for _,gr := range r.Options.ListGeoRestrictors() {
		r.Debug(fmt.Sprintf("---- %s\nSources: %v\n", f.IdentityString(), r.ListPreferredDataSources()))
		satisfies,intersection,deb := f.SatisfiesGeoRestriction(gr, r.ListPreferredDataSources())
		_=deb
/*
		r.Debug(fmt.Sprintf("---- %s --[%v]--\n", f.IdentityString(), r.TrackSpec))
		for _,tName := range f.ListTracks() {
			r.Debug(fmt.Sprintf(" [%-6.6s] %s\n", tName, f.Tracks[tName]))
		}
		r.Debug(fmt.Sprintf("\n%s\n", deb))
*/
		if satisfies {
			intersections = append(intersections, intersection)
		} else {
			r.I["[B] Eliminated: did not satisfy "+gr.String()]++
			failed = true
			break
		}
	}
	if failed { return false, intersections }

	r.I["[B] <b>Satisfied geo restrictions</b> "]++

	dataSrc := "non-ADSB"
	if f.HasTrack("ADSB") { dataSrc = "ADSB" }
	r.I["[Bz] track source: "+dataSrc]++
	
	return true, intersections
}

func (r *Report)Process(f *fdb.Flight) (FlightReportOutcome, error) {
	wasOK,intersections := r.PreProcess(f)
	if !wasOK { return RejectedByGeoRestriction,nil }
	return r.Func(r, f, intersections)
}

func (r *Report)FinishSummary() {
	r.Info("**** Stage: finish summary\n")
	r.Debug("* (DEBUG)\n")
	if r.SummarizeFunc != nil { r.SummarizeFunc(r) }
}

func (r *Report)MetadataTable()[][]template.HTML {
	all := map[string]string{}

	for k,v := range r.I { all[k] = fmt.Sprintf("%d", v) }
	for k,v := range r.F { all[k] = fmt.Sprintf("%.1f", v) }
	for k,v := range r.S { all[k] = v }

	if stats,valid := r.H.Stats(); valid {
		all["[Z] stats, <b>N</b>"] = fmt.Sprintf("%d", stats.N)
		all["[Z] stats, Mean"] = fmt.Sprintf("%.0f", stats.Mean)
		all["[Z] stats, Stddev"] = fmt.Sprintf("%.0f", stats.Stddev)
	}
	
	keys := []string{}
	for k,_ := range all { keys = append(keys, k) }
	sort.Strings(keys)
	
	out := [][]template.HTML{}
	for _,k := range keys {
		out = append(out, []template.HTML{ template.HTML(k), template.HTML(all[k]) })
	}
	
	return out
}