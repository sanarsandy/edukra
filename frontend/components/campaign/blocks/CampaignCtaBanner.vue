<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  buttonStyle: Record<string, any>
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0
})

const emit = defineEmits<{
  (e: 'buy'): void
}>()

const handleBuy = () => {
  emit('buy')
}
</script>

<template>
  <section 
    class="py-10 sm:py-16 text-white text-center animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || styles.buttonColor,
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div class="max-w-4xl mx-auto px-4">
      <h2 
        class="text-2xl sm:text-3xl md:text-4xl font-bold mb-3" 
        :style="{fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.headline || '🚀 Siap Untuk Memulai?' }}
      </h2>
      <p class="text-white/80 mb-8 max-w-2xl mx-auto">
        {{ block.data.subheadline || 'Jangan lewatkan kesempatan untuk meningkatkan skill Anda' }}
      </p>
      <button 
        @click="handleBuy"
        class="px-8 sm:px-12 py-4 bg-white font-bold rounded-xl transition-all transform hover:scale-105 active:scale-95 shadow-xl"
        :style="{
          color: block.data.backgroundColor || styles.buttonColor,
          borderRadius: buttonStyle.borderRadius
        }"
      >
        {{ block.data.cta_text || 'Daftar Sekarang' }}
      </button>
    </div>
  </section>
</template>
