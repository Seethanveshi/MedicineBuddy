package main

import (
	"MedicineBuddy/config"
	"MedicineBuddy/database"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	database.Connect(config.DBUrl)

	r := gin.Default()

	r.Run(":" + config.AppPort)
}
