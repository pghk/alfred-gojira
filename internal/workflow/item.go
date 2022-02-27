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
		err  error
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

func getIcon(issueType string) *aw.Icon {
	icon := &aw.Icon{Type: ""}
	path := "./issueIcons/"
	switch issueType {
	case "Bug":
		icon.Value = path + "bug.png"
	case "Epic":
		icon.Value = path + "epic.png"
	case "Improvement":
		icon.Value = path + "new_feature.png"
	case "Task":
		icon.Value = path + "task.png"
	case "Story", "Suggestion":
		icon.Value = path + "story.png"
	case "Sub-task":
		icon.Value = path + "subtask.png"
	case "Question/Research":
		icon.Value = path + "question.png"
	default:
		icon.Value = path + "generic_issue.png"
	}
	return icon
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

	issueUrlBase := "https://" + GetJiraHostname() + "/browse/"

	toWorkflow.NewItem(fmt.Sprintf("%s · %s · %s", key, status, assignee)).
		Subtitle(fmt.Sprintf("%s: %s", issueType, summary)).
		Match(fmt.Sprintf("%s %s %s %s %s", assignee, issueType, status, summary, key)).
		Arg(fmt.Sprint(issueUrlBase + key)).
		Autocomplete(key).
		Copytext(key).
		Valid(true).
		Icon(getIcon(issueType)).
		UID("")
}
