package p

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/api/iterator"

	firebase "firebase.google.com/go"

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

	// Get context and config.
	ctx := context.Background()

	// Create the Firebase app.
	app, err := firebase.NewApp(ctx, model.FirebaseConfig)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	// Create client.
	client, err := app.Firestore(ctx)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	defer client.Close()

	query := client.Collection(model.KindUser).Where("Email", "==", req.Email).Documents(ctx)

	answer, err := query.Next()
	if err != nil {
		if err == iterator.Done {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	var user model.User
	err = answer.DataTo(&user)
	if err != nil {
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

	auth, refresh, err := jwt.NewTokens(answer.Ref.ID, client)
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
