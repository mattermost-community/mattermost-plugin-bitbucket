package templaterenderer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonTemplates(t *testing.T) {
	tr := templateRenderer{}
	tr.init()
	tr.RegisterBitBucketAccountIDToUsernameMappingCallback(bitBucketAccountIDToUsernameMappingTestCallback)

	for name, tc := range map[string]struct {
		expected string
		template string
		payload  interface{}
	}{
		"repo": {
			payload:  getTestRepository(),
			template: "repo",
			expected: "[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket)",
		},
		"repoPullRequestWithTitle": {
			payload:  getTestPullRequestCreatedPayload(),
			template: "repoPullRequestWithTitle",
			expected: "[mattermost-plugin-bitbucket#1](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title",
		},
		"pullRequest": {
			payload:  getTestPullRequest(),
			template: "pullRequest",
			expected: "[#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1)",
		},
		"issue": {
			payload:  getTestIssueCreatedPayload(),
			template: "issue",
			expected: "[\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)",
		},
		"non MM user": {
			payload:  getTestOwnerThatDoesntHaveAccount(),
			template: "user",
			expected: "[testnickname](https://bitbucket.org/test-testnickname-url/)",
		},
		"MM user": {
			payload:  getTestOwnerThatHasMmAccount(),
			template: "user",
			expected: "@testMmUser",
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual, err := tr.renderTemplateWithName(tc.template, tc.payload)

			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
