terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-northeast-2"
}

resource "aws_ecs_cluster" "gigbridge" {
  name = "gigbridge-cluster"
}

resource "aws_rds_cluster" "gigbridge" {
  cluster_identifier = "gigbridge-db"
  engine             = "aurora-postgresql"
  engine_version     = "15.4"
  database_name      = "gigbridge"
  master_username    = "postgres"
  master_password    = var.db_password
}

variable "db_password" {
  type      = string
  sensitive = true
}
