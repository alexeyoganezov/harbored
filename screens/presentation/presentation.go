package presentation

import (
  "fyne.io/fyne"
  "fyne.io/fyne/canvas"
  "fyne.io/fyne/container"
  "fyne.io/fyne/dialog"
  "fyne.io/fyne/widget"
  "github.com/skip2/go-qrcode"
  "harbored/config"
  "harbored/i18n"
  "harbored/models/presentation"
  "harbored/models/presentation-storage"
  stored_settings "harbored/models/stored-settings"
  "harbored/models/wifi-settings"
  "harbored/services/presentations"
  "harbored/utils"
  "image"
  "strings"
)

type PresentationView struct {
	utils.View
	currentPageNumber   int
	slides              []*canvas.Image
	currentSlide        *canvas.Image
	presentations       []*presentation.Presentation
	currentPresentation *presentation.Presentation
	imageContainer      *fyne.Container
	splashscreen        *fyne.Container
	controls            *fyne.Container
	isCloseDialogOpen   bool
	storage             *presentation_storage.PresentationStorage
	slideChanges        chan bool
	presentationChanges chan bool
	qrImage             image.Image
	speakerView         *SpeakerView
}

func NewPresentationView(a *fyne.App, w *fyne.Window) *PresentationView {
	view := &PresentationView{
		View: utils.View{
			App:    a,
			Window: w,
		},
		slides:              make([]*canvas.Image, 0),
		storage:             presentation_storage.NewPresentationStorage(),
		slideChanges:        make(chan bool),
		presentationChanges: make(chan bool),
	}
	go view.storage.Init()
	return view
}

func (view *PresentationView) Render() *fyne.Container {
	app := *view.App
	window := *view.Window
	window.SetTitle("Audience view")
	window.SetPadded(false)

	wrapper := container.NewCenter()

	// Splash init
	view.splashscreen = getSplashscreen(view.Window, view)
	view.splashscreen.Hide()
	wrapper.Add(view.splashscreen)

	// Controls init
	view.controls = ControlsView(view.Window)
	view.controls.Hide()
	wrapper.Add(view.controls)

	// Loader init
	view.imageContainer = container.NewCenter()
	loader := widget.NewProgressBarInfinite()
	view.imageContainer.Add(loader)
	wrapper.Add(view.imageContainer)

	// Render content
	view.El = wrapper
	window.SetContent(view.El)

	// Load and show a slide
	view.presentations = presentations.Get()
	view.currentPresentation = view.presentations[0]
	view.storage.Request <- view.currentPresentation
	<-view.storage.Response
	view.prepareSlide()
	view.imageContainer.Remove(loader)
	view.imageContainer.Add(view.currentSlide)

	// Speaker screen
	s := app.NewWindow("Speaker view")
	s.SetPadded(false)
	s.Resize(fyne.NewSize(1024, 768))
	s.Show()
	view.speakerView = NewSpeakerView(view.App, &s, view)
	s.SetContent(view.speakerView.Render())

	view.initialize()

	return view.El
}

func getSplashscreen(w *fyne.Window, view *PresentationView) *fyne.Container {
	window := *w
	content := container.NewCenter()
	qr, _ := qrcode.New("http://"+wifi_settings.WifiSettings.PresenterIp+config.Config.ServerPort, qrcode.Highest)
	qrHeight := window.Canvas().Size().Height
	view.qrImage = qr.Image(qrHeight)
	fyneQrImage := canvas.NewImageFromImage(view.qrImage)
	fyneQrImage.SetMinSize(fyne.Size{
		Width:  qrHeight,
		Height: qrHeight,
	})
	content.Add(fyneQrImage)
	return content
}

func (view *PresentationView) prepareSlide() {
	window := *view.Window
	img := *view.storage.Slides[view.currentPageNumber]
	fyneImage := canvas.NewImageFromImage(img)
	imageBounds := img.Bounds()
	canvasSize := window.Canvas().Size()
	var width int
	var height int
	if imageBounds.Max.X > canvasSize.Width {
		width = canvasSize.Width
		height = canvasSize.Height
	} else {
		width = imageBounds.Max.X
		height = imageBounds.Max.Y
	}
	fyneImage.SetMinSize(fyne.Size{
		Width:  width,
		Height: height,
	})
	fyneImage.FillMode = canvas.ImageFillContain
	view.currentSlide = fyneImage
}

func (view *PresentationView) initialize() {
	window := *view.Window
	window.Canvas().SetOnTypedRune(func(r rune) {
		// Skip "Space" press
		if r == 32 {
			return
		}
		// Convert rune to fyne.KeyName string
		keyname := string(r)
		if utils.Keymap[keyname] != "" {
			keyname = utils.Keymap[keyname]
		} else {
			keyname = strings.ToUpper(keyname)
		}
		// Handle key press
		view.handleKeyPress(keyname)
	})
	window.Canvas().SetOnTypedKey(func(e *fyne.KeyEvent) {
		keyname := string(e.Name)
		if utils.StringArrayIncludes(utils.NonprintableKeys, keyname) {
			view.handleKeyPress(keyname)
		}
	})
	presentations.Start(view.currentPresentation.ID)
}

func (view *PresentationView) handleKeyPress(keyname string) bool {
	window := *view.Window
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysNextSlide, keyname) {
		if view.currentPageNumber < view.storage.PageCount-1 {
			view.showNextPage()
		}
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysPrevSlide, keyname) {
		if view.currentPageNumber > 0 {
			view.showPrevPage()
		}
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysNextPresentation, keyname) {
		presentations.Stop(view.currentPresentation.ID)
		nextPresentationIndex := view.currentPresentation.ID
		if nextPresentationIndex >= len(view.presentations) {
			nextPresentationIndex = 0
		}
		view.imageContainer.Remove(view.currentSlide)
		loader := widget.NewProgressBarInfinite()
		view.imageContainer.Add(loader)
		view.speakerView.Wrapper.Hide()
		view.speakerView.Loader.Show()
		view.storage.Request <- view.presentations[nextPresentationIndex]
		<-view.storage.Response
		view.currentPageNumber = 0
		view.currentPresentation = view.storage.Presentation
		view.presentationChanges <- true
		view.prepareSlide()
		view.imageContainer.Remove(loader)
		view.speakerView.Wrapper.Show()
		view.speakerView.Loader.Hide()
		view.imageContainer.Add(view.currentSlide)
		view.imageContainer.Refresh()
		presentations.Start(view.currentPresentation.ID)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysPrevPresentation, keyname) {
		presentations.Stop(view.currentPresentation.ID)
		nextPresentationIndex := view.currentPresentation.ID - 2
		if nextPresentationIndex < 0 {
			nextPresentationIndex = len(view.presentations) - 1
		}
		view.imageContainer.Remove(view.currentSlide)
		loader := widget.NewProgressBarInfinite()
		view.imageContainer.Add(loader)
		view.speakerView.Wrapper.Hide()
		view.speakerView.Loader.Show()
		view.storage.Request <- view.presentations[nextPresentationIndex]
		<-view.storage.Response
		view.currentPageNumber = 0
		view.currentPresentation = view.storage.Presentation
		view.presentationChanges <- true
		view.prepareSlide()
		view.imageContainer.Remove(loader)
		view.speakerView.Wrapper.Show()
		view.speakerView.Loader.Hide()
		view.imageContainer.Add(view.currentSlide)
		view.imageContainer.Refresh()
		presentations.Start(view.currentPresentation.ID)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleQR, keyname) {
		if view.splashscreen.Visible() {
			view.splashscreen.Hide()
			view.imageContainer.Show()
		} else {
			view.controls.Hide()
			view.splashscreen.Show()
			view.imageContainer.Hide()
		}
		view.imageContainer.Refresh()
		view.splashscreen.Refresh()
		view.speakerView.toggleQR()
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleHelp, keyname) {
		if view.controls.Visible() {
			view.controls.Hide()
			view.imageContainer.Show()
		} else {
			view.splashscreen.Hide()
			if view.speakerView.qr.Visible() {
				view.speakerView.toggleQR()
			}
			view.controls.Show()
			view.imageContainer.Hide()
		}
		view.imageContainer.Refresh()
		view.controls.Refresh()
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleFullscreen, keyname) {
		window.SetFullScreen(!window.FullScreen())
		return true
	}
	if keyname == "Escape" && view.isCloseDialogOpen == false {
		view.isCloseDialogOpen = true
		d := dialog.NewConfirm(i18n.T("presentation:closeConfirmTitle"), i18n.T("presentation:closeConfirmMessage"), func(result bool) {
			if result {
				window.Close()
			} else {
				view.isCloseDialogOpen = false
			}
		}, *view.Window)
		d.SetConfirmText(i18n.T("common:close"))
		d.SetDismissText(i18n.T("common:cancel"))
		d.Show()
		return true
	}
	return false
}

func (view *PresentationView) showNextPage() {
	prevSlide := view.currentSlide
	view.currentPageNumber += 1
	view.slideChanges <- true
	view.prepareSlide()
	view.imageContainer.Remove(prevSlide)
	view.imageContainer.Add(view.currentSlide)
	view.imageContainer.Refresh()

	presentations.ChangeSlide(view.currentPresentation.ID, view.currentPageNumber)
}

func (view *PresentationView) showPrevPage() {
	prevSlide := view.currentSlide
	view.currentPageNumber -= 1
	view.slideChanges <- true
	view.prepareSlide()
	view.imageContainer.Remove(prevSlide)
	view.imageContainer.Add(view.currentSlide)
	view.imageContainer.Refresh()

	presentations.ChangeSlide(view.currentPresentation.ID, view.currentPageNumber)
}
