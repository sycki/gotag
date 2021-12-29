package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

func main() {
	config := NewConfig()
	config.ParseFromCommand()

	for _, filePattern := range config.FilePatterns {
		fileNames, err := filepath.Glob(filePattern)
		if err != nil {
			println(errors.Wrap(err, "filepath.Glob").Error())
			return
		}

		for _, fileName := range fileNames {
			if err := FormatFile(config, fileName); err != nil {
				println(errors.Wrap(err, "FormatFile").Error())
				return
			}
		}
	}
}

func FormatFile(config *Config, fileName string) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parser.ParseFile")
	}
	ast.Inspect(astFile, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.StructType:
			FormatStruct(config, t)
			return false
		}
		return true
	})

	write, err := os.Create(fileName)
	if err != nil {
		return errors.Wrap(err, "os.Create")
	}
	defer write.Close()
	err = format.Node(write, fset, astFile)
	if err != nil {
		return errors.Wrap(err, "format.Node")
	}

	return write.Close()
}

func FormatStruct(config *Config, node *ast.StructType) {
	for _, field := range node.Fields.List {
		if len(field.Names) < 1 {
			continue
		}
		if !unicode.IsUpper(rune(field.Names[0].String()[0])) {
			continue
		}

		if field.Tag == nil {
			field.Tag = &ast.BasicLit{}
			field.Tag.ValuePos = field.Type.Pos() + 1
			field.Tag.Kind = token.STRING
		}
		FormatField(config, field)
	}

	return
}

func FormatField(config *Config, field *ast.Field) {
	var fieldName = field.Names[0].String()
	var oldTagsStr = field.Tag.Value
	var oldTags = strings.Fields(strings.Trim(oldTagsStr, "`"))
	var addTags []string
	var addedTags []string
	var removedTags []string

	index := config.Index
	if index < 0 {
		oldTagsStr = ""
		oldTags = nil
		index = 0
	}
	if index > len(oldTags) {
		index = len(oldTags)
	}

	// add tags by --add
	for _, addTag := range config.AddTags {
		tagReg := regexp.MustCompile(fmt.Sprintf(`%s:"[^"]+"`, addTag.TagKey))
		matchedTag := tagReg.FindString(oldTagsStr)
		if matchedTag == "" {
			var tagValue string
			switch addTag.TagValueStyle {
			case TagValueStyleCamel:
				tagValue = ToCamel(fieldName)
			case TagValueStyleSnake:
				tagValue = ToSnake(fieldName)
			case TagValueStyleGo:
				tagValue = fieldName
			case TagValueStyleUpper:
				tagValue = strings.ToUpper(fieldName)
			case TagValueStyleLower:
				tagValue = strings.ToLower(fieldName)
			default:
				fmt.Fprintf(os.Stderr, "unspported tag style %s", addTag.TagValueStyle)
			}
			newTag := fmt.Sprintf(`%s:"%s"`, addTag.TagKey, tagValue)
			addTags = append(addTags, newTag)
		}
	}
	addedTags = append(addedTags, oldTags[:index]...)
	addedTags = append(addedTags, addTags...)
	addedTags = append(addedTags, oldTags[index:]...)

	// remove tags by --remove
	for _, addedTag := range addedTags {
		words := strings.Split(addedTag, ":")
		if len(words) < 1 {
			continue
		}
		var isContain bool
		for _, removeTag := range config.RemoveTags {
			if words[0] == removeTag.TagKey {
				isContain = true
				break
			}
		}
		if !isContain {
			removedTags = append(removedTags, addedTag)
		}
	}

	if len(removedTags) > 0 {
		field.Tag.Value = fmt.Sprintf("`%s`", strings.Join(removedTags, " "))
	} else {
		field.Tag.Value = ""
	}
	return
}

var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnake(str string) string {
	if len(str) < 1 {
		return str
	}
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToCamel(str string) string {
	runes := []rune(str)
	var i int
	for i = 0; i < len(runes); i++ {
		if unicode.IsLower(runes[i]) {
			break
		}
		runes[i] = unicode.ToLower(runes[i])
	}
	if i != 1 && i != len(runes) {
		i--
		runes[i] = unicode.ToUpper(runes[i])
	}
	return string(runes)
}
