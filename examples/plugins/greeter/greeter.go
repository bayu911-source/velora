
package main

import (
	"fmt"
)

// Agent is the exported symbol that the plugin loader will look for.
// It must be a variable of a type that implements the plugin.Agent interface.
var Agent greeter

type greeter struct{}

func (g greeter) Name() string {
	return "greeter"
}

func (g greeter) Run(input string) (string, error) {
	return fmt.Sprintf("Hello, %s!", input), nil
}
