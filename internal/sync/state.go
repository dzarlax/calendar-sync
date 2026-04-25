package sync

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type StateManager struct {
	path   string
	mu     sync.RWMutex
	state  *SyncState
}

func NewStateManager(path string) (*StateManager, error) {
	sm := &StateManager{path: path, state: &SyncState{Mappings: make(map[string]MappingEntry)}}
	if err := sm.load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	return sm, nil
}

func (sm *StateManager) load() error {
	data, err := os.ReadFile(sm.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, sm.state)
}

func (sm *StateManager) save() error {
	data, err := json.MarshalIndent(sm.state, "", "  ")
	if err != nil {
		return err
	}
	tmp := sm.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, sm.path)
}

func (sm *StateManager) GetState() SyncState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	// return a shallow copy so callers can't mutate state through the pointer
	s := *sm.state
	mappings := make(map[string]MappingEntry, len(sm.state.Mappings))
	for k, v := range sm.state.Mappings {
		mappings[k] = v
	}
	s.Mappings = mappings
	return s
}

func (sm *StateManager) SetLastSync(t time.Time) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.LastSync = t
	return sm.save()
}

func (sm *StateManager) SetMapping(m365ID, googleID, hash string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.Mappings[m365ID] = MappingEntry{GoogleID: googleID, Hash: hash}
	return sm.save()
}

func (sm *StateManager) DeleteMapping(m365ID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.state.Mappings, m365ID)
	return sm.save()
}

func (sm *StateManager) GetMapping(m365ID string) (MappingEntry, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	entry, ok := sm.state.Mappings[m365ID]
	return entry, ok
}
