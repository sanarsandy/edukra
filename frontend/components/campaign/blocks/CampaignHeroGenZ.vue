<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted, watch } from 'vue'

const props = defineProps<{
  block: any
  styles?: any
  globalEndDate?: string | null
}>()

// Countdown Logic
const targetDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  return dateStr ? new Date(dateStr).getTime() : new Date('2026-01-27T19:00:00+07:00').getTime()
})

const countdown = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
const displayDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '27 Jan 2026'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }).format(date)
})
const displayTime = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '19:00 WIB'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { hour: '2-digit', minute: '2-digit', timeZoneName: 'short' }).format(date)
})

let timer: any = null

const updateCountdown = () => {
  const now = new Date().getTime()
  const distance = targetDate.value - now

  if (distance < 0) {
    countdown.value = { days: '00', hours: '00', minutes: '00', seconds: '00' }
    if (timer) clearInterval(timer)
    return
  }

  const d = Math.floor(distance / (1000 * 60 * 60 * 24))
  const h = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const m = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60))
  const s = Math.floor((distance % (1000 * 60)) / 1000)

  countdown.value = {
    days: d < 10 ? `0${d}` : `${d}`,
    hours: h < 10 ? `0${h}` : `${h}`,
    minutes: m < 10 ? `0${m}` : `${m}`,
    seconds: s < 10 ? `0${s}` : `${s}`
  }
}

watch(() => [props.block.data.end_date, props.globalEndDate], () => {
    updateCountdown()
})

onMounted(() => {
  updateCountdown()
  if (typeof window !== 'undefined') {
    timer = setInterval(updateCountdown, 1000)
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <header class="relative min-h-[80vh] flex flex-col justify-center items-center px-4 pt-10 text-center bg-[#050505] overflow-hidden font-outfit">
      
      <!-- Top Badge -->
      <div v-if="block.data.badge" class="inline-flex items-center gap-2 bg-neutral-900 border border-neutral-800 rounded-full pl-2 pr-4 py-1.5 mb-12 hover:scale-105 transition-transform cursor-pointer shadow-[0_0_20px_rgba(255,51,153,0.15)] relative z-20">
         <span class="text-xl">✨</span>
         <span class="text-[#84cc16] text-sm font-bold uppercase tracking-wider">{{ block.data.badge }}</span>
         <span class="text-xl">🔥</span>
      </div>
      
      <!-- Floating Emojis (3D Style) -->
      <div class="absolute top-[20%] left-[10%] text-5xl md:text-7xl animate-float-slow opacity-90 hidden md:block select-none pointer-events-none">🚀</div>
      <div class="absolute top-[15%] right-[15%] text-5xl md:text-7xl animate-float-medium opacity-90 hidden md:block select-none pointer-events-none">🌍</div>
      <div class="absolute bottom-[30%] left-[15%] text-5xl md:text-7xl animate-float-fast opacity-90 hidden md:block select-none pointer-events-none">💡</div>
      <div class="absolute bottom-[35%] right-[10%] text-5xl md:text-7xl animate-float-slow opacity-90 hidden md:block select-none pointer-events-none">💎</div>

      <!-- Main Title -->
      <h1 class="relative z-10 text-6xl sm:text-7xl md:text-9xl font-black leading-tight tracking-tight mb-8 text-white">
        <span class="text-transparent bg-clip-text bg-gradient-to-r from-[#e879f9] via-[#d946ef] to-[#f97316]">
          {{ block.data.headline || 'KEPABEANAN' }}
        </span>
      </h1>

      <!-- Subtitle -->
      <p class="relative z-10 text-lg sm:text-xl text-neutral-400 max-w-2xl mx-auto mb-12 font-medium">
        {{ block.data.subheadline }}
      </p>

      <!-- Info Pills -->
      <div class="flex flex-wrap justify-center gap-4 mb-20 relative z-10">
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <span class="text-[#a855f7]">📅</span>
          <span class="font-medium text-sm text-neutral-300">{{ displayDate }}</span>
        </div>
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <span class="text-[#ec4899]">⏰</span>
          <span class="font-medium text-sm text-neutral-300">{{ displayTime }}</span>
        </div>
        <div class="flex items-center gap-2 bg-neutral-900/80 border border-neutral-800 px-5 py-2.5 rounded-full backdrop-blur-md hover:border-[#d946ef]/50 transition-colors">
          <img src="/gmeets_logo.png" alt="Google Meet" class="w-5 h-5 object-contain" />
          <span class="font-medium text-sm text-neutral-300">Google Meet</span>
        </div>
      </div>

       <!-- Box Icon -->
       <div class="text-7xl mb-4 animate-bounce relative z-10 select-none">📦</div>

      <!-- Countdown Title -->
      <p class="text-xs font-bold text-neutral-500 uppercase tracking-[0.2em] mb-8 relative z-10">{{ block.data.promo_text || 'WEBINAR DIMULAI DALAM' }}</p>

      <!-- Countdown -->
      <div class="flex flex-wrap justify-center gap-4 relative z-10 mb-10">
        <!-- Days -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(168,85,247,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#e879f9] to-[#a855f7] bg-clip-text text-transparent">{{ countdown.days }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">HARI</div>
        </div>
        <!-- Hours -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(236,72,153,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#f472b6] to-[#ec4899] bg-clip-text text-transparent">{{ countdown.hours }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">JAM</div>
        </div>
        <!-- Minutes -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(59,130,246,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#60a5fa] to-[#3b82f6] bg-clip-text text-transparent">{{ countdown.minutes }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">MENIT</div>
        </div>
        <!-- Seconds -->
        <div class="text-center">
            <div class="w-24 h-24 bg-neutral-900/50 backdrop-blur-md rounded-3xl border border-white/5 flex items-center justify-center shadow-[0_0_30px_-5px_rgba(132,204,22,0.3)] hover:scale-105 transition-transform">
                <span class="text-5xl font-black bg-gradient-to-b from-[#a3e635] to-[#84cc16] bg-clip-text text-transparent">{{ countdown.seconds }}</span>
            </div>
            <div class="text-[10px] font-bold text-neutral-500 uppercase mt-3 tracking-wider">DETIK</div>
        </div>
      </div>
  </header>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800;900&display=swap');

.font-outfit {
    font-family: 'Outfit', sans-serif;
}

@keyframes float-slow {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}
@keyframes float-medium {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}
@keyframes float-fast {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.animate-float-slow { animation: float-slow 4s ease-in-out infinite; }
.animate-float-medium { animation: float-medium 3s ease-in-out infinite; }
.animate-float-fast { animation: float-fast 2.5s ease-in-out infinite; }
</style>
