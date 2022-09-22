package response

import (
	"encoding/json"
	"log"
)

type ErrorResponse struct {
	ID         *string `json:"id"`
	Error      string  `json:"error"`
	Message    string  `json:"message"`
	Stacktrace *string `json:"stacktrace,omitempty"`
}

func ErrorResponseJSON(error string, message string, id *string) []byte {
	c := ErrorResponse{
		ID:      id,
		Error:   error,
		Message: message,
	}

	response, err := json.Marshal(c)

	if err != nil {
		log.Fatal(err)
	}

	return response
}
