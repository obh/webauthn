package main

import (
	"errors"
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
)

type UserId string

type Datastore struct {
	users    map[UserId]User
	sessions map[UserId]*webauthn.SessionData
}

func InitData() *Datastore {
	u := &User{Id: "101", Username: "rohit", Displayname: "Rohit Sharma"}
	users := make(map[UserId]User)

	sessions := make(map[UserId]*webauthn.SessionData)
	users[UserId(u.Id)] = *u

	return &Datastore{
		users:    users,
		sessions: sessions,
	}
}

func (d *Datastore) GetUser(userId UserId) (*User, error) {
	u, ok := d.users[userId]
	if !ok {
		return nil, errors.New("user not found!")
	}
	return &u, nil
}

func (d *Datastore) StoreSession(userId UserId, session *webauthn.SessionData) {
	d.sessions[userId] = session
}

func (d *Datastore) GetSession(userId UserId) (*webauthn.SessionData, error) {
	s, ok := d.sessions[userId]
	if !ok {
		return nil, errors.New("session not found!")
	}
	return s, nil
}

func (d *Datastore) SaveUser(user User) {
	fmt.Println("storing User permanently: ", user)
}
