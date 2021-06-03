package setup_files

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/sqweek/dialog"
	"harbored/i18n"
	"harbored/screens/setup-wifi"
	"harbored/services/presentations"
	"harbored/utils"
	"image/color"
	"net/url"
)

type FilesView struct {
	utils.View
}

func NewFilesView(a *fyne.App, w *fyne.Window) *FilesView {
	view := &FilesView{
		View: utils.View{
			App:    a,
			Window: w,
		},
	}
	return view
}

func (view *FilesView) Render() *fyne.Container {
	wrapper := container.NewCenter()
	content := container.NewVBox()

	// Screen title
	title := canvas.NewText(i18n.T("files:title"), color.RGBA{0, 0, 0, 255})
	title.TextSize = 36
	title.TextStyle = fyne.TextStyle{Bold: true}
	content.Add(title)

	// Select directory message
	upload := container.NewVBox()
	subtitle := canvas.NewText(i18n.T("files:setPfdDir"), color.RGBA{0, 0, 0, 255})
	upload.Add(subtitle)

	// Select button
	selectBtn := widget.NewButton(i18n.T("files:select"), func() {
		selectedPath, _ := dialog.Directory().Title(i18n.T("files:title")).Browse()
		if selectedPath != "" {
			title.Hide()
			upload.Hide()
			loader := widget.NewProgressBarInfinite()
			content.Add(loader)
			content.Refresh()
			presentations.Load(selectedPath)
			view := setup_wifi.NewWifiView(view.App, view.Window)
			view.Render()
		}
	})
	btnContainer := container.NewPadded(selectBtn)
	upload.Add(btnContainer)

	// Optimization message
	note := widget.NewLabel(i18n.T("files:optimize"))
	note.Wrapping = fyne.TextWrapWord
	note.TextStyle = fyne.TextStyle{
		Italic: false,
	}
	upload.Add(note)

	adobeUrl, _ := url.Parse("https://www.adobe.com/acrobat/online/compress-pdf.html")
	adobeLink := widget.NewHyperlink("Adobe Compress PDF", adobeUrl)
	upload.Add(adobeLink)

	ilovepdfUrl, _ := url.Parse("https://www.ilovepdf.com/ru/compress_pdf")
	ilovepdfLink := widget.NewHyperlink("iLovePDF", ilovepdfUrl)
	upload.Add(ilovepdfLink)

	content.Add(upload)
	wrapper.Add(content)

	view.El = wrapper
	window := *view.Window
	window.SetContent(view.El)
	return view.El
}
