// Stores and manages network parameters
package wifi_settings

import (
	"github.com/mitchellh/mapstructure"
)

var WifiSettings *WifiSettingsStruct

type WifiSettingsStruct struct {
	PresenterIp string `json:"presenterIp"`
}

func NewWifiSettings() *WifiSettingsStruct {
	return &WifiSettingsStruct{}
}

func (this *WifiSettingsStruct) Set(args map[string]interface{}) {
	var settings WifiSettingsStruct
	mapstructure.Decode(args, &settings)
	this.PresenterIp = settings.PresenterIp
}

func init() {
	WifiSettings = NewWifiSettings()
}
