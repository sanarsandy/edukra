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

// Comparison items
const items = computed(() => {
  return props.block.data.items || [
    { feature: 'Akses Materi', ours: '✅ Seumur Hidup', others: '❌ 1 Tahun' },
    { feature: 'Sertifikat', ours: '✅ Ya', others: '❌ Tidak' },
    { feature: 'Support', ours: '✅ 24/7', others: '⏰ Jam Kerja' },
    { feature: 'Update Materi', ours: '✅ Gratis', others: '💰 Bayar Lagi' }
  ]
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
        class="text-2xl sm:text-3xl font-bold text-center mb-8" 
        :style="{color: styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title || 'Kenapa Memilih Kami?' }}
      </h2>
      
      <!-- Comparison Table -->
      <div class="overflow-x-auto">
        <table class="w-full min-w-[500px]">
          <thead>
            <tr>
              <th class="text-left p-4" :style="{color: styles.textSecondaryColor}">Fitur</th>
              <th 
                class="text-center p-4 rounded-t-xl"
                :style="{backgroundColor: styles.primaryColor, color: '#ffffff'}"
              >
                {{ block.data.ours_label || '🏆 Kursus Kami' }}
              </th>
              <th class="text-center p-4 bg-neutral-100" :style="{color: styles.textSecondaryColor}">
                {{ block.data.others_label || 'Kompetitor' }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, idx) in items" :key="idx" class="border-b border-neutral-100">
              <td class="p-4 font-medium" :style="{color: styles.textPrimaryColor}">
                {{ item.feature }}
              </td>
              <td 
                class="p-4 text-center font-semibold"
                :style="{backgroundColor: styles.primaryColor + '10'}"
              >
                {{ item.ours }}
              </td>
              <td class="p-4 text-center bg-neutral-50 text-neutral-500">
                {{ item.others }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </section>
</template>
