resource "ejson_keypair" "key" {}

data "ejson_file" "secrets" {
  file        = "secrets.ejson"
  private_key = ejson_keypair.key.private_key
}
