package session

import (
	"fmt"
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

	fnOnConnet func()
	fnOnMsg    func(message []byte) error

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
			log.Error("Session Start() Error:", err)
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
			log.Error("Session start() Error:", err)
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
	//TODO: add flag to check exit or not

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

	if t.fnOnMsg != nil {
		return t.fnOnMsg(message)
	}

	return nil

}

func (t *Session) SetOnMsgHandler(fn func(message []byte) error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fnOnMsg = fn
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
