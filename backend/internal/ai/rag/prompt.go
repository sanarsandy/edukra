package rag

import (
	"fmt"
	"strings"

	"github.com/lman-kadiv-doti/secure-whitelabel-lms/backend/internal/ai/providers"
)

// PromptTemplate defines the structure for AI prompts
type PromptTemplate struct {
	SystemPrompt string
	ContextPrompt string
	QuestionPrompt string
}

// DefaultPromptTemplate returns the default Indonesian prompt template
func DefaultPromptTemplate() PromptTemplate {
	return PromptTemplate{
		SystemPrompt: `Kamu adalah AI Tutor Edukra yang membantu siswa memahami materi kursus.

## PERAN UTAMA
Menjawab pertanyaan siswa HANYA berdasarkan materi/konteks yang diberikan.

## ATURAN WAJIB (TIDAK BOLEH DILANGGAR)
1. **JAWAB HANYA DARI KONTEKS**: Gunakan HANYA informasi yang ada dalam materi referensi yang diberikan.
2. **JANGAN MENGARANG**: Jika informasi tidak ada dalam konteks, WAJIB katakan: "Maaf, informasi tersebut tidak tersedia dalam materi kursus ini."
3. **JANGAN GUNAKAN PENGETAHUAN EKSTERNAL**: Meskipun kamu tahu jawabannya dari sumber lain, JANGAN gunakan jika tidak ada di konteks.
4. **AKUI KETERBATASAN**: Lebih baik mengatakan "tidak tahu" daripada mengarang jawaban yang salah.

## CARA MENJAWAB
- Gunakan bahasa Indonesia yang jelas dan mudah dipahami
- Berikan jawaban ringkas namun lengkap
- Gunakan bullet points atau numbering jika membantu
- Sebutkan sumber materi jika relevan (contoh: "Berdasarkan materi [nama lesson]...")

## CONTOH RESPON KETIKA TIDAK TAHU
❌ SALAH: "Menurut saya..." atau "Biasanya..." atau mengarang jawaban
✅ BENAR: "Maaf, saya tidak menemukan informasi tentang [topik] dalam materi kursus ini. Silakan tanyakan kepada instruktur atau cek materi lain yang mungkin membahas topik tersebut."`,

		ContextPrompt: `Berikut adalah materi referensi dari kursus yang dapat kamu gunakan untuk menjawab pertanyaan:

%s

---
INGAT: Jawab HANYA berdasarkan materi di atas. Jangan mengarang!`,

		QuestionPrompt: `Pertanyaan dari siswa: %s

Berikan jawaban yang membantu berdasarkan materi di atas. Jika tidak ada informasinya, katakan dengan jujur.`,
	}
}

// PromptBuilder builds prompts for the AI
type PromptBuilder struct {
	template PromptTemplate
}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder(template PromptTemplate) *PromptBuilder {
	return &PromptBuilder{template: template}
}

// BuildMessages creates the message array for the AI
func (b *PromptBuilder) BuildMessages(context, question string, chatHistory []providers.Message) []providers.Message {
	messages := []providers.Message{
		{
			Role:    "system",
			Content: b.template.SystemPrompt,
		},
	}

	// Add context if available
	if context != "" {
		messages = append(messages, providers.Message{
			Role:    "user",
			Content: fmt.Sprintf(b.template.ContextPrompt, context),
		})
		messages = append(messages, providers.Message{
			Role:    "assistant",
			Content: "Baik, saya sudah membaca materi referensi tersebut. Silakan ajukan pertanyaan Anda.",
		})
	}

	// Add chat history (limited to last N messages)
	maxHistoryMessages := 10
	startIdx := 0
	if len(chatHistory) > maxHistoryMessages {
		startIdx = len(chatHistory) - maxHistoryMessages
	}
	
	for _, msg := range chatHistory[startIdx:] {
		// Skip system messages from history
		if msg.Role != "system" {
			messages = append(messages, msg)
		}
	}

	// Add current question
	messages = append(messages, providers.Message{
		Role:    "user",
		Content: fmt.Sprintf(b.template.QuestionPrompt, question),
	})

	return messages
}

// BuildSimpleMessages creates messages without RAG context (for simple questions)
func (b *PromptBuilder) BuildSimpleMessages(question string, chatHistory []providers.Message) []providers.Message {
	messages := []providers.Message{
		{
			Role:    "system",
			Content: b.template.SystemPrompt,
		},
	}

	// Add chat history
	for _, msg := range chatHistory {
		if msg.Role != "system" {
			messages = append(messages, msg)
		}
	}

	// Add current question
	messages = append(messages, providers.Message{
		Role:    "user",
		Content: question,
	})

	return messages
}

// FormatResponse formats the AI response with source citations
func FormatResponse(response string, sources []Source) string {
	if len(sources) == 0 {
		return response
	}

	var sb strings.Builder
	sb.WriteString(response)
	sb.WriteString("\n\n---\n📚 **Sumber:**\n")

	for _, source := range sources {
		icon := getContentTypeIcon(source.ContentType)
		relevancePercent := int(source.Relevance * 100)
		sb.WriteString(fmt.Sprintf("- %s %s (%d%% relevan)\n", icon, source.LessonTitle, relevancePercent))
	}

	return sb.String()
}

// getContentTypeIcon returns emoji for content type
func getContentTypeIcon(contentType string) string {
	switch contentType {
	case "video", "video_transcript":
		return "🎥"
	case "pdf":
		return "📄"
	case "text":
		return "📝"
	case "quiz":
		return "📋"
	default:
		return "📚"
	}
}

// GreetingPrompt returns the initial greeting message
func GreetingPrompt(courseName string) string {
	return fmt.Sprintf(`Halo! 👋 Saya AI Tutor untuk kursus **"%s"**.

Saya sudah mempelajari semua materi di kursus ini dan siap membantu Anda memahami kontennya. Silakan tanyakan apa saja tentang materi kursus!

Beberapa contoh pertanyaan yang bisa Anda ajukan:
- Apa itu [konsep dari materi]?
- Jelaskan perbedaan antara [A] dan [B]
- Bagaimana cara [melakukan sesuatu]?
- Berikan contoh [topik tertentu]

Saya akan menjawab berdasarkan materi yang tersedia di kursus ini. 🎓`, courseName)
}
