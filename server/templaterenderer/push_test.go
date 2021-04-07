package templaterenderer

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestPushTemplates(t *testing.T) {
	tr := MakeTemplateRenderer()
	tr.RegisterBitBucketAccountIDToUsernameMappingCallback(bitBucketAccountIDToUsernameMappingTestCallback)

	t.Run("pushed", func(t *testing.T) {
		t.Run("single commit", func(t *testing.T) {
			expected := `
User @testMmUser pushed [1 new commit](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e) to [\[mattermost/mattermost-plugin-bitbucket:master\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/master):
[\[dca554\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd) edit readme - [testnickname](https://bitbucket.org/test-testnickname-url/)
`

			pl := getTestRepoPushPayloadWithOneCommit()

			actual, err := tr.RenderRepoPushEventNotificationForSubscribedChannels(pl)

			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})

		t.Run("two commits", func(t *testing.T) {
			expected := `
User @testMmUser pushed [2 new commits](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branches/compare/dca5546b6b1419ff71adcada81b457caf3dcbdcd..54ec7b7ec732bc97278ec82e2c50cfc260918f3e) to [\[mattermost/mattermost-plugin-bitbucket:master\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/master):
[\[dca554\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/dca5546b6b1419ff71adcada81b457caf3dcbdcd) edit readme - [testnickname](https://bitbucket.org/test-testnickname-url/)
[\[fd84b9\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/commits/fd84b9bfa6b415461f2c6559c262ed0826f76e08) edit something else - @testMmUser
`

			pl := getTestRepoPushPayloadWithTwoCommits()

			actual, err := tr.RenderRepoPushEventNotificationForSubscribedChannels(pl)

			require.NoError(t, err)
			require.Equal(t, expected, actual)
		})
	})

	t.Run("RenderBranchOrTagCreatedEventNotificationForSubscribedChannels branch", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Branch [test-new-branch](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-new-branch) " +
			"was created by [mmUserBitbucketNickName](https://bitbucket.org/test-mmUserBitbucketNickName-url/)\n"

		pl := getTestRepoPushBranchCreated()

		actual, err := tr.RenderBranchOrTagCreatedEventNotificationForSubscribedChannels(pl)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderBranchOrTagCreatedEventNotificationForSubscribedChannels tag", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Tag [test-new-tag](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-new-tag) " +
			"was created by [mmUserBitbucketNickName](https://bitbucket.org/test-mmUserBitbucketNickName-url/)\n"

		pl := getTestRepoPushTagCreated()

		actual, err := tr.RenderBranchOrTagCreatedEventNotificationForSubscribedChannels(pl)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderBranchOrTagDeletedEventNotificationForSubscribedChannels branch", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Branch [test-old-branch](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/branch/test-old-branch) " +
			"was deleted by [mmUserBitbucketNickName](https://bitbucket.org/test-mmUserBitbucketNickName-url/)\n"

		pl := getTestRepoPushBranchDeleted()

		actual, err := tr.RenderBranchOrTagDeletedEventNotificationForSubscribedChannels(pl)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderBranchOrTagDeletedEventNotificationForSubscribedChannels tag", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Tag [test-old-tag](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/tag/test-old-tag) " +
			"was deleted by [mmUserBitbucketNickName](https://bitbucket.org/test-mmUserBitbucketNickName-url/)\n"

		pl := getTestRepoPushTagDeleted()

		actual, err := tr.RenderBranchOrTagDeletedEventNotificationForSubscribedChannels(pl)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
