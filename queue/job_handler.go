package queue

import (
	"crypto/tls"

	"github.com/solabsafrica/afrikanest/logger"
	"gopkg.in/gomail.v2"
)

type JobHandler interface {
	HandleSendEmailJob(email string, content string, textContent string) error
}

type jobHanderImpl struct {
}

func (handler *jobHanderImpl) HandleSendEmailJob(email string, content string, textContent string) error {
	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", "info@afrikanest.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)

	// Set E-Mail subject
	m.SetHeader("Subject", "New Charge on your account")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", textContent)

	m.SetBody("text/html", content)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.mailtrap.io", 2525, "58e23d819a6fd8", "3155190a3260ac")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		logger.Info("Error sending email: %v", err)
		return err
	}

	return nil

}

func NewJobHandler() JobHandler {
	return &jobHanderImpl{}
}
