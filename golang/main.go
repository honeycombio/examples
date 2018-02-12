package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"golang.org/x/crypto/bcrypt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/honeycombio/libhoney-go"
	"github.com/jmoiron/sqlx"
)

var (
	decoder      = schema.NewDecoder()
	sessionName  = "default"
	sessionStore = sessions.NewCookieStore([]byte("best-secret-in-the-world"))
	baseTmpl     = `
<!doctype html>
<html>
<head>
<title>Shoutr</title>
</head>
<body>
{{template "body" .}}
</body>
`
	signupTmpl = `
{{define "body"}}
<h2>Sign Up</h2>
{{if .ErrorMessage}}
<p style="color: red;">{{.ErrorMessage}}</p>
{{end}}
<form action="/signup" method="POST">
<div>
<label><b>First Name:</b></label>
<input type="text" placeholder="First Name" name="first_name" required>
</div>

<div>
<label><b>Last Name:</b></label>
<input type="text" placeholder="Last Name" name="last_name" required>
</div>

<div>
<label><b>Username:</b></label>
<input type="text" placeholder="Username" name="username" required>
</div>

<div>
<label><b>Email (optional):</b></label>
<input type="text" placeholder="Email" name="email" required>
</div>

<div>
<label><b>Password:</b></label>
<input type="password" placeholder="Password" name="password" required>
</div>

<div>
<label><b>Confirm Password:</b></label>
<input type="password" placeholder="Confirm Password" name="repeated_password" required>
</div>

<div>
<button type="submit">Sign Up</button>
</div>
</form>
{{end}}
`
	loginTmpl = `
{{define "body"}}
{{if .ErrorMessage}}
<p style="color: red;">{{.ErrorMessage}}</p>
{{end}}
<form action="/login" method="POST">
<div>
<label><b>Username:</b></label>
<input type="text" placeholder="Username" name="username" required>
</div>

<div>
<label><b>Password:</b></label>
<input type="password" placeholder="Password" name="password" required>
</div>

<div>
<button type="submit">Login</button>
</div>
</form>
{{end}}
`
	mainTmpl = `
{{define "body"}}
{{if ne .User.ID 0}}
<form action="/logout" method="POST">
<button type="submit">Logout</button>
</form>
<p>
Welcome {{.User.FirstName}}.
</p>
<h3>Get shoutin':</h3>
<form action="/shout" method="POST">
<textarea rows="4" cols="80" name="content" required>
</textarea>
<button type="submit">Shout!</button>
</form>
{{if .Shouts}}
{{range $shout := .Shouts}}
<div style="margin-bottom: 10px; font-size: 1.1rem;">
<div style="font-size: 0.8rem;"><b>{{$shout.FirstName}} {{$shout.LastName}}</b> @{{$shout.Username}} | {{$shout.CreatedAt.Time.Format "Jan 02, 2006 15:04:05"}}</div>
{{$shout.Content}}
</div>
{{end}}
{{else}}
<i>Once you or others do some shouting, the shouts will appear here.</i>
{{end}}
{{else}}
<h1>Shoutr</h1>
<p>Shoutr is a new kind of web 3.0 social media platform.</p>
<p>With Shoutr, you can shout your opinions on the Internet!</p>
<p>Sign up for an account today to access the content in our walled garden.</p>
<a href="/signup">Sign Up</a> |
<a href="/login">Login</a>
{{end}}
{{end}}
`
	db *sqlx.DB
)

func init() {
	var err error
	dbUser := "root"
	dbPass := ""
	dbName := "shoutr"
	db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPass, dbName))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INT NOT NULL AUTO_INCREMENT,
	password VARCHAR(64) NOT NULL,
	first_name VARCHAR(64) NOT NULL,
	last_name VARCHAR(64) NOT NULL,
	username VARCHAR(64) NOT NULL,
	email VARCHAR(64),
	PRIMARY KEY (id),
	UNIQUE KEY (username)
);`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS shouts (
	id INT NOT NULL AUTO_INCREMENT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	user_id INT,
	content VARCHAR(280) NOT NULL,
	PRIMARY KEY (id)
);
`)
	if err != nil {
		panic(err)
	}

	hcConfig := libhoney.Config{
		WriteKey: os.Getenv("HC_WRITEKEY"),
		Dataset:  "shoutr-main",
	}
	if err := libhoney.Init(hcConfig); err != nil {
		panic(err)
	}
	if _, err := libhoney.VerifyWriteKey(hcConfig); err != nil {
		panic(err)
	}
}

type User struct {
	ID               int    `db:"id"`
	FirstName        string `db:"first_name" schema:"first_name"`
	LastName         string `db:"last_name" schema:"last_name"`
	Username         string `db:"username" schema:"username"`
	Email            string `db:"email" schema:"email"`
	Password         string `db:"password" schema:"password"`
	RepeatedPassword string `schema:"repeated_password"`
}

type Shout struct {
	ID        int            `db:"int"`
	UserID    int            `db:"user_id"`
	Content   string         `db:"content"`
	CreatedAt mysql.NullTime `db:"created_at"`
}

// Used to read the data from a MySQL JOIN query and render it on the
// front-end.
type RenderedShout struct {
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name" schema:"last_name"`
	Username  string         `db:"username" schema:"username"`
	Content   string         `db:"content"`
	CreatedAt mysql.NullTime `db:"created_at"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(baseTmpl + signupTmpl))
	tmplData := struct {
		ErrorMessage string
	}{}
	if r.Method == "GET" {
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			tmplData.ErrorMessage = "Couldn't parse form"
			if err := tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		var user User
		if err := decoder.Decode(&user, r.PostForm); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			tmplData.ErrorMessage = "An error occurred"
			if err := tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		if user.Password != user.RepeatedPassword {
			w.WriteHeader(http.StatusBadRequest)
			tmplData.ErrorMessage = "Passwords don't match"
			if err := tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}

		if err := validation.ValidateStruct(&user,
			validation.Field(&user.FirstName, is.Alpha),
			validation.Field(&user.LastName, is.Alpha),
			validation.Field(&user.Username, is.Alphanumeric),
			validation.Field(&user.Email, is.Email),
		); err != nil {
			log.Print(err)
			tmplData.ErrorMessage = "Validation failure"
			if err := tmpl.Execute(w, tmplData); err != nil {
				log.Print(err)
			}
			return
		}

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

		session, _ := sessionStore.Get(r, sessionName)
		userID, err := res.LastInsertId()
		if err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		session.Values["user_id"] = int(userID)
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(baseTmpl + loginTmpl))
	tmplData := struct {
		ErrorMessage string
	}{}

	if r.Method == "GET" {
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	user := User{}

	if err := r.ParseForm(); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "Couldn't parse form properly"
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	pass := r.FormValue("password")
	username := r.FormValue("username")

	if err := db.Get(&user, `SELECT id, password FROM users WHERE username = ?`, username); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "Couldn't log you in."
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmplData.ErrorMessage = "That's not a valid password, you sneaky devil."
		if err := tmpl.Execute(w, tmplData); err != nil {
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
	session, _ := sessionStore.Get(r, sessionName)
	userID := session.Values["user_id"]
	if err := r.ParseForm(); err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	if _, err := db.Exec(`INSERT INTO shouts (content, user_id) VALUES (?, ?)`, content, userID); err != nil {
		log.Print(err)
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionName)
	delete(session.Values, "user_id")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("").Parse(baseTmpl + mainTmpl))
	session, _ := sessionStore.Get(r, sessionName)
	tmplData := struct {
		User   User
		Shouts []RenderedShout
	}{}
	// Not logged in
	if userID, ok := session.Values["user_id"]; !ok {
		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	} else {
		if err := db.Get(&tmplData.User, `SELECT * FROM users WHERE id = ?`, userID); err != nil {
			log.Print(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}

		if err := db.Select(&tmplData.Shouts, `
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

		if err := tmpl.Execute(w, tmplData); err != nil {
			log.Print(err)
		}
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/signup", signupHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/shout", shoutHandler)
	log.Print("Serving app on localhost:8888 ....")
	log.Fatal(http.ListenAndServe(":8888", r))
}
