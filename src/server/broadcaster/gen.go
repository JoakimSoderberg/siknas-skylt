// +build ignore

package main

import (
	"strings"
	"os"
	"text/template"
)

func main() {
	name := os.Args[1]
	templateName := os.Args[2]

	tmpl, err := template.ParseFiles(templateName)
	if err != nil { panic(err) }

	f, err := os.Create(strings.ToLower(name) + "_broadcaster.go")
	if err != nil { panic(err) }
	defer f.Close()

	err = tmpl.Execute(f, struct {
		Name string
	}{
		Name: name,
	})
	if err != nil { panic(err) }
} 
