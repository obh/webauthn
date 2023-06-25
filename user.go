package main

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	Id          string
	Username    string
	Displayname string
	Credentials []webauthn.Credential
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
	return u.Credentials
}

func (u *User) AddCredential(c *webauthn.Credential) {
	if len(u.Credentials) == 0 {
		u.Credentials = make([]webauthn.Credential, 0)
	}
	u.Credentials = append(u.Credentials, *c)
}

func (u *User) WebAuthnIcon() string {
	return ""
}
