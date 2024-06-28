package app

import (
	"github.com/agiledragon/gomonkey/v2"
	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_checkConf(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		p := gomonkey.ApplyMethod(pterm.DefaultInteractiveTextInput, "Show", func(_ pterm.InteractiveTextInputPrinter, text ...string) (string, error) {
			return "111", nil
		})

		defer p.Reset()
		defer os.RemoveAll("test.yaml")

		clientID, clientSecret := checkConf("test.yaml")
		assert.Equal(t, "111", clientID)
		assert.Equal(t, "111", clientSecret)
	})
	t.Run("test2", func(t *testing.T) {
		defer os.RemoveAll("test.yaml")

		os.Setenv(clientIDKey, "clientID")
		os.Setenv(clientSecretKey, "clientSecret")

		clientID, clientSecret := checkConf("test.yaml")
		assert.Equal(t, "clientID", clientID)
		assert.Equal(t, "clientSecret", clientSecret)
	})
}
