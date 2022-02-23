package workflow

import (
	aw "github.com/deanishe/awgo"
	"log"
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

func GetCredentials(wf *aw.Workflow) Auth {
	tokenID := "jira.atlassian.com"
	if wf.Config.Get("HOSTNAME") != "" {
		tokenID = wf.Config.Get("HOSTNAME")
	}

	username := "test"
	if wf.Config.Get("USERNAME") != "" {
		username = wf.Config.Get("USERNAME")
	}

	token, err := wf.Keychain.Get(tokenID)
	if err != nil {
		log.Fatal(err)
	}

	return Auth{
		Username: username,
		Password: token,
	}
}