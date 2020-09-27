package template_renderer

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	"github.com/PuerkitoBio/goquery"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

type BitBucketAccountIDToUsernameMappingCallbackType func(string) string

type TemplateRenderer interface {
	RegisterBitBucketAccountIDToUsernameMappingCallback(callback BitBucketAccountIDToUsernameMappingCallbackType)
	RenderBranchOrTagCreatedEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (string, error)
	RenderBranchOrTagDeletedEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (string, error)
	RenderIssueCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCreatedPayload) (string, error)
	RenderIssueUpdatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueUpdatedPayload) (string, error)
	RenderIssueAssignmentNotificationForAssignedUser(pl webhook_payload.IssueUpdatedPayload) (string, error)
	RenderIssueStatusUpdateNotificationForIssueReporter(pl webhook_payload.IssueUpdatedPayload) (string, error)
	RenderIssueDescriptionMentionNotification(pl webhook_payload.IssueCreatedPayload) (string, error)
	RenderIssueCommentCreatedEventNotificationForSubscribedChannels(pl webhook_payload.IssueCommentCreatedPayload) (string, error)
	RenderIssueCommentNotificationForIssueReporter(pl webhook_payload.IssueCommentCreatedPayload) (string, error)
	RenderIssueCommentMentionNotification(pl webhook_payload.IssueCommentCreatedPayload) (string, error)
	RenderPullRequestCreatedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestCreatedPayload) (string, error)
	RenderPullRequestDeclinedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestDeclinedPayload) (string, error)
	RenderPullRequestDeclinedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestDeclinedPayload) (string, error)
	RenderPullRequestApprovedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestApprovedPayload) (string, error)
	RenderPullRequestApprovedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestApprovedPayload) (string, error)
	RenderPullRequestAssignedNotification(pl webhook_payload.PullRequestUpdatedPayload) (string, error)
	RenderPullRequestCommentNotificationForPullRequestAuthor(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error)
	RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error)
	RenderPullRequestCommentMentionNotification(pl webhook_payload.PullRequestCommentCreatedPayload) (string, error)
	RenderPullRequestDescriptionMentionNotification(pl webhook_payload.PullRequestCreatedPayload) (string, error)
	RenderPullRequestMergedEventNotificationForPullRequestAuthor(pl webhook_payload.PullRequestMergedPayload) (string, error)
	RenderPullRequestMergedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestMergedPayload) (string, error)
	RenderPullRequestUnapprovedEventNotificationForSubscribedChannels(pl webhook_payload.PullRequestUnapprovedPayload) (string, error)
	RenderPullRequestUnapprovedNotificationForPullRequestAuthor(pl webhook_payload.PullRequestUnapprovedPayload) (string, error)
	RenderRepoPushEventNotificationForSubscribedChannels(pl webhook_payload.RepoPushPayload) (string, error)
}

type templateRenderer struct {
	masterTemplate                              *template.Template
	bitBucketAccountIDToUsernameMappingCallback BitBucketAccountIDToUsernameMappingCallbackType
}

func MakeTemplateRenderer() TemplateRenderer {
	tr := templateRenderer{}
	tr.init()
	return &tr
}

func (tr *templateRenderer) renderTemplate(payload interface{}, templateName string, text string) (string, error) {
	//checks whether a template with this name is already defined
	t := tr.masterTemplate.Lookup(templateName)
	if t == nil {
		//if the template is not defined, it will be defined now
		t = template.Must(tr.masterTemplate.New(templateName).Parse(text))
	}

	var output bytes.Buffer

	err := t.Execute(&output, payload)
	if err != nil {
		return "", errors.Wrapf(err, "Could not execute template named %s", t.Name())
	}

	return output.String(), nil
}

func (tr *templateRenderer) init() {
	var funcMap = sprig.TxtFuncMap()
	// Quote the body
	funcMap["quote"] = func(body string) string {
		return ">" + strings.ReplaceAll(body, "\n", "\n>")
	}

	// Resolve a BitBucket username to the corresponding Mattermost username, if linked.
	funcMap["lookupMattermostUsername"] = tr.lookupMattermostUsername

	// Remove \n
	funcMap["removeLineBreaks"] = func(body string) string {
		return strings.ReplaceAll(body, "\n", "")
	}

	// Replace any BitBucket username with its corresponding Mattermost username, if any
	funcMap["replaceAllBitBucketUsernames"] = func(body string) string {

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
		if err != nil {
			return body
		}

		doc.Find("span[class=\"ap-mention\"]").Each(func(i int, selection *goquery.Selection) {
			bitbucketNickname := selection.Text()
			bitbucketAcountID := selection.AttrOr("data-atlassian-id", "")

			var mattermostUsername string
			if bitbucketAcountID != "" {
				mattermostUsername = tr.lookupMattermostUsername(bitbucketAcountID)
				if mattermostUsername != "" {
					bitbucketNickname = "@" + mattermostUsername
				}
			}

			selection.SetText(bitbucketNickname)
		})

		//Text() returns only text, without HTML attributes or tags
		return doc.Text()
	}

	tr.masterTemplate = template.Must(template.New("master").Funcs(funcMap).Parse(""))

	// The repo template links to the corresponding repository.
	template.Must(tr.masterTemplate.New("repo").Parse(`[\[{{.FullName}}\]]({{.Links.HTML.Href}})`))

	// The repoPullRequestWithTitle links to the corresponding pull request.
	template.Must(tr.masterTemplate.New("repoPullRequestWithTitle").Parse(
		`[{{.Repository.FullName}}#{{.PullRequest.ID}}]({{.PullRequest.Links.HTML.Href}}) - {{.PullRequest.Title}}`,
	))

	// The pullRequest links to the corresponding pull request, skipping the repo title.
	template.Must(tr.masterTemplate.New("pullRequest").Parse(
		`[#{{.ID}} {{.Title}}]({{.Links.HTML.Href}})`,
	))

	// The issue links to the corresponding issue.
	template.Must(tr.masterTemplate.New("issue").Parse(
		`[\[{{.Repository.FullName}}#{{.Issue.ID}}\]]({{.Issue.Links.HTML.Href}})`,
	))

	//The user template links to the corresponding user in Mattermost or in BitBucket.
	template.Must(tr.masterTemplate.New("user").Parse(`
{{- $mattermostUsername := .AccountId | lookupMattermostUsername}}
{{- if $mattermostUsername }}@{{$mattermostUsername}}
{{- else}}[{{.NickName}}]({{.Links.HTML.Href}})
{{- end -}}
`))
}

func (tr *templateRenderer) renderTemplateWithName(name string, data interface{}) (string, error) {
	var output bytes.Buffer
	t := tr.masterTemplate.Lookup(name)
	if t == nil {
		return "", errors.Errorf("no template named %s", name)
	}

	err := t.Execute(&output, data)
	if err != nil {
		return "", errors.Wrapf(err, "Could not execute template named %s", name)
	}

	return output.String(), nil
}

func (tr *templateRenderer) RegisterBitBucketAccountIDToUsernameMappingCallback(callback BitBucketAccountIDToUsernameMappingCallbackType) {
	tr.bitBucketAccountIDToUsernameMappingCallback = callback
}

func (tr *templateRenderer) lookupMattermostUsername(bitbucketAccountID string) string {
	if tr.bitBucketAccountIDToUsernameMappingCallback == nil {
		return ""
	}

	return tr.bitBucketAccountIDToUsernameMappingCallback(bitbucketAccountID)
}
