<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  block: any
  styles?: any
  campaignId?: string
}>()

const form = ref({
  fullName: '',
  email: '',
  phone: ''
})

const submitting = ref(false)
const submitError = ref('')
const submitSuccess = ref('')

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

const handleRegister = async () => {
  submitError.value = ''
  submitSuccess.value = ''

  if (!props.campaignId) {
    submitError.value = 'Campaign ID tidak ditemukan.'
    return
  }
  if (!form.value.email || !form.value.phone) {
    submitError.value = 'Email dan Nomor WhatsApp wajib diisi.'
    return
  }

  submitting.value = true
  try {
    const response = await $fetch<{
      payment_url?: string
      is_free?: boolean
      message?: string
    }>('/api/campaign-checkout', {
      method: 'POST',
      baseURL: apiBase,
      body: {
        campaign_id: props.campaignId,
        email: form.value.email,
        phone: form.value.phone,
        full_name: form.value.fullName,
        payment_method: '',
      }
    })

    if (response.is_free || response.message) {
      submitSuccess.value = response.message || 'Pendaftaran berhasil! Cek WhatsApp Anda.'
      form.value = { fullName: '', email: '', phone: '' }
    } else if (response.payment_url) {
      window.location.href = response.payment_url
    }
  } catch (err: any) {
    console.error(err)
    submitError.value = err.data?.error || 'Gagal memproses pendaftaran.'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <section id="register" class="py-24 px-4 text-center bg-[#050505] font-outfit">
      <div class="max-w-2xl mx-auto bg-neutral-900/50 backdrop-blur-xl border border-white/10 p-10 rounded-[40px] relative overflow-hidden">
          <!-- Glows -->
          <div class="absolute top-0 left-0 w-full h-full bg-gradient-to-b from-purple-500/5 to-transparent pointer-events-none"></div>

          <h2 class="text-3xl sm:text-4xl font-black mb-4 relative z-10 text-white">
            {{ block.data.headline || 'Siap Jadi Master Kepabeanan? 🚀' }}
          </h2>
          <p class="text-neutral-400 mb-8 relative z-10">
            {{ block.data.subheadline || 'Slot terbatas! Amankan kursimu sekarang sebelum penuh.' }}
          </p>

          <form @submit.prevent="handleRegister" class="relative z-10 space-y-4 max-w-sm mx-auto">
              <input v-model="form.fullName" type="text" placeholder="Nama Lengkap (Optional)" class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
              <input v-model="form.email" type="email" placeholder="Email Address *" required class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
              <input v-model="form.phone" type="tel" placeholder="WhatsApp Number *" required class="w-full bg-black/40 border border-white/10 rounded-2xl px-6 py-4 text-white placeholder-neutral-500 focus:border-[#d946ef] focus:ring-1 focus:ring-[#d946ef] outline-none transition-all" />
              
              <div v-if="submitError" class="text-red-500 text-sm font-bold">{{ submitError }}</div>
              <div v-if="submitSuccess" class="text-green-500 text-sm font-bold">{{ submitSuccess }}</div>

              <button :disabled="submitting" class="w-full bg-gradient-to-r from-[#ec4899] to-[#f97316] text-white font-bold text-lg py-4 rounded-2xl shadow-[0_10px_30px_rgba(236,72,153,0.3)] hover:shadow-[0_10px_40px_rgba(236,72,153,0.5)] hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed">
                  <span v-if="submitting">Memproses...</span>
                  <span v-else>{{ block.data.button_text || 'Daftar Sekarang (Gratis)' }}</span>
              </button>
          </form>
          <p class="text-neutral-600 text-sm mt-6">No credit card required. Instant access.</p>
      </div>
  </section>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;500;600;700;800;900&display=swap');

.font-outfit {
    font-family: 'Outfit', sans-serif;
}
</style>
