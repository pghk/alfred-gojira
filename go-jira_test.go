package main

import (
	"fmt"
	"github.com/go-jira/jira"
	"testing"

	"github.com/coryb/oreo"
	"github.com/dnaeon/go-vcr/recorder"
)

func TestGetJira(t *testing.T) {
	// Start test fixture recorder
	vcr, err := recorder.New("testdata/fixtures/test")
	if err != nil {
		t.Fatal(err)
	}
	// Stop recorder once done
	defer func(vcr *recorder.Recorder) {
		err := vcr.Stop()
		if err != nil {
			t.Fatal(err)
		}
	}(vcr)

	// Create an instance of the http client the go-jira library normally uses, injected with our recorder
	recordingHttpClient := oreo.New()
	recordingHttpClient.Client.Transport = vcr

	// Inject our modified client into the new instance provided by the go-jira library
	jiraClient := *jira.NewJira("https://jira.atlassian.com")
	jiraClient.UA = recordingHttpClient

	list, err := GetJiraList(jiraClient)
	if err != nil {
		t.Fatal(err)
	}
	got := list.Total

	fmt.Println(list.Issues[0])

	if got <= 0  {
		t.Errorf("got %q ", got)
	}
}