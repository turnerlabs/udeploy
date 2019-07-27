package handler

import (
	"fmt"
	"net/http"

	"github.com/turnerlabs/udeploy/component/commit"

	"github.com/labstack/echo/v4"
	"github.com/turnerlabs/udeploy/component/cache"
)

const apiURL = "https://api.github.com"

// GetInstanceCommits ..
func GetInstanceCommits(c echo.Context) error {

	app, found := cache.Apps.Get(c.Param("app"))
	if !found {
		return fmt.Errorf("%s app not found", c.Param("app"))
	}

	if len(app.Repo.Org) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	inst, found := app.Instances[c.Param("instance")]
	if !found {
		return fmt.Errorf("%s instance not found", c.Param("instance"))
	}

	commits, err := commit.BuildRelease(app.Repo.Org, app.Repo.Name, inst.Version(), "", apiURL, app.Repo.AccessToken, 50, app.Repo.CommitConfig)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, commits)
}

// GetVersionCommitsByRange ..
func GetVersionCommitsByRange(c echo.Context) error {

	app, found := cache.Apps.Get(c.Param("app"))
	if !found {
		return fmt.Errorf("%s app not found", c.Param("app"))
	}

	if len(app.Repo.Org) == 0 {
		return c.JSON(http.StatusOK, []string{})
	}

	commits, err := commit.BuildRelease(app.Repo.Org, app.Repo.Name, c.Param("target"), c.Param("current"), apiURL, app.Repo.AccessToken, 50, app.Repo.CommitConfig)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, commits)
}
