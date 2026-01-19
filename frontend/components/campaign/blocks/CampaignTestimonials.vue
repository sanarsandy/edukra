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
    : 'w-full max-w-6xl mx-auto px-4'
})

// Items with fallback
const items = computed(() => {
  return props.block.data.items || []
})

// Check if grid variant
const isGrid = computed(() => props.block.variant === 'grid')
</script>

<template>
  <section 
    v-if="items.length" 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#fafafa', 
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3" 
        :style="{color: block.data.titleColor || styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title || 'Apa Kata Mereka?' }}
      </h2>
      <p class="text-center mb-8 sm:mb-12" :style="{color: block.data.textColor || styles.textSecondaryColor}">
        Kisah sukses alumni kami
      </p>
      
      <!-- Variant: Grid View -->
      <div v-if="isGrid" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <div 
          v-for="(item, idx) in items" 
          :key="idx" 
          class="bg-white p-6 rounded-2xl shadow-sm border border-neutral-100 hover:shadow-lg transition-all h-full flex flex-col testimonial-card animate-fade-in-up"
          :style="{animationDelay: `${animationDelay + (idx * 100)}ms`}"
        >
          <!-- Stars -->
          <div class="flex gap-1 mb-4">
            <span v-for="i in 5" :key="i" class="text-yellow-400">⭐</span>
          </div>
          
          <!-- Quote -->
          <p class="text-neutral-700 mb-4 flex-1 italic leading-relaxed">
            "{{ item.text }}"
          </p>
          
          <!-- Author -->
          <div class="flex items-center gap-3 mt-auto pt-4 border-t border-neutral-100">
            <div class="w-10 h-10 rounded-full bg-neutral-200 flex items-center justify-center text-lg">
              {{ item.name?.charAt(0) || '?' }}
            </div>
            <div>
              <div class="font-semibold text-neutral-900">{{ item.name }}</div>
              <div v-if="item.role" class="text-sm text-neutral-500">{{ item.role }}</div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Variant: Carousel (Default) -->
      <div v-else class="flex gap-4 sm:gap-6 overflow-x-auto pb-4 snap-x snap-mandatory scrollbar-hide sm:grid sm:grid-cols-2 lg:grid-cols-3 sm:overflow-visible">
        <div 
          v-for="(item, idx) in items" 
          :key="idx" 
          class="flex-shrink-0 w-[85vw] sm:w-auto snap-center bg-white p-5 sm:p-6 rounded-2xl shadow-sm border border-neutral-100 hover:shadow-lg transition-all testimonial-card animate-fade-in-up"
          :style="{animationDelay: `${animationDelay + (idx * 100)}ms`}"
        >
          <!-- Stars -->
          <div class="flex gap-1 mb-4">
            <span v-for="i in 5" :key="i" class="text-yellow-400">⭐</span>
          </div>
          
          <!-- Quote -->
          <p class="text-neutral-700 mb-4 italic leading-relaxed line-clamp-4">
            "{{ item.text }}"
          </p>
          
          <!-- Author -->
          <div class="flex items-center gap-3 pt-4 border-t border-neutral-100">
            <div class="w-10 h-10 rounded-full bg-neutral-200 flex items-center justify-center text-lg">
              {{ item.name?.charAt(0) || '?' }}
            </div>
            <div class="min-w-0">
              <div class="font-semibold text-neutral-900 truncate">{{ item.name }}</div>
              <div v-if="item.role" class="text-sm text-neutral-500 truncate">{{ item.role }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
