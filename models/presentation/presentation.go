package presentation

type Presentation struct {
	ID                int    `json:"id"`
	Filename          string `json:"filename"`
	IsOnline          bool   `json:"isOnline"`
	CurrentPageNumber int    `json:"currentPageNumber"`
}
