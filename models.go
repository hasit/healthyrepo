package main

import "gopkg.in/mgo.v2/bson"

type Health struct {
	RepositoryOwner    string `json:"repository_owner"`
	RepositoryName     string `json:"repository_name"`
	RepositoryURL      string `json:"repository_url"`
	RepositoryFullName string `json:"repository_full_name"`
	Timestamp          string `json:"timestamp"`
	Indicators         struct {
		AverageResponseTimes []AverageResponseTime `json:"average_response_times"`
		Commits              Commits               `json:"commits"`
		Issues               Issues                `json:"issues"`
		License              License               `json:"license"`
		PullRequests         PullRequests          `json:"pull_requests"`
		Readme               Readme                `json:"readme"`
	} `json:"indicators"`
}

type Readme struct {
	Exists bool   `json:"exists"`
	URL    string `json:"url"`
}

type License struct {
	Exists bool   `json:"exists"`
	URL    string `json:"url"`
	Name   string `json:"name"`
}

type PullRequests struct {
	TotalPullRequests  string `json:"total_pull_requests"`
	MergedPullRequests string `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}

type Issues struct {
	TotalPullRequests  string `json:"total_pull_requests"`
	MergedPullRequests string `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}

type Commits struct {
	TotalPullRequests  string `json:"total_pull_requests"`
	MergedPullRequests string `json:"merged_pull_requests"`
	SentVsMerged       struct {
		PerWeek  string `json:"per_week"`
		PerMonth string `json:"per_month"`
	} `json:"sent_vs_merged"`
}

type AverageResponseTime struct {
	UserName            string `json:"user_name"`
	AverageResponseTime string `json:"average_response_time"`
}

type Indicator struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Key  string        `bson:"key" json:"key"`
}
