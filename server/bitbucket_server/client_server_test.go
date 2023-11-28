package bitbucket_server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockAPIClient struct {
	mock.Mock
}

func (m *MockAPIClient) GetUser(username string) (*http.Response, error) {
	args := m.Called(username)
	mockResp := mockHTTPResponse(200, `{"id": 1, "name": "John Doe", "links": {"self": [{"href": "http://example.com"}]}}`)
	return mockResp, args.Error(1)
}

func mockHTTPResponse(statusCode int, responseBody string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
	}
}

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

func TestBitbucketServerClient_GetMe(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(*MockAPIClient)
		expectedUser *BitbucketUser
		expectedErr  error
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *MockAPIClient) {
				m.On("GetUser", mock.Anything).Return(mockHTTPResponse(200, `{"id": 1, "name": "John Doe", "links": {"self": [{"href": "http://example.com"}]}}`), nil)
			},
			expectedUser: &BitbucketUser{ /* expected data */ },
			expectedErr:  nil,
		},
		{
			name:         "getWhoAmI fails",
			setupMock:    func(m *MockAPIClient) {},
			expectedUser: nil,
			expectedErr:  errors.New("getWhoAmI error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPIClient := &MockAPIClient{}
			tt.setupMock(mockAPIClient)

			c := &BitbucketServerClient{
				apiClient: mockAPIClient,
			}

			user, err := c.GetMe()

			assert.Equal(t, tt.expectedUser, user)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockAPIClient.AssertExpectations(t)
		})
	}
}
