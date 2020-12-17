terraform {
  required_providers {
    ejson = {
      version = "1.1.1"
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
