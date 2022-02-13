package main

import (
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiradata"
)

func GetJiraList() (*jiradata.SearchResults, error) {
	j := jira.NewJira("https://jira.atlassian.com")
	o := &jira.SearchOptions{Project: "JRACLOUD"}
	return j.Search(o)
}
