package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const commandHelp = `* |/bitbucket connect| - Connect your Mattermost account to your Bitbucket account
* |/bitbucket disconnect| - Disconnect your Mattermost account from your Bitbucket account
* |/bitbucket todo| - Get a list of unread messages and pull requests awaiting your review
* |/bitbucket subscriptions list| - Will list the current channel subscriptions
* |/bitbucket subscriptions add owner [features]| - Subscribe the current channel to all available repositories within an organization and receive notifications about opened pull requests and issues
* |/bitbucket subscriptions add owner/repo [features]| - Subscribe the current channel to receive notifications about opened pull requests and issues for a repository
  * |features| is a comma-delimited list of one or more the following:
    * issues - includes new and closed issues
	* pulls - includes new and closed pull requests
    * pushes - includes pushes
    * creates - includes branch and tag creations
    * deletes - includes branch and tag deletions
    * issue_comments - includes new issue comments
    * pull_reviews - includes pull request reviews
  * Defaults to "pulls,issues,creates,deletes"
* |/bitbucket subscriptions delete owner/repo| - Unsubscribe the current channel from a repository
* |/bitbucket me| - Display the connected Bitbucket account
* |/bitbucket settings [setting] [value]| - Update your user settings
  * |setting| can be "notifications" or "reminders"
  * |value| can be "on" or "off"`

const (
	featureIssues        = "issues"
	featurePulls         = "pulls"
	featurePushes        = "pushes"
	featureCreates       = "creates"
	featureDeletes       = "deletes"
	featureIssueComments = "issue_comments"
	featurePullReviews   = "pull_reviews"
)

const (
	requiredErrorMessage = "Please specify an ogranization/repository."
)

var validFeatures = map[string]bool{
	featureIssues:        true,
	featurePulls:         true,
	featurePushes:        true,
	featureCreates:       true,
	featureDeletes:       true,
	featureIssueComments: true,
	featurePullReviews:   true,
}

// validateFeatures returns false when 1 or more given features
// are invalid along with a list of the invalid features.
func validateFeatures(features []string) (bool, []string) {
	var invalidFeatures []string
	valid := true
	hasLabel := false
	for _, f := range features {
		if _, ok := validFeatures[f]; ok {
			continue
		}
		if strings.HasPrefix(f, "label") {
			hasLabel = true
			continue
		}
		invalidFeatures = append(invalidFeatures, f)
		valid = false
	}
	if valid && hasLabel {
		// must have "pulls" or "issues" in features when using a label
		for _, f := range features {
			if f == featurePulls || f == featureIssues {
				return valid, invalidFeatures
			}
		}
		valid = false
	}
	return valid, invalidFeatures
}

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "bitbucket",
		DisplayName:      "Bitbucket",
		Description:      "Integration with Bitbucket.",
		AutoComplete:     true,
		AutocompleteData: getAutocompleteData(),
		AutoCompleteDesc: "Available commands: connect, disconnect, todo, me, settings, subscribe, unsubscribe, help",
		AutoCompleteHint: "[command]",
	}
}

func getAutocompleteData() *model.AutocompleteData {
	bitbucket := model.NewAutocompleteData("bitbucket", "[command]", "Available commands: connect, disconnect, todo, me, settings, subscribe, unsubscribe, help")

	connect := model.NewAutocompleteData("connect", "", "Connect your Mattermost account to your Bitbucket account")

	bitbucket.AddCommand(connect)

	disconnect := model.NewAutocompleteData("disconnect", "", "Disconnect your Mattermost account from your Bitbucket account")
	bitbucket.AddCommand(disconnect)

	help := model.NewAutocompleteData("help", "", "Display Slash Command help text")
	bitbucket.AddCommand(help)

	todo := model.NewAutocompleteData("todo", "", "Get a list of unread messages and pull requests awaiting your review")
	bitbucket.AddCommand(todo)

	me := model.NewAutocompleteData("me", "", "Display the connected Bitbucket account")
	bitbucket.AddCommand(me)

	subscriptions := model.NewAutocompleteData("subscriptions", "[command]", "Available commands: list, add")
	subscriptionsList := model.NewAutocompleteData("list", "", "List Subscription for this channel")
	subscriptions.AddCommand(subscriptionsList)

	subscriptionsAdd := model.NewAutocompleteData("add", "owner[/repo] features", "subscribe to org/[repo]")
	subscriptionsAdd.AddTextArgument("Owner/repo to subscribe to", "[owner/repo]", "")
	subscriptionsAdd.AddTextArgument("Comma-delimited list of one or more of: issues, pulls, pushes, creates, deletes, issue_comments, pull_reviews. Defaults to pulls,issues,creates,deletes", "[features] (optional)", `/[^,-\s]+(,[^,-\s]+)*/`)
	subscriptions.AddCommand(subscriptionsAdd)

	subscriptionsDelete := model.NewAutocompleteData("delete", "[owner/repo]", "Remove subscription for org/[repo]")
	subscriptions.AddCommand(subscriptionsDelete)

	bitbucket.AddCommand(subscriptions)

	settings := model.NewAutocompleteData("settings", "[setting] [value]", "Update your user settings")
	settingNotifications := model.NewAutocompleteData("notifications", "", "Turn notifications on/off")
	settingValue := []model.AutocompleteListItem{{
		HelpText: "Turn notifications on",
		Item:     "on",
	}, {
		HelpText: "Turn notifications off",
		Item:     "off",
	}}
	settingNotifications.AddStaticListArgument("", true, settingValue)
	settings.AddCommand(settingNotifications)
	bitbucket.AddCommand(settings)

	return bitbucket
}

func (p *Plugin) postCommandResponse(args *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: args.ChannelId,
		Message:   text,
	}
	_ = p.API.SendEphemeralPost(args.UserId, post)
}

func (p *Plugin) handleSubscribe(_ *plugin.Context, args *model.CommandArgs, parameters []string, userInfo *BitbucketUserInfo) string {
	features := "pulls,issues,creates,deletes"

	txt := ""
	switch parameters[0] {
	case "list":
		if len(parameters) > 1 {
			return "Invalid command."
		}

		subs, err := p.GetSubscriptionsByChannel(args.ChannelId)
		if err != nil {
			return err.Error()
		}

		if len(subs) == 0 {
			txt = "Currently there are no subscriptions in this channel"
		} else {
			txt = "### Subscriptions in this channel\n"
		}
		for _, sub := range subs {
			txt += fmt.Sprintf("* `%s` - %s", strings.Trim(sub.Repository, "/"), sub.Features)
			txt += "\n"
		}
		return txt
	case "delete":
		if len(parameters) != 2 {
			return requiredErrorMessage
		}

		repo := parameters[1]
		if err := p.Unsubscribe(args.ChannelId, repo); err != nil {
			p.API.LogError("Encountered an error trying to unsubscribe", "err", err.Error())
			return "Encountered an error trying to unsubscribe. Please try again."
		}

		return fmt.Sprintf("Successfully unsubscribed from %s.", repo)
	case "add":
		if len(parameters) < 2 {
			return requiredErrorMessage
		}

		parameters = parameters[1:]
		var optionList []string
		optionList = append(optionList, parameters[1:]...)

		if len(optionList) > 1 {
			return "Just one list of features is allowed"
		} else if len(optionList) == 1 {
			features = optionList[0]
			fs := strings.Split(features, ",")
			ok, ifs := validateFeatures(fs)
			if !ok {
				msg := fmt.Sprintf("Invalid feature(s) provided: %s", strings.Join(ifs, ","))
				if len(ifs) == 0 {
					msg = "Feature list must have \"pulls\" or \"issues\" when using a label."
				}
				return msg
			}
		}

		ctx := context.Background()
		bitbucketClient := p.bitbucketConnect(*userInfo.Token)
		owner, repo := parseOwnerAndRepo(parameters[0], BitbucketBaseURL)
		previousSubscribedEvents, err := p.findSubscriptionsEvents(args.ChannelId, owner, repo)
		if err != nil {
			return err.Error()
		}

		if previousSubscribedEvents == features {
			previousSubscribedEvents = ""
		}

		if repo == "" {
			if err = p.SubscribeOrg(ctx, bitbucketClient, args.UserId, owner, args.ChannelId, features); err != nil {
				return err.Error()
			}

			orgLink := fmt.Sprintf("%s%s", p.getBaseURL(), owner)
			msg := fmt.Sprintf("Successfully subscribed to organization [%s](%s) with events: %s", owner, orgLink, formattedString(features))
			if previousSubscribedEvents != "" {
				msg += fmt.Sprintf("\nThe previous subscription with: %s was overwritten.\n", formattedString(previousSubscribedEvents))
			}

			post := &model.Post{
				ChannelId: args.ChannelId,
				UserId:    p.BotUserID,
				Message:   msg,
			}

			if _, appErr := p.API.CreatePost(post); appErr != nil {
				p.API.LogWarn("error while creating post", "post", post, "error", appErr.Error())
				return fmt.Sprintf("%s Though there was an error creating the public post: %s", msg, appErr.Error())
			}

			return msg
		}

		if err = p.Subscribe(ctx, bitbucketClient, args.UserId, owner, repo, args.ChannelId, features); err != nil {
			return err.Error()
		}

		repoLink := fmt.Sprintf("%s%s/%s", p.getBaseURL(), owner, repo)

		msg := fmt.Sprintf("Successfully subscribed to [%s/%s](%s) with events: %s", owner, repo, repoLink, formattedString(features))
		if previousSubscribedEvents != "" {
			msg += fmt.Sprintf("\nThe previous subscription with: %s was overwritten.\n", formattedString(previousSubscribedEvents))
		}

		post := &model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.BotUserID,
			Message:   msg,
		}

		if _, appErr := p.API.CreatePost(post); appErr != nil {
			p.API.LogWarn("error while creating post", "post", post, "error", appErr.Error())
			return fmt.Sprintf("%s Though there was an error creating the public post: %s", msg, appErr.Error())
		}

		return msg
	}

	return "Invalid Command. commands available `add`, `delete` and `list`"
}

func (p *Plugin) findSubscriptionsEvents(channelID, owner, repo string) (string, error) {
	previouslySubscribed, err := p.GetSubscriptionsByChannel(channelID)
	if err != nil {
		return "", err
	}

	subscriptionName := owner
	if repo != "" {
		subscriptionName += "/" + repo
	}

	for _, subscribe := range previouslySubscribed {
		if subscribe.Repository == subscriptionName {
			return subscribe.Features, nil
		}
	}
	return "", nil
}

func formattedString(s string) string {
	return "`" + strings.Join(strings.Split(s, ","), "`, `") + "`"
}

func (p *Plugin) handleDisconnect(_ *plugin.Context, args *model.CommandArgs, _ []string, _ *BitbucketUserInfo) string {
	p.disconnectBitbucketAccount(args.UserId)
	return "Disconnected your Bitbucket account."
}

func (p *Plugin) handleTodo(_ *plugin.Context, _ *model.CommandArgs, _ []string, userInfo *BitbucketUserInfo) string {
	bitbucketClient := p.bitbucketConnect(*userInfo.Token)

	text, err := p.GetToDo(context.Background(), userInfo, bitbucketClient)
	if err != nil {
		p.API.LogError("Encountered an error getting your to do items", "err", err.Error())
		return "Encountered an error getting your to do items."
	}
	return text
}

func (p *Plugin) handleMe(_ *plugin.Context, _ *model.CommandArgs, _ []string, userInfo *BitbucketUserInfo) string {
	bitbucketClient := p.bitbucketConnect(*userInfo.Token)
	bitbucketUser, _, err := bitbucketClient.UsersApi.UserGet(context.Background())
	if err != nil {
		p.API.LogError("Encountered an error getting your Bitbucket profile", "err", err.Error())
		return "Encountered an error getting your Bitbucket profile."
	}

	text := fmt.Sprintf("You are connected to Bitbucket as:\n# [%s](%s)",
		bitbucketUser.Username, bitbucketUser.Links.Html.Href)
	return text
}

func (p *Plugin) handleHelp(_ *plugin.Context, _ *model.CommandArgs, _ []string, userInfo *BitbucketUserInfo) string {
	bitbucketClient := p.bitbucketConnect(*userInfo.Token)
	bitbucketUser, _, err := bitbucketClient.UsersApi.UserGet(context.Background())
	if err != nil {
		return "Encountered an error getting your Bitbucket profile info."
	}

	message := fmt.Sprintf("#### Welcome to the Mattermost Bitbucket Plugin!\n"+
		"You've connected your Mattermost account to [%s](%s) on Bitbucket. Read about the features of this plugin below:\n\n"+
		"##### Daily Reminders\n"+
		"The first time you log in each day, you will get a post right here letting you know what messages you need to read and what pull requests are awaiting your review.\n"+
		"Turn off reminders with `/bitbucket settings reminders off`.\n\n"+
		"##### Notifications\n"+
		"When someone mentions you, requests your review, comments on or modifies one of your pull requests/issues, or assigns you, you'll get a post here about it.\n"+
		"Turn off notifications with `/bitbucket settings notifications off`.\n\n"+
		"##### Sidebar Buttons\n"+
		"Check out the buttons in the left-hand sidebar of Mattermost.\n"+
		"* The first button tells you how many pull requests you have submitted.\n"+
		"* The second shows the number of PR that are awaiting your review.\n"+
		"* The third shows the number of PR and issues your are assiged to.\n"+
		"* The fourth will refresh the numbers.\n\n"+
		"Click on them!\n\n"+
		"##### Slash Commands\n"+
		strings.ReplaceAll(commandHelp, "|", "`"), bitbucketUser.Username, bitbucketUser.Links.Html.Href)

	return message
}

func (p *Plugin) handleSettings(_ *plugin.Context, _ *model.CommandArgs, parameters []string, userInfo *BitbucketUserInfo) string {
	if len(parameters) < 2 {
		return "Please specify both a setting and value. Use `/bitbucket help` for more usage information."
	}

	setting := parameters[0]
	if setting != SettingNotifications && setting != SettingReminders {
		return "Unknown setting."
	}

	strValue := parameters[1]
	value := false
	if strValue == SettingOn {
		value = true
	} else if strValue != SettingOff {
		return "Invalid value. Accepted values are: \"on\" or \"off\"."
	}

	if setting == SettingNotifications {
		if value {
			err := p.storeBitbucketAccountIDToMattermostUserIDMapping(userInfo.BitbucketAccountID, userInfo.UserID)
			if err != nil {
				p.API.LogError("Encountered an error storing Bitbucket account ID to Mattermost user ID mapping", "err", err.Error())
			}
		} else {
			err := p.API.KVDelete(userInfo.BitbucketUsername + BitbucketAccountIDKey)
			if err != nil {
				p.API.LogError("Encountered an error deleting Bitbucket account ID to Mattermost user ID mapping", "err", err.Error())
			}
		}

		userInfo.Settings.Notifications = value
	} else if setting == SettingReminders {
		userInfo.Settings.DailyReminder = value
	}

	err := p.storeBitbucketUserInfo(userInfo)
	if err != nil {
		p.API.LogError("Failed to store settings", "err", err.Error())
		return "Failed to store settings"
	}

	return "Settings updated."
}

type commandHandleFunc func(c *plugin.Context, args *model.CommandArgs, parameters []string, userInfo *BitbucketUserInfo) string

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand API.
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	var parameters []string
	action := ""
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		parameters = split[2:]
	}

	if command != "/bitbucket" {
		return &model.CommandResponse{}, nil
	}

	if action == "connect" {
		siteURL := p.API.GetConfig().ServiceSettings.SiteURL
		if siteURL == nil {
			p.postCommandResponse(args, "Encountered an error connecting to Bitbucket.")
			return &model.CommandResponse{}, nil
		}

		msg := fmt.Sprintf("[Click here to link your Bitbucket account.](%s/plugins/bitbucket/oauth/connect)", *siteURL)
		p.postCommandResponse(args, msg)
		return &model.CommandResponse{}, nil
	}

	info, apiErr := p.getBitbucketUserInfo(args.UserId)
	if apiErr != nil {
		text := "Unknown error."
		if apiErr.ID == APIErrorIDNotConnected {
			text = "You must connect your account to Bitbucket first. Either click on the Bitbucket logo in the bottom left of the screen or enter `/bitbucket connect`."
		}
		p.postCommandResponse(args, text)
		return &model.CommandResponse{}, nil
	}

	if f, ok := p.CommandHandlers[action]; ok {
		message := f(c, args, parameters, info)
		p.postCommandResponse(args, message)
		return &model.CommandResponse{}, nil
	}

	p.postCommandResponse(args, fmt.Sprintf("Unknown action %v", action))
	return &model.CommandResponse{}, nil
}
