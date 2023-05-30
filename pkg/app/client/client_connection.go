package client

import (
	"encoding/json"
	"fmt"
	"log"
	"nostr-ex/pkg/models"

	eventUCase "nostr-ex/pkg/app/event/usecase"
	"nostr-ex/pkg/app/session/server/session"
)

// import (

// )

type ClientConnection struct {
	session.Session
	subID        *string
	curDBID      int
	eventHandler *eventUCase.Handler
}

func (t *ClientConnection) setSubID(subID *string) {
	t.subID = subID
}

// func (t *ClientConnection) getSubID() *string {
// 	return t.subID
// }

func (t *ClientConnection) OnSocketMsg(message []byte) error {
	// Parse the message as a JSON array
	var msg []interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		e := fmt.Errorf("OnSocketMsg: json unmarshal error:%s", err.Error())
		return e
	}
	if len(msg) < 3 {
		e := fmt.Errorf("OnSocketMsg: len(msg) <3")
		return e
	}
	// Handle each message type
	switch msg[0] {
	case "EVENT":
		fmt.Printf("Received event in session ID %d : %s %s\n", t.ID(), msg[1], msg[2])
	case "REQ":

		fmt.Printf("Subscription %s req\n", msg[1])

		tmp, ok := msg[1].(string)
		if !ok {
			t.setSubID(nil)
			t.curDBID = -1
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

func (t *ClientConnection) OnEvent(fromID int, event models.Msg) error {

	return nil
	// subID := t.getSubID()
	// if t.ID() != fromID { //不是自己
	// 	if subID == nil { //沒訂閱
	// 		return nil
	// 	}
	// }

	// if subID == nil { //自己
	// 	return t.WriteJson(
	// 		[]interface{}{"EVENT", "0", event})
	// }
	// id := *subID
	// return t.WriteJson(
	// 	[]interface{}{"EVENT", id, event})

}

// func (t *Client) OnEvent(subID string, event []byte) {

// 	// jsonData, _ := json.Marshal(msg[2])
// 	// if err := json.Unmarshal(jsonData, &event); err != nil {
// 	// 	//if err := json.Unmarshal([]byte(msg[1].(string)), &event); err != nil {
// 	// 	//if err := json.Unmarshal([]byte(msg[2].(string)), &event); err != nil {
// 	// 	e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
// 	// 	return e
// 	// }
// 	//fmt.Printf("Received event: %s\n", string(jsonData))

// 	fmt.Printf("\nOnEvent [ID = %d] [my subID = %s]  : %s\n",
// 		t.ID, subID, event) //TODO: delete

// 	mq := mqRepo.GetPubManager()
// 	mq.Send(event) //TODO: handle error

// }

func (t *ClientConnection) OnDBDone() {
	fmt.Println("================= OnDBDone =================")
}
