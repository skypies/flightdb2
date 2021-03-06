package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/bigquery" // different API; see ffs below.

	// "google.golang.org/ appengine/taskqueue"

	"github.com/skypies/util/date"
	"github.com/skypies/util/gcp/gcs"
	"github.com/skypies/util/gcp/tasks"
	"github.com/skypies/util/widget"

	"github.com/skypies/flightdb/fgae"
)

var(
	// The bigquery dataset (dest) is an entirely different google cloud project.

	// This, the 'src' project, needs its service worker account to be
	// added as an 'editor' to the 'dest' project, so that we can submit
	// bigquery load requests.

	// Similarly, the service worker from the 'dest' project needs to be added to
	// the 'source' project, so that dest can read the GCS folders. I think.

	// This is in the 'src' project
	folderGCS = "bigquery-flights"

	// This is the 'dest' project
	bigqueryProject = "serfr0-1000"
	bigqueryDataset = "public"
	bigqueryTableName = "flights"
)

// {{{ publishAllFlightsHandler

// http://backend-dot-serfr0-fdb.appspot.com/batch/publish-all-flights?skipload=1&date=range&range_from=2016/07/01&range_to=2016/07/03

// /batch/publish-all-flights?date=range&range_from=2015/08/09&range_to=2015/08/10
//  ?skipload=1  (skip loading them into bigquery; it's better to bulk load all of them at once)

// Writes them all into a batch queue
func publishAllFlightsHandler(db fgae.FlightDB, w http.ResponseWriter, r *http.Request) {
	ctx := db.Ctx()
	str := ""

	s,e,_ := widget.FormValueDateRange(r)
	days := date.IntermediateMidnights(s.Add(-1 * time.Second),e) // decrement start, to include it
	taskUrl := "/batch/publish-flights"

	taskClient,err := tasks.GetClient(ctx)
	if err != nil {
		db.Errorf(" publishAllFlightsHandler: GetClient: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _,day := range days {
		dayStr := day.Format("2006.01.02")
		thisUrl := fmt.Sprintf("%s?datestring=%s", taskUrl, dayStr)
		if r.FormValue("skipload") != "" {
			thisUrl += "&skipload=" + r.FormValue("skipload")
		}
		params := url.Values{}

		if _,err := tasks.SubmitAETask(ctx, taskClient, ProjectID, LocationID, QueueName, 0, thisUrl, params); err != nil {
			db.Errorf(" publishAllFlightsHandler: enqueue: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
/*
		t := taskqueue.NewPOSTTask(thisUrl, map[string][]string{})

		if _,err := taskqueue.Add(ctx, t, "bigbatch"); err != nil {
			db.Errorf("publishAllFlightsHandler: enqueue: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
*/
		str += " * posting for " + thisUrl + "\n"
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("OK, enqueued %d\n--\n%s", len(days), str)))
}

// }}}
// {{{ publishFlightsHandler

// http://backend-dot-serfr0-fdb.appspot.com/backend/publish-flights?datestring=yesterday
// http://backend-dot-serfr0-fdb.appspot.com/backend/publish-flights?datestring=2015.09.15

// As well as writing the data into a file in Cloud Storage, it will submit a load
// request into BigQuery to load that file.

func publishFlightsHandler(db fgae.FlightDB, w http.ResponseWriter, r *http.Request) {
	tStart := time.Now()

	datestring := r.FormValue("datestring")
	if datestring == "yesterday" {
		datestring = date.NowInPdt().AddDate(0,0,-1).Format("2006.01.02")
	}

	filename := "flights-"+datestring+".json"
	db.Infof("Starting /batch/publish-flights: %s", filename)
	
	n,err := writeBigQueryFlightsGCSFile(db, r, datestring, folderGCS, filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("skipload") == "" {
		if err := submitLoadJob(db, r, folderGCS, filename); err != nil {
			http.Error(w, "submitLoadJob failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("OK!\n%d flights written to gs://%s/%s and job sent - took %s\n",
		n, folderGCS, filename, time.Since(tStart))))
}

// }}}

// {{{ writeBigQueryFlightsGCSFile

// Returns number of records written (which is zero if the file already exists)
func writeBigQueryFlightsGCSFile(db fgae.FlightDB, r *http.Request, datestring, foldername,filename string) (int,error) {
	ctx := db.Ctx()
	
	if exists,err := gcs.Exists(ctx, foldername, filename); err != nil {
		return 0,err
	} else if exists {
		return 0,nil
	}	
	gcsHandle,err := gcs.OpenRW(ctx, foldername, filename, "application/json")
	if err != nil {
		return 0,err
	}
	encoder := json.NewEncoder(gcsHandle.IOWriter())
	
	tags := widget.FormValueCommaSpaceSepStrings(r,"tags")
	s := date.Datestring2MidnightPdt(datestring)
	e := s.AddDate(0,0,1).Add(-1 * time.Second) // +23:59:59 (or 22:59 or 24:59 when going in/out DST)

	n := 0

	q := fgae.QueryForTimeRange(tags, s, e)
	it := db.NewIterator(q)
	for it.Iterate(ctx) {
		f := it.Flight()

		// A flight that straddles midnight will have timeslots either
		// side, and so will end up showing in the results for two
		// different days. We don't want dupes in the aggregate output, so
		// we should only include the flight in one of them; we pick the
		// first day. So if the flight's first timeslot does not start
		// after our start-time, skip it.
		if slots := f.Timeslots(); len(slots)>0 && slots[0].Before(s) {
			continue
		}
		
		n++
		fbq := f.ForBigQuery()

		if err := encoder.Encode(fbq); err != nil {
			return 0,err
		}
	}
	if it.Err() != nil {
		return 0,fmt.Errorf("iterator [%s,%s] failed at %s: %v", s,e, time.Now(), it.Err())
	}

	if err := gcsHandle.Close(); err != nil {
		return 0, err
	}

	db.Infof("GCS bigquery file '%s' successfully written", filename)

	return n,nil
}

// }}}
// {{{ submitLoadJob

// https://cloud.google.com/bigquery/docs/loading-data-cloud-storage#bigquery-import-gcs-file-go
func submitLoadJob(db fgae.FlightDB, r *http.Request, gcsfolder, gcsfile string) error {
	ctx := db.Ctx()

	client,err := bigquery.NewClient(ctx, bigqueryProject)
	if err != nil {
		return fmt.Errorf("Creating bigquery client: %v", err)
	}
	myDataset := client.Dataset(bigqueryDataset)
	destTable := myDataset.Table(bigqueryTableName)
	
	gcsSrc := bigquery.NewGCSReference(fmt.Sprintf("gs://%s/%s", gcsfolder, gcsfile))
	gcsSrc.SourceFormat = bigquery.JSON
	gcsSrc.AllowJaggedRows = true

	loader := destTable.LoaderFrom(gcsSrc)
	loader.CreateDisposition = bigquery.CreateNever
	job,err := loader.Run(ctx)	
	if err != nil {
		return fmt.Errorf("Submission of load job: %v", err)
	}

	// https://godoc.org/cloud.google.com/go/bigquery#Copier
/*
	tableDest := &bigquery.Table{
		ProjectID: bigqueryProject,
		DatasetID: bigqueryDataset,
		TableID:   bigqueryTableName,
	}
	copier := myDataset.Table(bigqueryTableName).CopierFrom(gcsSrc)
	copier.WriteDisposition = bigquery.WriteAppend
	job,err := copier.Run(ctx)
	if err != nil {
		return fmt.Errorf("Submission of load job: %v", err)
	}
*/	
	//job,err := client.Copy(ctx, tableDest, gcsSrc, bigquery.WriteAppend)
	//if err != nil {
	//	return fmt.Errorf("Submission of load job: %v", err)
	//}

	time.Sleep(5 * time.Second)
	
	if status, err := job.Status(ctx); err != nil {
		return fmt.Errorf("Failure determining status: %v", err)
	} else if err := status.Err(); err != nil {
		detailedErrStr := ""
		for i,innerErr := range status.Errors {
			detailedErrStr += fmt.Sprintf(" [%2d] %v\n", i, innerErr)
		}
		db.Errorf("BiqQuery LoadJob error: %v\n--\n%s", err, detailedErrStr)
		return fmt.Errorf("Job error: %v\n--\n%s", err, detailedErrStr)
	} else {
		db.Infof("BiqQuery LoadJob status: done=%v, state=%s, %s",
			status.Done(), status.State, status)
	}
	
	return nil
}

// }}}

// {{{ -------------------------={ E N D }=----------------------------------

// Local variables:
// folded-file: t
// end:

// }}}

