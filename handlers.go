package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/labstack/echo"
)

func (h *Handler) getGithubRepoHealth(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	owner := c.Param("owner")
	repo := c.Param("repo")
	indicators := c.QueryParam("indicators")

	health := &Health{
		RepositoryName:     repo,
		RepositoryOwner:    owner,
		RepositoryFullName: fmt.Sprintf("%s/%s", owner, repo),
		RepositoryURL:      fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		Timestamp:          time.Now().String(),
	}

	fillHealthIndicators(health, indicators)

	return c.JSON(http.StatusOK, health)
}

func fillHealthIndicators(health *Health, indicators string) {
	if indicators == "all" {
		fmt.Println("all")
	} else {
		indicatorsList := strings.Split(indicators, ",")
		for _, indicator := range indicatorsList {
			switch indicator {
			case "readme":
				readme, err := getReadme(health.RepositoryOwner, health.RepositoryName)
				if err != nil {
					log.Println(err)
				}
				if readme != nil {
					health.Indicators.Readme = *readme
				}
			case "license":
				license, err := getLicense(health.RepositoryOwner, health.RepositoryName)
				if err != nil {
					log.Println(err)
				}
				if license != nil {
					health.Indicators.License = *license
				}
			}
		}
	}
}

func getReadme(owner, repo string) (*Readme, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	opts := &github.RepositoryContentGetOptions{}

	repoReadme, _, err := client.Repositories.GetReadme(ctx, owner, repo, opts)
	if err != nil {
		return nil, err
	}

	readme := &Readme{
		Exists: true,
		URL:    repoReadme.GetHTMLURL(),
	}

	return readme, nil
}

func getLicense(owner, repo string) (*License, error) {
	client := github.NewClient(nil)
	ctx := context.Background()

	repoLicense, _, err := client.Repositories.License(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	license := &License{
		Exists: true,
		URL:    repoLicense.GetHTMLURL(),
		Name:   repoLicense.License.GetName(),
	}

	return license, nil
}

func (h *Handler) getIndicators(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	i := db.DB("healthyrepo").C("indicators")

	result := []Indicator{}

	err := i.Find(nil).All(&result)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}
