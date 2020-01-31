package main

import (
	"flag"
	"log"

	"github.com/Oppodelldog/pulli/src/pulli"
	"github.com/Oppodelldog/pulli/src/version"

	"os"
)

func main() {
	log.Printf("pulli (%v)", version.Number)

	var searchRoot string
	var filters arrayFlags
	var filterMode string

	flag.StringVar(&searchRoot, "dir", ".", "defines the folder where to find git repos")
	flag.Var(&filters, "filter", "filters the given folder. (can be absolute path or regex)")
	flag.StringVar(&filterMode, "filtermode", "", "whitelist or blacklist")
	flag.Parse()

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
