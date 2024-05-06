package model

type Manga struct {
	Source      Source
	Title       *string
	AllTitles   []string
	Description string
	Author      string
	Artist      string
	Status      string
	Country     string
	Chapters    []Chapter
}
