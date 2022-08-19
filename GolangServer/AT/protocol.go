package AT

type SessionNewCommandCapabilitiesRequest struct {
	AtName       *string `json:"atName"`
	AtVersion    *string `json:"atVersion"`
	PlatformName *string `json:"platformName"`
}

type SessionNewCommandCapabilitiesRequestParameters struct {
	AlwaysMatch *SessionNewCommandCapabilitiesRequest `json:"alwaysMatch"`
}

type SessionNewCommandParams struct {
	Capabilities SessionNewCommandCapabilitiesRequestParameters `json:"capabilities"`
}

type SessionNewCommand struct {
	Method string                  `json:"method"`
	Params SessionNewCommandParams `json:"params"`
}

type VendorSettingsGetSettingsParameter struct {
	Name string `json:"name"`
}

type VendorSettingsGetSettingsParameters []VendorSettingsGetSettingsParameter

type GetSettingsCommand struct {
	Method string                              `json:"method"`
	Params VendorSettingsGetSettingsParameters `json:"params"`
}

type AnyCommand struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type ErrorResponse struct {
	ID         *string `json:"id"`
	Error      string  `json:"error"`
	Message    string  `json:"message"`
	Stacktrace *string `json:"stacktrace,omitempty"`
}

type GetSettingsResponse struct {
	Settings Settings `json:"settings"`
}
