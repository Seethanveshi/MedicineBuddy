package main

import (
	"MedicineBuddy/config"
	"MedicineBuddy/database"
	"MedicineBuddy/jobs"
	"MedicineBuddy/repository"
	"MedicineBuddy/routes"
	"MedicineBuddy/service"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config := config.Load()
	db := database.Connect(config.DBUrl)

	doseRepo := repository.NewDoseRepository(db)
	doseService := service.NewDoseService(doseRepo)

	doseCron := jobs.NewDoseCron(doseService)
	doseCron.Start()

	reportRepository := repository.NewDoseRepository(db)
	reportService := service.NewReportService(reportRepository)
	mediTakerReposotory := repository.NewMediTakerRepository(db)
	mediTakerService := service.NewMediTakerService(mediTakerReposotory)
	emailService := service.NewSMTPEmailService(os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"), os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_FROM"))
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)

	weeklyCron := jobs.NewWeeklyReportJob(reportService, mediTakerService, emailService, userService)
	weeklyCron.Start()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8081"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	routes.Router(r, db)
	r.Run(":" + config.AppPort)
}
