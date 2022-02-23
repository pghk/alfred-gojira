package main

import (
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"gojiralfredo/internal/jira-client"
	"gojiralfredo/internal/workflow"
	"log"
)

var (
	client *oreo.Client
	consumer jira.Consumer
	query jira.Query
	wf *aw.Workflow
	hostname string
)

func init() {
	wf = aw.New()
	hostname = "jira.atlassian.com"
	if wf.Config.Get("HOSTNAME") != "" {
		hostname = wf.Config.Get("HOSTNAME")
	}
	query = *jira.BuildQuery()
	client = jira.BuildClient(jira.GetCredentials(wf))
	setConsumer()
}

func setClient(newClient *oreo.Client) {
	client = newClient
}

func setConsumer() {
	consumer = *jira.BuildConsumer("https://" + hostname, client)
}

func runQuery() *jiradata.SearchResults {
	query.QueryFields = "assignee,created,priority,reporter,status,summary,updated,issuetype"
	query.Project = "RAV"
	query.Sort = "priority desc, key desc"


	results, err := consumer.Search(&query)
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
