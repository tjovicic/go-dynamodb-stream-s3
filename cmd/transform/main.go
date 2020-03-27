package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(e events.KinesisFirehoseEvent) (events.KinesisFirehoseResponse, error) {
	var response events.KinesisFirehoseResponse

	for _, record := range e.Records {
		transformedRecord := responseRecord(record)
		response.Records = append(response.Records, transformedRecord)
	}

	return response, nil
}

func responseRecord(r events.KinesisFirehoseEventRecord) events.KinesisFirehoseResponseRecord {
	var transformedRecord events.KinesisFirehoseResponseRecord

	transformedRecord.RecordID = r.RecordID
	transformedRecord.Result = events.KinesisFirehoseTransformedStateOk
	transformedRecord.Data = []byte("test")

	return transformedRecord
}

func main() {
	lambda.Start(handler)
}
