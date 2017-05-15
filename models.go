package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Health struct {
	RepositoryOwner    string    `json:"repository_owner"`
	RepositoryName     string    `json:"repository_name"`
	RepositoryURL      string    `json:"repository_url"`
	RepositoryFullName string    `json:"repository_full_name"`
	Timestamp          time.Time `json:"timestamp"`
	Indicators         struct {
		AverageResponseTimes []struct {
			UserName            string `json:"user_name"`
			AverageResponseTime string `json:"average_response_time"`
		} `json:"average_response_times"`
		Commits struct {
			TimeSinceFirstCommit string `json:"time_since_first_commit"`
			TimeSinceLastCommit  string `json:"time_since_last_commit"`
			CommitFrequency      struct {
				TotalCommits string `json:"total_commits"`
				PerWeek      string `json:"per_week"`
				PerMonth     string `json:"per_month"`
			} `json:"commit_frequency"`
		} `json:"commits"`
		Issues struct {
			TotalIssues        string `json:"total_issues"`
			AverageTimeToClose string `json:"average_time_to_close"`
			OpenVsClosed       struct {
				PerWeek  string `json:"per_week"`
				PerMonth string `json:"per_month"`
			} `json:"open_vs_closed"`
		} `json:"issues"`
		License      License `json:"license"`
		PullRequests struct {
			TotalPullRequests  string `json:"total_pull_requests"`
			MergedPullRequests string `json:"merged_pull_requests"`
			SentVsMerged       struct {
				PerWeek  string `json:"per_week"`
				PerMonth string `json:"per_month"`
			} `json:"sent_vs_merged"`
		} `json:"pull_requests"`
		Readme Readme `json:"readme"`
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

type Indicator struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Key  string        `bson:"key" json:"key"`
}
