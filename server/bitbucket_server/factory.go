package bitbucket_server

import (
	"fmt"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

type ClientConfiguration struct {
	SelfHostedURL    string
	SelfHostedAPIURL string
	APIClient        *bitbucketv1.APIClient

	LogError func(msg string, keyValuePairs ...interface{})
}

func GetBitbucketClient(clientType string, config ClientConfiguration) (Client, error) {
	if clientType == "server" {
		return newServerClient(config), nil
	}
	return nil, fmt.Errorf("wrong client passed")
}
