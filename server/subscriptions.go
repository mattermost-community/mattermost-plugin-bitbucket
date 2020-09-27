package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/subscription"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook_payload"
	"github.com/wbrefvem/go-bitbucket"
)

const (
	SUBSCRIPTIONS_KEY = "subscriptions"
)

func (p *Plugin) Subscribe(ctx context.Context, bitbucketClient *bitbucket.APIClient, userId, owner, repo, channelID, features string) error {
	if owner == "" {
		return fmt.Errorf("Invalid repository")
	}

	if err := p.checkOrg(owner); err != nil {
		return err
	}

	result, _, err := bitbucketClient.RepositoriesApi.RepositoriesUsernameRepoSlugGet(ctx, owner, repo)
	fmt.Printf("result = %+v\n", result)
	if err != nil {
		p.API.LogError(err.Error())
		return fmt.Errorf("Unknown repository %s/%s", owner, repo)
	}
	// if result, _, err := bitbucketClient.RepositoriesApi.RepositoriesUsernameRepoSlugGet(ctx, owner, repo); result == nil || err != nil {
	// 	if err != nil {
	// 		mlog.Error(err.Error())
	// 	}
	// 	return fmt.Errorf("Unknown repository %s/%s", owner, repo)
	// }

	fmt.Println("--- LETS SUBSCRIBE --")
	fmt.Printf("--- ChannelID = %+v\n", channelID)
	fmt.Printf("--- userid = %+v\n", userId)

	sub := &subscription.Subscription{
		ChannelID:  channelID,
		CreatorID:  userId,
		Features:   features,
		Repository: fmt.Sprintf("%s/%s", owner, repo),
	}

	if err := p.AddSubscription(fmt.Sprintf("%s/%s", owner, repo), sub); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) SubscribeOrg(ctx context.Context, bitbucketClient *github.Client, userId, org, channelID, features string) error {
	if org == "" {
		return fmt.Errorf("Invalid organization")
	}
	if err := p.checkOrg(org); err != nil {
		return err
	}

	listOrgOptions := github.RepositoryListByOrgOptions{
		Type: "all",
	}
	repos, _, err := bitbucketClient.Repositories.ListByOrg(ctx, org, &listOrgOptions)
	if repos == nil || err != nil {
		if err != nil {
			p.API.LogError(err.Error())
		}
		return fmt.Errorf("Unknown organization %s", org)
	}

	for _, repo := range repos {
		sub := &subscription.Subscription{
			ChannelID:  channelID,
			CreatorID:  userId,
			Features:   features,
			Repository: fmt.Sprintf("%s/%s", org, repo.GetFullName()),
		}

		if err := p.AddSubscription(fmt.Sprintf("%s/%s", org, repo), sub); err != nil {
			continue
		}
	}

	return nil
}

func (p *Plugin) GetSubscriptionsByChannel(channelID string) ([]*subscription.Subscription, error) {
	var filteredSubs []*subscription.Subscription
	subs, err := p.GetSubscriptions()
	if err != nil {
		return nil, err
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

	return filteredSubs, nil
}

func (p *Plugin) AddSubscription(repo string, sub *subscription.Subscription) error {
	subs, err := p.GetSubscriptions()
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (p *Plugin) GetSubscriptions() (*subscription.Subscriptions, error) {
	var subscriptions *subscription.Subscriptions

	value, err := p.API.KVGet(SUBSCRIPTIONS_KEY)
	if err != nil {
		return nil, err
	}

	if value == nil {
		subscriptions = &subscription.Subscriptions{Repositories: map[string][]*subscription.Subscription{}}
	} else {
		json.NewDecoder(bytes.NewReader(value)).Decode(&subscriptions)
	}

	return subscriptions, nil
}

func (p *Plugin) StoreSubscriptions(s *subscription.Subscriptions) error {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	p.API.KVSet(SUBSCRIPTIONS_KEY, b)
	return nil
}

func (p *Plugin) GetSubscribedChannelsForRepository(pl webhook_payload.Payload) []*subscription.Subscription {
	fmt.Println("---- GetSubscribedChannelsForRepository ----")
	subs, err := p.GetSubscriptions()
	if err != nil {
		return nil
	}

	fmt.Printf("---> name = %+v\n", pl.GetRepository().FullName)
	subsForRepo := subs.Repositories[pl.GetRepository().FullName]
	if subsForRepo == nil {
		return nil
	}
	fmt.Printf("---> subsForRepo = %+v\n", subsForRepo)

	subsToReturn := []*subscription.Subscription{}

	for _, sub := range subsForRepo {
		// if repo.GetPrivate() && !p.permissionToRepo(sub.CreatorID, name) {
		fmt.Printf("----> sub = %+v\n", sub)
		// fmt.Printf("isprivate = %+v\n", isprivate)
		fmt.Printf("----> sub.CreatorID = %+v\n", sub.CreatorID)
		if pl.GetRepository().IsPrivate && !p.permissionToRepo(sub.CreatorID, pl.GetRepository().FullName) {
			continue
		}
		subsToReturn = append(subsToReturn, sub)
	}

	return subsToReturn
}

func (p *Plugin) Unsubscribe(channelID string, repo string) error {
	config := p.getConfiguration()

	repo, _, _ = parseOwnerAndRepo(repo, config.EnterpriseBaseURL)

	if repo == "" {
		return fmt.Errorf("Invalid repository")
	}

	subs, err := p.GetSubscriptions()
	if err != nil {
		return err
	}

	repoSubs := subs.Repositories[repo]
	if repoSubs == nil {
		return nil
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
		subs.Repositories[repo] = repoSubs
		if err := p.StoreSubscriptions(subs); err != nil {
			return err
		}
	}

	return nil
}
