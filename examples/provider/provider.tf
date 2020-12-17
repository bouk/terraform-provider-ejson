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
