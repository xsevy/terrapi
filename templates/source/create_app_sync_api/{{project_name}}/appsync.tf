resource "aws_appsync_graphql_api" "appsync" {
  name                = "${local.project_name}_appsync"
  schema              = file("schema.graphql")
  authentication_type = "AWS_LAMBDA"

  lambda_authorizer_config {
    authorizer_uri = ""
  }
}


