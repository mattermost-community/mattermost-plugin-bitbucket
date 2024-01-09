package bitbucket_server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// TODO: This will be changed to create the main structs when modularized
type Link struct {
	Href string `json:"href"`
}

type BitbucketUser struct {
	AccountID int    `json:"id"`
	Username  string `json:"name"`
	Links     struct {
		Self []Link `json:"self"`
	} `json:"links"`
}

type BitbucketServerClient struct {
	BitbucketClient
}

func newServerClient(config ClientConfiguration) Client {
	return &BitbucketServerClient{
		BitbucketClient: BitbucketClient{
			ClientConfiguration: ClientConfiguration{
				SelfHostedURL:    config.SelfHostedURL,
				SelfHostedAPIURL: config.SelfHostedAPIURL,
				APIClient:        config.APIClient,
				OAuthClient:      config.OAuthClient,
				LogError:         config.LogError,
			},
		},
	}
}

func (c *BitbucketServerClient) getWhoAmI() (string, error) {
	requestURL := fmt.Sprintf("%s/plugins/servlet/applinks/whoami", c.SelfHostedURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "unable to create request for getting whoami identity")
	}

	resp, err := c.OAuthClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to make the request for getting whoami identity")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := errors.Errorf("who am i returned non-200 status code: %d", resp.StatusCode)
		return "", err
	}

	user, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to make the request for getting whoami identity")
	}

	return string(user), nil
}

func (c *BitbucketServerClient) GetMe() (*BitbucketUser, error) {
	username, err := c.getWhoAmI()
	if err != nil {
		c.LogError("failed to get whoami identity", "error", err.Error())
		return nil, err
	}

	resp, err := c.APIClient.DefaultApi.GetUser(username)
	if err != nil {
		c.LogError("failed to get user from bitbucket server", "error", err.Error())
		return nil, err
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		c.LogError("failed to marshaling user from bitbucket server", "error", err.Error())
		return nil, err
	}

	var user BitbucketUser
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		c.LogError("failed to parse user from bitbucket server", "error", err.Error())
		return nil, err
	}

	return &user, nil
}
