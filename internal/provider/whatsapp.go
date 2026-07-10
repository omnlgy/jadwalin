package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WhatsAppProvider struct {
	baseURL string
	apiKey  string
}

func NewWhatsAppProvider(baseURL, apiKey string) *WhatsAppProvider {
	return &WhatsAppProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

func (p *WhatsAppProvider) SendMessage(ctx context.Context, to, message string) error {
	formatedPhone, _ := strings.CutPrefix(to, "+")

	payload := map[string]interface{}{
		"phone":        formatedPhone + "@s.whatsapp.net",
		"message":      message,
		"is_forwarded": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", p.baseURL+"/send/message", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-device-id", p.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
