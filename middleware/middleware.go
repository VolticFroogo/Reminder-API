package middleware

import (
	"net/http"

	"cloud.google.com/go/datastore"

	"github.com/VolticFroogo/Reminder-API/helper"
	"github.com/VolticFroogo/Reminder-API/jwt"
	"github.com/VolticFroogo/Reminder-API/model"
)

type response struct {
	Credentials model.Credentials
}

func User(credentials model.Credentials, client *datastore.Client, w http.ResponseWriter) (valid bool, user string, err error) {
	valid, user, err = jwt.CheckAuth(credentials.Auth)
	if err != nil || valid {
		return
	}

	valid, user, err = jwt.CheckRefresh(credentials.Refresh, client)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusUnauthorized)

	if !valid {
		return
	}

	key, err := datastore.DecodeKey(user)
	if err != nil {
		return
	}

	auth, refresh, err := jwt.NewTokens(key, client)
	if err != nil {
		return
	}

	helper.JSONResponse(response{
		Credentials: model.Credentials{
			Auth:    auth,
			Refresh: refresh,
		},
	}, w)

	valid = false
	return
}
