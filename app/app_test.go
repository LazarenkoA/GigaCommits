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

		authKey := checkConf("test.yaml")
		assert.Equal(t, "111", authKey)
	})
	t.Run("test2", func(t *testing.T) {
		defer os.RemoveAll("test.yaml")

		os.Setenv(clientAuthKey, "authKey")

		authKey := checkConf("test.yaml")
		assert.Equal(t, "authKey", authKey)
	})
}
