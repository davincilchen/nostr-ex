package usecase

import (
	"sync"
)

type RelayManager struct {
	id       int
	relayMap map[int]*RelayConnector //KEY: public key
	mux      sync.Mutex
}

var relayManager *RelayManager

func newRelayManager() *RelayManager {
	s := &RelayManager{}
	s.relayMap = make(map[int]*RelayConnector)
	s.id = -1
	return s
}

func GetRelayManager() *RelayManager {
	if relayManager == nil {
		relayManager = newRelayManager()

	}
	return relayManager
}

func (t *RelayManager) NewID() int {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.id++
	return t.id
}

func (t *RelayManager) AddRelay(url string) (
	*RelayConnector, error) {

	id := t.NewID()
	u, err := NewRelayConnector(id, url)
	if err != nil {
		return nil, err
	}

	t.mux.Lock()
	t.relayMap[id] = u
	t.mux.Unlock()
	return u, nil
}

func (t *RelayManager) GetRelay(id int) *RelayConnector {

	t.mux.Lock()
	defer t.mux.Unlock()
	ret, ok := t.relayMap[id]
	if ok {
		return ret
	}
	return nil
}
