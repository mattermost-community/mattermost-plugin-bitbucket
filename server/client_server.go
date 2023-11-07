package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type BitbucketServerClient struct {
	apiClient *bitbucketv1.APIClient

	selfHostedURL    string
	selfHostedAPIURL string

	token oauth2.Token
}

func NewClientServer() *BitbucketServerClient {
	return &BitbucketServerClient{}
}

// TODO: This method needs to be changed when Modularization is built
func (c *BitbucketServerClient) Connect(config Configuration, token oauth2.Token, ts oauth2.TokenSource) *BitbucketServerClient {
	// setup Oauth context
	auth := context.WithValue(context.Background(), bitbucketv1.ContextOAuth2, ts)

	tc := oauth2.NewClient(auth, ts)

	// create config for bitbucket API
	configBb := bitbucketv1.NewConfiguration(config.BitbucketAPISelfHostedURL)
	configBb.HTTPClient = tc

	c.apiClient = bitbucketv1.NewAPIClient(context.Background(), configBb)
	c.selfHostedURL = config.BitbucketSelfHostedURL
	c.selfHostedAPIURL = config.BitbucketAPISelfHostedURL
	c.token = token

	return c
}

func (c *BitbucketServerClient) getWhoAmI() (string, error) {
	baseUrl := fmt.Sprintf("%s/plugins/servlet/applinks/whoami", c.selfHostedURL)

	req, err := http.NewRequest("GET", baseUrl, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

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

func (c *BitbucketServerClient) GetMe() (string, error) {
	return c.getWhoAmI()
}
