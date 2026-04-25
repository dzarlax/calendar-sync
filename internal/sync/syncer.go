package sync

import (
	"context"
	"log"
	"time"
)

type Syncer struct {
	client  *RestClient
	state   *StateManager
	source  string
	target  string
}

func NewSyncer(client *RestClient, state *StateManager, source, target string) *Syncer {
	return &Syncer{
		client: client,
		state:  state,
		source: source,
		target: target,
	}
}

func (s *Syncer) Sync(ctx context.Context) error {
	start := time.Now()

	snap := s.state.GetState()
	syncStart := snap.LastSync
	if syncStart.IsZero() {
		syncStart = time.Now().AddDate(0, 0, -30)
	}
	events, err := s.client.GetEvents(ctx, s.source, syncStart, time.Now())
	if err != nil {
		return err
	}

	log.Printf("[syncer] fetched %d events from %s since %v", len(events), s.source, snap.LastSync)

	currentM365IDs := make(map[string]struct{})
	for _, evt := range events {
		currentM365IDs[evt.ID] = struct{}{}
		if err := s.syncEvent(ctx, evt); err != nil {
			log.Printf("[syncer] failed to sync event %s: %v", evt.ID, err)
		}
	}

	if err := s.cleanupDeleted(ctx, currentM365IDs); err != nil {
		log.Printf("[syncer] cleanup error: %v", err)
	}

	if err := s.state.SetLastSync(time.Now()); err != nil {
		return err
	}

	log.Printf("[syncer] sync completed in %v", time.Since(start))
	return nil
}

func (s *Syncer) syncEvent(ctx context.Context, evt Event) error {
	hash := HashEvent(evt)
	entry, exists := s.state.GetMapping(evt.ID)

	if !exists {
		created, err := s.client.CreateEvent(ctx, s.target, toEventCreate(evt))
		if err != nil {
			return err
		}
		log.Printf("[syncer] created event %s → %s", evt.ID, created.ID)
		return s.state.SetMapping(evt.ID, created.ID, hash)
	}

	if entry.Hash != hash {
		_, err := s.client.UpdateEvent(ctx, s.target, entry.GoogleID, toEventUpdate(evt))
		if err != nil {
			return err
		}
		log.Printf("[syncer] updated event %s (%s)", evt.ID, entry.GoogleID)
		return s.state.SetMapping(evt.ID, entry.GoogleID, hash)
	}

	return nil
}

func (s *Syncer) cleanupDeleted(ctx context.Context, currentIDs map[string]struct{}) error {
	snap := s.state.GetState()
	for m365ID, entry := range snap.Mappings {
		if _, present := currentIDs[m365ID]; !present {
			if err := s.client.DeleteEvent(ctx, s.target, entry.GoogleID); err != nil {
				log.Printf("[syncer] failed to delete %s: %v", entry.GoogleID, err)
				continue
			}
			log.Printf("[syncer] deleted stale event %s from Google", entry.GoogleID)
			if err := s.state.DeleteMapping(m365ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func toEventCreate(evt Event) EventCreate {
	return EventCreate{
		Title:       evt.Title,
		Start:       evt.Start,
		End:         evt.End,
		Description: evt.Description,
		Location:    evt.Location,
	}
}

func toEventUpdate(evt Event) EventUpdate {
	return EventUpdate{
		Title:       &evt.Title,
		Start:       &evt.Start,
		End:         &evt.End,
		Description: &evt.Description,
		Location:    &evt.Location,
	}
}
