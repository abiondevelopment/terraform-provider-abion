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

resource "abion_dns_srv_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      target   = "server1.example.com."
      port     = "443"
      priority = "1"
      weight   = "100"
    },
    {
      target   = "server2.example.com."
      port     = "443"
      priority = "100"
      weight   = "10"
      comments = "test comment"
    },
    {
      target   = "server3.example.com."
      port     = "443"
      priority = "100"
      weight   = "30"
      ttl      = "300"
    },
  ]
}
