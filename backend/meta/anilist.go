package meta

import (
	"log/slog"
	"strconv"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/tools"
)

type Anilist struct {
	url string
}

type AnilistMangaList struct {
	Data struct {
		Page struct {
			Media []struct {
				ID    int
				Title struct {
					Romaji  string
					English string
					Native  string
				}
				Synonyms []string
			}
		}
	}
}

type AnilistManga struct {
	Data struct {
		Media struct {
			ID    int
			Title struct {
				Romaji  string
				English string
				Native  string
			}
			Synonyms []string
		}
	}
}

func InitAnilist() *Anilist {
	return &Anilist{
		url: "https://graphql.anilist.co",
	}
}

func (a *Anilist) Search(manga string) []*model.Manga {
	page := 1
	query := `
        query ($search: String, $page: Int) {
            Page (page: $page, perPage: 50) {
                pageInfo {
                    total
                    currentPage
                    lastPage
                    hasNextPage
                    perPage
                }
                media (type: MANGA, search: $search) {
                    id
                    title {
                        romaji
                        english
                        native
                    }
                    synonyms
                }
            }
        }
    `

	variables := map[string]string{
		"search": manga,
		"page":   strconv.Itoa(page),
	}

	mangaList := AnilistMangaList{}

	err := tools.GQLQuery(a.url, query, variables, &mangaList)
	if err != nil {
		slog.Error("Anilist search error")
	}

	return mangaList.transform()
}

func (a *Anilist) GetManga(id string) *model.Manga {
	query := `
        query ($id: Int) {
            Media (type: MANGA, id: $id) {
                id
                title {
                    romaji
                    english
                    native
                }
                synonyms
            }
        }
    `

	variables := map[string]string{
		"id": id,
	}

	manga := AnilistManga{}

	err := tools.GQLQuery(a.url, query, variables, &manga)
	if err != nil {
		slog.Error("Anilist get manga error")
	}

	return manga.transform()
}

func (aml *AnilistMangaList) transform() []*model.Manga {
	mangas := []*model.Manga{}

	for _, manga := range aml.Data.Page.Media {
		title := manga.Title.English
		if title == "" {
			title = manga.Title.Romaji
		}

		allTitles := []string{manga.Title.English, manga.Title.Romaji, manga.Title.Native}
		allTitles = append(allTitles, manga.Synonyms...)

		mangas = append(mangas, &model.Manga{
			Title:     title,
			AllTitles: allTitles,
		})
	}
	return mangas
}

func (am *AnilistManga) transform() *model.Manga {
	title := am.Data.Media.Title.English
	if title == "" {
		title = am.Data.Media.Title.Romaji
	}

	allTitles := []string{am.Data.Media.Title.English, am.Data.Media.Title.Romaji, am.Data.Media.Title.Native}
	allTitles = append(allTitles, am.Data.Media.Synonyms...)

	manga := &model.Manga{
		Title:     title,
		AllTitles: allTitles,
	}

	return manga
}
