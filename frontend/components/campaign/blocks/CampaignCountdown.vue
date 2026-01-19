<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block?: Block
  countdown: { days: number; hours: number; minutes: number; seconds: number }
  styles: CampaignStyles
  label?: string
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0,
  label: '⏰ Penawaran berakhir dalam:'
})

// Get color with fallback to global styles
const numberColor = computed(() => props.block?.data?.numberColor || props.styles.primaryColor)
const labelColor = computed(() => props.block?.data?.labelColor || 'rgba(255,255,255,0.6)')
const bgColor = computed(() => props.block?.data?.backgroundColor || props.styles.backgroundColor)

const countdownItems = computed(() => ({
  Hari: props.countdown.days,
  Jam: props.countdown.hours,
  Menit: props.countdown.minutes,
  Detik: props.countdown.seconds
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
