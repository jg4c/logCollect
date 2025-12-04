import json
import base64
import os
from datetime import datetime

import boto3

dynamodb = boto3.resource("dynamodb")
table_name = os.environ.get("LOGS_TABLE_NAME", "logs-table")
table = dynamodb.Table(table_name)

def lambda_handler(event, context):
    # event["Records"] contains all the Kinesis records in this batch
    for record in event["Records"]:
        # Kinesis data is base64 encoded
        payload = record["kinesis"]["data"]
        data_bytes = base64.b64decode(payload)
        log_event = json.loads(data_bytes.decode("utf-8"))

        # Extract fields, apply defaults
        service = log_event.get("service", "unknown")
        level = log_event.get("level", "INFO")
        message = log_event.get("message", "")

        timestamp = log_event.get("timestamp")
        if not timestamp:
            # If no timestamp, use now in ISO-8601
            timestamp = datetime.utcnow().isoformat() + "Z"

        item = {
            "service": service,
            "timestamp": timestamp,
            "level": level,
            "message": message,
        }

        # Write to DynamoDB
        table.put_item(Item=item)

    return {"statusCode": 200}
