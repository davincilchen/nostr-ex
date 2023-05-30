package client

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	eventUCase "nostr-ex/pkg/app/event/usecase"
	"nostr-ex/pkg/app/session/server/session"
	"nostr-ex/pkg/models"

	"github.com/gorilla/websocket"
)

var id = 0
var mu sync.Mutex

func GenID() int {
	mu.Lock()
	defer mu.Unlock()
	id++
	return id
}

type ClientConnection struct {
	session.Session

	mutSub sync.RWMutex
	subID  *string

	curDBID      int
	eventHandler *eventUCase.Handler
}

func NewClientConnection(conn *websocket.Conn) *ClientConnection {
	id := GenID()
	fmt.Println("NewClientConnection id:", id)
	c := &ClientConnection{
		Session:      *session.NewSession(conn, id),
		curDBID:      -1,
		eventHandler: eventUCase.NewEventHandler(),
	}

	c.SetOnMsgHandler(c.OnSocketMsg)

	return c
}

func (t *ClientConnection) setSubID(subID *string) {
	t.mutSub.Lock()
	defer t.mutSub.Unlock()

	t.subID = subID
}

func (t *ClientConnection) getSubID() *string {
	t.mutSub.RLock()
	defer t.mutSub.RUnlock()

	return t.subID
}

func (t *ClientConnection) IsReq() bool {
	t.mutSub.RLock()
	defer t.mutSub.RUnlock()

	return t.subID != nil
}

func (t *ClientConnection) OnSocketMsg(message []byte) error {
	// Parse the message as a JSON array
	var msg []interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		e := fmt.Errorf("OnSocketMsg: json unmarshal error:%s", err.Error())
		return e
	}
	if len(msg) < 1 {
		e := fmt.Errorf("OnSocketMsg: len(msg) <1")
		return e
	}
	// Handle each message type
	switch msg[0] {
	case "EVENT":
		fmt.Printf("Received event in session ID %d : %s %s\n", t.ID(), msg[1], msg[2])
	case "REQ":
		if len(msg) < 2 {
			e := fmt.Errorf("OnSocketMsg: len(msg) <2")
			return e
		}
		fmt.Printf("Subscription %s req\n", msg[1])
		t.curDBID = -1
		tmp, ok := msg[1].(string)
		if !ok {
			t.setSubID(nil)
			break
		}
		t.setSubID(&tmp)
		event := t.eventHandler.GetLastEvent()
		if event != nil {
			t.curDBID = int(event.ID)
		}

		t.WriteJson([]interface{}{"EOSE", tmp})
	case "CLOSE":
		// Subscription has been closed
		if len(msg) < 2 {
			e := fmt.Errorf("OnSocketMsg: len(msg) <2")
			return e
		}
		fmt.Printf("Subscription %s closed\n", msg[1])
		t.setSubID(nil)
		t.curDBID = -1
	case "EOSE":
		fmt.Printf("EOSE  \n")
	default:
		log.Printf("Unknown message type: %s\n", msg[0])
	}

	return nil
}

func (t *ClientConnection) OnDBDone() {
	fmt.Println("================= OnDBDone =================")
	if t.curDBID < 0 {
		return
	}
	list := t.eventHandler.GetEventFrom(t.curDBID)
	for _, v := range list {
		// Parse the event JSON
		var msg models.Msg
		if err := json.Unmarshal([]byte(v.Data), &msg); err != nil {
			//TODO:
			continue
		}

		t.WriteJson( //use routine
			[]interface{}{"EVENT", t.getSubID(), msg})
	}
}
