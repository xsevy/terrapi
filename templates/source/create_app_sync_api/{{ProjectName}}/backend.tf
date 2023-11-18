terraform {
  backend "s3" {
    bucket         = "{{.BackendBucket}}"
    region         = "{{.AWSRegion}}"
    key            = "terraform.tfstate"
    dynamodb_table = "{{.BackendLockTable}}"
  }
}
