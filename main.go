package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"errors"
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
	stateDesc      string = "Optional; the state JSON output by 'terraform show -json'. Read from stdin if omitted"
	stateFileDesc  string = "Optional; the path to a local file containing 'terraform show -json' output"
	headingDesc    string = "Optional; the heading text for use in the printed output. Default: Outputs"
	outputDesc     string = "Optional; the output format. Supported values: md, html. Default: md"
	versionDesc    string = "Print the current version and exit"
	defaultHeading string = "Outputs"
	defaultOutput  string = "md"
	sensitive      string = "sensitive; redacted"
)

type data struct {
	Heading string
	Outputs map[string]*tfjson.StateOutput
}

func value(output tfjson.StateOutput) string {
	if output.Sensitive {
		return sensitive
	}

	return fmt.Sprintf("%v", output.Value)
}

func prettyPrintValue(output tfjson.StateOutput) template.HTML {
	if output.Sensitive {
		return template.HTML(sensitive)
	}

	pretty, _ := json.MarshalIndent(output.Value, "", "  ")

	return template.HTML(string(pretty))
}

func dataType(output tfjson.StateOutput) string {
	return fmt.Sprintf("%T", output.Value)
}

func main() {
	stateJSON := flag.String("state", "", stateDesc)
	stateFile := flag.String("state-file", "", stateFileDesc)
	heading := flag.String("heading", defaultHeading, headingDesc)
	output := flag.String("output", defaultOutput, outputDesc)
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

	if *stateJSON != "" && *stateFile != "" {
		exit(errors.New("'-state' and '-state-file' are mutually exclusive; specify just one"))
	}

	if *stateFile != "" {
		b, err := os.ReadFile(*stateFile)
		if err != nil {
			exit(err)
		}

		stateReader = strings.NewReader(string(b))
	}

	var state *tfjson.State
	if err := json.NewDecoder(stateReader).Decode(&state); err != nil {
		exit(err)
	}

	tmpl, err := getTemplatePath(*output)
	if err != nil {
		exit(err)
	}

	t, err := template.New(strings.Split(tmpl, "/")[1]).Funcs(template.FuncMap{
		"value":       value,
		"dataType":    dataType,
		"prettyPrint": prettyPrintValue,
	}).ParseFS(templates, tmpl)
	if err != nil {
		exit(err)
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
		exit(err)
	}
}

func getTemplatePath(output string) (string, error) {
	switch output {
	case "html":
		return "templates/html.tmpl", nil
	case "md":
		return "templates/markdown.tmpl", nil
	default:
		return "", fmt.Errorf("%s is not a supported output format. Supported formats: md (default), html", output)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s", err.Error())
	os.Exit(1)
}
