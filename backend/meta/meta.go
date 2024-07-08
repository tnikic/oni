package meta

import "github.com/tnikic/oni/model"

type Provider interface {
	Search(query string) []*model.Manga
	GetManga(id string) *model.Manga
}
