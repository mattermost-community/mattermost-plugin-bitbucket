package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOwnerAndRepo(t *testing.T) {
	tcs := []struct {
		Full          string
		BaseURL       string
		ExpectedOwner string
		ExpectedRepo  string
	}{
		{Full: "mattermost", BaseURL: "", ExpectedOwner: "mattermost", ExpectedRepo: ""},
		{Full: "mattermost", BaseURL: "https://github.com/", ExpectedOwner: "mattermost", ExpectedRepo: ""},
		{Full: "https://bitbucket.org/mattermost", BaseURL: "", ExpectedOwner: "mattermost", ExpectedRepo: ""},
		{Full: "https://github.com/mattermost", BaseURL: "https://github.com/", ExpectedOwner: "mattermost", ExpectedRepo: ""},
		{Full: "mattermost/mattermost-server", BaseURL: "", ExpectedOwner: "mattermost", ExpectedRepo: "mattermost-server"},
		{Full: "mattermost/mattermost-server", BaseURL: "https://github.com/", ExpectedOwner: "mattermost", ExpectedRepo: "mattermost-server"},
		{Full: "https://bitbucket.org/mattermost/mattermost-server", BaseURL: "", ExpectedOwner: "mattermost", ExpectedRepo: "mattermost-server"},
		{Full: "https://github.com/mattermost/mattermost-server", BaseURL: "https://github.com/", ExpectedOwner: "mattermost", ExpectedRepo: "mattermost-server"},
		{Full: "", BaseURL: "", ExpectedOwner: "", ExpectedRepo: ""},
		{Full: "mattermost/mattermost/invalid_repo_url", BaseURL: "", ExpectedOwner: "", ExpectedRepo: ""},
		{Full: "https://github.com/mattermost/mattermost/invalid_repo_url", BaseURL: "", ExpectedOwner: "", ExpectedRepo: ""},
	}

	for _, tc := range tcs {
		_, owner, repo := parseOwnerAndRepoAndReturnFullAlso(tc.Full, tc.BaseURL)

		assert.Equal(t, tc.ExpectedOwner, owner)
		assert.Equal(t, tc.ExpectedRepo, repo)
	}
}

func TestGetYourAssigneeSearchQuery(t *testing.T) {
	result := getYourAssigneeIssuesSearchQuery("123", "testworkspace/testrepo")
	assert.Equal(t, "https://api.bitbucket.org/2.0/repositories/testworkspace/testrepo/issues?q=assignee.account_id%3D%22123%22%20AND%20state%21%3D%22closed%22",
		result)
}
