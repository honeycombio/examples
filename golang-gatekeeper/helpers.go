package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// userFacingErr takes an error type and formats an appropriate HTTP response
// for that type of error.
func userFacingErr(ctx context.Context, err error, errType apierr, w http.ResponseWriter) {
	// if we got a user-safe error, use that. otherwise use errType
	if err, ok := err.(*SafeError); ok {
		w.WriteHeader(err.responseCode)
		w.Write([]byte(err.responseBody))
		return
	}
	w.WriteHeader(responses[errType].responseCode)
	w.Write([]byte(responses[errType].responseBody))
}

func validateWritekey(ctx context.Context, wk string) (*Team, error) {
	// writekeys are only [a-zA-Z0-9]
	for _, char := range wk {
		if !strings.Contains(ValidWriteKeyCharset, string(char)) {
			// send a user-safe error here to propagate the difference up
			return nil, responses[apierrAuthMishapen]
		}
	}

	// authenticate

	// here we would call out to the database to validate the writekey
	// but instead in the interests of simplicity
	// we're just going to check that it's one of the few valid writekeys
	// TODO add an in-memory cache and take a longer time to hit the "database"
	for _, team := range knownTeams {
		if team.WriteKey == wk {
			return team, nil
		}
	}
	return nil, responses[apierrAuthFailure]
}

func resolveDataset(ctx context.Context, datasetName string) (*Dataset, error) {
	// here we would call out to the database to fetch the dataset object
	// or create it if one didn't exist
	// instead we're just going to take one from a known set of datasets
	for _, ds := range knownDatasets {
		if ds.Name == datasetName {
			return ds, nil
		}
	}
	return nil, errors.New("dataset not found")
}

func unmarshal(r *http.Request, ev *Event) error {
	// always close the body when done reading
	defer r.Body.Close()

	// set up a plaintext reader to abstract out the gzipping
	var reader io.Reader
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		buf := bytes.Buffer{}
		var err error
		if _, err = io.Copy(&buf, r.Body); err != nil {
			return err
		}
		if reader, err = gzip.NewReader(&buf); err != nil {
			return err
		}
	default:
		reader = r.Body
	}

	body := make(map[string]interface{})
	if err := json.NewDecoder(reader).Decode(&body); err != nil {
		return err
	}
	ev.Data = body
	return nil
}

// getPartition returns a random partition chosen from the available partition list
func getPartition(ctx context.Context, ds *Dataset) (int, error) {
	parts := ds.PartitionList
	if len(parts) <= 0 {
		return 0, errors.New("no partitions found")
	}
	partIndex := rand.Intn(len(parts))
	return parts[partIndex], nil
}

// lastCacheTime lets us simulate having the schema cached in memory, and
// falling through to the database every so often (in this case every 10sec).
// The first call to getSchema will "fall through" (i.e. take 30-50ms longer),
// and the rest of the requests for the following 10 seconds will be fast.
var lastCacheTime time.Time

// getSchema pretends to fetch the schema from a database and cache it
func getSchema(ctx context.Context, dataset *Dataset) error {
	if time.Since(lastCacheTime) > CacheTimeout {
		// we fall through the cache every 10 seconds. In production this might
		// be closer to 5 minutes
		// pretend to hit a slow database that takes 30-50ms
		time.Sleep(time.Duration((rand.Intn(20) + 30)) * time.Millisecond)
		lastCacheTime = time.Now()
	}
	// let's just fail sometimes to pretend
	if rand.Intn(60) == 0 {
		return errors.New("failed to get dataset schema")
	}
	return nil
}

func writeEvent(ctx context.Context, ev *Event) error {
	// TODO should probably add a timer here too, right?
	fileNmae := fmt.Sprintf("/tmp/api%d.log", ev.ChosenPartition)
	serialized, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	fh, err := os.OpenFile(fileNmae, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fh.Close()
	fh.Write(serialized)
	fh.Write([]byte("\n"))

	return err
}
