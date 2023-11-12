terraform {
  backend "s3" {
    bucket         = "{{backend_bucket}}"
    region         = "{{aws_region}}"
    key            = "terraform.tfstate"
    dynamodb_table = "{{backend_lock_table}}"
  }
}
