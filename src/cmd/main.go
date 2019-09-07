package main

import (
	"flag"
	"github.com/Oppodelldog/pulli/src/pulli"
	"github.com/Oppodelldog/pulli/src/version"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	pulli.InitLogging()
	logrus.Infof("pulli (%v)", version.Number)

	var searchRoot string
	var filters arrayFlags
	var filterMode string
	var logLevel int

	flag.StringVar(&searchRoot, "dir", ".", "defines the folder where to find git repos")
	flag.Var(&filters, "filter", "filters the given folder. (can be absolute path or regex)")
	flag.StringVar(&filterMode, "filtermode", "", "whitelist or blacklist")
	flag.IntVar(&logLevel, "loglevel", int(logrus.InfoLevel), "0=panic, 1=fatal, 2=error, 3=warn, 4=info, 5=debug")
	flag.Parse()

	logrus.SetLevel(logrus.Level(logLevel))

	if ok := pulli.ValidateFlags(searchRoot, filterMode, filters); !ok {
		os.Exit(1)
	}

	pulli.PullAll(searchRoot, filters, filterMode)
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
