package main

import "github.com/go-sql-driver/mysql"

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name" schema:"first_name"`
	LastName  string `db:"last_name" schema:"last_name"`
	Username  string `db:"username" schema:"username"`
	Email     string `db:"email" schema:"email"`
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
