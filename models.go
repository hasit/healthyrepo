package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

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
	Repository           Repository            `json:"repository"`
	AverageResponseTimes []AverageResponseTime `json:"average_response_times"`
}

type AverageResponseTime struct {
	Contributor           string    `json:"contributor"`
	AverageResponseTime   float32   `json:"average_response_time"`
	FirstContributionWeek time.Time `json:"first_contribution_week"`
}

type Indicator struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Key  string        `bson:"key" json:"key"`
}

type PullRequests struct {
	Repository         Repository `json:"repository"`
	TotalPullRequests  string     `json:"total_pull_requests"`
	MergedPullRequests string     `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}

type Issues struct {
	Repository         Repository `json:"repository"`
	TotalPullRequests  string     `json:"total_pull_requests"`
	MergedPullRequests string     `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}

type Commits struct {
	Repository         Repository `json:"repository"`
	TotalPullRequests  string     `json:"total_pull_requests"`
	MergedPullRequests string     `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}
