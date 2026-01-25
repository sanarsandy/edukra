<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()

const defaultProblems = [
  {
    emoji: '😤',
    title: 'Kepancing Emosi',
    description: 'Tiba-tiba marah atau sedih tanpa alasan jelas, padahal situasinya biasa saja'
  },
  {
    emoji: '🌀',
    title: 'Overthinking',
    description: 'Pikiran berputar-putar, susah tidur, cemas tentang hal yang belum terjadi'
  },
  {
    emoji: '🚪',
    title: 'Nutup Diri',
    description: 'Menghindari orang lain, susah terbuka, merasa tidak ada yang mengerti'
  },
  {
    emoji: '⏰',
    title: 'Prokrastinasi',
    description: 'Menunda-nunda hal penting, padahal tahu harus dilakukan'
  }
]

const problems = props.block.data?.problems || defaultProblems
</script>

<template>
  <section class="relative py-16 md:py-24 bg-gradient-to-b from-[#e8e4dc] to-[#f5f0e8]">
    <!-- Decorative -->
    <div class="absolute top-0 left-0 w-full h-24 bg-gradient-to-b from-[#d4cfc4] to-transparent"></div>
    
    <div class="max-w-5xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-12">
        <p class="text-sm font-medium text-amber-700 uppercase tracking-widest mb-3">
          {{ block.data?.eyebrow || 'Apakah kamu merasakan ini?' }}
        </p>
        <h2 class="text-3xl md:text-4xl font-serif text-[#4a4540] mb-4">
          {{ block.data?.headline || 'Bukan Masalah Hari Ini...' }}
        </h2>
        <p class="text-lg text-[#6b635a] max-w-2xl mx-auto">
          {{ block.data?.subheadline || 'Tapi luka lama dan pola lama yang kebawa sampai sekarang.' }}
        </p>
      </div>

      <!-- Problem Cards -->
      <div class="grid md:grid-cols-2 gap-6 mb-12">
        <div 
          v-for="(problem, index) in problems" 
          :key="index"
          class="bg-white/70 backdrop-blur-sm rounded-2xl p-6 border border-[#d4cfc4] shadow-sm hover:shadow-md transition-all group"
        >
          <div class="flex items-start gap-4">
            <div class="text-4xl flex-shrink-0 group-hover:scale-110 transition-transform">
              {{ problem.emoji }}
            </div>
            <div>
              <h3 class="font-semibold text-[#4a4540] text-lg mb-1">{{ problem.title }}</h3>
              <p class="text-[#6b635a] text-sm leading-relaxed">{{ problem.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Reflection Box -->
      <div class="bg-gradient-to-r from-amber-50 to-orange-50 rounded-2xl p-8 border border-amber-200/50 text-center">
        <p class="text-xl md:text-2xl text-[#4a4540] font-serif italic leading-relaxed">
          {{ block.data?.reflection || '"Ingin pulih, tapi nggak tahu caranya? Atau sudah coba banyak hal, tapi reaksinya tetap berulang?"' }}
        </p>
        <div class="mt-6 inline-flex items-center gap-2 text-amber-700 text-sm font-medium">
          <span class="w-12 h-px bg-amber-400"></span>
          <span>🍃</span>
          <span class="w-12 h-px bg-amber-400"></span>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.font-serif {
  font-family: 'Playfair Display', 'Georgia', serif;
}
</style>
