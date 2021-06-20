terraform {
  required_providers {
    ops = {
      source = "nanovms/ops"
    }
  }
}

provider "ops" {

}

resource "ops_image_package" "hello" {
    name   = "hello"
    package_name = "node_v14.2.0"
    arguments = ["hello.js"]
    config = "./config.json"
}

output "path" {
  value = ops_image_package.hello.path
}

output "configchecksum" {
  value = ops_image_package.hello.config_checksum
}
