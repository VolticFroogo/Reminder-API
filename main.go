package main

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type User struct {
	Username, Email, Password string
}

func main() {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, "froogo-reminder")
	if err != nil {
		log.Println(err)
		return
	}

	query := client.Collection("User").Where("Email", "==", "harry@froogo.co.uk")

	answer := query.Documents(ctx)

	iter, err := answer.Next()
	if err != nil {
		log.Println(err)
		return
	}

	var user User

	iter.DataTo(&user)

	log.Println(user)
}
