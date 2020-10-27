package templaterenderer

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"
)

func (tr *templateRenderer) RenderRepoPushEventNotificationForSubscribedChannels(pl webhookpayload.RepoPushPayload) (string, error) {
	return tr.renderTemplate(pl, "repoPushEventNotificationForSubscribedChannels", `
User {{template "user" .Actor}} {{if (index .Push.Changes 0).Forced}}force-{{end}}pushed `+
		`[{{len (index .Push.Changes 0).Commits}} new commit{{if ne (len (index .Push.Changes 0).Commits) 1}}s{{end}}]`+
		`({{(index .Push.Changes 0).Links.HTML.Href}}) to [\[{{.Repository.FullName}}:{{(index .Push.Changes 0).New.Name}}\]]({{(index .Push.Changes 0).New.Links.HTML.Href}}):
{{range (index .Push.Changes 0).Commits -}} 
[\[{{.Hash | substr 0 6}}\]]({{.Links.HTML.Href}}) {{.Message | removeLineBreaks}} - {{template "user" .Author.User}}
{{end -}}
`)
}

func (tr *templateRenderer) RenderBranchOrTagCreatedEventNotificationForSubscribedChannels(pl webhookpayload.RepoPushPayload) (string, error) {
	return tr.renderTemplate(pl, "branchOrTagCreatedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} {{if eq (index .Push.Changes 0).New.Type "tag"}}Tag{{else}}Branch{{end}} [{{(index .Push.Changes 0).New.Name}}]({{(index .Push.Changes 0).New.Links.HTML.Href}}) was created by {{template "user" .Actor}}
`)
}

func (tr *templateRenderer) RenderBranchOrTagDeletedEventNotificationForSubscribedChannels(pl webhookpayload.RepoPushPayload) (string, error) {
	return tr.renderTemplate(pl, "dranchOrTagDeletedEventNotificationForSubscribedChannels", `
{{template "repo" .Repository}} {{if eq (index .Push.Changes 0).Old.Type "tag"}}Tag{{else}}Branch{{end}} [{{(index .Push.Changes 0).Old.Name}}]({{(index .Push.Changes 0).Old.Links.HTML.Href}}) was deleted by {{template "user" .Actor}}
`)
}
