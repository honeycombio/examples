// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/honeycombio/libhoney-go"
)

type key int

const (
	requestIDKey key = 0
	parentIDKey  key = 1
)

// Define some wrappers to propagate "trace" or "request" identifiers down the
// call stack, to unify the various spans within a trace.
func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get("X-Request-ID")
	if reqID == "" {
		reqID = newID()
	}
	return context.WithValue(ctx, requestIDKey, reqID)
}

func requestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

// Define some wrappers to propagate "parent" identifiers down the call stack,
// for defining relationships between spans within a trace.
func newContextWithParentID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, parentIDKey, id)
}

func parentIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(parentIDKey).(string); ok {
		return id
	}
	return ""
}

// Generate a new unique identifier for our spans and traces. This can be any
// unique string -- Zipkin uses hex-encoded base64 ints, as we do here; other
// folks may prefer to use their UUID library of choice.
func newID() string {
	return fmt.Sprintf("%x", rand.Int63())
}

// Page represents the data (and some basic operations) on a wiki page.
//
// While the tracing instrumentation in this example is constrained to the
// handlers, we could just as easily propagate context down directly into this
// class if needed.
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save(ctx context.Context) error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(ctx context.Context, title string) (*Page, error) {
	filename := title + ".txt"
	id := newID()
	start := time.Now()
	body, err := ioutil.ReadFile(filename)
	sendSpan("ioutil.ReadFile", id, start, ctx, map[string]interface{}{"title": title, "bodylen": len(body), "error": err})
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// Our "View" handler. Tries to load a page from disk and render it. Falls back
// to the Edit handler if the page does not yet exist.
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	loadPageID := newID()
	loadPageStart := time.Now()
	p, err := loadPage(newContextWithParentID(r.Context(), loadPageID), title)
	sendSpan("loadPage", loadPageID, loadPageStart, r.Context(), map[string]interface{}{"title": title, "error": err})
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderID := newID()
	renderStart := time.Now()
	renderTemplate(newContextWithParentID(r.Context(), renderID), w, "view", p)
	sendSpan("renderTemplate", renderID, renderStart, r.Context(), map[string]interface{}{"template": "view"})
}

// Our "Edit" handler. Tries to load a page from disk to seed the edit screen,
// then renders a form to allow the user to define the content of the requested
// wiki page.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	loadPageID := newID()
	loadPageStart := time.Now()
	p, err := loadPage(newContextWithParentID(r.Context(), loadPageID), title)
	sendSpan("loadPage", loadPageID, loadPageStart, r.Context(), map[string]interface{}{"title": title, "error": err})
	if err != nil {
		p = &Page{Title: title}
	}
	renderID := newID()
	renderStart := time.Now()
	renderTemplate(newContextWithParentID(r.Context(), renderID), w, "edit", p)
	sendSpan("renderTemplate", renderID, renderStart, r.Context(), map[string]interface{}{"template": "edit"})
}

// Our "Save" handler simply persists a page to disk.
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	id := newID()
	start := time.Now()
	err := p.save(newContextWithParentID(r.Context(), id))
	sendSpan("ioutil.WriteFile", id, start, r.Context(), map[string]interface{}{"title": title, "bodylen": len(body), "error": err})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))

func renderTemplate(ctx context.Context, w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// This middleware treats each HTTP request as a distinct "trace." Each trace
// begins with a top-level ("root") span indicating that the HTTP request has
// begun.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		start := time.Now()
		id := newID()
		ctx := newContextWithRequestID(r.Context(), r)
		fn(w, r.WithContext(newContextWithParentID(ctx, id)), m[2])
		sendSpan(m[1], id, start, ctx, nil)
	}
}

// This wrapper takes a span name and some optional metadata, then emits a
// "span" to Honeycomb as part of the trace begun in the HTTP middleware.
func sendSpan(name, id string, start time.Time, ctx context.Context, metadata map[string]interface{}) {
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	// Field keys to trigger Honeycomb's tracing functionality on this dataset
	// defined at:
	// https://www.honeycomb.io/docs/working-with-data/tracing/send-trace-data/#manual-tracing
	metadata["name"] = name
	metadata["id"] = id
	metadata["traceId"] = requestIDFromContext(ctx)
	metadata["serviceName"] = "wiki"
	metadata["durationMs"] = float64(time.Since(start)) / float64(time.Millisecond)
	if parentID := parentIDFromContext(ctx); parentID != "" {
		metadata["parentId"] = parentID
	}

	ev := libhoney.NewEvent()
	// NOTE: Don't forget to set the timestamp to `start` -- because spans are
	// emitted at the *end* of their execution, we want to be doubly sure that
	// our manually-emitted events are timestamped with the time that the work
	// (the span's actual execution) really begun.
	ev.Timestamp = start
	ev.Add(metadata)
	ev.Send()
}

// Let's go!
func main() {
	libhoney.Init(libhoney.Config{
		WriteKey: os.Getenv("HONEYCOMB_WRITEKEY"),
		Dataset:  "golang-wiki-tracing-example",
	})
	defer libhoney.Close()

	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	// Redirect to a default wiki page.
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		h := http.RedirectHandler("/view/Index", http.StatusTemporaryRedirect)
		h.ServeHTTP(w, req)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
