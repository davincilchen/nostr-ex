package usecase

import (
	"fmt"
	mqRepo "nostr-ex/pkg/app/rabbitmq/repo"
	"nostr-ex/pkg/app/session"
	"nostr-ex/pkg/token"
)

type RelayConnector struct {
	session        *session.Session
	ID             int
	SubscriptionID string
}

func NewRelayConnector(url, pubKey, privateKey string) (*RelayConnector, error) {

	s := session.NewSession(url)

	u := &RelayConnector{
		SubscriptionID: token.GenUUIDv4String(),
		session:        s,
	}

	s.SetOnEventHandler(u.OnEvent)
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
