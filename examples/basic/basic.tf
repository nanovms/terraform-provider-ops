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
