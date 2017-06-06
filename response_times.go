package main

import (
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

// ResponseTimes holds information about health related to response times of contributions.
type ResponseTimes struct {
	Repository                *Repository                 `json:"repository" bson:"repository"`
	Since                     time.Time                   `json:"since" bson:"since"`
	ContributorsResponseStats []*ContributorResponseStats `json:"contributors_response_stats" bson:"contributors_response_stats"`
}

// ContributorResponseStats holds information to calculate average response times and issue response times of contributor.
type ContributorResponseStats struct {
	Contributor           string               `json:"contributor" bson:"contributor"`
	FirstContributionWeek time.Time            `json:"first_contribution_week" bson:"first_contribution_week"`
	AverageResponseTime   float64              `json:"average_response_time" bson:"average_response_time"`
	IssuesResponseTimes   []*IssueResponseTime `json:"issues_response_times" bson:"issues_response_times"`
}

// IssueResponseTime holds information about response times corresponding to issues.
type IssueResponseTime struct {
	IssueCreatedAt time.Time `json:"issue_created_at" bson:"issue_created_at"`
	Internal       bool      `json:"internal" bson:"internal"`
	IssueNumber    int       `json:"issue_number" bson:"issue_number"`
	ResponseTime   int       `json:"response_time" bson:"response_time"`
}

func getReponseTimes(responseTimes *ResponseTimes) error {
	ctx, client := newGithubClient(true)

	repoOwner := responseTimes.Repository.Owner
	repoName := responseTimes.Repository.Name

	contributorsStats, _, err := client.Repositories.ListContributorsStats(ctx, repoOwner, repoName)
	if _, ok := err.(*github.AcceptedError); ok {
		return errors.Wrap(err, "scheduled on Github side")
	}

	contributorsResponseStats := []*ContributorResponseStats{}

	for _, contributorStats := range contributorsStats {
		contributorResponseStats := &ContributorResponseStats{}
		contributorResponseStats.Contributor = contributorStats.Author.GetLogin()
		contributorResponseStats.AverageResponseTime = -1
		contributorResponseStats.IssuesResponseTimes = []*IssueResponseTime{}
		for weeki, week := range contributorStats.Weeks {
			if week.GetCommits() != 0 {
				contributorResponseStats.FirstContributionWeek = contributorStats.Weeks[weeki].GetWeek().Time
				break
			}
		}
		contributorsResponseStats = append(contributorsResponseStats, contributorResponseStats)
	}

	since := time.Now().AddDate(0, -6, 0)

	issueListByRepoOpts := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		Since:       since,
		State:       "all",
	}
	var issues []*github.Issue
	for {
		paginatedIssues, resp, err := client.Issues.ListByRepo(ctx, repoOwner, repoName, issueListByRepoOpts)
		if err != nil {
			return errors.Wrap(err, "error getting repo issues")
		}
		issues = append(issues, paginatedIssues...)
		if resp.NextPage == 0 {
			break
		}
		issueListByRepoOpts.Page = resp.NextPage
	}

	for _, issue := range issues {
		issueCreatedAt := issue.GetCreatedAt()
		issueCreator := issue.User.GetLogin()
		issueNumber := issue.GetNumber()
		issueCreatorIsContributor := false

		for _, contributorResponseStats := range contributorsResponseStats {
			if issueCreator == contributorResponseStats.Contributor {
				issueCreatorIsContributor = true
			}
		}

		issueListCommentsOpts := &github.IssueListCommentsOptions{
			ListOptions: github.ListOptions{PerPage: 50},
			Sort:        "created",
		}
		var issueComments []*github.IssueComment
		for {
			paginatedIssueComments, resp, err := client.Issues.ListComments(ctx, repoOwner, repoName, issue.GetNumber(), issueListCommentsOpts)
			if err != nil {
				return errors.Wrap(err, "error getting issue comments")
			}
			issueComments = append(issueComments, paginatedIssueComments...)
			if resp.NextPage == 0 {
				break
			}
			issueListCommentsOpts.Page = resp.NextPage
		}

		for _, issueComment := range issueComments {
			issueCommentCreator := issueComment.User.GetLogin()
			issueCommentCreatedAt := issueComment.GetCreatedAt()
			responseTimeInDays := daysDiff(issueCommentCreatedAt, issueCreatedAt)

			if issueCommentCreator != issueCreator {
				for _, contributorResponseStats := range contributorsResponseStats {
					if issueCommentCreator == contributorResponseStats.Contributor {
						issueResponseTime := &IssueResponseTime{}
						issueResponseTime.IssueCreatedAt = issueCreatedAt
						issueResponseTime.IssueNumber = issueNumber

						gotIssueResponseTime := false
						for _, issuesResponseTimes := range contributorResponseStats.IssuesResponseTimes {
							if issuesResponseTimes.IssueNumber == issueNumber {
								gotIssueResponseTime = true
							}
						}

						if contributorResponseStats.FirstContributionWeek.Before(issueCreatedAt) && !gotIssueResponseTime {
							if issueCreatorIsContributor {
								issueResponseTime.Internal = true
							}
							issueResponseTime.ResponseTime = responseTimeInDays
							contributorResponseStats.IssuesResponseTimes = append(contributorResponseStats.IssuesResponseTimes, issueResponseTime)
							break
						}
					}
				}
			}
		}
	}

	for _, contributorResponseStats := range contributorsResponseStats {
		if len(contributorResponseStats.IssuesResponseTimes) > 0 {
			var totalResponseTime int
			for _, issueResponseTime := range contributorResponseStats.IssuesResponseTimes {
				if !issueResponseTime.Internal {
					totalResponseTime += issueResponseTime.ResponseTime
				}
			}
			contributorResponseStats.AverageResponseTime = float64(totalResponseTime) / float64(len(contributorResponseStats.IssuesResponseTimes))
		}
	}

	responseTimes.ContributorsResponseStats = contributorsResponseStats

	return nil
}
