<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0
})

// Default trust badges
const badges = computed(() => {
  return props.block.data.badges || [
    { icon: '🔒', text: 'Pembayaran Aman' },
    { icon: '💯', text: 'Garansi 30 Hari' },
    { icon: '🎓', text: 'Sertifikat' },
    { icon: '♾️', text: 'Akses Selamanya' }
  ]
})
</script>

<template>
  <section 
    class="py-8 sm:py-10 border-y border-neutral-100 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#ffffff',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div class="max-w-5xl mx-auto px-4">
      <div class="flex flex-wrap items-center justify-center gap-6 sm:gap-10">
        <div 
          v-for="(badge, idx) in badges" 
          :key="idx"
          class="flex items-center gap-2 text-neutral-600"
        >
          <span class="text-xl sm:text-2xl">{{ badge.icon }}</span>
          <span class="text-sm font-medium">{{ badge.text }}</span>
        </div>
      </div>
    </div>
  </section>
</template>
