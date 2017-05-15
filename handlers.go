package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func (h *Handler) getGithubRepoHealth(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	owner := c.Param("owner")
	repo := c.Param("repo")
	queryParams := c.QueryParams()

	health := &Health{
		RepositoryName: fmt.Sprintf("%s/%s", owner, repo),
		RepositoryURL:  fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		Timestamp:      time.Now(),
	}

	for queryKey, queryValues := range queryParams {
		switch queryKey {
		case "indicators":
			fmt.Println(queryValues)
		}
	}

	return c.JSON(http.StatusOK, health)
}

func (h *Handler) getIndicators(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	i := db.DB("healthyrepo").C("indicators")

	result := []Indicator{}

	err := i.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, result)
}
