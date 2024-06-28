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

func NewGigaClient(ctx context.Context, clientId, clientSecret string) (*Client, error) {
	client, err := gigachat.NewInsecureClient(clientId, clientSecret)
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

	//models, err := c.client.ModelsWithContext(c.ctx)
	//_ = models

	req := &gigachat.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachat.Message{
			{
				Role:    "system",
				Content: c.prompt(locale, maxLength),
			},
			{
				Role:    "user",
				Content: diff,
			},
		},
		Temperature: ptr(0.7),
		TopP:        ptr(.1),
		MaxTokens:   ptr[int64](200),
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

func (c *Client) prompt(locale string, maxLength int) string {
	return fmt.Sprintf("Ты программист. Сгенерируйте краткое описание изменений для git commit, "+
		"для следующего diff с учетом приведенных ниже спецификаций:\n"+
		"Язык сообщения: %s\n"+
		"В сообщении о фиксации должно быть не более символов %d\n"+
		"Исключите все ненужно. Весь ваш ответ будет передан непосредственно в git commit.", locale, maxLength)
}

func ptr[T any](v T) *T {
	return &v
}
