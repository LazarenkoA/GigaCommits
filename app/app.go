package app

import (
	"context"
	"github.com/LazarenkoA/GigaCommits/giga"
	"github.com/LazarenkoA/GigaCommits/git"
	"github.com/pkg/errors"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/signal"
	"syscall"
)

//go:generate mockgen -source=$GOFILE -destination=./mock/mock.go
type IGiga interface {
	GetCommitMsg(diff string, locale string, maxLength int, debug bool) (string, error)
}

type confMap map[string]string

type IGit interface {
	GitDiff(debug bool) (string, error)
}

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	git    IGit
	ai     IGiga
	debug  bool
}

const (
	clientAuthKey = "authKey"
	confPath      = "giga_conf.yaml"
)

func NewApp(parentCtx context.Context, debug bool) (*App, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	authKey := checkConf(confPath)

	giga, err := giga.NewGigaClient(ctx, authKey)
	if err != nil {
		return nil, errors.Wrap(err, "create giga client")
	}

	return &App{
		ctx:    ctx,
		cancel: cancel,
		git:    git.NewGitClient(ctx),
		ai:     giga,
		debug:  debug,
	}, nil
}

func (a *App) Run() {
	go a.shutdown(syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)
	a.sendRequest()
}

func checkConf(confPath string) string {
	conf := loadConf(confPath)
	defer func() { saveConf(conf, confPath) }()

	authKey := os.Getenv(clientAuthKey)
	if authKey == "" {
		authKey = conf.getVal(clientAuthKey)
	}

	if authKey == "" {
		authKey, _ = pterm.DefaultInteractiveTextInput.Show("введите authKey")
		os.Setenv(authKey, authKey)
		conf[clientAuthKey] = authKey
	}

	return authKey
}

func loadConf(confPath string) confMap {
	result := map[string]string{}
	if f, err := os.Open(confPath); os.IsNotExist(err) {
		return map[string]string{}
	} else {
		if b, err := io.ReadAll(f); err == nil {
			yaml.Unmarshal(b, &result)
		}
	}

	return result
}

func saveConf(c confMap, confPath string) {
	b, _ := yaml.Marshal(c)
	if f, err := os.Open(confPath); !os.IsNotExist(err) {
		f.Write(b)
		f.Close()
	} else {
		f, _ := os.Create(confPath)
		f.Write(b)
		f.Close()
	}
}

func (a *App) shutdown(signals ...os.Signal) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)

	for {
		select {
		case signal := <-sigChan:
			switch signal {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
				a.cancel()
			}
		}
	}
}

func (a *App) sendRequest() {
	multi := pterm.DefaultMultiPrinter
	spinner1, _ := pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start("обрабобтка запроса")

	multi.Start()
	defer multi.Stop()

	//if err := a.git.DisableAutoCRLF(); err != nil {
	//	spinner1.Warning(err.Error())
	//}

	diff, err := a.git.GitDiff(a.debug)
	if err != nil {
		spinner1.Fail(err.Error())
		return
	}

	msg, err := a.ai.GetCommitMsg(diff, "ru", 100, a.debug)
	if err != nil {
		spinner1.Fail(err.Error())
		return
	}

	spinner1.Success(msg)
}

func (c confMap) getVal(key string) string {
	if v, ok := c[key]; ok {
		return v
	}

	return ""
}
