package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Response struct {
	Query      string    `json:"query"`
	QueryHash  string    `json:"query_hash"`
	Subject    string    `json:"subject"`
	Topics     []string  `json:"topics"`
	NotebookID string    `json:"notebook_id"`
	Response   string    `json:"response"`
	CreatedAt  time.Time `json:"created_at"`
	TTLSeconds int       `json:"ttl_seconds"`
}

type Storage struct {
	client *s3.Client
	bucket string
}

func NewStorage(endpoint, accessKey, secretKey, bucket string) (*Storage, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion("garage"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &Storage{
		client: client,
		bucket: bucket,
	}, nil
}

func NewStorageFromEnv() (*Storage, error) {
	endpoint := os.Getenv("GARAGE_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:3900"
	}

	accessKey := os.Getenv("GARAGE_ACCESS_KEY")
	secretKey := os.Getenv("GARAGE_SECRET_KEY")
	bucket := os.Getenv("GARAGE_BUCKET")
	if bucket == "" {
		bucket = "nlm-cache"
	}

	return NewStorage(endpoint, accessKey, secretKey, bucket)
}

func (s *Storage) generateKey(queryHash, subject string) string {
	now := time.Now()
	return fmt.Sprintf("responses/%d/%02d/%s-%s-%s.json",
		now.Year(), now.Month(), queryHash, subject, now.Format("200601021504"))
}

func (s *Storage) StoreResponse(ctx context.Context, queryHash, subject, query, response string, notebookID string, ttlSeconds int) (string, error) {
	key := s.generateKey(queryHash, subject)

	data := S3Response{
		Query:      query,
		QueryHash:  queryHash,
		Subject:    subject,
		Topics:     []string{},
		NotebookID: notebookID,
		Response:   response,
		CreatedAt:  time.Now(),
		TTLSeconds: ttlSeconds,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(jsonData),
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *Storage) GetResponse(ctx context.Context, key string) (*S3Response, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	var response S3Response
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Storage) DeleteResponse(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *Storage) ListExpiredKeys(ctx context.Context, before time.Time) ([]string, error) {
	var keys []string

	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
		Prefix: aws.String("responses/"),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, obj := range page.Contents {
			if obj.LastModified.Before(before) {
				keys = append(keys, *obj.Key)
			}
		}
	}

	return keys, nil
}

func (s *Storage) BucketExists(ctx context.Context) (bool, error) {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (s *Storage) CreateBucketIfNotExists(ctx context.Context) error {
	exists, err := s.BucketExists(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.bucket),
	})
	return err
}
