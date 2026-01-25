<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  block: any
  styles?: any
  campaignId?: string
}>()

const form = ref({
  name: '',
  email: '',
  phone: ''
})

const isSubmitting = ref(false)
const isSubmitted = ref(false)
const errorMessage = ref('')
const successMessage = ref('')

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

const submitForm = async () => {
  if (!props.campaignId) {
    errorMessage.value = 'Campaign ID tidak ditemukan.'
    return
  }
  if (!form.value.name || !form.value.email || !form.value.phone) {
    errorMessage.value = 'Mohon lengkapi nama, email, dan nomor WhatsApp'
    return
  }

  isSubmitting.value = true
  errorMessage.value = ''
  successMessage.value = ''

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
        full_name: form.value.name,
        payment_method: '', // Free webinar usually doesn't need payment method or handled by backend if price 0
      }
    })

    if (response.is_free || response.message) {
      isSubmitted.value = true
      successMessage.value = response.message || 'Pendaftaran berhasil! Cek WhatsApp Anda untuk link webinar.'
      form.value = { name: '', email: '', phone: '' }
    } else if (response.payment_url) {
      window.location.href = response.payment_url
    }
  } catch (error: any) {
    console.error(error)
    errorMessage.value = error.data?.error || 'Terjadi kesalahan. Silakan coba lagi.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <section id="cta_healing" class="relative py-16 md:py-24 bg-gradient-to-b from-[#e8e4dc] to-[#d4cfc4]">
    <!-- Decorative Elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-1/4 left-0 w-64 h-64 bg-amber-200/20 rounded-full blur-3xl"></div>
      <div class="absolute bottom-1/4 right-0 w-80 h-80 bg-sky-200/20 rounded-full blur-3xl"></div>
    </div>

    <div class="relative max-w-xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-10">
        <div class="inline-block text-5xl mb-4">🌿</div>
        <h2 class="text-3xl md:text-4xl font-serif text-[#4a4540] mb-4">
          {{ block.data?.headline || 'Siap Memulai Perjalanan Pemulihan?' }}
        </h2>
        <p class="text-[#6b635a]">
          {{ block.data?.subheadline || 'Daftar sekarang untuk mengamankan tempatmu di webinar gratis ini.' }}
        </p>
      </div>

      <!-- Form Card -->
      <div class="bg-white rounded-3xl shadow-2xl p-8 md:p-10">
        <!-- Success State -->
        <div v-if="isSubmitted" class="text-center py-8">
          <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <svg class="w-10 h-10 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 class="text-2xl font-serif text-[#4a4540] mb-2">Pendaftaran Berhasil! 🎉</h3>
          <p class="text-[#6b635a] mb-6">{{ successMessage || 'Cek email kamu untuk detil webinar.' }}</p>
          <p class="text-sm text-amber-700 font-medium">Sampai jumpa di webinar P.U.L.I.H!</p>
        </div>

        <!-- Form -->
        <form v-else @submit.prevent="submitForm" class="space-y-5">
          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-[#4a4540] mb-2">Nama Lengkap *</label>
            <input 
              v-model="form.name"
              type="text" 
              required
              placeholder="Masukkan nama lengkap"
              class="w-full px-4 py-3 rounded-xl border border-[#d4cfc4] bg-[#faf8f5] focus:outline-none focus:ring-2 focus:ring-amber-400 focus:border-transparent transition-all"
            />
          </div>

          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-[#4a4540] mb-2">Email *</label>
            <input 
              v-model="form.email"
              type="email" 
              required
              placeholder="email@contoh.com"
              class="w-full px-4 py-3 rounded-xl border border-[#d4cfc4] bg-[#faf8f5] focus:outline-none focus:ring-2 focus:ring-amber-400 focus:border-transparent transition-all"
            />
          </div>

          <!-- Phone -->
          <div>
            <label class="block text-sm font-medium text-[#4a4540] mb-2">No. WhatsApp</label>
            <input 
              v-model="form.phone"
              type="tel"
              placeholder="08xxxxxxxxxx"
              class="w-full px-4 py-3 rounded-xl border border-[#d4cfc4] bg-[#faf8f5] focus:outline-none focus:ring-2 focus:ring-amber-400 focus:border-transparent transition-all"
            />
          </div>

          <!-- Error Message -->
          <p v-if="errorMessage" class="text-red-600 text-sm">{{ errorMessage }}</p>

          <!-- Submit Button -->
          <button 
            type="submit"
            :disabled="isSubmitting"
            class="w-full bg-gradient-to-r from-[#4a4540] to-[#6b635a] text-white font-semibold py-4 px-6 rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            :style="block.data?.button_color ? { background: block.data.button_color } : {}"
          >
            <span v-if="isSubmitting" class="flex items-center gap-2">
              <svg class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Mendaftar...
            </span>
            <span v-else>
              {{ block.data?.button_text || 'DAFTAR GRATIS SEKARANG' }}
            </span>
          </button>

          <!-- Privacy Note -->
          <p class="text-xs text-center text-[#8a8279]">
            {{ block.data?.privacy_note || 'Data kamu aman dan tidak akan dibagikan ke pihak lain.' }}
          </p>
        </form>
      </div>

      <!-- Additional Trust Elements -->
      <div class="mt-8 text-center">
        <div class="flex items-center justify-center gap-4 text-[#6b635a] text-sm">
          <span class="flex items-center gap-1.5">
            <span>✓</span> Gratis
          </span>
          <span class="flex items-center gap-1.5">
            <span>✓</span> Online
          </span>
          <span class="flex items-center gap-1.5">
            <span>✓</span> Kuota Terbatas
          </span>
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
