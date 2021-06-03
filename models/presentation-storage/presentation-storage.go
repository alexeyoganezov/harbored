// Manages presentation document and metadata
package presentation_storage

import (
	"github.com/gen2brain/go-fitz"
	"harbored/config"
	"harbored/models/presentation"
	"image"
)

var dpi = 72.0

type PresentationStorage struct {
	doc               *fitz.Document
	Presentation      *presentation.Presentation
	PageCount         int
	CurrentPageNumber int
	Slides            []*image.Image
	Request           chan *presentation.Presentation
	Response          chan bool
}

func NewPresentationStorage() *PresentationStorage {
	p := PresentationStorage{}
	p.Request = make(chan *presentation.Presentation)
	p.Response = make(chan bool)
	return &p
}

func (this *PresentationStorage) Init() {
	for {
		select {
		case presentation := <-this.Request:
			this.Load(presentation)
			this.Response <- true
		}
	}
}

func (this *PresentationStorage) Load(p *presentation.Presentation) {
	doc, err := fitz.New(config.Config.StaticDir + "/" + p.Filename)
	if err != nil {
		panic(err)
	}
	this.doc = doc
	defer this.doc.Close()
	this.Presentation = p
	this.PageCount = doc.NumPage()
	this.CurrentPageNumber = 0
	this.Slides = nil
	for i := 0; i < this.PageCount; i++ {
		img, err := doc.ImageDPI(i, dpi)
		if err != nil {
			panic(err)
		}
		this.Slides = append(this.Slides, &img)
	}
}
