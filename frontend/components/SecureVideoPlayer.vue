<template>
  <div 
    class="secure-video-player relative bg-black rounded-lg overflow-hidden"
    @contextmenu.prevent
  >
    <!-- Loading State -->
    <div v-if="loading" class="absolute inset-0 flex items-center justify-center bg-black">
      <div class="text-center text-white">
        <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full mx-auto mb-3"></div>
        <p class="text-sm opacity-80">Memuat video...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-if="error" class="absolute inset-0 flex items-center justify-center bg-black">
      <div class="text-center text-white p-4">
        <svg class="w-12 h-12 mx-auto mb-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <p class="text-sm">{{ error }}</p>
        <button @click="retryLoad" class="mt-3 px-4 py-2 bg-primary-600 rounded-lg text-sm hover:bg-primary-700 transition-colors">
          Coba Lagi
        </button>
      </div>
    </div>

    <!-- Processing State -->
    <div v-if="processing" class="absolute inset-0 flex items-center justify-center bg-black">
      <div class="text-center text-white p-4">
        <svg class="w-12 h-12 mx-auto mb-3 text-yellow-500 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
        </svg>
        <p class="text-sm">Video sedang diproses...</p>
        <p class="text-xs opacity-60 mt-1">Silakan tunggu beberapa saat</p>
      </div>
    </div>

    <!-- Video Element -->
    <video
      ref="videoRef"
      class="w-full h-full"
      :class="{ 'opacity-0': loading || error || processing }"
      controls
      controlsList="nodownload noplaybackrate"
      disablepictureinpicture
      playsinline
      @ended="$emit('ended')"
      @timeupdate="onTimeUpdate"
      @loadedmetadata="onLoadedMetadata"
      @playing="isPlaying = true"
      @pause="isPlaying = false"
    ></video>

    <!-- Watermark Overlay (shows when playing) -->
    <VideoWatermark 
      v-if="showWatermark && userEmail"
      :user-email="userEmail"
      :opacity="0.12"
      :rotate-interval="25"
    />

    <!-- Protected Overlay (prevents some interactions) -->
    <div class="absolute top-0 left-0 right-0 h-1 pointer-events-none"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import Hls from 'hls.js'

interface Props {
  lessonId: string
  fallbackUrl?: string
  userEmail?: string
  showWatermark?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showWatermark: true
})

const emit = defineEmits<{
  (e: 'ended'): void
  (e: 'timeupdate', time: number): void
  (e: 'duration', duration: number): void
  (e: 'error', message: string): void
}>()

const videoRef = ref<HTMLVideoElement | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const processing = ref(false)
const hlsInstance = ref<Hls | null>(null)
const isPlaying = ref(false)

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Get auth token for HLS key requests
const getAuthToken = (): string => {
  if (typeof localStorage !== 'undefined') {
    return localStorage.getItem('token') || ''
  }
  return ''
}

const loadHLS = async () => {
  if (!props.lessonId) return
  
  loading.value = true
  error.value = null
  processing.value = false

  try {
    // First check HLS status
    const statusRes = await fetch(`${apiBase}/api/content/${props.lessonId}/hls/status`, {
      headers: {
        'Authorization': `Bearer ${getAuthToken()}`
      }
    })
    
    if (!statusRes.ok) {
      throw new Error('Failed to check HLS status')
    }

    const status = await statusRes.json()
    
    if (status.status === 'not_started' || status.status === 'pending') {
      // HLS not available, use fallback
      if (props.fallbackUrl) {
        loadDirectVideo(props.fallbackUrl)
      } else {
        error.value = 'Video belum tersedia'
      }
      return
    }

    if (status.status === 'processing') {
      processing.value = true
      loading.value = false
      // Poll for completion
      setTimeout(loadHLS, 5000)
      return
    }

    if (status.status !== 'ready') {
      throw new Error(`Unexpected status: ${status.status}`)
    }

    // HLS is ready, load the player
    const manifestUrl = `${apiBase}/api/content/${props.lessonId}/hls/manifest`
    
    if (Hls.isSupported()) {
      const hls = new Hls({
        xhrSetup: (xhr: XMLHttpRequest, url: string) => {
          // Add auth token to key requests
          if (url.includes('/hls/key')) {
            xhr.setRequestHeader('Authorization', `Bearer ${getAuthToken()}`)
          }
        },
        // Enable fetch for manifest with auth
        pLoader: class extends Hls.DefaultConfig.loader {
          constructor(config: any) {
            super(config)
          }
        }
      })

      hls.loadSource(manifestUrl)
      hls.attachMedia(videoRef.value!)

      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        loading.value = false
      })

      hls.on(Hls.Events.ERROR, (event, data) => {
        console.error('HLS Error:', data)
        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              error.value = 'Gagal memuat video. Periksa koneksi internet.'
              break
            case Hls.ErrorTypes.MEDIA_ERROR:
              hls.recoverMediaError()
              break
            default:
              error.value = 'Terjadi kesalahan saat memutar video'
              break
          }
          loading.value = false
        }
      })

      hlsInstance.value = hls
    } else if (videoRef.value?.canPlayType('application/vnd.apple.mpegurl')) {
      // Native HLS support (Safari)
      videoRef.value.src = manifestUrl
      loading.value = false
    } else {
      // No HLS support, use fallback
      if (props.fallbackUrl) {
        loadDirectVideo(props.fallbackUrl)
      } else {
        error.value = 'Browser tidak mendukung HLS'
      }
    }
  } catch (err) {
    console.error('Failed to load HLS:', err)
    // Try fallback
    if (props.fallbackUrl) {
      loadDirectVideo(props.fallbackUrl)
    } else {
      error.value = 'Gagal memuat video'
      loading.value = false
    }
  }
}

const loadDirectVideo = (url: string) => {
  if (videoRef.value) {
    videoRef.value.src = url
    loading.value = false
  }
}

const retryLoad = () => {
  error.value = null
  loadHLS()
}

const onTimeUpdate = () => {
  if (videoRef.value) {
    emit('timeupdate', videoRef.value.currentTime)
  }
}

const onLoadedMetadata = () => {
  if (videoRef.value) {
    emit('duration', videoRef.value.duration)
  }
}

// Cleanup
onUnmounted(() => {
  if (hlsInstance.value) {
    hlsInstance.value.destroy()
  }
})

// Watch for lessonId changes
watch(() => props.lessonId, () => {
  if (hlsInstance.value) {
    hlsInstance.value.destroy()
    hlsInstance.value = null
  }
  loadHLS()
})

onMounted(() => {
  loadHLS()
})
</script>

<style scoped>
.secure-video-player {
  min-height: 200px;
}

/* Hide download button */
video::-webkit-media-controls-download-button {
  display: none !important;
}

video::-webkit-media-controls-enclosure {
  overflow: hidden !important;
}

video::-webkit-media-controls-picture-in-picture-button {
  display: none !important;
}
</style>
