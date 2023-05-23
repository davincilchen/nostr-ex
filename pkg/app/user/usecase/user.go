package usecase

import (
	"fmt"
	eventUCase "nostr-ex/pkg/app/event/usecase"
	"nostr-ex/pkg/app/session"
	"nostr-ex/pkg/models"
	"nostr-ex/pkg/token"
	"sync"
)

type User struct {
	pubKey     string
	privateKey string
}

type NorstrUser struct {
	User
	session        *session.Session
	SubscriptionID string
	mux            sync.Mutex
}

func NewNostrUser(url, pubKey, privateKey string) (*NorstrUser, error) {

	s := session.NewSession(url)

	err := s.Start()
	if err != nil {
		return nil, err
	}

	u := &NorstrUser{
		User: User{
			pubKey:     pubKey,
			privateKey: privateKey,
		},
		session: s,
	}

	s.SetOnEventHandler(u.OnEvent)
	s.SetOnConnetHandler(u.OnConnect)
	return u, nil
}

func (t *NorstrUser) UpdatePrivateKey(privateKey string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.privateKey = privateKey
}

func (t *NorstrUser) PostEvent(msg string) error {
	ms := models.NewMsg(t.pubKey, msg)
	key := t.GetPrivateKey()
	event, err := ms.MakeEvent(key)
	if err != nil {
		return err
	}

	e := t.session.WriteJson(event) //TODO: e
	if e != nil {
		fmt.Println("PostEvent WriteJson Error:", e.Error())
		return e
	}

	// err = t.session.WriteJson(event)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (t *NorstrUser) ReqEvent() error {
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

func (t *NorstrUser) CloseReq() error {
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

func (t *NorstrUser) GetPrivateKey() string {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.privateKey
}

func (t *NorstrUser) GetSubscriptionID() string {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.SubscriptionID
}

func (t *NorstrUser) OnEvent(subID string, event []byte) {

	// jsonData, _ := json.Marshal(msg[2])
	// if err := json.Unmarshal(jsonData, &event); err != nil {
	// 	//if err := json.Unmarshal([]byte(msg[1].(string)), &event); err != nil {
	// 	//if err := json.Unmarshal([]byte(msg[2].(string)), &event); err != nil {
	// 	e := fmt.Errorf("Session msgHandle: json unmarshal error:%s", err.Error())
	// 	return e
	// }
	//fmt.Printf("Received event: %s\n", string(jsonData))

	fmt.Printf("\nOnEvent [my subID = %s] [my pubKey = %s] : %s\n",
		subID, t.pubKey, event)

	eUCase := eventUCase.NewEventHandler()
	data := models.Event{
		SubID: subID,
		Data:  string(event),
	}
	eUCase.SaveEvent(data)
}

func (t *NorstrUser) OnConnect() {

	fmt.Printf("\nOnConnect [my subID = %s] [my pubKey = %s] \n",
		t.GetSubscriptionID(), t.pubKey)

	t.ReqEvent()
}
