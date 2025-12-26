<template>
  <div class="ai-chat-input">
    <div v-if="remainingQuota <= 5 && remainingQuota > 0" class="quota-warning">
      <Icon name="mdi:alert" size="14" />
      <span>Sisa {{ remainingQuota }} pertanyaan hari ini</span>
    </div>

    <div v-if="remainingQuota <= 0" class="quota-exhausted">
      <Icon name="mdi:clock-alert" size="14" />
      <span>Batas harian tercapai. Coba lagi besok.</span>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="input-form">
      <div class="input-wrapper">
        <textarea
          ref="inputRef"
          v-model="inputText"
          placeholder="Tanyakan sesuatu tentang materi..."
          :disabled="disabled || isSending"
          @keydown="handleKeydown"
          rows="1"
        ></textarea>

        <button
          type="submit"
          :disabled="!canSubmit"
          class="send-button"
        >
          <Icon v-if="isSending" name="mdi:loading" size="20" class="spinning" />
          <Icon v-else name="mdi:send" size="20" />
        </button>
      </div>

      <div class="input-footer">
        <span class="hint">
          <kbd>Enter</kbd> untuk kirim, <kbd>Shift+Enter</kbd> untuk baris baru
        </span>
        <span class="quota-info">
          {{ remainingQuota }}/{{ quota.limit }} tersisa
        </span>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch, onMounted } from 'vue'

interface QuotaInfo {
  used: number
  limit: number
  remaining: number
}

const props = defineProps<{
  isSending: boolean
  disabled?: boolean
  quota: QuotaInfo
}>()

const emit = defineEmits<{
  send: [message: string]
}>()

const inputText = ref('')
const inputRef = ref<HTMLTextAreaElement | null>(null)

const remainingQuota = computed(() => props.quota.remaining)

const canSubmit = computed(() => {
  return inputText.value.trim().length > 0 && 
         !props.isSending && 
         !props.disabled &&
         remainingQuota.value > 0
})

function handleSubmit() {
  if (!canSubmit.value) return

  const message = inputText.value.trim()
  inputText.value = ''
  emit('send', message)
  
  // Reset textarea height
  if (inputRef.value) {
    inputRef.value.style.height = 'auto'
  }
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSubmit()
  }
}

// Auto-resize textarea
watch(inputText, () => {
  nextTick(() => {
    if (inputRef.value) {
      inputRef.value.style.height = 'auto'
      inputRef.value.style.height = Math.min(inputRef.value.scrollHeight, 150) + 'px'
    }
  })
})

onMounted(() => {
  inputRef.value?.focus()
})

defineExpose({
  focus: () => inputRef.value?.focus()
})
</script>

<style scoped>
.ai-chat-input {
  padding: 16px;
  background: white;
  border-top: 1px solid #e5e7eb;
}

.quota-warning {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: #fef3c7;
  border: 1px solid #fcd34d;
  border-radius: 8px;
  font-size: 13px;
  color: #92400e;
}

.quota-exhausted {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 12px;
  background: #fee2e2;
  border: 1px solid #fca5a5;
  border-radius: 8px;
  font-size: 13px;
  color: #991b1b;
}

.input-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-wrapper {
  display: flex;
  gap: 8px;
  align-items: flex-end;
}

textarea {
  flex: 1;
  resize: none;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px 16px;
  font-size: 14px;
  font-family: inherit;
  line-height: 1.5;
  min-height: 48px;
  max-height: 150px;
  transition: border-color 0.2s, box-shadow 0.2s;
}

textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

textarea:disabled {
  background: #f9fafb;
  cursor: not-allowed;
}

textarea::placeholder {
  color: #9ca3af;
}

.send-button {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, opacity 0.2s;
}

.send-button:hover:not(:disabled) {
  transform: scale(1.05);
}

.send-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
  color: #9ca3af;
}

.hint kbd {
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: inherit;
}

.quota-info {
  color: #6b7280;
}
</style>
