package workflow

import (
	aw "github.com/deanishe/awgo"
	"github.com/pghk/alfred-gojira/internal/jira"
	"go.deanishe.net/fuzzy"
)

const ConfigTriggerName = "configure"

type Configuration struct {
	Hostname    string
	Username    string
	Privatehost bool
	MaxResults int
}

var (
	config   *Configuration
	workflow *aw.Workflow
	_        aw.MagicAction = configure{}
)

func init() {
	config = &Configuration{
		Hostname:    "",
		Username:    "",
		Privatehost: true,
		MaxResults: 100,
	}
}

// The configure action sends the user to the "configure" script filter
type configure struct {
	aw.MagicAction
}

func (configure) Keyword() string     { return "config" }
func (configure) Description() string { return "Edit workflow configuration" }
func (configure) RunText() string     { return "Config action registered." }
func (configure) Run() error {
	return workflow.Alfred.RunTrigger(ConfigTriggerName, "")
}

func BuildWorkflow(sortOptions []fuzzy.Option) *aw.Workflow {
	workflow = aw.New(
		aw.AddMagic(configure{}),
		aw.SortOptions(sortOptions...),
	)

	// Update default config from environment variables
	if err := workflow.Config.To(config); err != nil {
		panic(err)
	}

	return workflow
}

func GetJiraHostname() string {
	return config.Hostname
}

func GetJiraUsername() string {
	return config.Username
}

func GetMaxResultSetting() int {
	return config.MaxResults
}

func CredentialsRequired() bool {
	return config.Privatehost
}

func GetCredentials() (jira.Auth, error) {
	tokenStorageKey := config.Hostname
	token, err := workflow.Keychain.Get(tokenStorageKey)
	if err != nil {
		return jira.Auth{}, err
	}

	return jira.Auth{
		Username: config.Username,
		Password: token,
	}, nil
}
