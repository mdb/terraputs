# This is a sample Terraform configuration
# used in terraputs tests.
#
# 'make testdata' offers a convenience make
# target to generate terraform.tfstate and
# show.json state state files from this
# configuration.

data "template_file" "greeting" {
  template = <<-EOT
  ${var.greeting}
  EOT
}

resource "local_file" "greeting" {
  content  = data.template_file.greeting.rendered
  filename = "${path.module}/greeting.txt"
}
