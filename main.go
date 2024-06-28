package main

import (
	"context"
	"flag"
	"github.com/LazarenkoA/GigaCommits/app"
	"github.com/pterm/pterm"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "вывод детальной информации")
	flag.Parse()

	if a, err := app.NewApp(context.Background(), debug); err == nil {
		a.Run()
	} else {
		pterm.Error.Println(err)
	}
}
