package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func Connect(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("database ping failed: ", err)
	}

	log.Println("database connected")
	return db

}
