package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

var sess = session.Must(session.NewSession())
var firehoseClient = firehose.New(sess)

const MaxWaitInterval = 100
const MaxRetries = 3

func handler(ctx context.Context, e events.DynamoDBEvent) {
	var rs []*firehose.Record

	for _, r := range e.Records {
		rs = append(rs, &firehose.Record{Data: eventsToByte(r)})
	}

	if len(rs) == 0 {
		log.Println("empty event record. exiting...")
		return
	}

	out, err := firehoseClient.PutRecordBatch(&firehose.PutRecordBatchInput{
		Records:            rs,
		DeliveryStreamName: aws.String(os.Getenv("STREAM_NAME")),
	})

	if err != nil {
		log.Panicln(err)
	}

	if *out.FailedPutCount == 0 {
		log.Println("success!")
		return
	}

	// retry failed put records
	if err := retry(out, rs); err != nil {
		log.Panicln(err)
	}

	return
}

func eventsToByte(r events.DynamoDBEventRecord) []byte {
	var data []byte

	data, err := json.Marshal(r.Change.NewImage)
	if err != nil {
		log.Panicln(err)
	}

	data = append(data, "\n"...) // add new line between records
	return data
}

func retry(out *firehose.PutRecordBatchOutput, rs []*firehose.Record) error {
	retries := 0
	retry := true

	for retry && retries < MaxRetries {
		waitTime := min(getWaitTimeExp(retries), MaxWaitInterval)

		time.Sleep(time.Duration(waitTime))

		var retryRecords []*firehose.Record
		for i, r := range out.RequestResponses {
			if r.ErrorCode != nil {
				retryRecords = append(retryRecords, rs[i])
			}
		}

		out, err := firehoseClient.PutRecordBatch(&firehose.PutRecordBatchInput{
			Records:            retryRecords,
			DeliveryStreamName: aws.String(os.Getenv("STREAM_NAME")),
		})

		if err != nil {
			return err
		}

		if *out.FailedPutCount == 0 {
			retry = false
		}

		retries++
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func getWaitTimeExp(retryCount int) int {
	waitTime := math.Pow(2, float64(retryCount)) * 100
	return int(waitTime)
}

func main() {
	lambda.Start(handler)
}
