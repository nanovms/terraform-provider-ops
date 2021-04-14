terraform {
  required_providers {
    ops = {
      version = "0.1"
      source  = "hashicorp.com/ops/ops"
    }
  }
}

provider "ops" {

}

resource "ops_images" "walk_server" {
    name   = "walk-server"
    elf    = "./walk-server"
    config = "./config.json"
}

output "path" {
  value = ops_images.walk_server.path
}
