package git

import (
	"context"
	"github.com/k0kubun/pp/v3"
	"github.com/pkg/errors"
	"os/exec"
)

const (
	gitBin = "git"
)

type Client struct {
	ctx context.Context
}

func NewGitClient(ctx context.Context) *Client {
	return &Client{
		ctx: ctx,
	}
}

func (c *Client) GitDiff(debug bool) (string, error) {
	path, err := exec.LookPath(gitBin)
	if err != nil {
		return "", errors.Wrap(err, "lookPath error")
	}

	cmd := exec.CommandContext(c.ctx, path, "diff", "--diff-algorithm=minimal") // --cached
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, "execute git error")
	}

	if debug {
		pp.Println(cmd)
	}

	return string(output), nil
}
