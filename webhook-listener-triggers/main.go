package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnygorilla"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
)

// Config configures this application
type Config struct {
	// SharedSecrets is an authentication key you get from the Triggers UI when
	// you create a webhook trigger recipient. You should only accept POSTs to
	// this webhook that have the secret. Secrets are a per-webhook endpoint
	// config, so if multiple webhook recipients are configured to send to this
	// app, multiple shared secrets will be necessary.
	SharedSecretTokens []string
	// Port is the port on localhost on which this webhook will listen. Default
	// 8080
	Port int
	// Output is the place we'll write the record of receiving notifications
	// from Honeycomb.  Default STDOUT
	Output io.Writer
}

// TriggerResult represents a single row in the table of results that come back
// from a trigger. A notification will have two lists of TriggerResults, one
// that is all the rows in the results table and the other that is only the rows
// that have crossed the threshold specified in the trigger
type TriggerResult struct {
	// Groups are the breakdown columns and values in those columns that are
	// present in the results or have triggered the threshold.
	Group map[string]interface{}
	// Result is the value of the Calculation for this group of columns
	Result float64
}

// Notification is the message we'll get from the Honeycomb API
type Notification struct {
	// Version is the version of this notification - changes to the structure of
	// this message will trigger changes to this version string
	Version string `json:"version"`
	// SharedSecret is configured on a per-webhook basis
	SharedSecret string `json:"shared_secret"`
	// TriggerName is the name of this trigger, as configured in the UI
	TriggerName string `json:"name"`
        TriggerID string   `json:"id"`
	// Status will be TRIGGERED or OK
	Status          string          `json:"status"`
	Summary         string          `json:"summary"`
	Description     string          `json:"description"`
	Operator        string          `json:"operator"`
	Threshold       float64         `json:"threshold"`
	ResultURL       string          `json:"result_url"` // permalink to the trigger results
	ResultGroups    []TriggerResult `json:"result_groups"`
	GroupsTriggered []TriggerResult `json:"result_groups_triggered"`
	TriggerURL      string          `json:"trigger_url"`
	// Timestamp does not come with the notification but I want it to be
	// serialized in the output here so we can see when things came in
	Timestamp time.Time `json:"timestamp"`
}

type App struct {
	conf Config
}

func main() {

	beeline.Init(beeline.Config{
		WriteKey: "abcabc123123",
		Dataset:  "trigger-webhook",
		// for demonstration, send the event to STDOUT intead of Honeycomb.
		// Remove the STDOUT setting when filling in a real write key.
		// STDOUT: true,
	})

	var a = &App{}
	a.conf = getConfig()

	r := mux.NewRouter()
	r.Use(hnygorilla.Middleware)
	r.HandleFunc("/notify", a.rcvNotification).Methods("POST")
	r.HandleFunc("/", a.defaultPath)

	listenAddr := fmt.Sprintf(":%d", a.conf.Port)
	log.Fatal(http.ListenAndServe(listenAddr, hnynethttp.WrapHandler(r)))
}

func getConfig() Config {
	return Config{
		SharedSecretTokens: []string{"would you like to play a game"},
		Port:               8090,
		Output:             os.Stdout,
	}
}

func (a *App) defaultPath(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Behold the Kroot in the Forest!\n"))
}

func (a *App) rcvNotification(w http.ResponseWriter, r *http.Request) {
	// first validate that the shared secret is legit
	token := r.Header.Get("X-Honeycomb-Webhook-Token")
	var matched bool
	if token != "" {
		for _, ss := range a.conf.SharedSecretTokens {
			if token == ss {
				matched = true
			}
		}
	}
	if !matched {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("failed to authenticate notification\n"))
		return
	}

	// ok, it's valid, let's parse the notification and handle it.
	defer r.Body.Close()
	bod, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to read body\n"))
		return
	}
	var notif = Notification{}
	err = json.Unmarshal(bod, &notif)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to parse json\n"))
		return
	}
	// hooray, let's write out our notification
	notif.Timestamp = time.Now()
	serialized, err := json.Marshal(notif)
	a.conf.Output.Write(serialized)
	a.conf.Output.Write([]byte("\n"))
	// and respond happily
	w.Write([]byte("accepted\n"))
}
