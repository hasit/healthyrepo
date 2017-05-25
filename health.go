package main

import (
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func getHealthData(repoOwner, repoName, indicator string) (interface{}, error) {
	ctx, client := newGithubClient(true)

	githubRepo, resp, err := client.Repositories.Get(ctx, repoOwner, repoName)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting repository %s/%s", repoOwner, repoName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	repository := Repository{
		ID:         githubRepo.GetID(),
		Owner:      githubRepo.Owner.GetLogin(),
		Name:       githubRepo.GetName(),
		FullName:   githubRepo.GetFullName(),
		URL:        githubRepo.GetHTMLURL(),
		OwnerIsOrg: false,
	}

	if githubRepo.Organization != nil {
		repository.OwnerIsOrg = true
	}

	var health interface{}

	switch indicator {
	case "docs":
		docs, err := getDocs(repoOwner, repoName)
		if err != nil {
			return nil, errors.Wrap(err, "error getting docs")
		}
		if docs != nil {
			docs.Repository = repository
			health = docs
		}
	case "response_times":
		responseTimes, err := getReponseTimes(repoOwner, repoName)
		if err != nil {
			return nil, errors.Wrap(err, "error getting contributors")
		}
		if responseTimes != nil {
			responseTimes.Repository = repository
			health = responseTimes
		}
	}

	return health, nil
}

func getDocs(repoOwner, repoName string) (*Docs, error) {
	docs := &Docs{}

	ctx, client := newGithubClient(false)

	opts := &github.RepositoryContentGetOptions{}

	repoReadme, _, err := client.Repositories.GetReadme(ctx, repoOwner, repoName, opts)
	if err != nil {
		return nil, errors.Wrap(err, "error getting readme")
	}

	if repoReadme != nil {
		docs.ReadmeExists = true
		docs.ReadmeURL = repoReadme.GetHTMLURL()
	}

	repoLicense, _, err := client.Repositories.License(ctx, repoOwner, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "error getting license")
	}

	if repoLicense != nil {
		docs.LicenseExists = true
		docs.LicenseURL = repoLicense.GetHTMLURL()
		docs.LicenseName = repoLicense.License.GetName()
	}

	return docs, nil
}

func getReponseTimes(repoOwner, repoName string) (*ResponseTimes, error) {
	responseTimes := &ResponseTimes{}

	ctx, client := newGithubClient(true)

	contributorsStats, _, err := client.Repositories.ListContributorsStats(ctx, repoOwner, repoName)
	if err != nil {
		return nil, err
	}

	for _, contributorStats := range contributorsStats {
		averageResponseTime := AverageResponseTime{}
		averageResponseTime.Contributor = contributorStats.Author.GetLogin()
		for wi, week := range contributorStats.Weeks {
			if week.GetCommits() != 0 {
				averageResponseTime.FirstContributionWeek = contributorStats.Weeks[wi].GetWeek().Time
				break
			}
		}
		averageResponseTime.AverageResponseTime = -1
		responseTimes.AverageResponseTimes = append(responseTimes.AverageResponseTimes, averageResponseTime)
	}

	firstContributionWeeks := make(map[string]time.Time)
	for _, averageResponseTime := range responseTimes.AverageResponseTimes {
		firstContributionWeeks[averageResponseTime.Contributor] = averageResponseTime.FirstContributionWeek
	}

	issueListByRepoOpts := &github.IssueListByRepoOptions{
		Since: time.Now().AddDate(-1, 0, 0),
		State: "all",
	}
	issues, _, err := client.Issues.ListByRepo(ctx, repoOwner, repoName, issueListByRepoOpts)
	if err != nil {
		return nil, err
	}

	responseTimesList := make(map[string]map[int]int)

	for _, issue := range issues {
		issueCreatedAt := issue.GetCreatedAt()
		issueCreator := issue.User.GetLogin()
		issueNumber := issue.GetNumber()

		issueListCommentsOpts := &github.IssueListCommentsOptions{
			Sort: "created",
		}
		issueComments, _, err := client.Issues.ListComments(ctx, repoOwner, repoName, issue.GetNumber(), issueListCommentsOpts)
		if err != nil {
			return nil, err
		}

		for _, issueComment := range issueComments {
			issueCommentCreator := issueComment.User.GetLogin()

			if issueCommentCreator != issueCreator {
				if firstContributionWeeks[issueCommentCreator].Before(issueCreatedAt) {
					issueCommentCreatedAt := issueComment.GetCreatedAt()
					responseTimeInDays := daysDiff(issueCommentCreatedAt, issueCreatedAt)

					if responseTimesList[issueCommentCreator] != nil {
						issueCommentCreatorResponseTimes := responseTimesList[issueCommentCreator]

						if _, ok := issueCommentCreatorResponseTimes[issueNumber]; !ok {
							issueCommentCreatorResponseTimes[issueNumber] = responseTimeInDays
							responseTimesList[issueCommentCreator] = issueCommentCreatorResponseTimes
						}
					} else {
						issueCommentCreatorResponseTimes := make(map[int]int)
						issueCommentCreatorResponseTimes[issueNumber] = responseTimeInDays
						responseTimesList[issueCommentCreator] = issueCommentCreatorResponseTimes
					}
				}
			}
		}
	}

	averageResponseTimes := responseTimes.AverageResponseTimes

	for user, rts := range responseTimesList {
		totalTime := 0
		for _, rt := range rts {
			totalTime += rt
		}
		for i, arts := range averageResponseTimes {
			if arts.Contributor == user {
				averageResponseTimes[i].AverageResponseTime = float32(totalTime) / float32(len(rts))
			}
		}
	}

	responseTimes.AverageResponseTimes = averageResponseTimes

	return responseTimes, nil
}
