package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type chapterAPI struct{}

// --------------------
// Chapter API routes
// --------------------
func ChapterAPI() {
	chapter := &chapterAPI{}

	e.GET("/api/chapter/update/:id", chapter.Update)
}

// --------------------
// Chapter API functions
// --------------------
func (chapter *chapterAPI) Update(c echo.Context) error {
	param := c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID")
	}
	ctr.Chapter.UpdateList(id)
	return c.NoContent(http.StatusOK)
}
