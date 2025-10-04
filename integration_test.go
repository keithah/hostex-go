package hostex_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/keithah/hostex-go"
)

// Integration tests - only run when HOSTEX_API_KEY is set
func getTestClient(t *testing.T) *hostex.Client {
	apiKey := os.Getenv("HOSTEX_API_KEY")
	if apiKey == "" {
		t.Skip("HOSTEX_API_KEY not set, skipping integration tests")
	}

	client, err := hostex.NewClient(hostex.Config{
		AccessToken: apiKey,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

// ============================================================
// PROPERTIES TESTS
// ============================================================

func TestIntegration_ListProperties(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list all properties", func(t *testing.T) {
		resp, err := client.ListProperties(ctx, nil)
		if err != nil {
			t.Fatalf("ListProperties failed: %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		t.Logf("Found %d properties (total: %d)", len(resp.Properties), resp.Total)

		for _, prop := range resp.Properties {
			t.Logf("  - %s (ID: %d)", prop.Title, prop.ID)
			if len(prop.Channels) > 0 {
				for _, ch := range prop.Channels {
					t.Logf("    Channel: %s (%s)", ch.ChannelType, ch.ListingID)
				}
			}
		}
	})

	t.Run("list properties with limit", func(t *testing.T) {
		resp, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{
			Limit: 1,
		})
		if err != nil {
			t.Fatalf("ListProperties with limit failed: %v", err)
		}

		if len(resp.Properties) > 1 {
			t.Errorf("Expected max 1 property, got %d", len(resp.Properties))
		}
	})

	t.Run("list properties with offset", func(t *testing.T) {
		resp, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{
			Offset: 1,
			Limit:  1,
		})
		if err != nil {
			t.Fatalf("ListProperties with offset failed: %v", err)
		}

		t.Logf("Got %d properties with offset=1", len(resp.Properties))
	})

	t.Run("get specific property by ID", func(t *testing.T) {
		// First get a property ID
		allProps, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{Limit: 1})
		if err != nil || len(allProps.Properties) == 0 {
			t.Skip("No properties available for testing")
		}

		propertyID := allProps.Properties[0].ID

		resp, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{
			ID: propertyID,
		})
		if err != nil {
			t.Fatalf("ListProperties by ID failed: %v", err)
		}

		if len(resp.Properties) != 1 {
			t.Errorf("Expected 1 property, got %d", len(resp.Properties))
		}

		if resp.Properties[0].ID != propertyID {
			t.Errorf("Expected property ID %d, got %d", propertyID, resp.Properties[0].ID)
		}
	})
}

func TestIntegration_ListRoomTypes(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list all room types", func(t *testing.T) {
		resp, err := client.ListRoomTypes(ctx, nil)
		if err != nil {
			t.Fatalf("ListRoomTypes failed: %v", err)
		}

		t.Logf("Found %d room types (total: %d)", len(resp.RoomTypes), resp.Total)

		for _, rt := range resp.RoomTypes {
			t.Logf("  - %s (ID: %d)", rt.Title, rt.ID)
		}
	})

	t.Run("list room types with limit", func(t *testing.T) {
		resp, err := client.ListRoomTypes(ctx, &hostex.ListRoomTypesParams{
			Limit: 5,
		})
		if err != nil {
			t.Fatalf("ListRoomTypes with limit failed: %v", err)
		}

		if len(resp.RoomTypes) > 5 {
			t.Errorf("Expected max 5 room types, got %d", len(resp.RoomTypes))
		}
	})
}

// ============================================================
// RESERVATIONS TESTS
// ============================================================

func TestIntegration_ListReservations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list all reservations", func(t *testing.T) {
		resp, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("ListReservations failed: %v", err)
		}

		t.Logf("Found %d reservations (total: %d)", len(resp.Reservations), resp.Total)

		for _, res := range resp.Reservations {
			t.Logf("  - %s: %s (%s)", res.ReservationCode, res.GuestName, res.Status)
			t.Logf("    Check-in: %s, Check-out: %s", res.CheckInDate, res.CheckOutDate)
			t.Logf("    Channel: %s", res.ChannelType)
		}
	})

	t.Run("filter by status", func(t *testing.T) {
		resp, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
			Status: "accepted",
			Limit:  5,
		})
		if err != nil {
			t.Fatalf("ListReservations by status failed: %v", err)
		}

		for _, res := range resp.Reservations {
			if res.Status != "accepted" {
				t.Errorf("Expected status 'accepted', got '%s'", res.Status)
			}
		}
	})

	t.Run("filter by date range", func(t *testing.T) {
		resp, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
			StartCheckInDate: "2025-09-01",
			EndCheckInDate:   "2025-12-31",
			Limit:            5,
		})
		if err != nil {
			t.Fatalf("ListReservations by date range failed: %v", err)
		}

		t.Logf("Found %d reservations in date range", len(resp.Reservations))
	})

	t.Run("filter by property ID", func(t *testing.T) {
		// First get a property ID
		props, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{Limit: 1})
		if err != nil || len(props.Properties) == 0 {
			t.Skip("No properties available")
		}

		propertyID := props.Properties[0].ID

		resp, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
			PropertyID: propertyID,
			Limit:      5,
		})
		if err != nil {
			t.Fatalf("ListReservations by property failed: %v", err)
		}

		for _, res := range resp.Reservations {
			if res.PropertyID != propertyID {
				t.Errorf("Expected property ID %d, got %d", propertyID, res.PropertyID)
			}
		}
	})

	t.Run("filter by reservation code", func(t *testing.T) {
		// First get a reservation code
		allRes, err := client.ListReservations(ctx, &hostex.ListReservationsParams{Limit: 1})
		if err != nil || len(allRes.Reservations) == 0 {
			t.Skip("No reservations available")
		}

		code := allRes.Reservations[0].ReservationCode

		resp, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
			ReservationCode: code,
		})
		if err != nil {
			t.Fatalf("ListReservations by code failed: %v", err)
		}

		if len(resp.Reservations) != 1 {
			t.Errorf("Expected 1 reservation, got %d", len(resp.Reservations))
		}

		if resp.Reservations[0].ReservationCode != code {
			t.Errorf("Expected code %s, got %s", code, resp.Reservations[0].ReservationCode)
		}
	})
}

func TestIntegration_ReservationCustomFields(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Get an accepted reservation to test with
	reservations, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
		Status: "accepted",
		Limit:  1,
	})
	if err != nil || len(reservations.Reservations) == 0 {
		t.Skip("No accepted reservations available for testing")
	}

	stayCode := reservations.Reservations[0].StayCode

	t.Run("get custom fields", func(t *testing.T) {
		resp, err := client.GetCustomFields(ctx, stayCode)
		if err != nil {
			t.Fatalf("GetCustomFields failed: %v", err)
		}

		t.Logf("Custom fields for %s: %v", stayCode, resp.CustomFields)
	})

	t.Run("update custom fields", func(t *testing.T) {
		testFields := map[string]interface{}{
			"test_field":       "test_value",
			"test_timestamp":   time.Now().Format(time.RFC3339),
			"integration_test": "true",
		}

		err := client.UpdateCustomFields(ctx, stayCode, testFields)
		if err != nil {
			t.Fatalf("UpdateCustomFields failed: %v", err)
		}

		// Verify the update
		resp, err := client.GetCustomFields(ctx, stayCode)
		if err != nil {
			t.Fatalf("GetCustomFields after update failed: %v", err)
		}

		if resp.CustomFields["test_field"] != "test_value" {
			t.Errorf("Expected test_field to be 'test_value', got %v", resp.CustomFields["test_field"])
		}

		t.Logf("Updated custom fields: %v", resp.CustomFields)
	})
}

// ============================================================
// CONVERSATIONS TESTS
// ============================================================

func TestIntegration_Conversations(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	var testConversationID string

	t.Run("list conversations", func(t *testing.T) {
		resp, err := client.ListConversations(ctx, &hostex.ListConversationsParams{
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("ListConversations failed: %v", err)
		}

		t.Logf("Found %d conversations (total: %d)", len(resp.Conversations), resp.Total)

		for _, conv := range resp.Conversations {
			t.Logf("  - ID: %s, Guest: %s, Property: %s", conv.ID, conv.Guest.Name, conv.PropertyTitle)

			if testConversationID == "" && conv.ID != "" {
				testConversationID = conv.ID
			}
		}
	})

	t.Run("list conversations with offset", func(t *testing.T) {
		resp, err := client.ListConversations(ctx, &hostex.ListConversationsParams{
			Offset: 2,
			Limit:  3,
		})
		if err != nil {
			t.Fatalf("ListConversations with offset failed: %v", err)
		}

		if len(resp.Conversations) > 3 {
			t.Errorf("Expected max 3 conversations, got %d", len(resp.Conversations))
		}
	})

	if testConversationID != "" {
		t.Run("get conversation details", func(t *testing.T) {
			details, err := client.GetConversation(ctx, testConversationID)
			if err != nil {
				t.Fatalf("GetConversation failed: %v", err)
			}

			t.Logf("Conversation %s:", testConversationID)
			t.Logf("  Guest: %s", details.Guest.Name)
			t.Logf("  Channel: %s", details.ChannelType)
			t.Logf("  Messages: %d", len(details.Messages))

			for i, msg := range details.Messages {
				if i >= 3 {
					break // Only show first 3
				}
				t.Logf("    [%s] %s: %s", msg.CreatedAt.Format("2006-01-02 15:04"), msg.SenderRole, msg.Content[:min(50, len(msg.Content))])
			}
		})

		// Note: We're not testing SendMessage in automated tests to avoid sending actual messages to guests
		t.Run("send message - skipped", func(t *testing.T) {
			t.Skip("Skipping SendMessage to avoid sending real messages to guests")
			// To manually test:
			// err := client.SendMessage(ctx, testConversationID, hostex.SendMessageData{
			// 	Message: "Test message from integration test",
			// })
		})
	}
}

// ============================================================
// REVIEWS TESTS
// ============================================================

func TestIntegration_Reviews(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list all reviews", func(t *testing.T) {
		resp, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("ListReviews failed: %v", err)
		}

		t.Logf("Found %d reviews (total: %d)", len(resp.Reviews), resp.Total)

		for _, review := range resp.Reviews {
			t.Logf("  - Reservation: %s, Status: %s", review.ReservationCode, review.ReviewStatus)
			if review.GuestReview != nil {
				t.Logf("    Guest: %d/5", review.GuestReview.Score)
			}
			if review.HostReview != nil {
				t.Logf("    Host: %d/5", review.HostReview.Score)
			}
		}
	})

	t.Run("filter by review status", func(t *testing.T) {
		resp, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
			ReviewStatus: "reviewed",
			Limit:        5,
		})
		if err != nil {
			t.Fatalf("ListReviews by status failed: %v", err)
		}

		t.Logf("Found %d reviewed reviews", len(resp.Reviews))
	})

	t.Run("filter pending host reviews", func(t *testing.T) {
		resp, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
			ReviewStatus: "pending_host_review",
			Limit:        5,
		})
		if err != nil {
			t.Fatalf("ListReviews pending failed: %v", err)
		}

		t.Logf("Found %d pending host reviews", len(resp.Reviews))
	})

	t.Run("filter by property", func(t *testing.T) {
		props, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{Limit: 1})
		if err != nil || len(props.Properties) == 0 {
			t.Skip("No properties available")
		}

		propertyID := props.Properties[0].ID

		resp, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
			PropertyID: propertyID,
			Limit:      5,
		})
		if err != nil {
			t.Fatalf("ListReviews by property failed: %v", err)
		}

		t.Logf("Found %d reviews for property %d", len(resp.Reviews), propertyID)
	})

	t.Run("filter by reservation code", func(t *testing.T) {
		allReviews, err := client.ListReviews(ctx, &hostex.ListReviewsParams{Limit: 1})
		if err != nil || len(allReviews.Reviews) == 0 {
			t.Skip("No reviews available")
		}

		code := allReviews.Reviews[0].ReservationCode

		resp, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
			ReservationCode: code,
		})
		if err != nil {
			t.Fatalf("ListReviews by code failed: %v", err)
		}

		if len(resp.Reviews) == 0 {
			t.Error("Expected at least 1 review")
		}
	})

	// Note: Not testing CreateReview in automated tests to avoid creating real reviews
	t.Run("create review - skipped", func(t *testing.T) {
		t.Skip("Skipping CreateReview to avoid creating real reviews")
		// To manually test with a pending review:
		// err := client.CreateReview(ctx, "reservation_code", hostex.CreateReviewData{
		// 	HostReviewScore:   5,
		// 	HostReviewContent: "Test review",
		// })
	})
}

// ============================================================
// AVAILABILITIES TESTS
// ============================================================

func TestIntegration_Availabilities(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	props, err := client.ListProperties(ctx, &hostex.ListPropertiesParams{Limit: 1})
	if err != nil || len(props.Properties) == 0 {
		t.Skip("No properties available for availability testing")
	}

	propertyID := props.Properties[0].ID

	t.Run("list availabilities", func(t *testing.T) {
		resp, err := client.ListAvailabilities(ctx, hostex.ListAvailabilitiesParams{
			PropertyIDs: fmt.Sprintf("%d", propertyID),
			StartDate:   "2025-11-01",
			EndDate:     "2025-11-30",
		})
		if err != nil {
			t.Fatalf("ListAvailabilities failed: %v", err)
		}

		t.Logf("Found %d listings with availability data", len(resp.Listings))

		for _, listing := range resp.Listings {
			t.Logf("  - Channel: %s, Listing: %s", listing.ChannelType, listing.ListingID)
			t.Logf("    Availabilities: %d days", len(listing.Availabilities))
		}
	})

	// Note: Not testing UpdateAvailabilities to avoid modifying real calendar
	t.Run("update availabilities - skipped", func(t *testing.T) {
		t.Skip("Skipping UpdateAvailabilities to avoid modifying real calendar")
		// To manually test:
		// err := client.UpdateAvailabilities(ctx, hostex.UpdateAvailabilitiesData{
		// 	PropertyIDs: []int{propertyID},
		// 	Dates:       []string{"2025-12-25"},
		// 	Available:   false,
		// })
	})
}

// ============================================================
// LISTINGS TESTS
// ============================================================

func TestIntegration_Listings(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	// Get a property with channels
	props, err := client.ListProperties(ctx, nil)
	if err != nil || len(props.Properties) == 0 {
		t.Skip("No properties available")
	}

	var testListings []hostex.Listing
	for _, prop := range props.Properties {
		for _, ch := range prop.Channels {
			testListings = append(testListings, hostex.Listing{
				ChannelType: ch.ChannelType,
				ListingID:   ch.ListingID,
			})
			if len(testListings) >= 2 {
				break
			}
		}
		if len(testListings) >= 2 {
			break
		}
	}

	if len(testListings) == 0 {
		t.Skip("No channel listings available")
	}

	t.Run("get listing calendar", func(t *testing.T) {
		resp, err := client.GetListingCalendar(ctx, hostex.GetListingCalendarData{
			StartDate: "2025-11-01",
			EndDate:   "2025-11-30",
			Listings:  testListings,
		})
		if err != nil {
			t.Fatalf("GetListingCalendar failed: %v", err)
		}

		t.Logf("Got calendar for %d listings", len(resp.Listings))

		for _, listing := range resp.Listings {
			t.Logf("  - %s/%s: %d days", listing.ChannelType, listing.ListingID, len(listing.Calendar))
		}
	})

	// Note: Not testing update operations to avoid modifying real listings
	t.Run("update listing prices - skipped", func(t *testing.T) {
		t.Skip("Skipping UpdateListingPrices to avoid modifying real prices")
	})

	t.Run("update listing inventories - skipped", func(t *testing.T) {
		t.Skip("Skipping UpdateListingInventories to avoid modifying real inventory")
	})

	t.Run("update listing restrictions - skipped", func(t *testing.T) {
		t.Skip("Skipping UpdateListingRestrictions to avoid modifying real restrictions")
	})
}

// ============================================================
// WEBHOOKS TESTS
// ============================================================

func TestIntegration_Webhooks(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list webhooks", func(t *testing.T) {
		resp, err := client.ListWebhooks(ctx)
		if err != nil {
			t.Fatalf("ListWebhooks failed: %v", err)
		}

		t.Logf("Found %d webhooks", len(resp.Webhooks))

		for _, wh := range resp.Webhooks {
			t.Logf("  - ID: %d, URL: %s, Manageable: %v", wh.ID, wh.URL, wh.Manageable)
		}
	})

	// Note: Not testing webhook create/delete to avoid modifying real webhooks
	t.Run("create webhook - skipped", func(t *testing.T) {
		t.Skip("Skipping CreateWebhook to avoid creating real webhooks")
		// To manually test:
		// resp, err := client.CreateWebhook(ctx, "https://example.com/webhook")
	})

	t.Run("delete webhook - skipped", func(t *testing.T) {
		t.Skip("Skipping DeleteWebhook to avoid deleting real webhooks")
		// To manually test:
		// err := client.DeleteWebhook(ctx, webhookID)
	})
}

// ============================================================
// UTILITIES TESTS
// ============================================================

func TestIntegration_CustomChannels(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list custom channels", func(t *testing.T) {
		resp, err := client.ListCustomChannels(ctx)
		if err != nil {
			t.Fatalf("ListCustomChannels failed: %v", err)
		}

		if len(resp.CustomChannels) == 0 {
			t.Error("Expected at least 1 custom channel")
		}

		t.Logf("Found %d custom channels", len(resp.CustomChannels))

		for _, ch := range resp.CustomChannels {
			t.Logf("  - %s (ID: %d)", ch.Name, ch.ID)
		}
	})
}

func TestIntegration_IncomeMethods(t *testing.T) {
	client := getTestClient(t)
	ctx := context.Background()

	t.Run("list income methods", func(t *testing.T) {
		resp, err := client.ListIncomeMethods(ctx)
		if err != nil {
			t.Fatalf("ListIncomeMethods failed: %v", err)
		}

		if len(resp.IncomeMethods) == 0 {
			t.Error("Expected at least 1 income method")
		}

		t.Logf("Found %d income methods", len(resp.IncomeMethods))

		for _, im := range resp.IncomeMethods {
			t.Logf("  - %s (ID: %d)", im.Name, im.ID)
		}
	})
}

// Helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
