# hostex-go

[![Go Reference](https://pkg.go.dev/badge/github.com/keithah/hostex-go.svg)](https://pkg.go.dev/github.com/keithah/hostex-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/keithah/hostex-go)](https://goreportcard.com/report/github.com/keithah/hostex-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the Hostex API v3.0.0 (Beta). Hostex is a property management platform that provides APIs for managing reservations, properties, availabilities, messaging, reviews, and more.

## Features

- **Complete API Coverage**: All Hostex API v3.0.0 endpoints
- **Type-Safe**: Comprehensive type definitions for all API entities
- **Context Support**: All methods accept context for cancellation and timeouts
- **Error Handling**: Detailed error responses with API error codes
- **Clean API**: Idiomatic Go interfaces and patterns
- **Well Documented**: Extensive documentation and examples
- **Production Ready**: Built for reliability and performance

## Installation

```bash
go get github.com/keithah/hostex-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/keithah/hostex-go"
)

func main() {
	// Create a new client
	client, err := hostex.NewClient(hostex.Config{
		AccessToken: "your_hostex_api_token",
	})
	if err != nil {
		log.Fatal(err)
	}

	// List properties
	properties, err := client.ListProperties(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, property := range properties.Properties {
		fmt.Printf("%s (ID: %d)\n", property.Title, property.ID)
	}
}
```

## API Coverage

### Properties
- `ListProperties` - List all properties
- `ListRoomTypes` - List room types

### Reservations
- `ListReservations` - Search and filter reservations
- `CreateReservation` - Create direct bookings
- `CancelReservation` - Cancel reservations
- `UpdateLockCode` - Update stay lock codes
- `GetCustomFields` - Get custom field values
- `UpdateCustomFields` - Update custom fields

### Availabilities
- `ListAvailabilities` - Check property availability
- `UpdateAvailabilities` - Block/open dates

### Conversations (Messaging)
- `ListConversations` - List guest conversations
- `GetConversation` - Get conversation details and messages
- `SendMessage` - Send messages to guests

### Reviews
- `ListReviews` - Query reviews
- `CreateReview` - Leave reviews or replies

### Webhooks
- `ListWebhooks` - List configured webhooks
- `CreateWebhook` - Register new webhooks
- `DeleteWebhook` - Remove webhooks

### Listings
- `GetListingCalendar` - Get listing calendars
- `UpdateListingPrices` - Update channel prices
- `UpdateListingInventories` - Update inventory levels
- `UpdateListingRestrictions` - Update restrictions

### Utilities
- `ListCustomChannels` - List custom channels
- `ListIncomeMethods` - List income methods

## Usage Examples

### List Recent Reservations

```go
reservations, err := client.ListReservations(ctx, &hostex.ListReservationsParams{
	Status:           "accepted",
	StartCheckInDate: "2024-01-01",
	Limit:            50,
})
if err != nil {
	log.Fatal(err)
}

for _, res := range reservations.Reservations {
	fmt.Printf("%s: %s (%s to %s)\n",
		res.ReservationCode,
		res.GuestName,
		res.CheckInDate,
		res.CheckOutDate,
	)
}
```

### Send a Message

```go
err := client.SendMessage(ctx, "conversation_id", hostex.SendMessageData{
	Message: "Thank you for your booking! Check-in is at 3 PM.",
})
if err != nil {
	log.Fatal(err)
}
```

### Update Availability

```go
err := client.UpdateAvailabilities(ctx, hostex.UpdateAvailabilitiesData{
	PropertyIDs: []int{12345},
	Dates:       []string{"2024-07-15", "2024-07-16"},
	Available:   false, // Block these dates
})
if err != nil {
	log.Fatal(err)
}
```

### Create a Direct Booking

```go
reservation, err := client.CreateReservation(ctx, hostex.CreateReservationData{
	PropertyID:       "12345",
	CustomChannelID:  1,
	CheckInDate:      "2024-07-01",
	CheckOutDate:     "2024-07-07",
	GuestName:        "Jane Smith",
	Currency:         "USD",
	RateAmount:       84000, // $840.00 in cents
	CommissionAmount: 8400,
	ReceivedAmount:   75600,
	IncomeMethodID:   1,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Created reservation: %s\n", reservation.Reservation.ReservationCode)
```

### Get Conversation Messages

```go
conversation, err := client.GetConversation(ctx, "conversation_id")
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Guest: %s\n", conversation.Guest.Name)
for _, msg := range conversation.Messages {
	fmt.Printf("[%s] %s: %s\n",
		msg.CreatedAt.Format("2006-01-02 15:04"),
		msg.SenderRole,
		msg.Content,
	)
}
```

## Configuration

### Custom HTTP Client

```go
import (
	"net/http"
	"time"
)

client, err := hostex.NewClient(hostex.Config{
	AccessToken: "your_token",
	HTTPClient: &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			// Custom transport settings
		},
	},
})
```

### Custom Base URL

```go
client, err := hostex.NewClient(hostex.Config{
	AccessToken: "your_token",
	BaseURL:     "https://api-staging.hostex.io/v3",
})
```

### Custom Timeout

```go
client, err := hostex.NewClient(hostex.Config{
	AccessToken: "your_token",
	Timeout:     60 * time.Second,
})
```

## Error Handling

All API methods return errors that include the Hostex API error code and message:

```go
reservations, err := client.ListReservations(ctx, nil)
if err != nil {
	// Error includes API error code and message
	log.Printf("API error: %v", err)
	return
}
```

## Context Usage

All API methods accept a `context.Context` parameter for cancellation and timeouts:

```go
// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

properties, err := client.ListProperties(ctx, nil)

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Cancel the request if needed
cancel()
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/keithah/hostex-go)
- [Hostex API Documentation](https://hostex-openapi.readme.io/)

## Testing

The library includes comprehensive integration tests covering all API endpoints.

### Running Tests

```bash
# Run unit tests (no API key needed)
make test-unit

# Run integration tests (requires API key)
export HOSTEX_API_KEY=your_key
make test-integration

# Or source .env file
source .env
make test-integration

# Run all tests with coverage
make test-all

# Check code quality
make check
```

### Test Coverage

All endpoints are tested with comprehensive integration tests:

- ✅ **Properties** - List, filter, pagination
- ✅ **Room Types** - List, filter
- ✅ **Reservations** - List, filter by status/date/property/code
- ✅ **Custom Fields** - Get and update
- ✅ **Conversations** - List, get details, pagination
- ✅ **Reviews** - List, filter by status/property/code
- ✅ **Availabilities** - List
- ✅ **Listings** - Get calendar
- ✅ **Webhooks** - List
- ✅ **Custom Channels** - List
- ✅ **Income Methods** - List

Tests that modify data (create/update/delete) are marked as skipped in automated runs to avoid affecting production data.

### Continuous Integration

GitHub Actions runs tests automatically on every push:
- Unit tests run on all PRs and pushes
- Integration tests run on pushes to main (when `HOSTEX_API_KEY` secret is configured)
- Code quality checks (formatting, vet, staticcheck)

## Requirements

- Go 1.18 or higher

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes and add tests
4. Run tests: `go test ./...`
5. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/keithah/hostex-go/issues)
- **Hostex API Support**: contact@hostex.io

---

**Note**: This library is not officially supported by Hostex. It's a community-driven project to provide Go developers with easy access to the Hostex API.
