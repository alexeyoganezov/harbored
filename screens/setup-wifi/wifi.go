package setup_wifi

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"harbored/config"
	"harbored/i18n"
	"harbored/models/network"
	"harbored/models/wifi-settings"
	"harbored/screens/setup-controls"
	"harbored/utils"
	"image/color"
	"strings"
)

type WifiView struct {
	utils.View
	selectedIp string
}

func NewWifiView(a *fyne.App, w *fyne.Window) *WifiView {
	view := &WifiView{
		View: utils.View{
			App:    a,
			Window: w,
		},
	}
	return view
}

func (view *WifiView) Render() *fyne.Container {
	wrapper := container.NewCenter()
	content := container.NewVBox()
	content.Resize(fyne.Size{
		Width:  600,
		Height: 300,
	})

	// Title
	title := canvas.NewText(i18n.T("wifi:title"), color.RGBA{0, 0, 0, 255})
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(title)

	// Get IP-addresses
	networks := network.NetScan()
	ips := make([]string, 0)
	for _, n := range *networks {
		split := strings.Split(n.Ip, "/")
		ips = append(ips, split[0]+config.Config.ServerPort)
	}

	if len(ips) == 0 {
		// Error message
		message := canvas.NewText(i18n.T("wifi:noNetwork"), color.RGBA{0, 0, 0, 255})
		content.Add(message)

		// Refresh button
		refreshBtn := widget.NewButton(i18n.T("wifi:refresh"), func() {
			//wifi_settings.WifiSettings.NetScan()
			view.Render()
		})
		refreshBtn.Importance = widget.MediumImportance
		content.Add(refreshBtn)

		// Skip message
		note := widget.NewLabel(i18n.T("wifi:skipMessage"))
		note.Wrapping = fyne.TextWrapWord
		note.TextStyle = fyne.TextStyle{
			Italic: false,
		}
		content.Add(note)

		// Skip button
		nextBtn := widget.NewButton(i18n.T("wifi:skipBtnText"), func() {
			v := setup_controls.NewControlsView(view.App, view.Window)
			v.Render()
			//view.Window.SetContent(setup_controls.ControlsView(view.App, view.Window))
		})
		nextBtn.Importance = widget.LowImportance
		nextBtnContainer := container.NewPadded(nextBtn)
		content.Add(nextBtnContainer)
	} else {
		// IP select
		view.selectedIp = ips[0]
		ipsSelect := widget.NewSelect(ips, func(value string) {
			view.selectedIp = value
		})
		ipsSelect.SetSelectedIndex(0)
		ipSelectItem := widget.FormItem{
			Text:   i18n.T("wifi:hostAddr"),
			Widget: ipsSelect,
		}

		form := widget.NewForm()

		if len(ips) > 1 {
			form.AppendItem(&ipSelectItem)
		}

		formContainer := container.NewPadded(form)
		content.Add(formContainer)

		if len(ips) == 1 {
			view.GoToNextStep()
			return wrapper
		}

		if len(ips) > 1 {
			note := widget.NewLabel(i18n.T("wifi:selectIpMessage"))
			note.Wrapping = fyne.TextWrapWord
			content.Add(note)
		}

		// Next button
		nextBtn := widget.NewButton(i18n.T("wifi:next"), view.GoToNextStep)
		nextBtn.Importance = widget.HighImportance
		nextBtnContainer := container.NewPadded(nextBtn)
		content.Add(nextBtnContainer)
	}

	wrapper.Add(content)
	view.El = wrapper
	window := *view.Window
	window.SetContent(view.El)
	return view.El
}

func (view *WifiView) GoToNextStep() {
	values := make(map[string]interface{})
	values["presenterIp"] = strings.Split(view.selectedIp, ":")[0]
	wifi_settings.WifiSettings.Set(values)
	v := setup_controls.NewControlsView(view.App, view.Window)
	v.Render()
}
