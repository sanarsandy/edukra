<template>
  <Transition name="fade">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-50 bg-black/95 flex flex-col"
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
          <!-- Page Navigation -->
          <div v-if="totalPages > 0" class="flex items-center gap-2 mr-4">
            <button 
              @click="prevPage" 
              :disabled="currentPage <= 1"
              class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg disabled:opacity-30 disabled:cursor-not-allowed"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
            </button>
            <span class="text-neutral-300 text-sm">
              {{ currentPage }} / {{ totalPages }}
            </span>
            <button 
              @click="nextPage" 
              :disabled="currentPage >= totalPages"
              class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg disabled:opacity-30 disabled:cursor-not-allowed"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
            </button>
          </div>
          <!-- Zoom Controls -->
          <button 
            @click="zoomOut" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg"
            title="Perkecil"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
            </svg>
          </button>
          <span class="text-neutral-300 text-sm min-w-[4rem] text-center">{{ Math.round(scale * 100) }}%</span>
          <button 
            @click="zoomIn" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg"
            title="Perbesar"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
          </button>
          <!-- Close Button -->
          <button 
            @click="close" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg ml-4"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Page Image Container -->
      <div 
        class="flex-1 overflow-auto flex items-start justify-center p-4 protected-content"
        @contextmenu.prevent
      >
        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center h-full">
          <div class="text-center text-white">
            <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
            <p class="text-sm opacity-80">Memuat halaman...</p>
          </div>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="flex items-center justify-center h-full">
          <div class="text-center text-white p-4">
            <svg class="w-12 h-12 mx-auto mb-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <p class="text-sm mb-4">{{ error }}</p>
            <button @click="loadDocument" class="px-4 py-2 bg-primary-600 rounded-lg text-sm hover:bg-primary-700 transition-colors">
              Coba Lagi
            </button>
          </div>
        </div>

        <!-- Page Image -->
        <div v-else class="relative">
          <img 
            v-if="currentPageUrl"
            :src="currentPageUrl"
            :style="{ transform: `scale(${scale})`, transformOrigin: 'top center' }"
            class="shadow-2xl transition-transform duration-200"
            @load="onImageLoad"
            @error="onImageError"
            @contextmenu.prevent
            @dragstart.prevent
          />
          
          <!-- Watermark Overlay -->
          <div 
            v-if="userEmail"
            class="absolute inset-0 pointer-events-none select-none overflow-hidden"
            :style="{ transform: `scale(${scale})`, transformOrigin: 'top center' }"
          >
            <div class="watermark-pattern" :style="{ opacity: 0.06 }">
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
        </div>
      </div>

      <!-- Footer -->
      <div class="p-3 bg-neutral-900 border-t border-neutral-800 flex justify-between items-center">
        <p class="text-xs text-neutral-500">
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          Dokumen ini dilindungi • PDF asli tidak dapat diunduh
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
  lessonId: string
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

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || ''

const loading = ref(true)
const error = ref<string | null>(null)
const totalPages = ref(0)
const currentPage = ref(1)
const scale = ref(1)
const pageLoading = ref(false)

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

// Get auth token
const getAuthToken = () => {
  const tokenCookie = useCookie('token')
  return tokenCookie.value || ''
}

// Current page image URL
const currentPageUrl = computed(() => {
  if (!props.lessonId || totalPages.value === 0) return ''
  const token = getAuthToken()
  return `${apiBase}/api/content/${props.lessonId}/pdf/page/${currentPage.value}?token=${token}`
})

// Load document info (page count)
const loadDocument = async () => {
  if (!props.lessonId) return
  
  loading.value = true
  error.value = null
  
  try {
    const token = getAuthToken()
    const response = await fetch(`${apiBase}/api/content/${props.lessonId}/pdf/pages`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Gagal memuat dokumen')
    }
    
    const data = await response.json()
    totalPages.value = data.total_pages
    currentPage.value = 1
  } catch (err: any) {
    console.error('Failed to load PDF info:', err)
    error.value = err.message || 'Gagal memuat dokumen'
  } finally {
    loading.value = false
  }
}

const onImageLoad = () => {
  pageLoading.value = false
}

const onImageError = () => {
  pageLoading.value = false
  error.value = 'Gagal memuat halaman'
}

const prevPage = () => {
  if (currentPage.value > 1) {
    pageLoading.value = true
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    pageLoading.value = true
    currentPage.value++
  }
}

const zoomIn = () => {
  if (scale.value < 2) {
    scale.value = Math.min(2, scale.value + 0.25)
  }
}

const zoomOut = () => {
  if (scale.value > 0.5) {
    scale.value = Math.max(0.5, scale.value - 0.25)
  }
}

const close = () => {
  emit('close')
}

const markComplete = () => {
  emit('complete')
}

// Handle keyboard shortcuts
const handleKeydown = (e: KeyboardEvent) => {
  // Block print shortcut
  if ((e.ctrlKey || e.metaKey) && e.key === 'p') {
    e.preventDefault()
    return
  }
  // Block save shortcut  
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    return
  }
  // Navigation
  if (e.key === 'ArrowLeft' || e.key === 'PageUp') {
    prevPage()
  } else if (e.key === 'ArrowRight' || e.key === 'PageDown') {
    nextPage()
  } else if (e.key === 'Escape') {
    close()
  }
}

// Watch for open state changes
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    loading.value = true
    error.value = null
    currentPage.value = 1
    scale.value = 1
    loadDocument()
    document.addEventListener('keydown', handleKeydown)
  } else {
    document.removeEventListener('keydown', handleKeydown)
    totalPages.value = 0
    currentPage.value = 1
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

.protected-content img {
  pointer-events: none;
  -webkit-user-drag: none;
  user-drag: none;
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
  color: #333;
  white-space: nowrap;
  text-shadow: 0 0 2px rgba(255,255,255,0.8);
}

/* Hide when printing */
@media print {
  * {
    display: none !important;
  }
}
</style>
