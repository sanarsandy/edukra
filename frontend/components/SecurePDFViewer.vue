<template>
  <Transition name="fade">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-50 bg-black/90 flex flex-col"
      @contextmenu.prevent
    >
      <!-- Header -->
      <div class="flex items-center justify-between p-4 bg-neutral-900 border-b border-neutral-800">
        <div class="flex items-center gap-3">
          <svg class="w-6 h-6 text-red-500" viewBox="0 0 24 24" fill="currentColor">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 2l5 5h-5V4z"/>
          </svg>
          <span class="text-white font-medium truncate max-w-md">{{ title }}</span>
        </div>
        <div class="flex items-center gap-2">
          <!-- Close Button -->
          <button 
            @click="close" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- PDF Viewer Container -->
      <div class="flex-1 relative overflow-hidden protected-content">
        <!-- Loading State -->
        <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-black z-10">
          <div class="text-center text-white">
            <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
            <p class="text-sm opacity-80">Memuat dokumen...</p>
          </div>
        </div>

        <!-- Error State -->
        <div v-if="error" class="absolute inset-0 flex items-center justify-center bg-black z-10">
          <div class="text-center text-white p-4">
            <svg class="w-12 h-12 mx-auto mb-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <p class="text-sm mb-2">{{ error }}</p>
            <p class="text-xs text-neutral-400 mb-4">Dokumen ini mungkin tidak bisa ditampilkan secara inline.</p>
            <a :href="pdfUrl" target="_blank" class="px-4 py-2 bg-neutral-800 rounded-lg text-sm hover:bg-neutral-700 transition-colors">
              Buka di Tab Baru
            </a>
          </div>
        </div>

        <!-- PDF via iframe -->
        <iframe
          v-if="viewerUrl"
          :src="viewerUrl"
          class="w-full h-full border-0"
          @load="onIframeLoad"
          @error="onIframeError"
        ></iframe>

        <!-- Watermark Overlay (always on top) -->
        <div 
          v-if="userEmail"
          class="absolute inset-0 pointer-events-none select-none overflow-hidden z-20"
        >
          <div class="watermark-pattern" :style="{ opacity: 0.08 }">
            <template v-for="row in 6" :key="row">
              <div class="watermark-row" :style="{ top: `${row * 16 - 8}%` }">
                <template v-for="col in 4" :key="col">
                  <span 
                    class="watermark-text"
                    :style="{ left: `${col * 25 - 12}%` }"
                  >
                    {{ maskedEmail }}
                  </span>
                </template>
              </div>
            </template>
          </div>
        </div>

        <!-- Interaction Blocker Overlay -->
        <div 
          class="absolute inset-0 z-10"
          style="pointer-events: none;"
          @contextmenu.prevent
        ></div>
      </div>

      <!-- Footer -->
      <div class="p-3 bg-neutral-900 border-t border-neutral-800 flex justify-between items-center">
        <p class="text-xs text-neutral-500">
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          Dokumen ini dilindungi
        </p>
        <button 
          v-if="!isCompleted"
          @click="markComplete"
          class="px-4 py-2 bg-accent-600 hover:bg-accent-700 text-white text-sm font-medium rounded-lg transition-colors"
        >
          Tandai Selesai
        </button>
        <span v-else class="text-accent-400 text-sm flex items-center gap-1">
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z"/>
          </svg>
          Selesai
        </span>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'

interface Props {
  isOpen: boolean
  pdfUrl: string
  title?: string
  userEmail?: string
  isCompleted?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Dokumen',
  isCompleted: false
})

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'complete'): void
}>()

const loading = ref(true)
const error = ref<string | null>(null)

const maskedEmail = computed(() => {
  if (!props.userEmail) return ''
  const [local, domain] = props.userEmail.split('@')
  if (local && domain) {
    const maskedLocal = local.length > 3 
      ? `${local.slice(0, 2)}***${local.slice(-1)}`
      : local
    return `${maskedLocal}@${domain}`
  }
  return props.userEmail
})

// Determine the viewer URL
const viewerUrl = computed(() => {
  if (!props.pdfUrl) return ''
  
  // For external URLs - use Google Docs Viewer to bypass CORS
  if (!props.pdfUrl.startsWith('/') && !props.pdfUrl.startsWith(window.location.origin) && !props.pdfUrl.startsWith('blob:')) {
    return `https://docs.google.com/viewer?url=${encodeURIComponent(props.pdfUrl)}&embedded=true`
  }
  
  return props.pdfUrl
})

const onIframeLoad = () => {
  loading.value = false
  error.value = null
}

const onIframeError = () => {
  loading.value = false
  error.value = 'Gagal memuat dokumen'
}

const close = () => {
  emit('close')
}

const markComplete = () => {
  emit('complete')
}

// Handle keyboard shortcuts
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
    e.preventDefault()
    return
  }
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    return
  }
  if (e.key === 'Escape') {
    close()
  }
}

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    loading.value = true
    error.value = null
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
    loading.value = true
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { 
  transition: opacity 0.2s ease; 
}
.fade-enter-from, .fade-leave-to { 
  opacity: 0; 
}

.protected-content {
  -webkit-user-select: none !important;
  -moz-user-select: none !important;
  -ms-user-select: none !important;
  user-select: none !important;
  -webkit-touch-callout: none !important;
}

.watermark-pattern {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.watermark-row {
  position: absolute;
  width: 100%;
  height: 16%;
}

.watermark-text {
  position: absolute;
  transform: rotate(-25deg);
  font-size: 13px;
  font-family: monospace;
  color: #000;
  white-space: nowrap;
  text-shadow: 0 0 1px rgba(255,255,255,0.5);
}

@media print {
  .protected-content {
    display: none !important;
  }
}
</style>
