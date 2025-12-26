<template>
  <div class="ai-chat-message" :class="messageClass">
    <div class="message-avatar">
      <div v-if="message.role === 'user'" class="avatar user-avatar">
        <Icon name="mdi:account" size="20" />
      </div>
      <div v-else class="avatar ai-avatar">
        <Icon name="mdi:robot-happy" size="20" />
      </div>
    </div>

    <div class="message-content">
      <div class="message-header">
        <span class="message-sender">
          {{ message.role === 'user' ? 'Anda' : 'AI Tutor' }}
        </span>
        <span class="message-time">
          {{ formatTime(message.createdAt) }}
        </span>
      </div>

      <div class="message-body" v-html="formattedContent"></div>

      <!-- Sources -->
      <div v-if="message.sources && message.sources.length > 0" class="message-sources">
        <div class="sources-header">
          <Icon name="mdi:book-open-page-variant" size="14" />
          <span>Sumber Referensi</span>
        </div>
        <div class="sources-list">
          <div 
            v-for="source in message.sources" 
            :key="source.lessonId"
            class="source-item"
          >
            <span class="source-icon">{{ getContentTypeIcon(source.contentType) }}</span>
            <span class="source-title">{{ source.lessonTitle }}</span>
            <span class="source-relevance">{{ Math.round(source.relevance * 100) }}%</span>
          </div>
        </div>
      </div>

      <!-- Token usage (optional) -->
      <div v-if="message.tokensUsed && showTokens" class="message-tokens">
        <Icon name="mdi:lightning-bolt" size="12" />
        <span>{{ message.tokensUsed }} tokens</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { marked } from 'marked'

interface Source {
  lessonId: string
  lessonTitle: string
  contentType: string
  preview: string
  relevance: number
}

interface ChatMessage {
  id: string
  role: 'user' | 'assistant' | 'system'
  content: string
  sources?: Source[]
  tokensUsed?: number
  createdAt: string
}

const props = defineProps<{
  message: ChatMessage
  showTokens?: boolean
}>()

const messageClass = computed(() => ({
  'user-message': props.message.role === 'user',
  'assistant-message': props.message.role === 'assistant',
  'system-message': props.message.role === 'system'
}))

const formattedContent = computed(() => {
  try {
    // Configure marked for safe rendering
    marked.setOptions({
      breaks: true,
      gfm: true
    })
    return marked(props.message.content)
  } catch {
    return props.message.content
  }
})

function formatTime(dateStr: string): string {
  if (!dateStr) return ''
  try {
    const date = new Date(dateStr)
    return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
  } catch {
    return ''
  }
}

function getContentTypeIcon(contentType: string): string {
  switch (contentType) {
    case 'video':
    case 'video_transcript':
      return '🎥'
    case 'pdf':
      return '📄'
    case 'text':
      return '📝'
    case 'quiz':
      return '📋'
    default:
      return '📚'
  }
}
</script>

<style scoped>
.ai-chat-message {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  margin-bottom: 12px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.user-message {
  background: linear-gradient(135deg, #e8f4fd 0%, #dbeafe 100%);
  margin-left: 24px;
  flex-direction: row-reverse;
}

.assistant-message {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  margin-right: 24px;
  border: 1px solid #e2e8f0;
}

.system-message {
  background: #fef3c7;
  font-style: italic;
}

.message-avatar {
  flex-shrink: 0;
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-avatar {
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
}

.ai-avatar {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.message-content {
  flex: 1;
  min-width: 0;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.message-sender {
  font-weight: 600;
  font-size: 14px;
  color: #374151;
}

.message-time {
  font-size: 12px;
  color: #9ca3af;
}

.message-body {
  font-size: 14px;
  line-height: 1.6;
  color: #374151;
}

.message-body :deep(p) {
  margin-bottom: 8px;
}

.message-body :deep(p:last-child) {
  margin-bottom: 0;
}

.message-body :deep(ul),
.message-body :deep(ol) {
  padding-left: 20px;
  margin-bottom: 8px;
}

.message-body :deep(code) {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
}

.message-body :deep(pre) {
  background: #1f2937;
  color: #f9fafb;
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  margin: 8px 0;
}

.message-body :deep(pre code) {
  background: transparent;
  padding: 0;
}

.message-sources {
  margin-top: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.sources-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 8px;
}

.sources-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.source-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.source-icon {
  font-size: 14px;
}

.source-title {
  flex: 1;
  color: #4b5563;
}

.source-relevance {
  font-size: 11px;
  padding: 2px 6px;
  background: #dcfce7;
  color: #16a34a;
  border-radius: 4px;
  font-weight: 500;
}

.message-tokens {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  font-size: 11px;
  color: #9ca3af;
}
</style>
