package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Indicator struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Key  string        `bson:"key" json:"key"`
}

type Repository struct {
	ID         int    `json:"id"`
	Owner      string `json:"owner"`
	OwnerIsOrg bool   `json:"owner_is_org"`
	Name       string `json:"name"`
	FullName   string `json:"full_name"`
	URL        string `json:"url"`
}

type Docs struct {
	Repository    Repository `json:"repository"`
	ReadmeExists  bool       `json:"readme_exists"`
	ReadmeURL     string     `json:"readme_url"`
	LicenseExists bool       `json:"license_exists"`
	LicenseURL    string     `json:"license_url"`
	LicenseName   string     `json:"license_name"`
}

type ResponseTimes struct {
	Repository               Repository                `json:"repository"`
	ContributorResponseTimes []ContributorResponseTime `json:"contributor_response_times"`
}

type ContributorResponseTime struct {
	Contributor           string                `json:"contributor"`
	AverageResponseTime   float32               `json:"average_response_time"`
	ResponseTimesByIssues []ResponseTimeByIssue `json:"response_times_by_issues"`
	FirstContributionWeek time.Time             `json:"first_contribution_week"`
}

type ResponseTimeByIssue struct {
	IssueNumber       int `json:"issue_number"`
	IssueResponseTime int `json:"issue_response_time"`
}

type Commits struct {
	Repository      Repository `json:"repository"`
	TotalCommits    int        `json:"total_commits"`
	FirstCommitAt   time.Time  `json:"first_commit_at"`
	FirstCommitBy   string     `json:"first_commit_by"`
	LastCommitAt    time.Time  `json:"last_commit_at"`
	LastCommitBy    string     `json:"last_commit_by"`
	CommitFrequency []PerWeek  `json:"commit_frequency"`
}

type PullRequests struct {
	Repository         Repository `json:"repository"`
	TotalPullRequests  string     `json:"total_pull_requests"`
	MergedPullRequests string     `json:"merged_pull_requests"`
	SentVsMerged       []PerWeek  `json:"sent_vs_merged"`
}

type Issues struct {
	Repository         Repository `json:"repository"`
	TotalIssues        string     `json:"total_issues"`
	AverageTimeToClose string     `json:"average_time_to_close"`
	OpenVsClosed       []PerWeek  `json:"open_vs_closed"`
}

type PerWeek struct {
	Week   uint `json:"week"`
	Number uint `json:"number"`
}
