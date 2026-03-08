resource "aws_ecs_cluster" "main" {
  name = "${var.project}-cluster"
}

resource "aws_security_group" "app" {
  name = "${var.project}-app-sg"

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
