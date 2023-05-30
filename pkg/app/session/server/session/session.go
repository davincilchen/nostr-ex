package session

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// .. //
type Session struct {
	fnOnMsg func(message []byte) error

	id    int
	conn  *websocket.Conn
	mutex sync.Mutex
}

func NewSession(conn *websocket.Conn, id int) *Session {
	return &Session{
		id:   id,
		conn: conn,
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

func (t *Session) SetOnMsgHandler(fn func(message []byte) error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.fnOnMsg = fn
}

func (t *Session) ID() int {
	return t.id
}
func (t *Session) Start() {
	//enableKeepActive //TODO:
	//trackSession(t, true) //TODO: 2023.05.30
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

// TODO: add structure
func (t *Session) msgHandle(message []byte) error {

	if t.fnOnMsg == nil {
		panic("fnOnMsg == nil")
	}
	return t.fnOnMsg(message)
}
