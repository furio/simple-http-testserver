package main

import (
	"github.com/furio/simple-http-testserver/cmd"
)

var (
	Version string
	Build   string
)

func main() {
	cmd.Execute()
}
