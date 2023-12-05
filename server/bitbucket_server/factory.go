package bitbucket_server

import (
	"fmt"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
)

func GetBitbucketClient(clientType string, apiClient *bitbucketv1.APIClient) (Client, error) {
	if clientType == "server" {
		return newServerClient(apiClient), nil
	}
	return nil, fmt.Errorf("wrong client passed")
}
