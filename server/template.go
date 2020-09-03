package main

import (
	"bytes"
	"github.com/pkg/errors"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

var masterTemplate *template.Template

func init() {
	var funcMap = sprig.TxtFuncMap()

	masterTemplate = template.Must(template.New("master").Funcs(funcMap).Parse(""))

	//general
	template.Must(masterTemplate.New("user").Parse(`[{{.NickName}}]({{.Links.HTML.Href}})`))
	template.Must(masterTemplate.New("repo").Parse(`[\[{{.FullName}}\]]({{.Links.HTML.Href}})`))
	template.Must(masterTemplate.New("pr").Parse(`[#{{.ID}} {{.Title}}]({{.Links.HTML.Href}})`))
	template.Must(masterTemplate.New("repoIssue").Parse(`[\[{{.Repository.FullName}}#{{.Issue.ID}}\]]({{.Issue.Links.HTML.Href}})`))

	//issues
	template.Must(masterTemplate.New("issueCreated").Parse(`
#### {{.Issue.Title}}
##### {{template "repoIssue" .}}
#new-issue by {{template "user" .Actor}}

{{.Issue.Content.Raw}}
`))

	template.Must(masterTemplate.New("issueUpdated").Parse(`
#### {{.Issue.Title}}
##### {{template "repoIssue" .}}
#updated-issue by {{template "user" .Actor}}

{{.Issue.Content.Raw}}
`))

	template.Must(masterTemplate.New("issueCommentCreated").Parse(`
{{template "repo" .Repository}} New comment by {{template "user" .Actor}} on {{template "repoIssue" .}}:

{{.Comment.Content.Raw}}
`))

	//pull requests
	template.Must(masterTemplate.New("prCreated").Funcs(funcMap).Parse(`
{{template "repo" .Repository}} Pull request {{template "pr" .PullRequest}} was created by {{template "user" .Actor}}.
`))

	template.Must(masterTemplate.New("prApproved").Funcs(funcMap).Parse(`
{{template "repo" .Repository}} Pull request {{template "pr" .PullRequest}} was approved by {{template "user" .Actor}}.
`))

	template.Must(masterTemplate.New("prMerged").Funcs(funcMap).Parse(`
{{template "repo" .Repository}} Pull request {{template "pr" .PullRequest}} was merged by {{template "user" .Actor}}.
`))

	template.Must(masterTemplate.New("prCommentCreated").Funcs(funcMap).Parse(`
{{template "repo" .Repository}} New review comment by {{template "user" .Actor}} on {{template "pr" .PullRequest}}:

{{.Comment.Content.Raw}}
`))

	//pushes
	template.Must(masterTemplate.New("pushed").Funcs(funcMap).Parse(`
User {{template "user" .Actor}} {{if (index .Push.Changes 0).Forced}}force-{{end}}pushed ` +
		`[{{len (index .Push.Changes 0).Commits}} new commit{{if ne (len (index .Push.Changes 0).Commits) 1}}s{{end}}]` +
		`({{(index .Push.Changes 0).Links.HTML.Href}}):
{{range (index .Push.Changes 0).Commits -}} 
[\[{{.Hash | substr 0 6}}\]]({{.Links.HTML.Href}}) {{.Message}}
{{end -}}
`))
}

func renderTemplate(name string, data interface{}) (string, error) {
	var output bytes.Buffer
	t := masterTemplate.Lookup(name)
	if t == nil {
		return "", errors.Errorf("no template named %s", name)
	}

	err := t.Execute(&output, data)
	if err != nil {
		return "", errors.Wrapf(err, "Could not execute template named %s", name)
	}

	return output.String(), nil
}
