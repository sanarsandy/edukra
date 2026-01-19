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
    : 'w-full max-w-3xl mx-auto px-4'
})

// Items with fallback
const items = computed(() => {
  return props.block.data.items || []
})

// Track open state
const openIndex = ref<number | null>(0)

const toggle = (idx: number) => {
  openIndex.value = openIndex.value === idx ? null : idx
}
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
        {{ block.data.title || 'Pertanyaan Umum' }}
      </h2>
      <p class="text-center mb-8 sm:mb-12" :style="{color: block.data.textColor || styles.textSecondaryColor}">
        {{ block.data.subtitle || 'Temukan jawaban untuk pertanyaan yang sering diajukan' }}
      </p>
      
      <!-- FAQ Accordion -->
      <div class="space-y-3 max-w-2xl mx-auto">
        <div 
          v-for="(item, idx) in items" 
          :key="idx" 
          class="border rounded-xl overflow-hidden transition-all"
          :class="openIndex === idx ? 'border-neutral-300 shadow-sm' : 'border-neutral-200'"
        >
          <!-- Question -->
          <button 
            @click="toggle(idx)"
            class="w-full text-left p-4 sm:p-5 flex items-center justify-between gap-4 transition-colors"
            :class="openIndex === idx ? 'bg-neutral-50' : 'bg-white hover:bg-neutral-50'"
            :style="{backgroundColor: openIndex === idx ? (block.data.backgroundColor || '#f9fafb') : '#ffffff'}"
          >
            <span class="font-semibold" :style="{color: block.data.titleColor || styles.textPrimaryColor}">
              {{ item.question }}
            </span>
            <svg 
              class="w-5 h-5 flex-shrink-0 transition-transform duration-200" 
              :class="openIndex === idx ? 'rotate-180' : ''"
              :style="{color: block.data.accentColor || styles.primaryColor}"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>
          
          <!-- Answer -->
          <div 
            class="overflow-hidden transition-all duration-200"
            :class="openIndex === idx ? 'max-h-96' : 'max-h-0'"
          >
            <div class="p-4 sm:p-5 pt-0 text-sm sm:text-base leading-relaxed" :style="{color: block.data.textColor || styles.textSecondaryColor}">
              {{ item.answer }}
            </div>
          </div>
        </div>
        
        <!-- Empty state -->
        <div v-if="!items.length" class="text-center text-neutral-400 py-8">
          Belum ada FAQ
        </div>
      </div>
    </div>
  </section>
</template>
