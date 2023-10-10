package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestPlugin() *Plugin {
	return &Plugin{}
}

func TestGetBitbucketBaseURL(t *testing.T) {
	p := setupTestPlugin()

	tests := []struct {
		name            string
		selfHostedURL   string
		expectedBaseURL string
	}{
		{
			name:            "With self hosted URL",
			selfHostedURL:   "https://selfhosted.bitbucket.com",
			expectedBaseURL: "https://selfhosted.bitbucket.com",
		},
		{
			name:            "Without self hosted URL",
			selfHostedURL:   "",
			expectedBaseURL: BitbucketBaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.configuration = &Configuration{
				BitbucketSelfHostedURL: tt.selfHostedURL,
			}

			assert.Equal(t, p.getBitbucketBaseURL(), tt.expectedBaseURL)
		})
	}
}

func TestGetBitbucketAPIBaseURL(t *testing.T) {
	p := setupTestPlugin()

	tests := []struct {
		name             string
		apiSelfHostedURL string
		expectedAPIURL   string
	}{
		{
			name:             "With self hosted API URL",
			apiSelfHostedURL: "https://api.selfhosted.example.com",
			expectedAPIURL:   "https://api.selfhosted.example.com",
		},
		{
			name:             "Without self hosted API URL",
			apiSelfHostedURL: "",
			expectedAPIURL:   BitbucketAPIBaseURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.configuration = &Configuration{
				BitbucketAPISelfHostedURL: tt.apiSelfHostedURL,
			}

			assert.Equal(t, p.getBitbucketAPIBaseURL(), tt.expectedAPIURL)
		})
	}
}
