package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zmb3/spotify"
)

var (
	ptrue  = true
	pfalse = false
)

type ParamList []Parameter

type Base struct {
	Name  string
	Help  string
	Short rune
}

func (b *Base) Extract(arg string) (string, bool) {
	parts := strings.SplitN(arg, "=", 2)
	re := regexp.MustCompile(fmt.Sprintf("^(--%s|-%c)$", b.Name, b.Short))
	if match := re.Match([]byte(parts[0])); match {
		if len(parts) == 2 {
			return parts[1], true
		}
		return "", true
	}
	return "", false
}

func (pl *ParamList) Parse(args []string) error {
	// fmt.Println(args)
	// for _, arg := range args {
	// 	for _, p := range pl {
	// 		base := p.GetBase()
	// 		base.Match(arg)
	// 	}
	// }
	return nil
}

type Parameter interface {
	Parse(string) error
	GetBase() Base
}

type BoolParam struct {
	Base
	Var *bool
}

func (p *BoolParam) Parse(arg string) error {
	if arg == "" {
		p.Var = &ptrue
		return nil
	}
	return fmt.Errorf("boolean arg `%s` does not take value", p.Name)
}

type IntParam struct {
	Base
	Var     *int
	Default int
	Min     *int
	Max     *int
}

func (p *IntParam) GetBase() Base {
	return p.Base
}
func (p *IntParam) Parse(arg string) error {
	val, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Errorf("`%s` must be a number", p.Name)
	}
	if p.Min != nil && val < *p.Min {
		return fmt.Errorf("`%s` must be greater than %d", p.Name, *p.Min)
	}
	if p.Max != nil && val > *p.Min {
		return fmt.Errorf("`%s` must be less than %d", p.Name, *p.Max)
	}
	p.Var = &val
	return err
}

type StringParam struct {
	Base
	Var     *string
	Default string
}

func (p *StringParam) GetBase() Base {
	return p.Base
}
func (p *StringParam) Parse(arg string) error {
	p.Var = &arg
	return nil
}

type IDParam struct {
	Base
	Var     *spotify.ID
	Default string
}

func (p *IDParam) GetBase() Base {
	return p.Base
}
func (p *IDParam) Parse(arg string) error {
	id := spotify.ID(arg)
	p.Var = &id
	return nil
}
