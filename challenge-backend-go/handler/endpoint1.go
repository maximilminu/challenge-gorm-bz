package handler

import (
	"challenge/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Endpoint1 struct {
	gorm.Model
	ID         uint      `gorm:"primaryKey"`
	Ticker     string    `json:"ticker"`
	Code       int32     `json:"code" binding:"gte=0"`
	Total      float64   `json:"total" binding:"gte=0"`
	Exceeded   bool      `json:"exceeded"`
	LastChange time.Time `json:"last_change"`
}

func HandleEndpoint1(c *gin.Context) {

	type Payload struct {
		Ticker     string  `json:"ticker"`
		Code       int32   `json:"code" binding:"gte=1"`
		Total      float64 `json:"total" binding:"gte=0"`
		Exceeded   bool    `json:"exceeded"`
		LastChange string  `json:"last_change"`
	}

	payload := make([]Payload, 0)
	response := make([]Endpoint1, 0)

	err := c.ShouldBindJSON(&payload)
	if err != nil {
		fmt.Printf("ERROR BINDING ENDPOINT 1 %+v", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"BINDING_ERROR": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		fmt.Printf("BEGIN TRANSACTION ERROR %+v", tx.Error.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"TRANSACTION_ERROR": tx.Error.Error()})
		return
	}

	for _, e := range payload {
		endpoint1 := Endpoint1{}

		layout := "2006-01-02 15:04:05"
		lastChange, err := time.Parse(layout, e.LastChange)

		if err != nil {
			fmt.Println(err)
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"PARSING_ERROR": err.Error()})
			return
		}

		endpoint1 = Endpoint1{
			//ID:         uuid.New(),
			Ticker:     e.Ticker,
			Code:       e.Code,
			Total:      e.Total,
			Exceeded:   e.Exceeded,
			LastChange: lastChange,
		}

		insert := tx.Create(&endpoint1)

		if insert.Error != nil {
			fmt.Println(insert.Error)
			c.JSON(http.StatusBadRequest, gin.H{"INSERT_ERROR": insert.Error})
			tx.Rollback()
			return
		}
		log.Printf("Endpoint1 created %v", insert.RowsAffected)

		response = append(response, endpoint1)
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}
