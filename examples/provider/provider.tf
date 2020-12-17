terraform {
  required_providers {
    ejson = {
      version = "1.1.1"
      source  = "bouk/ejson"
    }
  }
}

provider "ejson" {
  keydir = "keys" # Optional, defaults to $EJSON_KEYDIR or /opt/ejson/keys
}

resource "ejson_keypair" "key" {}

data "ejson_file" "secrets" {
  file        = "secrets.ejson"
  private_key = "12312..." # Optional, reads from keydir by default
}
