variable "key_name" {
    default = "key"
}

resource "tls_private_key" "pk" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = var.key_name
  public_key = tls_private_key.pk.public_key_openssh

  provisioner "local-exec" { 
    command = "echo '${tls_private_key.pk.private_key_pem}' > ./${var.key_name}.pem"
  }

  provisioner "local-exec" { 
    command = "chmod 400 ${var.key_name}.pem"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.nano"
  key_name      = aws_key_pair.generated_key.key_name
  vpc_security_group_ids = [aws_security_group.main.id]

  connection {
    host        = self.public_ip
    user        = "ubuntu"
    type        = "ssh"
    private_key = file("${aws_key_pair.generated_key.key_name}.pem")
    timeout     = "5m"
  }

    provisioner "remote-exec" {
        inline = [
            "mkdir app",
        ]
    }

    # provisioner "file" {
    #     source      = "../api/.env"
    #     destination = "~/app/api/.env"
    # }

    provisioner "file" {
        source      = "../go.mod"
        destination = "~/app/go.mod"
    }

    provisioner "file" {
        source      = "../go.sum"
        destination = "~/app/go.sum"
    }

    provisioner "file" {
        source      = abspath("../api")
        destination = "~/app"
    }

    provisioner "remote-exec" {
        inline = [
            "sudo snap install go --classic",
            "cd app/api",
            "go build -o api",
            "sudo su <<EOF",
            "nohup ./api -port :80 &",
            "sleep 5",
            "EOF",
            "touch finished.txt"
        ]
    }
}


resource "aws_security_group" "main" {
  name = "main"

  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = [var.anyCidr]
  }

  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = [var.anyCidr]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = [var.anyCidr]
    ipv6_cidr_blocks = ["::/0"]
  }
}

