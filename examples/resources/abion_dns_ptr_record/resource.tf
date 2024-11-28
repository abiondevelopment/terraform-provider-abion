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

resource "abion_dns_ptr_record" "example" {
  zone = "example.com"
  name = "10"
  records = [
    {
      ptr = "www.abiontest.com."
    },
    {
      ptr      = "www.abiontest2.com."
      comments = "test comment"
    },
    {
      ptr = "www.abiontest3.com."
      ttl = "300"
    },
  ]
}
