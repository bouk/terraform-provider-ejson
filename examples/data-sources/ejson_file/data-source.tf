data "ejson_file" "config" {
  file        = "secrets.ejson"
  private_key = "12312..." # Optional, reads from keydir by default
}
