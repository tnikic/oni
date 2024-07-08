package controller

import (
	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/source"
	"github.com/tnikic/oni/storage"
)

type MangaController struct {
	storage   storage.Provider
	sourceMap map[string]source.Provider

	cc *ChapterController
}

func InitMangaController(
	storage storage.Provider,
	sourceMap map[string]source.Provider,
) *MangaController {
	return &MangaController{
		storage:   storage,
		sourceMap: sourceMap,
	}
}

func (mc *MangaController) Add(manga *model.Manga) {
	mc.storage.StoreManga(manga)
}

func (mc *MangaController) Get(mangaID int) *model.Manga {
	manga := mc.storage.GetManga(mangaID)
	return manga
}

func (mc *MangaController) Delete(mangaID int) {
	mc.storage.DeleteManga(mangaID)
}

func (mc *MangaController) List() []*model.Manga {
	mangas := mc.storage.GetAllManga()
	return mangas
}

func (mc *MangaController) ListEntries(mangaID int) []*model.MangaEntry {
	manga := mc.Get(mangaID)
	entries := mc.storage.GetAllMangaEntries()

	var foundEntries []*model.MangaEntry

entry:
	for _, entry := range entries {
		for _, title := range manga.AllTitles {
			if title == entry.Name {
				foundEntries = append(foundEntries, entry)
				continue entry
			}
		}
	}

	return foundEntries
}
