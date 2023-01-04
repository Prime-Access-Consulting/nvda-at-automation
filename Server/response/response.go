package response

import (
	"Server/client"
	"encoding/json"
	"log"
)

type ErrorResponse struct {
	ID         *string `json:"id"`
	Error      string  `json:"error"`
	Message    string  `json:"message"`
	Stacktrace *string `json:"stacktrace,omitempty"`
}

type NewSessionResponse struct {
	ID     *string          `json:"id"`
	Result NewSessionResult `json:"result"`
}

type NewSessionResult struct {
	SessionID    string                 `json:"sessionId"`
	Capabilities NewSessionCapabilities `json:"capabilities"`
}

type NewSessionCapabilities struct {
	ATName       string `json:"atName"`
	ATVersion    string `json:"atVersion"`
	PlatformName string `json:"platformName"`
}

type GetSettingsResponse struct {
	ID     string            `json:"id"`
	Result RetrievedSettings `json:"result"`
}

type RetrievedSettings []RetrievedSetting

type RetrievedSetting struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type SetSettingsResponse struct {
	ID     string      `json:"id"`
	Result EmptyResult `json:"result"`
}

type PressKeysResponse struct {
	ID     string      `json:"id"`
	Result EmptyResult `json:"result"`
}

type EmptyResult interface{}

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

func NewSessionResponseJSON(info *client.Capabilities, sessionKey string) []byte {
	r := NewSessionResponse{
		Result: NewSessionResult{
			SessionID: sessionKey,
			Capabilities: NewSessionCapabilities{
				ATName:       info.Name,
				ATVersion:    info.Version,
				PlatformName: info.Platform,
			},
		},
	}

	response, err := json.Marshal(r)

	if err != nil {
		log.Fatal(err)
	}

	return response
}
