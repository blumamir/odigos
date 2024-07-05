package server

type ConfigSectionName string

// This interface is the message sent to opamp client to configure aspects of the SDK
const (
	RemoteConfigSdkConfigSectionName                      ConfigSectionName = "SDK"
	RemoteConfigInstrumentationLibrariesConfigSectionName ConfigSectionName = "InstrumentationLibraries"
)

type RemoteConfigSdk struct {
	RemoteResourceAttributes []ResourceAttribute `json:"remoteResourceAttributes"`
}

type RemoteConfigInstrumentationLibrary struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
