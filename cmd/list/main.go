package main

import (
	"fmt"
	"gojiralfredo/internal/jira-client"
)

func main() {
	fmt.Println(jira.ListIssues())
}
