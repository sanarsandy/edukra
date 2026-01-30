<template>
  <div class="secure-pdf-viewer-container">
    <!-- Inline Mode: No modal wrapper -->
    <div 
      v-if="inlineMode && isOpen" 
      class="h-full flex flex-col bg-neutral-50"
      @contextmenu.prevent
    >
    <!-- Inline Header -->
    <div class="flex items-center justify-between px-4 py-2 bg-white border-b border-neutral-200">
      <div class="flex items-center gap-3">
        <svg class="w-5 h-5 text-red-500" viewBox="0 0 24 24" fill="currentColor">
          <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6zm-1 2l5 5h-5V4z"/>
        </svg>
        <span class="text-neutral-900 font-medium truncate max-w-md text-sm">{{ title }}</span>
      </div>
      <div class="flex items-center gap-2">
        <!-- Page Navigation -->
        <div class="flex items-center gap-2" v-if="totalPages > 0">
          <button 
            @click="prevPage" 
            :disabled="currentPage <= 1"
            class="p-1.5 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
          </button>
          <span class="text-neutral-600 text-xs">
            {{ currentPage }} / {{ totalPages }}
          </span>
          <button 
            @click="nextPage" 
            :disabled="currentPage >= totalPages"
            class="p-1.5 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
        
        <div class="h-4 w-px bg-neutral-200 mx-1"></div>
        
        <!-- Zoom Controls -->
        <button 
          @click="zoomOut" 
          :disabled="zoomLevel <= 0"
          class="p-1.5 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg disabled:opacity-30"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
          </svg>
        </button>
        <span class="text-neutral-600 text-xs min-w-[3rem] text-center">{{ zoomOptions[zoomLevel].label }}</span>
        <button 
          @click="zoomIn" 
          :disabled="zoomLevel >= zoomOptions.length - 1"
          class="p-1.5 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg disabled:opacity-30"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
        </button>
        <button 
          @click="fitToWidth" 
          class="p-1.5 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg"
          title="Sesuaikan Lebar"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/>
          </svg>
        </button>
      </div>
    </div>
    
    <!-- PDF Container for Inline Mode -->
    <div 
      ref="containerRef"
      class="flex-1 relative overflow-auto bg-neutral-100 flex justify-center protected-content"
      :class="loading ? 'items-center' : 'items-start'"
    >
      <!-- Loading -->
      <div v-if="loading" class="absolute inset-0 flex items-center justify-center z-10 bg-neutral-100">
        <div class="text-center text-neutral-600">
          <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
          <p class="text-sm">Memuat dokumen...</p>
        </div>
      </div>
      
      <!-- Error -->
      <div v-else-if="error" class="absolute inset-0 flex items-center justify-center z-10">
        <div class="text-center text-red-500 p-4">
          <svg class="w-12 h-12 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <p class="text-sm">{{ error }}</p>
        </div>
      </div>
      
      <!-- PDF Render -->
      <div v-show="!loading && !error" class="py-6">
        <VuePdfEmbed
          v-if="effectivePdfUrl"
          ref="pdfEmbed"
          :source="effectivePdfUrl"
          :page="currentPage"
          :width="pdfWidth"
          @loaded="onLoaded"
          @rendered="onRendered"
          @loading-failed="onError"
        />
      </div>
      
      <!-- Watermark -->
      <div 
        v-if="userEmail && !loading && !error" 
        class="absolute inset-0 pointer-events-none select-none overflow-hidden z-20"
      >
        <div class="absolute inset-0 flex items-center justify-center" style="transform: rotate(-30deg);">
          <div class="flex flex-wrap gap-40 opacity-[0.08]">
            <span v-for="i in 20" :key="i" class="text-neutral-900 text-sm font-medium whitespace-nowrap">
              {{ maskedEmail }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
  
  <!-- Modal Mode: Original behavior -->
  <Transition v-else name="fade">
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
            :disabled="zoomLevel <= 0"
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg disabled:opacity-30"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"/>
            </svg>
          </button>
          <span class="text-neutral-300 text-sm min-w-[4rem] text-center">{{ zoomOptions[zoomLevel].label }}</span>
          <button 
            @click="zoomIn" 
            :disabled="zoomLevel >= zoomOptions.length - 1"
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg disabled:opacity-30"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
          </button>
          
          <!-- Fit Width Button -->
          <button 
            @click="fitToWidth" 
            class="p-2 text-neutral-400 hover:text-white hover:bg-neutral-800 rounded-lg ml-2"
            title="Sesuaikan Lebar"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5v-4m0 4h-4m4 0l-5-5"/>
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
      <div 
        ref="containerRef"
        class="flex-1 relative overflow-auto bg-neutral-800 flex justify-center protected-content"
        :class="loading ? 'items-center' : 'items-start'"
      >
        <!-- Loading State -->
        <div v-if="loading" class="absolute inset-0 flex items-center justify-center z-10 bg-neutral-800">
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

        <!-- PDF Viewer with Native Width -->
        <div 
          class="relative my-6 shadow-2xl bg-white"
          :style="{ width: pdfWidth + 'px' }"
        >
          <ClientOnly>
            <VuePdfEmbed
              v-if="effectivePdfUrl"
              ref="pdfEmbed"
              :source="effectivePdfUrl"
              :page="currentPage"
              :width="pdfWidth"
              @loaded="onLoaded"
              @rendered="onRendered"
              @loading-failed="onError"
            />
          </ClientOnly>

          <!-- Watermark Overlay (always on top of canvas) -->
          <div 
            v-if="userEmail"
            class="absolute inset-0 pointer-events-none select-none overflow-hidden z-20"
          >
            <div class="watermark-pattern">
              <template v-for="row in 12" :key="row">
                <div class="watermark-row" :style="{ top: `${row * 8 - 4}%` }">
                  <template v-for="col in 6" :key="col">
                    <span 
                      class="watermark-text"
                      :style="{ left: `${col * 16 - 8}%` }"
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
  </div>
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
  inlineMode?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Dokumen',
  isCompleted: false,
  inlineMode: false
})

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'complete'): void
}>()

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Get auth token from cookie (Nuxt way)
const tokenCookie = useCookie('token')

// Zoom options with preset widths for better rendering quality
const zoomOptions = [
  { label: '50%', width: 400 },
  { label: '75%', width: 600 },
  { label: '100%', width: 800 },
  { label: '125%', width: 1000 },
  { label: '150%', width: 1200 },
  { label: '175%', width: 1400 },
  { label: '200%', width: 1600 },
]

const loading = ref(true)
const error = ref<string | null>(null)
const currentPage = ref(1)
const totalPages = ref(0)
const zoomLevel = ref(3) // Default to 125% (index 3)
const pdfEmbed = ref(null)
const containerRef = ref<HTMLElement | null>(null)
const scale = ref(1.25) // For backward compatibility
const proxiedPdfUrl = ref<string | null>(null)

// Computed PDF width based on zoom level
const pdfWidth = computed(() => zoomOptions[zoomLevel.value].width)

// Determine if URL is external (needs proxy)
const isExternalUrl = (url: string): boolean => {
  if (!url) return false
  // External URLs start with http:// or https:// and are not from our domain
  if (url.startsWith('http://') || url.startsWith('https://')) {
    try {
      const urlObj = new URL(url)
      const currentHost = typeof window !== 'undefined' ? window.location.hostname : ''
      // If it's from a different domain, it's external
      return urlObj.hostname !== currentHost && 
             !url.includes(apiBase) && 
             !urlObj.hostname.includes('localhost')
    } catch {
      return true
    }
  }
  return false
}

// Get the actual URL to use for the PDF viewer
const effectivePdfUrl = computed(() => {
  // If we have a proxied URL (blob), use that
  if (proxiedPdfUrl.value) return proxiedPdfUrl.value
  
  // If this is an external URL and we don't have a proxied URL yet, return null
  // This prevents pdf.js from trying to fetch the external URL directly
  if (props.pdfUrl && isExternalUrl(props.pdfUrl)) {
    return null
  }
  
  // For local/internal URLs, use directly
  return props.pdfUrl
})

// Fetch external PDF through proxy
const fetchExternalPdf = async (url: string) => {
  try {
    loading.value = true
    error.value = null
    
    // Get auth token from cookie
    const token = tokenCookie.value
    if (!token) {
      error.value = 'Sesi login tidak valid. Silakan login ulang.'
      loading.value = false
      return
    }
    
    const response = await fetch(`${apiBase}/api/content/proxy-pdf`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ url })
    })
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}))
      throw new Error(errorData.error || `Failed to fetch PDF (${response.status})`)
    }
    
    // Create blob URL from response
    const blob = await response.blob()
    proxiedPdfUrl.value = URL.createObjectURL(blob)
  } catch (err: any) {
    console.error('Failed to proxy PDF:', err)
    error.value = err.message || 'Gagal memuat dokumen PDF eksternal.'
    loading.value = false
  }
}

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
  if (zoomLevel.value < zoomOptions.length - 1) {
    zoomLevel.value++
    scale.value = zoomOptions[zoomLevel.value].width / 800
  }
}

const zoomOut = () => {
  if (zoomLevel.value > 0) {
    zoomLevel.value--
    scale.value = zoomOptions[zoomLevel.value].width / 800
  }
}

// Fit PDF to container width
const fitToWidth = () => {
  if (containerRef.value) {
    const containerWidth = containerRef.value.clientWidth - 80 // Account for padding
    // Find closest zoom level
    let closestIdx = 0
    let minDiff = Math.abs(zoomOptions[0].width - containerWidth)
    
    for (let i = 1; i < zoomOptions.length; i++) {
      const diff = Math.abs(zoomOptions[i].width - containerWidth)
      if (diff < minDiff) {
        minDiff = diff
        closestIdx = i
      }
    }
    zoomLevel.value = closestIdx
    scale.value = zoomOptions[closestIdx].width / 800
  }
}

const close = () => {
  // Cleanup blob URL if created
  if (proxiedPdfUrl.value) {
    URL.revokeObjectURL(proxiedPdfUrl.value)
    proxiedPdfUrl.value = null
  }
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
  if (e.key === '+' || e.key === '=') { e.preventDefault(); zoomIn() }
  if (e.key === '-') { e.preventDefault(); zoomOut() }
  if (e.key === 'Escape') close()
}

watch(() => props.isOpen, async (isOpen) => {
  if (isOpen) {
    document.addEventListener('keydown', handleKeydown)
    // Reset state
    currentPage.value = 1
    loading.value = true
    error.value = null
    proxiedPdfUrl.value = null
    
    // Check if external URL needs proxy
    if (props.pdfUrl && isExternalUrl(props.pdfUrl)) {
      await fetchExternalPdf(props.pdfUrl)
    }
    
    // Auto fit to width on next tick
    setTimeout(() => {
      fitToWidth()
    }, 100)
  } else {
    document.removeEventListener('keydown', handleKeydown)
    // Cleanup blob URL
    if (proxiedPdfUrl.value) {
      URL.revokeObjectURL(proxiedPdfUrl.value)
      proxiedPdfUrl.value = null
    }
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  // Cleanup blob URL
  if (proxiedPdfUrl.value) {
    URL.revokeObjectURL(proxiedPdfUrl.value)
  }
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
  font-size: 12px;
  font-family: monospace;
  color: rgba(128, 128, 128, 0.25);
  white-space: nowrap;
  font-weight: 600;
  letter-spacing: 1px;
}

@media print {
  * { display: none !important; }
}
</style>
