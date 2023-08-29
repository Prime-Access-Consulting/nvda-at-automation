package command

const (
	NewSessionCommandMethod           = "session.new"
	GetSettingsCommandMethod          = "nvda:settings.getSettings"
	GetSupportedSettingsCommandMethod = "nvda:settings.getSupportedSettings"
	SetSettingsCommandMethod          = "nvda:settings.setSettings"
	PressKeysCommandMethod            = "interaction.pressKeys"
)

type EmptyParams interface{}

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

type GetSupportedSettingsCommand struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params EmptyParams `json:"params"`
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

type PressKeysCommand struct {
	ID     string              `json:"id"`
	Method string              `json:"method"`
	Params PressKeysParameters `json:"params"`
}

type PressKeysParameters struct {
	Keys []string `json:"keys"`
}
