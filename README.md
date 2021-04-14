# terraform-provider-ops

A terraform provider for OPS. Specify your nanos images and deploy to cloud providers like Google Cloud, Aws, Azure, Oracle Cloud Infrastructure, Open Stack, Vsphere, Upcloud and Digital Ocean.

```
provider "ops" {

}

resource "ops_images" "walk_server_image" {
  name        = "walk-server"
  elf         = "./walk-server"
  config      = "./config.json"
  targetcloud = "gcp"
}
```

## Requirements

* (OPS)[https://github.com/nanovms/ops]

## Build

Run the following command to build the provider

```shell
go build -o terraform-provider-ops
```

## Get Started

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```

If you want to destroy the resources created run `terraform destroy`.

Check our [examples](https://github.com/nanovms/terraform-provider-ops/examples).
