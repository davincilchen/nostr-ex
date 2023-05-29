package usecase

import (
	"fmt"
	mqRepo "nostr-ex/pkg/app/rabbitmq/repo"
	"nostr-ex/pkg/app/session"
	"nostr-ex/pkg/token"
	"sync"
)

type RelayConnector struct {
	session        *session.Session
	ID             int
	SubscriptionID string
	mux            sync.Mutex
}

func NewRelayConnector(url, pubKey, privateKey string) (*RelayConnector, error) {

	s := session.NewSession(url)

	u := &RelayConnector{
		session: s,
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
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.SubscriptionID != "" {
		req := []interface{}{"REQ", t.SubscriptionID, ""}
		e := t.session.WriteJson(req)
		if e != nil {
			return e
		}
		return nil
	}

	id := token.GenUUIDv4String()
	req := []interface{}{"REQ", id, ""}

	e := t.session.WriteJson(req)

	if e != nil {
		return e
	}

	t.SubscriptionID = id
	return nil
}

func (t *RelayConnector) CloseReq() error {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.SubscriptionID == "" {
		return nil
	}

	req := []interface{}{"CLOSE", t.SubscriptionID}

	e := t.session.WriteJson(req)

	if e != nil {
		return e
	}

	t.SubscriptionID = ""

	return nil
}

func (t *RelayConnector) GetSubscriptionID() string {
	t.mux.Lock()
	defer t.mux.Unlock()
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

	// eUCase := eventUCase.NewEventHandler()
	// data := models.Event{
	// 	SubID: subID,
	// 	Data:  string(event),
	// }
	// eUCase.SaveEvent(data)
}

func (t *RelayConnector) OnConnect() {

	fmt.Printf("\nOnConnect [ID = %d] [my subID = %s]  \n",
		t.ID, t.GetSubscriptionID())

	t.ReqEvent()
}
