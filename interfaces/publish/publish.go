package publish

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	sss "github.com/crowdmob/goamz/s3"
	"log"
	"time"
)

const (
	PATH_FORMAT = "raw/%s/%04d/%02d/%02d/%s"
)

type S3Publisher struct {
	s3     *sss.S3
	bucket *sss.Bucket
}

func NewS3Publisher(auth aws.Auth, region aws.Region, bucketName string) *S3Publisher {
	s3 := sss.New(auth, region)
	return &S3Publisher{
		s3:     s3,
		bucket: s3.Bucket(bucketName),
	}
}

func (t S3Publisher) Publish(namespace string, data []byte) {
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

func GzipData(data []byte) ([]byte, error) {
	buffer := new(bytes.Buffer)

	writer := gzip.NewWriter(buffer)
	if _, err := writer.Write(data); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func computeMD5(data []byte) (string, error) {
	h := md5.New()
	if _, err := h.Write(data); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
