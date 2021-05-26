package templaterenderer

import (
	"errors"

	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"
)

func (tr *templateRenderer) RenderIssueCreatedEventNotificationForSubscribedChannels(pl webhookpayload.Payload) (string, error) {
	typedPayload, ok := pl.(webhookpayload.IssueCreatedPayload)
	if !ok {
		return "", errors.New("invalid type")
	}
	return tr.renderTemplate(typedPayload, "issueCreatedEventNotificationForSubscribedChannels", `
#### {{.Issue.Title}}
##### {{template "issue" .}}
#new-issue by {{template "user" .Actor}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueUpdatedEventNotificationForSubscribedChannels(pl webhookpayload.Payload) (string, error) {
	typedPayload, ok := pl.(webhookpayload.IssueUpdatedPayload)
	if !ok {
		return "", errors.New("invalid type")
	}
	return tr.renderTemplate(typedPayload, "issueUpdatedEventNotificationForSubscribedChannels", `
#### {{.Issue.Title}}
##### {{template "issue" .}}
#updated-issue by {{template "user" .Actor}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueAssignmentNotificationForAssignedUser(pl webhookpayload.IssueUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueAssignmentNotificationForAssignedUser", `
{{template "user" .Actor}} assigned you to issue {{template "issue" .}}
`)
}

func (tr *templateRenderer) RenderIssueStatusUpdateNotificationForIssueReporter(pl webhookpayload.IssueUpdatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueStatusUpdateNotificationForIssueReporter", `
{{template "user" .Actor}} set status to {{.Changes.Status.New}} of your issue {{template "issue" .}}
`)
}

func (tr *templateRenderer) RenderIssueDescriptionMentionNotification(pl webhookpayload.IssueCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueDescriptionMentionNotification", `
{{template "user" .Actor}} mentioned you on {{template "issue" .}}:
{{.Issue.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentCreatedEventNotificationForSubscribedChannels(pl webhookpayload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} New comment by {{template "user" .Actor}} on {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentNotificationForIssueReporter(pl webhookpayload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentNotificationForIssueReporter", `
{{template "user" .Actor}} commented on your issue {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}

func (tr *templateRenderer) RenderIssueCommentMentionNotification(pl webhookpayload.IssueCommentCreatedPayload) (string, error) {
	return tr.renderTemplate(pl, "issueCommentMentionNotification", `
{{template "user" .Actor}} mentioned you on {{template "issue" .}}:
{{.Comment.Content.HTML | replaceAllBitBucketUsernames | quote}}
`)
}
