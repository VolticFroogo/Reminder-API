package p

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/jwt"
	"github.com/VolticFroogo/Reminder-API/middleware"
	"github.com/VolticFroogo/Reminder-API/model"
)

// Request is the JSON request.
type Request struct {
	Reminder    model.ReminderWithKey
	Credentials model.Credentials
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

	// Create a client.
	client, err := datastore.NewClient(ctx, model.ProjectID)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Load the JWT public key for checking authentication.
	err = jwt.LoadPublic(PublicKey)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Check the user's authentication.
	valid, userString, err := middleware.User(req.Credentials, client, w)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// If it's invalid, return.
	if !valid {
		return
	}

	// Decode the user key.
	user, err := datastore.DecodeKey(userString)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Decode the reminder key.
	key, err := datastore.DecodeKey(req.Reminder.Key)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Make sure the reminder they are updating is owned by them.
	if key.Parent.ID != user.ID {
		helper.ThrowErr(fmt.Errorf("reminder not owned by authorised user"), http.StatusForbidden, w)
		return
	}

	// Convert the ReminderWithKey to a Reminder.
	reminder := req.Reminder.ToReminder()

	// Save the new entity.
	_, err = client.Put(ctx, key, &reminder)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Send 200 OK response.
	w.WriteHeader(http.StatusOK)
}
