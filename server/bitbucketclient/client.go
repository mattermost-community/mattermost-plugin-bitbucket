package bitbucketclient

type Client interface {
	GetMe() (*BitbucketUser, error)
}

type BitbucketClient struct {
	ClientConfiguration
}
