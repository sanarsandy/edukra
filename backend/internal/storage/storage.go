package storage

import (
	"context"
	"io"
	"time"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Upload uploads a file to the specified bucket
	Upload(ctx context.Context, bucket, objectName string, reader io.Reader, size int64, contentType string) error

	// GetPresignedURL generates a pre-signed URL for downloading an object
	GetPresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error)

	// Delete removes an object from the storage
	Delete(ctx context.Context, bucket, objectName string) error

	// Exists checks if an object exists in the storage
	Exists(ctx context.Context, bucket, objectName string) (bool, error)
}

// ContentType constants for common file types
const (
	ContentTypeMP4  = "video/mp4"
	ContentTypeWebM = "video/webm"
	ContentTypePDF  = "application/pdf"
	ContentTypeJPEG = "image/jpeg"
	ContentTypePNG  = "image/png"
)

// BucketName for course content
const (
	BucketCourseContent = "course-content"
)
