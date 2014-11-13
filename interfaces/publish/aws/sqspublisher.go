package aws

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	sqs "github.com/crowdmob/goamz/sqs"
	"log"
	"time"
)

const (
	PATH_FORMAT = "raw/%s/%04d/%02d/%02d/%s"
)

type SQSPublisher struct {
	queue  *sqs.Queue
	bucket *sss.Bucket
}

func NewSQSPublisher(auth aws.Auth, region aws.Region, bucketName string) *SQSPublisher {
	queue := sqs.New(auth, region)
	return &SQSPublisher{
		queue:  queue,
		bucket: s3.Bucket(bucketName),
	}
}

func (t SQSPublisher) Publish(namespace string, data []byte) {
	payload, err := GzipData(data)
	if err != nil {
		return //err
	}

	md5sum, err := computeMD5(payload)
	if err != nil {
		return //err
	}

	now := time.Now().UTC()
	path := fmt.Sprintf(PATH_FORMAT, namespace, now.Year(), now.Month(), now.Day(), md5sum)
	//return t.bucket.Put(path, data, "application/json", sss.Private, sss.Options{})
	log.Println(path)

	return //nil
}
