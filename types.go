package hostex

import "time"

// Property represents a Hostex property
type Property struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Address   string    `json:"address,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	Latitude  float64   `json:"latitude,omitempty"`
	Channels  []Channel `json:"channels,omitempty"`
}

// Channel represents a booking channel for a property
type Channel struct {
	ChannelType string `json:"channel_type"`
	ListingID   string `json:"listing_id"`
}

// RoomType represents a room type in Hostex
type RoomType struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Properties []Property `json:"properties,omitempty"`
	Channels   []Channel  `json:"channels,omitempty"`
}

// Reservation represents a booking reservation
type Reservation struct {
	ReservationCode  string     `json:"reservation_code"`
	StayCode         string     `json:"stay_code"`
	ChannelID        string     `json:"channel_id,omitempty"`
	PropertyID       int        `json:"property_id"`
	ChannelType      string     `json:"channel_type"`
	ListingID        string     `json:"listing_id,omitempty"`
	CheckInDate      string     `json:"check_in_date"`
	CheckOutDate     string     `json:"check_out_date"`
	NumberOfGuests   int        `json:"number_of_guests,omitempty"`
	NumberOfAdults   int        `json:"number_of_adults,omitempty"`
	NumberOfChildren int        `json:"number_of_children,omitempty"`
	NumberOfInfants  int        `json:"number_of_infants,omitempty"`
	NumberOfPets     int        `json:"number_of_pets,omitempty"`
	Status           string     `json:"status"`
	GuestName        string     `json:"guest_name,omitempty"`
	GuestPhone       string     `json:"guest_phone,omitempty"`
	GuestEmail       string     `json:"guest_email,omitempty"`
	CancelledAt      *time.Time `json:"cancelled_at,omitempty"`
	BookedAt         *time.Time `json:"booked_at,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	Creator          string     `json:"creator,omitempty"`
	ConversationID   string     `json:"conversation_id,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	CustomFields     any        `json:"custom_fields,omitempty"`
	InReservationBox bool       `json:"in_reservation_box,omitempty"`
}

// CreateReservationData contains data for creating a new reservation
type CreateReservationData struct {
	PropertyID       string `json:"property_id"`
	CustomChannelID  int    `json:"custom_channel_id"`
	CheckInDate      string `json:"check_in_date"`
	CheckOutDate     string `json:"check_out_date"`
	GuestName        string `json:"guest_name"`
	Currency         string `json:"currency"`
	RateAmount       int    `json:"rate_amount"`
	CommissionAmount int    `json:"commission_amount"`
	ReceivedAmount   int    `json:"received_amount"`
	IncomeMethodID   int    `json:"income_method_id"`
	NumberOfGuests   int    `json:"number_of_guests,omitempty"`
	Email            string `json:"email,omitempty"`
	Mobile           string `json:"mobile,omitempty"`
	Remarks          string `json:"remarks,omitempty"`
}

// Conversation represents a guest conversation
type Conversation struct {
	ID            string    `json:"id"`
	ChannelType   string    `json:"channel_type"`
	Guest         Guest     `json:"guest"`
	PropertyID    int       `json:"property_id,omitempty"`
	PropertyTitle string    `json:"property_title,omitempty"`
	CheckInDate   string    `json:"check_in_date,omitempty"`
	CheckOutDate  string    `json:"check_out_date,omitempty"`
	LastMessageAt time.Time `json:"last_message_at,omitempty"`
	UnreadCount   int       `json:"unread_count,omitempty"`
}

// Guest represents a guest
type Guest struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Mobile string `json:"mobile,omitempty"`
}

// Message represents a conversation message
type Message struct {
	ID         string    `json:"id"`
	SenderRole string    `json:"sender_role"` // "host" or "guest"
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	ImageURL   string    `json:"image_url,omitempty"`
}

// SendMessageData contains data for sending a message
type SendMessageData struct {
	Message    string `json:"message,omitempty"`
	JpegBase64 string `json:"jpeg_base64,omitempty"`
}

// Review represents a guest or host review
type Review struct {
	ReservationCode string      `json:"reservation_code"`
	PropertyID      int         `json:"property_id"`
	ChannelType     string      `json:"channel_type"`
	CheckOutDate    string      `json:"check_out_date"`
	ReviewStatus    string      `json:"review_status"`
	HostReview      *ReviewData `json:"host_review,omitempty"`
	GuestReview     *ReviewData `json:"guest_review,omitempty"`
	HostReply       *ReplyData  `json:"host_reply,omitempty"`
}

// ReviewData contains review details
type ReviewData struct {
	Score     int       `json:"score,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// ReplyData contains reply details
type ReplyData struct {
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// CreateReviewData contains data for creating a review
type CreateReviewData struct {
	HostReviewScore   int    `json:"host_review_score,omitempty"`
	HostReviewContent string `json:"host_review_content,omitempty"`
	HostReplyContent  string `json:"host_reply_content,omitempty"`
}

// Webhook represents a webhook configuration
type Webhook struct {
	ID         int       `json:"id"`
	URL        string    `json:"url"`
	Manageable bool      `json:"manageable"`
	CreatedAt  time.Time `json:"created_at"`
}

// CustomChannel represents a custom booking channel
type CustomChannel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IncomeMethod represents a payment method
type IncomeMethod struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Availability represents property availability
type Availability struct {
	Date      string `json:"date"`
	Available bool   `json:"available"`
}

// UpdateAvailabilitiesData contains data for updating availability
type UpdateAvailabilitiesData struct {
	PropertyIDs []int    `json:"property_ids"`
	StartDate   string   `json:"start_date,omitempty"`
	EndDate     string   `json:"end_date,omitempty"`
	Dates       []string `json:"dates,omitempty"`
	Available   bool     `json:"available"`
}

// ListingCalendarData contains data for getting listing calendar
type GetListingCalendarData struct {
	StartDate string    `json:"start_date"`
	EndDate   string    `json:"end_date"`
	Listings  []Listing `json:"listings"`
}

// Listing represents a listing on a channel
type Listing struct {
	ChannelType string `json:"channel_type"`
	ListingID   string `json:"listing_id"`
}

// UpdateListingPricesData contains data for updating listing prices
type UpdateListingPricesData struct {
	ChannelType string  `json:"channel_type"`
	ListingID   string  `json:"listing_id"`
	Prices      []Price `json:"prices"`
}

// Price represents a price for a specific date
type Price struct {
	Date  string `json:"date"`
	Price int    `json:"price"`
}

// UpdateListingInventoriesData contains data for updating listing inventories
type UpdateListingInventoriesData struct {
	ChannelType string      `json:"channel_type"`
	ListingID   string      `json:"listing_id"`
	Inventories []Inventory `json:"inventories"`
}

// Inventory represents inventory for a specific date
type Inventory struct {
	Date      string `json:"date"`
	Inventory int    `json:"inventory"`
}

// UpdateListingRestrictionsData contains data for updating listing restrictions
type UpdateListingRestrictionsData struct {
	ChannelType  string        `json:"channel_type"`
	ListingID    string        `json:"listing_id"`
	Restrictions []Restriction `json:"restrictions"`
}

// Restriction represents restrictions for a specific date
type Restriction struct {
	Date              string `json:"date"`
	MinStay           int    `json:"min_stay,omitempty"`
	MaxStay           int    `json:"max_stay,omitempty"`
	ClosedToArrival   bool   `json:"closed_to_arrival,omitempty"`
	ClosedToDeparture bool   `json:"closed_to_departure,omitempty"`
}
