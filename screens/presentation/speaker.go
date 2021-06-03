package presentation

import (
	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"harbored/i18n"
	stored_settings "harbored/models/stored-settings"
	"harbored/utils"
	"image/color"
	"math"
	"strconv"
	"strings"
	"time"
)

type SpeakerView struct {
	utils.View
	parent             *PresentationView
	ticker             *time.Ticker
	start              time.Time
	leftSideContainer  *fyne.Container
	rightSideContainer *fyne.Container
	nextSlideContainer *fyne.Container
	controls           *fyne.Container
	Loader             *fyne.Container
	Wrapper            *fyne.Container
	currentSlide       *canvas.Image
	nextSlide          *canvas.Image
	qr                 *canvas.Image
	counter            *canvas.Text
	isCloseDialogOpen  bool
}

func NewSpeakerView(a *fyne.App, w *fyne.Window, parent *PresentationView) *SpeakerView {
	view := &SpeakerView{
		View: utils.View{
			App:    a,
			Window: w,
		},
		parent: parent,
		ticker: time.NewTicker(time.Second),
		start:  time.Now(),
	}
	return view
}

func (view *SpeakerView) Render() *fyne.Container {
	// Current slide
	img := *view.parent.storage.Slides[view.parent.currentPageNumber]
	view.currentSlide = canvas.NewImageFromImage(img)
	view.currentSlide.FillMode = canvas.ImageFillContain

	// Next slide
	nextImg := *view.parent.storage.Slides[view.parent.currentPageNumber+1]
	view.nextSlide = canvas.NewImageFromImage(nextImg)
	view.nextSlide.FillMode = canvas.ImageFillContain

	// QR Code
	fyneQrImage := canvas.NewImageFromImage(view.parent.qrImage)
	fyneQrImage.SetMinSize(fyne.Size{
		Width:  256,
		Height: 256,
	})
	fyneQrImage.Hide()
	view.qr = fyneQrImage

	// Controls init
	view.controls = ControlsView(view.Window)
	view.controls.Hide()

	view.counter = canvas.NewText("00:00", color.RGBA{0, 0, 0, 255})
	view.counter.TextSize = 64
	view.counter.TextStyle = fyne.TextStyle{Monospace: true}

	resetBtn := widget.NewButton(i18n.T("presentation:resetTimer"), view.resetTimer)
	resetBtn.Importance = widget.LowImportance

	vcontainer := container.NewVBox(view.counter, resetBtn)
	ccontainer := container.NewCenter(vcontainer)

	go func() {
		for {
			select {
			case <-view.ticker.C:
				since := time.Since(view.start)
				total := int(since.Seconds())
				//hours := int(total / (60 * 60) % 24)
				minutes := int(total/60) % 60
				seconds := int(total % 60)
				strMin := strconv.Itoa(minutes)
				if minutes < 10 {
					strMin = "0" + strMin
				}
				strSec := strconv.Itoa(seconds)
				if seconds < 10 {
					strSec = "0" + strSec
				}
				view.counter.Text = strMin + ":" + strSec
				view.counter.Refresh()
			case <-view.parent.slideChanges:
				view.updateCurrentSlide()
				view.updateNextSlide()
			case <-view.parent.presentationChanges:
				view.updateCurrentSlide()
				view.updateNextSlide()
				view.resetTimer()
			}
		}
	}()

	view.nextSlideContainer = container.NewMax(view.nextSlide)
	view.rightSideContainer = container.NewGridWithRows(2, view.nextSlideContainer, ccontainer)
	view.leftSideContainer = container.NewMax(view.currentSlide, view.qr)

	view.Wrapper = fyne.NewContainerWithLayout(NewSpeakerLayout(), view.leftSideContainer, view.rightSideContainer, view.controls)

	view.Loader = container.NewCenter(widget.NewProgressBarInfinite())
	view.Loader.Hide()

	view.initialize()

	return container.NewMax(view.Wrapper, container.NewCenter(view.Loader))
}

func (view *SpeakerView) initialize() {
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
}

func (view *SpeakerView) handleKeyPress(keyname string) bool {
	window := *view.Window
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysNextSlide, keyname) {
		view.parent.handleKeyPress(keyname)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysPrevSlide, keyname) {
		view.parent.handleKeyPress(keyname)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysNextPresentation, keyname) {
		view.parent.handleKeyPress(keyname)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysPrevPresentation, keyname) {
		view.parent.handleKeyPress(keyname)
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleQR, keyname) {
		view.parent.handleKeyPress(keyname)
		if view.controls.Visible() {
			view.controls.Hide()
			view.leftSideContainer.Show()
			view.rightSideContainer.Show()
		}
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleHelp, keyname) {
		if view.controls.Visible() {
			view.controls.Hide()
			view.leftSideContainer.Show()
			view.rightSideContainer.Show()
		} else {
			view.controls.Show()
			view.leftSideContainer.Hide()
			view.rightSideContainer.Hide()
		}
		view.controls.Refresh()
		return true
	}
	if utils.StringArrayIncludes(stored_settings.StoredSettings.KeysToggleFullscreen, keyname) {
		view.parent.handleKeyPress(keyname)
		return true
	}
	if keyname == "Escape" && view.isCloseDialogOpen == false {
		view.isCloseDialogOpen = true
		d := dialog.NewConfirm(i18n.T("presentation:closeConfirmTitle"), i18n.T("presentation:closeConfirmMessage"), func(result bool) {
			if result {
				window.Close()
				parent := *view.parent
				parentWindow := *parent.Window
				parentWindow.Close()
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

func (view *SpeakerView) toggleQR() {
	if view.qr.Visible() {
		view.qr.Hide()
		view.currentSlide.Show()
	} else {
		view.qr.Show()
		view.currentSlide.Hide()
	}
	view.leftSideContainer.Refresh()
}

// Custom Fyne layout
type SpeakerLayout struct{}

func NewSpeakerLayout() *SpeakerLayout {
	//func NewSpeakerLayout() *fyne.Layout {
	return &SpeakerLayout{}
}

func (m *SpeakerLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	gutter := 15

	leftWidth := int(math.Floor(float64(size.Width)/100*75)) - (gutter * 2)
	rightWidth := size.Width - leftWidth - (gutter * 3)

	left := objects[0]

	left.Resize(fyne.Size{
		Width:  leftWidth,
		Height: size.Height,
	})
	left.Move(fyne.NewPos(gutter, 0))

	right := objects[1]
	right.Resize(fyne.Size{
		Width:  rightWidth,
		Height: size.Height,
	})
	right.Move(fyne.NewPos(leftWidth+gutter*2, 0))

	controls := objects[2]
	controls.Resize(size)
	controls.Move(fyne.NewPos(0, 0))
}

func (m *SpeakerLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Union(child.MinSize())
	}

	return minSize
}

func (view *SpeakerView) resetTimer() {
	view.start = time.Now()
	view.counter.Text = "00:00"
	view.counter.Refresh()
}

func (view *SpeakerView) updateCurrentSlide() {
	prevSlide := view.currentSlide
	//
	img := *view.parent.storage.Slides[view.parent.currentPageNumber]
	view.currentSlide = canvas.NewImageFromImage(img)
	view.currentSlide.FillMode = canvas.ImageFillContain
	//
	view.leftSideContainer.Remove(prevSlide)
	view.leftSideContainer.Add(view.currentSlide)
}

func (view *SpeakerView) updateNextSlide() {
	if view.parent.currentPageNumber == view.parent.storage.PageCount-1 {
		view.nextSlide.Hide()
		return
	} else {
		view.nextSlide.Show()
	}
	prevSlide := view.nextSlide

	nextImg := *view.parent.storage.Slides[view.parent.currentPageNumber+1]
	view.nextSlide = canvas.NewImageFromImage(nextImg)
	view.nextSlide.FillMode = canvas.ImageFillContain

	view.nextSlideContainer.Remove(prevSlide)
	view.nextSlideContainer.Add(view.nextSlide)
}
