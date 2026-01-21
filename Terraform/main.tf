data "aws_ami" "amazon_linux"{
    most_recent = true

    filter{
        name    = "name"
        values  = ["amzn2-ami-hvm-*-x86_64-gp2"]
    }

    owners = ["amazon"]
}

resource "aws_instance" "go_service_vm"{
    ami             = data.aws_ami.amazon_linux.id
    instance_type   = var.instance_type
    key_name        = var.key_name

    tags = {
        Name            = "go-web-service-vm"
        Environment     = "dev"
    }
}