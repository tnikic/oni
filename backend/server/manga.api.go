package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type mangaAPI struct{}

// --------------------
// Manga API routes
// --------------------
func MangaAPI() {
	manga := &mangaAPI{}

	e.GET("/api/manga/list", manga.List)
}

// --------------------
// Manga API functions
// --------------------
func (m *mangaAPI) List(c echo.Context) error {
	manga := ctr.Manga.List()
	return c.JSON(http.StatusOK, manga)
}
