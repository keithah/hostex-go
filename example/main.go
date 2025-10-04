package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/keithah/hostex-go"
)

func main() {
	// Get API token from environment
	token := os.Getenv("HOSTEX_API_KEY")
	if token == "" {
		log.Fatal("HOSTEX_API_KEY environment variable is required")
	}

	// Create a new Hostex client
	client, err := hostex.NewClient(hostex.Config{
		AccessToken: token,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: List Properties
	fmt.Println("=== Properties ===")
	properties, err := client.ListProperties(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to list properties: %v", err)
	}

	for _, property := range properties.Properties {
		fmt.Printf("- %s (ID: %d)\n", property.Title, property.ID)
		if property.Address != "" {
			fmt.Printf("  Address: %s\n", property.Address)
		}
		if len(property.Channels) > 0 {
			fmt.Printf("  Channels: ")
			for i, ch := range property.Channels {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", ch.ChannelType)
			}
			fmt.Println()
		}
	}

	// Example 2: List Recent Reservations
	fmt.Println("\n=== Recent Reservations ===")
	reservations, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
		Status: "accepted",
		Limit:  5,
	})
	if err != nil {
		log.Fatalf("Failed to list reservations: %v", err)
	}

	for _, res := range reservations.Reservations {
		fmt.Printf("- %s\n", res.ReservationCode)
		fmt.Printf("  Guest: %s\n", res.GuestName)
		fmt.Printf("  Check-in: %s -> Check-out: %s\n", res.CheckInDate, res.CheckOutDate)
		fmt.Printf("  Channel: %s\n", res.ChannelType)
	}

	// Example 3: List Conversations
	fmt.Println("\n=== Conversations ===")
	conversations, err := client.ListConversations(ctx, &hostex.ListConversationsParams{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Failed to list conversations: %v", err)
	}

	for _, conv := range conversations.Conversations {
		fmt.Printf("- %s (ID: %s)\n", conv.Guest.Name, conv.ID)
		fmt.Printf("  Property: %s\n", conv.PropertyTitle)
		fmt.Printf("  Channel: %s\n", conv.ChannelType)
		if !conv.LastMessageAt.IsZero() {
			fmt.Printf("  Last message: %s\n", conv.LastMessageAt.Format("2006-01-02 15:04"))
		}
	}

	// Example 4: List Reviews
	fmt.Println("\n=== Reviews ===")
	reviews, err := client.ListReviews(ctx, &hostex.ListReviewsParams{
		ReviewStatus: "reviewed",
		Limit:        5,
	})
	if err != nil {
		log.Fatalf("Failed to list reviews: %v", err)
	}

	for _, review := range reviews.Reviews {
		fmt.Printf("- Reservation: %s\n", review.ReservationCode)
		if review.GuestReview != nil {
			fmt.Printf("  Guest rating: %d/5\n", review.GuestReview.Score)
			if review.GuestReview.Content != "" {
				fmt.Printf("  Guest comment: %s\n", review.GuestReview.Content)
			}
		}
		if review.HostReview != nil {
			fmt.Printf("  Host rating: %d/5\n", review.HostReview.Score)
			if review.HostReview.Content != "" {
				fmt.Printf("  Host comment: %s\n", review.HostReview.Content)
			}
		}
	}

	// Example 5: List Custom Channels
	fmt.Println("\n=== Custom Channels ===")
	channels, err := client.ListCustomChannels(ctx)
	if err != nil {
		log.Fatalf("Failed to list custom channels: %v", err)
	}

	for _, channel := range channels.CustomChannels {
		fmt.Printf("- %s (ID: %d)\n", channel.Name, channel.ID)
	}

	// Example 6: List Income Methods
	fmt.Println("\n=== Income Methods ===")
	methods, err := client.ListIncomeMethods(ctx)
	if err != nil {
		log.Fatalf("Failed to list income methods: %v", err)
	}

	for _, method := range methods.IncomeMethods {
		fmt.Printf("- %s (ID: %d)\n", method.Name, method.ID)
	}

	fmt.Println("\n=== Done ===")
}
