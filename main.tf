resource "aws_instance" "showcase" {
  instance_type = var.ec2-instance
  resource_id   = "showcase"
}

resource "aws_instance" "showcase-2" {
  instance_type = "cr1.8xlarge"
  resource_id   = "showcase-2"
}
