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
          <div class="flex items-center gap-2 mr-4" v-if="totalPages > 0">
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
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
            </svg>
          </button>
          <span class="text-neutral-300 text-sm min-w-[3rem] text-center">{{ Math.round(scale * 100) }}%</span>
          <button 
            @click="zoomIn" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg"
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

      <!-- PDF Viewer Container -->
      <div class="flex-1 relative overflow-auto bg-neutral-900 flex justify-center p-8 protected-content">
        <!-- Loading State -->
        <div v-if="loading" class="absolute inset-0 flex items-center justify-center z-10">
          <div class="text-center text-white">
            <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
            <p class="text-sm opacity-80">Memuat dokumen...</p>
          </div>
        </div>

        <!-- Error State -->
        <div v-if="error" class="absolute inset-0 flex items-center justify-center z-10">
          <div class="text-center text-white p-4">
            <svg class="w-12 h-12 mx-auto mb-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
            <p class="text-sm mb-2">{{ error }}</p>
          </div>
        </div>

        <!-- Canvas based PDF Viewer -->
        <div class="relative shadow-2xl" :style="{ transform: `scale(${scale})`, transformOrigin: 'top center' }">
          <ClientOnly>
            <VuePdfEmbed
              v-if="pdfUrl"
              ref="pdfEmbed"
              :source="pdfUrl"
              :page="currentPage"
              @loaded="onLoaded"
              @rendered="onRendered"
              @load-error="onError"
            />
          </ClientOnly>

          <!-- Watermark Overlay (always on top of canvas) -->
          <div 
            v-if="userEmail"
            class="absolute inset-0 pointer-events-none select-none overflow-hidden z-20"
          >
            <div class="watermark-pattern">
              <template v-for="row in 8" :key="row">
                <div class="watermark-row" :style="{ top: `${row * 12 - 6}%` }">
                  <template v-for="col in 5" :key="col">
                    <span 
                      class="watermark-text"
                      :style="{ left: `${col * 20 - 10}%` }"
                    >
                      {{ maskedEmail }}
                    </span>
                  </template>
                </div>
              </template>
            </div>
          </div>

          <!-- Interaction Blocker: blocks context menu, drag, selection -->
          <div 
            class="absolute inset-0 z-30"
            style="cursor: default;"
            @contextmenu.prevent
            @dragstart.prevent
            @selectstart.prevent
          ></div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-3 bg-neutral-900 border-t border-neutral-800 flex justify-between items-center">
        <p class="text-xs text-neutral-500">
          <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
          </svg>
          Dokumen dilindungi
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
import VuePdfEmbed from 'vue-pdf-embed'

// IMPORTANT: Import via specific path or just global to avoid TS issues if types missing
import * as pdfjsLib from 'pdfjs-dist'

// Set worker source to CDN to avoid bundler issues in Nuxt
if (typeof window !== 'undefined' && !pdfjsLib.GlobalWorkerOptions.workerSrc) {
  const version = pdfjsLib.version || '3.11.174'
  pdfjsLib.GlobalWorkerOptions.workerSrc = `https://cdn.jsdelivr.net/npm/pdfjs-dist@${version}/build/pdf.worker.min.js`
}

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
const currentPage = ref(1)
const totalPages = ref(0)
const scale = ref(1.2)
const pdfEmbed = ref(null)

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

const onLoaded = (doc: any) => {
  totalPages.value = doc.numPages
  loading.value = false
  error.value = null
}

const onRendered = () => {
  loading.value = false
}

const onError = (e: any) => {
  console.error('PDF Error:', e)
  loading.value = false
  error.value = 'Gagal memuat dokumen PDF. Pastikan file valid.'
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loading.value = true
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loading.value = true
  }
}

const zoomIn = () => {
  if (scale.value < 3) scale.value += 0.2
}

const zoomOut = () => {
  if (scale.value > 0.5) scale.value -= 0.2
}

const close = () => {
  emit('close')
}

const markComplete = () => {
  emit('complete')
}

// Keyboard shortcuts block
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && (e.key === 'p' || e.key === 's')) {
    e.preventDefault() // Block print/save
    return
  }
  if (e.key === 'ArrowRight') nextPage()
  if (e.key === 'ArrowLeft') prevPage()
  if (e.key === 'Escape') close()
}

watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
    // Reset state
    currentPage.value = 1
    loading.value = true
  } else {
    document.removeEventListener('keydown', handleKeydown)
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }

.protected-content {
  -webkit-user-select: none;
  user-select: none;
}

.watermark-pattern {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.watermark-row {
  position: absolute;
  width: 100%;
  height: 50px;
}

.watermark-text {
  position: absolute;
  transform: rotate(-25deg);
  font-size: 14px;
  font-family: monospace;
  color: rgba(0, 0, 0, 0.15); /* More visible watermark */
  white-space: nowrap;
  font-weight: bold;
}

@media print {
  * { display: none !important; }
}
</style>
