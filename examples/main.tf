terraform {
  required_providers {
    ejson = {
      version = "0.1.0"
      source  = "bouk/ejson"
    }
  }
}

provider "ejson" {
  keydir = "keys"
}

data "ejson_file" "config" {
  file = "secrets.ejson"
}

output "something" {
  value = data.ejson_file.config
}
