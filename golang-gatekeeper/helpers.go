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
	"runtime"
	"strconv"
	"strings"
	"time"

	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/timer"
	libhoney "github.com/honeycombio/libhoney-go"
)

// addCommonLibhoneyFields adds a few fields we want in all events
func addCommonLibhoneyFields() {
	// TODO what other fields should we add here for extra color?
	libhoney.AddDynamicField("meta.num_goroutines",
		func() interface{} { return runtime.NumGoroutine() })
	getAlloc := func() interface{} {
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem.Alloc
	}
	libhoney.AddDynamicField("meta.memory_inuse", getAlloc)

	startTime := time.Now()
	libhoney.AddDynamicField("meta.process_uptime_sec", func() interface{} {
		return time.Now().Sub(startTime) / time.Second
	})
}

// getHeaders pulls the three available headers out of the HTTP request and type
// asserts them to the right type.  It does no additional validation.
func getHeaders(r *http.Request, ev *Event) error {
	// add a timer around getting headers
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(r.Context(), "timer.get_headers_dur_ms", dur)
	}(timer.Start())

	// pull raw values from headers
	wk := r.Header.Get(HeaderWriteKey)
	beeline.AddField(r.Context(), HeaderWriteKey, wk)
	ts := r.Header.Get(HeaderTimestamp)
	beeline.AddField(r.Context(), HeaderTimestamp, ts)
	sr := r.Header.Get(HeaderSampleRate)
	beeline.AddField(r.Context(), HeaderSampleRate, sr)

	// assert types

	// writekeys are strings, so no assertion needed
	ev.WriteKey = wk

	// Timestamps should be RFC3339Nano. If we get the zero time that means
	// parsing failed. Leave it at zero until we do real time stuff later
	evTime, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		// it's fine if we can't parse the time (maybe it's missing!)
		// but we should note that we failed and continue
		beeline.AddField(r.Context(), "error_time_parsing", err)
	}
	ev.Timestamp = evTime

	// sample rate should be a positive integer. Defaults to 1 if empty.
	if sr == "" {
		sr = "1"
	}
	sampleRate, err := strconv.Atoi(sr)
	if err != nil {
		return err
	}
	ev.SampleRate = sampleRate
	beeline.AddField(r.Context(), "sample_rate", sampleRate)
	return nil
}

// userFacingErr takes an error type and formats an appropriate HTTP response
// for that type of error.
func userFacingErr(ctx context.Context, err error, errType apierr, w http.ResponseWriter) {
	beeline.AddField(ctx, "error", err.Error())
	beeline.AddField(ctx, "error_desc", responses[errType].responseBody)
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
	// add a timer around validation
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(ctx, "timer.validate_writekey_dur_ms", dur)
	}(timer.Start())

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
	// add a timer around validation
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(ctx, "timer.resolve_dataset_dur_ms", dur)
	}(timer.Start())

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
	// add a timer around unmarshalling json
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(r.Context(), "timer.unmarshal_json_dur_ms", dur)
	}(timer.Start())

	// always close the body when done reading
	defer r.Body.Close()

	// include whether the content was gzipped in the event
	var gzipped bool
	defer beeline.AddField(r.Context(), "gzipped", gzipped)

	// set up a plaintext reader to abstract out the gzipping
	var reader io.Reader
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		gzipped = true
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
	// add a timer around getting the right partition
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(ctx, "timer.get_partition_dur_ms", dur)
	}(timer.Start())

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
	// add a timer around getting the schema
	defer func(t timer.Timer) {
		dur := t.Finish()
		beeline.AddField(ctx, "timer.get_schema_dur_ms", dur)
	}(timer.Start())

	hitCache := true
	if time.Since(lastCacheTime) > CacheTimeout {
		// we fall through the cache every 10 seconds. In production this might
		// be closer to 5 minutes
		hitCache = false
		// pretend to hit a slow database that takes 30-50ms
		time.Sleep(time.Duration((rand.Intn(20) + 30)) * time.Millisecond)
		lastCacheTime = time.Now()
	}
	beeline.AddField(ctx, "hitSchemaCache", hitCache)
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
