package storage

import "github.com/tnikic/oni/model"

type Provider interface {
	// Manage stored manga
	StoreManga(manga *model.Manga)
	GetManga(mangaID int) *model.Manga
	DeleteManga(mangaID int)

	GetAllManga() []*model.Manga

	// Manage stored chapters
	StoreChapter(chapter *model.Chapter)
	GetChapter(chapterID int) *model.Chapter
	DeleteChapter(chapterID int)

	GetMangaChapters(mangaID int) []*model.Chapter

	// Manage stored platforms
	StorePlatform(platform *model.Platform)
	GetPlatform(platformID string) *model.Platform
	DeletePlatform(platformID string)

	GetAllPlatforms() []*model.Platform

	// Manage stored manga entries
	StoreMangaEntry(manga *model.MangaEntry)
	GetMangaEntry(mangaEntryID int) *model.MangaEntry
	DeleteMangaEntry(mangaEntryID int)

	GetMangaEntries(mangaID int) []*model.MangaEntry
	GetPlatformEntries(platformID string) []*model.MangaEntry
	GetAllMangaEntries() []*model.MangaEntry

	// Manage stored chapter entries
	StoreChapterEntry(chapter *model.ChapterEntry)
	GetChapterEntry(chapterEntryID int) *model.ChapterEntry
	DeleteChapterEntry(chapterEntryID int)

	GetChapterEntryForChapter(mangaEntryID int, chapterNumber int) *model.ChapterEntry
	GetChapterEntries(mangaEntryID int) []*model.ChapterEntry
}
