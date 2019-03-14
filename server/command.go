package main

import (
	"context"
	"fmt"
	"strings"

	// "github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/plugin"

	// "github.com/google/go-github/github"
	"github.com/mattermost/mattermost-server/model"
	"github.com/wbrefvem/go-bitbucket"
)

const COMMAND_HELP = `* |/bitbucket connect| - Connect your Mattermost account to your Bitbucket account
* |/bitbucket disconnect| - Disconnect your Mattermost account from your * Bitbucket account
* |/bitbucket todo| - Get a list of unread messages and pull requests awaiting your review
* |/bitbucket subscribe list| - Will list the current channel subscriptions
* |/bitbucket subscribe owner [features]| - Subscribe the current channel to all available repositories within an organization and receive notifications about opened pull requests and issues
* |/bitbucket subscribe owner/repo [features]| - Subscribe the current channel to receive notifications about opened pull requests and issues for a repository
  * |features| is a comma-delimited list of one or more the following:
    * issues - includes new and closed issues
	* pulls - includes new and closed pull requests
    * pushes - includes pushes
    * creates - includes branch and tag creations
    * deletes - includes branch and tag deletions
    * issue_comments - includes new issue comments
    * pull_reviews - includes pull request reviews
	* label:"<labelname>" - must include "pulls" or "issues" in feature list when using a label
  * Defaults to "pulls,issues,creates,deletes"
* |/bitbucket unsubscribe owner/repo| - Unsubscribe the current channel from a repository
* |/bitbucket me| - Display the connected Bitbucket account
* |/bitbucket settings [setting] [value]| - Update your user settings
  * |setting| can be "notifications" or "reminders"
  * |value| can be "on" or "off"`

func getCommand() *model.Command {
	return &model.Command{
		Trigger:          "bitbucket",
		DisplayName:      "Bitbucket",
		Description:      "Integration with Bitbucket.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: connect, disconnect, todo, me, settings, subscribe, unsubscribe, help",
		AutoCompleteHint: "[command]",
	}
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     GITHUB_USERNAME,
		IconURL:      GITHUB_ICON_URL,
		Type:         model.POST_DEFAULT,
	}
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	fmt.Printf("----- BB command.ExecuteCommand \n    --> args=%+v", args)
	split := strings.Fields(args.Command)
	command := split[0]
	parameters := []string{}
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
		fmt.Println("----- BB command.ExecuteCommand action=connect")
		config := p.API.GetConfig()
		// fmt.Printf("----- config.ServiceSettings.SiteURL = %+v\n", config.ServiceSettings.SiteURL)
		if config.ServiceSettings.SiteURL == nil {
			// fmt.Printf("----- ServiceSettings.SiteURL = %+v\n", config.ServiceSettings.SiteURL)
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Encountered an error connecting to GitHub."), nil
		}
		resp := getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("[Click here to link your Bitbucket account.](%s/plugins/bitbucket/oauth/connect)", *config.ServiceSettings.SiteURL))
		fmt.Printf("----- BB command.ExecuteCommand resp = %+v\n", resp)
		return resp, nil
	}

	ctx := context.Background()
	fmt.Printf("----- BB command.ExecuteCommand --> \nctx = %+v\n", ctx)
	var githubClient *bitbucket.APIClient

	info, apiErr := p.getGitHubUserInfo(args.UserId)
	fmt.Printf("----- BB commnad.ExecuteCOmmand --> \ninfo = %+v\n", info)
	if apiErr != nil {
		text := "Unknown error."
		if apiErr.ID == API_ERROR_ID_NOT_CONNECTED {
			text = "You must connect your account to GitHub first. Either click on the GitHub logo in the bottom left of the screen or enter `/bitbucket connect`."
		}
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}

	githubClient = p.githubConnect(*info.Token)
	// gitUser, _, err := githubClient.UsersApi.UsersUsernameGet(ctx, "jfrerich")

	switch action {
	case "subscribe":
		/*
			/bitbucket subscribe jfrerich/mattermost-bitbucket-readme
		*/
		fmt.Println("----- BB command.ExecuteCommand [action=subscribe]")
		config := bitbucket.NewConfiguration()
		// config := p.getConfiguration()
		features := "pulls,issues,creates,deletes"
		fmt.Printf("config = %+v\n", config)
		fmt.Printf("features = %+v\n", features)

		txt := ""
		if len(parameters) == 0 {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please specify a repository or 'list' command."), nil
		} else if len(parameters) == 1 && parameters[0] == "list" {
			subs, err := p.GetSubscriptionsByChannel(args.ChannelId)
			if err != nil {
				return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), nil
			}

			if len(subs) == 0 {
				txt = "Currently there are no subscriptions in this channel"
			} else {
				txt = "### Subscriptions in this channel\n"
			}
			for _, sub := range subs {
				txt += fmt.Sprintf("* `%s` - %s\n", sub.Repository, sub.Features)
			}
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, txt), nil
		} else if len(parameters) > 1 {
			features = strings.Join(parameters[1:], " ")
		}

		// _, owner, repo := parseOwnerAndRepo(parameters[0], config.EnterpriseBaseURL)
		// if repo == "" {
		// 	if err := p.SubscribeOrg(context.Background(), githubClient, args.UserId, owner, args.ChannelId, features); err != nil {
		// 		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), nil
		// 	}
		//
		// 	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Successfully subscribed to organization %s.", owner)), nil
		// }
		//
		// if err := p.Subscribe(context.Background(), githubClient, args.UserId, owner, repo, args.ChannelId, features); err != nil {
		// 	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, err.Error()), nil
		// }

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Successfully subscribed to %s.", "jfrerich repo")), nil
		// return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Successfully subscribed to %s.", repo)), nil
	case "unsubscribe":
		// if len(parameters) == 0 {
		// 	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please specify a repository."), nil
		// }
		//
		// repo := parameters[0]
		//
		// if err := p.Unsubscribe(args.ChannelId, repo); err != nil {
		// 	mlog.Error(err.Error())
		// 	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Encountered an error trying to unsubscribe. Please try again."), nil
		// }
		//
		// return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, fmt.Sprintf("Succesfully unsubscribed from %s.", repo)), nil
	case "disconnect":
		fmt.Println("----- BB command.ExecuteCommand action=disconnect")
		p.disconnectGitHubAccount(args.UserId)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Disconnected your Bitbucket account."), nil
	case "todo":
		// text, err := p.GetToDo(ctx, info.GitHubUsername, githubClient)
		// if err != nil {
		// 	mlog.Error(err.Error())
		// 	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Encountered an error getting your to do items."), nil
		// }
		// return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	case "me":
		fmt.Printf("----- BB command.ExecuteCommand action=me")
		// gitUser, _, err := githubClient.Users.Get(ctx, "")
		//
		// gitUser, _, err := githubClient.UsersApi.UsersUsernameGet(ctx, "jfrerich")
		gitUser, _, err := githubClient.UsersApi.UserGet(ctx)
		// fmt.Printf("----- BB command.ExecuteCommand ctx = %+v\n", ctx)
		fmt.Printf("----- BB command.ExecuteCommand action=me\n\n gitUser -> %+v", gitUser)
		avatar := gitUser.Links.Avatar.Href
		html := gitUser.Links.Html.Href
		username := gitUser.Username
		if err != nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Encountered an error getting your GitHub profile."), nil
		}

		// text := fmt.Sprintf("You are connected to GitHub as:\n# [![image](%s =40x40)](%s) [%s](%s)", gitUser.AvatarGet(), gitUser.GetHTMLURL(), gitUser.GetLogin(), gitUser.GetHTMLURL())
		text := fmt.Sprintf("You are connected to Bitbucket as:\n# [![image](%s =40x40)](%s) [%s](%s)", avatar, html, username, html)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	case "help":
		text := "###### Mattermost Bitbucket Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	case "":
		text := "###### Mattermost Bitbucket Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	case "settings":
		if len(parameters) < 2 {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Please specify both a setting and value. Use `/github help` for more usage information."), nil
		}

		setting := parameters[0]
		if setting != SETTING_NOTIFICATIONS && setting != SETTING_REMINDERS {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Unknown setting."), nil
		}

		strValue := parameters[1]
		value := false
		if strValue == SETTING_ON {
			value = true
		} else if strValue != SETTING_OFF {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Invalid value. Accepted values are: \"on\" or \"off\"."), nil
		}

		if setting == SETTING_NOTIFICATIONS {
			if value {
				p.storeGitHubToUserIDMapping(info.GitHubUsername, info.UserID)
			} else {
				p.API.KVDelete(info.GitHubUsername + GITHUB_USERNAME_KEY)
			}

			info.Settings.Notifications = value
		} else if setting == SETTING_REMINDERS {
			info.Settings.DailyReminder = value
		}

		p.storeGitHubUserInfo(info)

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Settings updated."), nil
	}

	return &model.CommandResponse{}, nil
}
