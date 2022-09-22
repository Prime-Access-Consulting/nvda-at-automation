package command

const NewSessionCommandMethod = "session.new"

type AnyCommand struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type NewSessionCommand struct {
	ID     string                  `json:"id"`
	Method string                  `json:"method"`
	Params NewSessionCommandParams `json:"params"`
}

type NewSessionCommandParams struct {
	Capabilities NewSessionCommandCapabilitiesRequestParameters `json:"capabilities"`
}

type NewSessionCommandCapabilitiesRequest struct {
	AtName       *string `json:"atName"`
	AtVersion    *string `json:"atVersion"`
	PlatformName *string `json:"platformName"`
}

type NewSessionCommandCapabilitiesRequestParameters struct {
	AlwaysMatch *NewSessionCommandCapabilitiesRequest `json:"alwaysMatch"`
}
