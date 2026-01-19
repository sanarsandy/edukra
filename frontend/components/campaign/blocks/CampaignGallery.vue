<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  courseImages?: string[] | null
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

// Images from block data or course
const images = computed(() => {
  return props.block.data.images || props.courseImages || []
})

// Lightbox state
const lightboxOpen = ref(false)
const lightboxIndex = ref(0)

const openLightbox = (idx: number) => {
  lightboxIndex.value = idx
  lightboxOpen.value = true
}

const closeLightbox = () => {
  lightboxOpen.value = false
}

const nextImage = () => {
  lightboxIndex.value = (lightboxIndex.value + 1) % images.value.length
}

const prevImage = () => {
  lightboxIndex.value = (lightboxIndex.value - 1 + images.value.length) % images.value.length
}
</script>

<template>
  <section 
    v-if="images.length"
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#fafafa',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        v-if="block.data.title"
        class="text-2xl sm:text-3xl font-bold text-center mb-8" 
        :style="{color: styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title }}
      </h2>
      
      <!-- Gallery Grid -->
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3 sm:gap-4">
        <div 
          v-for="(image, idx) in images" 
          :key="idx"
          class="aspect-square rounded-xl overflow-hidden cursor-pointer group hover-lift"
          @click="openLightbox(idx)"
        >
          <img 
            :src="image" 
            :alt="`Gallery ${idx + 1}`"
            class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-110"
          />
        </div>
      </div>
    </div>
    
    <!-- Lightbox -->
    <Teleport to="body">
      <div 
        v-if="lightboxOpen"
        class="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4"
        @click.self="closeLightbox"
      >
        <!-- Close Button -->
        <button 
          @click="closeLightbox"
          class="absolute top-4 right-4 w-10 h-10 bg-white/10 hover:bg-white/20 rounded-full flex items-center justify-center text-white transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
        
        <!-- Prev Button -->
        <button 
          v-if="images.length > 1"
          @click="prevImage"
          class="absolute left-4 w-12 h-12 bg-white/10 hover:bg-white/20 rounded-full flex items-center justify-center text-white transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
        </button>
        
        <!-- Image -->
        <img 
          :src="images[lightboxIndex]" 
          class="max-w-full max-h-[85vh] object-contain rounded-lg"
        />
        
        <!-- Next Button -->
        <button 
          v-if="images.length > 1"
          @click="nextImage"
          class="absolute right-4 w-12 h-12 bg-white/10 hover:bg-white/20 rounded-full flex items-center justify-center text-white transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
          </svg>
        </button>
        
        <!-- Counter -->
        <div class="absolute bottom-4 left-1/2 -translate-x-1/2 text-white/60 text-sm">
          {{ lightboxIndex + 1 }} / {{ images.length }}
        </div>
      </div>
    </Teleport>
  </section>
</template>
