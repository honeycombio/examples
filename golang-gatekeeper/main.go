package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/1/events/{datasetName}", handleEvent).Methods("POST")
	r.HandleFunc("/x/alive", healthHandler)
	r.HandleFunc("/", home)

	// Bind to a port and pass our router in
	fmt.Printf("Serving on http://localhost:8080/ ...\n")
	log.Fatal(http.ListenAndServe(":8080", r))
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
	w.Write([]byte(`{"alive": "yes"}`))
}

// handleEvent is the real purpose of this example API server
func handleEvent(w http.ResponseWriter, r *http.Request) {
	ev := &Event{}

	// authenticate writekey or return 401
	_, err := validateWritekey(r.Context(), ev.WriteKey)
	if err != nil {
		userFacingErr(r.Context(), err, apierrAuthFailure, w)
		return
	}

	// use the dataset name to get back a dataset object
	vars := mux.Vars(r)
	datasetName := vars["datasetName"]
	dataset, err := resolveDataset(r.Context(), datasetName)
	if err != nil {
		userFacingErr(r.Context(), err, apierrDatasetLookupFailure, w)
		return
	}

	// parse JSON body
	err = unmarshal(r, ev)
	if err != nil {
		userFacingErr(r.Context(), err, apierrJSONFailure, w)
		return
	}

	// get partition info - stub out
	partition, err := getPartition(r.Context(), dataset)
	if err != nil {
		userFacingErr(r.Context(), err, apierrDatasetLookupFailure, w)
		return
	}
	ev.ChosenPartition = partition

	// check time - use or set to now if broken
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now()
	}

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
