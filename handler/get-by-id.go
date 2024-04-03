package handler

import (
	"challenge/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetByID(c *gin.Context) {
	ID := c.Param("id")

	endpoint1 := Endpoint1{}

	database.DB.First(&endpoint1, ID)

	if endpoint1.ID == 0 {
		c.JSON(http.StatusNotFound, endpoint1)
		return
	}

	c.JSON(http.StatusOK, endpoint1)
}
