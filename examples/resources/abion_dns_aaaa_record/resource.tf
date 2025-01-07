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

resource "abion_dns_aaaa_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff0"
    },
    {
      ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff1"
      comments   = "test comment"
    },
    {
      ip_address = "2001:db8:ffff:ffff:ffff:ffff:ffff:fff2"
      ttl        = "300"
    },
  ]
}
