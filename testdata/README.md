# testdata

The `testdata/**/terraform.tfstate` and `testdata/**/show.json` files in this directory are generated
from the `testdata/**/main.tf` Terraform configuration via `make tfstate`. The `testdata/**/show.json`
files are used by `terraputs`'s tests.
