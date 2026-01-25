package video

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// HLSConfig contains configuration for HLS encoding
type HLSConfig struct {
	SegmentDuration int    // Segment duration in seconds
	KeyInfoPath     string // Path to key info file
	OutputDir       string // Output directory for HLS files
}

// DefaultHLSConfig returns default HLS configuration
func DefaultHLSConfig() HLSConfig {
	return HLSConfig{
		SegmentDuration: 10, // 10 second segments
	}
}

// EncryptionKey holds the encryption key and IV for a video
type EncryptionKey struct {
	Key []byte // 16-byte AES-128 key
	IV  []byte // 16-byte initialization vector
}

// GenerateEncryptionKey generates a new AES-128 encryption key and IV
func GenerateEncryptionKey() (*EncryptionKey, error) {
	key := make([]byte, 16) // AES-128 = 16 bytes
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}

	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	return &EncryptionKey{Key: key, IV: iv}, nil
}

// HLSResult contains the result of HLS processing
type HLSResult struct {
	ManifestPath string   // Path to master .m3u8 file
	SegmentPaths []string // Paths to .ts segment files
	KeyPath      string   // Path to key file
	Duration     float64  // Total duration in seconds
}

// Processor handles video processing operations
type Processor struct {
	tempDir string
}

// NewProcessor creates a new video processor
func NewProcessor() (*Processor, error) {
	tempDir, err := os.MkdirTemp("", "hls-processing-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	return &Processor{tempDir: tempDir}, nil
}

// Cleanup removes temporary files
func (p *Processor) Cleanup() {
	if p.tempDir != "" {
		os.RemoveAll(p.tempDir)
	}
}

// ConvertToHLS converts a video file to encrypted HLS format
func (p *Processor) ConvertToHLS(ctx context.Context, inputPath string, encKey *EncryptionKey, outputPrefix string) (*HLSResult, error) {
	log.Printf("Starting HLS conversion for: %s", inputPath)

	// Create output directory
	outputDir := filepath.Join(p.tempDir, outputPrefix)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output dir: %w", err)
	}

	// Write key file
	keyPath := filepath.Join(outputDir, "key.bin")
	if err := os.WriteFile(keyPath, encKey.Key, 0600); err != nil {
		return nil, fmt.Errorf("failed to write key file: %w", err)
	}

	// Create key info file for FFmpeg
	// Format: key_uri\nkey_path\nIV
	ivHex := hex.EncodeToString(encKey.IV)
	keyInfoPath := filepath.Join(outputDir, "keyinfo.txt")
	// Use relative URI - will be replaced during serving
	keyInfoContent := fmt.Sprintf("key.bin\n%s\n%s", keyPath, ivHex)
	if err := os.WriteFile(keyInfoPath, []byte(keyInfoContent), 0600); err != nil {
		return nil, fmt.Errorf("failed to write keyinfo file: %w", err)
	}

	// Manifest and segment paths
	manifestPath := filepath.Join(outputDir, "playlist.m3u8")
	segmentPattern := filepath.Join(outputDir, "segment_%03d.ts")

	// Build FFmpeg command
	args := []string{
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-preset", "fast",
		"-crf", "23",
		"-sc_threshold", "0",
		"-g", "48",
		"-keyint_min", "48",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", segmentPattern,
		"-hls_key_info_file", keyInfoPath,
		"-hls_playlist_type", "vod",
		"-f", "hls",
		manifestPath,
	}

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Capture stderr for debugging
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	log.Printf("Running FFmpeg: ffmpeg %s", strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		log.Printf("FFmpeg error output: %s", stderr.String())
		return nil, fmt.Errorf("FFmpeg failed: %w, stderr: %s", err, stderr.String())
	}

	// Collect segment files
	segments, err := filepath.Glob(filepath.Join(outputDir, "segment_*.ts"))
	if err != nil {
		return nil, fmt.Errorf("failed to list segments: %w", err)
	}

	log.Printf("HLS conversion complete: %d segments created", len(segments))

	return &HLSResult{
		ManifestPath: manifestPath,
		SegmentPaths: segments,
		KeyPath:      keyPath,
	}, nil
}

// GetOutputDir returns the output directory path
func (p *Processor) GetOutputDir() string {
	return p.tempDir
}

// RewriteManifestKeyURI rewrites the key URI in the manifest to use a custom URL
func RewriteManifestKeyURI(manifestPath, keyURI string) error {
	content, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	// Replace key.bin with the actual key URI
	newContent := strings.Replace(string(content), 
		`URI="key.bin"`, 
		fmt.Sprintf(`URI="%s"`, keyURI), 
		-1)

	if err := os.WriteFile(manifestPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// EncryptData encrypts data using AES-128-CBC (for testing/utility)
func EncryptData(key, iv, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Pad to block size
	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintext = append(plaintext, padtext...)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// ValidateFFmpeg checks if FFmpeg is installed and accessible
func ValidateFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg not found or not executable: %w", err)
	}
	return nil
}

// GetVideoDuration returns the duration of a video file in seconds
func GetVideoDuration(inputPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get video duration: %w", err)
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}
