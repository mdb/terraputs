# This is a sample Terraform configuration
# used in terraputs tests.
#
# 'make testdata' offers a convenience make
# target to generate terraform.tfstate and
# show.json state state files from this
# configuration.

output "a_string" {
  value = "foo"
}

output "a_sensitive_value" {
  sensitive = true
  value     = "foo"
}

output "a_basic_map" {
  value = {
    foo    = "bar"
    number = 42
  }
}

output "a_nested_map" {
  value = {
    foo    = "bar"
    number = 42
    baz = {
      bar = "baz"
      id  = "123"
    }
  }
}

output "a_list" {
  value = ["foo", "bar"]
}
