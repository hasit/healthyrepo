package main

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

func (h *DBHandler) getIndicators(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	indicators := &Indicators{}

	i := db.DB("healthyrepo").C("indicators")
	err := i.Find(nil).All(indicators)
	if err != nil {
		httpError := &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: errors.Wrap(err, "Error getting indicators list from database").Error(),
		}
		return c.JSON(http.StatusInternalServerError, httpError)
	}

	return c.JSON(http.StatusOK, indicators)
}

func (h *DBHandler) getIndicator(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	indicatorKey := c.Param("indicator")

	indicator := &Indicator{}

	i := db.DB("healthyrepo").C("indicators")
	err := i.Find(bson.M{"key": indicatorKey}).One(indicator)
	if err != nil {
		httpError := &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: errors.Wrapf(err, "Error getting indicator '%s' from database", indicatorKey).Error(),
		}
		return c.JSON(http.StatusInternalServerError, httpError)
	}

	return c.JSON(http.StatusOK, indicator)
}

func (h *DBHandler) getRepository(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	repoOwner := c.Param("owner")
	repoName := c.Param("repo")
	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)

	repository, err := getRepo(repoOwner, repoName)
	if err != nil {
		httpError := &HTTPError{
			Code:    http.StatusInternalServerError,
			Message: errors.Wrapf(err, "Error getting information for repository '%s' from Github API", repoFullName).Error(),
		}
		return c.JSON(http.StatusInternalServerError, httpError)
	}

	return c.JSON(http.StatusOK, repository)
}

func (h *DBHandler) getRepositoryHealthDocs(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	repoOwner := c.Param("owner")
	repoName := c.Param("repo")

	repository, err := getRepo(repoOwner, repoName)
	if err != nil {
		return err
	}

	docs := &Docs{}

	r := db.DB("healthyrepo").C("docs")
	err = r.Find(bson.M{"repository.full_name": repository.FullName}).One(docs)
	if err != nil {
		if err.Error() == "not found" {
			docs.Repository = repository

			err := getDocs(docs)
			if err != nil {
				httpError := &HTTPError{
					Code:    http.StatusInternalServerError,
					Message: errors.Wrapf(err, "Error getting health 'docs' of repo '%s' in Github API", repository.FullName).Error(),
				}
				return c.JSON(http.StatusInternalServerError, httpError)
			}

			err = r.Insert(docs)
			if err != nil {
				httpError := &HTTPError{
					Code:    http.StatusInternalServerError,
					Message: errors.Wrapf(err, "Error inserting health 'docs' of repo '%s' in database", repository.FullName).Error(),
				}
				return c.JSON(http.StatusInternalServerError, httpError)
			}
			log.Infof("Success inserting health 'docs' of repo '%s' in database", repository.FullName)
		} else {
			httpError := &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: errors.Wrapf(err, "Error getting health 'docs' of repo '%s' in database", repository.FullName).Error(),
			}
			return c.JSON(http.StatusInternalServerError, httpError)
		}
	}

	return c.JSON(http.StatusOK, docs)
}

func (h *DBHandler) getRepositoryHealthResponseTimes(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	repoOwner := c.Param("owner")
	repoName := c.Param("repo")

	repository, err := getRepo(repoOwner, repoName)
	if err != nil {
		return err
	}

	responseTimes := &ResponseTimes{}

	r := db.DB("healthyrepo").C("response_times")
	err = r.Find(bson.M{"repository.full_name": repository.FullName}).One(responseTimes)
	if err != nil {
		if err.Error() == "not found" {
			responseTimes.Repository = repository

			err := getReponseTimes(responseTimes)
			if err != nil {
				httpError := &HTTPError{
					Code:    http.StatusInternalServerError,
					Message: errors.Wrapf(err, "Error getting health 'response_times' of repo '%s' in Github API", repository.FullName).Error(),
				}
				return c.JSON(http.StatusInternalServerError, httpError)
			}

			err = r.Insert(responseTimes)
			if err != nil {
				httpError := &HTTPError{
					Code:    http.StatusInternalServerError,
					Message: errors.Wrapf(err, "Error inserting health 'response_times' of repo '%s' in database", repository.FullName).Error(),
				}
				return c.JSON(http.StatusInternalServerError, httpError)
			}

			log.Infof("Success inserting health 'response_times' of repo '%s' in database", repository.FullName)

		} else {
			httpError := &HTTPError{
				Code:    http.StatusInternalServerError,
				Message: errors.Wrapf(err, "Error getting health 'response_times' of repo '%s' in database", repository.FullName).Error(),
			}
			return c.JSON(http.StatusInternalServerError, httpError)
		}
	}

	return c.JSON(http.StatusOK, responseTimes)
}
