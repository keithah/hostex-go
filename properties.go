package hostex

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// PropertiesResponse represents the response from listing properties
type PropertiesResponse struct {
	Properties []Property `json:"properties"`
	Total      int        `json:"total"`
}

// ListPropertiesParams contains optional parameters for listing properties
type ListPropertiesParams struct {
	Offset int
	Limit  int
	ID     int
}

// ListProperties retrieves a list of properties
func (c *Client) ListProperties(ctx context.Context, params *ListPropertiesParams) (*PropertiesResponse, error) {
	urlParams := url.Values{}

	if params != nil {
		if params.Offset > 0 {
			urlParams.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Limit > 0 {
			urlParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.ID > 0 {
			urlParams.Set("id", strconv.Itoa(params.ID))
		}
	}

	resp, err := c.doRequest(ctx, "GET", "/properties", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result PropertiesResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal properties response: %w", err)
	}

	return &result, nil
}

// RoomTypesResponse represents the response from listing room types
type RoomTypesResponse struct {
	RoomTypes []RoomType `json:"room_types"`
	Total     int        `json:"total"`
}

// ListRoomTypesParams contains optional parameters for listing room types
type ListRoomTypesParams struct {
	Offset int
	Limit  int
}

// ListRoomTypes retrieves a list of room types
func (c *Client) ListRoomTypes(ctx context.Context, params *ListRoomTypesParams) (*RoomTypesResponse, error) {
	urlParams := url.Values{}

	if params != nil {
		if params.Offset > 0 {
			urlParams.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Limit > 0 {
			urlParams.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	resp, err := c.doRequest(ctx, "GET", "/room_types", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result RoomTypesResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal room types response: %w", err)
	}

	return &result, nil
}
