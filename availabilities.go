package hostex

import (
	"context"
	"fmt"
	"net/url"
)

// AvailabilitiesResponse represents the response from listing availabilities
type AvailabilitiesResponse struct {
	Listings []struct {
		ID              int            `json:"id"`
		ChannelType     string         `json:"channel_type"`
		ListingID       string         `json:"listing_id"`
		Availabilities  []Availability `json:"availabilities,omitempty"`
	} `json:"listings"`
}

// ListAvailabilitiesParams contains required parameters for listing availabilities
type ListAvailabilitiesParams struct {
	PropertyIDs string // Comma-separated property IDs
	StartDate   string // YYYY-MM-DD
	EndDate     string // YYYY-MM-DD
}

// ListAvailabilities retrieves availability information for properties
func (c *Client) ListAvailabilities(ctx context.Context, params ListAvailabilitiesParams) (*AvailabilitiesResponse, error) {
	urlParams := url.Values{}
	urlParams.Set("property_ids", params.PropertyIDs)
	urlParams.Set("start_date", params.StartDate)
	urlParams.Set("end_date", params.EndDate)

	resp, err := c.doRequest(ctx, "GET", "/availabilities", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result AvailabilitiesResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal availabilities response: %w", err)
	}

	return &result, nil
}

// UpdateAvailabilities updates property availability status
func (c *Client) UpdateAvailabilities(ctx context.Context, data UpdateAvailabilitiesData) error {
	_, err := c.doRequest(ctx, "POST", "/availabilities", nil, data)
	return err
}
