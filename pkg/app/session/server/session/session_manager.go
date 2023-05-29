package session

import (
	"fmt"
	"nostr-ex/pkg/models"
	"sync"
)

type SessionF interface {
	ID() int
	Start()
	Close()
	OnEvent(fromID int, event models.Msg) error
}

var allSession map[SessionF]struct{}
var allSessionMu sync.RWMutex

func init() {
	allSession = make(map[SessionF]struct{})
}

func trackSession(s SessionF, add bool) {

	allSessionMu.Lock()
	defer allSessionMu.Unlock()

	if add {
		allSession[s] = struct{}{}
	} else {
		delete(allSession, s)
	}
}

func ForEachSession(fn func(SessionF)) {

	allSe := make(map[SessionF]struct{}) //avoid deadlock

	allSessionMu.RLock()
	for s := range allSession {
		allSe[s] = struct{}{}
	}
	allSessionMu.RUnlock()

	for s := range allSe {
		fn(s)
	}
}

func CountSession() int {
	allSessionMu.RLock()
	defer allSessionMu.RUnlock()
	return len(allSession)
}

func DeleteSession(s SessionF) error {
	fmt.Println("DeleteSession")
	trackSession(s, false)
	return nil
}
