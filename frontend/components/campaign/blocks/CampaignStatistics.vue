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
    : 'w-full max-w-5xl mx-auto px-4'
})

// Stats items
const items = computed(() => {
  return props.block.data.items || []
})
</script>

<template>
  <section 
    class="py-12 sm:py-16 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#ffffff',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        v-if="block.data.title"
        class="text-xl sm:text-2xl font-bold text-center mb-8" 
        :style="{color: block.data.titleColor || styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title }}
      </h2>
      
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4 sm:gap-6">
        <div 
          v-for="(stat, idx) in items" 
          :key="idx"
          class="text-center p-4 sm:p-6 bg-neutral-50 rounded-2xl hover-lift transition-all"
        >
          <div 
            class="text-3xl sm:text-4xl font-black mb-1" 
            :style="{color: block.data.numberColor || block.data.accentColor || styles.primaryColor}"
          >
            {{ stat.value }}{{ stat.suffix || '' }}
          </div>
          <div 
            class="text-xs sm:text-sm font-medium" 
            :style="{color: block.data.textColor || styles.textSecondaryColor}"
          >
            {{ stat.label }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
