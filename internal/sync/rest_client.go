package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type RestClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewRestClient(baseURL, apiKey string) *RestClient {
	return &RestClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RestClient) do(req *http.Request, out interface{}) error {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		var e struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&e)
		if e.Error != "" {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, e.Error)
		}
		return fmt.Errorf("API error %d", resp.StatusCode)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func (c *RestClient) GetEvents(ctx context.Context, calendarID string, start, end time.Time) ([]Event, error) {
	q := url.Values{"calendar_id": {calendarID}}
	if !start.IsZero() {
		q.Set("start", start.Format(time.RFC3339))
	}
	if !end.IsZero() {
		q.Set("end", end.Format(time.RFC3339))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/events?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var events []Event
	return events, c.do(req, &events)
}

func (c *RestClient) CreateEvent(ctx context.Context, calendarID string, event EventCreate) (*Event, error) {
	body := struct {
		CalendarID  string     `json:"calendar_id"`
		Title       string     `json:"title"`
		Start       string     `json:"start"`
		End         string     `json:"end"`
		Description string     `json:"description,omitempty"`
		Location    string     `json:"location,omitempty"`
		Attendees   []Attendee `json:"attendees,omitempty"`
	}{
		CalendarID:  calendarID,
		Title:       event.Title,
		Start:       event.Start.Format(time.RFC3339),
		End:         event.End.Format(time.RFC3339),
		Description: event.Description,
		Location:    event.Location,
		Attendees:   event.Attendees,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/events", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var ev Event
	return &ev, c.do(req, &ev)
}

func (c *RestClient) UpdateEvent(ctx context.Context, calendarID, eventID string, event EventUpdate) (*Event, error) {
	body := struct {
		Title       *string    `json:"title,omitempty"`
		Start       *string    `json:"start,omitempty"`
		End         *string    `json:"end,omitempty"`
		Description *string    `json:"description,omitempty"`
		Location    *string    `json:"location,omitempty"`
		Attendees   []Attendee `json:"attendees,omitempty"`
	}{
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		Attendees:   event.Attendees,
	}
	if event.Start != nil {
		s := event.Start.Format(time.RFC3339)
		body.Start = &s
	}
	if event.End != nil {
		e := event.End.Format(time.RFC3339)
		body.End = &e
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	q := url.Values{"calendar_id": {calendarID}, "event_id": {eventID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.baseURL+"/api/events?"+q.Encode(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var ev Event
	return &ev, c.do(req, &ev)
}

func (c *RestClient) DeleteEvent(ctx context.Context, calendarID, eventID string) error {
	q := url.Values{"calendar_id": {calendarID}, "event_id": {eventID}}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/api/events?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
