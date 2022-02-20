package main

import (
	"github.com/coryb/oreo"
	"github.com/dnaeon/go-vcr/recorder"
	"gojiralfredo/internal/jira-client"
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
	recordingHttpClient := oreo.New()
	recordingHttpClient.Client.Transport = vcr

	// Inject our modified client into the new instance provided by the go-jira library
	client = *jira.BuildClient("https://jira.atlassian.com")
	client.UA = recordingHttpClient

	return vcr
}

func shutdown(vcr *recorder.Recorder) {
	err := vcr.Stop()
	if err != nil {
		log.Fatal(err)
	}
}

func TestRunQuery(t *testing.T) {
	got := runQuery()
	//fmt.Println(got.Issues[0])
	if got.Total <= 0  {
		t.Errorf("Received %q results", got.Total)
	}
}

func TestMain(m *testing.M) {
	vcr := setup()
	code := m.Run()
	shutdown(vcr)
	os.Exit(code)
}