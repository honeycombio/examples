package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnygorilla"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize beeline. The only required field is WriteKey.
	wk := os.Getenv("HONEYCOMB_WRITEKEY")
	if wk == "" {
		logrus.Error("got empty writekey from the environment. Please set HONEYCOMB_WRITEKEY")
	}
	beeline.Init(beeline.Config{
		WriteKey: wk,
		Dataset:  "gatekeeper-go",
		// for demonstration, send the event to STDOUT intead of Honeycomb.
		// Remove the STDOUT setting when filling in a real write key.
		// STDOUT: true,
	})
	// augment our dataset with a few extra useful fields
	addCommonLibhoneyFields()

	r := mux.NewRouter()
	r.Use(hnygorilla.Middleware)
	// Routes consist of a path and a handler function.
	r.HandleFunc("/1/events/{datasetName}", handleEvent).Methods("POST")
	r.HandleFunc("/x/alive", healthHandler)
	r.HandleFunc("/", home)

	// Bind to a port and pass our router in
	fmt.Printf("Serving on http://localhost:8080/ ...\n")
	log.Fatal(http.ListenAndServe(":8080", hnynethttp.WrapHandler(r)))
}

// home serves a mild but reasonable HTML response for requests to /
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>
	<body>
		Thanks for playing with <a href="https://www.honeycomb.io">Honeycomb.io</a>
	</body>
</html>
`))
}

// healthHandler provides a simple endpoint at which to point a health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
	beeline.AddField(r.Context(), "alive", true)
	w.Write([]byte(`{"alive": "yes"}`))
}

// handleEvent is the real purpose of this example API server
func handleEvent(w http.ResponseWriter, r *http.Request) {
	ev := &Event{}

	// get writekey, timestamp, and sample rate out of HTTP headers
	if err := getHeaders(r, ev); err != nil {
		userFacingErr(r.Context(), err, apierrParseFailure, w)
		return
	}

	// authenticate writekey or return 401
	team, err := validateWritekey(r.Context(), ev.WriteKey)
	if err != nil {
		userFacingErr(r.Context(), err, apierrAuthFailure, w)
		return
	}
	beeline.AddField(r.Context(), "team", team)

	// use the dataset name to get back a dataset object
	vars := mux.Vars(r)
	datasetName := vars["datasetName"]
	dataset, err := resolveDataset(r.Context(), datasetName)
	if err != nil {
		userFacingErr(r.Context(), err, apierrDatasetLookupFailure, w)
		return
	}
	beeline.AddField(r.Context(), "dataset", dataset)

	// parse JSON body
	err = unmarshal(r, ev)
	if err != nil {
		userFacingErr(r.Context(), err, apierrJSONFailure, w)
		return
	}
	beeline.AddField(r.Context(), "event_columns", len(ev.Data))

	// get partition info - stub out
	partition, err := getPartition(r.Context(), dataset)
	if err != nil {
		userFacingErr(r.Context(), err, apierrDatasetLookupFailure, w)
		return
	}
	ev.ChosenPartition = partition
	beeline.AddField(r.Context(), "chosen_partition", partition)

	// check time - use or set to now if broken
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now()
	} else {
		// record the difference between the event's timestamp and now to help identify lagging events
		eventTimeDelta := float64(time.Since(ev.Timestamp)) / float64(time.Second)
		beeline.AddField(r.Context(), "event_time_delta_sec", eventTimeDelta)
	}
	beeline.AddField(r.Context(), "event_time", ev.Timestamp)

	// verify schema - stub out
	err = getSchema(r.Context(), dataset)
	if err != nil {
		userFacingErr(r.Context(), err, apierrSchemaLookupFailure, w)
		return
	}

	// ok, everything checks out. Hand off to external service (aka write to
	// local disk)
	writeEvent(r.Context(), ev)
	return
}
