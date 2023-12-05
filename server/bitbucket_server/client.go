package bitbucket_server

import bitbucketv1 "github.com/gfleury/go-bitbucket-v1"

type Client interface {
	GetMe(accessToken string) (*BitbucketUser, error)
}

type BitbucketClient struct {
	apiClient *bitbucketv1.APIClient

	selfHostedURL    string
	selfHostedAPIURL string
}
