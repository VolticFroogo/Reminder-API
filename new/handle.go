package p

import (
	"encoding/json"
	"net/http"

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
	var req Request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		helper.ThrowErr(err, http.StatusInternalServerError, w)
		return
	}

	helper.JSONResponse(Response{
		ID: 123,
	}, w)
}
