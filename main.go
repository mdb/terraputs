package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

var (
	// Embed the entire templates directory in the compiled binary.
	//go:embed templates
	templates embed.FS

	// passed in at build time
	version string
)

const (
	stateDesc   string = "Required; the state JSON output by 'terraform show -json'"
	versionDesc string = "Print the current version and exit"
)

type data struct {
	Outputs map[string]*tfjson.StateOutput
}

func value(output tfjson.StateOutput) string {
	if output.Sensitive {
		return "sensitive; redacted"
	}

	return fmt.Sprintf("%v", output.Value)
}

func dataType(output tfjson.StateOutput) string {
	return fmt.Sprintf("%T", output.Value)
}

func main() {
	stateJSON := flag.String("state", "", stateDesc)
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && strings.ToLower(strings.TrimSpace(args[0])) == "version" {
		fmt.Println(version)
		return
	}

	if *stateJSON == "" {
		panic(errors.New("Mising required '-state' value"))
	}

	var state *tfjson.State
	if err := json.NewDecoder(strings.NewReader(*stateJSON)).Decode(&state); err != nil {
		panic(err)
	}

	t, err := template.New("markdown.tmpl").Funcs(template.FuncMap{
		"value":    value,
		"dataType": dataType,
	}).ParseFS(templates, "templates/markdown.tmpl")
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, data{
		Outputs: state.Values.Outputs,
	})
	if err != nil {
		panic(err)
	}
}
