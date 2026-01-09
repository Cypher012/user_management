package email

import (
	"log"
	"os"
	"strconv"
)

func EmailConfig() *Service {
	smtpPass := os.Getenv("BREVO_SMTP_PASS")
	port := os.Getenv("BREVO_SMTP_PORT")
	username := os.Getenv("BREVO_SMTP_USERNAME")
	host := os.Getenv("BREVO_SMTP_HOST")

	requiredEnvs := map[string]string{
		"BREVO_SMTP_PASS":     smtpPass,
		"BREVO_SMTP_PORT":     port,
		"BREVO_SMTP_USERNAME": username,
		"BREVO_SMTP_HOST":     host,
	}

	for name, value := range requiredEnvs {
		if value == "" {
			log.Fatalf("environment variable %s is not set", name)
		}
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("BREVO_SMTP_PORT must be a valid number")
	}

	smtpMailer := NewSTMPMailer(host, portInt, username, smtpPass, "My App <no-reply@myapp.com>")

	templates, err := LoadTemplates()
	if err != nil {
		log.Fatal()
	}

	return NewService(smtpMailer, templates)
}
