provider "aws" {
}

locals {
  spire_server_port = 8081
  echo_server_port  = 8443
}

resource "aws_key_pair" "developer" {
  key_name   = "service-id-demo-developer"
  public_key = var.developer_public_key
}

output "client_role_arn" {
  value = aws_iam_role.client.arn
}

output "client_public_ip" {
  value = aws_instance.client.public_ip
}

output "server_role_arn" {
  value = aws_iam_role.server.arn
}

output "server_public_ip" {
  value = aws_instance.server.public_ip
}

output "spire_server_public_ip" {
  value = aws_instance.spire_server.public_ip
}
