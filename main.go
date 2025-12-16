package main

import (
	"flag"
	"github.com/huangdijia/ccswitch/cmd"
)

// These variables are set by goreleaser during build
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	flag.StringVar(&version, "version", version, "version information")
}

func main() {
	// Set the version in the cmd package before executing
	cmd.SetVersion(version, commit, date)
	cmd.Execute()
}
