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

resource "abion_dns_caa_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      flag  = "0"
      tag   = "iodef"
      value = "mailto:webmaster@test.test"
      ttl   = "300"
    },
    {
      flag  = "0"
      tag   = "issue"
      value = "letsencrypt.test"
    },
    {
      flag     = "0"
      tag      = "issuewild"
      value    = ";"
      comments = "test comment"
    },
  ]
}
