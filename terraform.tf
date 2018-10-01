resource "aws_iam_role" "cloudwatch_cleaner_iam_role" {
  name = "cloudwatch_cleaner_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_lambda_function" "cloudwatch_cleaner_lambda" {
  filename         = "cloudwatch-cleaner.zip"
  function_name    = "cloudwatch-cleaner"
  role             = "${aws_iam_role.cloudwatch_cleaner_iam_role.arn}"
  handler          = "cloudwatch-cleaner"
  source_code_hash = "${base64sha256(file("cloudwatch-cleaner.zip"))}"
  runtime          = "go1.x"
  memory_size      = 128
  environment {
    variables = {
      RETENTION_DAYS = "${var.retention_days}"
    }
  }
}

data "aws_iam_policy" "AmazonEC2ReadOnlyAccess" {
  arn = "arn:aws:iam::aws:policy/AmazonEC2ReadOnlyAccess"
}

data "aws_iam_policy" "CloudWatchFullAccess" {
  arn = "arn:aws:iam::aws:policy/CloudWatchFullAccess"
}

resource "aws_iam_role_policy_attachment" "cloudwatch_cleaner_role_AmazonEC2ReadOnlyAccess_policy_attach" {
  role = "${aws_iam_role.cloudwatch_cleaner_iam_role.name}"
  policy_arn = "${data.aws_iam_policy.AmazonEC2ReadOnlyAccess.arn}"
}

resource "aws_iam_role_policy_attachment" "cloudwatch_cleaner_role_CloudWatchFullAccess_policy_attach" {
  role = "${aws_iam_role.cloudwatch_cleaner_iam_role.name}"
  policy_arn = "${data.aws_iam_policy.CloudWatchFullAccess.arn}"
}

resource "aws_cloudwatch_event_rule" "cloudwatch_cleaner_rule" {
  name                = "cloudwatch_cleaner_rule"
  schedule_expression = "${var.lambda_rate}"
  is_enabled          = true
}

resource "aws_cloudwatch_event_target" "cloudwatch_cleaner_target" {
  rule      = "${aws_cloudwatch_event_rule.cloudwatch_cleaner_rule.name}"
  arn       = "${aws_lambda_function.cloudwatch_cleaner_lambda.arn}"
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_cloudwatch_cleaner" {
    statement_id = "AllowExecutionFromCloudWatch"
    action = "lambda:InvokeFunction"
    function_name = "${aws_lambda_function.cloudwatch_cleaner_lambda.function_name}"
    principal = "events.amazonaws.com"
    source_arn = "${aws_cloudwatch_event_rule.cloudwatch_cleaner_rule.arn}"
}

variable "retention_days" {
  default = "30"  
}

variable "lambda_rate" {
  default = "rate(1 day)" 
}
