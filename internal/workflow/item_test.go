package workflow

import (
	"encoding/json"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira/jiradata"
	"github.com/nsf/jsondiff"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		issue      jiradata.Issue
		toWorkflow *aw.Workflow
	}

	wf := aw.New()

	unassignedIssue := jiradata.Issue{
		Key: "TEST-1234",
		Fields: map[string]interface{}{
			"summary": "No one has fixed this yet",
			"issuetype": map[string]interface{}{
				"name": "Bug",
			},
			"status": map[string]interface{}{
				"name": "Needs Triage",
			},
		},
	}
	unassignedIssueJSON := `{
		  "title": "TEST-1234 · Needs Triage · Unassigned",
          "subtitle": "Bug: No one has fixed this yet",
		  "text": {
			  "copy": "TEST-1234"
		  },
		  "match": "Unassigned Bug Needs Triage No one has fixed this yet TEST-1234",
  		  "arg": "https://jira.atlassian.com/browse/TEST-1234",
		  "autocomplete": "TEST-1234",
		  "icon": {
		      "path": "./issueIcons/bug.png"
		  },
		  "uid": "",
          "valid": true
        }`

	assignedIssue := jiradata.Issue{
		Key: "TEST-5678",
		Fields: map[string]interface{}{
			"summary": "Someone is working on this task at the moment",
			"issuetype": map[string]interface{}{
				"name": "Task",
			},
			"status": map[string]interface{}{
				"name": "In Progress",
			},
			"assignee": map[string]interface{}{
				"displayname": "Real Person",
			},
		},
	}
	assignedIssueJSON := `{
		  "title": "TEST-5678 · In Progress · Real Person",
          "subtitle": "Task: Someone is working on this task at the moment",
		  "text": {
			  "copy": "TEST-5678"
		  },
		  "match": "Real Person Task In Progress Someone is working on this task at the moment TEST-5678",
  		  "arg": "https://jira.atlassian.com/browse/TEST-5678",
		  "autocomplete": "TEST-5678",
		  "icon": {
			  "path": "./issueIcons/task.png"
		  },
		  "uid": "",
          "valid": true
        }`

	tests := []struct {
		name string
		args args
		want string
	}{
		{"unassigned", args{unassignedIssue, wf}, unassignedIssueJSON},
		{"assigned", args{assignedIssue, wf}, assignedIssueJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Add(&tt.args.issue, tt.args.toWorkflow)

			items := tt.args.toWorkflow.Feedback.Items
			newItem := items[len(items)-1]
			jsonData, _ := json.MarshalIndent(newItem, "", "  ")

			got := string(jsonData)

			diffOpts := jsondiff.DefaultConsoleOptions()
			res, diff := jsondiff.Compare([]byte(tt.want), []byte(got), &diffOpts)

			if res != jsondiff.FullMatch {
				t.Errorf("unexpected result: %s", diff)
			}
		})
	}
}
