package jobs

import (
	"MedicineBuddy/service"
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type WeeklyReportJob struct {
	reportService    *service.ReportService
	mediTakerService *service.MediTakerService
	emailService     *service.SMTPEmailService
	userService      *service.UserService
}

func NewWeeklyReportJob(rs *service.ReportService, ms *service.MediTakerService, es *service.SMTPEmailService, us *service.UserService) *WeeklyReportJob {
	return &WeeklyReportJob{
		reportService:    rs,
		mediTakerService: ms,
		emailService:     es,
		userService:      us,
	}
}

func (wrj *WeeklyReportJob) Start() {
	c := cron.New(cron.WithSeconds())

	c.AddFunc("0 59 23 * * 0", func() {
		wrj.RunWeeklyReportJob()
	})

	c.Start()
}

func (wrj *WeeklyReportJob) RunWeeklyReportJob() {
	ctx := context.Background()

	patients, err := wrj.userService.ListOfUsers(ctx)
	if err != nil {
		log.Println("Failed to fetch patients", err)
		return
	}

	for _, patient := range patients {
		report, err := wrj.reportService.WeeklyDetailed(ctx, patient.ID)
		if err != nil {
			log.Println("Failed to fetch weekly Details", err)
			continue
		}

		meditakers, err := wrj.mediTakerService.ListByPatient(ctx, patient.ID)
		if err != nil {
			log.Println("Failed to fetch meditakers", err)
			continue
		}

		body := wrj.emailService.BuildWeeklyEmail(report)

		for _, medimeditaker := range meditakers {
			err := wrj.emailService.Send(
				medimeditaker.Email,
				"Weekly medicine report",
				body,
			)

			if err != nil {
				log.Println("Failed to send Email", err)
			}
		}
	}
}
