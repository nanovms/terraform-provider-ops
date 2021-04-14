terraform {
  required_providers {
    ops = {
      version = "0.1"
      source  = "hashicorp.com/nanovms/ops"
    }
  }
}

provider "google" {
  project = "prod-1033"
  region  = "us-west2"
  zone    = "us-west2-a"
}

provider "ops" {

}

resource "ops_images" "walk_server_image" {
  name        = "walk-server"
  elf         = "./walk-server"
  config      = "./config.json"
  targetcloud = "gcp"
}

resource "google_storage_bucket" "images_bucket" {
  name          = "terraform-images"
  location      = "us"
  force_destroy = true
}

resource "google_storage_bucket_object" "walk_server_raw_disk" {
  name   = "walk-server.tar.gz"
  source = ops_images.walk_server_image.path
  bucket = google_storage_bucket.images_bucket.name
}

resource "google_compute_image" "walk_server_image" {
  name = "walk-server-img"

  raw_disk {
    source = google_storage_bucket_object.walk_server_raw_disk.self_link
  }

  labels = {
    "createdby" = "ops"
  }

}

resource "google_compute_instance" "walk_server_instance" {
  name         = "walk-server"
  machine_type = "f1-micro"

  boot_disk {
    initialize_params {
      image = google_compute_image.walk_server_image.self_link
    }
  }

  labels = {
    "createdby" = "ops"
  }

  tags = ["walk-server"]

  network_interface {
    # A default network is created for all GCP projects
    network = "default"
    access_config {
    }
  }
}

resource "google_compute_firewall" "walk_server_firewall" {
  name    = "walk-server-firewall"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["8080"]
  }

  target_tags = ["walk-server"]
}

output "image_path" {
  value = ops_images.walk_server_image.path
}

output "instance_ip" {
  value = google_compute_instance.walk_server_instance.network_interface[0].access_config[0].nat_ip
}
