package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"github.com/kushalktamang/blueprint-fullstack/internal/config"
	"github.com/pkg/errors"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

//go:embed templates/emails/*.html
var templatesFS embed.FS

type Client struct {
	client    *resend.Client
	logger    *zerolog.Logger
	templates map[Template]*template.Template
	fromName  string
	fromAddr  string
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) (*Client, error) {
	templates, err := loadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return &Client{
		client:    resend.NewClient(cfg.Integration.ResendAPIKey),
		logger:    logger,
		templates: templates,
		fromName:  cfg.Integration.EmailFromName,
		fromAddr:  cfg.Integration.EmailFromAddress,
	}, nil
}

func loadTemplates() (map[Template]*template.Template, error) {
	entries, err := templatesFS.ReadDir("templates/emails")
	if err != nil {
		return nil, fmt.Errorf("reading embedded templates dir: %w", err)
	}

	templates := make(map[Template]*template.Template, len(entries))
	for _, entry := range entries {
		name := entry.Name() // e.g. "welcome.html"
		path := "templates/emails/" + name

		tmpl, err := template.ParseFS(templatesFS, path)
		if err != nil {
			return nil, fmt.Errorf("parsing template %s: %w", name, err)
		}

		key := Template(name[:len(name)-len(".html")])
		templates[key] = tmpl
	}
	return templates, nil
}

func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName)

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse email template %s", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", "Blueprint", "onboarding@resend.dev"),
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}

	_, err = c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
