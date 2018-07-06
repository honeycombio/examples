package main

import (
	"net/http"
	"time"
)

const (
	HeaderWriteKey   string = "X-Honeycomb-Team"
	HeaderTimestamp  string = "X-Honeycomb-Event-Time"
	HeaderSampleRate string = "X-Honeycomb-Samplerate"
)

const ValidWriteKeyCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const CacheTimeout = 10 * time.Second

// apierr lets us create a bunch of sanitized user-visible errors
type apierr int

const (
	apierrUnknownError apierr = iota
	apierrParseFailure
	apierrAuthMishapen
	apierrAuthFailure
	apierrJSONFailure
	apierrDatasetLookupFailure
	apierrSchemaLookupFailure
)

type SafeError struct {
	responseCode int
	responseBody string
}

var responses = map[apierr]*SafeError{
	apierrParseFailure:         &SafeError{http.StatusBadRequest, `{"error":"unable to parse request headers"}`},
	apierrAuthMishapen:         &SafeError{http.StatusUnauthorized, `{"error":"writekey malformed - expect only letters and numbers"}`},
	apierrAuthFailure:          &SafeError{http.StatusUnauthorized, `{"error":"writekey didn't match valid credentials"}`},
	apierrJSONFailure:          &SafeError{http.StatusBadRequest, `{"error":"failed to unmarshal JSON body"}`},
	apierrDatasetLookupFailure: &SafeError{http.StatusBadRequest, `{"error":"failed to resolve dataset object"}`},
	apierrSchemaLookupFailure:  &SafeError{http.StatusInternalServerError, `{"error":"failed to resolve schema"}`},
}

// Error lets SafeError implement the error interface
func (s *SafeError) Error() string {
	return s.responseBody
}

type Event struct {
	WriteKey        string
	Timestamp       time.Time
	SampleRate      int
	Data            map[string]interface{}
	ChosenPartition int
}

type Team struct {
	ID       int
	Name     string
	WriteKey string
}

// knownTeams maps writekey to a Team struct
var knownTeams = []*Team{
	&Team{1, "RPO", "abcd123EFGH"},
	&Team{2, "b&w", "ijkl456MNOP"},
	&Team{3, "Third", "qrst789UVWX"},
}

type Dataset struct {
	ID            int
	Name          string
	PartitionList []int
}

// knownDatasets maps datasetName to a Dataset struct
var knownDatasets = []*Dataset{
	&Dataset{1, "wade", []int{1, 2, 3}},
	&Dataset{2, "james", []int{1, 2, 4}},
	&Dataset{3, "helen", []int{1, 3, 4}},
	&Dataset{4, "peter", []int{1, 2, 4}},
	&Dataset{5, "valentine", []int{1, 3, 4}},
	&Dataset{6, "andrew", []int{2, 3, 4}},
}
