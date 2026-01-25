<template>
  <div 
    class="relative w-full h-full bg-black"
    @contextmenu.prevent
  >
    <!-- YouTube IFrame (hidden behind overlay) -->
    <div ref="playerContainer" class="absolute inset-0">
      <div :id="playerId" class="w-full h-full"></div>
    </div>

    <!-- Permanent transparent overlay - always blocks right click -->
    <div 
      class="absolute inset-0 z-10"
      @contextmenu.prevent
      @click="togglePlayPause"
      @dblclick="toggleFullscreen"
    >
      <!-- Play/Pause indicator (shows briefly on click) -->
      <transition name="fade">
        <div 
          v-if="showPlayIndicator"
          class="absolute inset-0 flex items-center justify-center bg-black/20 pointer-events-none"
        >
          <div class="w-20 h-20 bg-black/50 rounded-full flex items-center justify-center">
            <svg v-if="isPlaying" class="w-10 h-10 text-white" fill="currentColor" viewBox="0 0 24 24">
              <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z"/>
            </svg>
            <svg v-else class="w-10 h-10 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
              <path d="M8 5v14l11-7z"/>
            </svg>
          </div>
        </div>
      </transition>
    </div>

    <!-- Bottom controls bar -->
    <div 
      class="absolute bottom-0 left-0 right-0 z-20 bg-gradient-to-t from-black/80 to-transparent p-4"
      @contextmenu.prevent
    >
      <div class="flex items-center justify-between text-white">
        <!-- Play/Pause button -->
        <button 
          @click.stop="togglePlayPause"
          class="p-2 hover:bg-white/20 rounded-full transition-colors"
        >
          <svg v-if="isPlaying" class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
            <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z"/>
          </svg>
          <svg v-else class="w-6 h-6 ml-0.5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </button>

        <!-- Time display -->
        <span class="text-sm font-mono">{{ formatTime(currentTime) }} / {{ formatTime(duration) }}</span>

        <!-- Fullscreen button -->
        <button 
          @click.stop="toggleFullscreen"
          class="p-2 hover:bg-white/20 rounded-full transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div 
      v-if="!playerReady"
      class="absolute inset-0 z-30 flex items-center justify-center bg-black"
    >
      <div class="text-center text-white">
        <div class="w-12 h-12 border-4 border-white/30 border-t-white rounded-full animate-spin mx-auto mb-3"></div>
        <p class="text-sm opacity-80">Memuat video...</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps<{
  videoId: string
  autoplay?: boolean
}>()

const emit = defineEmits<{
  (e: 'play'): void
  (e: 'pause'): void
  (e: 'ended'): void
}>()

const playerId = ref(`youtube-player-${Date.now()}`)
const playerContainer = ref<HTMLElement | null>(null)
const player = ref<any>(null)
const playerReady = ref(false)
const isPlaying = ref(false)
const showPlayIndicator = ref(false)
const currentTime = ref(0)
const duration = ref(0)

let timeUpdateInterval: ReturnType<typeof setInterval> | null = null
let indicatorTimeout: ReturnType<typeof setTimeout> | null = null

// Load YouTube IFrame API
const loadYouTubeAPI = (): Promise<void> => {
  return new Promise((resolve) => {
    if ((window as any).YT && (window as any).YT.Player) {
      resolve()
      return
    }

    const existingScript = document.querySelector('script[src*="youtube.com/iframe_api"]')
    if (existingScript) {
      // API is loading, wait for it
      const checkReady = setInterval(() => {
        if ((window as any).YT && (window as any).YT.Player) {
          clearInterval(checkReady)
          resolve()
        }
      }, 100)
      return
    }

    const tag = document.createElement('script')
    tag.src = 'https://www.youtube.com/iframe_api'
    const firstScript = document.getElementsByTagName('script')[0]
    firstScript.parentNode?.insertBefore(tag, firstScript)

    ;(window as any).onYouTubeIframeAPIReady = () => {
      resolve()
    }
  })
}

// Initialize player
const initPlayer = async () => {
  await loadYouTubeAPI()

  player.value = new (window as any).YT.Player(playerId.value, {
    videoId: props.videoId,
    playerVars: {
      autoplay: props.autoplay ? 1 : 0,
      controls: 0, // Hide native controls
      disablekb: 1, // Disable keyboard
      fs: 0, // Disable fullscreen button
      iv_load_policy: 3, // No annotations
      modestbranding: 1,
      rel: 0, // No related videos
      showinfo: 0,
      origin: window.location.origin
    },
    events: {
      onReady: onPlayerReady,
      onStateChange: onPlayerStateChange
    }
  })
}

const onPlayerReady = (event: any) => {
  playerReady.value = true
  duration.value = event.target.getDuration() || 0
  
  // Start time update interval
  timeUpdateInterval = setInterval(() => {
    if (player.value && player.value.getCurrentTime) {
      currentTime.value = player.value.getCurrentTime() || 0
    }
  }, 1000)
}

const onPlayerStateChange = (event: any) => {
  const state = event.data
  const YT = (window as any).YT

  if (state === YT.PlayerState.PLAYING) {
    isPlaying.value = true
    emit('play')
  } else if (state === YT.PlayerState.PAUSED) {
    isPlaying.value = false
    emit('pause')
  } else if (state === YT.PlayerState.ENDED) {
    isPlaying.value = false
    emit('ended')
  }
}

const togglePlayPause = () => {
  if (!player.value) return

  if (isPlaying.value) {
    player.value.pauseVideo()
  } else {
    player.value.playVideo()
  }

  // Show indicator
  showPlayIndicator.value = true
  if (indicatorTimeout) clearTimeout(indicatorTimeout)
  indicatorTimeout = setTimeout(() => {
    showPlayIndicator.value = false
  }, 500)
}

const toggleFullscreen = () => {
  const container = playerContainer.value?.parentElement
  if (!container) return

  if (document.fullscreenElement) {
    document.exitFullscreen()
  } else {
    container.requestFullscreen()
  }
}

const formatTime = (seconds: number): string => {
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// Watch for videoId changes
watch(() => props.videoId, (newId) => {
  if (player.value && player.value.loadVideoById) {
    player.value.loadVideoById(newId)
  }
})

onMounted(() => {
  initPlayer()
})

onUnmounted(() => {
  if (timeUpdateInterval) clearInterval(timeUpdateInterval)
  if (indicatorTimeout) clearTimeout(indicatorTimeout)
  if (player.value && player.value.destroy) {
    player.value.destroy()
  }
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
