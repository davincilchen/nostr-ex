package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	mqRepo "nostr-ex/pkg/app/rabbitmq/repo"
	"nostr-ex/pkg/app/session"
	"nostr-ex/pkg/token"
)

type Relay struct {
	ID             int    `json:"id"`
	SubscriptionID string `json:"sub_id"`
	URL            string `json:"url"`
}

type RelayConnector struct {
	Relay
	session *session.Session
}

func NewRelayConnector(id int, url string) (*RelayConnector, error) {

	s := session.NewSession(url)

	u := &RelayConnector{
		Relay: Relay{
			ID:             id,
			SubscriptionID: token.GenUUIDv4String(),
			URL:            url,
		},
		session: s,
	}

	s.SetOnMsgHandler(u.OnSocketMsg)
	s.SetOnConnetHandler(u.OnConnect)

	err := s.Start()
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (t *RelayConnector) ReqEvent() error {

	req := []interface{}{"REQ", t.GetSubscriptionID(), ""}
	e := t.session.WriteJson(req)
	if e != nil {
		return e
	}
	return nil

}

func (t *RelayConnector) GetSubscriptionID() string {
	return t.SubscriptionID
}

func (t *RelayConnector) OnSocketMsg(message []byte) error {
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
		t.OnEvent(subID, jsonData)

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

func (t *RelayConnector) OnEvent(subID string, event []byte) {

	// jsonData, _ := json.Marshal(msg[2])
	// if err := json.Unmarshal(jsonData, &event); err != nil {
	// 	//if err := json.Unmarshal([]byte(msg[1].(string)), &event); err != nil {
	// 	//if err := json.Unmarshal([]byte(msg[2].(string)), &event); err != nil {
	// 	e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
	// 	return e
	// }
	//fmt.Printf("Received event: %s\n", string(jsonData))

	fmt.Printf("\nOnEvent [ID = %d] [my subID = %s]  : %s\n",
		t.ID, subID, event) //TODO: delete

	mq := mqRepo.GetPubManager()
	mq.Send(event) //TODO: handle error

}

func (t *RelayConnector) OnConnect() {

	fmt.Printf("\nOnConnect [ID = %d] [my subID = %s]  \n",
		t.ID, t.GetSubscriptionID())

	t.ReqEvent()
}
