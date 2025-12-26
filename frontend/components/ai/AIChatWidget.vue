<template>
  <!-- Floating Button -->
  <div class="ai-tutor-widget">
    <Transition name="slide-up">
      <div v-if="isOpen" class="chat-container">
        <!-- Header -->
        <div class="chat-header">
          <div class="header-info">
            <div class="ai-icon">
              <Icon name="mdi:robot-happy" size="24" />
            </div>
            <div class="header-text">
              <h3>AI Tutor</h3>
              <span class="status">
                <span v-if="hasEmbeddings" class="status-dot online"></span>
                <span v-else class="status-dot offline"></span>
                {{ hasEmbeddings ? 'Siap membantu' : 'Materi belum diproses' }}
              </span>
            </div>
          </div>
          <div class="header-actions">
            <button @click="handleClear" class="action-btn" title="Bersihkan chat">
              <Icon name="mdi:broom" size="18" />
            </button>
            <button @click="closeChat" class="action-btn close-btn" title="Tutup">
              <Icon name="mdi:close" size="18" />
            </button>
          </div>
        </div>

        <!-- Messages -->
        <div ref="messagesContainer" class="chat-messages">
          <div v-if="isLoading" class="loading-state">
            <Icon name="mdi:loading" size="24" class="spinning" />
            <span>Memuat...</span>
          </div>

          <template v-else>
            <AIChatMessage
              v-for="msg in messages"
              :key="msg.id"
              :message="msg"
            />
          </template>

          <div v-if="isSending" class="typing-indicator">
            <div class="typing-dot"></div>
            <div class="typing-dot"></div>
            <div class="typing-dot"></div>
          </div>
        </div>

        <!-- Error -->
        <Transition name="fade">
          <div v-if="error" class="chat-error">
            <Icon name="mdi:alert-circle" size="16" />
            <span>{{ error }}</span>
            <button @click="error = null">
              <Icon name="mdi:close" size="14" />
            </button>
          </div>
        </Transition>

        <!-- Input -->
        <AIChatInput
          :is-sending="isSending"
          :disabled="!isAIEnabled"
          :quota="quota"
          @send="handleSend"
        />
      </div>
    </Transition>

    <!-- Floating Button -->
    <button 
      v-if="isAIEnabled"
      @click="toggleChat" 
      class="floating-button"
      :class="{ 'is-open': isOpen }"
    >
      <Icon v-if="isOpen" name="mdi:close" size="28" />
      <Icon v-else name="mdi:robot-happy-outline" size="28" />
      
      <!-- Pulse animation when not open -->
      <span v-if="!isOpen && hasEmbeddings" class="pulse-ring"></span>
    </button>

    <!-- Disabled notice -->
    <div v-if="!isAIEnabled && showDisabledNotice" class="disabled-notice">
      <Icon name="mdi:robot-confused" size="20" />
      <span>AI Tutor belum aktif</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import { useAITutor } from '~/composables/useAITutor'

const props = defineProps<{
  courseId: string
  showDisabledNotice?: boolean
}>()

const {
  messages,
  session,
  quota,
  aiStatus,
  isLoading,
  isSending,
  error,
  isOpen,
  isAIEnabled,
  hasEmbeddings,
  remainingQuota,
  canChat,
  loadSession,
  sendMessage,
  clearSession,
  toggleChat,
  openChat,
  closeChat
} = useAITutor(props.courseId)

const messagesContainer = ref<HTMLElement | null>(null)

// Auto-scroll to bottom when new messages
watch(messages, () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}, { deep: true })

// Scroll to bottom when opening
watch(isOpen, (open) => {
  if (open) {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
  }
})

async function handleSend(message: string) {
  try {
    await sendMessage(message)
  } catch (err) {
    // Error handled by composable
  }
}

async function handleClear() {
  if (confirm('Yakin ingin menghapus semua pesan chat?')) {
    await clearSession()
  }
}

// Load session when first opened
onMounted(() => {
  if (isOpen.value) {
    loadSession()
  }
})
</script>

<style scoped>
.ai-tutor-widget {
  position: fixed;
  bottom: 24px;
  right: 24px;
  z-index: 1000;
}

/* Chat Container */
.chat-container {
  position: absolute;
  bottom: 80px;
  right: 0;
  width: 420px;
  max-height: calc(100vh - 150px);
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

@media (max-width: 480px) {
  .chat-container {
    width: calc(100vw - 32px);
    right: -8px;
    max-height: calc(100vh - 120px);
  }
}

/* Header */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.header-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ai-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-text h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  opacity: 0.9;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-dot.online {
  background: #4ade80;
  animation: pulse 2s infinite;
}

.status-dot.offline {
  background: #fbbf24;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.header-actions {
  display: flex;
  gap: 4px;
}

.action-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* Messages */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  min-height: 300px;
  max-height: 400px;
  background: #f8fafc;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 200px;
  color: #6b7280;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.typing-indicator {
  display: flex;
  gap: 4px;
  padding: 16px;
  margin-left: 48px;
}

.typing-dot {
  width: 8px;
  height: 8px;
  background: #9ca3af;
  border-radius: 50%;
  animation: bounce 1.4s infinite ease-in-out;
}

.typing-dot:nth-child(1) { animation-delay: 0s; }
.typing-dot:nth-child(2) { animation-delay: 0.2s; }
.typing-dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-6px); }
}

/* Error */
.chat-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: #fee2e2;
  border-top: 1px solid #fca5a5;
  font-size: 13px;
  color: #991b1b;
}

.chat-error button {
  margin-left: auto;
  background: none;
  border: none;
  color: #991b1b;
  cursor: pointer;
  padding: 4px;
}

/* Floating Button */
.floating-button {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: none;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(16, 185, 129, 0.4);
  transition: transform 0.3s, box-shadow 0.3s;
  position: relative;
}

.floating-button:hover {
  transform: scale(1.1);
  box-shadow: 0 12px 32px rgba(16, 185, 129, 0.5);
}

.floating-button.is-open {
  background: linear-gradient(135deg, #64748b 0%, #475569 100%);
  box-shadow: 0 8px 24px rgba(71, 85, 105, 0.4);
}

.pulse-ring {
  position: absolute;
  inset: -4px;
  border-radius: 50%;
  border: 2px solid #10b981;
  animation: pulse-ring 2s infinite;
}

@keyframes pulse-ring {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(1.3);
    opacity: 0;
  }
}

/* Disabled notice */
.disabled-notice {
  position: absolute;
  bottom: 80px;
  right: 0;
  background: #f3f4f6;
  padding: 12px 16px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  white-space: nowrap;
}

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
