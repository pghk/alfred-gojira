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
	doDownload    bool
	workflowQuery string
	sopts         []fuzzy.Option
)

func init() {
	flag.BoolVar(&doDownload, "download", false, "retrieve list of issues from Jira")
	sopts = []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(
		aw.AddMagic(configMagic{}),
		aw.SortOptions(sopts...),
	)
	hostname = workflow.GetJiraHostname(wf)
	jiraQuery = *jira.BuildQuery()
	client = jira.BuildClient(workflow.GetCredentials(wf), true, true)
	setConsumer()
}

func setClient(newClient *oreo.Client) {
	client = newClient
}

func setConsumer() {
	consumer = *jira.BuildConsumer("https://" + hostname, client)
}

func runQuery() (*jiradata.SearchResults, error) {
	jiraQuery.QueryFields = strings.Join([]string{
		"assignee",
		"created",
		"priority",
		"status",
		"summary",
		"updated",
		"issuetype",
	},",")

	jiraQuery.MaxResults = 100

	if wf.Config.Get("Project") != "" {
		jiraQuery.Project = wf.Config.Get("Project")
	}
	jiraQuery.Sort = "updated desc, sprint desc, priority desc, key desc"

	return consumer.Search(&jiraQuery)
}

func parseResults(results *jiradata.SearchResults) {
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

	if args := flag.Args(); len(args) > 0 {
		workflowQuery = args[0]
	}

	if doDownload {
		wf.Configure(aw.TextErrors(true))
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
			cmd := exec.Command(os.Args[0], "-download")
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

	if workflowQuery != "" {
		wf.Filter(workflowQuery)
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
