package main

import (
	"fmt"
	"strings"
)

type Command struct {
	Name     string
	Help     string
	Alias    []string
	Commands []*Command
	Fn       func()
}

type Subcommands []*Command

func (c *Command) Run(args []string) error {
	fmt.Println(args)
	var visit func(remaining []string, c *Command) error
	visit = func(remaining []string, c *Command) error {
		fmt.Printf("visiting %s, %s\n", c.Name, remaining)

		if len(remaining) == 0 {
			if c.Fn != nil {
				c.Fn()
				return nil
			}
			return fmt.Errorf("%s", c.Help)
		}
		for _, sub := range c.Commands {
			fmt.Printf("comparing %s to %s\n", sub.Name, args[0])

			if sub.Name == remaining[0] {
				return visit(remaining[1:], sub)
			}
		}
		return fmt.Errorf("command not found")
	}
	return visit(args, c)
}

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
