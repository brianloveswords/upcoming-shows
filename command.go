package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

type Command struct {
	Name     string
	Help     string
	Alias    []string
	Commands []*Command
	Params   []Param
	Examples []Example
	Func     func()
}

type Example struct {
	Args     []string
	Comment  string
	Desc     string
	Negative bool
}

type Param struct {
	Name     string
	Help     string
	Alias    []string
	Raw      string
	Implicit bool
	ParseFn  func(string) error
}

func (p *Param) Consume(params []string) ([]string, error) {
	// fmt.Printf("param %s: checking params %s\n", p.Name, params)

	var names []string
	var value string
	var found bool
	names = append(p.Alias, p.Name)
	remaining := make([]string, 0, len(params))
	for _, param := range params {
		// add param to the remaining list and skip checking if we've
		// already found a param match
		if found {
			remaining = append(remaining, param)
			continue
		}

		for _, name := range names {
			// strategy 1: direct match
			if name == param {
				// fmt.Printf("FOUND PARAM MATCH: %s {%s}\n", name, param)
				found = true
				p.Raw = param
				break
			}

			// strategy 2: with =
			if strings.HasPrefix(param, name+"=") {
				// fmt.Printf("FOUND PARAM MATCH: %s {%s}\n", name, param)
				found = true
				p.Raw = param
				value = strings.Replace(param, name+"=", "", 1)
				break
			}
		}
		// at this point if we haven't found a match, add the param back
		// to the remaining list and skip trying to parse
		if !found {
			remaining = append(remaining, param)
			continue
		}

		err := p.ParseFn(value)
		if err != nil {
			return nil, err
		}
	}

	// if the param is "implicit", we want to parse it with an empty
	// value even if it wasn't found so we can get the default value
	if !found && p.Implicit {
		p.ParseFn("")
	}

	return remaining, nil
}

type Subcommands []*Command

func (c *Command) UsageAndExit() {
	var visit func(c *Command, prefix string)
	visit = func(c *Command, prefix string) {
		fullname := prefix + c.Name
		glog.Log("%s -- %s", color.GreenString(fullname), c.Help)

		for _, param := range c.Params {
			glog.Log("    %s :: %s", color.BlueString(param.Name), param.Help)
		}

		for _, sub := range c.Commands {
			visit(sub, fullname+" ")
		}
	}
	visit(c, "")
	os.Exit(1)
}

func (c *Command) ExamplesAndExit() {
	var visit func(c *Command, prefix string)
	visit = func(c *Command, prefix string) {
		fullname := prefix + c.Name
		for _, example := range c.Examples {
			var crayon func(string, ...interface{}) string
			var comment string
			crayon = color.GreenString
			if example.Negative {
				crayon = color.RedString
			}
			if example.Comment != "" {
				comment = color.YellowString(" # " + example.Comment)
			}
			command := fullname + " " + strings.Join(example.Args, " ")
			glog.Log("%s%s\n  %s\n", crayon(command), comment, example.Desc)
		}

		for _, sub := range c.Commands {
			visit(sub, fullname+" ")
		}
	}
	visit(c, "")
	os.Exit(1)
}

func (c *Command) Run(args []string) (err error) {
	// var spent []string
	var visit func(remaining []string, c *Command) error
	visit = func(remaining []string, c *Command) error {
		// fmt.Printf("visiting %s, %s\n", c.Name, remaining)

		// endpath 1: out of arguments and the command has a
		// corresponding function
		if len(remaining) == 0 {
			if c.Func == nil {

				// if this subcommand is merely the parent of further
				// subcommands and doesn't have a function, show the
				// usage and bounnnnnnce
				c.UsageAndExit()
			}
			c.Func()
			return nil
		}

		// endpath 2: help is in the way
		if remaining[0] == "--help" {
			c.UsageAndExit()
		}

		// endpath 3: examples are in the way
		if remaining[0] == "--examples" {
			c.ExamplesAndExit()
		}

		// endpath 4: out of children, has corresponding argument, and
		// param parser says everything's golden
		if len(c.Commands) == 0 {
			if c.Func == nil {
				return fmt.Errorf("%s", c.Help)
			}
			for _, param := range c.Params {
				remaining, err = param.Consume(remaining)
				if err != nil {
					return fmt.Errorf("error parsing param: %s", err)
				}
			}

			if len(remaining) > 0 {
				return fmt.Errorf("unknown params remaining: %s", remaining)
			}

			c.Func()
			return nil
		}

		for _, sub := range c.Commands {
			// fmt.Printf("comparing %s to %s\n", sub.Name, args[0])
			names := sub.Alias
			names = append(names, sub.Name)
			for _, name := range names {
				if name == remaining[0] {
					return visit(remaining[1:], sub)
				}
			}
		}
		return ErrCommandNotFound
	}
	return visit(args, c)
}

var ErrCommandNotFound = errors.New("command not found")

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
