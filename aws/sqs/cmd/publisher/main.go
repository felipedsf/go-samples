package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

type SqsClient struct {
	queueUrl string
	sqs      *sqs.SQS
}

func NewSqsClient(queueUrl string, sqs *sqs.SQS) *SqsClient {
	return &SqsClient{
		queueUrl: queueUrl,
		sqs:      sqs,
	}
}

func (c SqsClient) SendMessage(msg string) {
	log.Println("Publishing messages on sqs queue...")
	result, err := c.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(c.queueUrl),
	})
	if err != nil {
		panic(err)
	}
	log.Printf("Message published: %s", *result.MessageId)
}

func main() {
	sess := session.Must(session.NewSession())
	client := NewSqsClient(os.Getenv("SQS_QUEUE_URL"), sqs.New(sess))

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("{\"message\": \"%d - it's all right, from sqs! time: %s\"}", i, time.Now().Format(time.DateTime))
		client.SendMessage(msg)
	}

}
