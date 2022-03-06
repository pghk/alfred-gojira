# Alfred Gojira
Navigate your Jira issues in [Alfred](https://www.alfredapp.com/).

![ScreenShot](https://raw.github.com/pghk/alfred-gojira/main/assets/ScreenShot.png)

## Usage
- Launch the workflow to view a list of issues obtained according to your configured query
- Start typing input to filter the list
- Press enter `↩` to open the url of the highlighed issue in your default browser
- Press command `⌘` + `c` to copy the key of the highlighed issue
- Type `workflow:config` (with the workflow open) to enter configuration options

## Configuration
This workflow provides an external trigger so that it can be launched from a hotkey of your choice without needing to alter the workflow itself. To set this up, in a workflow of your own, create a hotkey trigger or keyword input, and connect it to an external trigger output to Workflow ID `com.pghk.gojira` and Trigger ID `issues`.

When launched for the first time, a list of configuration options will appear. On subsequent launches this can be re-visted via a query of `workflow:config`. You can also edit these options from Alfred's workflow config UI.

### Options
**Hostname**: the base URL of your Jira host, i.e. `jira.atlassian.com`

If your Jira host is public (as in the above example), set the **Private Host** variable to `false` or `0`, and you'll be able to proceed to the main workflow without any further setup.

**Username**: this should be your email.

**API Token**: in order to authenticate to a private Jira host, you must [create an API token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/) and allow this workflow to use it. This will not be stored in the workflows settings; the workflow expects to find it your macOS keychain (under an "account" value equal to what you provided in the **Hostname** variable, and `com.pghk.gojira` as "name" and "where"). You can place your token in the Keychain Access application yourself, and grant access once Alfred requests it, or you can provide your token to the configuration option to have it placed in the keychain for you.

**Max results**: the total number of issues matching your query to request from the Jira API. Whatever you set here, issues will be loaded in the background in pages of 100 at a time, and held in the workflow's cache for 3 hours.

### Query
To customize the scope of issues listed, modify the script value of this workflow's "gojira" script filter [with your own custom query](https://support.atlassian.com/jira-software-cloud/docs/use-advanced-search-with-jira-query-language-jql/):

`/list $1 -query "your JQL here"`

The default query is simply for open issues: `resolution = unresolved ORDER BY updated ASC`.

## Roadmap
1. Move the query configuration to a variable, or 
    1. as input from an external trigger, allowing users to configure multiple queries via different triggers 

## Big thanks 🙏🏼
- The [AwGo library](https://github.com/deanishe/awgo), used to provide cached, filterable feedback to Alfred.
- The [go-jira command line tool](https://github.com/go-jira/jira) (implemented as a library), used to communicate with the Jira API.
