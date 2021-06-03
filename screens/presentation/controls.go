package presentation

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"harbored/i18n"
	stored_settings "harbored/models/stored-settings"
	"harbored/utils"
	"image/color"
	"runtime"
)

func KeyItemView(key string) *fyne.Container {
	nextKeysContainer := container.NewHBox()
	el := canvas.NewText(utils.GetKeyName(key), color.RGBA{0, 0, 0, 255})
	nextKeysContainer.Add(el)
	return nextKeysContainer
}

func KeyListView(keys *[]string) *fyne.Container {
	container := container.NewVBox()
	for i := 0; i < len(*keys); i++ {
		arr := *keys
		value := arr[i]
		keyItem := KeyItemView(value)
		container.Add(keyItem)
	}
	return container
}

func ActionControlsView(w *fyne.Window, actionName string, keys *[]string) *fyne.Container {
	c := container.NewVBox()
	title := canvas.NewText(actionName, color.RGBA{0, 0, 0, 255})
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	c.Add(title)
	keyList := KeyListView(keys)
	c.Add(keyList)
	wrapper := container.NewPadded(c)
	return wrapper
}

func ControlsView(w *fyne.Window) *fyne.Container {
	wrapper := container.NewCenter()

	content := container.NewVBox()

	title := canvas.NewText(i18n.T("controls:title"), color.RGBA{0, 0, 0, 255})
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(container.NewPadded(container.NewPadded(title)))

	controls := container.NewAdaptiveGrid(3)

	// Next slide
	keysNextSlide := make([]string, len(stored_settings.StoredSettings.KeysNextSlide))
	copy(keysNextSlide, stored_settings.StoredSettings.KeysNextSlide)
	nextSlideControls := ActionControlsView(w, i18n.T("controls:nextSlide"), &keysNextSlide)
	controls.Add(nextSlideControls)

	// Prev slide
	keysPrevSlide := make([]string, len(stored_settings.StoredSettings.KeysPrevSlide))
	copy(keysPrevSlide, stored_settings.StoredSettings.KeysPrevSlide)
	prevSlideControls := ActionControlsView(w, i18n.T("controls:prevSlide"), &keysPrevSlide)
	controls.Add(prevSlideControls)

	// QR
	keysToggleQR := make([]string, len(stored_settings.StoredSettings.KeysToggleQR))
	copy(keysToggleQR, stored_settings.StoredSettings.KeysToggleQR)
	toggleQRControls := ActionControlsView(w, i18n.T("controls:showQr"), &keysToggleQR)
	controls.Add(toggleQRControls)

	// Next presentation
	keysNextPresentation := make([]string, len(stored_settings.StoredSettings.KeysNextPresentation))
	copy(keysNextPresentation, stored_settings.StoredSettings.KeysNextPresentation)
	nextPresentationControls := ActionControlsView(w, i18n.T("controls:nextPresentation"), &keysNextPresentation)
	controls.Add(nextPresentationControls)

	// Prev presentation
	keysPrevPresentation := make([]string, len(stored_settings.StoredSettings.KeysPrevPresentation))
	copy(keysPrevPresentation, stored_settings.StoredSettings.KeysPrevPresentation)
	prevPresentationControls := ActionControlsView(w, i18n.T("controls:prevPresentation"), &keysPrevPresentation)
	controls.Add(prevPresentationControls)

	// Various
	various := []string{
		"Esc - " + i18n.T("controls:closeApp"),
		"F1 - " + i18n.T("controls:showControls"),
	}
	if runtime.GOOS == "darwin" {
		various = append(various, "F12 - "+i18n.T("controls:fullscreen"))
	} else {
		various = append(various, "F11 - "+i18n.T("controls:fullscreen"))
	}
	variousControls := ActionControlsView(w, i18n.T("controls:various"), &various)
	controls.Add(variousControls)

	controlsContainer := container.NewPadded(controls)
	content.Add(controlsContainer)

	wrapper.Add(content)
	return wrapper
}
