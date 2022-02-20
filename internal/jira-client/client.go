package jira

import (
	"github.com/go-jira/jira"
)

type Client struct {
	*jira.Jira
}

type Query struct {
	*jira.SearchOptions
}

func BuildQuery() *Query {
	return &Query{&jira.SearchOptions{}}
}

func BuildClient(endpoint string) *Client {
	return &Client{jira.NewJira(endpoint)}
}
