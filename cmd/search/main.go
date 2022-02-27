package main

import (
	"flag"
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"go.deanishe.net/fuzzy"
	"gojiralfredo/internal/jira-client"
	"gojiralfredo/internal/workflow"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	client        *oreo.Client
	consumer      jira.Consumer
	jiraQuery     jira.Query
	wf            *aw.Workflow
	hostname      string
	cacheName     = "issues.json"
	maxCacheAge   = 180 * time.Minute
	doDownload  bool
	filter      string
	query		string
	sortOptions []fuzzy.Option
)

func openConfigEditor(wf *aw.Workflow) {
	if err := wf.Alfred.RunTrigger("settings", ""); err != nil {
		wf.FatalError(err)
	}
}

func init() {
	flag.BoolVar(&doDownload, "download", false, "retrieve list of issues from Jira")
	flag.StringVar(&query, "query", "", "JQL string, used to query Jira for issues")
	sortOptions = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.SeparatorBonus(5.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(
		aw.AddMagic(configMagic{}),
		aw.SortOptions(sortOptions...),
	)
	hostname = workflow.GetJiraHostname(wf)
	jiraQuery = *jira.BuildQuery()
	credentials := workflow.GetCredentials(wf, openConfigEditor)
	client = jira.BuildClient(credentials, true, false)
	setConsumer()
}

func setClient(newClient *oreo.Client) {
	client = newClient
}

func setConsumer() {
	consumer = *jira.BuildConsumer("https://"+hostname, client)
}

func setQuery() {
	jiraQuery.QueryFields = strings.Join([]string{
		"assignee",
		"created",
		"priority",
		"status",
		"summary",
		"updated",
		"issuetype",
		"resolution",
	}, ",")

	jiraQuery.MaxResults = 1000

	jiraQuery.SearchOptions.Query = query
}

func runQuery() (*jiradata.SearchResults, error) {
	return consumer.Search(&jiraQuery)
}

func parseResults(results *jiradata.SearchResults) {
	for _, issue := range results.Issues {
		workflow.Add(issue, wf)
	}
}

// Magic Action to enter the "settings" script filter
type configMagic struct{}

func (configMagic) Keyword() string     { return "config" }
func (configMagic) Description() string { return "Edit workflow configuration" }
func (configMagic) RunText() string     { return "Config action registered." }
func (configMagic) Run() error {
	return wf.Alfred.RunTrigger("settings", "")
}

var _ aw.MagicAction = configMagic{}

func run() {
	wf.Args()
	flag.Parse()

	if args := flag.Args(); len(args) > 0 {
		filter = args[0]
	}

	if doDownload {
		wf.Configure(aw.TextErrors(true))
		setQuery()
		results, err := runQuery()
		if err != nil {
			wf.FatalError(err)
		}
		if err := wf.Cache.StoreJSON(cacheName, results); err != nil {
			wf.FatalError(err)
		}
		return
	}

	var results *jiradata.SearchResults
	if wf.Cache.Exists(cacheName) {
		if err := wf.Cache.LoadJSON(cacheName, &results); err != nil {
			wf.FatalError(err)
		}
	}

	if wf.Cache.Expired(cacheName, maxCacheAge) {
		wf.Rerun(0.3)
		if !wf.IsRunning("download") {
			cmd := exec.Command(os.Args[0], "-download", "-query", query)
			if err := wf.RunInBackground("download", cmd); err != nil {
				wf.FatalError(err)
			}
		} else {
			log.Printf("query job already running.")
		}

		if results == nil {
			wf.NewItem("Fetching issues…").
				Icon(aw.IconInfo)
			wf.SendFeedback()
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
