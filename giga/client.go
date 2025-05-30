package giga

import (
	"context"
	"fmt"
	"github.com/k0kubun/pp/v3"
	"github.com/paulrzcz/go-gigachat"
	"github.com/pkg/errors"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/mock.go
type IGigaClient interface {
	AuthWithContext(ctx context.Context) error
	ChatWithContext(ctx context.Context, in *gigachat.ChatRequest) (*gigachat.ChatResponse, error)
}

type Client struct {
	ctx    context.Context
	client IGigaClient
}

func NewGigaClient(ctx context.Context, authKey string) (*Client, error) {
	client, err := gigachat.NewInsecureClientWithAuthKey(authKey)
	if err != nil {
		return nil, errors.Wrap(err, "newGigaClient error")
	}

	return &Client{
		ctx:    ctx,
		client: client,
	}, nil
}

func (c *Client) GetCommitMsg(diff string, locale string, maxLength int, debug bool) (string, error) {
	err := c.client.AuthWithContext(c.ctx)
	if err != nil {
		return "", errors.Wrap(err, "auth error")
	}

	if diff == "" {
		return "", errors.New("diff is not defined")
	}

	req := &gigachat.ChatRequest{
		Model: "GigaChat-Pro",
		Messages: []gigachat.Message{
			{
				Role:    "system",
				Content: c.prompt(locale),
			},
			{
				Role:    "user",
				Content: "вот результат команды \"git diff\"\n" + diff,
			},
		},
		Temperature: ptr(0.7),
		MaxTokens:   ptr[int64](1000),
	}

	if debug {
		pp.Println("REQUEST:")
		pp.Println(req)
	}

	resp, err := c.client.ChatWithContext(c.ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "request error")
	}

	if debug {
		pp.Println("RESPONSE:")
		pp.Println(resp)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("response does not contain data")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) prompt(locale string) string {
	return fmt.Sprintf("Ты программист. Проанализируй отправленый тебе результат команды \"git diff\" и сгенерируйте КРАТКОЕ описание ИЗ ОДНОГО ПРЕДЛОЖЕНИЯ для команды git commit. "+
		"Необходимо учесть приведенные ниже спецификации:\n"+
		"Язык сообщения: %s\n"+
		"Исключите все ненужно. Весь ваш ответ будет передан непосредственно в git commit.", locale)
}

func ptr[T any](v T) *T {
	return &v
}
