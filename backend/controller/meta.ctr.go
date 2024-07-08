package controller

import (
	"github.com/tnikic/oni/meta"
	"github.com/tnikic/oni/model"
)

type MetaController struct {
	searchMeta meta.Provider
	metaMap    map[string]meta.Provider
}

func InitMetaController(metaMap map[string]meta.Provider) *MetaController {
	mc := &MetaController{
		metaMap: metaMap,
	}

	mc.SetSearchMeta("anilist")

	return mc
}

func (mc *MetaController) Search(query string) []*model.Manga {
	mangas := mc.searchMeta.Search(query)
	return mangas
}

func (mc *MetaController) Get(id string) *model.Manga {
	manga := mc.searchMeta.GetManga(id)
	return manga
}

func (mc *MetaController) SetSearchMeta(id string) {
	mc.searchMeta = mc.metaMap[id]
}
