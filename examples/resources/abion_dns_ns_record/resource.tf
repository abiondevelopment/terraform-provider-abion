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

resource "abion_dns_ns_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      nameserver = "ns1.testabiondns.se."
    },
    {
      nameserver = "ns2.testabiondns.se."
      comments   = "test comment"
    },
    {
      nameserver = "ns3.testabiondns.se."
      ttl        = "300"
    },
  ]
}
