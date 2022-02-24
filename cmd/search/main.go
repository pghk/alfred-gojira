package main

import (
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"gojiralfredo/internal/jira-client"
	"gojiralfredo/internal/workflow"
	"log"
	"strings"
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
	hostname = workflow.GetJiraHostname(wf)
	query = *jira.BuildQuery()
	client = jira.BuildClient(workflow.GetCredentials(wf), true, false)
	setConsumer()
}

func setClient(newClient *oreo.Client) {
	client = newClient
}

func setConsumer() {
	consumer = *jira.BuildConsumer("https://" + hostname, client)
}

func runQuery() (*jiradata.SearchResults, error) {
	query.QueryFields = strings.Join([]string{
		"assignee",
		"created",
		"priority",
		"status",
		"summary",
		"updated",
		"issuetype",
	},",")

	if wf.Config.Get("Project") != "" {
		query.Project = wf.Config.Get("Project")
	}
	query.Sort = "priority desc, key desc"

	return consumer.Search(&query)
}

func parseResults() {
	results, err := runQuery()
	if err != nil {
		log.Fatal(err)
	}
	for _, issue := range results.Issues {
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
