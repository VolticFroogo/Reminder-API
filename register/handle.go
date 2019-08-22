package p

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/jwt"
	"github.com/VolticFroogo/Reminder-API/model"
)

// Request is the JSON request.
type Request struct {
	User model.User
}

// Response is the JSON response if the function was successful.
type Response struct {
	Auth, Refresh string
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

	// Create a key.
	key := datastore.IncompleteKey(model.KindUser, nil)

	// Hash the user's password before storing it in Datastore.
	req.User.Password, err = helper.HashPassword(req.User.Password)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Save the new entity.
	key, err = client.Put(ctx, key, &req.User)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	err = jwt.LoadPrivateKey()
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	auth, refresh, err := jwt.NewTokens(key, client)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Send the JSON response.
	helper.JSONResponse(Response{
		Auth:    auth,
		Refresh: refresh,
	}, w)
}
