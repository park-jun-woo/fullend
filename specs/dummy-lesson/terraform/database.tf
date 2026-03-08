resource "aws_db_instance" "main" {
  identifier             = "${var.project}-db"
  engine                 = "postgres"
  engine_version         = "16.4"
  instance_class         = "db.t4g.micro"
  allocated_storage      = 20
  db_name                = "lesson"
  username               = "admin"
  password               = var.db_password
  skip_final_snapshot    = true
  publicly_accessible    = false
  vpc_security_group_ids = [aws_security_group.db.id]
}

variable "db_password" {
  sensitive = true
}

resource "aws_security_group" "db" {
  name = "${var.project}-db-sg"

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.app.id]
  }
}
