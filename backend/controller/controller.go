package controller

import (
	"log/slog"

	"github.com/tnikic/oni/meta"
	"github.com/tnikic/oni/source"
	"github.com/tnikic/oni/storage"
)

type Controller struct {
	Manga    *MangaController
	Chapter  *ChapterController
	Platform *PlatformController
	Meta     *MetaController
}

func InitController(storage storage.Provider) *Controller {
	// Meta Providers
	slog.Info("Initializing Meta Providers")
	metaMap := make(map[string]meta.Provider)

	metaMap["anilist"] = meta.InitAnilist()

	// Source Providers
	slog.Info("Initializing Source Providers")
	sourceMap := make(map[string]source.Provider)

	sourceMap["mangasee"] = source.InitMangasee()

	// Controller
	mangaController := InitMangaController(storage, sourceMap)
	chapterController := InitChapterController(storage, sourceMap, mangaController)
	platformController := InitPlatformController(storage, sourceMap, mangaController)
	metaController := InitMetaController(metaMap)

	// Update all on startup
	slog.Info("Updating all sources on startup")
	go platformController.UpdateAll()

	return &Controller{
		Manga:    mangaController,
		Chapter:  chapterController,
		Platform: platformController,
		Meta:     metaController,
	}
}
