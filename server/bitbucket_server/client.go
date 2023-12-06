package bitbucket_server

type Client interface {
	GetMe(accessToken string) (*BitbucketUser, error)
}

type BitbucketClient struct {
	ClientConfiguration
}
