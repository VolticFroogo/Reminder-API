package helper

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse is the type used for error JSON responses.
type ErrorResponse struct {
	Error string
}

// ThrowErr is used for throwing errors via a JSON response.
func ThrowErr(err error, status int, w http.ResponseWriter) {
	// Set the status header of the response.
	w.WriteHeader(status)

	// Send the error as a JSON response.
	JSONResponse(ErrorResponse{
		Error: err.Error(),
	}, w)
}

// JSONResponse sends a client a JSON response.
func JSONResponse(data interface{}, w http.ResponseWriter) (err error) {
	// Encode response into JSON.
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return
	}

	// Write JSON data to response writer.
	w.Write(dataJSON)
	return
}
