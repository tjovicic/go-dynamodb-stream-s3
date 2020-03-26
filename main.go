package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.DynamoDBEvent) {
	fmt.Println(event.Records)
}

func main() {
	lambda.Start(handler)
}
