package model

import (
	"time"

	firebase "firebase.google.com/go"
	"github.com/dgrijalva/jwt-go"
)

const (
	ProjectID = "froogo-reminder"

	KindUser     = "User"
	KindReminder = "Reminder"
	KindJTI      = "JTI"

	AuthDuration    = time.Hour * 24
	RefreshDuration = time.Hour * 24 * 14
)

var (
	FirebaseConfig = &firebase.Config{ProjectID: ProjectID}
)

type Reminder struct {
	Name, Description string `json:",omitempty"`
	Activation        int64  `json:",omitempty"`
}

type ReminderWithKey struct {
	Name, Description, Key string `json:",omitempty"`
	Activation             int64  `json:",omitempty"`
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
