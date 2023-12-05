package bitbucket_server

import (
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"golang.org/x/oauth2"
)

type Client interface {
	GetMe() (*BitbucketUser, error)
}

type BitbucketClient struct {
	apiClient *bitbucketv1.APIClient

	selfHostedURL    string
	selfHostedAPIURL string

	token oauth2.Token
}
