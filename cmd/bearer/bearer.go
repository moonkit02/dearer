package main

import (
	"github.com/moonkit02/dearer/cmd/bearer/build"
	"github.com/moonkit02/dearer/external/run"
)

func main() {
	run.Run(build.Version, build.CommitSHA, run.NewEngine(run.DefaultLanguages()))
}
