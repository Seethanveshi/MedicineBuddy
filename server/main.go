package main

import (
	"MedicineBuddy/config"
	"MedicineBuddy/database"
	"MedicineBuddy/jobs"
	"MedicineBuddy/repository"
	"MedicineBuddy/routes"
	"MedicineBuddy/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	db := database.Connect(config.DBUrl)

	doseRepo := repository.NewDoseRepository(db)
	doseService := service.NewDoseService(doseRepo)

	doseCron := jobs.NewDoseCron(doseService)
	doseCron.Start()

	r := gin.Default()
	routes.Router(r, db)
	r.Run(":" + config.AppPort)
}
