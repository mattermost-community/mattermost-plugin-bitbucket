package template_renderer

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
)

func (tr *templateRenderer) RenderIssueCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCreatedEventNotificationForSubscribedChannels", `
#### {{.Issue.Title}}
##### {{template "issue" .}}
#new-issue by {{template "user" .Actor}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueUpdatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueUpdatedEventNotificationForSubscribedChannels", `
#### {{.Issue.Title}}
##### {{template "issue" .}}
#updated-issue by {{template "user" .Actor}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueAssignmentNotificationForAssignedUser(pl webhook_payload.IssueUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueAssignmentNotificationForAssignedUser", `
{{template "user" .Actor}} assigned you to issue {{template "issue" .}}
`)
}

func (tr *templateRenderer) RenderIssueStatusUpdateNotificationForIssueReporter(pl webhook_payload.IssueUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueStatusUpdateNotificationForIssueReporter", `
{{template "user" .Actor}} set status to {{.Changes.Status.New}} of your issue {{template "issue" .}}
`)
}

func (tr *templateRenderer) RenderIssueDescriptionMentionNotification(pl webhook_payload.IssueCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueDescriptionMentionNotification", `
{{template "user" .Actor}} mentioned you on {{template "issue" .}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} New comment by {{template "user" .Actor}} on {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentNotificationForIssueReporter(pl webhook_payload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentNotificationForIssueReporter", `
{{template "user" .Actor}} commented on your issue {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentMentionNotification(pl webhook_payload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentMentionNotification", `
{{template "user" .Actor}} mentioned you on {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}
