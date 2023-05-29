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

	rl := t.GetRelay(pubKey)
	if rl != nil {
		if privateKey != "" {
			rl.UpdatePrivateKey(privateKey)
		}
		return rl, nil
	}

	u, err := NewRelayConnector(url, pubKey, privateKey)
	if err != nil {
		return nil, err
	}

	t.mux.Lock()
	t.relayMap[pubKey] = u
	t.mux.Unlock()
	return u, nil
}

func (t *RelayManager) GetRelay(pubKey string) *RelayConnector {

	t.mux.Lock()
	defer t.mux.Unlock()
	ret, ok := t.relayMap[pubKey]
	if ok {
		return ret
	}
	return nil
}

func (t *RelayManager) ReqEvent(url, pubKey string) error { //TODO: url

	rl := t.GetRelay(pubKey)
	if rl == nil {
		u, err := t.AddRelay(url, pubKey, "")
		if err != nil {
			return err
		}
		rl = u
	}

	return rl.ReqEvent()

}

func (t *RelayManager) CloseReq(url, pubKey string) error {

	rl := t.GetRelay(pubKey)
	if rl == nil {
		tmp, err := t.AddRelay(url, pubKey, "")
		if err != nil {
			return err
		}
		rl = tmp
	}

	return rl.CloseReq()

}
