package main

import (
	"database/sql"

	"github.com/TheSgtPepper23/coordinator/handlers"
	"github.com/TheSgtPepper23/coordinator/models"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

type Config struct {
	DB     *sql.DB
	Models models.Models
}

// Custom middleware to enable cors
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	//Creates database connection and initialize the models package
	conn, err := connectToDB()
	if err != nil {
		panic(err)
	}
	models.New(conn)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	router.GET("/ping", handlers.Ping)

	maps := router.Group("/map")
	{
		maps.GET("/", handlers.GetMaps)
		maps.POST("/", handlers.CreateMap)
		maps.DELETE("/:mapid", handlers.DeleteMap)
		maps.PUT("/:mapid", handlers.EditMap)
		maps.PUT("/addCoordinate/:mapid", handlers.AddCooridnate)
		maps.GET("coordinates/:mapid", handlers.GetCoordinates)
	}

	coordinates := router.Group("/coordinate")
	{
		coordinates.DELETE("/:coordid", handlers.DeleteCoordinate)
		coordinates.PUT("/:coordid", handlers.EditCoordinate)
	}

	router.Run(":50600")

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
