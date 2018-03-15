package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"text/template"

	"github.com/domodwyer/mailyak"
)

type emailConfig struct {
	URL              string
	SMTPHost         string
	SMTPPort         string
	SMTPUser         string
	SMTPPassword     string
	MailTemplatePath string
}

type emailData struct {
	User        User
	Subject     string
	Body        string
	EmailConfig emailConfig
}

func (ed *emailData) setBody(content string) {
	ed.Body = content
}

func extractEmailConfigFromConfig(config Config) emailConfig {
	return emailConfig{
		URL:              config.URL,
		MailTemplatePath: config.MailTemplatePath,
		SMTPHost:         config.SMTPHost,
		SMTPPort:         config.SMTPPort,
		SMTPUser:         config.SMTPUser,
		SMTPPassword:     config.SMTPPassword,
	}
}

func loadTemplate(name string) (tmpl *template.Template, err error) {

	tl := filepath.Join(
		config.MailTemplatePath,
		fmt.Sprintf("%s.txt", name),
	)

	Info.Println("loading mail template", tl)

	raw, err := ioutil.ReadFile(tl)
	if err != nil {
		return tmpl, err
	}

	tmpl = template.New(name)
	template.Must(tmpl.Parse(string(raw)))

	return tmpl, err
}

func sendEmailConfirmation(user User) error {
	var ed emailData
	var body bytes.Buffer

	tmpl, err := loadTemplate("confirmation")
	if err != nil {
		return err
	}

	// Build emailData and execute the template
	ed = emailData{
		User:        user,
		Subject:     "Welcome to Graphia CMS",
		EmailConfig: extractEmailConfigFromConfig(config),
	}

	err = tmpl.ExecuteTemplate(&body, "confirmation", ed)
	if err != nil {
		return err
	}

	ed.setBody(body.String())

	return mailer.send(ed)
}

// Mailer is the default method of sending mails
type Mailer struct {
	send func(ed emailData) error
}

// DefaultSender is Mailer's default send allowing testing to be split
func DefaultSender(ed emailData) error {

	my := mailyak.New(
		fmt.Sprintf("%s:%s", ed.EmailConfig.SMTPHost, ed.EmailConfig.SMTPPort),
		smtp.PlainAuth(
			"local.smtp",
			ed.EmailConfig.SMTPUser,
			ed.EmailConfig.SMTPPassword,
			ed.EmailConfig.SMTPPort,
		),
	)

	my.To(ed.User.Email)

	my.From("noreply@graphia.co.uk")
	my.FromName("Graphia CMS")
	my.Subject(ed.Subject)

	my.Plain().Set(ed.Body)

	return my.Send()

}
