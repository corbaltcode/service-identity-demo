resource "aws_instance" "client" {
  ami                         = "ami-04cb4ca688797756f" # Amazon Linux 2023 AMI
  instance_type               = "t2.micro"
  key_name                    = aws_key_pair.developer.key_name
  vpc_security_group_ids      = [aws_security_group.client.id]
  subnet_id                   = var.subnet_id
  associate_public_ip_address = true
  iam_instance_profile        = aws_iam_instance_profile.client.name

  tags = {
    Name = "service-id-demo-client"
  }
}

resource "aws_security_group" "client" {
  name   = "service-id-demo-client"
  vpc_id = var.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "TCP"
    cidr_blocks = var.developer_cidrs
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_iam_instance_profile" "client" {
  name = "service-id-demo-client"
  role = aws_iam_role.client.name
}

resource "aws_iam_role" "client" {
  name               = "service-id-demo-client"
  assume_role_policy = data.aws_iam_policy_document.client_trust.json
}

data "aws_iam_policy_document" "client_trust" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ec2.amazonaws.com"]
    }
  }
}
