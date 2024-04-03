package main

import (
	"challenge/database"
	"challenge/handler"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {
	// db := database.GetConnection()
	// if err := database.MakeMigrations(db); err != nil {
	// 	panic(err)
	// }
	database.GetConnection()

	fmt.Println("CHALLENGE DB CONNECTED")

	err := database.DB.AutoMigrate(&handler.Endpoint1{}, &handler.File{})
	if err != nil {
		panic("failed to migrate database")
	}

	r := gin.Default()

	r.POST("/endpoint1", handler.HandleEndpoint1)
	r.POST("/endpoint2", handler.HandleEndpoint2)
	r.GET("/by-id/:id", handler.GetByID)
	r.GET("/by-ticker/:ticker", handler.RawGet)

	port := ":" + goDotEnvVariable("PORT")

	if port == ":" {
		port = ":8080"
	}

	r.Run(port)
}
