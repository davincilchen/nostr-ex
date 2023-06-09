package repo

import "fmt"

type PubManager struct {
	connector *Connector
}

var pubManager *PubManager

func Init(url, qName string) error {

	manager, err := newPubManager(url, qName)
	if err != nil {
		return err
	}
	pubManager = manager

	//TODO: 暫時做法
	m2, err := newPubManager(url, "DBEvent")
	if err != nil {
		return err
	}
	pubDBWriteDone = m2

	return nil
}

func newPubManager(url, qName string) (*PubManager, error) {
	c := NewConnector(url, qName)
	err := c.Connect()
	if err != nil {
		return nil, err
	}

	s := &PubManager{
		connector: c,
	}

	return s, nil
}

func GetPubManager() *PubManager {
	return pubManager
}

func (t *PubManager) Send(data []byte) error {
	if t.connector == nil {
		return fmt.Errorf("connector == nil")
	}
	return t.connector.Send(data)

}

func (t *PubManager) Close() error {
	if t.connector == nil {
		return fmt.Errorf("connector == nil")
	}
	t.connector.DisConnect()
	return nil

}

//TODO: 暫時做法

var pubDBWriteDone *PubManager

func GetDBPublisher() *PubManager {
	return pubDBWriteDone
}
