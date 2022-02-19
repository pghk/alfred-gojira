// Copyright (c) 2018 Dean Jackson <deanishe@deanishe.net>
// MIT Licence - http://opensource.org/licenses/MIT

/*
Workflow settings demonstrates binding a struct to Alfred's settings.

The workflow's settings are stored in info.plist/the workflow's
configuration sheet in Alfred Preferences.

These are imported into the Server struct using Config.To().

The Script Filter displays these settings, and you can select one
to change its value.

If you enter a new value, this is saved to info.plist/the configuration
sheet via Config.Set(), and the workflow is run again by calling
its "settings" External Trigger via Alfred.RunTrigger().
*/
package main

import (
	"flag"
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
)

// Server contains the configuration loaded from the workflow's settings
// in the configuration sheet.
type Server struct {
	Hostname   string
	Username   string
}

var (
	srv *Server
	wf  *aw.Workflow
	// Command-line arguments
	setKey, getKey string
)

func init() {
	wf = aw.New()
	flag.StringVar(&setKey, "set", "", "save a value")
	flag.StringVar(&getKey, "get", "", "enter a new value")
}

// save setting to info.plist via Alfred's AppleScript API
func runSet(key, value string) {
	wf.Configure(aw.TextErrors(true))

	log.Printf("saving %#v to %s ...", value, key)

	if key == "API_TOKEN" {
		if err := wf.Keychain.Set(srv.Hostname, value); err != nil {
			wf.FatalError(err)
		}
	} else {
		if err := wf.Config.Set(key, value, false).Do(); err != nil {
			wf.FatalError(err)
		}
	}

	if err := wf.Alfred.RunTrigger("settings", ""); err != nil {
		wf.FatalError(err)
	}

	log.Printf("saved %#v to %s", value, key)
}

// get new value for setting from user via Script Filter
func runGet(key, value string) {
	log.Printf("getting new %s ...", key)

	if value != "" {
		var varname string
		switch key {
		case "API token":
			varname = "API_TOKEN"
		case "hostname":
			varname = "HOSTNAME"
		case "username":
			varname = "USERNAME"
		}

		wf.NewItem(fmt.Sprintf("Set %s to “%s”", key, value)).
			Subtitle("↩ to save").
			Arg(value).
			Valid(true).
			Var("value", value).
			Var("varname", varname)
	}

	wf.WarnEmpty(fmt.Sprintf("Enter %s", key), "")
	wf.SendFeedback()
}

func run() {
	wf.Args() // call to handle magic actions

	// ----------------------------------------------------------------
	// Load configuration

	// Default configuration
	srv = &Server{
		Hostname:   "localhost",
		Username:   "anonymous",
	}

	// Update config from environment variables
	if err := wf.Config.To(srv); err != nil {
		panic(err)
	}

	log.Printf("loaded: %#v", srv)

	// ----------------------------------------------------------------
	// Parse command-line flags and decide what to do

	flag.Parse()
	query := flag.Arg(0)

	if setKey != "" {
		runSet(setKey, query)
		return
	}

	if getKey != "" {
		runGet(getKey, query)
		return
	}

	// ----------------------------------------------------------------
	// Show available settings.

	wf.NewItem("Hostname: "+srv.Hostname).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "hostname")

	wf.NewItem("Username: "+srv.Username).
		Subtitle("↩ to edit").
		Valid(true).
		Var("name", "username")

	_, err := wf.Keychain.Get(srv.Hostname)
	if err != nil {
		wf.NewItem("API token: not set").
			Subtitle("↩ to add to Keychain").
			Valid(true).
			Var("name", "API token")
	} else {
		wf.NewItem("API token: set in macOS Keychain").
			Subtitle("↩ to edit in keychain").
			Valid(true).
			Var("name", "API token")
	}


	if query != "" {
		wf.Filter(query)
	}

	wf.WarnEmpty("No Matching Items", "Try a different query?")
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
