# Terraform Provider Ejson

Run the following command to build the provider

```shell
go build -o terraform-provider-ejson
```

## Example

```terraform
terraform {
  required_providers {
    ejson = {
      version = "1.1.0"
      source  = "bouk/ejson"
    }
  }
}

provider "ejson" {
  keydir = "keys" # Optional, defaults to $EJSON_KEYDIR or /opt/ejson/keys
}

data "ejson_file" "config" {
  file = "secrets.ejson"
  private_key = "12312..." # Optional, reads from keydir by default
}
```

`data.ejson_file.config` will contain a `data` attribute containing the decrypted JSON, and a `map` attribute with all the string key: values. 
