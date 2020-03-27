# Example of using AWS Lambda and DynamoDB Streams

## Quick start
`make deploy`

### Description
Lambda function defined in `main.go` listens to DynamoDB stream and writes them to S3

### Techstack:
- serverless
- AWS Lambda
- DynamoDB
- Go
- CloudFormation

### References
https://serverless.com/framework/docs/providers/aws/events/streams/

https://itonaut.com/2019/02/06/implement-kinesis-firehose-s3-delivery-preprocessed-by-lambda-in-aws-cloudformation/

https://itiskj.hatenablog.com/entry/2018/08/28/121536


https://stackoverflow.com/questions/44032664/reference-function-from-within-serverless-yml