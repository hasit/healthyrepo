package main

import (
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h := Handler{}
	err := h.initDB()
	if err != nil {
		e.Logger.Panic(err)
	}

	e.File("/favicon.ico", "images/favicon.png")

	e.GET("/github.com/:owner/:repo/health", h.getGithubRepoHealth)
	e.GET("/indicators", h.getIndicators)

	e.Server.Addr = ":1323"

	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
