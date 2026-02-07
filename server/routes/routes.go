package routes

import (
	"MedicineBuddy/handler"
	"MedicineBuddy/repository"
	"MedicineBuddy/service"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine, db *sql.DB) {
	doseRepository := repository.NewDoseRepository(db)
	doseService := service.NewDoseService(doseRepository)
	doseHandler := handler.NewDoseHandler(doseService)

	medicineRepository := repository.NewMedicineRepository(db)
	medicineService := service.NewMedicineService(medicineRepository, doseService)
	medicineHandler := handler.NewMedicineHandler(medicineService)

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
}
