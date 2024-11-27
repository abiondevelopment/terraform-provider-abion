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

data "abion_dns_txt_record" "example" {
  zone = "example.com"
  name = "www"
}

output "example_txt_data" {
  value = data.abion_dns_txt_record.example.records
}