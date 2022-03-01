package main

import (
	"flag"
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"github.com/pghk/alfred-gojira/internal/jira"
	"github.com/pghk/alfred-gojira/internal/workflow"
	"go.deanishe.net/fuzzy"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	wf          *aw.Workflow
	cacheName   = "issues.json"
	maxCacheAge = 180 * time.Minute
	getData     bool
	filter      string
	query       string
	sortOptions []fuzzy.Option
	results     *jiradata.SearchResults
)

func openConfigEditor() {
	if err := wf.Alfred.RunTrigger(workflow.ConfigTriggerName, ""); err != nil {
		wf.FatalError(err)
	}
}

func init() {
	flag.BoolVar(&getData, "fresh", false, "Populate cache with list of issues from Jira")
	flag.StringVar(&query, "query", "", "JQL string, used to query Jira for issues")
	sortOptions = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.SeparatorBonus(5.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = workflow.BuildWorkflow(sortOptions)
}

func runQuery(client *oreo.Client, jiraQuery *jira.Query) (*jiradata.SearchResults, error) {
	hostname := workflow.GetJiraHostname()
	if client == nil {
		if workflow.CredentialsRequired() {
			credentials := workflow.GetCredentials(openConfigEditor)
			client = jira.BuildClient(credentials, true, true)
		} else {
			client = jira.BuildClient(jira.Auth{}, false, true)
		}
	}

	if jiraQuery == nil {
		fields := strings.Join([]string{
			"assignee",
			"created",
			"priority",
			"status",
			"summary",
			"updated",
			"issuetype",
			"resolution",
		}, ",")
		jiraQuery = jira.BuildQuery(query, fields, 1000)
	}

	consumer := *jira.BuildConsumer("https://"+hostname, client)
	return consumer.Search(jiraQuery)
}

func refreshCache(results *jiradata.SearchResults) bool {
	wf.Rerun(0.3)
	if !wf.IsRunning("fetchIssues") {
		fetchIssues()
	} else {
		log.Printf("query job already running.")
	}

	if results == nil {
		wf.NewItem("Fetching issues…").
			Icon(aw.IconInfo)
		wf.SendFeedback()
		return true
	}
	return false
}

func fetchIssues() {
	cmd := exec.Command(os.Args[0], "-fresh", "-query", query)
	if err := wf.RunInBackground("fetchIssues", cmd); err != nil {
		wf.FatalError(err)
	}
}

func parseResults(results *jiradata.SearchResults) {
	for _, issue := range results.Issues {
		workflow.Add(issue, wf)
	}
}

func run() {
	wf.Args()
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		filter = args[0]
	}

	if getData {
		wf.Configure(aw.TextErrors(true))
		results, err := runQuery(nil, nil)
		if err != nil {
			wf.FatalError(err)
		}
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			wf.FatalError(err)
		}
		return
	}

	if wf.Cache.Exists(cacheName) {
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			wf.FatalError(err)
		}
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		if refreshCache(results) {
			return
		}
	}

	parseResults(results)

	if filter != "" {
		wf.Filter(filter)
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
