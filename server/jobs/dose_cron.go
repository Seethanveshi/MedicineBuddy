package jobs

import (
	"MedicineBuddy/service"
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type DoseCron struct {
	doseService *service.DoseService
}

func NewDoseCron(doseService *service.DoseService) *DoseCron {
	return &DoseCron{doseService: doseService}
}

func (d *DoseCron) Start() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("* */5 * * * *", func() {
		ctx := context.Background()

		if err := d.doseService.MarkMissedDose(ctx); err != nil {
			log.Println("failed to mark missed doses:", err)
		}
	})
	c.Start()
}
