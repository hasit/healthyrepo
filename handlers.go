package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/labstack/echo"
)

func (h *Handler) getGithubRepoHealth(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	owner := c.Param("owner")
	repo := c.Param("repo")
	queryParams := c.QueryParams()

	health := &Health{
		RepositoryName:     repo,
		RepositoryOwner:    owner,
		RepositoryFullName: fmt.Sprintf("%s/%s", owner, repo),
		RepositoryURL:      fmt.Sprintf("https://github.com/%s/%s", owner, repo),
		Timestamp:          time.Now(),
	}

	for queryKey, queryValues := range queryParams {
		switch queryKey {
		case "indicators":
			fillHealthIndicators(health, queryValues)
		}
	}

	return c.JSON(http.StatusOK, health)
}

func fillHealthIndicators(health *Health, indicators []string) {
	for _, indicator := range indicators {
		switch indicator {
		case "readme":
			// readme, err := getReadme(health.RepositoryOwner, health.RepositoryName)
			// if err != nil {
			// 	log.Println(err)
			// }
			// if readme != nil {
			// 	health.Indicators.Readme = *readme
			// }
			go health.fillReadme()
		case "license":
			// license, err := getLicense(health.RepositoryOwner, health.RepositoryName)
			// if err != nil {
			// 	log.Println(err)
			// }
			// if license != nil {
			// 	health.Indicators.License = *license
			// }
			go health.fillLicense()
		}
	}
}

func (health *Health) fillReadme() {
	client := github.NewClient(nil)
	ctx := context.Background()

	opts := &github.RepositoryContentGetOptions{}

	repoReadme, _, err := client.Repositories.GetReadme(ctx, health.RepositoryOwner, health.RepositoryName, opts)
	if err != nil {
		log.Println(err)
	}

	health.Indicators.Readme.Exists = true
	health.Indicators.Readme.URL = repoReadme.GetHTMLURL()
}

func (health *Health) fillLicense() {
	client := github.NewClient(nil)
	ctx := context.Background()

	repoLicense, _, err := client.Repositories.License(ctx, health.RepositoryOwner, health.RepositoryName)
	if err != nil {
		log.Println(err)
	}

	health.Indicators.License.Exists = true
	health.Indicators.License.URL = repoLicense.GetHTMLURL()
	health.Indicators.License.Name = repoLicense.License.GetName()
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
