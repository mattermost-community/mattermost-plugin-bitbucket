package webhook

import (
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhookpayload"

	"github.com/pkg/errors"
)

func (w *webhook) HandlePullRequestCreatedEvent(pl webhookpayload.PullRequestCreatedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestDescriptionMentionNotification(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandlePullRequestApprovedEvent(pl webhookpayload.PullRequestApprovedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestApprovedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestApprovedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandlePullRequestDeclinedEvent(pl webhookpayload.PullRequestDeclinedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestDeclinedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestDeclinedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandlePullRequestUnapprovedEvent(pl webhookpayload.PullRequestUnapprovedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestUnapprovedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestUnapprovedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandlePullRequestMergedEvent(pl webhookpayload.PullRequestMergedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestMergedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestMergedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2)), nil
}

func (w *webhook) HandlePullRequestCommentCreatedEvent(pl webhookpayload.PullRequestCommentCreatedPayload) ([]*HandleWebhook, error) {
	var handlers []*HandleWebhook

	handler1, err := w.createPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler2, err := w.createPullRequestCommentMentionNotification(pl)
	if err != nil {
		return nil, err
	}

	handler3, err := w.createPullRequestCommentNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, err
	}

	return cleanWebhookHandlers(append(handlers, handler1, handler2, handler3)), nil
}

func (w *webhook) HandlePullRequestUpdatedEvent(pl webhookpayload.PullRequestUpdatedPayload) ([]*HandleWebhook, error) {
	// ignore if there are no reviewers
	if len(pl.PullRequest.Reviewers) == 0 {
		return []*HandleWebhook{}, nil
	}

	thisPullRequestReviewers, err := w.reviewConfiguration.GetAlreadyNotifiedUsers(pl.PullRequest.ID)
	if err != nil {
		return nil, err
	}

	message, templateErr := w.templateRenderer.RenderPullRequestAssignedNotification(pl)
	if templateErr != nil {
		return nil, templateErr
	}

	handler := &HandleWebhook{Message: message}

	// if reviewers are not empty, send them notifications
	for _, reviewer := range pl.PullRequest.Reviewers {
		// check if the user had been already notified
		if contains(thisPullRequestReviewers, reviewer.AccountID) {
			continue
		}

		thisPullRequestReviewers = append(thisPullRequestReviewers, reviewer.AccountID)
		handler.ToBitbucketUsers = append(handler.ToBitbucketUsers, reviewer.AccountID)
	}

	// save information about users that had been notified
	w.reviewConfiguration.SaveNotifiedUsers(pl.PullRequest.ID, thisPullRequestReviewers)

	return cleanWebhookHandlers([]*HandleWebhook{handler}), nil
}

func (w *webhook) createPullRequestCreatedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.Pulls() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestApprovedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestApprovedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestApprovedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestDeclinedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestDeclinedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestDeclinedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestUnapprovedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestUnapprovedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestUnapprovedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestMergedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestMergedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestMergedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl webhookpayload.PullRequestCommentCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestCommentCreatedEventNotificationForSubscribedChannels(pl)
	if err != nil {
		return nil, err
	}

	handler := &HandleWebhook{Message: message}

	subs := w.subscriptionConfiguration.GetSubscribedChannelsForRepository(&pl)
	if len(subs) == 0 {
		return handler, nil
	}

	for _, sub := range subs {
		if !sub.PullReviews() {
			continue
		}
		handler.ToChannels = append(handler.ToChannels, sub.ChannelID)
	}

	return handler, nil
}

func (w *webhook) createPullRequestDescriptionMentionNotification(pl webhookpayload.PullRequestCreatedPayload) (*HandleWebhook, error) {
	mentionedAccountIDs := w.parseBitbucketAcountIDsFromHTML(pl.PullRequest.Rendered.Description.HTML)
	message, err := w.templateRenderer.RenderPullRequestDescriptionMentionNotification(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, mentionedAccountIDs), nil
}

func (w *webhook) createPullRequestCommentMentionNotification(pl webhookpayload.PullRequestCommentCreatedPayload) (*HandleWebhook, error) {
	mentionedAccountIDs := w.parseBitbucketAcountIDsFromHTML(pl.Comment.Content.HTML)
	message, err := w.templateRenderer.RenderPullRequestCommentMentionNotification(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, mentionedAccountIDs), nil
}

func (w *webhook) createPullRequestCommentNotificationForPullRequestAuthor(pl webhookpayload.PullRequestCommentCreatedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestCommentNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.PullRequest.Author.AccountID}), nil
}

func (w *webhook) createPullRequestApprovedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestApprovedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestApprovedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.PullRequest.Author.AccountID}), nil
}

func (w *webhook) createPullRequestDeclinedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestDeclinedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestDeclinedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.PullRequest.Author.AccountID}), nil
}

func (w *webhook) createPullRequestUnapprovedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestUnapprovedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestUnapprovedNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.PullRequest.Author.AccountID}), nil
}

func (w *webhook) createPullRequestMergedNotificationForPullRequestAuthor(pl webhookpayload.PullRequestMergedPayload) (*HandleWebhook, error) {
	message, err := w.templateRenderer.RenderPullRequestMergedEventNotificationForPullRequestAuthor(pl)
	if err != nil {
		return nil, errors.Wrap(err, TemplateErrorText)
	}

	return w.createPrivateMessageHandleWebhook(&pl, message, []string{pl.PullRequest.Author.AccountID}), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
