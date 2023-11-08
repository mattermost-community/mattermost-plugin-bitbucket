package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGetWhoAmI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer valid-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("mocked-user"))
	}))
	defer server.Close()

	tests := []struct {
		name           string
		token          oauth2.Token
		expectedResult string
		expectedError  string
	}{
		{
			name:           "valid token",
			token:          oauth2.Token{AccessToken: "valid-token"},
			expectedResult: "mocked-user",
			expectedError:  "",
		},
		{
			name:           "invalid token",
			token:          oauth2.Token{AccessToken: "invalid-token"},
			expectedResult: "",
			expectedError:  "who am i returned non-200 status code: 401",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClientServer()
			client.selfHostedURL = server.URL
			client.token = tc.token

			result, err := client.getWhoAmI()

			if tc.expectedError == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}
