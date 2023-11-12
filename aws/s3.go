package aws

import (
	s3_sdk "github.com/aws/aws-sdk-go/service/s3"
)

type s3 struct {
	client *s3_sdk.S3
}

type S3 interface {
	ListBuckets() ([]string, error)
}

func NewS3(aws AWS) S3 {
	return &s3{
		client: s3_sdk.New(aws.GetSession()),
	}
}

func (s *s3) ListBuckets() ([]string, error) {
	var buckets []string

	result, err := s.client.ListBuckets(nil)
	if err != nil {
		return buckets, err
	}

	for _, b := range result.Buckets {
		buckets = append(buckets, *b.Name)
	}

	return buckets, nil
}
