package model

type Chapter struct {
	// Keys
	ID             int
	MangaID        int
	ChapterEntryID int
	// Data
	Number int
}

type ChapterEntry struct {
	// Keys
	ID           int
	ChapterID    int
	MangaEntryID int
	MangaID      int
	// Data
	ChapterNumber int
	ChapterUrl    string
}
