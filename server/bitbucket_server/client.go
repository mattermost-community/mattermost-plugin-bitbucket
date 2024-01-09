package bitbucket_server

type Client interface {
	GetMe() (*BitbucketUser, error)
}

type BitbucketClient struct {
	ClientConfiguration
}
