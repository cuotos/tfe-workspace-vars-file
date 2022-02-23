When using Terraform Cloud to run TF, the variables are configured in the Workspace on terraform.io.

If you want to run a plan locally, you need to have all the same variables configured in a local tfvars file.

This tool will connecto to app.terraform.io, and output all the variables associated with a workspace, to a file `tfe.auto.tfvars`
which will automatically be picked up and used by `terraform plan | apply`

`go run . -w <workspace ID>`