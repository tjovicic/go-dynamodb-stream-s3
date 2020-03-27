package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

var sess = session.Must(session.NewSession())
var firehoseClient = firehose.New(sess)

func handler(ctx context.Context, e events.DynamoDBEvent) {
	var rs []*firehose.Record

	for _, r := range e.Records {
		var data []byte

		data, err := json.Marshal(r.Change.NewImage)
		if err != nil {
			log.Panicln(err)
		}

		data = append(data, "\n"...) // add new line between records

		rs = append(rs, &firehose.Record{Data: data})
	}

	out, err := firehoseClient.PutRecordBatch(&firehose.PutRecordBatchInput{
		Records:            rs,
		DeliveryStreamName: aws.String(os.Getenv("STREAM_NAME")),
	})

	log.Println(out)

	if err != nil {
		log.Panicln(err)
	}

	return
}

func main() {
	lambda.Start(handler)
}
