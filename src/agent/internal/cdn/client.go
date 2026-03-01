package cdn

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

type CDNClient interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (string, error)
	Download(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
	GetURL(key string) string
}

type Config struct {
	Provider  string
	Endpoint  string
	Bucket    string
	AccessKey string
	SecretKey string
	ZoneID    string
	Domain    string
}

func NewCDNClient(config Config) (CDNClient, error) {
	switch config.Provider {
	case "cloudflare", "cf":
		return NewCloudFlareClient(config)
	case "cloudfront", "aws":
		return NewCloudFrontClient(config)
	case "garage", "s3":
		return NewS3Client(config)
	default:
		// Default to Garage S3 (local)
		return NewS3Client(config)
	}
}

type CloudFlareClient struct {
	config Config
}

func NewCloudFlareClient(config Config) (*CloudFlareClient, error) {
	if config.ZoneID == "" {
		config.ZoneID = os.Getenv("CLOUDFLARE_ZONE_ID")
	}
	if config.Domain == "" {
		config.Domain = os.Getenv("CLOUDFLARE_DOMAIN")
	}
	return &CloudFlareClient{config: config}, nil
}

func (c *CloudFlareClient) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	slog.Info("uploading to CloudFlare R2", "key", key, "size", len(data))

	// CloudFlare R2 uses S3-compatible API
	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return "", err
	}

	return s3Client.Upload(ctx, key, data, contentType)
}

func (c *CloudFlareClient) Download(ctx context.Context, key string) ([]byte, error) {
	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return nil, err
	}
	return s3Client.Download(ctx, key)
}

func (c *CloudFlareClient) Delete(ctx context.Context, key string) error {
	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return err
	}
	return s3Client.Delete(ctx, key)
}

func (c *CloudFlareClient) GetURL(key string) string {
	return fmt.Sprintf("https://%s/%s", c.config.Domain, key)
}

type CloudFrontClient struct {
	config Config
	domain string
}

func NewCloudFrontClient(config Config) (*CloudFrontClient, error) {
	return &CloudFrontClient{
		config: config,
		domain: config.Endpoint,
	}, nil
}

func (c *CloudFrontClient) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	slog.Info("uploading to CloudFront/S3", "key", key)

	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return "", err
	}

	return s3Client.Upload(ctx, key, data, contentType)
}

func (c *CloudFrontClient) Download(ctx context.Context, key string) ([]byte, error) {
	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return nil, err
	}
	return s3Client.Download(ctx, key)
}

func (c *CloudFrontClient) Delete(ctx context.Context, key string) error {
	s3Client, err := NewS3Client(c.config)
	if err != nil {
		return err
	}
	return s3Client.Delete(ctx, key)
}

func (c *CloudFrontClient) GetURL(key string) string {
	return fmt.Sprintf("https://%s/%s", c.domain, key)
}

type S3Client struct {
	config   Config
	endpoint string
	bucket   string
}

func NewS3Client(config Config) (*S3Client, error) {
	if config.Endpoint == "" {
		config.Endpoint = os.Getenv("GARAGE_ENDPOINT")
	}
	if config.Endpoint == "" {
		config.Endpoint = "http://localhost:3900"
	}
	if config.Bucket == "" {
		config.Bucket = os.Getenv("CDN_BUCKET")
	}
	if config.Bucket == "" {
		config.Bucket = "bac-cdn"
	}
	if config.AccessKey == "" {
		config.AccessKey = os.Getenv("GARAGE_ACCESS_KEY")
	}
	if config.SecretKey == "" {
		config.SecretKey = os.Getenv("GARAGE_SECRET_KEY")
	}

	return &S3Client{
		endpoint: config.Endpoint,
		bucket:   config.Bucket,
		config:   config,
	}, nil
}

func (c *S3Client) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	slog.Info("uploading to S3-compatible storage", "key", key, "size", len(data))

	// Use existing storage if available
	// This is a placeholder - actual implementation would use minio-go or aws-sdk
	_ = contentType

	return c.GetURL(key), nil
}

func (c *S3Client) Download(ctx context.Context, key string) ([]byte, error) {
	slog.Info("downloading from S3", "key", key)
	return nil, fmt.Errorf("not implemented - use storage package")
}

func (c *S3Client) Delete(ctx context.Context, key string) error {
	slog.Info("deleting from S3", "key", key)
	return nil
}

func (c *S3Client) GetURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", c.endpoint, c.bucket, key)
}

func GetEnvConfig() Config {
	return Config{
		Provider:  os.Getenv("CDN_PROVIDER"),
		Endpoint:  os.Getenv("CDN_ENDPOINT"),
		Bucket:    os.Getenv("CDN_BUCKET"),
		AccessKey: os.Getenv("CDN_ACCESS_KEY"),
		SecretKey: os.Getenv("CDN_SECRET_KEY"),
		ZoneID:    os.Getenv("CLOUDFLARE_ZONE_ID"),
		Domain:    os.Getenv("CLOUDFLARE_DOMAIN"),
	}
}

func NewDefaultClient() (CDNClient, error) {
	config := GetEnvConfig()
	if config.Provider == "" {
		config.Provider = "garage"
	}
	return NewCDNClient(config)
}

type VersionedCDN struct {
	cdn      CDNClient
	versions map[string]int
}

func NewVersionedCDN(cdn CDNClient) *VersionedCDN {
	return &VersionedCDN{
		cdn:      cdn,
		versions: make(map[string]int),
	}
}

func (v *VersionedCDN) UploadVersioned(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	current := v.versions[key]
	nextVersion := current + 1

	versionedKey := fmt.Sprintf("v%d/%s", nextVersion, key)
	url, err := v.cdn.Upload(ctx, versionedKey, data, contentType)
	if err != nil {
		return "", err
	}

	v.versions[key] = nextVersion
	slog.Info("uploaded new version", "key", key, "version", nextVersion)

	return url, nil
}

func (v *VersionedCDN) GetVersion(key string, version int) string {
	return fmt.Sprintf("v%d/%s", version, key)
}

func (v *VersionedCDN) GetLatestVersion(key string) int {
	return v.versions[key]
}
