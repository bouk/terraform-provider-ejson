terraform {
  required_providers {
    ejson = {
      version = "0.1"
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

output something {
  value = ejson_file.config.data.hi
}
