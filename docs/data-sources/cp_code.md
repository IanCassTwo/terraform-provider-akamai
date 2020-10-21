---
layout: "akamai"
page_title: "Akamai: cp_code"
subcategory: "Provisioning"
description: |-
 CP Code
---

# akamai_cp_code


Use `akamai_cp_code` data source to retrieve a cpcode id.

## Example Usage

Basic usage:

```hcl
data "akamai_cp_code" "example" {
     name = "cpcode name"
     group = "grp_#####"
     contract = "ctr_#####"
}

resource "akamai_property" "example" {
    contract = "${data.akamai_cpcode.example.id}"
    ...
}
```

## Argument Reference

The following arguments are supported:

* `name` — (Required) The CP code name.
* `group` — (Required) The group ID
* `contract` — (Required) The contract ID

## Attributes Reference

The following are the return attributes:

* `id` — The CP code ID.
* `product_ids` - An array of product ids associated with this cpcode
