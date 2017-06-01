package main

import (
	"fmt"
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
	IssueNumber    int       `json:"issue_number" bson:"issue_number"`
	ResponseTime   int       `json:"response_time" bson:"response_time"`
}

func getReponseTimes(responseTimes *ResponseTimes) error {
	ctx, client := newGithubClient(true)

	repoOwner := responseTimes.Repository.Owner
	repoName := responseTimes.Repository.Name

	contributorsStats, _, err := client.Repositories.ListContributorsStats(ctx, repoOwner, repoName)
	if err != nil {
		return errors.Wrap(err, "error getting contributor stats")
	}

	contributorsResponseStats := []*ContributorResponseStats{}

	for _, contributorStats := range contributorsStats {
		contributorResponseStats := &ContributorResponseStats{}
		contributorResponseStats.Contributor = contributorStats.Author.GetLogin()
		contributorResponseStats.AverageResponseTime = -1
		for weeki, week := range contributorStats.Weeks {
			if week.GetCommits() != 0 {
				contributorResponseStats.FirstContributionWeek = contributorStats.Weeks[weeki].GetWeek().Time
				break
			}
		}
		contributorsResponseStats = append(contributorsResponseStats, contributorResponseStats)
	}

	since := time.Now().AddDate(0, -6, 0)
	responseTimes.Since = since

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

	// wg := sync.WaitGroup{}
	// wg.Add(len(contributorsResponseStats))
	for _, contributorResponseStats := range contributorsResponseStats {
		getContributorResponseStats(repoOwner, repoName, issues, contributorResponseStats)
	}
	// wg.Wait()

	responseTimes.ContributorsResponseStats = contributorsResponseStats

	return nil
}

func getContributorResponseStats(repoOwner string, repoName string, issues []*github.Issue, contributorResponseStats *ContributorResponseStats) {
	// defer wg.Done()
	ctx, client := newGithubClient(true)
	contributor := contributorResponseStats.Contributor
	issuesResponseTimes := []*IssueResponseTime{}

	for _, issue := range issues {
		issueCreatedAt := issue.GetCreatedAt()
		issueCreator := issue.User.GetLogin()
		issueNumber := issue.GetNumber()

		issueResponseTime := &IssueResponseTime{}
		issueResponseTime.IssueCreatedAt = issueCreatedAt
		issueResponseTime.IssueNumber = issueNumber

		if issueCreator != contributor {
			issueListCommentsOpts := &github.IssueListCommentsOptions{
				ListOptions: github.ListOptions{PerPage: 50},
				Sort:        "created",
			}
			var issueComments []*github.IssueComment
			for {
				paginatedIssueComments, resp, err := client.Issues.ListComments(ctx, repoOwner, repoName, issue.GetNumber(), issueListCommentsOpts)
				if err != nil {
					fmt.Println(err)
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

				if issueCommentCreator == contributor {
					if contributorResponseStats.FirstContributionWeek.Before(issueCreatedAt) {
						issueResponseTime.ResponseTime = responseTimeInDays
						issuesResponseTimes = append(issuesResponseTimes, issueResponseTime)
						break
					}
				}
			}
		}
	}

	if len(issuesResponseTimes) > 0 {
		var totalResponseTime int
		for _, issueResponseTime := range issuesResponseTimes {
			totalResponseTime += issueResponseTime.ResponseTime
		}
		contributorResponseStats.AverageResponseTime = float64(totalResponseTime) / float64(len(issuesResponseTimes))
	}

	contributorResponseStats.IssuesResponseTimes = issuesResponseTimes
}
