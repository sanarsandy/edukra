package rag

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ChunkConfig holds configuration for text chunking
type ChunkConfig struct {
	MaxChunkSize    int  // Maximum characters per chunk
	ChunkOverlap    int  // Characters to overlap between chunks
	MinChunkSize    int  // Minimum characters for a chunk
	PreserveWords   bool // Avoid splitting words
	PreserveSentences bool // Try to split at sentence boundaries
}

// DefaultChunkConfig returns sensible default chunking config
func DefaultChunkConfig() ChunkConfig {
	return ChunkConfig{
		MaxChunkSize:      1500, // ~350-400 tokens
		ChunkOverlap:      200,
		MinChunkSize:      100,
		PreserveWords:     true,
		PreserveSentences: true,
	}
}

// Chunk represents a text chunk with metadata
type Chunk struct {
	Index      int    `json:"index"`
	Text       string `json:"text"`
	StartChar  int    `json:"start_char"`
	EndChar    int    `json:"end_char"`
	TokenCount int    `json:"token_count"` // Estimated
}

// Chunker handles text chunking for embeddings
type Chunker struct {
	config ChunkConfig
}

// NewChunker creates a new chunker
func NewChunker(config ChunkConfig) *Chunker {
	return &Chunker{config: config}
}

// ChunkText splits text into chunks with overlap
func (c *Chunker) ChunkText(text string) []Chunk {
	fmt.Printf("[Chunker] Starting with text length: %d\n", len(text))
	
	// Clean and normalize text
	text = c.cleanText(text)
	fmt.Printf("[Chunker] After cleaning, text length: %d\n", len(text))
	
	if len(text) <= c.config.MaxChunkSize {
		fmt.Printf("[Chunker] Text fits in single chunk\n")
		return []Chunk{{
			Index:      0,
			Text:       text,
			StartChar:  0,
			EndChar:    len(text),
			TokenCount: c.estimateTokens(text),
		}}
	}

	var chunks []Chunk
	currentPos := 0
	chunkIndex := 0
	iterCount := 0
	maxIter := 1000 // Safety limit

	for currentPos < len(text) && iterCount < maxIter {
		iterCount++
		
		// Calculate end position
		endPos := currentPos + c.config.MaxChunkSize
		if endPos > len(text) {
			endPos = len(text)
		}

		// Adjust to sentence/word boundary
		if c.config.PreserveSentences && endPos < len(text) {
			endPos = c.findSentenceBoundary(text, currentPos, endPos)
		} else if c.config.PreserveWords && endPos < len(text) {
			endPos = c.findWordBoundary(text, endPos)
		}

		chunkText := text[currentPos:endPos]
		
		// Skip chunks that are too small (except the last one)
		if len(chunkText) < c.config.MinChunkSize && currentPos+len(chunkText) < len(text) {
			currentPos = endPos
			continue
		}

		chunks = append(chunks, Chunk{
			Index:      chunkIndex,
			Text:       strings.TrimSpace(chunkText),
			StartChar:  currentPos,
			EndChar:    endPos,
			TokenCount: c.estimateTokens(chunkText),
		})
		chunkIndex++

		// Move position with overlap
		currentPos = endPos - c.config.ChunkOverlap
		if currentPos < 0 || currentPos >= endPos {
			currentPos = endPos
		}
	}
	
	if iterCount >= maxIter {
		fmt.Printf("[Chunker] WARNING: Hit max iterations!\n")
	}

	fmt.Printf("[Chunker] Created %d chunks\n", len(chunks))
	return chunks
}

// cleanText normalizes text for processing
func (c *Chunker) cleanText(text string) string {
	// Remove invalid UTF-8 byte sequences
	text = strings.ToValidUTF8(text, "")
	
	// Remove multiple whitespaces
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	
	// Remove zero-width characters and control characters
	text = regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`).ReplaceAllString(text, "")
	
	// Trim
	text = strings.TrimSpace(text)
	
	return text
}

// findSentenceBoundary finds the nearest sentence end before maxPos
func (c *Chunker) findSentenceBoundary(text string, minPos, maxPos int) int {
	if maxPos >= len(text) {
		return len(text)
	}

	// Look backwards for sentence endings
	sentenceEnders := []string{". ", ".\n", "! ", "!\n", "? ", "?\n", "。", "！", "？"}
	
	bestPos := maxPos
	for _, ender := range sentenceEnders {
		idx := strings.LastIndex(text[minPos:maxPos], ender)
		if idx > 0 {
			pos := minPos + idx + len(ender)
			if pos > minPos+c.config.MinChunkSize && pos < bestPos {
				bestPos = pos
			}
		}
	}

	// If no sentence boundary found, fall back to word boundary
	if bestPos == maxPos {
		return c.findWordBoundary(text, maxPos)
	}

	return bestPos
}

// findWordBoundary finds the nearest word end before pos
func (c *Chunker) findWordBoundary(text string, pos int) int {
	if pos >= len(text) {
		return len(text)
	}

	// Look backwards for whitespace
	for i := pos; i > pos-50 && i > 0; i-- {
		if unicode.IsSpace(rune(text[i])) {
			return i
		}
	}

	return pos
}

// estimateTokens gives rough token estimate (1 token ≈ 4 chars for English)
func (c *Chunker) estimateTokens(text string) int {
	// This is a rough estimate
	// More accurate would use tiktoken or similar
	charCount := len(text)
	
	// Average chars per token varies by language
	// English: ~4 chars/token
	// Indonesian: ~3-4 chars/token
	return (charCount + 3) / 4
}

// ChunkHTML strips HTML and chunks the text
func (c *Chunker) ChunkHTML(html string) []Chunk {
	// Strip HTML tags
	text := c.stripHTML(html)
	return c.ChunkText(text)
}

// stripHTML removes HTML tags from text
func (c *Chunker) stripHTML(html string) string {
	// Remove script and style tags entirely
	scriptRe := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	html = scriptRe.ReplaceAllString(html, "")
	
	styleRe := regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`)
	html = styleRe.ReplaceAllString(html, "")
	
	// Convert block elements to newlines
	blockRe := regexp.MustCompile(`(?i)</?(p|div|br|h[1-6]|li|tr)[^>]*>`)
	html = blockRe.ReplaceAllString(html, "\n")
	
	// Remove all other HTML tags
	tagRe := regexp.MustCompile(`<[^>]+>`)
	text := tagRe.ReplaceAllString(html, "")
	
	// Decode common HTML entities
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	
	return text
}
