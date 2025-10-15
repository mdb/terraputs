[![CI/CD](https://github.com/mdb/terraputs/actions/workflows/main.yml/badge.svg)](https://github.com/mdb/terraputs/actions/workflows/main.yml)

# terraputs

A CLI to generate Terraform outputs documentation.

## Why?

`terraputs` analyzes the contents of a [terraform state](https://www.terraform.io/docs/language/state/index.html)
file and generates Markdown or HTML documentation from its [outputs](https://www.terraform.io/docs/language/values/outputs.html).

A common workflow executes `terraputs -state $(terraform show -json) > outputs.md`
after each invocation of `terraform apply`, then commits `outputs.md` to source
control or publishes its contents to GitHub release notes, offering up-to-date
documentation about resources managed by a Terraform project.

<a style="display: block;" href="https://asciinema.org/a/lFUVfdhes0i1cVbtFUvzwLMKd"><img style="width: 500px;" src="demo.svg"></a>

## Usage

```
terraputs -h
Usage of terraputs:
  -description string
        Optional; a contextual description preceding the outputs. (default "Terraform state outputs.")
  -heading string
        Optional; the heading text for use in the printed output. (default "Outputs")
  -output string
        Optional; the output format. Supported values: md, html. (default "md")
  -state string
        Optional; the state JSON output by 'terraform show -json'. Read from stdin if omitted
  -state-file string
        Optional; the path to a local file containing 'terraform show -json' output
```

Examples:

```
# terraputs can accept arbitrary state as a string argument on the command line:
terraputs -state "$(terraform show -json)" -heading "Terraform Outputs"

# ...Or directly from standard input, if the -state option is omitted. To read
# directly from terraform, consider:
terraform show -json | terraputs -heading "Terraform Outputs"

# ...Or read a tfstate JSON file from the filesystem:
terraform show -json > show.json
terraputs < show.json
```

### Results examples

<details>

<summary>Example markdown output</summary>

```
# Terraform Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string
```

</details>

<details>

<summary>Markdown output rendered via GitHub-flavored markdown</summary>

# Terraform Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

</details>

<details>

<summary>Example HTML output</summary>

```html
<h2>Outputs</h2>
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
```

</details>

<details>
<summary>HTML output rendered via GitHub-flavored markdown</summary>

<h2>Outputs</h2>
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

</details>
