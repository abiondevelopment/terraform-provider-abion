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

resource "abion_dns_mx_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      host     = "mail1.example.com."
      priority = "10"
    },
    {
      host     = "mail2.example.com."
      priority = "20"
      comments = "test comment"
    },
    {
      host     = "mail3.example.com."
      priority = "30"
      ttl      = "300"
    },
  ]
}
