package handlers

import (
	"net/http"
	"wmenjoy.com/iptv/models"

	"github.com/labstack/echo"
)

func GetGroups(channels []models.Channel) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch tasks using our new model
		return c.JSON(http.StatusOK, models.GetGroups(channels))
	}
}
func GetChannels(channels []models.Channel) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Fetch tasks using our new model
		return c.JSON(http.StatusOK, models.GetChannels(channels, c.Param("id")))
	}
}
