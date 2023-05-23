package session

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

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

type Session struct {
	id         int
	serverAddr string
	fnOnEvent  func(subID string, event []byte)
	fnOnConnet func()

	conn  *websocket.Conn
	mutex sync.Mutex
}

func NewSession(url string) *Session {

	return &Session{
		id:         GenID(),
		serverAddr: url,
	}
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

func (t *Session) Start() error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(string(debug.Stack()))
		}
	}()

	log.Infof(" %s | dial", t.basicInfo())

	conn, _, err := websocket.DefaultDialer.Dial(t.serverAddr, nil)
	if err != nil {
		log.Error("websocket dial error:", err)
		return err
	}

	log.Infof(" %s | dial success", t.basicInfo())

	t.mutex.Lock()
	t.conn = conn
	t.mutex.Unlock()

	handler := t.getOnConnetHandler()
	if handler != nil {
		handler()
	}
	go t.start()

	return nil
}

func (t *Session) start() {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			log.Error(string(debug.Stack()))
		}
	}()

	for {
		_, data, err := t.conn.ReadMessage()
		if err != nil {
			log.Infof(" %s | read err: %v", t.basicInfo(), err)
			break
		}

		fmt.Println()
		log.Infof("ReadMessage %s", string(data))
		//log.Infof("ReadMessage %v", data)

		err = t.msgHandle(data)
		if err != nil {
			log.Infof(" %s | msgHandle err: %v", t.basicInfo(), err)
			break
		}

	}

	log.Infof(" %s | closed", t.basicInfo())

	// 暫時作法
	log.Infof(" %s | retry to connect", t.basicInfo())

	for {
		err := t.Start()
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}

}

func (t *Session) Close() {
	t.conn.Close()
}

func (t *Session) basicInfo() string {
	return fmt.Sprintf("ID:%3d ,%15s", t.id, t.serverAddr)
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
		if len(msg) != 3 {
			break
		}

		//fmt.Printf("Received event: %+v\n", event)
		//fmt.Printf("Received event: %s\n", string(jsonData))
		subID, _ := msg[1].(string)
		jsonData, _ := json.Marshal(msg[2])

		handler := t.getOnEventHandler()
		if handler != nil {
			handler(subID, jsonData)
		}
	case "CLOSE":
		// Subscription has been closed
		fmt.Printf("Subscription %s closed\n", msg[1])
	case "EOSE":
		fmt.Printf("EOSE %s \n", msg[1])
	default:
		log.Printf("Unknown message type: %s\n", msg[0])
	}

	return nil
}

func (t *Session) SetOnEventHandler(fn func(subID string, event []byte)) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fnOnEvent = fn
}

func (t *Session) getOnEventHandler() func(subID string, event []byte) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.fnOnEvent
}

func (t *Session) SetOnConnetHandler(fn func()) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fnOnConnet = fn
}

func (t *Session) getOnConnetHandler() func() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.fnOnConnet

}
