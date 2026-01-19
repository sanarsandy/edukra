<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()
</script>

<template>
  <section class="py-20 px-4 bg-white font-inter">
      <div class="max-w-5xl mx-auto">
          <div class="bg-gradient-to-br from-slate-50 to-slate-100 rounded-3xl p-8 lg:p-12 border border-slate-200">
              <div class="grid lg:grid-cols-3 gap-8 items-center">
                  <!-- Photo -->
                  <div class="lg:col-span-1">
                      <div class="relative">
                          <div class="aspect-square rounded-2xl overflow-hidden shadow-xl shadow-slate-300/50">
                              <img 
                                :src="block.data.image_url || 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=400'" 
                                :alt="block.data.name" 
                                class="w-full h-full object-cover"
                              />
                          </div>
                          <!-- Decorative Badge -->
                          <div class="absolute -bottom-4 -right-4 bg-amber-500 text-white px-4 py-2 rounded-xl font-semibold shadow-lg text-sm">
                              Narasumber
                          </div>
                      </div>
                  </div>
                  
                  <!-- Content -->
                  <div class="lg:col-span-2">
                      <p class="text-amber-600 font-semibold text-sm uppercase tracking-wider mb-2">Belajar dari Ahlinya</p>
                      <h3 class="text-3xl lg:text-4xl font-bold text-slate-800 mb-2">{{ block.data.name || 'Nama Narasumber' }}</h3>
                      <p class="text-blue-900 font-semibold text-lg mb-6">{{ block.data.title || 'Jabatan / Keahlian' }}</p>
                      
                      <p class="text-slate-600 leading-relaxed text-lg mb-8">
                          {{ block.data.bio || 'Deskripsi singkat tentang narasumber dan pengalamannya di bidang kepabeanan.' }}
                      </p>
                      
                      <!-- Stats/Credentials -->
                      <div v-if="block.data.stats && block.data.stats.length" class="flex flex-wrap gap-4">
                          <div v-for="(stat, idx) in block.data.stats" :key="idx" class="flex items-center gap-3 bg-white rounded-xl px-5 py-3 shadow-sm border border-slate-100">
                              <span class="text-2xl">{{ stat.icon }}</span>
                              <span class="font-medium text-slate-700">{{ stat.label }}</span>
                          </div>
                      </div>
                  </div>
              </div>
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
