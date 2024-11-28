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

data "abion_dns_cname_record" "example" {
  zone = "example.com"
  name = "@"
}

output "example_ip_addresses" {
  value = data.abion_dns_cname_record.example.record
}
