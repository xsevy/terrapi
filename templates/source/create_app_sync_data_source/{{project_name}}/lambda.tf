resource "null_resource" "install_dependencies" {
  provisioner "local-exec" {
    command = "pip install -r ${local.lambda_source_dir}/requirements.txt -t ${local.lambda_layer_output_dir}/python"
  }

  triggers = {
    always_run = "${timestamp()}"
  }
}

resource "aws_lambda_layer_version" "lambda_layer" {
  filename            = data.archive_file.zip_layer.output_path
  layer_name          = "${local.project_name}-lambda-layer"
  compatible_runtimes = [local.lambda_runtime]
  source_code_hash    = base64sha256(data.archive_file.zip_layer.output_path)
}

resource "aws_lambda_function" "lambda_function" {
  filename         = data.archive_file.zip_the_python_code.output_path
  function_name    = local.project_name
  role             = aws_iam_role.lambda_role.arn
  handler          = "index.lambda_handler"
  runtime          = local.lambda_runtime
  layers           = [aws_lambda_layer_version.lambda_layer.arn]
  source_code_hash = base64sha256(data.archive_file.zip_the_python_code.output_path)
  depends_on = [
    aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role
  ]
}
