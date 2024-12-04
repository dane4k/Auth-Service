package main

import (
	"AuthService/db"
	"AuthService/internal/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Println(err)
	}

	db.InitDB()

	router := gin.Default()

	routes.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
