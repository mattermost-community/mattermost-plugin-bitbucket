package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	// "github.com/google/go-github/github"
	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/wbrefvem/go-bitbucket"

	"golang.org/x/oauth2"
)

const (
	API_ERROR_ID_NOT_CONNECTED = "not_connected"
	BITBUCKET_ICON_URL         = "https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png"
	BITBUCKET_USERNAME         = "Bitbucket Plugin"
)

type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writeAPIError(w http.ResponseWriter, err *APIErrorResponse) {
	b, _ := json.Marshal(err)
	w.WriteHeader(err.StatusCode)
	w.Write(b)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/webhook":
		fmt.Println("----- BB api.ServeHttp *** /webhook -----")
		p.handleWebhook(w, r)
	case "/oauth/connect":
		fmt.Println("----- BB api.ServeHttp *** /oauth/connect -----")
		p.connectUserToBitbucket(w, r)
	case "/oauth/complete":
		fmt.Println("----- BB *** api.ServeHttp *** /oauth/complete -----")
		p.completeConnectUserToBitbucket(w, r)
	case "/api/v1/connected":
		fmt.Println("----- BB *** api.ServeHttp *** /apt/v1/connected -----")
		p.getConnected(w, r)
	// case "/api/v1/todo":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/todo -----")
	// 	p.postToDo(w, r)
	// case "/api/v1/reviews":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/reviews -----")
	// 	p.getReviews(w, r)
	// case "/api/v1/yourprs":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/yourprs -----")
	// 	p.getYourPrs(w, r)
	// case "/api/v1/yourassignments":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/yourassignments -----")
	// 	p.getYourAssignments(w, r)
	// case "/api/v1/mentions":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/mentions -----")
	// 	p.getMentions(w, r)
	// case "/api/v1/unreads":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/unreads -----")
	// 	p.getUnreads(w, r)
	// case "/api/v1/settings":
	// 	fmt.Println("----- BB api.ServeHttp *** /apt/v1/settings -----")
	// 	p.updateSettings(w, r)
	case "/api/v1/user":
		fmt.Println("----- BB api.ServeHttp *** /apt/v1/user -----")
		p.getBitbucketUser(w, r)
	default:
		fmt.Println("----- BB api.ServeHttp *** default -----")
		http.NotFound(w, r)
	}
}

func (p *Plugin) connectUserToBitbucket(w http.ResponseWriter, r *http.Request) {

	// get MM userId from request
	userID := r.Header.Get("Mattermost-User-ID")
	fmt.Printf("---- connectUser  ->  userID = %+v\n", userID)
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	// get ClienID, ClientSecret, AuthURL and TokenURL endpoints
	conf := p.getOAuthConfig()

	// create state for KV store
	state := fmt.Sprintf("%v_%v", model.NewId()[0:15], userID)

	// save state in KV store
	p.API.KVSet(state, []byte(state))

	// get Auth Code and Link including ClientSecret and state
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

	fmt.Printf("---- connectUser  ->  conf = %+v\n", conf)
	fmt.Printf("---- connectUser  ->  state = %+v\n", state)
	fmt.Printf("---- connectUser  ->  url	 = %+v\n", url)

	// redirect to Authorization Code Link
	http.Redirect(w, r, url, http.StatusFound)
}

func (p *Plugin) completeConnectUserToBitbucket(w http.ResponseWriter, r *http.Request) {

	config := p.getConfiguration()
	//TODO

	ctx := context.Background()

	// get ClienID, ClientSecret, AuthURL and TokenURL endpoints
	conf := p.getOAuthConfig()

	fmt.Printf("---- completeConnectUser  ->  r.URL.Query() = %+v\n", r.URL.Query())
	fmt.Printf("---- completeConnectUser  ->  conf = %+v\n", conf)

	// get Authorization Code from url
	code := r.URL.Query().Get("code")
	fmt.Printf("---- completeConnectUser  ->  code = %+v\n", code)

	if len(code) == 0 {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	// get state from url
	state := r.URL.Query().Get("state")
	fmt.Printf("---- completeConnectUser  ->  state = %+v\n", state)

	// check stored state value is equal to return value from bitbucket callback
	if storedState, err := p.API.KVGet(state); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "missing stored state", http.StatusBadRequest)
		return
	} else if string(storedState) != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// get userID from state
	userID := strings.Split(state, "_")[1]
	fmt.Printf("---- completeConnectUser  ->  userID = %+v\n", userID)

	p.API.KVDelete(state)

	// converts auth code into a token.
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Printf("---- completeConnectUser  ->  tok = %+v\n", tok)

	var bitbucketClient *bitbucket.APIClient

	// connect to bitbucket API with authorization token
	bitbucketClient = p.bitbucketConnect(*tok)

	// get bitbucket user from Authorized bitbucket API request
	bitbucketUser, _, err := bitbucketClient.UsersApi.UserGet(ctx)
	fmt.Printf("---- completeConnectUser  ->  bitbucketUser = %+v\n", bitbucketUser)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := &BitbucketUserInfo{
		UserID: userID,
		Token:  tok,
		// BitbucketUsername: bitbucketUser.GetLogin(),
		BitbucketUsername: bitbucketUser.Username,
		LastToDoPostAt:    model.GetMillis(),
		Settings: &UserSettings{
			SidebarButtons: SETTING_BUTTONS_TEAM,
			DailyReminder:  true,
			Notifications:  true,
		},
		// AllowedPrivateRepos: config.EnablePrivateRepo,
	}

	fmt.Printf("---- completeConnectUser  ->  userInfo = %+v\n", userInfo)

	// Store User Info to Db
	if err := p.storeBitbucketUserInfo(userInfo); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Unable to connect user to Bitbucket", http.StatusInternalServerError)
		return
	}

	//TODO - need methods for getting Username and Link
	if err := p.storeBitbucketToUserIDMapping(bitbucketUser.Username, userID); err != nil {
		fmt.Println(err.Error())
	}

	// Post intro post
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
		"* The fourth tracks the number of unread messages you have.\n"+
		"* The fifth will refresh the numbers.\n\n"+
		"Click on them!\n\n"+
		"##### Slash Commands\n"+

		//TODO - need methods for getting Username and Link
		strings.Replace(COMMAND_HELP, "|", "`", -1), bitbucketUser.Username, bitbucketUser.Links.Html.Href)

	p.CreateBotDMPost(userID, message, "custom_git_welcome")

	p.API.PublishWebSocketEvent(
		WS_EVENT_CONNECT,
		map[string]interface{}{
			"connected":           true,
			"bitbucket_username":  userInfo.BitbucketUsername,
			"bitbucket_client_id": config.BitbucketOAuthClientID,
		},
		&model.WebsocketBroadcast{UserId: userID},
	)

	html := `
	<!DOCTYPE html>
	<html>
		<head>
			<script>
				window.close();
			</script>
		</head>
		<body>
			<p>Completed connecting to Bitbucket. Please close this window.</p>
		</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

type ConnectedResponse struct {
	Connected         bool          `json:"connected"`
	BitbucketUsername string        `json:"bitbucket_username"`
	BitbucketClientID string        `json:"bitbucket_client_id"`
	EnterpriseBaseURL string        `json:"enterprise_base_url,omitempty"`
	Organization      string        `json:"organization"`
	Settings          *UserSettings `json:"settings"`
}

type BitbucketUserRequest struct {
	UserID string `json:"user_id"`
}

type BitbucketUserResponse struct {
	Username string `json:"username"`
}

func (p *Plugin) getBitbucketUser(w http.ResponseWriter, r *http.Request) {
	requestorID := r.Header.Get("Mattermost-User-ID")
	if requestorID == "" {
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
		return
	}
	req := &BitbucketUserRequest{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil || req.UserID == "" {
		if err != nil {
			mlog.Error("Error decoding JSON body: " + err.Error())
		}
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object with a non-blank user_id field.", StatusCode: http.StatusBadRequest})
		return
	}

	userInfo, apiErr := p.getBitbucketUserInfo(req.UserID)
	if apiErr != nil {
		if apiErr.ID == API_ERROR_ID_NOT_CONNECTED {
			writeAPIError(w, &APIErrorResponse{ID: "", Message: "User is not connected to a Bitbucket account.", StatusCode: http.StatusNotFound})
		} else {
			writeAPIError(w, apiErr)
		}
		return
	}

	if userInfo == nil {
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "User is not connected to a Bitbucket account.", StatusCode: http.StatusNotFound})
		return
	}

	resp := &BitbucketUserResponse{Username: userInfo.BitbucketUsername}
	b, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		mlog.Error("Error encoding JSON response: " + jsonErr.Error())
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an unexpected error. Please try again.", StatusCode: http.StatusInternalServerError})
	}
	w.Write(b)
}

func (p *Plugin) getConnected(w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
		return
	}

	resp := &ConnectedResponse{
		Connected:         false,
		EnterpriseBaseURL: config.EnterpriseBaseURL,
		Organization:      config.BitbucketOrg,
	}

	info, _ := p.getBitbucketUserInfo(userID)
	if info != nil && info.Token != nil {
		resp.Connected = true
		resp.BitbucketUsername = info.BitbucketUsername
		resp.BitbucketClientID = config.BitbucketOAuthClientID
		resp.Settings = info.Settings

		if info.Settings.DailyReminder && r.URL.Query().Get("reminder") == "true" {
			lastPostAt := info.LastToDoPostAt

			var timezone *time.Location
			offset, _ := strconv.Atoi(r.Header.Get("X-Timezone-Offset"))
			timezone = time.FixedZone("local", -60*offset)

			// Post to do message if it's the next day and been more than an hour since the last post
			now := model.GetMillis()
			nt := time.Unix(now/1000, 0).In(timezone)
			lt := time.Unix(lastPostAt/1000, 0).In(timezone)
			if nt.Sub(lt).Hours() >= 1 && (nt.Day() != lt.Day() || nt.Month() != lt.Month() || nt.Year() != lt.Year()) {
				p.PostToDo(info)
				info.LastToDoPostAt = now
				p.storeBitbucketUserInfo(info)
			}
		}

		privateRepoStoreKey := info.UserID + BITBUCKET_PRIVATE_REPO_KEY
		if config.EnablePrivateRepo && !info.AllowedPrivateRepos {
			hasBeenNotified := false
			if val, err := p.API.KVGet(privateRepoStoreKey); err == nil {
				hasBeenNotified = val != nil
			} else {
				mlog.Error("Unable to get private repo key value, err=" + err.Error())
			}

			if !hasBeenNotified {
				p.CreateBotDMPost(info.UserID, "Private repositories have been enabled for this plugin. To be able to use them you must disconnect and reconnect your Bitbucket account. To reconnect your account, use the following slash commands: `/bitbucket disconnect` followed by `/bitbucket connect`.", "")
				if err := p.API.KVSet(privateRepoStoreKey, []byte("1")); err != nil {
					mlog.Error("Unable to set private repo key value, err=" + err.Error())
				}
			}
		}
	}

	b, _ := json.Marshal(resp)
	w.Write(b)
}

func (p *Plugin) getMentions(w http.ResponseWriter, r *http.Request) {
	// config := p.getConfiguration()
	//
	// userID := r.Header.Get("Mattermost-User-ID")
	// if userID == "" {
	// 	http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 	return
	// }
	//
	// ctx := context.Background()
	//
	// var bitbucketClient *github.Client
	// username := ""
	//
	// if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 	writeAPIError(w, err)
	// 	return
	// } else {
	// 	bitbucketClient = p.bitbucketConnect(*info.Token)
	// 	username = info.BitbucketUsername
	// }
	//
	// result, _, err := bitbucketClient.Search.Issues(ctx,
	// getMentionSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	// if err != nil {
	// 	mlog.Error(err.Error())
	// }
	//
	// resp, _ := json.Marshal(result.Issues)
	// w.Write(resp)
}

func (p *Plugin) getUnreads(w http.ResponseWriter, r *http.Request) {
	// userID := r.Header.Get("Mattermost-User-ID")
	// if userID == "" {
	// 	http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 	return
	// }
	//
	// ctx := context.Background()
	//
	// var bitbucketClient *github.Client
	//
	// if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 	writeAPIError(w, err)
	// 	return
	// } else {
	// 	bitbucketClient = p.bitbucketConnect(*info.Token)
	// }
	//
	// notifications, _, err := bitbucketClient.Activity.ListNotifications(ctx, &github.NotificationListOptions{})
	// if err != nil {
	// 	mlog.Error(err.Error())
	// }
	//
	// filteredNotifications := []*github.Notification{}
	// for _, n := range notifications {
	// 	if n.GetReason() == "subscribed" {
	// 		continue
	// 	}
	//
	// 	if p.checkOrg(n.GetRepository().GetOwner().GetLogin()) != nil {
	// 		continue
	// 	}
	//
	// 	filteredNotifications = append(filteredNotifications, n)
	// }
	//
	// resp, _ := json.Marshal(filteredNotifications)
	// w.Write(resp)
}

func (p *Plugin) getReviews(w http.ResponseWriter, r *http.Request) {
	// config := p.getConfiguration()
	//
	// userID := r.Header.Get("Mattermost-User-ID")
	// if userID == "" {
	// 	http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 	return
	// }
	//
	// ctx := context.Background()
	//
	// var bitbucketClient *github.Client
	// username := ""
	//
	// if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 	writeAPIError(w, err)
	// 	return
	// } else {
	// 	bitbucketClient = p.bitbucketConnect(*info.Token)
	// 	username = info.BitbucketUsername
	// }
	//
	// result, _, err := bitbucketClient.Search.Issues(ctx,
	// getReviewSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	// if err != nil {
	// 	mlog.Error(err.Error())
	// }
	//
	// resp, _ := json.Marshal(result.Issues)
	// w.Write(resp)
}

func (p *Plugin) getYourPrs(w http.ResponseWriter, r *http.Request) {
	// config := p.getConfiguration()
	//
	// userID := r.Header.Get("Mattermost-User-ID")
	// if userID == "" {
	// 	http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 	return
	// }
	//
	// ctx := context.Background()
	//
	// var bitbucketClient *github.Client
	// username := ""
	//
	// if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 	writeAPIError(w, err)
	// 	return
	// } else {
	// 	bitbucketClient = p.bitbucketConnect(*info.Token)
	// 	username = info.BitbucketUsername
	// }
	//
	// result, _, err := bitbucketClient.Search.Issues(ctx,
	// getYourPrsSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	// if err != nil {
	// 	mlog.Error(err.Error())
	// }
	//
	// resp, _ := json.Marshal(result.Issues)
	// w.Write(resp)
}

func (p *Plugin) getYourAssignments(w http.ResponseWriter, r *http.Request) {
	// 	config := p.getConfiguration()
	//
	// 	userID := r.Header.Get("Mattermost-User-ID")
	// 	if userID == "" {
	// 		http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 		return
	// 	}
	//
	// 	ctx := context.Background()
	//
	// 	var bitbucketClient *github.Client
	// 	username := ""
	//
	// 	if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 		writeAPIError(w, err)
	// 		return
	// 	} else {
	// 		bitbucketClient = p.bitbucketConnect(*info.Token)
	// 		username = info.BitbucketUsername
	// 	}
	//
	// 	result, _, err := bitbucketClient.Search.Issues(ctx,
	// 	getYourAssigneeSearchQuery(username, config.BitbucketOrg), &github.SearchOptions{})
	// 	if err != nil {
	// 		mlog.Error(err.Error())
	// 	}
	//
	// 	resp, _ := json.Marshal(result.Issues)
	// 	w.Write(resp)
	// }
	//
	// func (p *Plugin) postToDo(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("")
	// 	userID := r.Header.Get("Mattermost-User-ID")
	// 	if userID == "" {
	// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
	// 		return
	// 	}
	//
	// 	var bitbucketClient *github.Client
	// 	username := ""
	//
	// 	if info, err := p.getBitbucketUserInfo(userID); err != nil {
	// 		writeAPIError(w, err)
	// 		return
	// 	} else {
	// 		bitbucketClient = p.bitbucketConnect(*info.Token)
	// 		username = info.BitbucketUsername
	// 	}
	//
	// 	text, err := p.GetToDo(context.Background(), username, bitbucketClient)
	// 	if err != nil {
	// 		mlog.Error(err.Error())
	// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an error getting the to do items.", StatusCode: http.StatusUnauthorized})
	// 		return
	// 	}
	//
	// 	if err := p.CreateBotDMPost(userID, text, "custom_git_todo"); err != nil {
	// 		writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an error posting the to do items.", StatusCode: http.StatusUnauthorized})
	// 	}
	//
	// 	w.Write([]byte("{\"status\": \"OK\"}"))
}

func (p *Plugin) updateSettings(w http.ResponseWriter, r *http.Request) {
	// userID := r.Header.Get("Mattermost-User-ID")
	// if userID == "" {
	// 	http.Error(w, "Not authorized", http.StatusUnauthorized)
	// 	return
	// }
	//
	// var settings *UserSettings
	// json.NewDecoder(r.Body).Decode(&settings)
	// if settings == nil {
	// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
	// 	return
	// }
	//
	// info, err := p.getBitbucketUserInfo(userID)
	// if err != nil {
	// 	writeAPIError(w, err)
	// 	return
	// }
	//
	// info.Settings = settings
	//
	// if err := p.storeBitbucketUserInfo(info); err != nil {
	// 	mlog.Error(err.Error())
	// 	http.Error(w, "Encountered error updating settings", http.StatusInternalServerError)
	// }
	//
	// resp, _ := json.Marshal(info.Settings)
	// w.Write(resp)
}
