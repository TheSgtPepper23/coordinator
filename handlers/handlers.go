package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/TheSgtPepper23/coordinator/models"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Body    any    `json:"body,omitempty"`
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusBadRequest, Response{Error: false, Message: "ping"})
}

func GetMaps(c *gin.Context) {
	mapList, err := models.GetMaps()

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
	} else {
		c.JSON(http.StatusOK, Response{Error: false, Body: mapList})
	}
}

func CreateMap(c *gin.Context) {
	var newMap models.Map

	if err := c.ShouldBindJSON(&newMap); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	lid, err := models.CreateMap(newMap.Name, newMap.Version)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: fmt.Sprint(lid)})

}

func DeleteMap(c *gin.Context) {
	textId := c.Param("mapid")
	mapId, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	tempMap := models.Map{
		ID: mapId,
	}

	err = tempMap.DeleteMap()

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: "Map deleted"})
}

func EditMap(c *gin.Context) {
	textId := c.Param("mapid")
	mapId, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	currentMap := models.Map{
		ID: mapId,
	}

	var newMap models.Map

	if err := c.ShouldBindJSON(&newMap); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	err = currentMap.EditMap(newMap.Name, newMap.Version)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: "Map successfully modified"})
}

func AddCooridnate(c *gin.Context) {
	textId := c.Param("mapid")
	mapId, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	currentMap := models.Map{
		ID: mapId,
	}

	var newCoordinate models.Coordinate

	if err := c.ShouldBindJSON(&newCoordinate); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	_, err = currentMap.AddCoordinate(newCoordinate.Name, newCoordinate.XValue, newCoordinate.YValue, newCoordinate.ZValue)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: "New coordinate added"})
}

func GetCoordinates(c *gin.Context) {
	textId := c.Param("mapid")
	mapId, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	currentMap := models.Map{
		ID: mapId,
	}

	err = currentMap.GetCoordinates()

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Body: currentMap.Coordinates})
}

func EditCoordinate(c *gin.Context) {
	textId := c.Param("coordid")
	coordid, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	currentCoord := models.Coordinate{
		ID: coordid,
	}

	var newCoordinate models.Coordinate

	if err := c.ShouldBindJSON(&newCoordinate); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	err = currentCoord.EditCoordinate(newCoordinate.Name, newCoordinate.XValue, newCoordinate.YValue, newCoordinate.ZValue)

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: "Coordinate successfully modified"})
}

func DeleteCoordinate(c *gin.Context) {
	textId := c.Param("coordid")
	coordid, err := strconv.ParseInt(textId, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: true, Message: err.Error()})
		return
	}

	currentCoord := models.Coordinate{
		ID: coordid,
	}

	err = currentCoord.DeleteCoordinate()

	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{Error: true, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Error: false, Message: "Coordinate deleted"})
}
