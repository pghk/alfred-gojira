package main

import (
	"flag"
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
	wf = aw.New(aw.AddMagic(configMagic{}))
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

	query.MaxResults = 200

	if wf.Config.Get("Project") != "" {
		query.Project = wf.Config.Get("Project")
	}
	query.Sort = "updated desc, sprint desc, priority desc, key desc"

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

// Magic Action to enter the "settings" script filter
type configMagic struct {}
func (configMagic) Keyword() string {return "config" }
func (configMagic) Description() string {return "Edit workflow configuration" }
func (configMagic) RunText() string {return "Config action registered." }
func (configMagic) Run() error {
	return wf.Alfred.RunTrigger("settings", "")
}
var _ aw.MagicAction = configMagic{}

func run() {
	wf.Args()
	flag.Parse()
	query := flag.Arg(0)
	parseResults()
	if query != "" {
		wf.Filter(query)
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
