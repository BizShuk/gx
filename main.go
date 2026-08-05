package main

import (
	"github.com/bizshuk/gx/cmd"
	"github.com/bizshuk/gx/config"
)

func main() {
	config.Default()
	cmd.Execute()
}
