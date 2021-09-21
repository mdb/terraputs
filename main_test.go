package main

import (
	"errors"
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
		"-state-file string",
		stateFileDesc,
		"-heading string",
		headingDesc,
		"-output string",
		outputDesc,
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
		expectedError  error
		expectedOutput string
	}{{
		command: `./terraputs -state $(cat testdata/basic/show.json)`,
		expectedOutput: `## Outputs

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
		expectedOutput: `## Outputs

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
		expectedOutput: `## Outputs

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
		expectedOutput: `## foo

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
		command: `./terraputs -state $(cat testdata/basic/show.json) -description "A custom description."`,
		expectedOutput: `## Outputs

A custom description.

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
		expectedOutput: `## foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}, {
		command: `./terraputs -state $(cat testdata/emptyconfig/show.json) -heading foo`,
		expectedOutput: `## foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}, {
		command: `./terraputs -state-file testdata/basic/show.json -heading foo`,
		expectedOutput: `## foo

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
		command: `./terraputs -state-file testdata/basic/show.json -heading foo -output html`,
		expectedOutput: `<h2>foo</h2>
<p>Terraform state outputs.</p>
<table>
  <tr>
    <th>Output</th>
    <th>Value</th>
    <th>Type</th>
  </tr>

  <tr>
    <td>a_basic_map</td>
    <td><pre>{
  "foo": "bar",
  "number": 42
}</pre></td>
    <td>map[string]interface {}</td>
  </tr>

  <tr>
    <td>a_list</td>
    <td><pre>[
  "foo",
  "bar"
]</pre></td>
    <td>[]interface {}</td>
  </tr>

  <tr>
    <td>a_nested_map</td>
    <td><pre>{
  "baz": {
    "bar": "baz",
    "id": "123"
  },
  "foo": "bar",
  "number": 42
}</pre></td>
    <td>map[string]interface {}</td>
  </tr>

  <tr>
    <td>a_sensitive_value</td>
    <td><pre>sensitive; redacted</pre></td>
    <td>string</td>
  </tr>

  <tr>
    <td>a_string</td>
    <td><pre>"foo"</pre></td>
    <td>string</td>
  </tr>

</table>
`,
	}, {
		command: `./terraputs -state-file testdata/basic/show.json -description "A custom description." -output html`,
		expectedOutput: `<h2>Outputs</h2>
<p>A custom description.</p>
<table>
  <tr>
    <th>Output</th>
    <th>Value</th>
    <th>Type</th>
  </tr>

  <tr>
    <td>a_basic_map</td>
    <td><pre>{
  "foo": "bar",
  "number": 42
}</pre></td>
    <td>map[string]interface {}</td>
  </tr>

  <tr>
    <td>a_list</td>
    <td><pre>[
  "foo",
  "bar"
]</pre></td>
    <td>[]interface {}</td>
  </tr>

  <tr>
    <td>a_nested_map</td>
    <td><pre>{
  "baz": {
    "bar": "baz",
    "id": "123"
  },
  "foo": "bar",
  "number": 42
}</pre></td>
    <td>map[string]interface {}</td>
  </tr>

  <tr>
    <td>a_sensitive_value</td>
    <td><pre>sensitive; redacted</pre></td>
    <td>string</td>
  </tr>

  <tr>
    <td>a_string</td>
    <td><pre>"foo"</pre></td>
    <td>string</td>
  </tr>

</table>
`,
	}, {
		command: `./terraputs -state-file testdata/nooutputs/show.json -heading foo`,
		expectedOutput: `## foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}, {
		command: `./terraputs -state-file testdata/emptyconfig/show.json -heading foo`,
		expectedOutput: `## foo

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |

`,
	}, {
		command:        `./terraputs -state-file testdata/basic/i-do-not-exist.json -heading foo`,
		expectedError:  errors.New("exit status 1"),
		expectedOutput: "open testdata/basic/i-do-not-exist.json: no such file or directory",
	}, {
		command:        `./terraputs -state $(cat testdata/basic/show.json) -state-file testdata/basic/show.json -heading foo`,
		expectedError:  errors.New("exit status 1"),
		expectedOutput: "'-state' and '-state-file' are mutually exclusive; specify just one",
	}, {
		command:        `./terraputs -state $(cat testdata/basic/show.json) -output foo`,
		expectedError:  errors.New("exit status 1"),
		expectedOutput: "'foo' is not a supported output format. Supported formats: 'md' (default), 'html'",
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when terraputs is passed '%q'", test.command), func(t *testing.T) {
			output, err := exec.Command("/bin/sh", "-c", test.command).CombinedOutput()
			if err != nil && test.expectedError == nil {
				t.Errorf("expected '%s' not to error; got '%v'", test.command, err)
			}

			if test.expectedError != nil && err == nil {
				t.Errorf("expected error '%s'; got '%v'", test.expectedError.Error(), err)
			}

			if test.expectedError != nil && err != nil && test.expectedError.Error() != err.Error() {
				t.Errorf("expected error '%s'; got '%v'", test.expectedError.Error(), err.Error())
			}

			if string(output) != test.expectedOutput {
				t.Logf("expected output: \n%s", test.expectedOutput)
				t.Logf("got output: \n%s", output)
				t.Errorf("received unexpected output from '%s'", test.command)
			}
		})
	}
}
