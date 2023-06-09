package main

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	Id          string
	Username    string
	Displayname string
}

func (u *User) WebAuthnID() []byte {
	return []byte(u.Id)
}

func (u *User) WebAuthnName() string {
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	return u.Displayname
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return nil
}

func (u *User) WebAuthnIcon() string {
	return ""
}
