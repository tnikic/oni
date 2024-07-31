package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type metaAPI struct{}

// --------------------
// Meta API routes
// --------------------
func MetaAPI() {
	meta := &metaAPI{}

	e.GET("/api/meta/search", meta.Search)
}

// --------------------
// Meta API functions
// --------------------
func (m *metaAPI) Search(c echo.Context) error {
	query := c.QueryParam("query")
	results := ctr.Meta.Search(query)
	return c.JSON(http.StatusOK, results)
}
