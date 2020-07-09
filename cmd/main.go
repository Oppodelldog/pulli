package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Oppodelldog/pulli/internal/pulli"
	"github.com/Oppodelldog/pulli/internal/version"

	"os"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var osArgsProvider = func() []string {
	return os.Args
}

func main() {
	log.Printf("pulli (%v)", version.Number)

	args := osArgsProvider()
	subCommandIndicator := 2

	if len(args) >= subCommandIndicator {
		buildCommand(args)
	} else {
		pullAll(args)
	}
}

var exitProgram = os.Exit

var pulliValidateFlags = pulli.ValidateFlags
var pulliPullAllFunc = pulli.PullAll
var pullAll = func(args []string) {
	var searchRoot string
	var filters arrayFlags
	var filterMode string

	fs := flag.NewFlagSet("pulli", flag.ContinueOnError)
	fs.StringVar(&searchRoot, pulli.ArgNameDir, ".", "defines the folder where to find git repos")
	fs.Var(&filters, pulli.ArgNameFilter, "filters the given folder. (can be absolute path or regex)")
	fs.StringVar(&filterMode, pulli.ArgNameFilterMode, "", "whitelist or blacklist")
	err := fs.Parse(args[1:])
	if err != nil {
		fmt.Println(err)
		exitProgram(1)
		return
	}

	if ok := pulliValidateFlags(searchRoot, filterMode, filters); !ok {
		flag.PrintDefaults()
		exitProgram(1)
		return
	}

	pulliPullAllFunc(searchRoot, filters, filterMode)
}

var pulliBuildCommandFunc = pulli.BuildCommand
var buildCommand = func(args []string) {
	var searchRoot string
	subCommand := args[1]
	if subCommand == pulli.SubCommandNameBuildCommand {
		fs := flag.NewFlagSet(pulli.SubCommandNameBuildCommand, flag.ContinueOnError)
		fs.StringVar(&searchRoot, pulli.ArgNameDir, ".", "defines the folder where to find git repos")
		err := fs.Parse(args[2:])
		if err != nil {
			fmt.Println(err)
			exitProgram(1)
			return
		}

		pulliBuildCommandFunc(searchRoot)
	}
}
