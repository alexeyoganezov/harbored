package setup_controls

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"harbored/i18n"
	stored_settings "harbored/models/stored-settings"
	"harbored/screens/presentation"
	"harbored/utils"
	"image/color"
	"runtime"
	"strings"
)

type ControlsView struct {
	utils.View
}

func NewControlsView(a *fyne.App, w *fyne.Window) *ControlsView {
	view := &ControlsView{
		View: utils.View{
			App:    a,
			Window: w,
		},
	}
	return view
}

func (view *ControlsView) Render() *fyne.Container {
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
	nextSlideControls := ActionControlsView(view.Window, i18n.T("controls:nextSlide"), &keysNextSlide, false)
	controls.Add(nextSlideControls)

	// Prev slide
	keysPrevSlide := make([]string, len(stored_settings.StoredSettings.KeysPrevSlide))
	copy(keysPrevSlide, stored_settings.StoredSettings.KeysPrevSlide)
	prevSlideControls := ActionControlsView(view.Window, i18n.T("controls:prevSlide"), &keysPrevSlide, false)
	controls.Add(prevSlideControls)

	// QR
	keysToggleQR := make([]string, len(stored_settings.StoredSettings.KeysToggleQR))
	copy(keysToggleQR, stored_settings.StoredSettings.KeysToggleQR)
	toggleQRControls := ActionControlsView(view.Window, i18n.T("controls:showQr"), &keysToggleQR, false)
	controls.Add(toggleQRControls)

	// Next presentation
	keysNextPresentation := make([]string, len(stored_settings.StoredSettings.KeysNextPresentation))
	copy(keysNextPresentation, stored_settings.StoredSettings.KeysNextPresentation)
	nextPresentationControls := ActionControlsView(view.Window, i18n.T("controls:nextPresentation"), &keysNextPresentation, false)
	controls.Add(nextPresentationControls)

	// Prev presentation
	keysPrevPresentation := make([]string, len(stored_settings.StoredSettings.KeysPrevPresentation))
	copy(keysPrevPresentation, stored_settings.StoredSettings.KeysPrevPresentation)
	prevPresentationControls := ActionControlsView(view.Window, i18n.T("controls:prevPresentation"), &keysPrevPresentation, false)
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
	variousControls := ActionControlsView(view.Window, i18n.T("controls:various"), &various, true)
	controls.Add(variousControls)

	controlsContainer := container.NewPadded(controls)
	content.Add(controlsContainer)

	btnContainer := container.NewHBox()

	resetBtn := widget.NewButton(i18n.T("controls:reset"), func() {
		stored_settings.StoredSettings.Reset()
		v := NewControlsView(view.App, view.Window)
		v.Render()
		//view.Window.SetContent(v)
	})
	resetBtn.Importance = widget.MediumImportance

	nextBtn := widget.NewButton(i18n.T("common:next"), func() {
		title.Hide()
		controls.Hide()
		btnContainer.Hide()
		loader := widget.NewProgressBarInfinite()
		content.Add(loader)
		values := make(map[string]interface{})
		values["keysNextSlide"] = keysNextSlide
		values["keysPrevSlide"] = keysPrevSlide
		values["keysToggleQR"] = keysToggleQR
		values["keysNextPresentation"] = keysNextPresentation
		values["keysPrevPresentation"] = keysPrevPresentation
		stored_settings.StoredSettings.SetControls(values)
		p := presentation.NewPresentationView(view.App, view.Window)
		p.Render()
	})
	nextBtn.Importance = widget.HighImportance

	btnContainer.Add(resetBtn)
	btnContainer.Add(layout.NewSpacer())
	btnContainer.Add(nextBtn)

	btnWrapper := container.NewPadded(btnContainer)
	content.Add(container.NewPadded(btnWrapper))

	wrapper.Add(content)

	view.El = wrapper
	window := *view.Window
	window.SetContent(view.El)
	return view.El
}

func KeyItemView(key string, keys *[]string, uncontrolled bool) *fyne.Container {
	nextKeysContainer := container.NewHBox()
	el := canvas.NewText(utils.GetKeyName(key), color.RGBA{0, 0, 0, 255})
	nextKeysContainer.Add(el)
	if !uncontrolled {
		nextKeysContainer.Add(layout.NewSpacer())
		btn := widget.NewButton(i18n.T("common:remove"), func() {
			*keys = utils.StringFilter(*keys, func(s string) bool {
				return s != key
			})
			nextKeysContainer.Hide()
		})
		btn.Importance = widget.LowImportance
		nextKeysContainer.Add(btn)
	}
	return nextKeysContainer
}

func KeyListView(keys *[]string, uncontrolled bool) *fyne.Container {
	container := container.NewVBox()
	for i := 0; i < len(*keys); i++ {
		arr := *keys
		value := arr[i]
		keyItem := KeyItemView(value, keys, uncontrolled)
		container.Add(keyItem)
	}
	return container
}

func ActionControlsView(w *fyne.Window, actionName string, keys *[]string, uncontrolled bool) *fyne.Container {
	c := container.NewVBox()
	title := canvas.NewText(actionName, color.RGBA{0, 0, 0, 255})
	title.TextSize = 20
	title.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	c.Add(title)
	keyList := KeyListView(keys, uncontrolled)
	c.Add(keyList)
	if !uncontrolled {
		addBtn := widget.NewButton(i18n.T("common:add"), func() {})
		addBtn.OnTapped = func() {
			addBtn.SetText(i18n.T("controls:pressBtn"))
			window := *w
			window.Canvas().SetOnTypedRune(func(r rune) {
				// Skip "Space" press
				if r == 32 {
					return
				}
				// Convert rune to fyne.KeyName string
				key := string(r)
				if utils.Keymap[key] != "" {
					key = utils.Keymap[key]
				} else {
					key = strings.ToUpper(key)
				}
				// Handle key press (copy-paste)
				*keys = append(*keys, key)
				keyItem := KeyItemView(key, keys, uncontrolled)
				keyList.Add(keyItem)
				window.Canvas().SetOnTypedRune(func(r rune) {})
				window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {})
				addBtn.SetText(i18n.T("common:add"))
				addBtn.Refresh()
				c.Refresh()
				keyList.Refresh()
			})
			window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
				key := string(e.Name)
				if !utils.StringArrayIncludes(utils.NonprintableKeys, key) {
					return
				}
				// Handle key press (copy-paste)
				*keys = append(*keys, key)
				keyItem := KeyItemView(key, keys, uncontrolled)
				keyList.Add(keyItem)
				window.Canvas().SetOnTypedRune(func(r rune) {})
				window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {})
				addBtn.SetText(i18n.T("common:add"))
				addBtn.Refresh()
				c.Refresh()
				keyList.Refresh()
			})
		}
		c.Add(addBtn)
	}
	wrapper := container.NewPadded(c)
	return wrapper
}
