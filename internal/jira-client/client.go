package jira

import (
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiradata"
)

func BuildClient() *jira.Jira {
	return jira.NewJira("https://jira.atlassian.com")
}

func ListIssues() (*jiradata.SearchResults, error) {
	return GetJiraList(*BuildClient())
}

func GetJiraList(j jira.Jira) (*jiradata.SearchResults, error) {
	o := &jira.SearchOptions{Project: "JRACLOUD"}
	return j.Search(o)
}


