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

resource "abion_dns_txt_record" "example" {
  zone = "example.com"
  name = "www"
  records = [
    {
      txt_data = "txt 1"
    },
    {
      txt_data = "txt 2"
      comments = "test comment"
    },
    {
      txt_data = "txt 3"
      ttl      = "300"
    },
  ]
}
