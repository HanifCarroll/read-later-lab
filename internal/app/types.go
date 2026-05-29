package app

import "time"

type ReadingMode string

const (
	ModeReadFully ReadingMode = "read_fully"
	ModeSkim      ReadingMode = "skim"
	ModeReference ReadingMode = "reference"
	ModeSkip      ReadingMode = "skip"
)

type ItemStatus string

const (
	StatusInbox     ItemStatus = "inbox"
	StatusReadSoon  ItemStatus = "read_soon"
	StatusSkimLater ItemStatus = "skim_later"
	StatusReference ItemStatus = "reference"
	StatusSkipped   ItemStatus = "skipped"
	StatusArchived  ItemStatus = "archived"
)

type Item struct {
	ID              string      `json:"id"`
	URL             string      `json:"url"`
	Title           string      `json:"title"`
	Site            string      `json:"site"`
	Author          string      `json:"author,omitempty"`
	PublishedAt     string      `json:"publishedAt,omitempty"`
	SavedReason     string      `json:"savedReason"`
	Summary         string      `json:"summary"`
	ReadTimeMinutes int         `json:"readTimeMinutes"`
	RecommendedMode ReadingMode `json:"recommendedMode"`
	BestSections    []string    `json:"bestSections"`
	ReadIf          []string    `json:"readIf"`
	SkipIf          []string    `json:"skipIf"`
	Status          ItemStatus  `json:"status"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
}

type Article struct {
	URL         string
	Title       string
	Site        string
	Author      string
	PublishedAt string
	Text        string
}

type Analysis struct {
	Summary         string      `json:"summary"`
	ReadTimeMinutes int         `json:"readTimeMinutes"`
	RecommendedMode ReadingMode `json:"recommendedMode"`
	BestSections    []string    `json:"bestSections"`
	ReadIf          []string    `json:"readIf"`
	SkipIf          []string    `json:"skipIf"`
}
