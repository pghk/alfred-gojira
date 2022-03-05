package main

import (
	"flag"
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/pghk/alfred-gojira/internal/workflow"
	"go.deanishe.net/fuzzy"
	"log"
	"strconv"
)

var (
	wf              *aw.Workflow
	tokenStorageKey string
	setKey, getKey  string
)

func init() {
	wf = workflow.BuildWorkflow([]fuzzy.Option{})
	tokenStorageKey = workflow.GetJiraHostname()
	flag.StringVar(&setKey, "set", "", "save a value")
	flag.StringVar(&getKey, "get", "", "enter a new value")
}

func runSet(key, value string) {
	wf.Configure(aw.TextErrors(true))

	log.Printf("saving %#v to %s ...", value, key)

	if key == "API_TOKEN" {
		if err := wf.Keychain.Set(tokenStorageKey, value); err != nil {
			wf.FatalError(err)
		}
	} else if key == "PRIVATEHOST" {
		_, err := strconv.ParseBool(value)
		if err != nil {
			wf.FatalError(err)
		}
		if err := wf.Config.Set(key, value, false).Do(); err != nil {
			wf.FatalError(err)
		}
	} else {
		if err := wf.Config.Set(key, value, false).Do(); err != nil {
			wf.FatalError(err)
		}
	}

	if err := wf.Alfred.RunTrigger(workflow.ConfigTriggerName, ""); err != nil {
		wf.FatalError(err)
	}

	log.Printf("saved %#v to %s", value, key)
}

// get new value for setting from user via Script Filter
func runGet(key, value string) {
	log.Printf("getting new %s ...", key)

	if value != "" {
		var varname string
		switch key {
		case "API token":
			varname = "API_TOKEN"
		case "hostname":
			varname = "HOSTNAME"
		case "username":
			varname = "USERNAME"
		case "private host setting":
			varname = "PRIVATEHOST"
		case "max results setting":
			varname = "MAX_RESULTS"
		}

		wf.NewItem(fmt.Sprintf("Set %s to “%s”", key, value)).
			Subtitle("↩ to save").
			Arg(value).
			Valid(true).
			Var("value", value).
			Var("varname", varname)
	}

	wf.WarnEmpty(fmt.Sprintf("Enter %s", key), "")
	wf.SendFeedback()
}

func run() {
	needCredentials := workflow.CredentialsRequired()
	haveCredentials := false
	if needCredentials {
		_, err := wf.Keychain.Get(tokenStorageKey)
		if err == nil {
			haveCredentials = true
		}
	}

	wf.Args()

	flag.Parse()
	query := flag.Arg(0)

	if setKey != "" {
		runSet(setKey, query)
		return
	}

	if getKey != "" {
		runGet(getKey, query)
		return
	}

	wf.NewItem("Hostname: "+workflow.GetJiraHostname()).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "hostname")

	wf.NewItem("Username: "+workflow.GetJiraUsername()).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "username")

	wf.NewItem("Max results: "+strconv.Itoa(workflow.GetMaxResultSetting())).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "max results setting")

	wf.NewItem("Private host: "+strconv.FormatBool(needCredentials)).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "private host setting")

	if needCredentials && !haveCredentials {
		wf.NewItem("API token: not set").
			Subtitle("↩ to add to Keychain").
			Valid(true).
			Var("name", "API token")
	}

	if haveCredentials {
		wf.NewItem("API token: set in macOS Keychain").
			Subtitle("↩ to edit in keychain").
			Valid(true).
			Var("name", "API token")
	}

	if haveCredentials || !needCredentials {
		wf.NewItem("Search Jira Issues").
			Subtitle("↩ to use workflow as configured").
			Valid(true).
			Var("trigger", "issues")
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Items", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
