package hostex

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ReviewsResponse represents the response from listing reviews
type ReviewsResponse struct {
	Reviews []Review `json:"reviews"`
	Total   int      `json:"total"`
}

// ListReviewsParams contains optional parameters for listing reviews
type ListReviewsParams struct {
	ReservationCode    string
	PropertyID         int
	ReviewStatus       string
	StartCheckOutDate  string
	EndCheckOutDate    string
	Offset             int
	Limit              int
}

// ListReviews retrieves a list of reviews
func (c *Client) ListReviews(ctx context.Context, params *ListReviewsParams) (*ReviewsResponse, error) {
	urlParams := url.Values{}

	if params != nil {
		if params.ReservationCode != "" {
			urlParams.Set("reservation_code", params.ReservationCode)
		}
		if params.PropertyID > 0 {
			urlParams.Set("property_id", strconv.Itoa(params.PropertyID))
		}
		if params.ReviewStatus != "" {
			urlParams.Set("review_status", params.ReviewStatus)
		}
		if params.StartCheckOutDate != "" {
			urlParams.Set("start_check_out_date", params.StartCheckOutDate)
		}
		if params.EndCheckOutDate != "" {
			urlParams.Set("end_check_out_date", params.EndCheckOutDate)
		}
		if params.Offset > 0 {
			urlParams.Set("offset", strconv.Itoa(params.Offset))
		}
		if params.Limit > 0 {
			urlParams.Set("limit", strconv.Itoa(params.Limit))
		}
	}

	resp, err := c.doRequest(ctx, "GET", "/reviews", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result ReviewsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reviews response: %w", err)
	}

	return &result, nil
}

// CreateReview creates a review or reply for a reservation
func (c *Client) CreateReview(ctx context.Context, reservationCode string, data CreateReviewData) error {
	_, err := c.doRequest(ctx, "POST", "/reviews/"+reservationCode, nil, data)
	return err
}
