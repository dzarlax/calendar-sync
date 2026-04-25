package sync

import "time"

type Event struct {
	ID           string    `json:"id"`
	CalendarID   string    `json:"calendar_id"`
	Provider     string    `json:"provider"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	Location     string    `json:"location,omitempty"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	AllDay       bool      `json:"all_day,omitempty"`
	Status       string    `json:"status,omitempty"`
	Attendees    []Attendee `json:"attendees,omitempty"`
	OnlineMeeting string   `json:"online_meeting,omitempty"`
}

type Attendee struct {
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

type SyncState struct {
	LastSync time.Time      `json:"last_sync"`
	Mappings map[string]MappingEntry `json:"mappings"`
}

type MappingEntry struct {
	GoogleID string `json:"google_id"`
	Hash     string `json:"hash"`
}

type EventCreate struct {
	Title       string     `json:"title"`
	Start       time.Time  `json:"start"`
	End         time.Time  `json:"end"`
	Description string     `json:"description,omitempty"`
	Location    string     `json:"location,omitempty"`
	Attendees   []Attendee `json:"attendees,omitempty"`
}

type EventUpdate struct {
	Title       *string    `json:"title,omitempty"`
	Start       *time.Time `json:"start,omitempty"`
	End         *time.Time `json:"end,omitempty"`
	Description *string    `json:"description,omitempty"`
	Location    *string    `json:"location,omitempty"`
	Attendees   []Attendee `json:"attendees,omitempty"`
}

type Config struct {
	APIBaseURL   string
	APIKey       string
	SyncSource   string
	SyncTarget   string
	SyncInterval time.Duration
	StateFile    string
}
