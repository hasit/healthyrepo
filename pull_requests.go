package main

import (
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type PullRequests struct {
	Repository   *Repository   `json:"repository" bson:"repository"`
	Total        int           `json:"total" bson:"total"`
	TotalMerged  int           `json:"total_merged" bson:"total_merged"`
	SentVsMerged []*PRPerMonth `json:"sent_vs_merged" bson:"sent_vs_merged"`
}

type PRPerMonth struct {
	MonthYear string `json:"month_year" bson:"month_year"`
	Sent      int    `json:"sent" bson:"sent"`
	Merged    int    `json:"merged" bson:"merged"`
}

type prStats struct {
	sent   int
	merged int
}

func getPullRequests(prs *PullRequests) error {
	ctx, client := newGithubClient(false)

	repoOwner := prs.Repository.Owner
	repoName := prs.Repository.Name

	pullOpts := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "all",
	}
	pulls := []*github.PullRequest{}
	for {
		paginatedPulls, resp, err := client.PullRequests.List(ctx, repoOwner, repoName, pullOpts)
		if err != nil {
			return errors.Wrap(err, "error getting repo pull requests")
		}
		pulls = append(pulls, paginatedPulls...)
		if resp.NextPage == 0 {
			break
		}
		pullOpts.Page = resp.NextPage
	}

	prs.Total = len(pulls)
	prs.TotalMerged = 0

	// prPerMonth := []*PRPerMonth{}

	// var monthYearPRStats map[string]*prStats

	for _, pull := range pulls {
		// pullCreatedAt := pull.GetCreatedAt()
		// merged := false

		// month := pullCreatedAt.Month()
		// year := pullCreatedAt.Year()
		// monthyear := fmt.Sprintf("%d/%d", int(month), int(year))

		if pull.MergedAt != nil {
			prs.TotalMerged++
			// merged = true
		}

		// if monthYearPRStats[monthyear] == nil {
		// 	monthYearPRStats[monthyear] = &prStats{merged: 0, sent: 0}
		// }
		// if merged {
		// 	monthYearPRStats[monthyear].merged++
		// }
		// monthYearPRStats[monthyear].sent++
	}

	// fmt.Println(monthYearPRStats)

	// for my, monthYearPRStat := range monthYearPRStats {
	// 	prPerMonth = append(prPerMonth, &PRPerMonth{MonthYear: my, Merged: monthYearPRStat.merged, Sent: monthYearPRStat.sent})
	// }

	// prs.SentVsMerged = prPerMonth

	return nil
}
