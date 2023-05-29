package usecase

import (
	"sync"
)

type RelayManager struct {
	relayMap map[string]*RelayConnector //KEY: public key
	mux      sync.Mutex
}

var relayManager *RelayManager

func newRelayManager() *RelayManager {
	s := &RelayManager{}
	s.relayMap = make(map[string]*RelayConnector)

	return s
}

func GetRelayManager() *RelayManager {
	if relayManager == nil {
		relayManager = newRelayManager()

	}
	return relayManager
}

func (t *RelayManager) AddRelay(url, pubKey, privateKey string) (
	*RelayConnector, error) {

	// rl := t.GetRelay(pubKey)
	// if rl != nil {

	// 	return rl, nil
	// }

	u, err := NewRelayConnector(url, pubKey, privateKey)
	if err != nil {
		return nil, err
	}

	t.mux.Lock()
	t.relayMap[url] = u
	t.mux.Unlock()
	return u, nil
}

func (t *RelayManager) GetRelay(url string) *RelayConnector {

	t.mux.Lock()
	defer t.mux.Unlock()
	ret, ok := t.relayMap[url]
	if ok {
		return ret
	}
	return nil
}
