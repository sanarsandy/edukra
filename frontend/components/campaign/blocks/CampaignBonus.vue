<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  isMobileLayout?: boolean
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0
})

// Container class
const containerClass = computed(() => {
  return props.isMobileLayout 
    ? 'w-full max-w-[480px] mx-auto px-4' 
    : 'w-full max-w-4xl mx-auto px-4'
})

// Bonus items
const items = computed(() => {
  return props.block.data.items || []
})

// Format price
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('id-ID').format(price || 0)
}
</script>

<template>
  <section 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#fffbeb',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <!-- Header -->
      <div class="text-center mb-8">
        <span 
          class="inline-block text-xs font-bold px-3 py-1 rounded-full mb-3"
          :style="{
            backgroundColor: block.data.accentColor || '#facc15', 
            color: '#78350f'
          }"
        >
          🎁 BONUS
        </span>
        <h2 
          class="text-2xl sm:text-3xl font-bold" 
          :style="{color: styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
        >
          {{ block.data.title || 'Bonus Eksklusif' }}
        </h2>
        <p v-if="block.data.subtitle" class="mt-2" :style="{color: styles.textSecondaryColor}">
          {{ block.data.subtitle }}
        </p>
      </div>
      
      <!-- Bonus Items -->
      <div class="space-y-4">
        <div 
          v-for="(bonus, idx) in items" 
          :key="idx"
          class="flex items-center gap-4 p-4 sm:p-5 bg-white rounded-2xl shadow-sm border border-orange-100 hover:shadow-lg transition-all hover-lift animate-fade-in-up"
          :style="{animationDelay: `${animationDelay + (idx * 50)}ms`}"
        >
          <!-- Icon -->
          <div class="w-12 h-12 sm:w-14 sm:h-14 rounded-xl bg-orange-100 flex items-center justify-center text-2xl flex-shrink-0">
            {{ bonus.emoji || '🎁' }}
          </div>
          
          <!-- Content -->
          <div class="flex-1 min-w-0">
            <h4 class="font-bold text-neutral-900 text-sm sm:text-base">{{ bonus.title }}</h4>
            <p v-if="bonus.description" class="text-neutral-500 text-xs sm:text-sm mt-0.5 line-clamp-2">
              {{ bonus.description }}
            </p>
          </div>
          
          <!-- Value Badge -->
          <div class="text-green-600 font-bold text-xs sm:text-sm whitespace-nowrap">
            <span v-if="bonus.value">Rp {{ formatPrice(bonus.value) }}</span>
            <span v-else>GRATIS</span>
          </div>
        </div>
      </div>
      
      <!-- Total Value -->
      <div 
        v-if="block.data.total_value"
        class="mt-8 p-5 sm:p-6 bg-gradient-to-r from-green-500 to-emerald-600 rounded-2xl text-white text-center"
      >
        <p class="text-sm opacity-80 mb-1">Total Nilai Bonus</p>
        <p class="text-2xl sm:text-3xl font-black">Rp {{ formatPrice(block.data.total_value) }}</p>
      </div>
    </div>
  </section>
</template>
