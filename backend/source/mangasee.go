package source

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-rod/rod"

	"github.com/tnikic/oni/model"
	"github.com/tnikic/oni/tools"
)

type Mangasee struct {
	id      string
	mainUrl string
	listUrl string
}

func InitMangasee() *Mangasee {
	return &Mangasee{
		id:      "mangasee",
		mainUrl: "https://mangasee123.com",
		listUrl: "https://mangasee123.com/directory/",
	}
}

func (m *Mangasee) GetList() []*model.MangaEntry {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(m.listUrl).MustWaitStable()
	elements := page.MustElementsX("//div[@class='top-15 ng-scope']/a")
	if elements.Empty() {
		slog.Error("No elements found on mangasee list page")
	}

	foundMangas := []*model.MangaEntry{}
	for _, element := range elements {
		manga := model.MangaEntry{
			PlatformID: m.id,
			Name:       element.MustText(),
			MangaUrl:   element.MustProperty("href").String(),
		}

		foundMangas = append(foundMangas, &manga)
	}
	return foundMangas
}

func (m *Mangasee) GetChapters(manga *model.MangaEntry) []*model.ChapterEntry {
	mangaRSSUrl := strings.Replace(manga.MangaUrl, "/manga/", "/rss/", 1)
	mangaRSSUrl += ".xml"

	var chapterEntries []*model.ChapterEntry
	var chapters xmlChapters
	tools.XMLUnmarshal(mangaRSSUrl, &chapters)

	for _, chapter := range chapters.Chapters {
		words := strings.Fields(chapter.Title)
		number, _ := strconv.Atoi(words[len(words)-1])

		link := strings.Replace(chapter.Link, "-page-1", "", 1)

		chapterEntry := model.ChapterEntry{
			MangaEntryID:  manga.ID,
			ChapterUrl:    link,
			ChapterNumber: number,
		}

		chapterEntries = append(chapterEntries, &chapterEntry)
	}

	return chapterEntries
}

func (m *Mangasee) GetPages(chapter *model.ChapterEntry) []string {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(chapter.ChapterUrl).MustWaitStable()
	elements := page.MustElementsX("//div[@class='ImageGallery']/div/div/img")
	if elements.Empty() {
		slog.Error("No elements found on mangasee chapter page")
	}

	urls := []string{}
	for _, element := range elements {
		pageUrl := element.MustProperty("src").String()
		urls = append(urls, pageUrl)
	}

	return urls
}

type xmlChapters struct {
	Chapters []xmlChapter `xml:"channel>item"`
}

type xmlChapter struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}
