package hostex

import (
	"context"
	"fmt"
)

// ListingCalendarResponse represents the response from getting listing calendar
type ListingCalendarResponse struct {
	Listings []struct {
		ChannelType string `json:"channel_type"`
		ListingID   string `json:"listing_id"`
		Calendar    []struct {
			Date               string `json:"date"`
			Price              int    `json:"price,omitempty"`
			Inventory          int    `json:"inventory,omitempty"`
			Available          bool   `json:"available,omitempty"`
			MinStay            int    `json:"min_stay,omitempty"`
			MaxStay            int    `json:"max_stay,omitempty"`
			ClosedToArrival    bool   `json:"closed_to_arrival,omitempty"`
			ClosedToDeparture  bool   `json:"closed_to_departure,omitempty"`
		} `json:"calendar"`
	} `json:"listings"`
}

// GetListingCalendar retrieves calendar information for multiple listings
func (c *Client) GetListingCalendar(ctx context.Context, data GetListingCalendarData) (*ListingCalendarResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/listings/calendar", nil, data)
	if err != nil {
		return nil, err
	}

	var result ListingCalendarResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal listing calendar response: %w", err)
	}

	return &result, nil
}

// UpdateListingPrices updates listing prices for channel listings
func (c *Client) UpdateListingPrices(ctx context.Context, data UpdateListingPricesData) error {
	_, err := c.doRequest(ctx, "POST", "/listings/prices", nil, data)
	return err
}

// UpdateListingInventories updates inventory levels for channel listings
func (c *Client) UpdateListingInventories(ctx context.Context, data UpdateListingInventoriesData) error {
	_, err := c.doRequest(ctx, "POST", "/listings/inventories", nil, data)
	return err
}

// UpdateListingRestrictions updates listing restrictions for channel listings
func (c *Client) UpdateListingRestrictions(ctx context.Context, data UpdateListingRestrictionsData) error {
	_, err := c.doRequest(ctx, "POST", "/listings/restrictions", nil, data)
	return err
}
