package ui

import(
	"fmt"
	"net/http"

	"github.com/skypies/flightdb/fgae"
)

func init() {
}

// ?idspec=XX,YY,...  (or ?idspec=XX&idspec=YYY&...)
//  &viewtype={vector,descent,sideview,track}
//  &sample=N        (sample the track every N seconds)

//  &alt=30000       (max altitude for graph)
//  &length=80       (max distance from origin; in nautical miles)
//  &dist=from       (for distance axis, use dist from airport; by default, uses dist along path)
//  &colorby=delta   (delta groundspeed, instead of groundspeed)

func VisualizeHandler(db fgae.FlightDB, w http.ResponseWriter, r *http.Request) {
	if r.FormValue("debug") != "" {
		str := "OK\n"
		for k, v := range r.Form {
			str += fmt.Sprintf(" %-20.20s: '%s'\n", k, v)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(str))
	}
	
	switch r.FormValue("viewtype") {
	case "vector":   TracksetHandler(db,w,r)
	case "sideview": SideviewHandler(db,w,r)
	case "track":    TrackHandler(db,w,r)
	default:         http.Error(w, "Specify viewtype={vector|sideview|track}", http.StatusBadRequest)
	}		
}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
