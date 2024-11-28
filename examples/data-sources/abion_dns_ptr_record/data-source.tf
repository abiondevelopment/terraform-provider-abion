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

data "abion_dns_ptr_record" "example" {
  zone = "example.com"
  name = "10"
}

output "example_ip_addresses" {
  value = data.abion_dns_ptr_record.example.records
}
