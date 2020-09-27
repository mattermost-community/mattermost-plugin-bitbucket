package webhook

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/subscription"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/template_renderer"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
	"strings"
)

const (
	TemplateErrorText = "failed to render template"
)

type HandleWebhook struct {
	Message          string
	ToBitbucketUsers []string
	ToChannels       []string
}

type SubscriptionHandler interface {
	GetSubscribedChannelsForRepository(webhook_payload.Payload) []*subscription.Subscription
}

type PullRequestReviewHandler interface {
	GetAlreadyNotifiedUsers(pullRequestID int64) ([]string, error)
	SaveNotifiedUsers(int64, []string)
}

type Webhook interface {
	HandleRepoPushEvent(webhook_payload.RepoPushPayload) ([]*HandleWebhook, error)
	HandleIssueCreatedEvent(webhook_payload.IssueCreatedPayload) ([]*HandleWebhook, error)
	HandleIssueUpdatedEvent(webhook_payload.IssueUpdatedPayload) ([]*HandleWebhook, error)
	HandleIssueCommentCreatedEvent(webhook_payload.IssueCommentCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestCreatedEvent(webhook_payload.PullRequestCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestApprovedEvent(webhook_payload.PullRequestApprovedPayload) ([]*HandleWebhook, error)
	HandlePullRequestDeclinedEvent(webhook_payload.PullRequestDeclinedPayload) ([]*HandleWebhook, error)
	HandlePullRequestUnapprovedEvent(webhook_payload.PullRequestUnapprovedPayload) ([]*HandleWebhook, error)
	HandlePullRequestMergedEvent(webhook_payload.PullRequestMergedPayload) ([]*HandleWebhook, error)
	HandlePullRequestCommentCreatedEvent(webhook_payload.PullRequestCommentCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestUpdatedEvent(webhook_payload.PullRequestUpdatedPayload) ([]*HandleWebhook, error)
}

type webhook struct {
	subscriptionConfiguration SubscriptionHandler
	reviewConfiguration       PullRequestReviewHandler
	templateRenderer          template_renderer.TemplateRenderer
}

func NewWebhook(s SubscriptionHandler, r PullRequestReviewHandler, t template_renderer.TemplateRenderer) Webhook {
	return &webhook{subscriptionConfiguration: s, reviewConfiguration: r, templateRenderer: t}
}

func (w *webhook) createPrivateMessageHandleWebhook(pl webhook_payload.Payload, message string, accountIDs []string) *HandleWebhook {
	handler := &HandleWebhook{Message: message}

	for _, accountID := range accountIDs {
		if accountID == pl.GetActor().AccountId {
			continue
		}

		handler.ToBitbucketUsers = append(handler.ToBitbucketUsers, accountID)
	}

	return handler

}

func (w *webhook) parseBitbucketAcountIDsFromHTML(html string) []string {
	accountIdMap := map[string]bool{}
	var accountIds []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return accountIds
	}

	//looking for span tags in the HTML with user account IDs
	doc.Find("span[class=\"ap-mention\"]").Each(func(i int, selection *goquery.Selection) {
		bitbucketUserAccountID := selection.AttrOr("data-atlassian-id", "")
		//put the found accountID in the map if it doesn't exist there yet
		if bitbucketUserAccountID != "" && !accountIdMap[bitbucketUserAccountID] {
			accountIds = append(accountIds, bitbucketUserAccountID)
			accountIdMap[bitbucketUserAccountID] = true
		}
	})

	return accountIds
}

func cleanWebhookHandlers(handlers []*HandleWebhook) []*HandleWebhook {
	res := make([]*HandleWebhook, 0)
	for _, handler := range handlers {
		//don't pass nil handlers
		if handler == nil {
			continue
		}
		// don't send handlers with empty messages
		if handler.Message == "" {
			continue
		}
		res = append(res, handler)
	}
	return res
}
