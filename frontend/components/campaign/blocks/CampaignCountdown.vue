<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block?: Block
  countdown?: { days: number; hours: number; minutes: number; seconds: number } | null
  styles: CampaignStyles
  label?: string
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0,
  label: '⏰ Penawaran berakhir dalam:'
})

// Internal countdown state
const internalCountdown = ref<{ days: number; hours: number; minutes: number; seconds: number } | null>(null)
let countdownInterval: ReturnType<typeof setInterval> | null = null

// Get end date from block.data or fall back to global countdown
const endDateString = computed(() => props.block?.data?.end_date || props.block?.data?.endDate)

// Start the countdown timer
const startTimer = () => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }

  const endDateStr = endDateString.value
  
  // If no end_date in block, use props.countdown directly
  if (!endDateStr) {
    internalCountdown.value = props.countdown || { days: 0, hours: 0, minutes: 0, seconds: 0 }
    return
  }

  const endDate = new Date(endDateStr).getTime()

  const updateCountdown = () => {
    const now = new Date().getTime()
    const distance = endDate - now

    if (distance <= 0) {
      internalCountdown.value = { days: 0, hours: 0, minutes: 0, seconds: 0 }
      if (countdownInterval) clearInterval(countdownInterval)
      return
    }

    internalCountdown.value = {
      days: Math.floor(distance / (1000 * 60 * 60 * 24)),
      hours: Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
      minutes: Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)),
      seconds: Math.floor((distance % (1000 * 60)) / 1000)
    }
  }

  updateCountdown()
  countdownInterval = setInterval(updateCountdown, 1000)
}

// Watch for changes in end date and restart timer
watch(endDateString, () => {
  startTimer()
}, { immediate: false })

// Watch props.countdown for fallback - sync continuously
watch(() => props.countdown, (newVal) => {
  // Only use props.countdown if block doesn't have its own end_date
  if (!endDateString.value && newVal) {
    internalCountdown.value = { ...newVal }
  }
}, { deep: true, immediate: true })

onMounted(() => {
  startTimer()
})

onUnmounted(() => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
})

// Get color with fallback to global styles
const numberColor = computed(() => props.block?.data?.numberColor || props.styles.primaryColor)
const labelColor = computed(() => props.block?.data?.labelColor || 'rgba(255,255,255,0.8)')
const bgColor = computed(() => props.block?.data?.backgroundColor || props.styles.backgroundColor)

// Use internal countdown or props fallback
const displayCountdown = computed(() => internalCountdown.value || props.countdown || { days: 0, hours: 0, minutes: 0, seconds: 0 })

const countdownItems = computed(() => ({
  Hari: displayCountdown.value.days,
  Jam: displayCountdown.value.hours,
  Menit: displayCountdown.value.minutes,
  Detik: displayCountdown.value.seconds
}))
</script>

<template>
  <section 
    class="py-8 sm:py-12 text-white animate-fade-in-up" 
    :style="{
      backgroundColor: bgColor, 
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div class="max-w-4xl mx-auto px-4 text-center">
      <p class="text-base sm:text-lg mb-4 sm:mb-6" :style="{color: labelColor}">
        {{ block?.data?.label || label }}
      </p>
      <div class="flex justify-center gap-2 sm:gap-4">
        <div 
          v-for="(value, labelItem) in countdownItems" 
          :key="labelItem" 
          class="bg-black/30 backdrop-blur rounded-xl p-3 sm:p-5 min-w-[60px] sm:min-w-[80px] countdown-unit"
        >
          <div 
            class="text-2xl sm:text-4xl md:text-5xl font-bold tabular-nums countdown-number" 
            :style="{color: numberColor}"
          >
            {{ String(value).padStart(2, '0') }}
          </div>
          <div class="text-xs sm:text-sm mt-1" :style="{color: labelColor}">{{ labelItem }}</div>
        </div>
      </div>
    </div>
  </section>
</template>

