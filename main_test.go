package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const testVersion string = "test"

func TestMain(m *testing.M) {
	// compile an 'terraputs' for for use in running tests
	exe := exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.version=%s", testVersion), "-o", "terraputs")
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	m.Run()

	// delete the compiled terraputs
	err = os.Remove("terraputs")
	if err != nil {
		log.Fatal(err)
	}
}

func TestHelpFlag(t *testing.T) {
	help := []string{
		"-state string",
		stateDesc,
		"-heading string",
		headingDesc,
	}

	tests := []struct {
		arg     string
		outputs []string
	}{{
		arg:     "-h",
		outputs: help,
	}, {
		arg:     "-help",
		outputs: help,
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when terraputs is passed '%s'", test.arg), func(t *testing.T) {
			output, err := exec.Command("./terraputs", test.arg).CombinedOutput()

			if err != nil {
				t.Errorf("expected '%s' not to error; got '%v'", test.arg, err)
			}

			for _, o := range test.outputs {
				if !strings.Contains(string(output), o) {
					t.Errorf("expected '%s' to include output '%s'; got '%s'", test.arg, o, output)
				}
			}
		})
	}
}

func TestVersionArg(t *testing.T) {
	args := []string{
		"version",
		" version ",
		"VERSION",
		"vERsION",
	}

	for _, arg := range args {
		t.Run(fmt.Sprintf("when terraputs is passed '%s'", arg), func(t *testing.T) {
			output, err := exec.Command("./terraputs", arg).CombinedOutput()

			if err != nil {
				t.Errorf("expected '%s' not to cause error; got '%v'", arg, err)
			}

			if !strings.Contains(string(output), testVersion) {
				t.Errorf("expected '%s' to output version '%s'; got '%s'", arg, testVersion, output)
			}
		})
	}
}

// These should all be ok and functionally equivalent:
//   terraputs -state "$(cat stateFile)"
//   terraputs < stateFile
//   cat stateFile | terraputs
func TestTerraputs(t *testing.T) {
	tests := []struct {
		command        string
		shouldError    bool
		expectedOutput string
	}{{
		command: `./terraputs -state $(cat testdata/basic/show.json)`,
		expectedOutput: `# Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

`,
	}, {
		command: `./terraputs < testdata/basic/show.json`,
		expectedOutput: `# Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

`,
	}, {
		command: `cat testdata/basic/show.json | ./terraputs`,
		expectedOutput: `# Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

`,
	}, {
		command: `./terraputs -state $(cat testdata/basic/show.json) -heading foo`,
		expectedOutput: `# foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

`,
	}, {
		command: `./terraputs -state $(cat testdata/nooutputs/show.json) -heading foo`,
		expectedOutput: `# foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}, {
		command: `./terraputs -state $(cat testdata/emptyconfig/show.json) -heading foo`,
		expectedOutput: `# foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when terraputs is passed '%q'", test.command), func(t *testing.T) {
			output, err := exec.Command("/bin/sh", "-c", test.command).CombinedOutput()

			if err != nil {
				t.Errorf("expected '%s' not to error; got '%v'", test.command, err)
			}

			if string(output) != test.expectedOutput {
				t.Logf("expected output: \n%s", test.expectedOutput)
				t.Logf("got output: \n%s", output)
				t.Errorf("received unexpected output from '%s'", test.command)
			}
		})
	}
}
