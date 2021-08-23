package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
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
	stateDesc      string = "Optional; the state JSON output by 'terraform show -json', read from stdin if omitted"
	headingDesc    string = "Optional; the heading text for use in the printed markdown"
	versionDesc    string = "Print the current version and exit"
	defaultHeading string = "Outputs"
)

type data struct {
	Heading string
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
	heading := flag.String("heading", defaultHeading, headingDesc)
	flag.Parse()

	args := flag.Args()
	if len(args) > 0 && strings.ToLower(strings.TrimSpace(args[0])) == "version" {
		fmt.Println(version)
		return
	}

	var stateReader io.Reader = bufio.NewReader(os.Stdin)
	if *stateJSON != "" {
		stateReader = strings.NewReader(*stateJSON)
	}

	var state *tfjson.State
	if err := json.NewDecoder(stateReader).Decode(&state); err != nil {
		panic(err)
	}

	t, err := template.New("markdown.tmpl").Funcs(template.FuncMap{
		"value":    value,
		"dataType": dataType,
	}).ParseFS(templates, "templates/markdown.tmpl")
	if err != nil {
		panic(err)
	}

	outputs := map[string]*tfjson.StateOutput{}
	if state.Values != nil {
		outputs = state.Values.Outputs
	}

	err = t.Execute(os.Stdout, data{
		Outputs: outputs,
		Heading: *heading,
	})
	if err != nil {
		panic(err)
	}
}
