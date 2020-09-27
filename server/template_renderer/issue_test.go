package template_renderer

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestIssueTemplates(t *testing.T) {

	tr := MakeTemplateRenderer()
	tr.RegisterBitBucketAccountIDToUsernameMappingCallback(bitBucketAccountIDToUsernameMappingTestCallback)

	t.Run("RenderIssueCreatedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n#### README.md is outdated" +
			"\n##### [\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)" +
			"\n#new-issue by @testMmUser:" +
			"\n>README.md should be updated\n"

		actual, err := tr.RenderIssueCreatedEventNotificationForSubscribedChannels(getTestIssueCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueDescriptionMentionNotification", func(t *testing.T) {
		expected := "\n@testMmUser " +
			"mentioned you on [\\[mattermost-plugin-bitbucket#1\\]]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated):" +
			"\n>README.md should be updated\n"

		actual, err := tr.RenderIssueDescriptionMentionNotification(getTestIssueCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueUpdatedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n#### README.md is outdated" +
			"\n##### [\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)" +
			"\n#updated-issue by @testMmUser:" +
			"\n>README.md should be updated\n"

		actual, err := tr.RenderIssueUpdatedEventNotificationForSubscribedChannels(getTestIssueUpdatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueAssignmentNotificationForAssignedUser", func(t *testing.T) {
		expected := "\n@testMmUser assigned you to issue " +
			"[\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)\n"

		actual, err := tr.RenderIssueAssignmentNotificationForAssignedUser(getTestIssueUpdatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueStatusUpdateNotificationForIssueReporter", func(t *testing.T) {
		expected := "\n@testMmUser set status to " +
			"closed of your issue [\\[mattermost-plugin-bitbucket#1\\]]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated)\n"

		actual, err := tr.RenderIssueStatusUpdateNotificationForIssueReporter(getTestIssueUpdatedPayloadWithStatusChange())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueCommentCreatedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"New comment by @testMmUser on " +
			"[\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated):" +
			"\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderIssueCommentCreatedEventNotificationForSubscribedChannels(getTestIssueCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueCommentNotificationForIssueReporter", func(t *testing.T) {
		expected := "\n@testMmUser commented on your issue " +
			"[\\[mattermost-plugin-bitbucket#1\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated):" +
			"\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderIssueCommentNotificationForIssueReporter(getTestIssueCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderIssueCommentMentionNotification", func(t *testing.T) {
		expected := "\n@testMmUser mentioned you on [\\[mattermost-plugin-bitbucket#1\\]]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/issues/1/readme-is-outdated):" +
			"\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderIssueCommentMentionNotification(getTestIssueCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
