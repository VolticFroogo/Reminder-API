package p

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/jwt"
	"github.com/VolticFroogo/Reminder-API/middleware"
	"github.com/VolticFroogo/Reminder-API/model"
)

// Request is the JSON request.
type Request struct {
	Reminder    model.Reminder
	Credentials model.Credentials
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

	// Create a client.
	client, err := datastore.NewClient(ctx, model.ProjectID)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	err = jwt.LoadPublic(PublicKey)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	valid, userString, err := middleware.User(req.Credentials, client, w)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	if !valid {
		return
	}

	user, err := datastore.DecodeKey(userString)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Create a key.
	key := datastore.IncompleteKey(model.KindReminder, user)

	// Save the new entity.
	key, err = client.Put(ctx, key, &req.Reminder)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Send the JSON response.
	helper.JSONResponse(Response{
		ID: key.ID,
	}, w)
}
