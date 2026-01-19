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

// Get YouTube embed URL
const embedUrl = computed(() => {
  const url = props.block.data.youtube_url
  if (!url) return null
  
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]+)/,
    /youtube\.com\/shorts\/([a-zA-Z0-9_-]+)/
  ]
  
  for (const pattern of patterns) {
    const match = url.match(pattern)
    if (match && match[1]) {
      return `https://www.youtube.com/embed/${match[1]}`
    }
  }
  
  return null
})
</script>

<template>
  <section 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#111827',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        v-if="block.data.title"
        class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-3 text-white" 
        :style="{fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title }}
      </h2>
      <p 
        v-if="block.data.subtitle"
        class="text-center mb-8 sm:mb-10 text-white/70"
      >
        {{ block.data.subtitle }}
      </p>
      
      <!-- Video Embed -->
      <div class="max-w-3xl mx-auto">
        <div 
          v-if="embedUrl"
          class="relative aspect-video rounded-2xl overflow-hidden shadow-2xl bg-black"
        >
          <iframe
            :src="embedUrl"
            class="absolute inset-0 w-full h-full"
            frameborder="0"
            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
          ></iframe>
        </div>
        
        <!-- Placeholder when no video -->
        <div 
          v-else
          class="aspect-video rounded-2xl bg-neutral-800 flex items-center justify-center"
        >
          <div class="text-center text-white/50">
            <div class="text-5xl mb-3">🎬</div>
            <p class="text-sm">Tambahkan YouTube URL untuk menampilkan video</p>
          </div>
        </div>
      </div>
      
      <!-- Optional Caption -->
      <p 
        v-if="block.data.caption"
        class="text-center text-white/60 text-sm mt-4"
      >
        {{ block.data.caption }}
      </p>
    </div>
  </section>
</template>
