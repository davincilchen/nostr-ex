package session

import (
	"fmt"
	"sync"
)

type SessionF interface {
	ID() int
	Start()
	Close()
	OnDBDone()
}

var allSession map[SessionF]struct{}
var allSessionMu sync.RWMutex

func init() {
	allSession = make(map[SessionF]struct{})
}

func TrackSession(s SessionF, add bool) {

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
	TrackSession(s, false)
	return nil
}
