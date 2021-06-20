terraform {
  required_providers {
    ops = {
      source = "nanovms/ops"
    }
  }
}

provider "ops" {

}

resource "ops_image_executable" "walk_server" {
    name   = "walk-server"
    elf    = "./walk-server"
    config = "./config.json"
}

output "path" {
  value = ops_image_executable.walk_server.path
}

output "configchecksum" {
  value = ops_image_executable.walk_server.config_checksum
}

output "elfchecksum" {
  value = ops_image_executable.walk_server.elf_checksum
}
