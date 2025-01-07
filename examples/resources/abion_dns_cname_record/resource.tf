terraform {
  required_providers {
    abion = {
      source = "abiondevelopment/abion"
    }
  }
}

provider "abion" {
  apikey = "<api key>"
}

resource "abion_dns_cname_record" "example" {
  zone = "example.com"
  name = "@"
  record = {
    cname    = "www.devabion.se."
    ttl      = "3600"
    comments = "test comment"
  }
}
