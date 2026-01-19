<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  blockIndex: number
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
    : 'w-full max-w-5xl mx-auto px-4'
})

// Items with fallback
const items = computed(() => {
  return props.block.data.items || []
})
</script>

<template>
  <section 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#ffffff', 
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3" 
        :style="{color: block.data.titleColor || styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title || 'Yang Akan Anda Dapatkan' }}
      </h2>
      <p 
        class="text-center mb-8 sm:mb-12 max-w-2xl mx-auto" 
        :style="{color: block.data.textColor || styles.textSecondaryColor}"
      >
        {{ block.data.subtitle }}
      </p>
      
      <!-- Benefits Grid -->
      <div class="grid sm:grid-cols-2 gap-4 sm:gap-6">
        <div 
          v-for="(item, idx) in items" 
          :key="idx" 
          class="flex items-start gap-4 p-4 sm:p-5 bg-gradient-to-br from-neutral-50 to-white rounded-2xl border border-neutral-100 hover:border-neutral-200 transition-all hover:shadow-lg group hover-lift animate-fade-in-up"
          :style="{animationDelay: `${animationDelay + (idx * 50)}ms`}"
        >
          <!-- Icon -->
          <div 
            class="w-12 h-12 sm:w-14 sm:h-14 rounded-2xl flex items-center justify-center flex-shrink-0 transition-transform group-hover:scale-110" 
            :style="{backgroundColor: (block.data.iconColor || block.data.accentColor || styles.primaryColor) + '15'}"
          >
            <svg 
              class="w-6 h-6 sm:w-7 sm:h-7" 
              :style="{color: block.data.iconColor || block.data.accentColor || styles.primaryColor}" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          
          <!-- Text -->
          <div class="flex-1 min-w-0">
            <h4 v-if="item.title" class="font-bold mb-1" :style="{color: block.data.titleColor || styles.textPrimaryColor}">
              {{ item.title }}
            </h4>
            <p :style="{color: block.data.textColor || styles.textSecondaryColor}" class="text-sm sm:text-base leading-relaxed">
              {{ item.text || item.description }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
