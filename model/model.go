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

// ToReminder converts a ReminderWithKey to a Reminder.
func (input ReminderWithKey) ToReminder() (output Reminder) {
	output = Reminder{
		Name:         input.Name,
		Description:  input.Description,
		Creation:     input.Creation,
		Modification: input.Modification,
		Activation:   input.Activation,
	}
}

type User struct {
	Username, Email, Password string `json:",omitempty"`
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
