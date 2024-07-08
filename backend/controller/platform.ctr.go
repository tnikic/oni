package controller

import (
	"slices"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/source"
	"github.com/tnikic/oni/storage"
)

type PlatformController struct {
	storage   storage.Provider
	sourceMap map[string]source.Provider

	mc *MangaController
}

func InitPlatformController(
	storage storage.Provider,
	sourceMap map[string]source.Provider,
	mc *MangaController,
) *PlatformController {
	return &PlatformController{
		storage:   storage,
		sourceMap: sourceMap,
		mc:        mc,
	}
}

func (pc *PlatformController) Get(platformID string) *model.Platform {
	platform := pc.storage.GetPlatform(platformID)
	return platform
}

func (pc *PlatformController) List() []*model.Platform {
	platforms := pc.storage.GetAllPlatforms()
	return platforms
}

func (pc *PlatformController) ListActive() []*model.Platform {
	platforms := pc.storage.GetAllPlatforms()

	for _, platform := range platforms {
		if platform.Active {
			platforms = append(platforms, platform)
		}
	}

	return platforms
}

func (pc *PlatformController) Toggle(platformID string) *model.Platform {
	platform := pc.Get(platformID)
	platform.Active = !platform.Active

	pc.storage.StorePlatform(platform)
	return platform
}

func (pc *PlatformController) Update(platformID string) {
	list := pc.sourceMap[platformID].GetList()
	mangaEntries := pc.storage.GetPlatformEntries(platformID)

	for _, newEntry := range list {
		idx := slices.IndexFunc(mangaEntries, func(oldEntry *model.MangaEntry) bool {
			return oldEntry.MangaUrl == newEntry.MangaUrl
		})

		if idx == -1 {
			pc.storage.StoreMangaEntry(newEntry)
			continue
		}
	}
}

func (pc *PlatformController) UpdateAll() {
	activePlatforms := pc.ListActive()

	for _, platform := range activePlatforms {
		pc.Update(platform.ID)
	}
}

func (pc *PlatformController) ListSources(manga *model.Manga) []*model.MangaEntry {
	entries := pc.storage.GetAllMangaEntries()
	var foundSources []*model.MangaEntry

entry:
	for _, entry := range entries {
		for _, title := range manga.AllTitles {
			if title == entry.Name {
				foundSources = append(foundSources, entry)
				continue entry
			}
		}
	}

	return foundSources
}
