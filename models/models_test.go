package models

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

var testDB *sql.DB

/*
Test IDS:

Maps:
*9000
*10000

Coordinates:
*9000
*10000
*11000
*12000
*13000
*14000
*/

func setupDatabase() {
	conn, err := sql.Open("sqlite", "../database/test.db")

	if err != nil {
		panic(err)
	}

	allTables := "CREATE TABLE IF NOT EXISTS maps (id integer primary key, name text not null, created_at datetime not null, version text not null);"
	allTables = allTables + "CREATE TABLE IF NOT EXISTS coordinates (id integer primary key, name text not null, created_at datetime not null, xvalue float not null,  yvalue float not null, zvalue float not null, mapid integer not null, FOREIGN KEY(mapid) REFERENCES maps(id));"
	statement, err := conn.Prepare(allTables)

	if err != nil {
		panic(err)
	}

	_, err = statement.Exec()

	if err != nil {
		panic(err)
	}

	//Adds some elements for the tests
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := conn.BeginTx(ctx, nil)

	if err != nil {
		panic(err)
	}

	defer tx.Rollback()

	//Inserts some maps
	insertMaps, err := tx.PrepareContext(ctx, `
	insert into maps (id, name, created_at, version) values 
	($1, $2, $3, $4), 
	($5, $6, $7, $8);
	`)

	if err != nil {
		panic(err)
	}

	_, err = insertMaps.Exec(9000, "testMap", time.Now(), "Java", 10000, "secondTestMap", time.Now(), "Bedrock")

	if err != nil {
		panic(err)
	}

	//INserts some coordinates
	insertCoordinates, err := tx.PrepareContext(ctx, `
	insert into coordinates (id, name, xvalue, yvalue, zvalue, mapid, created_at) values 
	($1, $2, $3, $4, $5, $6, $14),
	($7, $8, $3, $4, $5, $6, $14),
	($9, $10, $3, $4, $5, $6, $14),
	($11, $12, $3, $4, $5, $13, $14),
	($14, $15, $3, $4, $5, $13, $14),
	($16, $17, $3, $4, $5, $13, $14)
	`)

	if err != nil {
		panic(err)
	}

	_, err = insertCoordinates.Exec(9000, "fisrt", 200.5, 68, 12.4, 9000, 10000, "second", 11000, "third", 12000, "fourth", 10000, 13000, "fifth", 14000, "sixth", time.Now())

	if err != nil {
		panic(err)
	}

	if err = tx.Commit(); err != nil {
		panic(err)
	}

	testDB = conn
	New(testDB)
}

func shutDownDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	//Deletes all the tables so IDs are not a problem
	clearMaps, err := testDB.PrepareContext(ctx, " drop table coordinates; drop table maps;")

	if err != nil {
		panic(err)
	}

	_, err = clearMaps.Exec()

	if err != nil {
		panic(err)
	}

	testDB.Close()
}

func TestMain(m *testing.M) {
	setupDatabase()
	code := m.Run()
	shutDownDatabase()
	os.Exit(code)
}

func TestCreateMap(t *testing.T) {
	_, err := CreateMap("Nuevas aventuras", "Java")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMaps(t *testing.T) {
	maps, err := GetMaps()

	if err != nil {
		t.Fatal(err)
	}

	if len(maps) < 2 {
		t.Fail()
	}
}

func TestEditMap(t *testing.T) {
	testMap := Map{
		ID: 9000,
	}

	err := testMap.EditMap("Modified name by EditaMap test", "Java")

	if err != nil || testMap.Version != "Java" {
		t.Fail()
	}
}

func TestDeleteMap(t *testing.T) {
	testMap := Map{
		ID: 10000,
	}

	err := testMap.DeleteMap()

	if err != nil {
		t.Fail()
	}

	//CHecks if the coordinates where eliminated to
	err = testMap.GetCoordinates()
	if err != nil || len(testMap.Coordinates) != 0 {
		t.Fail()
	}
}

func TestAddCoordinate(t *testing.T) {
	testMap := Map{
		ID: 9000,
	}

	_, err := testMap.AddCoordinate("Test Coordinate from the AddCoordinate funciton", 200, 64, 1000)

	//Althoug it has 4 coordinates registered total, it only checks if the new has been added to the list
	if err != nil || len(testMap.Coordinates) != 1 {
		fmt.Println(err)
		t.Fail()
	}
}

func TestGetCoordinates(t *testing.T) {
	testMap := Map{
		ID: 9000,
	}

	err := testMap.GetCoordinates()

	//Only checks for 2 of the 3 created in the setup due to tests order execution being random. It could be runned after 1 is added or deleted
	if err != nil || len(testMap.Coordinates) < 2 {
		t.Fail()
	}
}

func TestDeleteCoordinate(t *testing.T) {
	testMap := Map{
		ID: 9000,
	}

	//First Gets all the coordinates. Doesnt check for an error
	testMap.GetCoordinates()
	//Current ammount of coordinates
	currentSize := len(testMap.Coordinates)

	err := testMap.Coordinates[0].DeleteCoordinate()

	//The function doesnt delete the coordinate from the map list, so the list size shouldnt be different
	if err != nil || len(testMap.Coordinates) != currentSize {
		t.Fail()
	}

	testMap.GetCoordinates()

	//After the coordinates are updated the size is now different
	if len(testMap.Coordinates) == currentSize {
		t.Fail()
	}

}

func TestEditCoordinate(t *testing.T) {
	testMap := Map{
		ID: 9000,
	}

	//First Gets all the coordinates. Doesnt check for an error
	testMap.GetCoordinates()

	err := testMap.Coordinates[0].EditCoordinate("Edited name", 100, 100, 100)

	if err != nil || testMap.Coordinates[0].Name != "Edited name" {
		fmt.Print(err)
		t.Fail()
	}
}
