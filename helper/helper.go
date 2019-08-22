package helper

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

func ThrowErr(err error, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	JSONResponse(ErrorResponse{
		Error: err.Error(),
	}, w)
}

// JSONResponse sends a client a JSON response.
func JSONResponse(data interface{}, w http.ResponseWriter) (err error) {
	dataJSON, err := json.Marshal(data) // Encode response into JSON.
	if err != nil {
		return
	}
	w.Write(dataJSON) // Write JSON data to response writer.
	return
}
