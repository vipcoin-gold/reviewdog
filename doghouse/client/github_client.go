package client

import (
	"context"

	"github.com/google/go-github/v49/github"

	"github.com/vipcoin-gold/reviewdog/doghouse"
	"github.com/vipcoin-gold/reviewdog/doghouse/server"
)

// GitHubClient is client which talks to GitHub directly instead of talking to
// doghouse server.
type GitHubClient struct {
	Client *github.Client
}

func (c *GitHubClient) Check(ctx context.Context, req *doghouse.CheckRequest) (*doghouse.CheckResponse, error) {
	return server.NewChecker(req, c.Client).Check(ctx)
}
