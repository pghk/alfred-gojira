package jira

import (
	"encoding/base64"
	"fmt"
	"github.com/coryb/oreo"
	aw "github.com/deanishe/awgo"
	"github.com/go-jira/jira"
	"log"
	"net/http"
)

type Consumer struct {
	*jira.Jira
}

type Query struct {
	*jira.SearchOptions
}

type Auth struct {
	Username string
	Password string
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

func BuildClient(auth Auth) *oreo.Client {
	rawAuth := fmt.Sprintf("%s:%s", auth.Username, auth.Password)
	encodedAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(rawAuth)))
	return oreo.New().WithPreCallback(func(req *http.Request) (*http.Request, error) {
		req.Header.Add("Authorization", encodedAuth)
		return req, nil
	})
}

func BuildQuery() *Query {
	return &Query{&jira.SearchOptions{}}
}

func BuildConsumer(endpoint string, client *oreo.Client) *Consumer {
	consumer := &Consumer{jira.NewJira(endpoint)}
	consumer.UA = client
	return consumer
}
