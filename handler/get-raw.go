package handler

import (
	"challenge/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RawGet(c *gin.Context) {
	ticker := c.Param("ticker")

	// endpoint1 := make([]Endpoint1, 0)
	var total float64

	database.DB.Raw("SELECT SUM(total) FROM endpoint1 WHERE ticker = ?", ticker).Scan(&total)

	// if endpoint1.ID == 0 {
	// 	c.JSON(http.StatusNotFound, endpoint1)
	// 	return
	// }

	c.JSON(http.StatusOK, total)
}
