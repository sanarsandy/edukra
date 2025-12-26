package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorage implements Storage interface using MinIO
type MinioStorage struct {
	client           *minio.Client
	publicClient     *minio.Client // Client with public endpoint for URL signing
	defaultBucket    string
	presignExpiry    time.Duration
	internalEndpoint string
	publicEndpoint   string
}

// MinioConfig holds configuration for MinIO connection
type MinioConfig struct {
	Endpoint        string
	PublicEndpoint  string // External/public endpoint for browser access
	AccessKey       string
	SecretKey       string
	UseSSL          bool
	BucketContent   string
	PresignedExpiry int // in seconds
}

// LoadMinioConfigFromEnv loads MinIO configuration from environment variables
func LoadMinioConfigFromEnv() *MinioConfig {
	useSSL := false
	if os.Getenv("MINIO_USE_SSL") == "true" {
		useSSL = true
	}

	expiry := 7200 // default 2 hours
	if exp := os.Getenv("MINIO_PRESIGNED_EXPIRY"); exp != "" {
		if parsed, err := strconv.Atoi(exp); err == nil {
			expiry = parsed
		}
	}

	return &MinioConfig{
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
		PublicEndpoint:  getEnvOrDefault("MINIO_PUBLIC_ENDPOINT", os.Getenv("MINIO_ENDPOINT")),
		AccessKey:       os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey:       os.Getenv("MINIO_SECRET_KEY"),
		UseSSL:          useSSL,
		BucketContent:   getEnvOrDefault("MINIO_BUCKET_CONTENT", BucketCourseContent),
		PresignedExpiry: expiry,
	}
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// NewMinioStorage creates a new MinIO storage client
func NewMinioStorage(cfg *MinioConfig) (*MinioStorage, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT is required")
	}

	// Main client for internal operations (connects via Docker network)
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Public client for URL signing (uses public endpoint for correct signature)
	var publicClient *minio.Client
	publicEndpoint := cfg.PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.Endpoint
	}
	
	publicClient, err = minio.New(publicEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Printf("Warning: Failed to create public MinIO client, using internal: %v", err)
		publicClient = client
	}

	storage := &MinioStorage{
		client:           client,
		publicClient:     publicClient,
		defaultBucket:    cfg.BucketContent,
		presignExpiry:    time.Duration(cfg.PresignedExpiry) * time.Second,
		internalEndpoint: cfg.Endpoint,
		publicEndpoint:   publicEndpoint,
	}

	// Ensure bucket exists
	ctx := context.Background()
	if err := storage.ensureBucket(ctx, cfg.BucketContent); err != nil {
		log.Printf("Warning: Failed to ensure bucket %s: %v", cfg.BucketContent, err)
	}

	return storage, nil
}

// ensureBucket creates a bucket if it doesn't exist
func (m *MinioStorage) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := m.client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		log.Printf("Created MinIO bucket: %s", bucket)
	}
	return nil
}

// Upload uploads a file to MinIO
func (m *MinioStorage) Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	opts := minio.PutObjectOptions{
		ContentType: contentType,
	}

	_, err := m.client.PutObject(ctx, bucket, objectName, reader, size, opts)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	log.Printf("Uploaded file to MinIO: %s/%s", bucket, objectName)
	return nil
}

// GetPresignedURL generates a pre-signed URL for downloading an object
func (m *MinioStorage) GetPresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	if expiry == 0 {
		expiry = m.presignExpiry
	}

	// Set request parameters for content-disposition
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=\"%s\"", objectName))

	// Use internal client to generate URL (publicClient can't connect from inside Docker)
	presignedURL, err := m.client.PresignedGetObject(ctx, bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Replace internal endpoint with public endpoint for browser access
	// Note: MinIO signature does NOT include Host header, so this is safe
	urlStr := presignedURL.String()
	if m.publicEndpoint != "" && m.publicEndpoint != m.internalEndpoint {
		urlStr = strings.Replace(urlStr, m.internalEndpoint, m.publicEndpoint, 1)
	}

	return urlStr, nil
}

// GetPresignedURLWithDefaultExpiry generates a pre-signed URL with default expiry
func (m *MinioStorage) GetPresignedURLWithDefaultExpiry(ctx context.Context, bucket, objectName string) (string, time.Time, error) {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	expiresAt := time.Now().Add(m.presignExpiry)

	presignedURL, err := m.GetPresignedURL(ctx, bucket, objectName, m.presignExpiry)
	if err != nil {
		return "", time.Time{}, err
	}

	return presignedURL, expiresAt, nil
}

// Delete removes an object from MinIO
func (m *MinioStorage) Delete(ctx context.Context, bucket, objectName string) error {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	err := m.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	log.Printf("Deleted file from MinIO: %s/%s", bucket, objectName)
	return nil
}

// GetObject retrieves an object from MinIO and returns it as an io.ReadCloser
func (m *MinioStorage) GetObject(ctx context.Context, bucket, objectName string) (*minio.Object, error) {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	obj, err := m.client.GetObject(ctx, bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return obj, nil
}

// GetObjectInfo retrieves object info including content type
func (m *MinioStorage) GetObjectInfo(ctx context.Context, bucket, objectName string) (minio.ObjectInfo, error) {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	info, err := m.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, fmt.Errorf("failed to get object info: %w", err)
	}

	return info, nil
}

// Exists checks if an object exists in MinIO
func (m *MinioStorage) Exists(ctx context.Context, bucket, objectName string) (bool, error) {
	if bucket == "" {
		bucket = m.defaultBucket
	}

	_, err := m.client.StatObject(ctx, bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetDefaultBucket returns the default bucket name
func (m *MinioStorage) GetDefaultBucket() string {
	return m.defaultBucket
}

// GetDefaultExpiry returns the default presigned URL expiry duration
func (m *MinioStorage) GetDefaultExpiry() time.Duration {
	return m.presignExpiry
}

// Singleton instance for global access
var defaultStorage *MinioStorage

// InitStorage initializes the global storage instance
func InitStorage() error {
	cfg := LoadMinioConfigFromEnv()
	
	// If MinIO is not configured, skip initialization
	if cfg.Endpoint == "" {
		log.Println("MinIO not configured, skipping storage initialization")
		return nil
	}

	storage, err := NewMinioStorage(cfg)
	if err != nil {
		return err
	}

	defaultStorage = storage
	log.Printf("MinIO storage initialized: endpoint=%s, bucket=%s", cfg.Endpoint, cfg.BucketContent)
	return nil
}

// GetStorage returns the global storage instance
func GetStorage() *MinioStorage {
	return defaultStorage
}

// IsConfigured returns true if storage is configured and initialized
func IsConfigured() bool {
	return defaultStorage != nil
}
