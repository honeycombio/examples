package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	libhoney "github.com/honeycombio/libhoney-go"
)

var (
	hostname       string
	hnyDatasetName = "examples.golang-webapp"
	hnyContextKey  = "honeycombEvent"
)

func init() {
	hcConfig := libhoney.Config{
		WriteKey: os.Getenv("HONEYCOMB_WRITEKEY"),
		Dataset:  hnyDatasetName,
	}

	// This will ensure that our libhoney events get printed to the
	// console. This allows for easier iterating and debugging of
	// instrumentation.
	if os.Getenv("ENV") == "dev" {
		hcConfig.Output = &libhoney.WriterOutput{}
	}

	if err := libhoney.Init(hcConfig); err != nil {
		log.Print(err)
		os.Exit(1)
	}

	if hnyTeam, err := libhoney.VerifyWriteKey(hcConfig); err != nil {
		log.Print(err)
		log.Print("Please make sure the HONEYCOMB_WRITEKEY environment variable is set.")
		os.Exit(1)
	} else {
		log.Print(fmt.Sprintf("Sending Honeycomb events to the %q dataset on %q team", hnyDatasetName, hnyTeam))
	}

	// Initialize fields that every sent event will have.
	if hostname, err := os.Hostname(); err == nil {
		libhoney.AddField("system.hostname", hostname)
	}
	libhoney.AddDynamicField("runtime.num_goroutines", func() interface{} {
		return runtime.NumGoroutine()
	})
	libhoney.AddDynamicField("runtime.memory_inuse", func() interface{} {
		// Could leave this out if performance sensitive. However, it's
		// used here to demonstrate dynamic event fields. This will
		// ensure that every event includes informaiton about the
		// memory usage of the process at the time the event was sent.
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		return mem.Alloc
	})
}

type HoneyResponseWriter struct {
	*libhoney.Event
	http.ResponseWriter
}

func (hrw *HoneyResponseWriter) WriteHeader(status int) {
	hrw.Event.AddField("response.status_code", status)
	hrw.ResponseWriter.WriteHeader(status)
}

func addRequestProps(req *http.Request, ev *libhoney.Event) {
	// Add a variety of details about the HTTP request, such as user agent
	// and method, to any created libhoney event.
	ev.AddField("request.method", req.Method)
	ev.AddField("request.path", req.URL.Path)
	ev.AddField("request.host", req.URL.Host)
	ev.AddField("request.proto", req.Proto)
	ev.AddField("request.content_length", req.ContentLength)
	ev.AddField("request.remote_addr", req.RemoteAddr)
	ev.AddField("request.user_agent", req.UserAgent())
}

// HoneycombMiddleware will wrap our HTTP handle funcs to automatically
// generate an event-per-request and set properties on them.
func HoneycombMiddleware(fn func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// We'll time each HTTP request and add that as a property to
		// the sent Honeycomb event, so start the timer for that.
		startHandler := time.Now()
		ev := libhoney.NewEvent()

		defer func() {
			if err := ev.Send(); err != nil {
				log.Print("Error sending libhoney event: ", err)
			}
		}()

		addRequestProps(r, ev)

		// Create a context where we will store the libhoney event. We
		// will add default values to this event for every HTTP
		// request, and the user can access it to add their own
		// (powerful, custom) fields.
		ctx := context.WithValue(r.Context(), hnyContextKey, ev)
		reqWithContext := r.WithContext(ctx)

		// TODO: Probably not quite right
		if _, ok := ev.Fields()["response.status_code"]; !ok {
			ev.AddField("response.status_code", 200)
		}

		// HoneyResponseWriter will capture the HTTP status code for us
		// and add it to the libhoney per-request event.
		fn(&HoneyResponseWriter{
			Event:          ev,
			ResponseWriter: w,
		}, reqWithContext)

		handlerDuration := time.Since(startHandler)
		ev.AddField("timers.total_time_ms", handlerDuration/time.Millisecond)
	}
}
