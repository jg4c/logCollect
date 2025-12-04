resource "aws_dynamodb_table" "log_table" {
  name           = "log_table"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "service"
  range_key      = "timestamp"

  attribute {
    name = "service"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  tags = {
    Environment = "dev"
    Project     = "log-collect"
  }
}