package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"nostr-ex/pkg/app/session/server/session"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func SocketHandler(c *gin.Context) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err) //TODO:
	}

	s := session.NewSession(ws)
	//s.Start()

	defer func() {
		closeSocketErr := ws.Close()
		if closeSocketErr != nil {
			panic(err)
		}
	}()

	fmt.Printf("session %d is connected\n", s.ID())
	defer fmt.Printf("session %d is disconnected\n", s.ID())

	session.WaitGroup.Add(1)
	defer func() {
		if err := recover(); err != nil {
			logrus.Error("session id:", s.ID(), "", err)
			logrus.Error(string(debug.Stack()))

		}
		s.Close()
		err = session.DeleteSession(s)
		if err != nil {
			logrus.Warningf("DeleteSession failed, session: %d, err: %v",
				s.ID(), err)
		}
		session.WaitGroup.Done()
	}()

	s.Start()

}
