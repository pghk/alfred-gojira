package main

import (
	"fmt"
	"testing"
)

func TestGetJira(t *testing.T) {
	list, err := GetJiraList()
	if err != nil {
		t.Fatal(err)
	}
	got := list.Total

	fmt.Println(list.Issues[0])

	if got <= 0  {
		t.Errorf("got %q ", got)
	}
}