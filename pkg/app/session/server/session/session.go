package session

import (
	"encoding/json"
	"fmt"
	"sync"

	"nostr-ex/pkg/models"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var id = 0
var mu sync.Mutex

func GenID() int {
	mu.Lock()
	defer mu.Unlock()
	id++
	return id
}

// .. //
type Session struct {
	id    int
	conn  *websocket.Conn
	mutex sync.Mutex

	subID  *string
	mutSub sync.RWMutex
}

func NewSession(conn *websocket.Conn) *Session {
	id := GenID()
	fmt.Println("NewSession id:", id)
	return &Session{
		id:   id,
		conn: conn,
	}
}

func (t *Session) OnEvent(fromID int, event models.Msg) error {

	subID := t.getSubID()
	if t.ID() != fromID { //不是自己
		if subID == nil { //沒訂閱
			return nil
		}
	}

	if subID == nil { //自己
		return t.WriteJson(
			[]interface{}{"EVENT", "0", event})
	}
	id := *subID
	return t.WriteJson(
		[]interface{}{"EVENT", id, event})

}

// func (t *Session) WriteMessage(messageType int, data []byte) error {
// 	t.mutex.Lock()
// 	defer t.mutex.Unlock()
// 	return t.conn.WriteMessage(messageType, data)
// }

func (t *Session) WriteMessage(data []byte) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.conn.WriteMessage(websocket.TextMessage, data)
}

func (t *Session) WriteJson(v interface{}) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.conn.WriteJSON(v)
}

func (t *Session) ID() int {
	return t.id
}
func (t *Session) Start() {

	trackSession(t, true)
	for {
		_, data, err := t.conn.ReadMessage()
		if err != nil {
			log.Infof(" %s | read err: %v", t.basicInfo(), err)
			break
		}
		log.Infof("ReadMessage %s", string(data))
		//log.Infof("ReadMessage %v", data)

		err = t.msgHandle(data)
		if err != nil {
			log.Infof(" %s | msgHandle err: %v", t.basicInfo(), err)
			break
		}

	}

	log.Infof(" %s | closed", t.basicInfo())
}

func (t *Session) Close() {
	t.conn.Close()
}

func (t *Session) basicInfo() string {
	return fmt.Sprintf("%15d", t.ID())
}

func (t *Session) msgHandle(message []byte) error {

	// Parse the message as a JSON array
	var msg []interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
		return e
	}
	// Handle each message type
	switch msg[0] {
	case "EVENT":
		// Parse the event JSON
		var event models.Msg
		jsonData, _ := json.Marshal(msg[1])
		if err := json.Unmarshal(jsonData, &event); err != nil {
			e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
			return e
		}
		// Print the event data
		fmt.Printf("Received event in session ID %d : %+v\n", t.ID(), event)

		ForEachSession(func(s SessionF) {
			s.OnEvent(t.ID(), event)
		})

	case "REQ":
		// Subscription has been closed
		fmt.Printf("Subscription %s req\n", msg[1])

		tmp, ok := msg[1].(string)
		if ok {
			t.setSubID(&tmp)
		}

		t.WriteJson([]interface{}{"EOSE", tmp})
	case "CLOSE":
		// Subscription has been closed
		fmt.Printf("Subscription %s closed\n", msg[1])
		t.setSubID(nil)
	case "EOSE":
		fmt.Printf("EOSE  \n")
	default:
		log.Printf("Unknown message type: %s\n", msg[0])
	}

	return nil
}

func (t *Session) setSubID(subID *string) {
	t.mutSub.Lock()
	defer t.mutSub.Unlock()

	t.subID = subID
}

func (t *Session) getSubID() *string {
	t.mutSub.RLock()
	defer t.mutSub.RUnlock()

	return t.subID
}

func (t *Session) IsReq() bool {
	t.mutSub.RLock()
	defer t.mutSub.RUnlock()

	return t.subID != nil
}
