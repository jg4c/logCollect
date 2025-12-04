resource "aws_kinesis_stream" "log_stream" {
  name        = "log_stream"
  shard_count = 1

  # Retention config
  retention_period = 168

  # Monitoring config
  shard_level_metrics = [
    "IncomingBytes",
    "OutgoingBytes",
    "IncomingRecords",
    "OutgoingRecords"
  ]

  tags = {
    Environment = "dev"
    Project     = "log-collect"
  }
}
