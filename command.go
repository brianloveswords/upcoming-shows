package main

import (
	"strings"
)

type Command struct {
	Name     string
	Help     string
	Alias    []string
	Commands []*Command
}

type Subcommands []*Command

func (c *Command) String() string {
	var results []string
	var visit func(prefix string, c *Command)
	visit = func(prefix string, c *Command) {
		results = append(results, prefix+c.Name)
		for _, sub := range c.Commands {
			visit(prefix+c.Name+":", sub)
		}
	}
	visit("", c)
	return strings.Join(results, "\n")
}
