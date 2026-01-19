<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  block: any
  styles?: any
}>()

const defaultItems = [
    { title: 'Dasar Kepabeanan', text: 'Memahami konsep Bea dan Cukai dengan bahasa sederhana', icon: '📘', image: 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&q=80&w=400' },
    { title: 'Prosedur Impor', text: 'Langkah-langkah impor barang dan dokumen yang diperlukan', icon: '📦', image: 'https://images.unsplash.com/photo-1578575437130-527eed3abbec?auto=format&fit=crop&q=80&w=400' },
    { title: 'Perhitungan Pajak', text: 'Cara menghitung bea masuk, PPN, dan PPh impor', icon: '💰', image: 'https://images.unsplash.com/photo-1554224155-6726b3ff858f?auto=format&fit=crop&q=80&w=400' },
    { title: 'Studi Kasus', text: 'Contoh nyata belanja online dan pengiriman barang', icon: '✅', image: 'https://images.unsplash.com/photo-1553413077-190dd305871c?auto=format&fit=crop&q=80&w=400' }
]

const items = computed(() => {
    return (props.block.data.items && props.block.data.items.length > 0) ? props.block.data.items : defaultItems
})
</script>

<template>
  <section class="py-20 px-4 bg-slate-50 font-inter">
      <div class="max-w-6xl mx-auto">
          <!-- Header -->
          <div class="text-center mb-14">
              <div class="inline-flex items-center gap-2 bg-amber-100 text-amber-700 rounded-full px-4 py-2 text-sm font-semibold mb-4">
                  <span>📚</span>
                  <span>Materi Webinar</span>
              </div>
              <h2 class="text-3xl lg:text-4xl font-bold text-slate-800 mb-4">
                  {{ block.data.title || 'Apa yang Akan Kamu Pelajari?' }}
              </h2>
              <p class="text-slate-600 max-w-2xl mx-auto text-lg">
                  {{ block.data.subtitle || 'Materi disusun khusus untuk pemula dengan bahasa yang mudah dipahami' }}
              </p>
          </div>
          
          <!-- Feature Grid with Images -->
          <div class="grid md:grid-cols-2 gap-6">
              <div v-for="(item, idx) in items" :key="idx" 
                   class="bg-white rounded-2xl overflow-hidden border border-slate-100 shadow-sm hover:shadow-xl transition-all group">
                  <!-- Image -->
                  <div class="h-40 overflow-hidden relative">
                      <img 
                        :src="item.image || 'https://images.unsplash.com/photo-1586528116311-ad8dd3c8310d?auto=format&fit=crop&q=80&w=400'"
                        :alt="item.title"
                        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                      />
                      <div class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
                      <div class="absolute bottom-3 left-4 text-4xl">{{ item.icon }}</div>
                  </div>
                  
                  <!-- Content -->
                  <div class="p-6">
                      <h3 class="text-xl font-bold text-slate-800 mb-2">{{ item.title }}</h3>
                      <p class="text-slate-600 leading-relaxed">{{ item.text || item.desc }}</p>
                  </div>
              </div>
          </div>
          
          <!-- Bottom Banner -->
          <div class="mt-12 bg-gradient-to-r from-blue-900 to-blue-800 rounded-2xl p-8 flex flex-col md:flex-row items-center justify-between gap-6">
              <div class="flex items-center gap-4">
                  <div class="w-16 h-16 bg-amber-500/20 rounded-2xl flex items-center justify-center text-4xl">
                      🎓
                  </div>
                  <div class="text-white">
                      <p class="font-bold text-xl">Bonus: E-Sertifikat</p>
                      <p class="text-blue-200">Dapatkan sertifikat kehadiran untuk semua peserta</p>
                  </div>
              </div>
              <a href="#register" class="bg-amber-500 hover:bg-amber-400 text-white font-semibold px-6 py-3 rounded-xl transition-colors whitespace-nowrap">
                  Daftar Gratis →
              </a>
          </div>
      </div>
  </section>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&display=swap');

.font-inter {
    font-family: 'Inter', sans-serif;
}
</style>
