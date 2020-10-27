package templaterenderer

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"
)

func (tr *templateRenderer) RenderPullRequestCreatedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was created by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestDeclinedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestDeclinedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDeclinedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was declined by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestDeclinedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestDeclinedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDeclinedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} declined your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestApprovedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestApprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestApprovedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} approved your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestApprovedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestApprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestApprovedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was approved by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestAssignedNotification(pl webhookpayload.PullRequestUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestAssignedNotification", `
{{template "user" .Actor}} assigned you to pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentNotificationForPullRequestAuthor(pl webhookpayload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentNotificationForPullRequestAuthor", `
{{template "user" .Actor}} commented on your pull request {{template "repoPullRequestWithTitle" .}}
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} New comment by {{template "user" .Actor}} on {{template "pullRequest" .PullRequest}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentMentionNotification(pl webhookpayload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentMentionNotification", `
{{template "user" .Actor}} mentioned you on [{{.Repository.FullName}}#{{.PullRequest.ID}}]({{.Comment.Links.HTML.Href}}) - {{.PullRequest.Title}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestDescriptionMentionNotification(pl webhookpayload.PullRequestCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDescriptionMentionNotification", `
{{template "user" .Actor}} mentioned you in pull request {{template "repoPullRequestWithTitle" .}}:
{{.PullRequest.Rendered.Description.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestMergedEventNotificationForPullRequestAuthor(pl webhookpayload.PullRequestMergedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestMergedEventNotificationForPullRequestAuthor", `
{{template "user" .Actor}} merged your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestMergedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestMergedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestMergedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was merged by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestUnapprovedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestUnapprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestUnapprovedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was unapproved by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestUnapprovedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestUnapprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestUnapprovedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} unapproved your pull request {{template "repoPullRequestWithTitle" .}}
`)
}
