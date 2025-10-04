package hostex

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ConversationsResponse represents the response from listing conversations
type ConversationsResponse struct {
	Conversations []Conversation `json:"conversations"`
	Total         int            `json:"total"`
}

// ListConversationsParams contains optional parameters for listing conversations
type ListConversationsParams struct {
	Offset int
	Limit  int
}

// ListConversations retrieves a list of conversations
func (c *Client) ListConversations(ctx context.Context, params *ListConversationsParams) (*ConversationsResponse, error) {
	urlParams := url.Values{}

	// offset is required by the API
	offset := 0
	limit := 20

	if params != nil {
		if params.Offset > 0 {
			offset = params.Offset
		}
		if params.Limit > 0 {
			limit = params.Limit
		}
	}

	urlParams.Set("offset", strconv.Itoa(offset))
	urlParams.Set("limit", strconv.Itoa(limit))

	resp, err := c.doRequest(ctx, "GET", "/conversations", urlParams, nil)
	if err != nil {
		return nil, err
	}

	var result ConversationsResponse
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversations response: %w", err)
	}

	return &result, nil
}

// ConversationDetails represents detailed conversation information
type ConversationDetails struct {
	Guest       Guest     `json:"guest"`
	ChannelType string    `json:"channel_type"`
	Messages    []Message `json:"messages"`
}

// GetConversation retrieves detailed information about a specific conversation
func (c *Client) GetConversation(ctx context.Context, conversationID string) (*ConversationDetails, error) {
	resp, err := c.doRequest(ctx, "GET", "/conversations/"+conversationID, nil, nil)
	if err != nil {
		return nil, err
	}

	var result ConversationDetails
	if err := unmarshalData(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conversation details: %w", err)
	}

	return &result, nil
}

// SendMessage sends a message to a conversation
func (c *Client) SendMessage(ctx context.Context, conversationID string, data SendMessageData) error {
	_, err := c.doRequest(ctx, "POST", "/conversations/"+conversationID, nil, data)
	return err
}
