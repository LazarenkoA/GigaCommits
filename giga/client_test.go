package giga

import (
	"context"
	mock_giga "github.com/LazarenkoA/GigaCommits/giga/mock"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/golang/mock/gomock"
	"github.com/paulrzcz/go-gigachat"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_prompt(t *testing.T) {
	prompt := new(Client).prompt("ru", 100)
	assert.Equal(t, `Ты программист. Сгенерируйте краткое описание изменений для git commit, для следующего diff с учетом приведенных ниже спецификаций:
Язык сообщения: ru
В сообщении о фиксации должно быть не более символов 100
Исключите все ненужно. Весь ваш ответ будет передан непосредственно в git commit.`, prompt)
}

func Test_GetCommitMsg(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	t.Run("error create", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*Client, error) {
			return nil, errors.New("error")
		})
		defer p.Reset()

		cli, err := NewGigaClient(context.Background(), "111", "222")
		assert.Nil(t, cli)
		assert.EqualError(t, err, "newGigaClient error: error")
	})
	t.Run("auth error", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*gigachat.Client, error) {
			return new(gigachat.Client), nil
		})
		defer p.Reset()

		client := mock_giga.NewMockIGigaClient(c)
		client.EXPECT().AuthWithContext(gomock.Any()).Return(errors.New("error"))

		cli, _ := NewGigaClient(context.Background(), "111", "222")
		cli.client = client

		_, err := cli.GetCommitMsg("", "ru", 100, false)
		assert.EqualError(t, err, "auth error: error")
	})
	t.Run("req error", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*gigachat.Client, error) {
			return new(gigachat.Client), nil
		})
		defer p.Reset()

		client := mock_giga.NewMockIGigaClient(c)
		client.EXPECT().AuthWithContext(gomock.Any()).Return(nil)
		client.EXPECT().ChatWithContext(gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))

		cli, _ := NewGigaClient(context.Background(), "111", "222")
		cli.client = client

		_, err := cli.GetCommitMsg("tyuyu", "ru", 100, false)
		assert.EqualError(t, err, "request error: error")
	})
	t.Run("response does not contain data", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*gigachat.Client, error) {
			return new(gigachat.Client), nil
		})
		defer p.Reset()

		client := mock_giga.NewMockIGigaClient(c)
		client.EXPECT().AuthWithContext(gomock.Any()).Return(nil)
		client.EXPECT().ChatWithContext(gomock.Any(), gomock.Any()).Return(&gigachat.ChatResponse{}, nil)

		cli, _ := NewGigaClient(context.Background(), "111", "222")
		cli.client = client

		_, err := cli.GetCommitMsg("ghgh", "ru", 100, false)
		assert.EqualError(t, err, "response does not contain data")
	})
	t.Run("diff is not defined", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*gigachat.Client, error) {
			return new(gigachat.Client), nil
		})
		defer p.Reset()

		client := mock_giga.NewMockIGigaClient(c)
		client.EXPECT().AuthWithContext(gomock.Any()).Return(nil)

		cli, _ := NewGigaClient(context.Background(), "111", "222")
		cli.client = client

		_, err := cli.GetCommitMsg("", "ru", 100, false)
		assert.EqualError(t, err, "diff is not defined")
	})
	t.Run("pass", func(t *testing.T) {
		p := gomonkey.ApplyFunc(gigachat.NewInsecureClient, func(clientId string, clientSecret string) (*gigachat.Client, error) {
			return new(gigachat.Client), nil
		})
		defer p.Reset()

		client := mock_giga.NewMockIGigaClient(c)
		client.EXPECT().AuthWithContext(gomock.Any()).Return(nil)
		client.EXPECT().ChatWithContext(gomock.Any(), gomock.Any()).Return(&gigachat.ChatResponse{
			Choices: []gigachat.Choice{{Message: gigachat.Message{Content: "test"}}},
		}, nil)

		cli, _ := NewGigaClient(context.Background(), "111", "222")
		cli.client = client

		msg, err := cli.GetCommitMsg("hjhj", "ru", 100, false)
		assert.NoError(t, err)
		assert.Equal(t, "test", msg)
	})
}
