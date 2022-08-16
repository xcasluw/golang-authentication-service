package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"xcasluw/golang-authentication-service/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting authentication service")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to postgress")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgress not yet ready ...")
			counts++
		} else {
			log.Println("connected to postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for tow seconds...")
		time.Sleep(time.Second * 2)
		continue
	}
}
