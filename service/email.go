package service

import (
	"context"
	"fmt"

	"github.com/matcornic/hermes/v2"
	"github.com/solabsafrica/afrikanest/config"
	"github.com/solabsafrica/afrikanest/logger"
	"github.com/solabsafrica/afrikanest/model"
	"github.com/solabsafrica/afrikanest/queue"
)

type EmailServiceWithContext func(ctx context.Context) EmailService

type EmailService interface {
	SendEmail(to string) error
	SendLeaseChargePaymentEmail(to model.Tenant, lease model.Lease, ammount int64) error
	SendTenantWelcomeEmail(to model.Tenant) error
}

type emailServiceImpl struct {
	ctx    context.Context
	config *config.Config
	q      queue.JobQueue
}

func NewEmailServiceWithContext(q queue.JobQueue) EmailServiceWithContext {
	config := config.Get()
	return func(ctx context.Context) EmailService {
		return &emailServiceImpl{
			ctx:    ctx,
			config: config,
			q:      q,
		}
	}
}

func (service *emailServiceImpl) SendEmail(to string) error {

	return nil
}

func (service *emailServiceImpl) SendLeaseChargePaymentEmail(to model.Tenant, lease model.Lease, ammount int64) error {

	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "AFRIKANEST",
			Link: "https://afrikanest.com",
			// Optional product logo
			Logo: "https://www.afrikanest.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Flogo.fcf91aaf.svg&w=96&q=75",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: to.Email,
			Intros: []string{
				fmt.Sprintf("You have a new lease charge payment for %s of KES %d", lease.Unit.Name, ammount),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click here to view the charge",
					Button: hermes.Button{
						Color: "#0F0AEE", // Optional action button color
						Text:  "View Charge",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		logger.Error(err)
	}

	// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		logger.Error(err)
	}

	service.q.Publish(service.ctx, queue.NewSendEmailJob(to.Email, emailBody, emailText))

	return nil
}

func (service *emailServiceImpl) SendTenantWelcomeEmail(to model.Tenant) error {

	// Configure hermes by setting a theme and your product info
	h := hermes.Hermes{
		// Optional Theme
		// Theme: new(Default)
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "AFRIKANEST",
			Link: "https://afrikanest.com",
			// Optional product logo
			Logo: "https://www.afrikanest.com/_next/image?url=%2F_next%2Fstatic%2Fmedia%2Flogo.fcf91aaf.svg&w=96&q=75",
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: to.FirstName,
			Intros: []string{
				fmt.Sprintf("Welcome to AFRIKANEST, %s", to.FirstName),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click here to setup your account",
					Button: hermes.Button{
						Color: "#0F0AEE", // Optional action button color
						Text:  "Setup Account",
						Link:  "https://hermes-example.com/confirm?token=d9729feb74992cc3482b350163a1a010",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		logger.Error(err)
	}

	emailText, err := h.GeneratePlainText(email)
	if err != nil {
		logger.Error(err)
	}

	service.q.Publish(service.ctx, queue.NewSendEmailJob(to.Email, emailBody, emailText))
	return nil
}
