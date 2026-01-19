<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  buttonStyle: Record<string, any>
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
}>()

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

// Container class
const containerClass = computed(() => {
  return props.isMobileLayout 
    ? 'w-full max-w-[480px] mx-auto px-4' 
    : 'w-full max-w-lg mx-auto px-4'
})

// Section background
const sectionBackground = computed(() => {
  const bgColor = props.block.data.backgroundColor || props.styles.backgroundColor
  return `linear-gradient(135deg, ${bgColor}, ${adjustColor(bgColor, -30)})`
})

// Features list
const features = computed(() => {
  return props.block.data.features || [
    'Akses materi seumur hidup',
    'Sertifikat kelulusan',
    'Grup komunitas eksklusif',
    'Update materi gratis'
  ]
})

const handleBuy = () => {
  emit('buy')
}
</script>

<template>
  <section 
    id="pricing" 
    class="py-12 sm:py-20 text-white animate-fade-in-up" 
    :style="{
      background: sectionBackground,
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-2" :style="{fontFamily: styles.fontFamilyHeading}">
        🎁 Penawaran Spesial
      </h2>
      <p class="text-center text-white/70 mb-8">Khusus untuk Anda yang serius ingin berkembang</p>
      
      <!-- Variant: Highlight (Card) - Default -->
      <div v-if="block.variant !== 'minimal'" class="bg-white/10 backdrop-blur-xl rounded-3xl p-6 sm:p-8 border border-white/20 relative overflow-hidden animate-scale-in">
        <!-- Popular Badge -->
        <div 
          class="absolute -top-1 -right-8 text-xs font-bold px-8 py-1 rotate-45 transform origin-center" 
          :style="{backgroundColor: block.data.accentColor || styles.accentColor, color: block.data.accentColor ? '#000' : '#78350f'}"
        >
          TERLARIS
        </div>
        
        <!-- Course/Webinar Name -->
        <h3 class="text-xl sm:text-2xl font-bold text-center mb-4">{{ block.data.product_name || 'Paket Lengkap' }}</h3>
        
        <!-- Price Section -->
        <div class="text-center my-6">
          <div v-if="coursePrice && coursePrice > (displayPrice || 0)" class="text-white/50 line-through text-lg mb-1">
            Rp {{ formatPrice(coursePrice) }}
          </div>
          <div class="flex items-baseline justify-center gap-2">
            <span class="text-4xl sm:text-5xl font-black" :style="{color: styles.primaryColor}">
              Rp {{ formatPrice(displayPrice || 0) }}
            </span>
            <span v-if="!isFree" class="text-white/60 text-sm">/sekali bayar</span>
          </div>
          <div v-if="coursePrice && coursePrice > (displayPrice || 0)" class="mt-2 inline-block bg-green-500/20 text-green-300 px-3 py-1 rounded-full text-sm font-medium">
            Hemat {{ Math.round(((coursePrice - (displayPrice || 0)) / coursePrice) * 100) }}%
          </div>
        </div>
        
        <!-- Features -->
        <ul class="space-y-3 mb-8">
          <li v-for="(feature, idx) in features" :key="idx" class="flex items-center gap-3 text-white/90">
            <svg class="w-5 h-5 flex-shrink-0" :style="{color: styles.primaryColor}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
            </svg>
            <span class="text-sm sm:text-base">{{ feature }}</span>
          </li>
        </ul>
        
        <!-- CTA Button -->
        <button 
          @click="handleBuy" 
          class="w-full py-4 sm:py-5 text-white text-lg sm:text-xl font-bold rounded-2xl transition-all transform hover:scale-[1.02] active:scale-[0.98] shadow-xl soft-shadow-btn" 
          :style="buttonStyle"
        >
          {{ isFree ? '🎉 Daftar Gratis' : (block.data.cta_text || '🔥 Daftar Sekarang') }}
        </button>
        
        <!-- Trust Badge -->
        <p class="text-center text-white/60 text-sm mt-4 flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
          </svg>
          Garansi 30 Hari Uang Kembali
        </p>
      </div>
      
      <!-- Variant: Minimal -->
      <div v-else class="text-center">
        <div v-if="coursePrice && coursePrice > (displayPrice || 0)" class="text-white/50 line-through text-xl mb-2">
          Rp {{ formatPrice(coursePrice) }}
        </div>
        <div class="text-5xl sm:text-6xl font-black mb-6" :style="{color: styles.primaryColor}">
          Rp {{ formatPrice(displayPrice || 0) }}
        </div>
        <button 
          @click="handleBuy" 
          class="w-full sm:w-auto px-12 py-5 text-white text-xl font-bold rounded-2xl transition-all transform hover:scale-105 active:scale-95 shadow-2xl" 
          :style="buttonStyle"
        >
          {{ isFree ? '🎉 Daftar Gratis' : (block.data.cta_text || '🚀 Daftar Sekarang') }}
        </button>
      </div>
    </div>
  </section>
</template>
