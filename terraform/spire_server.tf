resource "aws_instance" "spire_server" {
  ami                         = "ami-04cb4ca688797756f" # Amazon Linux 2023 AMI
  instance_type               = "t2.micro"
  key_name                    = aws_key_pair.developer.key_name
  vpc_security_group_ids      = [aws_security_group.spire_server.id]
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.spire_server.name

  tags = {
    Name = "service-id-demo-spire-server"
  }
}

resource "aws_security_group" "spire_server" {
  name   = "service-id-demo-spire-server"
  vpc_id = var.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "TCP"
    cidr_blocks = var.developer_cidrs
  }

  ingress {
    from_port       = local.spire_server_port
    to_port         = local.spire_server_port
    protocol        = "TCP"
    security_groups = [aws_security_group.client.id, aws_security_group.server.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_instance_profile" "spire_server" {
  name = "service-id-demo-spire-server"
  role = aws_iam_role.spire_server.name
}

resource "aws_iam_role" "spire_server" {
  name               = "service-id-demo-spire-server"
  assume_role_policy = data.aws_iam_policy_document.spire_server_trust.json

  inline_policy {
    name   = "node-attestation"
    policy = data.aws_iam_policy_document.spire_server_node_attestation.json
  }
}

data "aws_iam_policy_document" "spire_server_trust" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "spire_server_node_attestation" {
  statement {
    actions = [
      "ec2:DescribeInstances",
      "iam:GetInstanceProfile"
    ]

    resources = ["*"]
  }
}

