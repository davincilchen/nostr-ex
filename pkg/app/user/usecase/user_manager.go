package usecase

import (
	"sync"
)

type UserManager struct {
	userMap map[string]*NorstrUser //KEY: public key
	mux     sync.Mutex
}

var userManager *UserManager

func newUserManager() *UserManager {
	s := &UserManager{}
	s.userMap = make(map[string]*NorstrUser)

	return s
}

func GetUserManager() *UserManager {
	if userManager == nil {
		userManager = newUserManager()

	}
	return userManager
}

func (t *UserManager) AddDefaultListener(url, pubKey, privateKey string) (
	*NorstrUser, error) {
	u, err := t.AddUser(url, pubKey, privateKey) //TODO: graceful shutdown, delete thread
	if err != nil {
		return nil, err
	}

	u.ReqEvent()
	return u, nil

}

func (t *UserManager) AddUser(url, pubKey, privateKey string) (
	*NorstrUser, error) {

	user := t.GetUser(pubKey)
	if user != nil {
		if privateKey != "" {
			user.UpdatePrivateKey(privateKey)
		}
		return user, nil
	}

	u, err := NewNostrUser(url, pubKey, privateKey)
	if err != nil {
		return nil, err
	}

	t.mux.Lock()
	t.userMap[pubKey] = u
	t.mux.Unlock()
	return u, nil
}

func (t *UserManager) GetUser(pubKey string) *NorstrUser {

	t.mux.Lock()
	defer t.mux.Unlock()
	user, ok := t.userMap[pubKey]
	if ok {
		return user
	}
	return nil
}

func (t *UserManager) ReqEvent(url, pubKey string) error { //TODO: url

	user := t.GetUser(pubKey)
	if user == nil {
		u, err := t.AddUser(url, pubKey, "")
		if err != nil {
			return err
		}
		user = u
	}

	return user.ReqEvent()

}

func (t *UserManager) CloseReq(url, pubKey string) error {

	user := t.GetUser(pubKey)
	if user == nil {
		u, err := t.AddUser(url, pubKey, "")
		if err != nil {
			return err
		}
		user = u
	}

	return user.CloseReq()

}
