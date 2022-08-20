package main

import (
	"flag"
	"fmt"
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"github.com/pghk/alfred-gojira/internal/jira"
	"github.com/pghk/alfred-gojira/internal/workflow"
	"go.deanishe.net/fuzzy"
	"log"
	"net/url"
	"os"
	"os/exec"
	"regexp"
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
			credentials, err := workflow.GetCredentials()
			if err != nil {
				openConfigEditor()
				log.Fatal(err)
			} else {
				client = jira.BuildClient(credentials, true, false)
			}
		} else {
			client = jira.BuildClient(jira.Auth{}, false, false)
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
		jiraQuery = jira.BuildQuery(query, fields, workflow.GetMaxResultSetting())
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

func addFallbackItems() {
	keySubStrings := parseIssueKey(filter)
	queryIsIssueKey := keySubStrings != nil
	if queryIsIssueKey {
		addIssueKeyFallbacks(keySubStrings)
	} else {
		addSearchFallback()
	}
}

func parseIssueKey(query string) []string {
	jiraKeyPattern := regexp.MustCompile("([A-Z0-9]+-)?(\\d+)")
	return jiraKeyPattern.FindStringSubmatch(query)
}

func addIssueKeyFallbacks(keySubStrings []string) {
	fullString, projectKey, issueNumber := keySubStrings[0], keySubStrings[1], keySubStrings[2]
	if projectKey != "" {
		issueKey := fullString
		addIssueKeyFallbackItem(issueKey)
	} else {
		projectsInScope := workflow.GetProjectList()
		for _, projectChoice := range projectsInScope {
			issueKey := fmt.Sprintf("%s-%s", projectChoice, issueNumber)
			addIssueKeyFallbackItem(issueKey)
		}
	}
}

func addIssueKeyFallbackItem(issueKey string) {
	issueURL := fmt.Sprintf("https://%s/browse/%s", workflow.GetJiraHostname(), issueKey)
	wf.NewItem(fmt.Sprintf("Go to issue %s", issueKey)).Subtitle(issueURL).Arg(issueURL).Valid(true)
}

func addSearchFallback() {
	searchQuery := fmt.Sprintf("(summary ~ %[1]s OR description ~ %[1]s)", filter)
	if scope := workflow.GetProjectString(); scope != "" {
		searchQuery += fmt.Sprintf(" AND project IN (%s)", scope)
	}
	searchURL := fmt.Sprintf(
		"https://%s/issues/?jql=%s",
		workflow.GetJiraHostname(),
		url.QueryEscape(searchQuery),
	)
	wf.NewItem(fmt.Sprintf("Search issues for %s", searchURL)).Subtitle(searchURL).Arg(searchURL).Valid(true)
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

	if wf.IsEmpty() {
		addFallbackItems()
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
