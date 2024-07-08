package source

import "github.com/tnikic/oni/model"

type Provider interface {
	GetList() []*model.MangaEntry
	GetChapters(manga *model.MangaEntry) []*model.ChapterEntry
	GetPages(chapter *model.ChapterEntry) []string
}
