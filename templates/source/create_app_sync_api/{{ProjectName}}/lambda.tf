data "aws_lambda_function" "authorizer" {
  function_name = "{{.AuthorizerLambdaFunction}}"
}

resource "aws_lambda_permission" "appsync_lambda_authorizer" {
  statement_id  = "${local.project_name}-appsync_lambda_authorizer"
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.authorizer.function_name
  principal     = "appsync.amazonaws.com"
  source_arn    = aws_appsync_graphql_api.appsync.arn
}
