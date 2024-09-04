package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func NewAwsSqsListener(workers int, queue string, client *sqs.SQS) *AwsSqsListener {
	return &AwsSqsListener{
		workers:   workers,
		wg:        &sync.WaitGroup{},
		queueUrl:  queue,
		sqsClient: client,
	}
}

type AwsSqsListener struct {
	workers   int
	wg        *sync.WaitGroup
	queueUrl  string
	sqsClient *sqs.SQS
}

func (l AwsSqsListener) receiveMessageWorker(id int) {
	log.Printf("starting worker %d", id)
	for {
		result, err := l.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: &l.queueUrl,
		})
		if err != nil {
			log.Fatalf("Error worker %d on receiving message, %s", id, err)
		}

		for _, message := range result.Messages {
			log.Printf("Worker %d, message reveiced: %s", id, *message.Body)
			_, err = l.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &l.queueUrl,
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				log.Fatalf("Error worker %d on delete message, %s", id, err)
			}
		}
	}
}

func (l AwsSqsListener) Listen() {
	for i := 0; i < l.workers; i++ {
		l.wg.Add(1)
		go l.receiveMessageWorker(i)
	}
	l.wg.Wait()
}

func (l AwsSqsListener) Shutdown() {
	log.Printf("shutting down aws sqs listener")
	for i := 0; i < l.workers; i++ {
		l.wg.Done()
	}
}

func main() {
	done := make(chan bool)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	sess := session.Must(session.NewSession())
	sqsListener := NewAwsSqsListener(3, os.Getenv("SQS_QUEUE_URL"), sqs.New(sess))
	go sqsListener.Listen()
	go func() {
		sig := <-interrupt
		log.Printf("Signal intercepted: %v\n", sig.String())
		sqsListener.Shutdown()
		done <- true
	}()
	<-done
}
