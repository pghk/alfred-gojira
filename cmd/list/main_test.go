package main

import (
	"encoding/json"
	"fmt"
	"github.com/coryb/oreo"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/go-jira/jira/jiradata"
	"github.com/pghk/alfred-gojira/internal/jira"
	"log"
	"net/http"
	"testing"
)

func setup(name string, withAuthorization bool, withLog bool) (*oreo.Client, *recorder.Recorder) {
	// Start test fixture recorder
	vcr, err := recorder.New("../../test/_fixtures/" + name)
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of the http client the go-jira library normally uses, injected with our recorder
	spyClient := jira.BuildClient(jira.Auth{
		Username: "testing",
		Password: "hello this is password",
	}, withAuthorization, withLog)
	spyClient.Client.Transport = vcr

	return spyClient, vcr
}

func teardown(vcr *recorder.Recorder) {
	err := vcr.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func TestAuthorization(t *testing.T) {
	client, vcr := setup("auth", true, false)

	vcr.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		authPresent := len(r.Header["Authorization"]) != 0
		authMatches := r.Header["Authorization"][0] == i.Headers.Get("Authorization")
		return cassette.DefaultMatcher(r, i) && (authPresent && authMatches)
	})

	_, response := runQuery(client, nil)
	jsonResult, _ := json.Marshal(response)
	result := jiradata.ErrorCollection{}
	if err := json.Unmarshal(jsonResult, &result); err != nil {
		t.Fatal(err)
	}

	if result.Status != 401 {
		t.Fatal(response)
	}

	teardown(vcr)
}

func TestRunQuery(t *testing.T) {
	client, vcr := setup("search", false, false)
	jiraQuery := jira.BuildQuery("project = JRACLOUD ORDER BY assignee", "", 101)

	got, err := runQuery(client, jiraQuery)
	if err != nil {
		t.Fatal(err)
	}

	if got.Total <= 0 {
		t.Errorf("Received %d results", got.Total)
	}

	fieldChecks := []struct {
		name  string
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
	teardown(vcr)
}
