package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	styleList      = []string{"camel", "snake", "go", "upper", "lower"}
	styleAliasList = map[string][]string{
		"camel": {"aA"},
		"snake": {"a_a"},
		"go":    {"Aa"},
		"upper": {"AA"},
		"lower": {"aa"},
	}
)

type Config struct {
	FilePatterns   []string
	Index          int
	AddTagFlags    []string
	RemoveTagFlags []string
	AddTags        []TagDescribe
	RemoveTags     []TagDescribe
}

func NewConfig() *Config {
	return &Config{}
}

func (p *Config) ParseFromCommand() {
	cmd := pflag.NewFlagSet("gotag", pflag.ExitOnError)
	cmd.Usage = func() {
		templ := `
gotag is a command tool that can automatically generate tags for golang struct.

Usage: gotag [OPTIONS] filename ...

Options:
    -a, --add <TagName:Style>   Add tags to struct
                                TagName:     can be any string
                                Style:       "camel", "snake", "go", "upper", "lower"
                                Style alias: "aA", "a_a", "Aa", "AA", "aa"
    -i, --index <Number>        Specify the position for the --add option
                                Negative number: means cover all
                                0 ~ max integer: insert to the specified position
    -r, --remove <TagName>      Remove tags from struct

Examples:
    gotag -a json:camel -a gorm:snake model/*.go

    gotag -a json:aA,gorm:a_a -i -1 model/*.go

    gotag -r json model/*.go

`

		fmt.Fprintf(os.Stderr, templ)
		os.Exit(0)
	}
	cmd.StringSliceVarP(&p.AddTagFlags, "add", "a", nil, `For example: -a "json:snake myTag:camel"`)
	cmd.StringSliceVarP(&p.RemoveTagFlags, "remove", "r", nil, `For example: -r "json myTag"`)
	cmd.IntVarP(&p.Index, "index", "i", 0, `Effective value: -1, 0, 1, 2 ...`)
	cmd.Parse(os.Args)
	p.FilePatterns = cmd.Args()[1:]

	if len(p.FilePatterns) < 1 ||
		(len(p.AddTagFlags) < 1 && len(p.RemoveTagFlags) < 1) {
		cmd.Usage()
	}

	for _, flg := range p.AddTagFlags {
		words := strings.Split(flg, ":")
		if len(words) != 2 {
			println("invalid tag:", flg)
			os.Exit(3)
		}
		var isSupported bool
		for _, v := range styleList {
			if words[1] == v {
				isSupported = true
				break
			}
		}
		if !isSupported {
			var isStyleAlias bool
			for tagKey, styleAlias := range styleAliasList {
				for _, styleAlia := range styleAlias {
					if words[1] == styleAlia {
						isStyleAlias = true
						words[1] = tagKey
						break
					}
				}
				if isStyleAlias {
					break
				}
			}
			if !isStyleAlias {
				println("unsupported style:", words[1])
				os.Exit(5)
			}
		}
		p.AddTags = append(p.AddTags, TagDescribe{TagKey: words[0], TagValueStyle: words[1]})
	}
	for _, flg := range p.RemoveTagFlags {
		p.RemoveTags = append(p.RemoveTags, TagDescribe{TagKey: flg})
	}

	return
}
