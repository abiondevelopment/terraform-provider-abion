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

data "abion_dns_ns_record" "example" {
  zone = "example.com"
  name = "www"
}

output "example_records" {
  value = data.abion_dns_ns_record.example.records
}