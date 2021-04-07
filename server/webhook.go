package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/kosgrz/mattermost-plugin-bitbucket/server/subscription"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"
)

const (
	KeyAssignUserPr          = "pr_assigned_"
	BitbucketWebhookPostType = "custom_bb_webhook"
)

func (p *Plugin) handleWebhook(w http.ResponseWriter, r *http.Request) {
	hook, _ := webhookpayload.New()
	payload, err := hook.Parse(r,
		webhookpayload.RepoPushEvent,
		webhookpayload.IssueCreatedEvent,
		webhookpayload.IssueUpdatedEvent,
		webhookpayload.IssueCommentCreatedEvent,
		webhookpayload.PullRequestCreatedEvent,
		webhookpayload.PullRequestUpdatedEvent,
		webhookpayload.PullRequestApprovedEvent,
		webhookpayload.PullRequestUnapprovedEvent,
		webhookpayload.PullRequestDeclinedEvent,
		webhookpayload.PullRequestMergedEvent,
		webhookpayload.PullRequestCommentCreatedEvent)

	if err != nil {
		p.API.LogError(err.Error())
		return
	}

	var handlers []*webhook.HandleWebhook
	var handlerError error

	switch typedPayload := payload.(type) {
	case webhookpayload.RepoPushPayload:
		handlers, handlerError = p.webhookHandler.HandleRepoPushEvent(typedPayload)
	case webhookpayload.IssueCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueCreatedEvent(typedPayload)
	case webhookpayload.IssueUpdatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueUpdatedEvent(typedPayload)
	case webhookpayload.IssueCommentCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandleIssueCommentCreatedEvent(typedPayload)
	case webhookpayload.PullRequestCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestCreatedEvent(typedPayload)
	case webhookpayload.PullRequestUpdatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestUpdatedEvent(typedPayload)
	case webhookpayload.PullRequestApprovedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestApprovedEvent(typedPayload)
	case webhookpayload.PullRequestCommentCreatedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestCommentCreatedEvent(typedPayload)
	case webhookpayload.PullRequestDeclinedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestDeclinedEvent(typedPayload)
	case webhookpayload.PullRequestUnapprovedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestUnapprovedEvent(typedPayload)
	case webhookpayload.PullRequestMergedPayload:
		handlers, handlerError = p.webhookHandler.HandlePullRequestMergedEvent(typedPayload)
	}

	if handlerError != nil {
		p.API.LogError(handlerError.Error())
		return
	}

	p.executeHandlers(handlers, payload.(webhookpayload.Payload))
}

func (p *Plugin) executeHandlers(webhookHandlers []*webhook.HandleWebhook, pl webhookpayload.Payload) {
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

			userInfo, userInfoErr := p.getBitbucketUserInfo(userID)
			if userInfoErr != nil {
				continue
			}

			if !userInfo.Settings.Notifications {
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
	_, owner, repo := parseOwnerAndRepoAndReturnFullAlso(ownerAndRepo, p.getBaseURL())

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

	if _, httpResponse, err := bitbucketClient.RepositoriesApi.RepositoriesUsernameRepoSlugGet(context.Background(), owner, repo); err != nil {
		if httpResponse != nil {
			_ = httpResponse.Body.Close()
		}
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

func (s subscriptionHandler) GetSubscribedChannelsForRepository(pl webhookpayload.Payload) []*subscription.Subscription {
	return s.p.GetSubscribedChannelsForRepository(pl)
}

func (r *pullRequestReviewHandler) GetAlreadyNotifiedUsers(pullRequestID int64) ([]string, error) {
	bytesThisPrReviewers, err := r.p.API.KVGet(KeyAssignUserPr + strconv.FormatInt(pullRequestID, 10))
	if err != nil {
		return nil, err
	}

	// if nil, then return empty list
	if bytesThisPrReviewers == nil {
		return []string{}, nil
	}

	var prReviewers pullRequestReviewers
	appErr := json.Unmarshal(bytesThisPrReviewers, &prReviewers)
	if appErr != nil {
		r.p.API.LogError("Couldn't read information about notified users",
			"pl.PullRequest.ID", pullRequestID, "err", appErr)
		return nil, appErr
	}

	return prReviewers.Users, nil
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

	apiErr := r.p.API.KVSet(KeyAssignUserPr+strconv.FormatInt(pullRequestID, 10), bytesThisPrReviewers)
	if apiErr != nil {
		// err is nil, but it's still going here don't know why todo
		r.p.API.LogWarn("Couldn't save information about notified users for PR",
			"thisPrReviewers", thisPrReviewers, "apiErr", apiErr)
	}
}
