package main

import (
	"database/sql"

	"github.com/TheSgtPepper23/coordinator/models"
	_ "modernc.org/sqlite"
)

type Config struct {
	DB     *sql.DB
	Models models.Models
}

func main() {
	conn, err := connectToDB()

	if err != nil {
		panic(err)
	}

	models.New(conn)

	maps, err := models.GetMaps()

	if err != nil {
		panic(err)
	}

	map2 := maps[0]

	map2.DeleteMap()

	if err != nil {
		panic(err)
	}
}

// Return a connection to the local database and if not created, creates the tables needed
func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./database/main.db")

	if err != nil {
		return nil, err
	}

	allTables := "CREATE TABLE IF NOT EXISTS maps (id integer primary key, name text not null, created_at datetime not null, version text not null);"
	allTables = allTables + "CREATE TABLE IF NOT EXISTS coordinates (id integer primary key, name text not null, created_at datetime not null, xvalue float not null,  yvalue float not null, zvalue float not null, mapid integer not null, FOREIGN KEY(mapid) REFERENCES maps(id));"
	statement, err := db.Prepare(allTables)

	if err != nil {
		return nil, err
	}

	_, err = statement.Exec()

	if err != nil {
		return nil, err
	}

	return db, nil
}
