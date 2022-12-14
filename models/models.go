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

func (m *Map) DeleteMap() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	statement, err := db.PrepareContext(ctx, "delete from maps where id = $1")

	if err != nil {
		return err
	}

	_, err = statement.Exec(m.ID)

	if err != nil {
		return err
	}

	return nil
}

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
