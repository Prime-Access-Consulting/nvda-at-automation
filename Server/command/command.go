package command

const (
	NewSessionCommandMethod  = "session.new"
	GetSettingsCommandMethod = "nvda:settings.getSettings"
	SetSettingsCommandMethod = "nvda:settings.setSettings"
)

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

type GetSettingsCommand struct {
	ID     string                              `json:"id"`
	Method string                              `json:"method"`
	Params VendorSettingsGetSettingsParameters `json:"params"`
}

type VendorSettingsGetSettingsParameter struct {
	Name string `json:"name"`
}

type VendorSettingsGetSettingsParameters struct {
	Settings []VendorSettingsGetSettingsParameter `json:"settings"`
}

type SetSettingsCommand struct {
	ID     string                              `json:"id"`
	Method string                              `json:"method"`
	Params VendorSettingsSetSettingsParameters `json:"params"`
}

type VendorSettingsSetSettingsParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type VendorSettingsSetSettingsParameters struct {
	Settings []VendorSettingsSetSettingsParameter `json:"settings"`
}
