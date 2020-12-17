---
page_title: "ejson_file Data Source - terraform-provider-ejson"
subcategory: ""
description: |-
  
---

# Data Source `ejson_file`



## Example Usage

```terraform
data "ejson_file" "config" {
  file        = "secrets.ejson"
  private_key = "12312..." # Optional, reads from keydir by default
}
```

## Schema

### Required

- **file** (String) ejson file to decrypt.

### Optional

- **id** (String) The ID of this resource.
- **private_key** (String, Sensitive) Private key to use for decryption. The provider-level config keydir is used to find a key by default.

### Read-only

- **data** (String, Sensitive) Decrypted contents of ejson file. Use jsondecode to get an object from the JSON blob.
- **map** (Map of String, Sensitive) Mapping of decrypted keys to values, only top-level string values are included. The public key is stripped out and underscore prefixes are removed.


