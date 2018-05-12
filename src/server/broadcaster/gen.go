// +build ignore

package main

import (
	"os"
	"regexp"
	"strings"
	"text/template"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func main() {
	name := os.Args[1]
	templateName := os.Args[2]

	tmpl, err := template.ParseFiles(templateName)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(ToSnakeCase(name) + "_broadcaster.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, struct {
		Name string
	}{
		Name: name,
	})
	if err != nil {
		panic(err)
	}
}
