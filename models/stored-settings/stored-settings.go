// Manages app configuration file with settings.
// It stores only controls right now.
package stored_settings

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shibukawa/configdir"
)

var StoredSettings *StoredSettingsStruct
var SettingsFolder *configdir.Config

type StoredSettingsStruct struct {
	KeysNextSlide        []string `json:"keysNextSlide"`
	KeysPrevSlide        []string `json:"keysPrevSlide"`
	KeysToggleQR         []string `json:"keysToggleQR"`
	KeysNextPresentation []string `json:"keysNextPresentation"`
	KeysPrevPresentation []string `json:"keysPrevPresentation"`
	KeysToggleHelp       []string `json:"keysToggleHelp"`
	KeysToggleFullscreen []string `json:"keysToggleFullscreen"`
}

func (s *StoredSettingsStruct) SetControls(args map[string]interface{}) {
	var settings StoredSettingsStruct
	mapstructure.Decode(args, &settings)
	StoredSettings.KeysNextSlide = settings.KeysNextSlide
	StoredSettings.KeysPrevSlide = settings.KeysPrevSlide
	StoredSettings.KeysToggleQR = settings.KeysToggleQR
	StoredSettings.KeysNextPresentation = settings.KeysNextPresentation
	StoredSettings.KeysPrevPresentation = settings.KeysPrevPresentation
	data, _ := json.Marshal(&StoredSettings)
	err := SettingsFolder.WriteFile("settings.json", data)
	if err != nil {
		fmt.Println("err", err)
	}
}

func (s *StoredSettingsStruct) Reset() {
	settings := NewStoredSettings()
	StoredSettings.KeysNextSlide = settings.KeysNextSlide
	StoredSettings.KeysPrevSlide = settings.KeysPrevSlide
	StoredSettings.KeysToggleQR = settings.KeysToggleQR
	StoredSettings.KeysNextPresentation = settings.KeysNextPresentation
	StoredSettings.KeysPrevPresentation = settings.KeysPrevPresentation
	data, _ := json.Marshal(&StoredSettings)
	err := SettingsFolder.WriteFile("settings.json", data)
	if err != nil {
		fmt.Println("err", err)
	}
}

func NewStoredSettings() *StoredSettingsStruct {
	return &StoredSettingsStruct{
		KeysNextSlide: []string{
			"Right", // ArrowRight
			"Down",  // ArrowDown
			"Next",  // PageDown
			"N",     // KeyN
			"Space",
		},
		KeysPrevSlide: []string{
			"Left",      // ArrowLeft
			"Up",        // ArrowUp
			"Prior",     // PageUp
			"P",         // KeyP
			"BackSpace", // Backspace
		},
		KeysToggleQR: []string{
			"Q",      // KeyQ
			"B",      // KeyB
			"Window", // KeyW
			".",      // Period
			",",      // Comma
		},
		KeysNextPresentation: []string{
			"RightShift",
			"End",
		},
		KeysPrevPresentation: []string{
			"LeftShift",
			"Home",
		},
		KeysToggleHelp: []string{
			"F1",
		},
		KeysToggleFullscreen: []string{
			"F11",
			"F12",
		},
	}
}

func init() {
	configDirs := configdir.New("harbored", "harbored")
	folders := configDirs.QueryFolders(configdir.Global)
	SettingsFolder = folders[0]
	if SettingsFolder.Exists("settings.json") {
		data, err := SettingsFolder.ReadFile("settings.json")
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(data, &StoredSettings)
	} else {
		StoredSettings = NewStoredSettings()
		data, _ := json.Marshal(&StoredSettings)
		err := SettingsFolder.WriteFile("settings.json", data)
		if err != nil {
			fmt.Println("err", err)
		}
	}
}
