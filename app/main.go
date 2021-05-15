package main

import (
	"github.com/jakubd/ipd"
	"github.com/jakubd/ipd/app/cmd"
)

func main() {
	ipd.CheckMaxmindEnvironment()
	cmd.Execute()
}
