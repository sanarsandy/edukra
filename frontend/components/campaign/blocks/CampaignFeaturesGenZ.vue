<script setup lang="ts">
import { computed } from 'vue'
const props = defineProps<{
  block: any
  styles?: any
}>()

// Default items if none provided
const defaultItems = [
    { title: 'Kepabeanan 101', desc: 'Dari nol sampai paham! Konsep dasar yang wajib tau sebelum terjun ke dunia ekspor-impor.', icon: '📖', color: 'text-[#e879f9]', bg_gradient: 'from-purple-900/20 to-neutral-900', border_hover: 'hover:border-purple-500/30', glow_bg: 'bg-purple-600/10', icon_gradient: 'from-[#e879f9] to-[#a855f7]', shadow_color: 'shadow-purple-500/20', corner_icon: '📚' },
    { title: 'Anti Ribet Cukai', desc: 'Step-by-step prosedur cukai yang simple. Goodbye dokumen berantakan!', icon: '📋', color: 'text-[#f472b6]', bg_gradient: 'from-pink-900/20 to-neutral-900', border_hover: 'hover:border-pink-500/30', glow_bg: 'bg-pink-600/10', icon_gradient: 'from-[#f472b6] to-[#ec4899]', shadow_color: 'shadow-pink-500/20', corner_icon: '✨' },
    { title: 'Go International!', desc: 'Tips praktis buat kamu yang mau mulai atau scale up bisnis ke pasar global.', icon: '🌐', color: 'text-[#60a5fa]', bg_gradient: 'from-blue-900/20 to-neutral-900', border_hover: 'hover:border-blue-500/30', glow_bg: 'bg-blue-600/10', icon_gradient: 'from-[#60a5fa] to-[#3b82f6]', shadow_color: 'shadow-blue-500/20', corner_icon: '🌍' },
    { title: 'Karir yang Cuan', desc: 'Peluang karir di bidang kepabeanan & cukai yang lagi hot banget di 2026!', icon: '📈', color: 'text-[#fbbf24]', bg_gradient: 'from-orange-900/20 to-neutral-900', border_hover: 'hover:border-orange-500/30', glow_bg: 'bg-orange-600/10', icon_gradient: 'from-[#fbbf24] to-[#f59e0b]', shadow_color: 'shadow-orange-500/20', corner_icon: '💰' }
]

const items = computed(() => {
    return (props.block.data.items && props.block.data.items.length > 0) ? props.block.data.items : defaultItems
})
</script>

<template>
  <section class="max-w-7xl mx-auto px-4 py-20 bg-[#050505] font-outfit">
      <div class="text-center mb-16">
          <div class="text-6xl mb-4">🎯</div>
          <h2 class="text-4xl md:text-5xl font-black mb-4 text-white">
              {{ block.data.title || 'Yang Bakal Kamu' }} <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#e879f9] to-[#f97316]">DAPETIN</span>
          </h2>
          <p class="text-neutral-400 max-w-2xl mx-auto">Materi super praktis yang langsung bisa dipake. Cocok buat pemula sampai yang udah pro! 💪</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div v-for="(item, idx) in items" :key="idx" 
             :class="[
                'border border-white/5 p-10 rounded-[32px] relative overflow-hidden group transition-all',
                item.bg_gradient || 'bg-gradient-to-br from-neutral-900/50 to-neutral-900',
                item.border_hover || 'hover:border-white/20'
             ]">
             
            <!-- Glow Effect -->
            <div :class="['absolute top-0 right-0 w-[300px] h-[300px] rounded-full blur-[80px] -translate-y-1/2 translate-x-1/2 transition-colors', item.glow_bg]"></div>
            
            <!-- Icon -->
            <div :class="['w-16 h-16 rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-lg transition-transform group-hover:scale-110', 
                item.icon_gradient ? `bg-gradient-to-br ${item.icon_gradient}` : 'bg-neutral-800', 
                item.shadow_color]">
                {{ item.icon }}
            </div>
            
            <h3 class="text-2xl font-bold mb-3 relative z-10 text-white">{{ item.title }}</h3>
            <p class="text-neutral-400 text-sm leading-relaxed relative z-10">
                {{ item.desc }}
            </p>
            <div class="absolute top-4 right-4 text-2xl select-none">{{ item.corner_icon || '✨' }}</div>
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
