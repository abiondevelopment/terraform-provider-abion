terraform {
  required_providers {
    abion = {
      source = "abion/abion"
    }
  }
}

provider "abion" {
  #dev
  # apikey = "+FOQDv1eG4BBFmYV3WwxtTI0pKZSK+g2rF+F9fmZGWJRZj/0qzM51ZYbcRl0vZuM1Hv9dbJj7eBmRG8ijNWASA=="
  #demo
  apikey = "p/cEqe8kBSF68Ft8+I7H39FmBALgaAswVXRog7Wo7ec="
}

resource "abion_dns_a_record" "example_x" {
  zone = "pmapitest10.com"
  name = "test20241203-1"
  records = [
    {
      ip_address = "203.0.113.0"
      comments   = "test comment"
      ttl        = "300"
    }
  ]
}

resource "abion_dns_txt_record" "example_x" {
  zone = "pmapitest10.com"
  name = "test20241205"
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

resource "abion_dns_cname_record" "example_x" {
  zone = "pmapitest10.com"
  name = "testcname20241205"
  record = {
    cname    = "www.devabion.se."
    ttl      = "3600"
    comments = "test comment"
  }
}


resource "abion_dns_mx_record" "example_x" {
  zone = "pmapitest10.com"
  name = "testzone20241205"
  records = [
    {
      host     = "mail1.example.com."
      priority = "10"
    },
    {
      host     = "mail2.example.com."
      priority = "20"
      comments = "test comment"
    },
    {
      host     = "mail3.example.com."
      priority = "30"
      ttl      = "300"
    },
  ]
}