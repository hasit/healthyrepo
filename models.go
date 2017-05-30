package main

import "time"

// Indicator holds human readable name and key for making requests to GET /health/:indicator endpoint.
type Indicator struct {
	Name string `json:"name" bson:"name"`
	Key  string `json:"key" bson:"key"`
}

// Repository holds information about the repository for which health queries are made.
type Repository struct {
	ID         int    `json:"id" bson:"id"`
	Host       string `json:"host" bson:"host"`
	Owner      string `json:"owner" bson:"owner"`
	OwnerIsOrg bool   `json:"org" bson:"org"`
	Name       string `json:"name" bson:"name"`
	FullName   string `json:"full_name" bson:"full_name"`
	URL        string `json:"url" bson:"url"`
}

// Docs holds information about health related to static documents like README and LICENSE.
type Docs struct {
	Repository    Repository `json:"repository" bson:"repository"`
	ReadmeExists  bool       `json:"readme_exists" bson:"readme_exists"`
	ReadmeURL     string     `json:"readme_url" bson:"readme_url"`
	LicenseExists bool       `json:"license_exists" bson:"license_exists"`
	LicenseURL    string     `json:"license_url" bson:"license_url"`
	LicenseName   string     `json:"license_name" bson:"license_name"`
}

// ResponseTimes holds information about health related to response times of contributions.
type ResponseTimes struct {
	Repository                Repository               `json:"repository" bson:"repository"`
	ContributorsResponseStats map[string]ResponseStats `json:"contributors_response_stats" bson:"contributors_response_stats"`
}

// ResponseStats holds information to calculate average response times and issue response times of contributor.
type ResponseStats struct {
	FirstContributionWeek time.Time   `json:"first_contribution_week" bson:"first_contribution_week"`
	AverageResponseTime   float32     `json:"avg_rt" bson:"avg_rt"`
	IssuesResponseTimes   map[int]int `json:"issues_rts" bson:"-"`
}

// type Commits struct {
// 	Repository      Repository `json:"repository" bson:"repository"`
// 	TotalCommits    int        `json:"total_commits" bson:"total_commits"`
// 	FirstCommitAt   time.Time  `json:"first_commit_at" bson:"first_commit_at"`
// 	FirstCommitBy   string     `json:"first_commit_by" bson:"first_commit_by"`
// 	LastCommitAt    time.Time  `json:"last_commit_at" bson:"last_commit_at"`
// 	LastCommitBy    string     `json:"last_commit_by" bson:"last_commit_by"`
// 	CommitFrequency []PerWeek  `json:"commit_frequency" bson:"commit_frequency"`
// }

// type PullRequests struct {
// 	Repository         Repository `json:"repository" bson:"repository"`
// 	TotalPullRequests  string     `json:"total_pull_requests" bson:"total_pull_requests"`
// 	MergedPullRequests string     `json:"merged_pull_requests" bson:"merged_pull_requests"`
// 	SentVsMerged       []PerWeek  `json:"sent_vs_merged" bson:"sent_vs_merged"`
// }

// type Issues struct {
// 	Repository       Repository `json:"repository" bson:"repository"`
// 	TotalIssues      string     `json:"total_issues" bson:"total_issues"`
// 	AverageCloseTime string     `json:"average_close_time" bson:"average_close_time"`
// 	OpenVsClosed     []PerWeek  `json:"open_vs_closed" bson:"open_vs_closed"`
// }

// type PerWeek struct {
// 	Week   uint `json:"week" bson:"week"`
// 	Number uint `json:"number" bson:"number"`
// }
