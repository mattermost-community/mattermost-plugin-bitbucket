package bitbucket_server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetWhoAmI(t *testing.T) {
	t.Run("successfuly get who am i from oauth", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("testuser"))
		}))
		defer ts.Close()

		client := newServerClient(ClientConfiguration{SelfHostedURL: ts.URL, OAuthClient: ts.Client()}).(*BitbucketServerClient)

		username, err := client.getWhoAmI()

		assert.Nil(t, err)
		assert.Equal(t, "testuser", username)
	})

	t.Run("return error when the status code is different than 200", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer ts.Close()

		client := newServerClient(ClientConfiguration{SelfHostedURL: ts.URL, OAuthClient: ts.Client()}).(*BitbucketServerClient)

		_, err := client.getWhoAmI()
		assert.Error(t, err)
	})
}
