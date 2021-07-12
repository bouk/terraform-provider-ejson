resource "ejson_keypair" "key" {}

resource "ejson_file" "secrets" {
  data = jsonencode({
    "hello" => "nice to meet you"
  })
  public_key = ejson_keypair.key.public_key
}

output "encrypted" {
  value = ejson_file.secrets.encrypted
}
