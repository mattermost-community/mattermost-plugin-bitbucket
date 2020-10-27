package templaterenderer

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestPrTemplates(t *testing.T) {
	tr := MakeTemplateRenderer()
	tr.RegisterBitBucketAccountIDToUsernameMappingCallback(bitBucketAccountIDToUsernameMappingTestCallback)

	t.Run("RenderPullRequestCreatedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) " +
			"was created by @testMmUser\n"

		actual, err := tr.RenderPullRequestCreatedEventNotificationForSubscribedChannels(getTestPullRequestCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestAssignedNotification", func(t *testing.T) {
		expected := "\n@testMmUser assigned you to pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title\n"

		actual, err := tr.RenderPullRequestAssignedNotification(getTestPullRequestUpdatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestDeclinedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) " +
			"was declined by @testMmUser\n"

		actual, err := tr.RenderPullRequestDeclinedEventNotificationForSubscribedChannels(getTestPullRequestDeclinedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestDeclinedNotificationForPullRequestAuthor", func(t *testing.T) {
		expected := "\n@testMmUser declined your pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title\n"

		actual, err := tr.RenderPullRequestDeclinedNotificationForPullRequestAuthor(getTestPullRequestDeclinedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestApprovedNotificationForPullRequestAuthor", func(t *testing.T) {
		expected := "\n@testMmUser approved your pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title\n"

		actual, err := tr.RenderPullRequestApprovedNotificationForPullRequestAuthor(getTestPullRequestApprovedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestApprovedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) " +
			"was approved by @testMmUser\n"

		actual, err := tr.RenderPullRequestApprovedEventNotificationForSubscribedChannels(getTestPullRequestApprovedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestCommentNotificationForPullRequestAuthor", func(t *testing.T) {
		expected := "\n@testMmUser commented on your pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - " +
			"Test title\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderPullRequestCommentNotificationForPullRequestAuthor(getTestPullRequestCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"New comment by @testMmUser on " +
			"[#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1):" +
			"\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels(getTestPullRequestCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestCommentMentionNotification", func(t *testing.T) {
		expected := "\n@testMmUser " +
			"mentioned you on [mattermost-plugin-bitbucket#1](https://bitbucket.org/test-comment-link/) - Test title:" +
			"\n>this issue should be fixed by @testMmUser\n"

		actual, err := tr.RenderPullRequestCommentMentionNotification(getTestPullRequestCommentCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestDescriptionMentionNotification", func(t *testing.T) {
		expected := "\n@testMmUser mentioned you in pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1)" +
			" - Test title:\n>Test description @testMmUser\n"

		actual, err := tr.RenderPullRequestDescriptionMentionNotification(getTestPullRequestCreatedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestMergedEventNotificationForPullRequestAuthor", func(t *testing.T) {
		expected := "\n@testMmUser merged your pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title\n"

		pl := getTestPullRequestMergedPayload()

		actual, err := tr.RenderPullRequestMergedEventNotificationForPullRequestAuthor(pl)

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestMergedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) " +
			"was merged by @testMmUser\n"

		actual, err := tr.RenderPullRequestMergedEventNotificationForSubscribedChannels(getTestPullRequestMergedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestUnapprovedEventNotificationForSubscribedChannels", func(t *testing.T) {
		expected := "\n[\\[mattermost-plugin-bitbucket\\]](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket) " +
			"Pull request [#1 Test title](https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) " +
			"was unapproved by @testMmUser\n"

		actual, err := tr.RenderPullRequestUnapprovedEventNotificationForSubscribedChannels(getTestPullRequestUnapprovedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("RenderPullRequestUnapprovedNotificationForPullRequestAuthor", func(t *testing.T) {
		expected := "\n@testMmUser unapproved your pull request [mattermost-plugin-bitbucket#1]" +
			"(https://bitbucket.org/mattermost/mattermost-plugin-bitbucket/pull-requests/1) - Test title\n"

		actual, err := tr.RenderPullRequestUnapprovedNotificationForPullRequestAuthor(getTestPullRequestUnapprovedPayload())

		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}
