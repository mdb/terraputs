[![CI/CD](https://github.com/mdb/terraputs/actions/workflows/main.yml/badge.svg)](https://github.com/mdb/terraputs/actions/workflows/main.yml)

# terraputs

A CLI to generate Terraform outputs documentation.

## What's the use case?

`terraputs` analyzes the contents of a [terraform state](https://www.terraform.io/docs/language/state/index.html)
file and generates Markdown documentation from its [outputs](https://www.terraform.io/docs/language/values/outputs.html).

A common workflow might execute `terraputs -state $(terraform show -json) > outputs.md` after each
invocation of `terraform apply`, then commit `outputs.md` to source control or publish its contents to
offer up-to-date documentation about resources managed by a Terraform project.

<a style="display: block;" href="https://asciinema.org/a/lFUVfdhes0i1cVbtFUvzwLMKd"><img style="width: 500px;" src="demo.svg"></a>

## Usage

A few typical usage examples:

```
# terraputs can accept arbitrary state as a string argument on the command line:
terraputs -state "$(terraform show -json)" -heading "Terraform Outputs"

# or directly from standard input, if the -state option is omitted. to read
# directly from terraform, consider:
terraform show -json | terraputs -heading "Terraform Outputs"

# or read a tfstate file from the filesystem:
terraputs < terraform.tfstate
```

Example output:

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

## More options

```
terraputs -h
Usage of terraputs:
  -heading string
        Optional; the heading text for use in the printed markdown (default "Outputs")
  -state string
        Optional; the state JSON output by 'terraform show -json', read from stdin if omitted
  -state-file string
        Optional; the path to a local file containing 'terraform show -json' output
```

## Example output table formatted by GitHub

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string

## TODO

* provide the ability to pass a custom template
* improve the formatting and readability of lists and maps
* create automated releases
* create a GitHub Action
* publish an OCI image
