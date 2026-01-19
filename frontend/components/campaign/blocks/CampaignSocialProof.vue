<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0
})

// Mock notifications for social proof
const notifications = ref<{name: string; action: string; time: string}[]>([])
const currentNotification = ref<{name: string; action: string; time: string} | null>(null)

// Default sample names
const sampleNames = [
  'Budi S.', 'Rina A.', 'Ahmad F.', 'Dewi P.', 'Joko W.',
  'Siti N.', 'Agus M.', 'Maya R.', 'Hendra K.', 'Putri L.'
]

const sampleCities = [
  'Jakarta', 'Bandung', 'Surabaya', 'Yogyakarta', 'Semarang',
  'Medan', 'Makassar', 'Bali', 'Malang', 'Solo'
]

// Generate random notification
const generateNotification = () => {
  const name = sampleNames[Math.floor(Math.random() * sampleNames.length)]
  const city = sampleCities[Math.floor(Math.random() * sampleCities.length)]
  const minutes = Math.floor(Math.random() * 30) + 1
  
  return {
    name: `${name} dari ${city}`,
    action: props.block.data.action_text || 'baru saja mendaftar',
    time: `${minutes} menit yang lalu`
  }
}

// Show notification periodically
let notificationInterval: ReturnType<typeof setInterval> | null = null

const showNotification = () => {
  currentNotification.value = generateNotification()
  setTimeout(() => {
    currentNotification.value = null
  }, 5000)
}

onMounted(() => {
  // Show first notification after delay
  setTimeout(showNotification, 3000)
  
  // Then show every 15-30 seconds
  notificationInterval = setInterval(() => {
    showNotification()
  }, (Math.random() * 15000) + 15000)
})

onUnmounted(() => {
  if (notificationInterval) clearInterval(notificationInterval)
})
</script>

<template>
  <Teleport to="body">
    <Transition name="slide-in">
      <div 
        v-if="currentNotification && block.enabled !== false"
        class="fixed bottom-4 left-4 z-40 bg-white rounded-xl shadow-2xl border border-neutral-200 p-4 max-w-xs animate-fade-in-up"
      >
        <div class="flex items-start gap-3">
          <!-- Avatar -->
          <div 
            class="w-10 h-10 rounded-full flex items-center justify-center text-lg flex-shrink-0"
            :style="{backgroundColor: styles.primaryColor + '20'}"
          >
            👤
          </div>
          
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <p class="font-semibold text-sm text-neutral-900 truncate">
              {{ currentNotification.name }}
            </p>
            <p class="text-sm text-neutral-600">
              {{ currentNotification.action }}
            </p>
            <p class="text-xs text-neutral-400 mt-1">
              {{ currentNotification.time }}
            </p>
          </div>
          
          <!-- Close -->
          <button 
            @click="currentNotification = null" 
            class="text-neutral-400 hover:text-neutral-600 p-1"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.slide-in-enter-active,
.slide-in-leave-active {
  transition: all 0.3s ease;
}
.slide-in-enter-from {
  opacity: 0;
  transform: translateX(-100%);
}
.slide-in-leave-to {
  opacity: 0;
  transform: translateX(-100%);
}
</style>
