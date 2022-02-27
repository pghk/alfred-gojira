package workflow

import (
	aw "github.com/deanishe/awgo"
)

type Auth struct {
	Username string
	Password string
}

func GetJiraHostname(wf *aw.Workflow) string {
	if wf.Config.Get("HOSTNAME") != "" {
		return wf.Config.Get("HOSTNAME")
	}
	return "jira.atlassian.com"
}

func GetJiraUsername(wf *aw.Workflow) string {
	if wf.Config.Get("USERNAME") != "" {
		return wf.Config.Get("HOSTNAME")
	}
	return ""
}

func GetCredentials(wf *aw.Workflow, fallback func(wf *aw.Workflow)) Auth {
	tokenID := GetJiraHostname(wf)
	username := GetJiraUsername(wf)

	token, tokenErr := wf.Keychain.Get(tokenID)
	if tokenErr != nil {
		fallback(wf)
	}

	return Auth{
		Username: username,
		Password: token,
	}
}
