resource "aws_iam_role" "lambda_role" {
  name               = "${local.project_name}-lambda-role"
  assume_role_policy = file("${path.module}/lambda_role_policy.json")
}

resource "aws_iam_policy" "iam_policy_for_lambda" {
  name        = "aws-iam-policy-for-${aws_iam_role.lambda_role.name}"
  path        = "/"
  description = "AWS IAM Policy for managing ${aws_iam_role.lambda_role.name}"
  policy      = templatefile("${path.module}/iam_policy_for_lambda.json", {})
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.iam_policy_for_lambda.arn
}
