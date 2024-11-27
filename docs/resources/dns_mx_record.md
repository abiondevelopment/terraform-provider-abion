---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "abion_dns_mx_record Resource - abion"
subcategory: ""
description: |-
  Use this resource to create, update and delete DNS MX records of a zone.
---

# abion_dns_mx_record (Resource)

Use this resource to create, update and delete DNS MX records of a zone.

## Example Usage

```terraform
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

resource "abion_dns_mx_record" "example" {
  zone = "example.com"
  name = "www"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name to create records for. For example `@`, `www`, `ftp`, `www.east`. The `@` character represents the root of the zone.
- `records` (Attributes List) The list of MX records. Records are sorted to avoid constant changing plans (see [below for nested schema](#nestedatt--records))
- `zone` (String) The zone the record belongs to.

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Required:

- `host` (String) The hostname of the mail server.
- `priority` (Number) The priority in which order mail servers are tried.

Optional:

- `comments` (String) Comments for the record.
- `ttl` (Number) Time-to-live (TTL) for the record, in seconds.

## Import

Import is supported using the following syntax:

```shell
# DNS MX records can be imported by specifying the string identifier. The import ID should be in the format: "zone/name". The `@` character represents the root of the zone, E.g., "example.com/@"
terraform import abion_dns_mx_record.example "example.com/www"
```