package main

import (
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiradata"
)

func GetJiraList(j jira.Jira) (*jiradata.SearchResults, error) {
	o := &jira.SearchOptions{Project: "JRACLOUD"}
	return j.Search(o)
}
