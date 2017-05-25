package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

func (h *DBHandler) getHealth(c echo.Context) error {
	// db := h.DB.Clone()
	// defer db.Close()

	repoOwner := c.Param("owner")
	repoName := c.Param("repo")
	indicator := c.Param("indicator")

	respData, err := getHealthData(repoOwner, repoName, indicator)
	if err != nil {
		return errors.Wrap(err, "error getting health")
	}

	return c.JSON(http.StatusOK, respData)
}

func (h *DBHandler) getIndicators(c echo.Context) error {
	db := h.DB.Clone()
	defer db.Close()

	i := db.DB("healthyrepo").C("indicators")

	respData := []Indicator{}

	err := i.Find(nil).All(&respData)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, respData)
}
