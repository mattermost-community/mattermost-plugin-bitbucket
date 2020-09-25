package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/template_renderer"
	"github.com/kosgrz/mattermost-plugin-bitbucket/server/webhook"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/google/go-github/github"
	"github.com/wbrefvem/go-bitbucket"
	"golang.org/x/oauth2"
)

const (
	BITBUCKET_TOKEN_KEY        = "_bitbuckettoken"
	BITBUCKET_USERNAME_KEY     = "_bitbucketusername"
	BITBUCKET_ACCOUNT_ID_KEY   = "_bitbucketaccountid"
	BITBUCKET_PRIVATE_REPO_KEY = "_bitbucketprivate"
	WS_EVENT_CONNECT           = "connect"
	WS_EVENT_DISCONNECT        = "disconnect"
	WS_EVENT_REFRESH           = "refresh"
	SETTING_BUTTONS_TEAM       = "team"
	SETTING_NOTIFICATIONS      = "notifications"
	SETTING_REMINDERS          = "reminders"
	SETTING_ON                 = "on"
	SETTING_OFF                = "off"
)

type Plugin struct {
	plugin.MattermostPlugin
	// bitbucketClient    *github.Client
	bitbucketClient *bitbucket.APIClient

	BotUserID string

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	// webhookHandler is responsible for handling webhook events.
	webhookHandler webhook.Webhook
}

// func (p *Plugin) bitbucketConnect(token oauth2.Token) (*bitbucket.APIClient, *context.valueCtx) {
func (p *Plugin) bitbucketConnect(token oauth2.Token) *bitbucket.APIClient {

	// get Oauth token source and client
	ts := oauth2.StaticTokenSource(&token)

	// setup Oauth context
	auth := context.WithValue(oauth2.NoContext, bitbucket.ContextOAuth2, ts)

	tc := oauth2.NewClient(auth, ts)

	// create config for bitbucket API
	configBb := bitbucket.NewConfiguration()
	configBb.HTTPClient = tc

	// create new bitbucket client API
	newClient := bitbucket.NewAPIClient(configBb)

	// TODO figure out how to add auth to client so dont' have to return it
	return newClient
}

func (p *Plugin) OnActivate() error {

	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		return errors.Wrap(err, "invalid config")
	}

	if p.API.GetConfig().ServiceSettings.SiteURL == nil {
		return errors.New("siteURL is not set. Please set a siteURL and restart the plugin")
	}

	err := p.API.RegisterCommand(getCommand())
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	botID, err := p.Helpers.EnsureBot(&model.Bot{
		Username:    "bitbucket",
		DisplayName: "BitBucket",
		Description: "Created by the BitBucket plugin.",
	})
	if err != nil {
		return errors.Wrap(err, "failed to ensure BitBucket bot")
	}

	p.BotUserID = botID

	templateRenderer := template_renderer.MakeTemplateRenderer()
	templateRenderer.RegisterBitBucketAccountIDToUsernameMappingCallback(p.getBitBucketAccountIDToMattermostUsernameMapping)
	p.webhookHandler = webhook.NewWebhook(&subscriptionHandler{p}, &pullRequestReviewHandler{p}, templateRenderer)

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
	fmt.Printf("TODO : determine proper repo scrope for bitbucket %v", repo)

	fmt.Println("TODO -> check Scopes statement -> differs from GH")
	return &oauth2.Config{
		ClientID:     config.BitbucketOAuthClientID,
		ClientSecret: config.BitbucketOAuthClientSecret,
		Scopes:       []string{"repository"},
		RedirectURL:  fmt.Sprintf("%s/plugins/%s/oauth/complete", *p.API.GetConfig().ServiceSettings.SiteURL, manifest.Id),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL.String(),
			TokenURL: tokenURL.String(),
		},
	}
}

type BitbucketUserInfo struct {
	UserID              string
	Token               *oauth2.Token
	BitbucketUsername   string
	LastToDoPostAt      int64
	Settings            *UserSettings
	AllowedPrivateRepos bool
}

type UserSettings struct {
	SidebarButtons string `json:"sidebar_buttons"`
	DailyReminder  bool   `json:"daily_reminder"`
	Notifications  bool   `json:"notifications"`
}

func (p *Plugin) storeBitbucketUserInfo(info *BitbucketUserInfo) error {
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

	if err := p.API.KVSet(info.UserID+BITBUCKET_TOKEN_KEY, jsonInfo); err != nil {
		return err
	}

	return nil
}

func (p *Plugin) getBitbucketUserInfo(userID string) (*BitbucketUserInfo, *APIErrorResponse) {
	config := p.getConfiguration()

	var userInfo BitbucketUserInfo

	if infoBytes, err := p.API.KVGet(userID + BITBUCKET_TOKEN_KEY); err != nil || infoBytes == nil {
		return nil, &APIErrorResponse{ID: API_ERROR_ID_NOT_CONNECTED, Message: "Must connect user account to Bitbucket first.", StatusCode: http.StatusBadRequest}
	} else if err := json.Unmarshal(infoBytes, &userInfo); err != nil {
		return nil, &APIErrorResponse{ID: "", Message: "Unable to parse token.", StatusCode: http.StatusInternalServerError}
	}

	unencryptedToken, err := decrypt([]byte(config.EncryptionKey), userInfo.Token.AccessToken)
	if err != nil {
		p.API.LogError(err.Error())
		return nil, &APIErrorResponse{ID: "", Message: "Unable to decrypt access token.", StatusCode: http.StatusInternalServerError}
	}

	userInfo.Token.AccessToken = unencryptedToken

	return &userInfo, nil
}

func (p *Plugin) storeBitbucketToUserIDMapping(bitbucketUsername, userID string) error {
	if err := p.API.KVSet(bitbucketUsername+BITBUCKET_USERNAME_KEY, []byte(userID)); err != nil {
		return fmt.Errorf("Encountered error saving BitBucket username mapping")
	}
	return nil
}

func (p *Plugin) storeBitbucketUserIDToMattermostUserIDMapping(bitbucketUserID, userID string) error {
	if err := p.API.KVSet(bitbucketUserID+BITBUCKET_ACCOUNT_ID_KEY, []byte(userID)); err != nil {
		return fmt.Errorf("Encountered error saving BitBucket user ID mapping")
	}
	return nil
}

func (p *Plugin) getBitbucketToUserIDMapping(bitbucketUsername string) string {
	userID, _ := p.API.KVGet(bitbucketUsername + BITBUCKET_USERNAME_KEY)
	return string(userID)
}

func (p *Plugin) getBitbucketAccountIDToMattermostUserIDMapping(bitbucketAccountID string) string {
	userID, _ := p.API.KVGet(bitbucketAccountID + BITBUCKET_ACCOUNT_ID_KEY)
	return string(userID)
}

func (p *Plugin) disconnectBitbucketAccount(userID string) {
	userInfo, _ := p.getBitbucketUserInfo(userID)
	if userInfo == nil {
		return
	}

	p.API.KVDelete(userID + BITBUCKET_TOKEN_KEY)
	p.API.KVDelete(userInfo.BitbucketUsername + BITBUCKET_USERNAME_KEY)

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
		p.API.LogError("Couldn't get bot's DM channel")
		return err
	}

	post := &model.Post{
		UserId:    p.BotUserID,
		ChannelId: channel.Id,
		Message:   message,
		Type:      postType,
		Props: map[string]interface{}{
			"from_webhook":      "true",
			"override_username": BITBUCKET_USERNAME,
			"override_icon_url": BITBUCKET_ICON_URL,
		},
	}

	if _, err := p.API.CreatePost(post); err != nil {
		p.API.LogError(err.Error())
		return err
	}

	return nil
}

func (p *Plugin) PostToDo(info *BitbucketUserInfo) {
	// text, err := p.GetToDo(context.Background(), info.BitbucketUsername,
	// p.bitbucketConnect(*info.Token))
	// if err != nil {
	// 	mlog.Error(err.Error())
	// 	return
	// }
	//
	// p.CreateBotDMPost(info.UserID, text, "custom_git_todo")
}

func (p *Plugin) GetToDo(ctx context.Context, username string, bitbucketClient *github.Client) (string, error) {
	config := p.getConfiguration()

	issueResults, _, err := bitbucketClient.Search.Issues(ctx, getReviewSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	if err != nil {
		return "", err
	}

	notifications, _, err := bitbucketClient.Activity.ListNotifications(ctx, &github.NotificationListOptions{})
	if err != nil {
		return "", err
	}

	yourPrs, _, err := bitbucketClient.Search.Issues(ctx, getYourPrsSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	if err != nil {
		return "", err
	}

	yourAssignments, _, err := bitbucketClient.Search.Issues(ctx, getYourAssigneeSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
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
			message := fmt.Sprintf("[Vulnerability Alert for %v](%v)", n.GetRepository().GetFullName(), fixBitbucketNotificationSubjectURL(n.GetSubject().GetURL()))
			notificationContent += fmt.Sprintf("* %v\n", message)
		default:
			subjectURL := fixBitbucketNotificationSubjectURL(n.GetSubject().GetURL())
			notificationContent += fmt.Sprintf("* %v\n", subjectURL)
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

	configOrg := strings.TrimSpace(config.BitbucketOrg)
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

func (p *Plugin) getBaseURL() string {
	config := p.getConfiguration()
	if config.EnterpriseBaseURL != "" {
		return config.EnterpriseBaseURL
	}

	return "https://bitbucket.org/"
}

// getBitBucketAccountIDToMattermostUsernameMapping maps a BitBucket account ID to the corresponding Mattermost username, if any.
func (p *Plugin) getBitBucketAccountIDToMattermostUsernameMapping(bitbucketAccountID string) string {
	user, _ := p.API.GetUser(p.getBitbucketAccountIDToMattermostUserIDMapping(bitbucketAccountID))
	if user == nil {
		return ""
	}

	return user.Username
}
