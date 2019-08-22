package p

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/model"
)

// Request is the JSON request.
type Request struct {
	Reminder model.Reminder
}

// Response is the JSON response if the function was successful.
type Response struct {
	ID int64
}

// Handle is the first function called handling the HTTP request.
func Handle(w http.ResponseWriter, r *http.Request) {
	// Create a variable to store the request when decoded.
	var req Request

	// Decode the body into the request.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Get context.
	ctx := context.Background()

	// Creates a client.
	client, err := datastore.NewClient(ctx, model.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Creates a Key instance.
	key := datastore.IncompleteKey(model.KindReminder, nil)

	// Saves the new entity.
	key, err = client.Put(ctx, key, &req.Reminder)
	if err != nil {
		log.Fatalf("Failed to save forum: %v", err)
	}

	// Send the JSON response.
	helper.JSONResponse(Response{
		ID: key.ID,
	}, w)
}
