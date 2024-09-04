package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var (
	s3Client *s3.S3
	s3Bucket string
	wg       *sync.WaitGroup
)

func init() {
	sess := session.Must(session.NewSession())
	s3Client = s3.New(sess)
	s3Bucket = os.Getenv("S3_BUCKET")
	wg = new(sync.WaitGroup)
}

func main() {
	dir, err := os.Open("./tmp")
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	uploadCh := make(chan struct{}, 100)
	errorCh := make(chan string)
	now := time.Now()
	for {
		files, err := dir.Readdir(1)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			continue
		}

		go func() {
			for {
				select {
				case filename := <-errorCh:
					log.Default().Printf("Retrying upload file: %s\n", filename)
					uploadCh <- struct{}{}
					uploadFile(filename, uploadCh, errorCh)
				}
			}
		}()

		wg.Add(1)
		uploadCh <- struct{}{}
		go uploadFile(
			files[0].Name(),
			uploadCh,
			errorCh,
		)
	}
	wg.Wait()
	log.Default().Printf("finished with success, total time: %v", time.Since(now).Seconds())
}

func uploadFile(filename string, uploadCh <-chan struct{}, errorCh chan<- string) {
	completeFileName := fmt.Sprintf("./tmp/%s", filename)
	log.Default().Printf("Uploading file %s to bucket %s\n", filename, s3Bucket)
	f, err := os.Open(completeFileName)
	if err != nil {
		log.Default().Printf("Error opening file %s: %s\n", completeFileName, err)
		<-uploadCh
		errorCh <- filename
		return
	}
	defer f.Close()
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		log.Default().Printf("Error uploading file %s: %s\n", filename, err)
		<-uploadCh
		errorCh <- filename
		return
	}
	wg.Done()
	<-uploadCh
	log.Default().Printf("Successfully uploaded file %s\n", filename)
}
