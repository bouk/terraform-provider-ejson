---
page_title: "ejson_keypair Resource - terraform-provider-ejson"
subcategory: ""
description: |-
  Generate an ejson keypair.
---

# Resource `ejson_keypair`

Generate an ejson keypair.

## Example Usage

```terraform
resource "ejson_keypair" "key" {}

data "ejson_file" "secrets" {
  file        = "secrets.ejson"
  private_key = ejson_keypair.key.private_key
}
```

## Schema

### Optional

- **id** (String) The ID of this resource.

### Read-only

- **private_key** (String, Sensitive) The private part of the key pair, for decrypting the ejson file.
- **public_key** (String) The public part of the key pair, for embedding into the ejson file.

## Import

Import is supported using the following syntax:

```shell
$ terraform import ejson_keypair.key 965aed709d63e22fd7cb6b4ee8f317bb0ad99a07fa68bc24bdcc349d0b2af130
```
