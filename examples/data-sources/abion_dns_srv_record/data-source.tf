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

data "abion_dns_srv_record" "example" {
  zone = "example.com"
  name = "www"
}

output "example_records" {
  value = data.abion_dns_srv_record.example.records
}