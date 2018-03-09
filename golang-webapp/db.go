package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	maxConnectRetries = 10
	db                *sqlx.DB
)

func init() {
	var err error
	dbUser := "root"
	dbPass := ""
	dbName := "shoutr"
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	for i := 0; i < maxConnectRetries; i++ {
		db, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName))
		if err != nil {
			log.Print("Error connecting to database: ", err)
		} else {
			break
		}
		if i == maxConnectRetries-1 {
			panic("Couldn't connect to DB")
		}
		time.Sleep(1 * time.Second)
	}

	log.Print("Bootstrapping database...")

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id INT NOT NULL AUTO_INCREMENT,
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
	content VARCHAR(140) NOT NULL,
	PRIMARY KEY (id)
);
`)
	if err != nil {
		panic(err)
	}
}
