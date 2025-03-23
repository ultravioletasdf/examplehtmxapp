package utils

import (
	"net/mail"

	"github.com/matcornic/hermes/v2"
	"github.com/valord577/mailx"
)

var h = hermes.Hermes{
	Product: hermes.Product{
		Name: "Example HTMX App",
		Link: "https://example.com",
		Logo: "https://htmx.org/img/topo.svg",
	},
}

var dialer *mailx.Dialer

func CreateEmailClient(cfg *Config) {
	dialer = &mailx.Dialer{
		Host:         cfg.SmtpHost,
		Port:         cfg.SmtpPort,
		Username:     cfg.SmtpUsername,
		Password:     cfg.SmtpPassword,
		SSLOnConnect: false,
	}
}
func SendVerification(cfg *Config, address string, code string) error {
	email := hermes.Email{
		Body: hermes.Body{
			Name:   address,
			Intros: []string{"Welcome to exampleHTMXapp! Please verify your email with the code below to access the app."},
			Actions: []hermes.Action{
				hermes.Action{Instructions: "Please enter this code in the website:", InviteCode: code},
			},
		},
	}
	html, err := h.GenerateHTML(email)
	if err != nil {
		return err
	}

	m := mailx.NewMessage()
	m.SetTo(address)
	m.SetFrom(&mail.Address{Name: "Runik", Address: cfg.SmtpAddress})
	m.SetSubject("Verify your email for Runik")
	m.SetHtmlBody(html)

	return dialer.DialAndSend(m)
}
