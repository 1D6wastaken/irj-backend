package smtp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"irj/pkg/glog"
)

type (
	SMTPService struct {
		Logger    *glog.Logger
		Host      string
		APIKey    string
		FromName  string
		FromEmail string
	}

	EmailRequest struct {
		Sender      EmailPerson   `json:"sender"`
		To          []EmailPerson `json:"to"`
		Subject     string        `json:"subject"`
		HTMLContent string        `json:"htmlContent"`
		TextContent string        `json:"textContent,omitempty"`
	}

	EmailPerson struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
)

func NewSMTPService(logger *glog.Logger, host, apiKey, fromName, fromEmail string) *SMTPService {
	return &SMTPService{
		Logger:    logger,
		Host:      host,
		APIKey:    apiKey,
		FromName:  fromName,
		FromEmail: fromEmail,
	}
}

func (s *SMTPService) send(ctx context.Context, to []EmailPerson, subject, htmlContent, textContent string) error {
	reqBody := EmailRequest{
		Sender: EmailPerson{
			Name:  s.FromName,
			Email: s.FromEmail,
		},
		To:          to,
		Subject:     subject,
		HTMLContent: htmlContent,
		TextContent: textContent,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		s.Logger.Err(err).Msg("unable to marshal request body")

		return err
	}

	// Création de la requête HTTP
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.Host, bytes.NewBuffer(payload))
	if err != nil {
		s.Logger.Err(err).Msg("unable to create request")

		return err
	}

	// En-têtes requis
	req.Header.Set("Api-Key", s.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		s.Logger.Err(err).Msg("Send email failed")

		return err
	}

	defer resp.Body.Close()

	// Vérification du résultat
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("send email failed: %s", resp.Status)
}
