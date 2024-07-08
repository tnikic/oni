package controller

import (
	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/storage"
)

func InitDatabase(storage storage.Provider) {
	// Setup platform Providers
	setupPlatform(storage, "mangasee", "Mangasee")
}

func setupPlatform(storage storage.Provider, id string, name string) {
	platform := storage.GetPlatform(id)

	if platform == nil {
		platform = &model.Platform{
			ID:      id,
			Name:    name,
			Active:  true,
			Ranking: 1,
		}
		storage.StorePlatform(platform)
	}
}
