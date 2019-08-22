package model

type Reminder struct {
	Name, Description                  string
	Creation, Modification, Activation int64
}
