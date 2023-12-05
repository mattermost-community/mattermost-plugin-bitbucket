package bitbucket_server

import (
	"fmt"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

func GetBitbucketClient(clientType string, selfHostedURL string, selfHostedAPIURL string, apiClient *bitbucketv1.APIClient) (Client, error) {
	if clientType == "server" {
		return newServerClient(selfHostedURL, selfHostedAPIURL, apiClient), nil
	}
	return nil, fmt.Errorf("wrong client passed")
}
