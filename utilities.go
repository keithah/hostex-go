package hostex

import (
	"context"
	"fmt"
)

// CustomChannelsResponse represents the response from listing custom channels
type CustomChannelsResponse struct {
	CustomChannels []CustomChannel `json:"custom_channels"`
}

// ListCustomChannels retrieves custom channels from Custom Options Page
func (c *Client) ListCustomChannels(ctx context.Context) (*CustomChannelsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/custom_channels", nil, nil)
	if err != nil {
		return nil, err
	}

	var result CustomChannelsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom channels response: %w", err)
	}

	return &result, nil
}

// IncomeMethodsResponse represents the response from listing income methods
type IncomeMethodsResponse struct {
	IncomeMethods []IncomeMethod `json:"income_methods"`
}

// ListIncomeMethods retrieves income methods from Custom Options Page
func (c *Client) ListIncomeMethods(ctx context.Context) (*IncomeMethodsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/income_methods", nil, nil)
	if err != nil {
		return nil, err
	}

	var result IncomeMethodsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal income methods response: %w", err)
	}

	return &result, nil
}
