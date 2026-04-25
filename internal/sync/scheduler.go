package sync

import (
	"context"
	"log"
	"sync/atomic"
	"time"
)

type Scheduler struct {
	syncer   *Syncer
	interval time.Duration
	running  atomic.Bool
	stopCh   chan struct{}
}

func NewScheduler(syncer *Syncer, interval time.Duration) *Scheduler {
	return &Scheduler{
		syncer:   syncer,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	log.Printf("[scheduler] started, interval=%v", s.interval)
	s.runSync(ctx)

	for {
		select {
		case <-ticker.C:
			s.runSync(ctx)
		case <-s.stopCh:
			log.Printf("[scheduler] stopped")
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) runSync(ctx context.Context) {
	if s.running.Swap(true) {
		log.Printf("[scheduler] sync already running, skipping")
		return
	}
	defer s.running.Store(false)

	if err := s.syncer.Sync(ctx); err != nil {
		log.Printf("[scheduler] sync error: %v", err)
	}
}
