package main

import (
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"gojiralfredo/internal/jira-client"
	"gojiralfredo/internal/workflow"
	"log"
)

var (
	client jira.Client
	query jira.Query
	wf *aw.Workflow
)

func init() {
	client = *jira.BuildClient("https://jira.atlassian.com")
	query = *jira.BuildQuery()
	wf = aw.New()
}

func runQuery() *jiradata.SearchResults {
	query.QueryFields = "issuetype,status,assignee"

	results, err := client.Search(&query)
	if err != nil {
		log.Fatal(err)
	}
	return results
}

func parseResults() {
	results := runQuery().Issues
	for _, issue := range results {
		workflow.Add(issue, wf)
	}

}

func run() {
	parseResults()
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
