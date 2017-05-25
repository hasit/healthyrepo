package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newGithubClient(auth bool) (context.Context, *github.Client) {
	ctx := context.Background()

	if auth {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: getEnv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		return ctx, github.NewClient(tc)
	}

	return ctx, github.NewClient(nil)
}
