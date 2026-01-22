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
  return dateStr ? new Date(dateStr).getTime() : new Date().getTime() + 86400000 // Default 24h
})

const countdown = ref({ days: '00', hours: '00', minutes: '00', seconds: '00' })
const displayDate = computed(() => {
  const dateStr = props.globalEndDate || props.block.data.end_date
  if (!dateStr) return '25 Januari 2025'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'long', year: 'numeric' }).format(date)
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

const scrollToForm = () => {
    const formEl = document.querySelector('#cta_clean') || document.querySelector('form')
    if (formEl) formEl.scrollIntoView({ behavior: 'smooth' })
}
</script>

<template>
  <section class="max-w-6xl mx-auto px-4 py-12 md:py-20 flex flex-col md:flex-row items-center gap-10 md:gap-16 font-inter">
    <!-- CONTENT LEFT -->
    <div class="flex-1 text-center md:text-left order-2 md:order-1">
        <!-- Badge -->
        <div class="inline-block bg-[#FDF2F8] text-[#DB2777] font-bold text-xs px-3 py-1 rounded-full mb-6 uppercase tracking-wider">
            {{ block.data.badge || 'WEBINAR GRATIS' }}
        </div>
        
        <!-- Headline -->
        <h1 class="text-3xl md:text-5xl font-bold text-neutral-900 leading-tight mb-4">
            {{ block.data.headline }}
        </h1>
        
        <!-- Subheadline -->
        <p class="text-neutral-500 mb-8 max-w-lg mx-auto md:mx-0">
             {{ block.data.subheadline }}
        </p>

        <!-- Event Details -->
        <div class="flex flex-wrap justify-center md:justify-start gap-4 mb-8 text-sm text-neutral-600 font-medium">
             <div class="flex items-center gap-2">
                 <span>📅</span> {{ displayDate }}
             </div>
              <div class="flex items-center gap-2">
                 <span>⏰</span> {{ displayTime }}
             </div>
              <div class="flex items-center gap-2">
                 <img src="/gmeets_logo.png" alt="Google Meet" class="w-5 h-5 object-contain" /> Google Meet
             </div>
        </div>

         <!-- Countdown Box -->
        <div class="bg-neutral-50 rounded-xl p-6 mb-8 max-w-md mx-auto md:mx-0 border border-neutral-100">
             <p class="text-xs font-semibold text-neutral-400 uppercase mb-4 tracking-wide text-center">Pendaftaran Ditutup Dalam:</p>
             <div class="flex justify-between items-center px-4">
                 <div class="text-center">
                     <span class="text-2xl md:text-3xl font-bold text-[#F472B6] block">{{ countdown.days }}</span>
                     <span class="text-[10px] text-neutral-400 uppercase font-medium">Hari</span>
                 </div>
                 <div class="text-2xl text-neutral-200">:</div>
                 <div class="text-center">
                     <span class="text-2xl md:text-3xl font-bold text-[#F472B6] block">{{ countdown.hours }}</span>
                     <span class="text-[10px] text-neutral-400 uppercase font-medium">Jam</span>
                 </div>
                 <div class="text-2xl text-neutral-200">:</div>
                 <div class="text-center">
                     <span class="text-2xl md:text-3xl font-bold text-[#F472B6] block">{{ countdown.minutes }}</span>
                     <span class="text-[10px] text-neutral-400 uppercase font-medium">Menit</span>
                 </div>
                 <div class="text-2xl text-neutral-200">:</div>
                 <div class="text-center">
                     <span class="text-2xl md:text-3xl font-bold text-[#F472B6] block">{{ countdown.seconds }}</span>
                     <span class="text-[10px] text-neutral-400 uppercase font-medium">Detik</span>
                 </div>
             </div>
        </div>
        
        <!-- CTA -->
        <button @click="scrollToForm" 
                class="text-white font-bold py-3 px-8 rounded-full shadow-lg hover:shadow-xl hover:scale-105 transition-all w-full md:w-auto"
                :class="{ 'bg-gradient-to-r from-blue-900 to-blue-600': !block.data.button_color }"
                :style="{ background: block.data.button_color }">
            DAFTAR SEKARANG
        </button>
        <p class="text-xs text-blue-800 mt-3 font-medium flex items-center justify-center md:justify-start gap-1">
            <span>ℹ️</span> Hanya tersisa 50 kursi lagi!
        </p>

    </div>

    <!-- IMAGE RIGHT -->
    <div class="flex-1 order-1 md:order-2">
         <div class="relative rounded-2xl overflow-hidden shadow-2xl mx-auto max-w-md aspect-[4/3] bg-neutral-100">
             <!-- Placeholder or actual image if implemented later -->
             <img :src="block.data.image_url || 'https://images.unsplash.com/photo-1573164713988-8665fc963095?auto=format&fit=crop&q=80&w=800'" alt="Webinar Preview" class="object-cover w-full h-full hover:scale-105 transition-transform duration-700" />
         </div>
    </div>

  </section>
</template>

<style scoped>
/* Standard Google Font import if needed, but assuming Inter is available globally or via Nuxt config */
.font-inter {
    font-family: 'Inter', sans-serif;
}
</style>
