package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/go-playground/webhooks.v5/bitbucket"
)

func getActor() bitbucket.Owner {
	actor := bitbucket.Owner{NickName: "testnickname"}
	actor.Links.HTML.Href = "https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/"
	return actor
}

func getRepository() bitbucket.Repository {
	repository := bitbucket.Repository{FullName: "mattermost-plugin-bitbucket"}
	repository.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket"
	return repository
}

func getPullRequest() bitbucket.PullRequest {
	pullRequest := bitbucket.PullRequest{
		ID:    1,
		Title: "Test title",
	}
	pullRequest.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1"

	return pullRequest
}

func getComment() bitbucket.Comment {
	comment := bitbucket.Comment{}
	comment.Content.Raw = "test comment"

	return comment
}

func getIssue() bitbucket.Issue {
	issue := bitbucket.Issue{}
	issue.ID = 1
	issue.Title = "README.md is outdated"
	issue.Content.Raw = "README.md should be updated"
	issue.Links.HTML.Href = "https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated"

	return issue
}

func TestSimpleTemplates(t *testing.T) {
	t.Run("user", func(t *testing.T) {
		expected := "[testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/)"

		actual, err := renderTemplate("user", getActor())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("repo", func(t *testing.T) {
		expected := `[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket)`

		actual, err := renderTemplate("repo", getRepository())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestIssueTemplates(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		expected := `
#### README.md is outdated
##### [\[mattermost-plugin-bitbucket#1\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)
#new-issue by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/)

README.md should be updated
`

		payload := bitbucket.IssueCreatedPayload{
			Repository: getRepository(),
			Actor:      getActor(),
			Issue:      getIssue(),
		}

		actual, err := renderTemplate("issueCreated", payload)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("updated", func(t *testing.T) {
		expected := `
#### README.md is outdated
##### [\[mattermost-plugin-bitbucket#1\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)
#updated-issue by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/)

README.md should be updated
`

		payload := bitbucket.IssueUpdatedPayload{
			Repository: getRepository(),
			Actor:      getActor(),
			Issue:      getIssue(),
		}

		actual, err := renderTemplate("issueUpdated", payload)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("comment created", func(t *testing.T) {
		expected := `
[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) New comment by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/) on [\[mattermost-plugin-bitbucket#1\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated):

test comment
`
		payload := bitbucket.IssueCommentCreatedPayload{
			Repository: getRepository(),
			Actor:      getActor(),
			Issue:      getIssue(),
			Comment:    getComment(),
		}

		actual, err := renderTemplate("issueCommentCreated", payload)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestPullRequestTemplates(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		expected := `
[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) was created by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/).
`

		payload := bitbucket.PullRequestCreatedPayload{
			Repository:  getRepository(),
			PullRequest: getPullRequest(),
			Actor:       getActor(),
		}

		actual, err := renderTemplate("prCreated", payload)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("comment created", func(t *testing.T) {
		expected := `
[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) New review comment by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/) on [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1):

test comment
`

		payload := bitbucket.PullRequestCommentCreatedPayload{
			Repository:  getRepository(),
			PullRequest: getPullRequest(),
			Actor:       getActor(),
			Comment:     getComment(),
		}

		actual, err := renderTemplate("prCommentCreated", payload)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("approved", func(t *testing.T) {
		expected := `
[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) was approved by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/).
`

		payload := bitbucket.PullRequestApprovedPayload{
			Repository:  getRepository(),
			PullRequest: getPullRequest(),
			Actor:       getActor(),
		}

		actual, err := renderTemplate("prApproved", payload)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("merged", func(t *testing.T) {
		expected := `
[\[mattermost-plugin-bitbucket\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) was merged by [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/).
`

		payload := bitbucket.PullRequestMergedPayload{
			Repository:  getRepository(),
			PullRequest: getPullRequest(),
			Actor:       getActor(),
		}

		actual, err := renderTemplate("prMerged", payload)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestPushTemplate(t *testing.T) {
	t.Run("single commit", func(t *testing.T) {
		expected := `
User [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/) pushed [1 new commit](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e):
[\[dca554\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd) edit readme
`

		payload, err := loadRepoPushPayloadFromFile("push_1_commit_payload.json")
		require.NoError(t, err)

		actual, err := renderTemplate("pushed", payload)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("two commits", func(t *testing.T) {
		expected := `
User [testnickname](https://bitbucket.org/%7B4f86ef3-b0a7-4a39-a118-d3e06e31981f%7D/) pushed [2 new commits](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e):
[\[dca554\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd) edit readme
[\[fd84b9\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/fd84b9bfa6b415461f2c6559c262ed0826f76e08) edit something else
`

		payload, err := loadRepoPushPayloadFromFile("push_2_commits_payload.json")
		require.NoError(t, err)

		actual, err := renderTemplate("pushed", payload)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

// had some issues with creating RepoPushPayload, so used json instead
func loadRepoPushPayloadFromFile(filename string) (*bitbucket.RepoPushPayload, error) {
	var payload bitbucket.RepoPushPayload
	f, err := ioutil.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(f, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
