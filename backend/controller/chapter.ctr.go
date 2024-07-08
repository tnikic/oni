package controller

import (
	"slices"
	"sort"
	"strconv"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/source"
	"github.com/tnikic/oni/storage"
	"github.com/tnikic/oni/tools"
)

type ChapterController struct {
	storage   storage.Provider
	sourceMap map[string]source.Provider

	mc *MangaController
}

func InitChapterController(
	storage storage.Provider,
	sourceMap map[string]source.Provider,
	mc *MangaController,
) *ChapterController {
	return &ChapterController{
		storage:   storage,
		sourceMap: sourceMap,
		mc:        mc,
	}
}

func (cc *ChapterController) UpdateList(mangaEntryID int) {
	mangaEntry := cc.storage.GetMangaEntry(mangaEntryID)
	oldChapterEntries := cc.storage.GetChapterEntries(mangaEntryID)
	newChapterEntries := cc.sourceMap[mangaEntry.PlatformID].GetChapters(mangaEntry)

	for _, newChapterEntry := range newChapterEntries {
		idx := slices.IndexFunc(oldChapterEntries, func(oldChapterEntry *model.ChapterEntry) bool {
			return oldChapterEntry.ChapterUrl == newChapterEntry.ChapterUrl
		})

		if idx == -1 {
			cc.storage.StoreChapterEntry(newChapterEntry)
			continue
		}
	}
}

func (cc *ChapterController) DownloadChapters(mangaID int) {
	chapters := cc.storage.GetMangaChapters(mangaID)
	mangaEntries := cc.storage.GetMangaEntries(mangaID)

	sort.Slice(mangaEntries, func(i, j int) bool {
		return mangaEntries[i].Ranking < mangaEntries[j].Ranking
	})

	for _, chapter := range chapters {
		for _, mangaEntry := range mangaEntries {
			chapterEntry := cc.storage.GetChapterEntryForChapter(mangaEntry.ID, chapter.Number)
			if chapterEntry == nil {
				continue
			}

			if chapter.ChapterEntryID == chapterEntry.ID {
				break
			}

			cc.DownloadChapter(chapterEntry.ID)
			chapter.ChapterEntryID = chapterEntry.ID
			break
		}
	}
}

func (cc *ChapterController) DownloadChapter(chapterEntryID int) {
	chapterEntry := cc.storage.GetChapterEntry(chapterEntryID)
	mangaEntry := cc.storage.GetMangaEntry(chapterEntry.MangaEntryID)
	manga := cc.mc.Get(mangaEntry.MangaID)

	pages := cc.sourceMap[mangaEntry.PlatformID].GetPages(chapterEntry)

	for index, page := range pages {
		tools.DownloadImage(page, manga.Path+"/"+strconv.Itoa(index)+".png")
	}
}
