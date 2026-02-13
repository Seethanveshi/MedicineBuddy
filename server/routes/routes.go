package routes

import (
	"MedicineBuddy/handler"
	"MedicineBuddy/repository"
	"MedicineBuddy/service"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, db *sql.DB) {
	doseRepository := repository.NewDoseRepository(db)
	doseService := service.NewDoseService(doseRepository)
	doseHandler := handler.NewDoseHandler(doseService)

	medicineRepository := repository.NewMedicineRepository(db)
	medicineService := service.NewMedicineService(medicineRepository, doseService)
	medicineHandler := handler.NewMedicineHandler(medicineService)

	mediTakerReposotory := repository.NewMediTakerRepository(db)
	mediTakerService := service.NewMediTakerService(mediTakerReposotory)
	mediTakerHandler := handler.NewMediTakerHandler(mediTakerService)

	reportRepository := repository.NewDoseRepository(db)
	reportService := service.NewReportService(reportRepository)
	emailService := service.NewSMTPEmailService(os.Getenv("SMTP_HOST"), os.Getenv("SMTP_PORT"), os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"), os.Getenv("SMTP_FROM"))
	reportHandler := handler.NewReportHandler(reportService, mediTakerService, emailService)

	r.POST("/meditakers", mediTakerHandler.CreateMediTaker)
	r.GET("/meditakers", mediTakerHandler.List)
	r.DELETE("/meditakers/:id", mediTakerHandler.Delete)
	r.PUT("/meditakers/:id", mediTakerHandler.Update)

	doses := r.Group("/doses")
	{
		doses.GET("/date", doseHandler.GetDosesByDate)
		doses.GET("/today", doseHandler.GetToday)
		doses.GET("/upcoming", doseHandler.GetUpcoming)
		doses.GET("/history", doseHandler.GetHistory)
		doses.POST("/:id/take", doseHandler.TakeDose)
		doses.POST("/:id/skip", doseHandler.SkipDose)
	}

	medicines := r.Group("/medicines")
	{
		medicines.POST("", medicineHandler.CreateMedicine)
	}

	reports := r.Group("/reports")
	{
		reports.GET("/weekly", reportHandler.WeeklyDetailed)
		reports.POST("/send-test", reportHandler.SendTest)
	}
}
