# TODO

Goal: move to go113 runtime
 - which means removing all / appengine


TODO
 - pi/airspace/memcache, then the rest of pi/airspace/realtime

 - remove flightdb/ref/*-old.go (pi uses them)
 - remove pi/airspace/memcache
 - mv flightdb/app/*/static/ and templates/ into a sensible place - but they are symlinked from other repos




















# flightdb - a database for flight tracks

Prerequisites:
* the [Go programming language](https://golang.org/dl/)
* define your Go workspace: `export GOPATH=~/go`

Download and run things locally
* `go get github.com/skypies/flightdb/app/frontend` (pulls down all dependencies)
* `goapp serve $GOPATH/github.com/skypies/flightdb/app/frontend` (build & run locally)
* Look at <http://localhost:8080/> (appengine admin panel is <http://localhost:8000/>)

To deploy everything into a Google Cloud project:

    $ gcloud app deploy --project=serfr0-fdb app/frontend --version=one
    $ gcloud app deploy --project=serfr0-fdb app/backend  --version=one

    $ gcloud app deploy --project=serfr0-fdb app/dispatch.yaml
    $ gcloud app deploy --project=serfr0-fdb app/queues.yaml
    $ gcloud app deploy --project=serfr0-fdb app/cron.yaml
    $ gcloud app deploy --project=serfr0-fdb app/index.yaml

If you want it to accumulate realtime flight track data, you'll also want to:
* deploy `github.com/skypies/pi/skypi` onto some Raspberry Pi receivers
* deploy `github.com/skypies/pi/consolidator` into a VM inside your project

The skypies will post bundles of received ADSB (and perhaps MLAT)
messages up to Google PubSub, every second or so. The consolidator
will read those bundles, group them by airframe, and add them into the
database.

If you have CSV dumps of historical flight track data (perhaps from
the FAA), you can import it using the code in
`github.com/skypies/flightdb/app/backend/foia.go`.
