// Use a shared workqueue ('batch') to do some processing against the entire database.
package fgae

import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"

	fdb "github.com/skypies/flightdb2"
)

// {{{ BatchHandler

// This enqueues tasks for each key in the DB.
func BatchHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	//c := appengine.Timeout(appengine.NewContext(r), 5*time.Minute)
	// db := FlightDB{C:c}

	str := "Kicking off the batch run\n"

	keys,err := datastore.NewQuery(kFlightKind).KeysOnly().GetAll(c, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	str += fmt.Sprintf("Hello - we found %d keys\n", len(keys))

	for _,k := range keys {
		//str += " k = "+k.Encode() + "\n"
	
		t := taskqueue.NewPOSTTask("/fdb/batch/instance", map[string][]string{
			"k": {k.Encode()},
		})
		if _,err := taskqueue.Add(c, t, "batch"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("OK\n%s", str)))
}

// }}}
// {{{ BatchInstanceHandler

// /fdb/batch/instance?k=agxzfnNlcmZyMC1mZGJyDgsSBmZsaWdodBiK8QQM

// This handler re-keys all the flight objects.
func BatchInstanceHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	str := fmt.Sprintf("OK\nbatch, for [%s]\n", r.FormValue("k"))
	
	key,err := datastore.DecodeKey(r.FormValue("k"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blob := fdb.IndexedFlightBlob{}
	if err := datastore.Get(c, key, &blob); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	f,err := blob.ToFlight(key.Encode())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str += fmt.Sprintf("* pulled up %s\n", f)

	icaoId := string(f.IcaoId)
	if icaoId == "" {
		str += "* No IcaoID in flight, bailing"
	} else {
		rootKey := datastore.NewKey(c, kFlightKind, string(f.IcaoId), 0, nil)
		newKey := datastore.NewIncompleteKey(c, kFlightKind, rootKey)
		str += fmt.Sprintf("** old: %#v\n**root: %#v\n** new: %#v\n", key, rootKey, newKey)

		if _,err := datastore.Put(c, newKey, &blob); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		str += "\n* added under new key!\n"
		
		if err := datastore.Delete(c, key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		str += "\n* deleted under old key\n"
	}

	db := FlightDB{C:c}
	db.Infof("%s", str)
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(str))
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}
