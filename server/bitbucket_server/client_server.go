package bitbucket_server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
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

func newServerClient(selfHostedURL string, selfHostedAPIURL string, apiClient *bitbucketv1.APIClient) Client {
	return &BitbucketServerClient{
		BitbucketClient: BitbucketClient{
			apiClient:        apiClient,
			selfHostedURL:    selfHostedURL,
			selfHostedAPIURL: selfHostedAPIURL,
		},
	}
}

func (c *BitbucketServerClient) getWhoAmI(accessToken string) (string, error) {
	requestURL := fmt.Sprintf("%s/plugins/servlet/applinks/whoami", c.selfHostedURL)

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.Errorf("who am i returned non-200 status code: %d", resp.StatusCode)
	}

	user, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response")
	}

	return string(user), nil
}

func (c *BitbucketServerClient) GetMe(accessToken string) (*BitbucketUser, error) {
	username, err := c.getWhoAmI(accessToken)
	if err != nil {
		return nil, err
	}

	resp, err := c.apiClient.DefaultApi.GetUser(username)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(resp.Values)
	if err != nil {
		return nil, err
	}

	var user BitbucketUser
	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
