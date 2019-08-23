package p

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/jwt"
	"github.com/VolticFroogo/Reminder-API/model"
)

// Request is the JSON request.
type Request struct {
	Email, Password string
}

// Response is the JSON response if the function was successful.
type Response struct {
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

	query := datastore.NewQuery(model.KindUser).Filter("Email =", req.Email)

	answer := client.Run(ctx, query)

	var user model.User

	key, err := answer.Next(&user)
	if err != nil {
		if err == iterator.Done {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = jwt.LoadPrivate(PrivateKey, Password)
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
		Credentials: model.Credentials{
			Auth:    auth,
			Refresh: refresh,
		},
	}, w)
}
