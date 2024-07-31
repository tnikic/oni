package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/tnikic/oni/controller"
	"github.com/tnikic/oni/storage"
	"github.com/tnikic/oni/tools"
)

var ctr *controller.Controller
var e *echo.Echo

func Start(storage storage.Provider) {
	config := tools.LoadConfig()
	controller.Startup(storage)

	ctr = controller.InitController(storage)

	e = echo.New()
	e.Use(middleware.Recover())

	// Different API routes
	PlatformAPI()
	MangaAPI()
	MetaAPI()

	e.Logger.Fatal(e.Start(config.Host + ":" + config.Port))
}
