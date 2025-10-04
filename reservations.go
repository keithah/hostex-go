package hostex

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ReservationsResponse represents the response from listing reservations
type ReservationsResponse struct {
	Reservations []Reservation `json:"reservations"`
	Total        int           `json:"total"`
}

// ListReservationsParams contains optional parameters for listing reservations
type ListReservationsParams struct {
	ReservationCode    string
	PropertyID         int
	Status             string
	StartCheckInDate   string
	EndCheckInDate     string
	StartCheckOutDate  string
	EndCheckOutDate    string
	OrderBy            string
	Offset             int
	Limit              int
}

// ListReservations retrieves a list of reservations
func (c *Client) ListReservations(ctx context.Context, params *ListReservationsParams) (*ReservationsResponse, error) {
	urlParams := url.Values{}

	if params != nil {
		if params.ReservationCode != "" {
			urlParams.Set("reservation_code", params.ReservationCode)
		}
		if params.PropertyID > 0 {
			urlParams.Set("property_id", strconv.Itoa(params.PropertyID))
		}
		if params.Status != "" {
			urlParams.Set("status", params.Status)
		}
		if params.StartCheckInDate != "" {
			urlParams.Set("start_check_in_date", params.StartCheckInDate)
		}
		if params.EndCheckInDate != "" {
			urlParams.Set("end_check_in_date", params.EndCheckInDate)
		}
		if params.StartCheckOutDate != "" {
			urlParams.Set("start_check_out_date", params.StartCheckOutDate)
		}
		if params.EndCheckOutDate != "" {
			urlParams.Set("end_check_out_date", params.EndCheckOutDate)
		}
		if params.OrderBy != "" {
			urlParams.Set("order_by", params.OrderBy)
		}
		if params.Offset > 0 {
			urlParams.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Limit > 0 {
			urlParams.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	resp, err := c.doRequest(ctx, "GET", "/reservations", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result ReservationsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reservations response: %w", err)
	}

	return &result, nil
}

// CreateReservationResponse represents the response from creating a reservation
type CreateReservationResponse struct {
	Reservation Reservation `json:"reservation"`
}

// CreateReservation creates a new direct booking reservation
func (c *Client) CreateReservation(ctx context.Context, data CreateReservationData) (*CreateReservationResponse, error) {
	resp, err := c.doRequest(ctx, "POST", "/reservations", nil, data)
	if err != nil {
		return nil, err
	}

	var result CreateReservationResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create reservation response: %w", err)
	}

	return &result, nil
}

// CancelReservation cancels a direct booking reservation
func (c *Client) CancelReservation(ctx context.Context, reservationCode string) error {
	_, err := c.doRequest(ctx, "DELETE", "/reservations/"+reservationCode, nil, nil)
	return err
}

// UpdateLockCode updates the lock code for a stay
func (c *Client) UpdateLockCode(ctx context.Context, stayCode, lockCode string) error {
	body := map[string]string{
		"lock_code": lockCode,
	}
	_, err := c.doRequest(ctx, "PATCH", "/reservations/"+stayCode+"/check_in_details", nil, body)
	return err
}

// CustomFieldsResponse represents the response from getting custom fields
type CustomFieldsResponse struct {
	CustomFields map[string]interface{} `json:"custom_fields"`
}

// GetCustomFields retrieves custom fields for a stay
func (c *Client) GetCustomFields(ctx context.Context, stayCode string) (*CustomFieldsResponse, error) {
	resp, err := c.doRequest(ctx, "GET", "/reservations/"+stayCode+"/custom_fields", nil, nil)
	if err != nil {
		return nil, err
	}

	var result CustomFieldsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal custom fields response: %w", err)
	}

	return &result, nil
}

// UpdateCustomFields updates custom fields for a stay
func (c *Client) UpdateCustomFields(ctx context.Context, stayCode string, customFields map[string]interface{}) error {
	body := map[string]interface{}{
		"custom_fields": customFields,
	}
	_, err := c.doRequest(ctx, "PATCH", "/reservations/"+stayCode+"/custom_fields", nil, body)
	return err
}
