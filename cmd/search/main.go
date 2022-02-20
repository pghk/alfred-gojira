package main

import (
	"fmt"
	"gojiralfredo/internal/jira-client"
	"log"
)

var (
	client jira.Client
	query jira.Query
)

func init() {
	client = *jira.BuildClient("https://jira.atlassian.com")
	query = *jira.BuildQuery()
}

func runQuery() jira.Results {
	results, err := client.Search(&query)
	if err != nil {
		log.Fatal(err)
	}
	return jira.Results{SearchResults: *results}
}

func main() {
	fmt.Println(runQuery())
}
