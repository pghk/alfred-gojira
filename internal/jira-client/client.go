package jira

import (
	"encoding/base64"
	"fmt"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"gojiralfredo/internal/workflow"
	"log"
	"net/http"
)

type Consumer struct {
	*jira.Jira
}

type Query struct {
	*jira.SearchOptions
}

type debugLogger struct{}

func (l *debugLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

var logger = debugLogger{}

type urlLogger struct {
	wrapped http.Transport
}

func (u *urlLogger) RoundTrip(req *http.Request) (res *http.Response, e error) {
	log.Printf("Requesting %v\n", req.URL)
	return u.wrapped.RoundTrip(req)
}

func BuildClient(auth workflow.Auth, authorize bool, verbose bool) *oreo.Client {
	rawAuth := fmt.Sprintf("%s:%s", auth.Username, auth.Password)
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(rawAuth))
	authHeaderVal := fmt.Sprintf("Basic %s", encodedAuth)

	client := oreo.New()

	if authorize {
		client = client.
			WithPreCallback(func(req *http.Request) (*http.Request, error) {
				req.Header.Add("Authorization", authHeaderVal)
				return req, nil
			})
	}

	if verbose {
		client = client.WithLogger(&logger).WithTrace(true)
	} else {
		client = client.WithTransport(&urlLogger{})
	}

	return client
}

func BuildQuery() *Query {
	return &Query{&jira.SearchOptions{}}
}

func BuildConsumer(endpoint string, client *oreo.Client) *Consumer {
	consumer := &Consumer{jira.NewJira(endpoint)}
	consumer.UA = client
	return consumer
}
