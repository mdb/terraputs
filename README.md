[![CI/CD](https://github.com/mdb/terraputs/actions/workflows/main.yml/badge.svg)](https://github.com/mdb/terraputs/actions/workflows/main.yml)

# terraputs

A CLI to generate Terraform outputs documentation.

## What's the use case?

Invoke `terraputs -state $(terraform show -json) > outputs.md` after each invocation of `terraform apply` and commit `outputs.md` to source control for up-to-date Terraform state documentation-as-code.

## Usage

Basic usage:

```
terraputs \
  -state $(terraform show -json)
```

Example output:

```
# Outputs

Terraform state outputs.

| Output | Value | Type
| --- | --- | --- |
| a_basic_map | map[foo:bar number:42] | map[string]interface {}
| a_list | [foo bar] | []interface {}
| a_nested_map | map[baz:map[bar:baz id:123] foo:bar number:42] | map[string]interface {}
| a_sensitive_value | sensitive; redacted | string
| a_string | foo | string
```

Example output table formatted by GitHub:

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
