package main

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-plugin-bitbucket/server/subscription"
)

func CheckError(t *testing.T, wantErr bool, err error) {
	message := "should return no error"
	if wantErr {
		message = "should return error"
	}
	assert.Equal(t, wantErr, err != nil, message)
}

// pluginWithMockedSubs returns mocked plugin for given subscriptions
func pluginWithMockedSubs(subscriptions []*subscription.Subscription) *Plugin {
	p := NewPlugin()
	mockPluginAPI := &plugintest.API{}

	subs := subscription.Subscriptions{Repositories: map[string][]*subscription.Subscription{}}
	subs.Repositories[""] = subscriptions
	jsn, _ := json.Marshal(subs)
	mockPluginAPI.On("KVGet", SubscriptionsKey).Return(jsn, nil)
	p.SetAPI(mockPluginAPI)
	return p
}

// wantedSubscriptions returns what should be returned after sorting by repo names
func wantedSubscriptions(repoNames []string, chanelID string) []*subscription.Subscription {
	var subs []*subscription.Subscription
	for _, st := range repoNames {
		subs = append(subs, &subscription.Subscription{
			ChannelID:  chanelID,
			Repository: st,
		})
	}
	return subs
}

func TestPlugin_GetSubscriptionsByChannel(t *testing.T) {
	type args struct {
		channelID string
	}
	tests := []struct {
		name    string
		plugin  *Plugin
		args    args
		want    []*subscription.Subscription
		wantErr bool
	}{
		{
			name: "basic test",
			args: args{channelID: "1"},
			plugin: pluginWithMockedSubs([]*subscription.Subscription{
				{
					ChannelID:  "1",
					Repository: "asd",
					CreatorID:  "1",
				},
				{
					ChannelID:  "1",
					Repository: "123",
					CreatorID:  "1",
				},
				{
					ChannelID:  "1",
					Repository: "",
					CreatorID:  "1",
				},
			}),
			want:    wantedSubscriptions([]string{"", "123", "asd"}, "1"),
			wantErr: false,
		},
		{
			name:    "test empty",
			args:    args{channelID: "1"},
			plugin:  pluginWithMockedSubs([]*subscription.Subscription{}),
			want:    wantedSubscriptions([]string{}, "1"),
			wantErr: false,
		},
		{
			name: "test shuffled",
			args: args{channelID: "1"},
			plugin: pluginWithMockedSubs([]*subscription.Subscription{
				{
					ChannelID:  "1",
					Repository: "c",
					CreatorID:  "3",
				},
				{
					ChannelID:  "1",
					Repository: "b",
					CreatorID:  "3",
				},
				{
					ChannelID:  "1",
					Repository: "ab",
					CreatorID:  "3",
				},
				{
					ChannelID:  "1",
					Repository: "a",
					CreatorID:  "3",
				},
			}),
			want:    wantedSubscriptions([]string{"a", "ab", "b", "c"}, "1"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.plugin.GetSubscriptionsByChannel(tt.args.channelID)

			CheckError(t, tt.wantErr, err)

			assert.Equal(t, tt.want, got, "they should be same")
		})
	}
}
