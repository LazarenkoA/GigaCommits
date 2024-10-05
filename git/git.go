package git

import (
	"bytes"
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

func (c *Client) DisableAutoCRLF() error {
	path, err := exec.LookPath(gitBin)
	if err != nil {
		return errors.Wrap(err, "lookPath error")
	}

	cmd := exec.CommandContext(c.ctx, path, "config", "core.autocrlf", "false") // "--global"
	err = cmd.Run()

	return errors.Wrap(err, "git exec error")
}

func (c *Client) GitDiff(debug bool) (string, error) {
	path, err := exec.LookPath(gitBin)
	if err != nil {
		return "", errors.Wrap(err, "lookPath error")
	}

	cmd := exec.CommandContext(c.ctx, path, "diff", "--diff-algorithm=minimal") // --cached

	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err = cmd.Run()
	if err != nil {
		pp.Println("STDERR:", stderr.String())
		return "", errors.Wrap(err, "execute git error")
	}

	if debug {
		pp.Println(cmd)
	}

	return stdout.String(), nil
}
