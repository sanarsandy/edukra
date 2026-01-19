<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()
</script>

<template>
  <section class="py-20 px-4 bg-[#050505] font-outfit relative overflow-hidden">
      <!-- Background Glows -->
      <div class="absolute top-1/2 left-1/4 w-[500px] h-[500px] bg-purple-600/10 rounded-full blur-[100px] -translate-y-1/2 pointer-events-none"></div>
      <div class="absolute bottom-0 right-1/4 w-[400px] h-[400px] bg-pink-600/10 rounded-full blur-[100px] pointer-events-none"></div>

      <div class="max-w-4xl mx-auto relative z-10 text-center">
          <div class="mb-12">
            <h2 class="text-4xl md:text-5xl font-black text-white mb-4">
                MEET THE <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#e879f9] to-[#f97316]">EXPERT</span>
            </h2>
            <p class="text-neutral-400">Belajar langsung dari ahlinya. No theory, full practical insight! 🧠</p>
          </div>

          <div class="bg-neutral-900/50 backdrop-blur-xl border border-white/10 rounded-[40px] p-8 md:p-12 relative group hover:border-[#d946ef]/30 transition-colors">
              <!-- Avatar -->
              <div class="w-40 h-40 mx-auto mb-8 relative">
                  <div class="absolute inset-0 bg-gradient-to-br from-[#e879f9] to-[#f97316] rounded-full blur-lg opacity-50 group-hover:opacity-100 transition-opacity"></div>
                  <img :src="block.data.image_url" :alt="block.data.name" class="w-full h-full object-cover rounded-full border-4 border-[#050505] relative z-10" />
                  <div class="absolute -bottom-2 -right-2 text-4xl animate-bounce z-20">👋</div>
              </div>

              <!-- Content -->
              <h3 class="text-3xl font-bold text-white mb-2">{{ block.data.name }}</h3>
              <p class="text-[#d946ef] font-bold tracking-wider uppercase text-sm mb-6">{{ block.data.title }}</p>
              
              <p class="text-neutral-300 leading-relaxed max-w-2xl mx-auto mb-10 text-lg">
                  "{{ block.data.bio }}"
              </p>

              <!-- Stats -->
              <div v-if="block.data.stats && block.data.stats.length" class="grid grid-cols-2 md:grid-cols-2 gap-4 max-w-lg mx-auto">
                  <div v-for="(stat, idx) in block.data.stats" :key="idx" class="bg-black/40 rounded-2xl p-4 border border-white/5 hover:border-white/20 transition-colors">
                      <div class="text-2xl mb-1">{{ stat.icon }}</div>
                      <div class="text-white font-bold">{{ stat.label }}</div>
                  </div>
              </div>
          </div>
      </div>
  </section>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800;900&display=swap');

.font-outfit {
    font-family: 'Outfit', sans-serif;
}
</style>
