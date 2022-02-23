package main

import (
	"fmt"
	"github.com/dnaeon/go-vcr/recorder"
	"gojiralfredo/internal/jira-client"
	"gojiralfredo/internal/workflow"
	"log"
	"os"
	"testing"
)

func setup() *recorder.Recorder {
	// Start test fixture recorder
	vcr, err := recorder.New("../../test/_data/fixtures/search")
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of the http client the go-jira library normally uses, injected with our recorder
	spyClient := jira.BuildClient(workflow.Auth{
		Username: "testing",
		Password: "hello this is password",
	}, false, true)
	spyClient.Client.Transport = vcr

	// Inject our modified client into the new instance provided by the go-jira library
	setClient(spyClient)
	setConsumer()

	return vcr
}

func shutdown(vcr *recorder.Recorder) {
	err := vcr.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func TestRunQuery(t *testing.T) {
	query.MaxResults = 20
	query.Project = "JRACLOUD"
	query.Sort = "assignee"

	got := runQuery()
	//fmt.Println(got.Issues[0])
	if got.Total <= 0  {
		t.Errorf("Received %d results", got.Total)
	}

	fieldChecks := []struct {
		name string
		field string
	}{
		{"Summary", "summary"},
		{"Issue Type", "issuetype"},
		{"Status", "status"},
		{"Assignee", "assignee"},
	}
	for _, tt := range fieldChecks {
		t.Run(fmt.Sprintf("Contains%sField", tt.name), func(t *testing.T) {
			for _, issue := range got.Issues {
				_, ok := issue.Fields[tt.field]
				if !ok {
					t.Fatalf("Issue doesn't contain %s field", tt.name)
				}
			}
		})
	}

}

func TestMain(m *testing.M) {
	vcr := setup()
	code := m.Run()
	shutdown(vcr)
	os.Exit(code)
}