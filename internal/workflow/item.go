package workflow

import (
	"encoding/json"
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"log"
)

func convertType(input interface{}, output interface{}) {
	var (
		err error
		data []byte
	)

	// Convert to JSON from map[string]interface{}
	data, err = json.Marshal(input)
	if err == nil {
		// Populate struct from JSON
		err = json.Unmarshal(data, &output)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func populate(defaultInstance interface{}, fromIssue jiradata.Issue, field string) {
	fields := fromIssue.Fields
	if _, present := fields[field]; present {
		convertType(fields[field], &defaultInstance)
	}
}

func Add(issue *jiradata.Issue, toWorkflow *aw.Workflow) {
	defaultIssueType := jiradata.IssueType{}
	defaultIssueStatus := jiradata.Status{}
	defaultAssignee := jiradata.User{
		DisplayName: "Unassigned",
	}

	populate(&defaultIssueType, *issue, "issuetype")
	populate(&defaultIssueStatus, *issue, "status")
	populate(&defaultAssignee, *issue, "assignee")

	key := issue.Key
	status := defaultIssueStatus.Name
	assignee := defaultAssignee.DisplayName
	issueType := defaultIssueType.Name
	summary := issue.Fields["summary"].(string)

	issueUrlBase := "https://" + GetJiraHostname(toWorkflow) + "/browse/"

	toWorkflow.NewItem(fmt.Sprintf("%s · %s · %s", key, status, assignee)).
		Subtitle(fmt.Sprintf("%s: %s", issueType, summary)).
		Match(fmt.Sprintf("%s %s", key, summary)).
		Arg(fmt.Sprint(issueUrlBase + key)).
		Autocomplete(key).
		Valid(true).
		UID("")
}