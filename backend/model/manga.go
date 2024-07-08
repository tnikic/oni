package model

type Manga struct {
	// Keys
	ID int
	// Data
	Title       string
	AllTitles   []string
	Description string
	Author      string
	Artist      string
	Status      string
	Country     string
	Path        string
}

type MangaEntry struct {
	// Keys
	ID         int
	MangaID    int
	PlatformID string
	// Data
	Name         string
	MangaUrl     string
	ChapterCount int
	Ranking      int
}
