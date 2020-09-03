package main

import (
	"context"
	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/mattermost-server/v5/model"
	bb_webhook "gopkg.in/go-playground/webhooks.v5/bitbucket"
	"net/http"
)

const (
	PostTypeIssueCreated     = "custom_bb_issue_create"
	PostTypePushed           = "custom_bb_push"
	PostTypeIssueUpdated     = "custom_bb_issue_update"
	PostTypePrCommentCreated = "custom_bb_pr_cc"
	PostTypePrMerged         = "custom_bb_pr_merged"
	PostTypePrApproved       = "custom_bb_pr_approved"
	PostTypePrCreated        = "custom_bb_pr_created"
	PostTypeCommentCreated   = "custom_bb_issue_cc"

	TemplateErrorText = "failed to render template"
)

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {

	hook, _ := bb_webhook.New()
	payload, err := hook.Parse(r,
		bb_webhook.RepoPushEvent,
		bb_webhook.IssueCreatedEvent,
		bb_webhook.IssueUpdatedEvent,
		bb_webhook.IssueCommentCreatedEvent,
		bb_webhook.PullRequestCreatedEvent,
		bb_webhook.PullRequestApprovedEvent,
		bb_webhook.PullRequestMergedEvent,
		bb_webhook.PullRequestCommentCreatedEvent)

	var handler func()

	switch payload.(type) {
	case bb_webhook.RepoPushPayload:
		handler = func() {
			p.repoPushEvent(payload)
			return
		}
	case bb_webhook.IssueUpdatedPayload:
		handler = func() {
			p.postIssueUpdatedEvent(payload)
			return
		}
	case bb_webhook.IssueCreatedPayload:
		handler = func() {
			p.postIssueCreatedEvent(payload)
			return
		}
	case bb_webhook.IssueCommentCreatedPayload:
		handler = func() {
			p.postIssueCommentCreatedEvent(payload)
			return
		}
	case bb_webhook.PullRequestCreatedPayload:
		handler = func() {
			p.postPullRequestCreatedEvent(payload)
			return
		}
	case bb_webhook.PullRequestCommentCreatedPayload:
		handler = func() {
			p.postPullRequestCommentCreatedEvent(payload)
			return
		}
	case bb_webhook.PullRequestApprovedPayload:
		handler = func() {
			p.postPullRequestApprovedEvent(payload)
			return
		}
	case bb_webhook.PullRequestMergedPayload:
		handler = func() {
			p.postPullRequestMergedEvent(payload)
			return
		}
	}

	if err != nil {
		p.API.LogError(err.Error())
		return
	}

	handler()
}

func (p *Plugin) repoPushEvent(pl interface{}) {
	r := pl.(bb_webhook.RepoPushPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("pushed", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypePushed,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.Pushes() {
			continue
		}

		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			mlog.Error(err.Error())
		}
	}
}

func (p *Plugin) postIssueUpdatedEvent(pl interface{}) {
	r := pl.(bb_webhook.IssueUpdatedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("issueUpdated", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypeIssueUpdated,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.Issues() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) postIssueCreatedEvent(pl interface{}) {
	r := pl.(bb_webhook.IssueCreatedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("issueCreated", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypeIssueCreated,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.Issues() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) postIssueCommentCreatedEvent(pl interface{}) {
	r := pl.(bb_webhook.IssueCommentCreatedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)

	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("issueCommentCreated", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypeCommentCreated,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.IssueComments() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}

}

func (p *Plugin) postPullRequestCreatedEvent(pl interface{}) {
	r := pl.(bb_webhook.PullRequestCreatedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("prCreated", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypePrCreated,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.Pulls() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) postPullRequestApprovedEvent(pl interface{}) {
	r := pl.(bb_webhook.PullRequestApprovedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("prApproved", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypePrApproved,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) postPullRequestMergedEvent(pl interface{}) {
	r := pl.(bb_webhook.PullRequestMergedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("prMerged", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypePrMerged,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) postPullRequestCommentCreatedEvent(pl interface{}) {
	r := pl.(bb_webhook.PullRequestCommentCreatedPayload)

	reponame := r.Repository.FullName
	isprivate := r.Repository.IsPrivate

	subs := p.GetSubscribedChannelsForRepository(reponame, isprivate)
	if subs == nil || len(subs) == 0 {
		return
	}

	message, err := renderTemplate("prCommentCreated", r)
	if err != nil {
		mlog.Error(TemplateErrorText, mlog.Err(err))
		return
	}

	post := &model.Post{
		Type:    PostTypePrCommentCreated,
		UserId:  p.BotUserID,
		Message: message,
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		post.ChannelId = sub.ChannelID
		if _, err := p.API.CreatePost(post); err != nil {
			p.API.LogError(err.Error())
		}
	}
}

func (p *Plugin) permissionToRepo(userID string, ownerAndRepo string) bool {
	_, owner, repo := parseOwnerAndRepo(ownerAndRepo, p.getBaseURL())

	if owner == "" {
		return false
	}
	if err := p.checkOrg(owner); err != nil {
		return false
	}

	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		return false
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	if _, _, err := bitbucketClient.RepositoriesApi.RepositoriesUsernameRepoSlugGet(context.Background(), owner, repo); err != nil {
		mlog.Error(err.Error())
		return false
	}
	return true
}
