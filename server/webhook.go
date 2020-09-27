package main

import (
	"context"
	"encoding/json"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/subscription"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
	"github.com/mattermost/mattermost-server/v5/model"
	"net/http"
	"strconv"
)

const (
	KeyAssignUserPr          = "pr_assigned_"
	BitbucketWebhookPostType = "custom_bb_webhook"
)

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {

	hook, _ := webhook_payload.New()
	payload, err := hook.Parse(r,
		webhook_payload.RepoPushEvent,
		webhook_payload.IssueCreatedEvent,
		webhook_payload.IssueUpdatedEvent,
		webhook_payload.IssueCommentCreatedEvent,
		webhook_payload.PullRequestCreatedEvent,
		webhook_payload.PullRequestUpdatedEvent,
		webhook_payload.PullRequestApprovedEvent,
		webhook_payload.PullRequestUnapprovedEvent,
		webhook_payload.PullRequestDeclinedEvent,
		webhook_payload.PullRequestMergedEvent,
		webhook_payload.PullRequestCommentCreatedEvent)

	if err != nil {
		p.API.LogError(err.Error())
		return
	}

	var handlers []*webhook.HandleWebhook
	var handlerError error

	switch payload.(type) {
	case webhook_payload.RepoPushPayload:
		handlers, handlerError = p.webhookHandler.HandleRepoPushEvent(payload.(webhook_payload.RepoPushPayload))
	case webhook_payload.IssueCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueCreatedEvent(payload.(webhook_payload.IssueCreatedPayload))
	case webhook_payload.IssueUpdatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueUpdatedEvent(payload.(webhook_payload.IssueUpdatedPayload))
	case webhook_payload.IssueCommentCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueCommentCreatedEvent(payload.(webhook_payload.IssueCommentCreatedPayload))
	case webhook_payload.PullRequestCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestCreatedEvent(payload.(webhook_payload.PullRequestCreatedPayload))
	case webhook_payload.PullRequestUpdatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestUpdatedEvent(payload.(webhook_payload.PullRequestUpdatedPayload))
	case webhook_payload.PullRequestApprovedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestApprovedEvent(payload.(webhook_payload.PullRequestApprovedPayload))
	case webhook_payload.PullRequestCommentCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestCommentCreatedEvent(payload.(webhook_payload.PullRequestCommentCreatedPayload))
	case webhook_payload.PullRequestDeclinedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestDeclinedEvent(payload.(webhook_payload.PullRequestDeclinedPayload))
	case webhook_payload.PullRequestUnapprovedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestUnapprovedEvent(payload.(webhook_payload.PullRequestUnapprovedPayload))
	case webhook_payload.PullRequestMergedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestMergedEvent(payload.(webhook_payload.PullRequestMergedPayload))
	}

	if handlerError != nil {
		p.API.LogError(handlerError.Error())
		return
	}

	p.executeHandlers(handlers, payload.(webhook_payload.Payload))
}

func (p *Plugin) executeHandlers(webhookHandlers []*webhook.HandleWebhook, pl webhook_payload.Payload) {
	for _, webhookHandler := range webhookHandlers {

		post := &model.Post{
			UserId:  p.BotUserID,
			Message: webhookHandler.Message,
			Type:    BitbucketWebhookPostType,
		}

		for _, channelID := range webhookHandler.ToChannels {
			post.ChannelId = channelID
			if _, err := p.API.CreatePost(post); err != nil {
				p.API.LogError(err.Error())
			}
		}

		for _, toBitbucketUser := range webhookHandler.ToBitbucketUsers {
			userID := p.getBitbucketAccountIDToMattermostUserIDMapping(toBitbucketUser)
			if userID == "" {
				continue
			}

			if pl.GetRepository().IsPrivate && !p.permissionToRepo(userID, pl.GetRepository().FullName) {
				continue
			}

			channel, err := p.API.GetDirectChannel(userID, p.BotUserID)
			if err != nil {
				continue
			}

			post.ChannelId = channel.Id
			_, err = p.API.CreatePost(post)
			if err != nil {
				p.API.LogError(err.Error())
			}
			p.sendRefreshEvent(userID)
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
		p.API.LogError("Couldn't fetch repositories info", "err", err)
		return false
	}
	return true
}

type pullRequestReviewHandler struct {
	p *Plugin
}

type pullRequestReviewers struct {
	Users []string
}

type subscriptionHandler struct {
	p *Plugin
}

func (s subscriptionHandler) GetSubscribedChannelsForRepository(pl webhook_payload.Payload) []*subscription.Subscription {
	return s.p.GetSubscribedChannelsForRepository(pl)
}

func (r *pullRequestReviewHandler) GetAlreadyNotifiedUsers(pullRequestID int64) ([]string, error) {
	bytesThisPrReviewers, err := r.p.API.KVGet(KeyAssignUserPr + strconv.FormatInt(pullRequestID, 10))
	if err != nil {
		return nil, err
	}

	//if nil, then return empty list
	if bytesThisPrReviewers == nil {
		return []string{}, nil
	}

	var pullRequestReviewers pullRequestReviewers
	appErr := json.Unmarshal(bytesThisPrReviewers, &pullRequestReviewers)
	if appErr != nil {
		r.p.API.LogError("Couldn't read information about notified users",
			"pl.PullRequest.ID", pullRequestID, "err", appErr)
		return nil, appErr
	}

	return pullRequestReviewers.Users, nil
}

func (r pullRequestReviewHandler) SaveNotifiedUsers(pullRequestID int64, notifiedUsers []string) {
	thisPrReviewers := pullRequestReviewers{}
	thisPrReviewers.Users = notifiedUsers
	bytesThisPrReviewers, err := json.Marshal(thisPrReviewers)
	if err != nil {
		r.p.API.LogWarn("Couldn't marshal notified users for PR",
			"thisPrReviewers", thisPrReviewers, "err", err)
		return
	}

	err = r.p.API.KVSet(KeyAssignUserPr+strconv.FormatInt(pullRequestID, 10), bytesThisPrReviewers)
	if err != nil {
		//err is nil, but it's still going here don't know why todo
		r.p.API.LogWarn("Couldn't save information about notified users for PR",
			"thisPrReviewers", thisPrReviewers, "err", err)
	}
}
