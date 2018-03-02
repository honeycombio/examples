package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	libhoney "github.com/honeycombio/libhoney-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	decoder      = schema.NewDecoder()
	sessionName  = "default"
	sessionStore = sessions.NewCookieStore([]byte("best-secret-in-the-world"))
)

const (
	templatesDir = "templates"
	// 140 is the proper amount of characters for a microblog. Any other
	// value is heresy.
	maxShoutLength = 140
)

func hnyEventFromRequest(r *http.Request) *libhoney.Event {
	ev, ok := r.Context().Value(hnyContextKey).(*libhoney.Event)
	if !ok {
		// We control the way this is being put on context anyway, so
		// this should not happen unless someone is tinkering with that
		// code.
		panic("Couldn't get libhoney event from request context")
	}
	return ev
}

func addFinalErr(err *error, ev *libhoney.Event) {
	if *err != nil {
		ev.AddField("error", (*err).Error())
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ev := hnyEventFromRequest(r)
	defer addFinalErr(&err, ev)

	tmpl := template.Must(template.
		ParseFiles(
			filepath.Join(templatesDir, "base.html"),
			filepath.Join(templatesDir, "signup.html"),
		))
	tmplData := struct {
		ErrorMessage string
	}{}
	if r.Method == "GET" {
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}
	if r.Method == "POST" {
		if err = r.ParseForm(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			tmplData.ErrorMessage = "Couldn't parse form"
			if err = tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		var user User
		if err = decoder.Decode(&user, r.PostForm); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)

			tmplData.ErrorMessage = "An error occurred"
			if err = tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		ev.AddField("user.email", user.Email)

		if user.Password != user.RepeatedPassword {
			w.WriteHeader(http.StatusBadRequest)
			tmplData.ErrorMessage = "Passwords don't match"
			if err = tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		var hashedPassword []byte

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		if err = validation.ValidateStruct(&user,
			validation.Field(&user.FirstName, is.Alpha),
			validation.Field(&user.LastName, is.Alpha),
			validation.Field(&user.Username, is.Alphanumeric),
			validation.Field(&user.Email, is.Email),
		); err != nil {
			log.Print(err)
			tmplData.ErrorMessage = "Validation failure"
			if err = tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		queryStart := time.Now()

		res, err := db.Exec(`INSERT INTO users
(first_name, last_name, username, password, email)
VALUES
(?, ?, ?, ?, ?)
`, user.FirstName, user.LastName, user.Username, hashedPassword, user.Email)
		if err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		ev.AddField("timers.db.users_insert_ms", time.Since(queryStart)/time.Millisecond)

		session, _ := sessionStore.Get(r, sessionName)
		userID, err := res.LastInsertId()
		if err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = int(userID)

		ev.AddField("user.id", int(userID))

		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ev := hnyEventFromRequest(r)
	defer addFinalErr(&err, ev)

	tmpl := template.Must(template.
		ParseFiles(
			filepath.Join(templatesDir, "base.html"),
			filepath.Join(templatesDir, "login.html"),
		))
	tmplData := struct {
		ErrorMessage string
	}{}

	if r.Method == "GET" {
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	user := User{}

	if err = r.ParseForm(); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "Couldn't parse form properly"
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	pass := r.FormValue("password")
	username := r.FormValue("username")

	if err = db.Get(&user, `SELECT id, password FROM users WHERE username = ?`, username); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "Couldn't log you in."
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "That's not a valid password, you sneaky devil."
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	session, _ := sessionStore.Get(r, sessionName)
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func shoutHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ev := hnyEventFromRequest(r)
	defer addFinalErr(&err, ev)

	session, _ := sessionStore.Get(r, sessionName)
	userID := session.Values["user_id"]
	if err = r.ParseForm(); err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ev.AddField("user.id", userID)
	content := r.FormValue("content")
	ev.AddField("shout.content_length", len(content))

	if len(content) > maxShoutLength {
		session, _ := sessionStore.Get(r, sessionName)
		session.AddFlash("Your shout is too long!")
		session.Save(r, w)
		ev.AddField("shout.content", content[:140])
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	ev.AddField("shout.content", content)

	if _, err = db.Exec(`INSERT INTO shouts (content, user_id) VALUES (?, ?)`, content, userID); err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ev := hnyEventFromRequest(r)
	defer addFinalErr(&err, ev)
	session, _ := sessionStore.Get(r, sessionName)
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ev := hnyEventFromRequest(r)
	defer addFinalErr(&err, ev)
	tmpl := template.Must(template.ParseFiles(
		filepath.Join(templatesDir, "base.html"),
		filepath.Join(templatesDir, "home.html"),
	))
	session, _ := sessionStore.Get(r, sessionName)
	tmplData := struct {
		User         User
		Shouts       []RenderedShout
		ErrorMessage string
	}{}

	flashes := session.Flashes()
	if len(flashes) == 1 {
		flash, ok := flashes[0].(string)
		if !ok {
			ev.AddField("flash.err", "Flash didn't assert to type string, got "+reflect.TypeOf(flash).String())
		} else {
			tmplData.ErrorMessage = flash
			ev.AddField("flash.value", flash)
		}
		session.Save(r, w)
	}

	// Not logged in
	if userID, ok := session.Values["user_id"]; !ok {
		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	} else {
		if err = db.Get(&tmplData.User, `SELECT * FROM users WHERE id = ?`, userID); err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if err = db.Select(&tmplData.Shouts, `
SELECT users.first_name, users.last_name, users.username, shouts.content, shouts.created_at
FROM shouts
INNER JOIN users
ON shouts.user_id = users.id
ORDER BY created_at DESC
`); err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if err = tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}
}
