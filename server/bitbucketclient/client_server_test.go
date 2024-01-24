package bitbucketclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	bitbucketv1 "github.com/gfleury/go-bitbucket-v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) GetUser(_ string) (*bitbucketv1.APIResponse, error) {
	args := m.Called()
	return args.Get(0).(*bitbucketv1.APIResponse), args.Error(1)
}

func TestGetWhoAmI(t *testing.T) {
	t.Run("successfully get who am i from oauth", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte("testuser"))
			if err != nil {
				t.Fatalf("Failed to write response: %v", err)
			}
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
