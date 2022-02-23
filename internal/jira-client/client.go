package jira

import (
	"encoding/base64"
	"fmt"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"gojiralfredo/internal/workflow"
	"net/http"
)

type Consumer struct {
	*jira.Jira
}

type Query struct {
	*jira.SearchOptions
}

func BuildClient(auth workflow.Auth) *oreo.Client {
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
