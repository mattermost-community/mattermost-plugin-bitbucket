package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"

	"github.com/google/go-github/github"
	"github.com/wbrefvem/go-bitbucket"
	"golang.org/x/oauth2"
)

const (
	GITHUB_TOKEN_KEY        = "_githubtoken"
	GITHUB_STATE_KEY        = "_githubstate"
	GITHUB_USERNAME_KEY     = "_githubusername"
	GITHUB_PRIVATE_REPO_KEY = "_githubprivate"
	WS_EVENT_CONNECT        = "connect"
	WS_EVENT_DISCONNECT     = "disconnect"
	WS_EVENT_REFRESH        = "refresh"
	SETTING_BUTTONS_TEAM    = "team"
	SETTING_BUTTONS_CHANNEL = "channel"
	SETTING_BUTTONS_OFF     = "off"
	SETTING_NOTIFICATIONS   = "notifications"
	SETTING_REMINDERS       = "reminders"
	SETTING_ON              = "on"
	SETTING_OFF             = "off"
)

type Plugin struct {
	plugin.MattermostPlugin
	// githubClient    *github.Client
	bitbucketClient *bitbucket.APIClient

	BotUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// func (p *Plugin) githubConnect(token oauth2.Token) (*bitbucket.APIClient, *context.valueCtx) {
func (p *Plugin) githubConnect(token oauth2.Token) (*bitbucket.APIClient, context.Context) {

	// plugin configuration
	// config := p.getConfiguration()
	// fmt.Println(reflect.TypeOf(config).String())
	// fmt.Printf("auth = %+v\n", auth)

	// fmt.Printf("----- BB plugin.githubConnect  ->  config = %v\n", config)

	// get Oauth token source and client
	// ctx := context.Background()
	ts := oauth2.StaticTokenSource(&token)
	// tc := oauth2.NewClient(ctx, ts)
	// fmt.Printf("----- BB plugin.githubConnect  ->  &token = %v\n", &token)
	// fmt.Printf("----- BB plugin.githubConnect  ->  tc = %v\n", tc)
	// fmt.Printf("----- BB plugin.githubConnect  ->  ts = %v\n", ts)

	// ---- setup Oauth context ----
	auth := context.WithValue(oauth2.NoContext, bitbucket.ContextOAuth2, ts)
	// auth := context.WithValue(context.Background(), bitbucket.ContextBasicAuth, bitbucket.BasicAuth{
	// 	UserName: "jasonfrerich@yahoo.com",
	// 	Password: "wrongpassword",
	// })

	// create config for bitbucket API
	config_bb := bitbucket.NewConfiguration()

	// create new bitbucket client API
	new_client := bitbucket.NewAPIClient(config_bb)

	// prove the bitbucket API client works
	// comment, _, err := new_client.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesIssueIdCommentsCommentIdGet(ctx, "51110066", "jfrerich", "mattermost-bitbucket-readme", "1")
	// fmt.Printf("----- 2. BB plugin.githubConnect  ->  comment = %+v\n", comment)
	// fmt.Printf("----- 2. BB plugin.githubConnect  ->  err = %+v\n", err)
	//
	// // ---- test bitbucket API with Oauth ----
	// comment_w, _, err_w := new_client.RepositoriesApi.RepositoriesUsernameRepoSlugGet(auth, "jfrerich", "solid")
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> comment_w = %+v\n", comment_w)
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> err_w = %+v\n", err_w)
	//
	// // ---- test bitbucket API without Oauth ----
	// comment_wo, _, err_wo := new_client.RepositoriesApi.RepositoriesUsernameRepoSlugGet(context.Background(), "jfrerich", "solid")
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> comment_wo = %+v\n", comment_wo)
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> err_wo = %+v\n", err_wo)
	// // ---- test bitbucket API with Oauth ----
	//
	// gitUser2, _, err2 := new_client.UsersApi.UserGet(auth)
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> gitUser2 = %+v\n", gitUser2)
	// fmt.Printf("----- #### BB api.completeConnectUserToGitHub  -> err2 = %+v\n", err2)

	// fmt.Printf("----- 2. BB plugin.githubConnect  ->  config_bb = %v", config_bb)
	// TODO figure out how to add auth to client so dont' have to return it
	return new_client, auth

}

func (p *Plugin) OnActivate() error {

	config := p.getConfiguration()
	fmt.Printf("----- BB plugin.OnActivate  ->  config = %+v \n", config)

	if err := config.IsValid(); err != nil {
		return err
	}
	p.API.RegisterCommand(getCommand())
	user, err := p.API.GetUserByUsername(config.Username)
	fmt.Printf("----- BB plugin.OnActivate  ->  BotUsername = %+v \n", user)
	if err != nil {
		mlog.Error(err.Error())
		return fmt.Errorf("Unable to find user with configured username: %v", config.Username)
	}

	p.BotUserID = user.Id
	return nil
}

func (p *Plugin) getOAuthConfig() *oauth2.Config {

	config := p.getConfiguration()

	authURL, _ := url.Parse("https://bitbucket.org/")
	tokenURL, _ := url.Parse("https://bitbucket.org/")

	if len(config.EnterpriseBaseURL) > 0 {
		authURL, _ = url.Parse(config.EnterpriseBaseURL)
		tokenURL, _ = url.Parse(config.EnterpriseBaseURL)
	}

	authURL.Path = path.Join(authURL.Path, "site", "oauth2", "authorize")
	tokenURL.Path = path.Join(tokenURL.Path, "site", "oauth2", "access_token")

	repo := "public_repo"
	if config.EnablePrivateRepo {
		// means that asks scope for privaterepositories
		repo = "repo"
	}

	fmt.Println("----- BB plugin.getOAuthconfig  ->  repo =", repo)

	fmt.Println("TODO -> check Scopes statement -> diffes from GH")
	return &oauth2.Config{
		ClientID:     config.BitbucketOAuthClientID,
		ClientSecret: config.BitbucketOAuthClientSecret,
		Scopes:       []string{"repository"},
		// Scopes:       []string{repo, "notifications"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL.String(),
			TokenURL: tokenURL.String(),
		},
	}
}

type GitHubUserInfo struct {
	UserID              string
	Token               *oauth2.Token
	GitHubUsername      string
	LastToDoPostAt      int64
	Settings            *UserSettings
	AllowedPrivateRepos bool
}

type UserSettings struct {
	SidebarButtons string `json:"sidebar_buttons"`
	DailyReminder  bool   `json:"daily_reminder"`
	Notifications  bool   `json:"notifications"`
}

func (p *Plugin) storeGitHubUserInfo(info *GitHubUserInfo) error {
	config := p.getConfiguration()

	encryptedToken, err := encrypt([]byte(config.EncryptionKey), info.Token.AccessToken)
	if err != nil {
		return err
	}

	info.Token.AccessToken = encryptedToken

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		return err
	}

	if err := p.API.KVSet(info.UserID+GITHUB_TOKEN_KEY, jsonInfo); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) getGitHubUserInfo(userID string) (*GitHubUserInfo, *APIErrorResponse) {
	config := p.getConfiguration()

	// fmt.Println("----- BB plugin.getGitHubUserInfo -----")
	var userInfo GitHubUserInfo

	if infoBytes, err := p.API.KVGet(userID + GITHUB_TOKEN_KEY); err != nil || infoBytes == nil {
		return nil, &APIErrorResponse{ID: API_ERROR_ID_NOT_CONNECTED, Message: "Must connect user account to GitHub first.", StatusCode: http.StatusBadRequest}
	} else if err := json.Unmarshal(infoBytes, &userInfo); err != nil {
		return nil, &APIErrorResponse{ID: "", Message: "Unable to parse token.", StatusCode: http.StatusInternalServerError}
	}

	unencryptedToken, err := decrypt([]byte(config.EncryptionKey), userInfo.Token.AccessToken)
	fmt.Printf("----- BB plugin.getGitHubUserInfo  -->  unencryptedToken = %+v\n", unencryptedToken)
	if err != nil {
		mlog.Error(err.Error())
		return nil, &APIErrorResponse{ID: "", Message: "Unable to decrypt access token.", StatusCode: http.StatusInternalServerError}
	}

	userInfo.Token.AccessToken = unencryptedToken

	return &userInfo, nil
}

func (p *Plugin) storeGitHubToUserIDMapping(githubUsername, userID string) error {
	if err := p.API.KVSet(githubUsername+GITHUB_USERNAME_KEY, []byte(userID)); err != nil {
		return fmt.Errorf("Encountered error saving github username mapping")
	}
	return nil
}

func (p *Plugin) getGitHubToUserIDMapping(githubUsername string) string {
	userID, _ := p.API.KVGet(githubUsername + GITHUB_USERNAME_KEY)
	return string(userID)
}

func (p *Plugin) disconnectGitHubAccount(userID string) {
	userInfo, _ := p.getGitHubUserInfo(userID)
	if userInfo == nil {
		return
	}

	p.API.KVDelete(userID + GITHUB_TOKEN_KEY)
	p.API.KVDelete(userInfo.GitHubUsername + GITHUB_USERNAME_KEY)

	if user, err := p.API.GetUser(userID); err == nil && user.Props != nil && len(user.Props["git_user"]) > 0 {
		delete(user.Props, "git_user")
		p.API.UpdateUser(user)
	}

	p.API.PublishWebSocketEvent(
		WS_EVENT_DISCONNECT,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}

func (p *Plugin) CreateBotDMPost(userID, message, postType string) *model.AppError {
	channel, err := p.API.GetDirectChannel(userID, p.BotUserID)
	if err != nil {
		mlog.Error("Couldn't get bot's DM channel", mlog.String("user_id", userID))
		return err
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message,
		Type:      postType,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": GITHUB_USERNAME,
			"override_icon_url": GITHUB_ICON_URL,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		mlog.Error(err.Error())
		return err
	}

	return nil
}

func (p *Plugin) PostToDo(info *GitHubUserInfo) {
	// text, err := p.GetToDo(context.Background(), info.GitHubUsername, p.githubConnect(*info.Token))
	// if err != nil {
	// 	mlog.Error(err.Error())
	// 	return
	// }
	//
	// p.CreateBotDMPost(info.UserID, text, "custom_git_todo")
}

func (p *Plugin) GetToDo(ctx context.Context, username string, githubClient *github.Client) (string, error) {
	config := p.getConfiguration()

	issueResults, _, err := githubClient.Search.Issues(ctx, getReviewSearchQuery(username, config.GitHubOrg), &github.SearchOptions{})
	if err != nil {
		return "", err
	}

	notifications, _, err := githubClient.Activity.ListNotifications(ctx, &github.NotificationListOptions{})
	if err != nil {
		return "", err
	}

	yourPrs, _, err := githubClient.Search.Issues(ctx, getYourPrsSearchQuery(username, config.GitHubOrg), &github.SearchOptions{})
	if err != nil {
		return "", err
	}

	yourAssignments, _, err := githubClient.Search.Issues(ctx, getYourAssigneeSearchQuery(username, config.GitHubOrg), &github.SearchOptions{})
	if err != nil {
		return "", err
	}

	text := "##### Unread Messages\n"

	notificationCount := 0
	notificationContent := ""
	for _, n := range notifications {
		if n.GetReason() == "subscribed" {
			continue
		}

		if n.GetRepository() == nil {
			p.API.LogError("Unable to get repository for notification in todo list. Skipping.")
			continue
		}

		if p.checkOrg(n.GetRepository().GetOwner().GetLogin()) != nil {
			continue
		}

		switch n.GetSubject().GetType() {
		case "RepositoryVulnerabilityAlert":
			message := fmt.Sprintf("[Vulnerability Alert for %v](%v)", n.GetRepository().GetFullName(), fixGithubNotificationSubjectURL(n.GetSubject().GetURL()))
			notificationContent += fmt.Sprintf("* %v\n", message)
		default:
			url := fixGithubNotificationSubjectURL(n.GetSubject().GetURL())
			notificationContent += fmt.Sprintf("* %v\n", url)
		}

		notificationCount++
	}

	if notificationCount == 0 {
		text += "You don't have any unread messages.\n"
	} else {
		text += fmt.Sprintf("You have %v unread messages:\n", notificationCount)
		text += notificationContent
	}

	text += "##### Review Requests\n"

	if issueResults.GetTotal() == 0 {
		text += "You have don't have any pull requests awaiting your review.\n"
	} else {
		text += fmt.Sprintf("You have %v pull requests awaiting your review:\n", issueResults.GetTotal())

		for _, pr := range issueResults.Issues {
			text += fmt.Sprintf("* %v\n", pr.GetHTMLURL())
		}
	}

	text += "##### Your Open Pull Requests\n"

	if yourPrs.GetTotal() == 0 {
		text += "You have don't have any open pull requests.\n"
	} else {
		text += fmt.Sprintf("You have %v open pull requests:\n", yourPrs.GetTotal())

		for _, pr := range yourPrs.Issues {
			text += fmt.Sprintf("* %v\n", pr.GetHTMLURL())
		}
	}

	text += "##### Your Assignments\n"

	if yourAssignments.GetTotal() == 0 {
		text += "You have don't have any assignments.\n"
	} else {
		text += fmt.Sprintf("You have %v assignments:\n", yourAssignments.GetTotal())

		for _, assign := range yourAssignments.Issues {
			text += fmt.Sprintf("* %v\n", assign.GetHTMLURL())
		}
	}

	return text, nil
}

func (p *Plugin) checkOrg(org string) error {
	config := p.getConfiguration()

	configOrg := strings.TrimSpace(config.GitHubOrg)
	if configOrg != "" && configOrg != org {
		return fmt.Errorf("Only repositories in the %v organization are supported", configOrg)
	}

	return nil
}

func (p *Plugin) sendRefreshEvent(userID string) {
	p.API.PublishWebSocketEvent(
		WS_EVENT_REFRESH,
		nil,
		&model.WebsocketBroadcast{UserId: userID},
	)
}
