# Skill: Garage S3 Storage

## Purpose

S3-compatible storage using Garage for BAC Unified.

## When to use

- Storing uploaded files (PDFs, images)
- Saving generated notebooks
- Caching animation outputs

## Why Garage?

- S3-compatible API
- Runs locally (development)
- Lightweight (Rust)
- Production-ready

## Go Integration

```go
import (
    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func NewStorage() (*minio.Client, error) {
    endpoint := os.Getenv("GARAGE_ENDPOINT") // localhost:3900
    client, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(
            os.Getenv("GARAGE_ACCESS_KEY"),
            os.Getenv("GARAGE_SECRET_KEY"),
            "",
        ),
        Secure: false, // HTTP for local
    })
    return client, err
}
```

## Operations

### Upload

```go
func Upload(ctx context.Context, client *minio.Client, key string, data []byte) error {
    _, err := client.PutObject(ctx, "bac-resources", key,
        bytes.NewReader(data), int64(len(data)),
        minio.PutObjectOptions{ContentType: "application/json"})
    return err
}
```

### Download

```go
func Download(ctx context.Context, client *minio.Client, key string) ([]byte, error) {
    obj, err := client.GetObject(ctx, "bac-resources", key, minio.GetObjectOptions{})
    defer obj.Close()
    
    data, err := io.ReadAll(obj)
    return data, err
}
```

### Delete

```go
func Delete(ctx context.Context, client *minio.Client, key string) error {
    return client.RemoveObject(ctx, "bac-resources", key, minio.RemoveObjectOptions{})
}
```

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `GARAGE_ENDPOINT` | Garage server (localhost:3900) |
| `GARAGE_ACCESS_KEY` | Access key |
| `GARAGE_SECRET_KEY` | Secret key |
| `GARAGE_BUCKET` | Bucket name (bac-resources) |

## Bucket Structure

```
bac-resources/
├── uploads/
│   └── {year}/{month}/{day}/{uuid}.{ext}
├── notebooks/
│   └── {user_id}/{notebook_id}.json
├── responses/
│   └── {year}/{month}/{query_hash}.json
└── animations/
    └── {year}/{month}/{id}.mp4
```

## Local Development

```bash
# Run Garage
garage server -g localhost -a localhost:3900

# Create bucket
garage bucket create bac-resources

# Add key
garage key create --name bac-admin

# Allow access
garage bucket allow bac-resources --key <key-id> --read --write
```

## Anti-Patterns

- ❌ Using local filesystem for production
- ❌ Storing secrets in bucket
- ❌ Not handling large files (use multipart)
- ❌ Not using content types
