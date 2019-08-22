package model

const (
	ProjectID = "froogo-reminder"

	KindReminder = "Reminder"
)

type Reminder struct {
	Name, Description                  string
	Creation, Modification, Activation int64
}
