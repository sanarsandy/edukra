<template>
  <Transition name="fade">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-50 bg-black/90 flex flex-col"
      @contextmenu.prevent
      @keydown.prevent="handleKeydown"
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
          <div class="flex items-center gap-2 mr-4">
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

      <!-- PDF Canvas Container -->
      <div 
        ref="containerRef"
        class="flex-1 overflow-auto flex items-start justify-center p-4 protected-content"
        @scroll="onScroll"
      >
        <div v-if="loading" class="flex items-center justify-center h-full">
          <div class="text-center text-white">
            <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
            <p class="text-sm opacity-80">Memuat dokumen...</p>
          </div>
        </div>

        <div v-else-if="error" class="flex items-center justify-center h-full">
          <div class="text-center text-white p-4">
            <svg class="w-12 h-12 mx-auto mb-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <p class="text-sm">{{ error }}</p>
            <button @click="loadDocument" class="mt-3 px-4 py-2 bg-primary-600 rounded-lg text-sm hover:bg-primary-700 transition-colors">
              Coba Lagi
            </button>
          </div>
        </div>

        <div v-else class="relative">
          <canvas ref="canvasRef" class="shadow-2xl"></canvas>
          
          <!-- Watermark Overlay -->
          <div 
            v-if="userEmail"
            class="absolute inset-0 pointer-events-none select-none overflow-hidden"
          >
            <div 
              class="watermark-pattern"
              :style="{ opacity: 0.08 }"
            >
              <template v-for="row in 5" :key="row">
                <div class="watermark-row" :style="{ top: `${row * 20 - 10}%` }">
                  <template v-for="col in 3" :key="col">
                    <span 
                      class="watermark-text"
                      :style="{ left: `${col * 33 - 16}%` }"
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
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import * as pdfjsLib from 'pdfjs-dist'

// Get version from the installed package and use matching worker from CDN
const pdfjsVersion = pdfjsLib.version
pdfjsLib.GlobalWorkerOptions.workerSrc = `https://cdnjs.cloudflare.com/ajax/libs/pdf.js/${pdfjsVersion}/pdf.worker.min.mjs`

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

const containerRef = ref<HTMLDivElement | null>(null)
const canvasRef = ref<HTMLCanvasElement | null>(null)

const loading = ref(true)
const error = ref<string | null>(null)
const pdfDoc = ref<any>(null)
const currentPage = ref(1)
const totalPages = ref(0)
const scale = ref(1.5)

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

const loadDocument = async () => {
  if (!props.pdfUrl) return
  
  loading.value = true
  error.value = null

  try {
    const loadingTask = pdfjsLib.getDocument(props.pdfUrl)
    pdfDoc.value = await loadingTask.promise
    totalPages.value = pdfDoc.value.numPages
    await renderPage(currentPage.value)
  } catch (err) {
    console.error('Failed to load PDF:', err)
    error.value = 'Gagal memuat dokumen'
  } finally {
    loading.value = false
  }
}

const renderPage = async (pageNum: number) => {
  if (!pdfDoc.value || !canvasRef.value) return

  try {
    const page = await pdfDoc.value.getPage(pageNum)
    const viewport = page.getViewport({ scale: scale.value })

    const canvas = canvasRef.value
    const context = canvas.getContext('2d')
    
    canvas.height = viewport.height
    canvas.width = viewport.width

    const renderContext = {
      canvasContext: context!,
      viewport: viewport
    }

    await page.render(renderContext).promise
  } catch (err) {
    console.error('Failed to render page:', err)
    error.value = 'Gagal menampilkan halaman'
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    renderPage(currentPage.value)
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    renderPage(currentPage.value)
  }
}

const zoomIn = () => {
  if (scale.value < 3) {
    scale.value = Math.min(3, scale.value + 0.25)
    renderPage(currentPage.value)
  }
}

const zoomOut = () => {
  if (scale.value > 0.5) {
    scale.value = Math.max(0.5, scale.value - 0.25)
    renderPage(currentPage.value)
  }
}

const close = () => {
  emit('close')
}

const markComplete = () => {
  emit('complete')
}

const onScroll = () => {
  // Could implement scroll-based page detection here
}

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
    nextTick(() => {
      loadDocument()
      document.addEventListener('keydown', handleKeydown)
    })
  } else {
    document.removeEventListener('keydown', handleKeydown)
    // Reset state
    currentPage.value = 1
    pdfDoc.value = null
    loading.value = true
    error.value = null
  }
})

// Watch for URL changes while open
watch(() => props.pdfUrl, () => {
  if (props.isOpen) {
    currentPage.value = 1
    loadDocument()
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
  height: 20%;
}

.watermark-text {
  position: absolute;
  transform: rotate(-30deg);
  font-size: 14px;
  font-family: monospace;
  color: #000;
  white-space: nowrap;
  text-shadow: 0 0 1px rgba(255,255,255,0.5);
}

/* Prevent dev tools inspection (basic) */
@media print {
  .protected-content {
    display: none !important;
  }
}
</style>
