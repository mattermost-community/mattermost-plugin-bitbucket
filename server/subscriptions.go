package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/wbrefvem/go-bitbucket"

	"github.com/mattermost/mattermost-plugin-bitbucket/server/subscription"
	"github.com/mattermost/mattermost-plugin-bitbucket/server/webhookpayload"
)

const (
	SubscriptionsKey         = "subscriptions"
	UnsubscribedErrorMessage = "Unable to unsubscribe from %s as it is not currently part of a subscription in this channel."
)

func (p *Plugin) Subscribe(ctx context.Context, bitbucketClient *bitbucket.APIClient, userID, owner, repo, channelID, features string) error {
	if owner == "" {
		return errors.Errorf("invalid repository")
	}

	if err := p.checkOrg(owner); err != nil {
		return errors.Wrap(err, "organization not supported")
	}

	var err error

	if repo == "" {
		_, _, err = bitbucketClient.UsersApi.UserGet(ctx)
		if err != nil {
			p.API.LogError("Cannot fetch user", "err", err.Error())
			return errors.Errorf("Unknown organization %s", owner)
		}
	} else {
		_, _, err = bitbucketClient.RepositoriesApi.RepositoriesUsernameRepoSlugGet(context.Background(), owner, repo)
		if err != nil {
			p.API.LogError("Cannot fetch repository", "err", err.Error())
			return errors.Errorf("unknown repository %s", fullNameFromOwnerAndRepo(owner, repo))
		}
	}

	sub := &subscription.Subscription{
		ChannelID:  channelID,
		CreatorID:  userID,
		Features:   features,
		Repository: fullNameFromOwnerAndRepo(owner, repo),
	}

	if err := p.AddSubscription(fullNameFromOwnerAndRepo(owner, repo), sub); err != nil {
		return errors.Wrap(err, "could not add subscription")
	}

	return nil
}

func (p *Plugin) SubscribeOrg(ctx context.Context, bitbucketClient *bitbucket.APIClient, userID, org, channelID, features string) error {
	if org == "" {
		return errors.New("invalid organization")
	}

	return p.Subscribe(ctx, bitbucketClient, userID, org, "", channelID, features)
}

func (p *Plugin) GetSubscriptionsByChannel(channelID string) ([]*subscription.Subscription, error) {
	var filteredSubs []*subscription.Subscription
	subs, err := p.GetSubscriptions()
	if err != nil {
		return nil, errors.Wrap(err, "could not get subscriptions")
	}

	for repo, v := range subs.Repositories {
		for _, s := range v {
			if s.ChannelID == channelID {
				// this is needed to be backwards compatible
				if len(s.Repository) == 0 {
					s.Repository = repo
				}
				filteredSubs = append(filteredSubs, s)
			}
		}
	}

	sort.Slice(filteredSubs, func(i, j int) bool {
		return filteredSubs[i].Repository < filteredSubs[j].Repository
	})

	return filteredSubs, nil
}

func (p *Plugin) AddSubscription(repo string, sub *subscription.Subscription) error {
	subs, err := p.GetSubscriptions()
	if err != nil {
		return errors.Wrap(err, "could not get subscriptions")
	}

	repoSubs := subs.Repositories[repo]
	if repoSubs == nil {
		repoSubs = []*subscription.Subscription{sub}
	} else {
		exists := false
		for index, s := range repoSubs {
			if s.ChannelID == sub.ChannelID {
				repoSubs[index] = sub
				exists = true
				break
			}
		}

		if !exists {
			repoSubs = append(repoSubs, sub)
		}
	}

	subs.Repositories[repo] = repoSubs

	err = p.StoreSubscriptions(subs)
	if err != nil {
		return errors.Wrap(err, "could not store subscriptions")
	}

	return nil
}

func (p *Plugin) GetSubscriptions() (*subscription.Subscriptions, error) {
	var subscriptions *subscription.Subscriptions

	value, appErr := p.API.KVGet(SubscriptionsKey)
	if appErr != nil {
		return nil, errors.Wrap(appErr, "could not get subscriptions from KVStore")
	}

	if value == nil {
		return &subscription.Subscriptions{Repositories: map[string][]*subscription.Subscription{}}, nil
	}

	err := json.NewDecoder(bytes.NewReader(value)).Decode(&subscriptions)
	if err != nil {
		return nil, errors.Wrap(err, "could not properly decode subscriptions key")
	}

	return subscriptions, nil
}

func (p *Plugin) StoreSubscriptions(s *subscription.Subscriptions) error {
	b, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "error while converting subscriptions map to json")
	}

	if appErr := p.API.KVSet(SubscriptionsKey, b); appErr != nil {
		return errors.Wrap(appErr, "could not store subscriptions in KV store")
	}

	return nil
}

func (p *Plugin) GetSubscribedChannelsForRepository(pl webhookpayload.Payload) []*subscription.Subscription {
	name := pl.GetRepository().FullName
	org := strings.Split(name, "/")[0]
	subs, err := p.GetSubscriptions()
	if err != nil {
		return nil
	}

	// Add subscriptions for the specific repo
	var subsForRepo []*subscription.Subscription
	if subs.Repositories[name] != nil {
		subsForRepo = append(subsForRepo, subs.Repositories[name]...)
	}

	// Add subscriptions for the organization
	orgKey := fullNameFromOwnerAndRepo(org, "")
	if subs.Repositories[orgKey] != nil {
		subsForRepo = append(subsForRepo, subs.Repositories[orgKey]...)
	}

	if len(subsForRepo) == 0 {
		return nil
	}

	var subsToReturn []*subscription.Subscription

	for _, sub := range subsForRepo {
		if !p.permissionToRepo(sub.CreatorID, name) {
			continue
		}
		subsToReturn = append(subsToReturn, sub)
	}

	return subsToReturn
}

func (p *Plugin) Unsubscribe(channelID, repo string) (string, error) {
	if len(strings.Split(repo, "/")) != 2 {
		return requiredErrorMessage, nil
	}

	owner, repo := parseOwnerAndRepo(repo, p.getBaseURL())
	if owner == "" && repo == "" {
		return requiredErrorMessage, nil
	}
	repoWithOwner := fmt.Sprintf("%s/%s", owner, repo)

	subs, err := p.GetSubscriptions()
	if err != nil {
		return "", errors.Wrap(err, "could not get subscriptions")
	}

	repoSubs := subs.Repositories[repoWithOwner]
	if repoSubs == nil {
		return fmt.Sprintf(UnsubscribedErrorMessage, repo), nil
	}

	removed := false
	for index, sub := range repoSubs {
		if sub.ChannelID == channelID {
			repoSubs = append(repoSubs[:index], repoSubs[index+1:]...)
			removed = true
			break
		}
	}

	if removed {
		subs.Repositories[repoWithOwner] = repoSubs
		if err := p.StoreSubscriptions(subs); err != nil {
			return "", errors.Wrap(err, "could not store subscriptions")
		}
		return fmt.Sprintf("Successfully unsubscribed from %s.", repo), nil
	}

	return fmt.Sprintf(UnsubscribedErrorMessage, repo), nil
}
