package webhook

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/subscription"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/templaterenderer"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"
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
	GetSubscribedChannelsForRepository(webhookpayload.Payload) []*subscription.Subscription
}

type PullRequestReviewHandler interface {
	GetAlreadyNotifiedUsers(pullRequestID int64) ([]string, error)
	SaveNotifiedUsers(int64, []string)
}

type Webhook interface {
	HandleRepoPushEvent(webhookpayload.RepoPushPayload) ([]*HandleWebhook, error)
	HandleIssueCreatedEvent(webhookpayload.IssueCreatedPayload) ([]*HandleWebhook, error)
	HandleIssueUpdatedEvent(webhookpayload.IssueUpdatedPayload) ([]*HandleWebhook, error)
	HandleIssueCommentCreatedEvent(webhookpayload.IssueCommentCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestCreatedEvent(webhookpayload.PullRequestCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestApprovedEvent(webhookpayload.PullRequestApprovedPayload) ([]*HandleWebhook, error)
	HandlePullRequestDeclinedEvent(webhookpayload.PullRequestDeclinedPayload) ([]*HandleWebhook, error)
	HandlePullRequestUnapprovedEvent(webhookpayload.PullRequestUnapprovedPayload) ([]*HandleWebhook, error)
	HandlePullRequestMergedEvent(webhookpayload.PullRequestMergedPayload) ([]*HandleWebhook, error)
	HandlePullRequestCommentCreatedEvent(webhookpayload.PullRequestCommentCreatedPayload) ([]*HandleWebhook, error)
	HandlePullRequestUpdatedEvent(webhookpayload.PullRequestUpdatedPayload) ([]*HandleWebhook, error)
}

type webhook struct {
	subscriptionConfiguration SubscriptionHandler
	reviewConfiguration       PullRequestReviewHandler
	templateRenderer          templaterenderer.TemplateRenderer
}

func NewWebhook(s SubscriptionHandler, r PullRequestReviewHandler, t templaterenderer.TemplateRenderer) Webhook {
	return &webhook{subscriptionConfiguration: s, reviewConfiguration: r, templateRenderer: t}
}

func (w *webhook) createPrivateMessageHandleWebhook(pl webhookpayload.Payload, message string, accountIDs []string) *HandleWebhook {
	handler := &HandleWebhook{Message: message}

	for _, accountID := range accountIDs {
		if accountID == pl.GetActor().AccountID {
			continue
		}

		handler.ToBitbucketUsers = append(handler.ToBitbucketUsers, accountID)
	}

	return handler
}

func (w *webhook) parseBitbucketAcountIDsFromHTML(html string) []string {
	accountIDMap := map[string]bool{}
	var accountIds []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return accountIds
	}

	//looking for span tags in the HTML with user account IDs
	doc.Find("span[class=\"ap-mention\"]").Each(func(i int, selection *goquery.Selection) {
		bitbucketUserAccountID := selection.AttrOr("data-atlassian-id", "")
		//put the found accountID in the map if it doesn't exist there yet
		if bitbucketUserAccountID != "" && !accountIDMap[bitbucketUserAccountID] {
			accountIds = append(accountIds, bitbucketUserAccountID)
			accountIDMap[bitbucketUserAccountID] = true
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
