package service

import (
	"MedicineBuddy/dto"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type SMTPEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPEmailService(host, port, username, password, from string) *SMTPEmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPEmailService) Send(to, subject, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", s.from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/plain", body)

	// fmt.Println(s.from, to, subject, body)

	port, _ := strconv.Atoi(s.port)

	d := gomail.NewDialer(s.host, port, s.from, s.password)

	return d.DialAndSend(mail)
}

func (s *SMTPEmailService) BuildWeeklyEmail(r dto.WeeklyDetailedReport) string {
	var b strings.Builder

	b.WriteString("Patient: " + r.PatientName + "\n\n")

	for _, day := range r.Days {
		b.WriteString(day.Date + "\n")

		for _, d := range day.Doses {
			line := fmt.Sprintf(
				"%s %s %s → %s",
				d.ScheduledAt.Format("15:04"),
				d.MedicineName,
				d.Dosage,
				strings.ToUpper(d.Status),
			)

			if d.TakenAt != nil {
				line += " at " + d.TakenAt.Format("15:04")
			}

			b.WriteString(line + "\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}
