package jira

import (
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiradata"
)

type Client struct {
	*jira.Jira
}

type Query struct {
	*jira.SearchOptions
}

type Results struct {
	jiradata.SearchResults
}

func BuildQuery() *Query {
	return &Query{&jira.SearchOptions{}}
}

func BuildClient(endpoint string) *Client {
	return &Client{jira.NewJira(endpoint)}
}
