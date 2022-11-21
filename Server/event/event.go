package event

import (
	"encoding/json"
	"log"
)

const (
	InteractionCapturedOutputEventMethod = "interaction.capturedOutput"
)

type InteractionCapturedOutputParameters struct {
	SessionID string `json:"sessionid"`
	Data      string `json:"data"`
}

type InteractionCapturedOutputEvent struct {
	Method string                              `json:"method"`
	Params InteractionCapturedOutputParameters `json:"params"`
}

func InteractionCapturedOutputEventJSON(sessionId *string, output *string) []byte {
	r := InteractionCapturedOutputEvent{
		Method: InteractionCapturedOutputEventMethod,
		Params: InteractionCapturedOutputParameters{
			SessionID: *sessionId,
			Data:      *output,
		},
	}

	response, err := json.Marshal(r)

	if err != nil {
		log.Fatal(err)
	}

	return response
}
