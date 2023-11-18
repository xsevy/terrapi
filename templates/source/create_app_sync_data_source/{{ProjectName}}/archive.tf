data "archive_file" "zip_the_python_code" {
  type        = "zip"
  source_dir  = "${local.lambda_source_dir}/"
  output_path = local.lambda_zip_path
}

data "archive_file" "zip_layer" {
  type        = "zip"
  source_dir  = "${local.lambda_layer_output_dir}/"
  output_path = local.lambda_layer_zip_path
  depends_on  = [null_resource.install_dependencies]
}
