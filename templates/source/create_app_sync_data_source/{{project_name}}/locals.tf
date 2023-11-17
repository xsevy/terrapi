locals {
  project_name = "{{.project_name}}"

  lambda_zip_file_name       = "lambda.zip"
  lambda_zip_path            = "${path.module}/${local.lambda_zip_file_name}"
  lambda_source_dir          = "${path.module}/lambda"
  lambda_layer_zip_file_name = "layer.zip"
  lambda_layer_zip_path      = "${path.module}/${local.lambda_layer_zip_file_name}"
  lambda_layer_output_dir    = "${path.module}/lambda_layer_files"
  lambda_runtime             = "{{.lambda_runtime}}"
}
