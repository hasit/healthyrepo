package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	dbh := DBHandler{}
	err := dbh.initDB()
	if err != nil {
		e.Logger.Panic(err)
	}

	e.File("/favicon.ico", "images/favicon.ico")

	e.GET("/indicators", dbh.getIndicators)
	e.GET("/health/:indicator/github.com/:owner/:repo", dbh.getHealth)

	e.Logger.Fatal(e.Start(":8080"))
}
