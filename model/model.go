package model

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	ProjectID = "froogo-reminder-api"

	KindUser     = "User"
	KindReminder = "Reminder"
	KindJTI      = "JTI"

	AuthDuration    = time.Hour * 24
	RefreshDuration = time.Hour * 24 * 14
)

type Reminder struct {
	Name, Description                  string `json:",omitempty"`
	Creation, Modification, Activation int64  `json:",omitempty"`
}

type ReminderWithKey struct {
	Name, Description, Key             string `json:",omitempty"`
	Creation, Modification, Activation int64  `json:",omitempty"`
}

type User struct {
	Username, Email string `json:",omitempty"`
	Password        string `json:"-"`
}

type Credentials struct {
	Auth, Refresh string `json:",omitempty"`
}

type Token struct {
	jwt.StandardClaims
}

type JTI struct {
	Expiry int64
}
