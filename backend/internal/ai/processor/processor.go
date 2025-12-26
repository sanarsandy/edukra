package processor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/rag"
)

// ContentProcessor handles extraction of content from various sources
type ContentProcessor struct {
	chunker    *rag.Chunker
	tempDir    string
}

// NewContentProcessor creates a new content processor
func NewContentProcessor() *ContentProcessor {
	tempDir := os.TempDir()
	return &ContentProcessor{
		chunker: rag.NewChunker(rag.DefaultChunkConfig()),
		tempDir: tempDir,
	}
}

// ProcessText processes raw text content
func (p *ContentProcessor) ProcessText(content string) []rag.Chunk {
	return p.chunker.ChunkText(content)
}

// ProcessHTML processes HTML content (strips tags first)
func (p *ContentProcessor) ProcessHTML(html string) []rag.Chunk {
	return p.chunker.ChunkHTML(html)
}

// ProcessPDFFromURL downloads and processes a PDF from URL
func (p *ContentProcessor) ProcessPDFFromURL(ctx context.Context, url string) ([]rag.Chunk, error) {
	// Download PDF to temp file
	tempFile := filepath.Join(p.tempDir, fmt.Sprintf("pdf_%d.pdf", os.Getpid()))
	defer os.Remove(tempFile)

	if err := p.downloadFile(ctx, url, tempFile); err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}

	return p.ProcessPDFFromFile(tempFile)
}

// ProcessPDFFromFile processes a local PDF file
func (p *ContentProcessor) ProcessPDFFromFile(filePath string) ([]rag.Chunk, error) {
	fmt.Printf("[PDF] Starting to extract text from: %s\n", filePath)
	text, err := p.extractTextFromPDF(filePath)
	if err != nil {
		fmt.Printf("[PDF] Extraction error: %v\n", err)
		return nil, err
	}
	fmt.Printf("[PDF] Extracted %d characters of text\n", len(text))

	chunks := p.chunker.ChunkText(text)
	fmt.Printf("[PDF] Created %d chunks\n", len(chunks))
	return chunks, nil
}

// extractTextFromPDF extracts text from PDF using pdftotext (poppler-utils)
func (p *ContentProcessor) extractTextFromPDF(filePath string) (string, error) {
	fmt.Printf("[PDF] Running pdftotext on: %s\n", filePath)
	// Try using pdftotext (poppler-utils)
	cmd := exec.Command("pdftotext", "-layout", "-enc", "UTF-8", filePath, "-")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("[PDF] pdftotext failed: %v, trying Python fallback\n", err)
		// Fallback: try using Python with pdfminer
		return p.extractWithPython(filePath)
	}
	
	text := string(output)
	fmt.Printf("[PDF] pdftotext succeeded, got %d bytes\n", len(output))
	
	// Limit text to prevent memory issues (50KB max)
	maxTextLen := 50000
	if len(text) > maxTextLen {
		fmt.Printf("[PDF] Truncating text from %d to %d chars\n", len(text), maxTextLen)
		text = text[:maxTextLen]
	}
	
	return text, nil
}

// extractWithPython uses Python as fallback for PDF extraction
func (p *ContentProcessor) extractWithPython(filePath string) (string, error) {
	pythonScript := `
import sys
try:
    from pdfminer.high_level import extract_text
    text = extract_text(sys.argv[1])
    print(text)
except ImportError:
    # Try with PyPDF2
    try:
        import PyPDF2
        with open(sys.argv[1], 'rb') as f:
            reader = PyPDF2.PdfReader(f)
            text = ''
            for page in reader.pages:
                text += page.extract_text() + '\n'
            print(text)
    except ImportError:
        print("ERROR: No PDF library available")
        sys.exit(1)
`
	
	cmd := exec.Command("python3", "-c", pythonScript, filePath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to extract PDF text: %w", err)
	}
	
	text := string(output)
	if strings.HasPrefix(text, "ERROR:") {
		return "", fmt.Errorf(text)
	}
	
	return text, nil
}

// downloadFile downloads a file from URL
func (p *ContentProcessor) downloadFile(ctx context.Context, url, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// ProcessYouTubeTranscript fetches YouTube auto-captions
func (p *ContentProcessor) ProcessYouTubeTranscript(ctx context.Context, videoURL string) ([]rag.Chunk, error) {
	// Extract video ID
	videoID := extractYouTubeID(videoURL)
	if videoID == "" {
		return nil, fmt.Errorf("invalid YouTube URL")
	}

	// Try to get transcript using yt-dlp
	text, err := p.getYouTubeTranscriptWithYTDLP(ctx, videoID)
	if err != nil {
		// Fallback to YouTube transcript API
		return nil, fmt.Errorf("could not get YouTube transcript: %w", err)
	}

	return p.chunker.ChunkText(text), nil
}

// getYouTubeTranscriptWithYTDLP uses yt-dlp to get subtitles
func (p *ContentProcessor) getYouTubeTranscriptWithYTDLP(ctx context.Context, videoID string) (string, error) {
	tempFile := filepath.Join(p.tempDir, fmt.Sprintf("yt_%s", videoID))
	defer os.Remove(tempFile + ".id.vtt")
	defer os.Remove(tempFile + ".en.vtt")

	// Try Indonesian first, then English
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--write-auto-sub",
		"--sub-lang", "id,en",
		"--skip-download",
		"-o", tempFile,
		fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID),
	)
	
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read the subtitle file
	var subtitleFile string
	if _, err := os.Stat(tempFile + ".id.vtt"); err == nil {
		subtitleFile = tempFile + ".id.vtt"
	} else if _, err := os.Stat(tempFile + ".en.vtt"); err == nil {
		subtitleFile = tempFile + ".en.vtt"
	} else {
		return "", fmt.Errorf("no subtitles found")
	}

	content, err := os.ReadFile(subtitleFile)
	if err != nil {
		return "", err
	}

	// Parse VTT to plain text
	return parseVTTToText(string(content)), nil
}

// parseVTTToText converts VTT format to plain text
func parseVTTToText(vtt string) string {
	var text strings.Builder
	lines := strings.Split(vtt, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip VTT headers and timestamps
		if line == "" || line == "WEBVTT" || strings.Contains(line, "-->") {
			continue
		}
		// Skip position/alignment tags
		if strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:") {
			continue
		}
		// Skip numeric lines (cue identifiers)
		if _, err := fmt.Sscanf(line, "%d", new(int)); err == nil && len(line) < 10 {
			continue
		}
		// Remove HTML tags
		line = strings.ReplaceAll(line, "<c>", "")
		line = strings.ReplaceAll(line, "</c>", "")
		
		if line != "" {
			text.WriteString(line)
			text.WriteString(" ")
		}
	}

	return text.String()
}

// extractYouTubeID extracts video ID from various YouTube URL formats
func extractYouTubeID(url string) string {
	// youtube.com/watch?v=VIDEO_ID
	if strings.Contains(url, "youtube.com/watch") {
		parts := strings.Split(url, "v=")
		if len(parts) > 1 {
			id := strings.Split(parts[1], "&")[0]
			return id
		}
	}
	// youtu.be/VIDEO_ID
	if strings.Contains(url, "youtu.be/") {
		parts := strings.Split(url, "youtu.be/")
		if len(parts) > 1 {
			id := strings.Split(parts[1], "?")[0]
			return id
		}
	}
	return ""
}
