package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type platformAPI struct{}

// --------------------
// Platform API routes
// --------------------
func PlatformAPI() {
	platform := &platformAPI{}

	e.GET("/api/platform/list", platform.List)
}

// --------------------
// Platform API functions
// --------------------
func (p *platformAPI) List(c echo.Context) error {
	platform := ctr.Platform.List()
	return c.JSON(http.StatusOK, platform)
}
