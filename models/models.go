package models

import (
	"context"
	"database/sql"
	"time"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) {
	db = dbPool
}

type Models struct {
	Map        Map
	Coordinate Coordinate
}

type Map struct {
	ID           int64
	Name         string
	CreationDate time.Time
	Version      string
	Coordinates  []*Coordinate
}

type Coordinate struct {
	ID     int64
	Name   string
	XValue float64
	YValue float64
	ZValue float64
}

// Creates a new map and stores it in the database Returns the inserted ID or an error
func CreateMap(name, version string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	statement, err := db.PrepareContext(ctx, "insert into maps (name, created_at, version) values ($1, $2, $3)")

	if err != nil {
		return -1, err
	}

	res, err := statement.Exec(name, time.Now(), version)

	if err != nil {
		return -1, err
	}

	lid, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	return lid, nil

}

// Returns all the maps stored in the database (doesn't include the locations)
func GetMaps() ([]*Map, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, "select id, name, version, created_at from maps")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	maps := []*Map{}

	for rows.Next() {
		var tempMap Map
		err := rows.Scan(&tempMap.ID,
			&tempMap.Name,
			&tempMap.Version,
			&tempMap.CreationDate,
		)

		if err != nil {
			return nil, err
		}

		maps = append(maps, &tempMap)
	}

	return maps, nil
}

// Modify the name or the version of a stored map
func (m *Map) EditMap(name, version string) error {
	m.Name = name
	m.Version = version

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	statement, err := db.PrepareContext(ctx, "update maps set name = $1, version = $2 where id = $3")

	if err != nil {
		return err
	}

	_, err = statement.Exec(name, version, m.ID)

	if err != nil {
		return err
	}

	return nil
}

// Deletes the map, and its locations from the database. This process is a transaction
func (m *Map) DeleteMap() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	coordsDelete, err := tx.PrepareContext(ctx, "delete from coordinates where mapid = $1")

	if err != nil {
		return err
	}

	_, err = coordsDelete.Exec(m.ID)

	if err != nil {
		return err
	}

	mapDelete, err := tx.PrepareContext(ctx, "delete from maps where id = $1")

	if err != nil {
		return err
	}

	_, err = mapDelete.Exec(m.ID)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Creates a new coordinate and stores it in the database related to the current map
func (m *Map) AddCoordinate(name string, xValue, yValue, zValue float64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	statement, err := db.PrepareContext(ctx, "insert into coordinates (name, created_at, xvalue, yvalue, zvalue, mapid) values ($1, $2, $3, $4, $5, $6)")

	if err != nil {
		return -1, err
	}

	res, err := statement.Exec(name, time.Now(), xValue, yValue, zValue, m.ID)

	if err != nil {
		return -1, err
	}

	lid, err := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	newCoord := Coordinate{
		ID:     lid,
		Name:   name,
		XValue: xValue,
		YValue: yValue,
		ZValue: zValue,
	}

	m.Coordinates = append(m.Coordinates, &newCoord)

	return lid, nil
}

// Gets all the coordinates for the map and set its Coordinates property with the results
func (m *Map) GetCoordinates() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	m.Coordinates = nil

	rows, err := db.QueryContext(ctx, "select id, name, xvalue, yvalue, zvalue from coordinates where mapid = $1 order by created_at desc", m.ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	coords := []*Coordinate{}

	for rows.Next() {
		var tempCoord Coordinate

		err := rows.Scan(
			&tempCoord.ID,
			&tempCoord.Name,
			&tempCoord.XValue,
			&tempCoord.YValue,
			&tempCoord.ZValue,
		)

		if err != nil {
			return err
		}

		coords = append(coords, &tempCoord)
	}

	m.Coordinates = coords
	return nil

}

// Deletes the coordinate from the database but not from the map list. GetCoordinates should be called to update the map list of coordinates
func (c *Coordinate) DeleteCoordinate() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	statement, err := db.PrepareContext(ctx, "delete from coordinates where id = $1")

	if err != nil {
		return err
	}

	_, err = statement.Exec(c.ID)

	if err != nil {
		return err
	}

	return nil
}

// Modifies the values of the coordinate object and in the database
func (c *Coordinate) EditCoordinate(name string, xValue, yValue, zValue float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	c.Name = name
	c.XValue = xValue
	c.YValue = yValue
	c.ZValue = zValue

	statement, err := db.PrepareContext(ctx, "update coordinates set name = $1, xvalue = $2, yvalue = $3, zvalue = $4 where id = $5")

	if err != nil {
		return err
	}

	_, err = statement.Exec(name, xValue, yValue, zValue, c.ID)

	if err != nil {
		return err
	}

	return nil
}
