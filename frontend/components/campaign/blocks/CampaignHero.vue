<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  buttonStyle: Record<string, any>
  countdown?: { days: number; hours: number; minutes: number; seconds: number } | null
  displayPrice?: number
  coursePrice?: number
  isFree?: boolean
  isMobileLayout?: boolean
  animationDelay?: number
  preview?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0,
  preview: false
})

const emit = defineEmits<{
  (e: 'buy'): void
  (e: 'scrollTo', target: string): void
}>()

// Computed helper for container class
const containerClass = computed(() => {
  return props.isMobileLayout 
    ? 'w-full max-w-[480px] mx-auto px-4' 
    : 'w-full max-w-4xl mx-auto px-4'
})

const containerClassWide = computed(() => {
  return props.isMobileLayout 
    ? 'w-full max-w-[480px] mx-auto px-4' 
    : 'w-full max-w-7xl mx-auto px-4'
})

// Hero style for background
const heroStyle = computed(() => {
  const bgImage = props.block.data.background_image
  return {
    backgroundImage: bgImage 
      ? `url(${bgImage})` 
      : `linear-gradient(135deg, ${props.styles.backgroundColor}, ${adjustColor(props.styles.backgroundColor, -30)})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    minHeight: props.block.data.minHeight ? `${props.block.data.minHeight}px` : undefined
  }
})

// Adjust color helper
const adjustColor = (hex: string, percent: number): string => {
  if (!hex || !hex.startsWith('#')) return hex
  const num = parseInt(hex.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + percent))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + percent))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + percent))
  return `#${(1 << 24 | r << 16 | g << 8 | b).toString(16).slice(1)}`
}

// Format price
const formatPrice = (price: number): string => {
  return new Intl.NumberFormat('id-ID').format(price || 0)
}

// Compute per-block button style with fallback
const blockButtonStyle = computed(() => {
  const base = { ...props.buttonStyle }
  if (props.block.data.buttonColor) {
    base.backgroundColor = props.block.data.buttonColor
  }
  if (props.block.data.buttonTextColor) {
    base.color = props.block.data.buttonTextColor
  }
  return base
})

const handleBuy = () => {
  emit('buy')
}

const scrollTo = (target: string) => {
  emit('scrollTo', target)
}
</script>

<template>
  <section 
    :class="[
      block.variant === 'bio_profile' 
        ? 'relative w-full bg-white pb-12 overflow-hidden' 
        : 'relative min-h-[100svh] flex items-center justify-center text-white overflow-hidden', 
      'animate-fade-in-up'
    ]" 
    :style="{
      ...(block.variant === 'bio_profile' ? {} : heroStyle), 
      animationDelay: `${animationDelay}ms`
    }"
  >
    <!-- Gradient Overlay (Only for non-bio layout) -->
    <div v-if="block.variant !== 'bio_profile'" class="absolute inset-0 bg-gradient-to-b from-black/70 via-black/50 to-black/80"></div>
    
    <!-- Urgency Badge -->
    <div v-if="countdown && !preview" class="absolute top-4 left-1/2 -translate-x-1/2 z-20 animate-bounce-in">
      <div 
        class="bg-red-500 text-white px-4 py-2 rounded-full text-sm font-bold animate-pulse-glow flex items-center gap-2 shadow-lg hover:scale-105 transition-transform cursor-pointer" 
        @click="scrollTo('#pricing')"
      >
        <span>🔥</span> Promo Berakhir dalam {{ countdown.days }}h {{ countdown.hours }}m
      </div>
    </div>
    
    <!-- Content: Bio Link Profile Variant -->
    <div v-if="block.variant === 'bio_profile'" class="w-full">
       <!-- Cover Image -->
       <div 
         class="h-48 md:h-64 w-full bg-cover bg-center rounded-b-[40px] relative overflow-hidden" 
         :style="{backgroundImage: `url(${block.data.background_image || 'https://via.placeholder.com/800x400'})`}"
       >
         <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent"></div>
       </div>
       
       <!-- Profile Info -->
       <div class="px-4 -mt-16 flex flex-col items-center relative z-20">
          <!-- Avatar -->
          <div class="p-1 bg-white rounded-full shadow-lg">
             <img 
               :src="block.data.profile_image || 'https://via.placeholder.com/150'" 
               class="w-32 h-32 rounded-full object-cover border-4 border-white bg-neutral-100" 
               alt="Profile" 
             />
          </div>
          
          <!-- Handle -->
          <div v-if="block.data.handle" class="mt-3 text-sm font-bold text-neutral-600 bg-white/90 backdrop-blur px-4 py-1.5 rounded-full shadow-sm border border-neutral-100 text-center min-w-[120px]">
             {{ block.data.handle.startsWith('@') ? block.data.handle : '@' + block.data.handle }}
          </div>
          
          <!-- Name & Bio -->
          <div class="mt-6 text-center max-w-md mx-auto px-4">
             <h1 class="text-2xl font-bold text-neutral-900 mb-3 leading-tight">{{ block.data.headline }}</h1>
             <p class="text-neutral-600 leading-relaxed text-sm md:text-base whitespace-pre-line">{{ block.data.subheadline }}</p>
          </div>

          <!-- Call to Action (Optional for Bio) -->
          <div v-if="block.data.cta_link" class="mt-6 w-full max-w-xs">
              <a :href="block.data.cta_link" class="block w-full text-center py-3 bg-primary-600 text-white font-semibold rounded-xl hover:bg-primary-700 transition-colors shadow-lg shadow-primary-500/20">
                  {{ block.data.cta_text || 'Contact Me' }}
              </a>
          </div>
       </div>
    </div>

    <!-- Content: Split Variant -->
    <div v-else-if="block.variant === 'split'" :class="[containerClassWide, 'relative z-10 py-20 sm:py-32 grid lg:grid-cols-2 gap-12 items-center']">
       <div class="text-left">
          <!-- Badge -->
          <div v-if="block.data.badge" class="inline-block bg-white/20 backdrop-blur px-4 py-1.5 rounded-full text-sm font-medium mb-6 animate-fade-in-up" :style="{color: block.data.textColor || '#ffffff'}">
            {{ block.data.badge }}
          </div>
          <!-- Headline -->
          <h1 class="text-3xl sm:text-5xl lg:text-6xl font-extrabold mb-6 leading-tight tracking-tight animate-fade-in-up delay-100" :style="{color: block.data.titleColor || '#ffffff'}">
            {{ block.data.headline }}
          </h1>
          <!-- Subheadline -->
          <p class="text-lg sm:text-xl mb-8 leading-relaxed animate-fade-in-up delay-200" :style="{color: block.data.textColor || 'rgba(255,255,255,0.9)'}">
            {{ block.data.subheadline }}
          </p>
          <!-- Price & CTA -->
          <div class="animate-fade-in-up delay-300">
             <div v-if="displayPrice" class="mb-6 flex items-baseline gap-3">
               <span class="text-3xl sm:text-4xl font-bold" :style="{color: styles.primaryColor}">Rp {{ formatPrice(displayPrice) }}</span>
               <span v-if="coursePrice && coursePrice > displayPrice" class="text-white/50 line-through text-lg">Rp {{ formatPrice(coursePrice) }}</span>
             </div>
             <div class="flex flex-col sm:flex-row gap-4">
               <button 
                 @click="handleBuy" 
                 class="px-8 py-4 text-white text-lg font-bold rounded-xl transition-all transform hover:scale-105 active:scale-95 shadow-xl" 
                 :style="blockButtonStyle"
               >
                 {{ isFree ? '🎉 Daftar Gratis' : (block.data.cta_text || '🚀 Daftar Sekarang') }}
               </button>
               <p class="text-white/60 text-sm flex items-center gap-2 mt-2 sm:mt-0">
                 <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                 Garansi 30 Hari
               </p>
             </div>
          </div>
       </div>
       <!-- Empty div for spacing on mobile, or could be image -->
       <div class="hidden lg:block relative h-full min-h-[400px]"></div>
    </div>

    <!-- Content: Centered Variant (Default) -->
    <div v-else :class="[containerClass, 'relative z-10 text-center py-16 sm:py-24']">
      <!-- Badge -->
      <div v-if="block.data.badge" class="inline-block bg-white/20 backdrop-blur px-4 py-1.5 rounded-full text-sm font-medium mb-6" :style="{color: block.data.textColor || '#ffffff'}">
        {{ block.data.badge }}
      </div>
      
      <!-- Headline -->
      <h1 class="text-3xl sm:text-4xl md:text-5xl lg:text-6xl font-extrabold mb-4 sm:mb-6 leading-tight tracking-tight" :style="{color: block.data.titleColor || '#ffffff'}">
        {{ block.data.headline }}
      </h1>
      
      <!-- Subheadline -->
      <p class="text-base sm:text-lg md:text-xl mb-6 sm:mb-8 max-w-2xl mx-auto leading-relaxed" :style="{color: block.data.textColor || 'rgba(255,255,255,0.9)'}">
        {{ block.data.subheadline }}
      </p>
      
      <!-- Price Preview -->
      <div v-if="displayPrice" class="mb-6">
        <span v-if="coursePrice && coursePrice > displayPrice" class="text-white/50 line-through text-lg mr-2">Rp {{ formatPrice(coursePrice) }}</span>
        <span class="text-2xl sm:text-3xl font-bold" :style="{color: styles.primaryColor}">Rp {{ formatPrice(displayPrice) }}</span>
      </div>
      
      <!-- CTA Button -->
      <button 
        @click="handleBuy" 
        class="w-full sm:w-auto px-8 sm:px-12 py-4 sm:py-5 text-white text-lg sm:text-xl font-bold rounded-2xl transition-all transform hover:scale-105 active:scale-95 shadow-2xl" 
        :style="blockButtonStyle"
      >
        {{ isFree ? '🎉 Daftar Gratis' : (block.data.cta_text || '🚀 Daftar Sekarang') }}
      </button>
      
      <!-- Trust Indicator -->
      <!-- Trust Indicator -->
      <p v-if="block.data.trust_text !== undefined ? block.data.trust_text : true" class="mt-4 text-white/60 text-sm">
        {{ block.data.trust_text || '✓ Akses Seumur Hidup • ✓ Garansi 30 Hari' }}
      </p>
    </div>
    
    <!-- Scroll Indicator -->
    <div v-if="!preview" class="absolute bottom-6 left-1/2 -translate-x-1/2 animate-bounce hidden sm:block">
      <svg class="w-6 h-6 text-white/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3"/>
      </svg>
    </div>
  </section>
</template>
