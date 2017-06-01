package main

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// Docs holds information about health related to static documents like README and LICENSE.
type Docs struct {
	Repository    *Repository `json:"repository" bson:"repository"`
	ReadmeExists  bool        `json:"readme_exists" bson:"readme_exists"`
	ReadmeURL     string      `json:"readme_url" bson:"readme_url"`
	LicenseExists bool        `json:"license_exists" bson:"license_exists"`
	LicenseURL    string      `json:"license_url" bson:"license_url"`
	LicenseName   string      `json:"license_name" bson:"license_name"`
}

func getDocs(docs *Docs) error {
	ctx, client := newGithubClient(false)

	repoOwner := docs.Repository.Owner
	repoName := docs.Repository.Name

	opts := &github.RepositoryContentGetOptions{}

	repoReadme, resp, err := client.Repositories.GetReadme(ctx, repoOwner, repoName, opts)
	if err != nil {
		if resp.StatusCode != 404 {
			return errors.Wrap(err, "error getting readme")
		}
		docs.ReadmeExists = false
		docs.ReadmeURL = ""
	}

	if repoReadme != nil {
		docs.ReadmeExists = true
		docs.ReadmeURL = repoReadme.GetHTMLURL()
	}

	repoLicense, resp, err := client.Repositories.License(ctx, repoOwner, repoName)
	if err != nil {
		if resp.StatusCode != 404 {
			return errors.Wrap(err, "error getting license")
		}
		docs.LicenseExists = false
		docs.LicenseURL = ""
		docs.LicenseName = ""
	}

	if repoLicense != nil {
		docs.LicenseExists = true
		docs.LicenseURL = repoLicense.GetHTMLURL()
		docs.LicenseName = repoLicense.License.GetName()
	}

	return nil
}
