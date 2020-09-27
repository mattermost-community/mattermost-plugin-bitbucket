package template_renderer

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
)

func (tr *templateRenderer) RenderPullRequestCreatedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was created by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestDeclinedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestDeclinedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDeclinedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was declined by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestDeclinedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestDeclinedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDeclinedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} declined your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestApprovedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestApprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestApprovedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} approved your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestApprovedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestApprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestApprovedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was approved by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestAssignedNotification(pl webhook_payload.PullRequestUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestAssignedNotification", `
{{template "user" .Actor}} assigned you to pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentNotificationForPullRequestAuthor(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentNotificationForPullRequestAuthor", `
{{template "user" .Actor}} commented on your pull request {{template "repoPullRequestWithTitle" .}}
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} New comment by {{template "user" .Actor}} on {{template "pullRequest" .PullRequest}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestCommentMentionNotification(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestCommentMentionNotification", `
{{template "user" .Actor}} mentioned you on [{{.Repository.FullName}}#{{.PullRequest.ID}}]({{.Comment.Links.HTML.Href}}) - {{.PullRequest.Title}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestDescriptionMentionNotification(pl webhook_payload.PullRequestCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestDescriptionMentionNotification", `
{{template "user" .Actor}} mentioned you in pull request {{template "repoPullRequestWithTitle" .}}:
{{.PullRequest.Rendered.Description.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderPullRequestMergedEventNotificationForPullRequestAuthor(pl webhook_payload.PullRequestMergedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestMergedEventNotificationForPullRequestAuthor", `
{{template "user" .Actor}} merged your pull request {{template "repoPullRequestWithTitle" .}}
`)
}

func (tr *templateRenderer) RenderPullRequestMergedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestMergedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestMergedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was merged by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestUnapprovedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestUnapprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestUnapprovedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} Pull request {{template "pullRequest" .PullRequest}} was unapproved by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderPullRequestUnapprovedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestUnapprovedPayload) (string, error) {
	return tr.renderTemplate(pl, "pullRequestUnapprovedNotificationForPullRequestAuthor", `
{{template "user" .Actor}} unapproved your pull request {{template "repoPullRequestWithTitle" .}}
`)
}
