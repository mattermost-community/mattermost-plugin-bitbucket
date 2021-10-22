package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/wbrefvem/go-bitbucket"

	"golang.org/x/oauth2"
)

const (
	APIErrorIDNotConnected = "not_connected"
	// TokenTTL is the OAuth token expiry duration in seconds
	TokenTTL = 10 * 60
)

// OAuthState is struct where OAuth state is stored
type OAuthState struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

// APIErrorResponse is object with error information that is sent to webapp
type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e *APIErrorResponse) Error() string {
	return e.Message
}

type PRDetails struct {
	URL          string                  `json:"url"`
	ID           int                     `json:"id"`
	Participants []bitbucket.Participant `json:"participants"`
}

// HTTPHandlerFuncWithUser is http.HandleFunc but userID is already exported
type HTTPHandlerFuncWithUser func(w http.ResponseWriter, r *http.Request, userID string)

// ResponseType indicates type of response returned by api
type ResponseType string

const (
	// ResponseTypeJSON indicates that response type is json
	ResponseTypeJSON ResponseType = "JSON_RESPONSE"
	// ResponseTypePlain indicates that response type is text plain
	ResponseTypePlain ResponseType = "TEXT_RESPONSE"
)

func (p *Plugin) writeJSON(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		p.API.LogWarn("Failed to marshal JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		p.API.LogWarn("Failed to write JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) writeAPIError(w http.ResponseWriter, apiErr *APIErrorResponse) {
	b, err := json.Marshal(apiErr)
	if err != nil {
		p.API.LogWarn("Failed to marshal API error", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(apiErr.StatusCode)

	_, err = w.Write(b)
	if err != nil {
		p.API.LogWarn("Failed to write JSON response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) initializeAPI() {
	p.router = mux.NewRouter()

	oauthRouter := p.router.PathPrefix("/oauth").Subrouter()
	apiRouter := p.router.PathPrefix("/api/v1").Subrouter()

	p.router.HandleFunc("/webhook", p.handleWebhook).Methods(http.MethodPost)

	oauthRouter.HandleFunc("/connect", p.extractUserMiddleWare(p.connectUserToBitbucket, ResponseTypePlain)).Methods(http.MethodGet)
	oauthRouter.HandleFunc("/complete", p.extractUserMiddleWare(p.completeConnectUserToBitbucket, ResponseTypePlain)).Methods(http.MethodGet)

	apiRouter.HandleFunc("/connected", p.getConnected).Methods(http.MethodGet)
	apiRouter.HandleFunc("/todo", p.extractUserMiddleWare(p.postToDo, ResponseTypeJSON)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/reviews", p.extractUserMiddleWare(p.getReviews, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/yourprs", p.extractUserMiddleWare(p.getYourPrs, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/prsdetails", p.extractUserMiddleWare(p.getPrsDetails, ResponseTypePlain)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/searchissues", p.extractUserMiddleWare(p.searchIssues, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/yourassignments", p.extractUserMiddleWare(p.getYourAssignments, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/createissue", p.extractUserMiddleWare(p.createIssue, ResponseTypePlain)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/createissuecomment", p.extractUserMiddleWare(p.createIssueComment, ResponseTypePlain)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/repositories", p.extractUserMiddleWare(p.getRepositories, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/settings", p.extractUserMiddleWare(p.updateSettings, ResponseTypePlain)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/user", p.extractUserMiddleWare(p.getBitbucketUser, ResponseTypeJSON)).Methods(http.MethodPost)
	apiRouter.HandleFunc("/issue", p.extractUserMiddleWare(p.getIssueByID, ResponseTypePlain)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/pr", p.extractUserMiddleWare(p.getPrByID, ResponseTypePlain)).Methods(http.MethodGet)

	apiRouter.HandleFunc("/config", checkPluginRequest(p.getConfig)).Methods(http.MethodGet)
	apiRouter.HandleFunc("/token", checkPluginRequest(p.getToken)).Methods(http.MethodGet)
}

func (p *Plugin) extractUserMiddleWare(handler HTTPHandlerFuncWithUser, responseType ResponseType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID == "" {
			switch responseType {
			case ResponseTypeJSON:
				p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Not authorized.", StatusCode: http.StatusUnauthorized})
			case ResponseTypePlain:
				http.Error(w, "Not authorized", http.StatusUnauthorized)
			default:
				p.API.LogError("Unknown ResponseType detected")
			}
			return
		}

		handler(w, r, userID)
	}
}

func checkPluginRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// All other plugins are allowed
		pluginID := r.Header.Get("Mattermost-Plugin-ID")
		if pluginID == "" {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
		return
	}

	r.Header.Set("Mattermost-Plugin-ID", c.SourcePluginId)
	w.Header().Set("Content-Type", "application/json")

	p.router.ServeHTTP(w, r)
}

func (p *Plugin) connectUserToBitbucket(w http.ResponseWriter, r *http.Request, userID string) {
	conf := p.getOAuthConfig()

	state := OAuthState{
		UserID: userID,
		Token:  model.NewId()[:15],
	}

	stateBytes, err := json.Marshal(state)
	if err != nil {
		http.Error(w, "json marshal failed", http.StatusInternalServerError)
		return
	}

	appErr := p.API.KVSetWithExpiry(state.Token, stateBytes, TokenTTL)
	if appErr != nil {
		http.Error(w, "error setting stored state", http.StatusBadRequest)
		return
	}

	url := conf.AuthCodeURL(state.Token, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusFound)
}

func (p *Plugin) completeConnectUserToBitbucket(w http.ResponseWriter, r *http.Request, authedUserID string) {
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	stateToken := r.URL.Query().Get("state")

	storedState, appErr := p.API.KVGet(stateToken)
	if appErr != nil {
		p.API.LogError("Missing stored state", "appErr", appErr.Error())
		http.Error(w, "missing stored state", http.StatusBadRequest)
		return
	}
	appErr = p.API.KVDelete(stateToken)
	if appErr != nil {
		p.API.LogError("Error deleting stored state", "appErr", appErr.Error())
		http.Error(w, "error deleting stored state", http.StatusBadRequest)
		return
	}

	var state OAuthState
	if err := json.Unmarshal(storedState, &state); err != nil {
		http.Error(w, "json unmarshal failed", http.StatusBadRequest)
		return
	}

	if state.Token != stateToken {
		http.Error(w, "invalid state token", http.StatusBadRequest)
		return
	}
	if state.UserID != authedUserID {
		http.Error(w, "Not authorized, incorrect user", http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	conf := p.getOAuthConfig()

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		p.API.LogError("Error while converting authorization code into token", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bitbucketClient := p.bitbucketConnect(*tok)
	bitbucketUser, httpResponse, err := bitbucketClient.UsersApi.UserGet(ctx)
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}
	if err != nil {
		p.API.LogError("Error converting authorization code int token", "err", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo := &BitbucketUserInfo{
		UserID:             state.UserID,
		Token:              tok,
		BitbucketUsername:  bitbucketUser.Username,
		BitbucketAccountID: bitbucketUser.AccountId,
		LastToDoPostAt:     model.GetMillis(),
		Settings: &UserSettings{
			SidebarButtons: SettingButtonsTeam,
			DailyReminder:  true,
			Notifications:  true,
		},
	}

	if err = p.storeBitbucketUserInfo(userInfo); err != nil {
		p.API.LogError("Error connecting user to Bitbucket", "err", err.Error())
		http.Error(w, "Unable to connect user to Bitbucket", http.StatusInternalServerError)
		return
	}

	if err = p.storeBitbucketAccountIDToMattermostUserIDMapping(bitbucketUser.AccountId, state.UserID); err != nil {
		p.API.LogError("Error storing Bitbucket account ID to Mattermost user ID mapping", "err", err.Error())
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
		"* The fourth will refresh the numbers.\n\n"+
		"Click on them!\n\n"+
		"##### Slash Commands\n"+
		strings.ReplaceAll(commandHelp, "|", "`"), bitbucketUser.Username, bitbucketUser.Links.Html.Href)

	p.CreateBotDMPost(state.UserID, message, "custom_bitbucket_welcome")

	config := p.getConfiguration()

	p.API.PublishWebSocketEvent(
		WsEventConnect,
		map[string]interface{}{
			"connected":           true,
			"bitbucket_username":  userInfo.BitbucketUsername,
			"bitbucket_client_id": config.BitbucketOAuthClientID,
			"organization":        config.BitbucketOrg,
		},
		&model.WebsocketBroadcast{UserId: state.UserID},
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
	_, err = w.Write([]byte(html))
	if err != nil {
		p.API.LogWarn("Failed to write HTML response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p *Plugin) getBitbucketUser(w http.ResponseWriter, r *http.Request, _ string) {
	type BitbucketUserRequest struct {
		UserID string `json:"user_id"`
	}

	req := &BitbucketUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.API.LogError("Error decoding BitbucketUserRequest from JSON body", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object.", StatusCode: http.StatusBadRequest})
		return
	}

	if req.UserID == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object with a non-blank user_id field.", StatusCode: http.StatusBadRequest})
		return
	}

	userInfo, apiErr := p.getBitbucketUserInfo(req.UserID)
	if apiErr != nil {
		if apiErr.ID == APIErrorIDNotConnected {
			p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "User is not connected to a Bitbucket account.", StatusCode: http.StatusNotFound})
		} else {
			p.writeAPIError(w, apiErr)
		}
		return
	}

	if userInfo == nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "User is not connected to a BitBucket account.", StatusCode: http.StatusNotFound})
		return
	}

	type BitbucketUserResponse struct {
		Username string `json:"username"`
	}

	resp := &BitbucketUserResponse{Username: userInfo.BitbucketUsername}
	p.writeJSON(w, resp)
}

func (p *Plugin) getConnected(w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()

	type ConnectedResponse struct {
		Connected         bool          `json:"connected"`
		BitbucketUsername string        `json:"bitbucket_username"`
		BitbucketClientID string        `json:"bitbucket_client_id"`
		Organization      string        `json:"organization"`
		Settings          *UserSettings `json:"settings"`
	}

	resp := &ConnectedResponse{
		Connected:    false,
		Organization: config.BitbucketOrg,
	}

	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		p.writeJSON(w, resp)
		return
	}

	info, _ := p.getBitbucketUserInfo(userID)
	if info == nil || info.Token == nil {
		p.writeJSON(w, resp)
		return
	}

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
			if p.HasUnreads(info) {
				p.PostToDo(info)
				info.LastToDoPostAt = now
				if err := p.storeBitbucketUserInfo(info); err != nil {
					p.API.LogWarn("Failed to store bitbucket info for new user", "userID", userID, "error", err.Error())
				}
			}
		}
	}

	p.writeJSON(w, resp)
}

func (p *Plugin) getReviews(w http.ResponseWriter, _ *http.Request, userID string) {
	userInfo, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*userInfo.Token)

	userRepos, err := p.getUserRepositories(context.Background(), bitbucketClient)
	if err != nil {
		p.API.LogError("Error occurred while searching for repositories", "err", err.Error())
		return
	}

	yourPrs, err := p.getAssignedPRs(context.Background(), userInfo, bitbucketClient, userRepos)
	if err != nil {
		p.API.LogError("Error occurred while searching for pull requests", "err", err.Error())
		return
	}

	p.writeJSON(w, yourPrs)
}

func (p *Plugin) getYourPrs(w http.ResponseWriter, _ *http.Request, userID string) {
	userInfo, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*userInfo.Token)

	userRepos, err := p.getUserRepositories(context.Background(), bitbucketClient)
	if err != nil {
		p.API.LogError("error occurred while searching for repositories", "err", err)
		return
	}

	openPRs, err := p.getOpenPRs(context.Background(), userInfo, bitbucketClient, userRepos)
	if err != nil {
		p.API.LogError("error occurred while searching for pull requests", "err", err)
		return
	}

	p.writeJSON(w, openPRs)
}

func (p *Plugin) getPrsDetails(w http.ResponseWriter, r *http.Request, userID string) {
	info, err := p.getBitbucketUserInfo(userID)
	if err != nil {
		p.writeAPIError(w, err)
		return
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	var prList []*PRDetails
	if err := json.NewDecoder(r.Body).Decode(&prList); err != nil {
		p.API.LogError("Error decoding PRDetails JSON body", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object.", StatusCode: http.StatusBadRequest})
		return
	}

	prDetails := make([]*PRDetails, len(prList))
	ctx := context.Background()
	var wg sync.WaitGroup
	for i, pr := range prList {
		i := i
		pr := pr
		wg.Add(1)
		go func() {
			defer wg.Done()
			prDetail := p.fetchPRDetails(ctx, bitbucketClient, pr.URL, pr.ID)
			prDetails[i] = prDetail
		}()
	}

	wg.Wait()

	p.writeJSON(w, prDetails)
}

func (p *Plugin) fetchPRDetails(ctx context.Context, client *bitbucket.APIClient, prURL string, prID int) *PRDetails {
	repoOwner, repoName := getRepoOwnerAndNameFromURL(prURL)

	prInfo, httpResponse, err := client.PullrequestsApi.RepositoriesUsernameRepoSlugPullrequestsPullRequestIdGet(ctx, repoOwner, repoName, int32(prID))
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}
	if err != nil {
		p.API.LogError("Error fetching pull request", "prID", prID, "err", err.Error())
		return nil
	}

	return &PRDetails{
		URL:          prURL,
		ID:           prID,
		Participants: prInfo.Participants,
	}
}

func (p *Plugin) searchIssues(w http.ResponseWriter, r *http.Request, userID string) {
	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	searchTerm := r.FormValue("term")

	result, err := p.getIssuesWithTerm(bitbucketClient, searchTerm)
	if err != nil {
		p.API.LogError("Error fetching issues with term", "searchTerm", searchTerm, "err", err.Error())
		return
	}

	p.writeJSON(w, result)
}

func (p *Plugin) getPermaLink(postID string) string {
	siteURL := *p.API.GetConfig().ServiceSettings.SiteURL
	return fmt.Sprintf("%v/_redirect/pl/%v", siteURL, postID)
}

func getFailReason(code int, repo string, username string) string {
	cause := ""
	switch code {
	case http.StatusInternalServerError:
		cause = "Internal server error"
	case http.StatusBadRequest:
		cause = "Bad request"
	case http.StatusNotFound:
		cause = fmt.Sprintf("Sorry, either you don't have access to the repo %s with the user %s or it is no longer available", repo, username)
	case http.StatusUnauthorized:
		cause = fmt.Sprintf("Sorry, your user %s is unauthorized to do this action", username)
	case http.StatusForbidden:
		cause = fmt.Sprintf("Sorry, you don't have enough permissions to comment in the repo %s with the user %s", repo, username)
	default:
		cause = fmt.Sprintf("Unknown status code %d", code)
	}
	return cause
}

func (p *Plugin) createIssueComment(w http.ResponseWriter, r *http.Request, userID string) {
	type CreateIssueCommentRequest struct {
		PostID  string `json:"post_id"`
		Owner   string `json:"owner"`
		Repo    string `json:"repo"`
		Number  int    `json:"number"`
		Comment string `json:"comment"`
	}

	req := &CreateIssueCommentRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.API.LogError("Error decoding CreateIssueCommentRequest JSON body", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object.", StatusCode: http.StatusBadRequest})
		return
	}

	if req.PostID == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid post id", StatusCode: http.StatusBadRequest})
		return
	}

	if req.Owner == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid repo owner.", StatusCode: http.StatusBadRequest})
		return
	}

	if req.Repo == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid repo.", StatusCode: http.StatusBadRequest})
		return
	}

	if req.Number == 0 {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid issue number.", StatusCode: http.StatusBadRequest})
		return
	}

	if req.Comment == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid non empty comment.", StatusCode: http.StatusBadRequest})
		return
	}

	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	post, appErr := p.API.GetPost(req.PostID)
	if appErr != nil {
		p.API.LogError("failed to load post", "postID", req.PostID)
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to load post " + req.PostID, StatusCode: http.StatusInternalServerError})
		return
	}
	if post == nil {
		p.API.LogError("failed to load post: not found", "postID", req.PostID)
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to load post " + req.PostID + ": not found", StatusCode: http.StatusNotFound})
		return
	}

	commentUsername, err := p.getUsername(post.UserId)
	if err != nil {
		p.API.LogError("failed to load post", "UserId", post.UserId)
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to get username", StatusCode: http.StatusInternalServerError})
		return
	}

	currentUsername := info.BitbucketUsername
	permalink := p.getPermaLink(req.PostID)
	permalinkMessage := fmt.Sprintf("*@%s attached a* [message](%s) *from %s*\n\n", currentUsername, permalink, commentUsername)

	req.Comment = permalinkMessage + req.Comment
	comment := bitbucket.IssueComment{}
	comment.Content = &bitbucket.IssueContent{
		Raw: req.Comment,
	}

	httpResponse, err := bitbucketClient.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesIssueIdCommentsPost(context.Background(), strconv.Itoa(req.Number), req.Owner, req.Repo, comment)
	if err != nil {
		if httpResponse != nil {
			p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to create an issue comment" + getFailReason(httpResponse.StatusCode, req.Repo, currentUsername), StatusCode: httpResponse.StatusCode})
			_ = httpResponse.Body.Close()
		} else {
			p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to create an issue comment"})
		}
		return
	}
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}

	locationURL := httpResponse.Header.Get("Location")
	splittedLocationURL := strings.Split(locationURL, "/")
	commentID := splittedLocationURL[len(splittedLocationURL)-1]

	issueComment, httpResponse, err := bitbucketClient.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesIssueIdCommentsCommentIdGet(context.Background(), commentID, req.Owner, req.Repo, strconv.Itoa(req.Number))
	if err != nil {
		if httpResponse != nil {
			p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to fetch the newly created comment: " + getFailReason(httpResponse.StatusCode, req.Repo, currentUsername), StatusCode: httpResponse.StatusCode})
			_ = httpResponse.Body.Close()
		} else {
			p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to fetch the newly created comment"})
		}
		return
	}
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}

	permalinkReplyMessage := fmt.Sprintf("[Message](%v) attached to Bitbucket issue [#%v](%v)", permalink, req.Number, issueComment.Links.Html.Href)
	reply := &model.Post{
		Message:   permalinkReplyMessage,
		ChannelId: post.ChannelId,
		UserId:    userID,
	}

	_, appErr = p.API.CreatePost(reply)
	if appErr != nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to create notification post " + req.PostID, StatusCode: http.StatusInternalServerError})
		return
	}

	p.writeJSON(w, issueComment)
}

func (p *Plugin) getYourAssignments(w http.ResponseWriter, _ *http.Request, userID string) {
	userInfo, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*userInfo.Token)

	userRepos, err := p.getUserRepositories(context.Background(), bitbucketClient)
	if err != nil {
		p.API.LogError("Error occurred while searching for repositories", "err", err)
		return
	}

	yourAssignments, err := p.getAssignedIssues(context.Background(), userInfo, bitbucketClient, userRepos)
	if err != nil {
		p.API.LogError("Error occurred while searching assigned issues", "err", err)
		return
	}

	p.writeJSON(w, yourAssignments)
}

func (p *Plugin) postToDo(w http.ResponseWriter, _ *http.Request, userID string) {
	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	text, err := p.GetToDo(context.Background(), info, bitbucketClient)
	if err != nil {
		p.API.LogError("Error fetching todos", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Encountered an error getting the to do items.", StatusCode: http.StatusUnauthorized})
		return
	}

	p.CreateBotDMPost(userID, text, "custom_bitbucket_todo")

	resp := struct {
		Status string
	}{"OK"}

	p.writeJSON(w, resp)
}

func (p *Plugin) updateSettings(w http.ResponseWriter, r *http.Request, userID string) {
	var settings *UserSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		p.API.LogError("Error decoding settings from JSON body", "err", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if settings == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	info, err := p.getBitbucketUserInfo(userID)
	if err != nil {
		p.writeAPIError(w, err)
		return
	}

	info.Settings = settings

	if err := p.storeBitbucketUserInfo(info); err != nil {
		p.API.LogError("Error updating settings", "err", err.Error())
		http.Error(w, "Encountered error updating settings", http.StatusInternalServerError)
		return
	}

	p.writeJSON(w, info.Settings)
}

func (p *Plugin) getIssueByID(w http.ResponseWriter, r *http.Request, userID string) {
	owner := r.FormValue("owner")
	repo := r.FormValue("repo")
	issueID := r.FormValue("id")
	_, err := strconv.Atoi(issueID)
	if err != nil {
		p.writeAPIError(w, &APIErrorResponse{Message: "Invalid param 'id'.", StatusCode: http.StatusBadRequest})
		return
	}

	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}
	bitbucketClient := p.bitbucketConnect(*info.Token)

	result, httpResponse, err := bitbucketClient.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesIssueIdGet(context.Background(), owner, issueID, repo)
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}
	if err != nil {
		p.API.LogDebug("Could not get issue", "owner", owner, "repo", repo, "number", issueID, "error", err.Error())
		p.writeAPIError(w, &APIErrorResponse{Message: "Could not get issue", StatusCode: http.StatusInternalServerError})
		return
	}

	p.writeJSON(w, result)
}

func (p *Plugin) getPrByID(w http.ResponseWriter, r *http.Request, userID string) {
	owner := r.FormValue("owner")
	repo := r.FormValue("repo")
	prID := r.FormValue("id")

	prIDInt, err := strconv.ParseInt(prID, 10, 64)
	if err != nil {
		p.writeAPIError(w, &APIErrorResponse{Message: "Invalid param 'id'.", StatusCode: http.StatusBadRequest})
		return
	}

	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}
	bitbucketClient := p.bitbucketConnect(*info.Token)

	result, httpResponse, err := bitbucketClient.PullrequestsApi.RepositoriesUsernameRepoSlugPullrequestsPullRequestIdGet(context.Background(), owner, repo, int32(prIDInt))
	if httpResponse != nil {
		_ = httpResponse.Body.Close()
	}
	if err != nil {
		p.API.LogDebug("Could not get pull request", "owner", owner, "repo", repo, "ID", prID, "error", err.Error())
		p.writeAPIError(w, &APIErrorResponse{Message: "Could not get pull request", StatusCode: http.StatusInternalServerError})
		return
	}

	p.writeJSON(w, result)
}

func (p *Plugin) getRepositories(w http.ResponseWriter, _ *http.Request, userID string) {
	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	bitbucketClient := p.bitbucketConnect(*info.Token)

	ctx := context.Background()

	repos, err := p.getUserRepositories(ctx, bitbucketClient)
	if err != nil {
		p.API.LogError("Failed to fetch repositories", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{Message: "Failed to fetch repositories", StatusCode: http.StatusInternalServerError})
		return
	}

	// Only send down fields to client that are needed
	type RepositoryResponse struct {
		Name     string `json:"name,omitempty"`
		FullName string `json:"full_name,omitempty"`
	}

	resp := make([]RepositoryResponse, len(repos))
	for i, r := range repos {
		resp[i].Name = r.Name
		resp[i].FullName = r.FullName
	}

	p.writeJSON(w, resp)
}

func (p *Plugin) createIssue(w http.ResponseWriter, r *http.Request, userID string) {
	type IssueRequest struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		Repo   string `json:"repo"`
		PostID string `json:"post_id"`
	}

	// get data for the issue from the request body and fill IssueRequest object
	issue := &IssueRequest{}
	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
		p.API.LogError("Error decoding JSON body", "err", err.Error())
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a JSON object.", StatusCode: http.StatusBadRequest})
		return
	}

	if issue.Title == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid issue title.", StatusCode: http.StatusBadRequest})
		return
	}

	if issue.Repo == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a valid repo name.", StatusCode: http.StatusBadRequest})
		return
	}

	if issue.PostID == "" {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "Please provide a postID", StatusCode: http.StatusBadRequest})
		return
	}

	// Make sure user has a connected Bitbucket account
	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		p.writeAPIError(w, apiErr)
		return
	}

	post, appErr := p.API.GetPost(issue.PostID)
	if appErr != nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to load post " + issue.PostID, StatusCode: http.StatusInternalServerError})
		return
	}
	if post == nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to load post " + issue.PostID + ": not found", StatusCode: http.StatusNotFound})
		return
	}

	username, err := p.getUsername(post.UserId)
	if err != nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to get username", StatusCode: http.StatusInternalServerError})
		return
	}

	bbIssue := bitbucket.Issue{Title: issue.Title}
	bbIssue.Content = &bitbucket.IssueContent{}
	bbIssue.Content.Raw = issue.Body

	permalink := p.getPermaLink(issue.PostID)

	mmMessage := fmt.Sprintf("_Issue created from a [Mattermost message](%v) *by %s*._", permalink, username)

	if bbIssue.Content.Raw != "" {
		mmMessage = "\n\n" + mmMessage
	}
	bbIssue.Content.Raw += mmMessage

	currentUser, appErr := p.API.GetUser(userID)
	if appErr != nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to load current user", StatusCode: http.StatusInternalServerError})
		return
	}

	splittedRepo := strings.Split(issue.Repo, "/")
	owner := splittedRepo[0]
	repoName := splittedRepo[1]
	bitbucketClient := p.bitbucketConnect(*info.Token)
	issuePostResult, issuePostResponse, err := bitbucketClient.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesPost(context.Background(), owner, repoName, bbIssue)
	if err != nil {
		if issuePostResponse != nil {
			p.writeAPIError(w,
				&APIErrorResponse{
					ID: "",
					Message: "failed to create issue: " + getFailReason(issuePostResponse.StatusCode,
						issue.Repo,
						currentUser.Username),
					StatusCode: issuePostResponse.StatusCode,
				})
			_ = issuePostResponse.Body.Close()
		} else {
			p.writeAPIError(w,
				&APIErrorResponse{
					ID:      "",
					Message: "failed to create issue: " + err.Error(),
				})
		}
		return
	}
	if issuePostResponse != nil {
		_ = issuePostResponse.Body.Close()
	}

	rootID := issue.PostID
	if post.RootId != "" {
		rootID = post.RootId
	}

	issueGetResult, issueGetResponse, err := bitbucketClient.IssueTrackerApi.RepositoriesUsernameRepoSlugIssuesIssueIdGet(context.Background(),
		owner, fmt.Sprint(issuePostResult.Id), repoName)
	if err != nil {
		var statusCode int
		if issueGetResponse != nil {
			statusCode = issueGetResponse.StatusCode
			_ = issueGetResponse.Body.Close()
		}
		p.writeAPIError(w,
			&APIErrorResponse{
				ID:         "",
				Message:    "issue created, but an error occurred while fetching a link: " + err.Error(),
				StatusCode: statusCode,
			},
		)
		return
	}
	if issueGetResponse != nil {
		_ = issueGetResponse.Body.Close()
	}

	message := fmt.Sprintf("Created Bitbucket issue [#%v](%v) from a [message](%s)", issuePostResult.Id, issueGetResult.Links.Html.Href, permalink)
	reply := &model.Post{
		Message:   message,
		ChannelId: post.ChannelId,
		RootId:    rootID,
		ParentId:  rootID,
		UserId:    userID,
	}

	_, appErr = p.API.CreatePost(reply)
	if appErr != nil {
		p.writeAPIError(w, &APIErrorResponse{ID: "", Message: "failed to create notification post " + issue.PostID, StatusCode: http.StatusInternalServerError})
		return
	}

	p.writeJSON(w, issueGetResult)
}

func (p *Plugin) getConfig(w http.ResponseWriter, _ *http.Request) {
	config := p.getConfiguration()

	p.writeJSON(w, config)
}

func (p *Plugin) getToken(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	if userID == "" {
		http.Error(w, "please provide a userID", http.StatusBadRequest)
		return
	}

	info, apiErr := p.getBitbucketUserInfo(userID)
	if apiErr != nil {
		http.Error(w, apiErr.Error(), apiErr.StatusCode)
		return
	}

	p.writeJSON(w, info.Token)
}

func getRepoOwnerAndNameFromURL(url string) (string, string) {
	splitted := strings.Split(url, "/")
	return splitted[len(splitted)-2], splitted[len(splitted)-1]
}
