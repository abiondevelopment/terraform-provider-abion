---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "abion_dns_ptr_record Data Source - abion"
subcategory: ""
description: |-
  Use this data source to get DNS PTR records of the zone.
---

# abion_dns_ptr_record (Data Source)

Use this data source to get DNS PTR records of the zone.

## Example Usage

```terraform
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

data "abion_dns_ptr_record" "example" {
  zone = "example.com"
  name = "10"
}

output "example_records" {
  value = data.abion_dns_ptr_record.example.records
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name to fetch records for. For example `@`, `www`, `ftp`, `www.east`. The `@` character represents the root of the zone.
- `zone` (String) The zone the record belongs to.

### Read-Only

- `records` (Attributes List) The list of PTR records. Records are sorted to avoid constant changing plans (see [below for nested schema](#nestedatt--records))

<a id="nestedatt--records"></a>
### Nested Schema for `records`

Required:

- `ptr` (String) The canonical name this record will point to.

Optional:

- `comments` (String) Comments for the record.
- `ttl` (Number) Time-to-live (TTL) for the record, in seconds.
