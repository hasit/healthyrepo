package main

import (
	"net/http"

	"github.com/pkg/errors"
)

// Repository holds information a repository hosted on Github.
type Repository struct {
	ID             int    `json:"github_id" bson:"github_id"`
	Owner          string `json:"owner" bson:"owner"`
	OwnerAvatarURL string `json:"owner_avatar_url" bson:"owner_avatar_url"`
	OwnerIsOrg     bool   `json:"org" bson:"org"`
	Name           string `json:"name" bson:"name"`
	FullName       string `json:"full_name" bson:"full_name"`
	RepoURL        string `json:"repo_url" bson:"repo_url"`
	OwnerURL       string `json:"owner_url" bson:"owner_url"`
}

func getRepo(repoOwner, repoName string) (*Repository, error) {
	ctx, client := newGithubClient(true)

	repo, resp, err := client.Repositories.Get(ctx, repoOwner, repoName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting repository %s/%s", repoOwner, repoName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	repository := &Repository{
		ID:             repo.GetID(),
		Owner:          repo.Owner.GetLogin(),
		Name:           repo.GetName(),
		FullName:       repo.GetFullName(),
		RepoURL:        repo.GetHTMLURL(),
		OwnerURL:       repo.Owner.GetHTMLURL(),
		OwnerAvatarURL: repo.Owner.GetAvatarURL(),
		OwnerIsOrg:     false,
	}

	if repo.Organization != nil {
		repository.OwnerIsOrg = true
	}

	return repository, nil
}
