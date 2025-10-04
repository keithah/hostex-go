package hostex

import (
	"context"
	"fmt"
	"strconv"
)

// WebhooksResponse represents the response from listing webhooks
type WebhooksResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

// ListWebhooks retrieves a list of configured webhooks
func (c *Client) ListWebhooks(ctx context.Context) (*WebhooksResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/webhooks", nil, nil)
	if err != nil {
		return nil, err
	}

	var result WebhooksResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal webhooks response: %w", err)
	}

	return &result, nil
}

// CreateWebhookResponse represents the response from creating a webhook
type CreateWebhookResponse struct {
	Webhook Webhook `json:"webhook"`
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(ctx context.Context, webhookURL string) (*CreateWebhookResponse, error) {
	body := map[string]string{
		"url": webhookURL,
	}

	resp, err := c.doRequest(ctx, "POST", "/webhooks", nil, body)
	if err != nil {
		return nil, err
	}

	var result CreateWebhookResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create webhook response: %w", err)
	}

	return &result, nil
}

// DeleteWebhook deletes a webhook by ID
func (c *Client) DeleteWebhook(ctx context.Context, webhookID int) error {
	_, err := c.doRequest(ctx, "DELETE", "/webhooks/"+strconv.Itoa(webhookID), nil, nil)
	return err
}
