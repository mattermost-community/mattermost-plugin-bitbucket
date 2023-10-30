package main

import (
	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

type BitbucketServerClient struct {
	configuration *bitbucketv1.Configuration
}

func NewClientServer() *BitbucketServerClient {
	return &BitbucketServerClient{}
}

func (c *BitbucketServerClient) NewConfiguration(config Configuration) *bitbucketv1.Configuration {
	c.configuration = bitbucketv1.NewConfiguration(config.BitbucketAPISelfHostedURL)
	return c.configuration
}

func (c *BitbucketServerClient) GetCurrentUser() string {
	return ""
}

func (c *BitbucketServerClient) GetMe() {

}
