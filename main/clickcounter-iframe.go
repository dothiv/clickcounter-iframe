package main

import (
	"os"

	clickcounteriframe "github.com/dothiv/clickcounter-iframe"
	"github.com/wsxiaoys/terminal/color"
)

func error(msg string) {
	color.Fprintln(os.Stderr, "@{!r}ERROR @{|}"+msg)
}

func main() {
	configFilename := "config.ini"
	if len(os.Args) > 2 {
		configFilename = os.Args[2]
	}
	c, err := clickcounteriframe.NewConfig(configFilename)
	if err != nil {
		error(err.Error())
		os.Exit(1)
	}
	err = clickcounteriframe.Serve(c)
	if err != nil {
		error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
