---
page_title: "ejson Provider"
subcategory: ""
description: |-
  
---

# ejson Provider



## Example Usage

```terraform
terraform {
  required_providers {
    ejson = {
      version = "0.1.0"
      source  = "bouk/ejson"
    }
  }
}

provider "ejson" {
  keydir = "keys" # Optional, defaults to $EJSON_KEYDIR or /opt/ejson/keys
}
```

## Schema

### Optional

- **keydir** (String) Directory to read private keys from. Defaults to $EJSON_KEYDIR or /opt/ejson/keys if not set.
