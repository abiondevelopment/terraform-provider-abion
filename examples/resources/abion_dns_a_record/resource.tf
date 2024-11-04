terraform {
  required_providers {
    abion = {
      source = "abion/abion"
    }
  }
}

provider "abion" {
  apikey = "<api key>"
}

resource "abion_dns_a_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      ip_address = "203.0.113.0"
    },
    {
      ip_address = "203.0.113.1"
      comments   = "test comment"
    },
    {
      ip_address = "203.0.113.2"
      ttl        = "300"
    },
  ]
}
