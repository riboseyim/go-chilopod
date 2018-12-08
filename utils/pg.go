package utils

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func initPGConn() (*sql.DB, error) {
	//driver
	db, err := sql.Open("postgres", "user=tianguan password=tianguan123 dbname=TTank sslmode=disable")
	checkDBErr(err)
	return db, err
}

func checkDBErr(err error) {
	if err != nil {
		panic(err)
	}
}
