<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()

const defaultLearnings = [
  {
    icon: '🧠',
    title: 'Upgrade Sistem Operasi Mental',
    description: 'Pelajari cara menggunakan pikiran dewasa untuk merapikan peta mental yang dipakai hari ini'
  },
  {
    icon: '🔍',
    title: 'Kenali Sumber Luka',
    description: 'Pahami akar masalah dari pola emosi dan perilaku yang terus berulang'
  },
  {
    icon: '🛠️',
    title: 'Teknik Praktis Pemulihan',
    description: 'Dapatkan tools yang bisa langsung dipraktikkan setelah webinar selesai'
  },
  {
    icon: '🌱',
    title: 'Langkah Pertama Healing',
    description: 'Mulai perjalanan pemulihan dengan fondasi yang tepat dan terarah'
  }
]

const learnings = props.block.data?.learnings || defaultLearnings
</script>

<template>
  <section class="relative py-16 md:py-24 bg-[#f5f0e8]">
    <div class="max-w-5xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-12">
        <p class="text-sm font-medium text-[#6b635a] uppercase tracking-widest mb-3">
          {{ block.data?.eyebrow || 'Di Webinar P.U.L.I.H' }}
        </p>
        <h2 class="text-3xl md:text-4xl font-serif text-[#4a4540] mb-4">
          {{ block.data?.headline || 'Yang Akan Kamu Pelajari' }}
        </h2>
        <p class="text-lg text-[#6b635a] max-w-2xl mx-auto">
          {{ block.data?.subheadline || 'Pulih bukan berarti memorinya hilang, tapi sistem operasi mental lo yang di-upgrade.' }}
        </p>
      </div>

      <!-- Learning Cards -->
      <div class="grid md:grid-cols-2 gap-6 mb-12">
        <div 
          v-for="(item, index) in learnings" 
          :key="index"
          class="relative bg-white rounded-2xl p-6 shadow-md border border-[#e8e4dc] overflow-hidden group hover:shadow-lg transition-all"
        >
          <!-- Decorative Corner -->
          <div class="absolute top-0 right-0 w-24 h-24 bg-gradient-to-br from-amber-100/50 to-transparent rounded-bl-full"></div>
          
          <div class="relative">
            <div class="text-4xl mb-4 group-hover:scale-110 transition-transform duration-300">
              {{ item.icon }}
            </div>
            <h3 class="font-semibold text-[#4a4540] text-xl mb-2">{{ item.title }}</h3>
            <p class="text-[#6b635a] leading-relaxed">{{ item.description }}</p>
          </div>
        </div>
      </div>

      <!-- Key Message -->
      <div class="bg-gradient-to-r from-[#4a4540] to-[#5a534a] rounded-2xl p-8 text-center text-white shadow-xl">
        <div class="text-5xl mb-4">✨</div>
        <p class="text-xl md:text-2xl font-light leading-relaxed max-w-2xl mx-auto">
          {{ block.data?.key_message || 'Kita akan bahas gimana cara pakai pikiran dewasa lo buat ngerapihin peta mental yang lo pakai hari ini.' }}
        </p>
      </div>
    </div>
  </section>
</template>

<style scoped>
.font-serif {
  font-family: 'Playfair Display', 'Georgia', serif;
}
</style>
