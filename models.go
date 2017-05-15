package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Health struct {
	RepositoryURL  string    `json:"repository_url"`
	RepositoryName string    `json:"repository_name"`
	Timestamp      time.Time `json:"timestamp"`
	Indicators     struct {
		AverageResponseTimes []struct {
			UserName            string `json:"user_name,omitempty"`
			AverageResponseTime string `json:"average_response_time,omitempty"`
		} `json:"average_response_times,omitempty"`
		Commits struct {
			TimeSinceFirstCommit string `json:"time_since_first_commit,omitempty"`
			TimeSinceLastCommit  string `json:"time_since_last_commit,omitempty"`
			CommitFrequency      struct {
				TotalCommits string `json:"total_commits,omitempty"`
				PerWeek      string `json:"per_week,omitempty"`
				PerMonth     string `json:"per_month,omitempty"`
			} `json:"commit_frequency,omitempty"`
		} `json:"commits,omitempty"`
		Issues struct {
			TotalIssues        string `json:"total_issues,omitempty"`
			AverageTimeToClose string `json:"average_time_to_close,omitempty"`
			OpenVsClosed       struct {
				PerWeek  string `json:"per_week,omitempty"`
				PerMonth string `json:"per_month,omitempty"`
			} `json:"open_vs_closed,omitempty"`
		} `json:"issues,omitempty"`
		License struct {
			Exists bool   `json:"exists,omitempty"`
			URL    string `json:"url,omitempty"`
			Name   string `json:"name,omitempty"`
		} `json:"license,omitempty"`
		PullRequests struct {
			TotalPullRequests  string `json:"total_pull_requests,omitempty"`
			MergedPullRequests string `json:"merged_pull_requests,omitempty"`
			SentVsMerged       struct {
				PerWeek  string `json:"per_week,omitempty"`
				PerMonth string `json:"per_month,omitempty"`
			} `json:"sent_vs_merged,omitempty"`
		} `json:"pull_requests,omitempty"`
		Readme struct {
			Exists bool   `json:"exists,omitempty"`
			URL    string `json:"url,omitempty"`
		} `json:"readme,omitempty"`
	} `json:"indicators,omitempty"`
}

type Indicator struct {
	ID   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Key  string        `bson:"key" json:"key"`
}
