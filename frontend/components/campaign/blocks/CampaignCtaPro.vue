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
      submitSuccess.value = response.message || 'Pendaftaran berhasil! Cek WhatsApp Anda untuk link webinar.'
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
  <section id="register" class="py-20 px-4 bg-white font-inter">
      <div class="max-w-5xl mx-auto">
          <div class="grid lg:grid-cols-2 gap-12 items-center">
              <!-- Left Content -->
              <div>
                  <p class="text-amber-600 font-semibold text-sm uppercase tracking-wider mb-3">Daftar Sekarang</p>
                  <h2 class="text-3xl lg:text-4xl font-bold text-slate-800 mb-6">
                      {{ block.data.headline || 'Amankan Tempatmu Sekarang!' }}
                  </h2>
                  <p class="text-slate-600 text-lg leading-relaxed mb-8">
                      {{ block.data.subheadline || 'Kuota terbatas! Daftar gratis dan dapatkan akses penuh ke webinar beserta e-sertifikat.' }}
                  </p>
                  
                  <!-- Benefits List -->
                  <div class="space-y-4">
                      <div class="flex items-center gap-3">
                          <div class="w-6 h-6 bg-green-100 rounded-full flex items-center justify-center">
                              <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                              </svg>
                          </div>
                          <span class="text-slate-700">100% Gratis tanpa biaya apapun</span>
                      </div>
                      <div class="flex items-center gap-3">
                          <div class="w-6 h-6 bg-green-100 rounded-full flex items-center justify-center">
                              <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                              </svg>
                          </div>
                          <span class="text-slate-700">E-Sertifikat untuk semua peserta</span>
                      </div>
                      <div class="flex items-center gap-3">
                          <div class="w-6 h-6 bg-green-100 rounded-full flex items-center justify-center">
                              <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                              </svg>
                          </div>
                          <span class="text-slate-700">Akses rekaman setelah acara</span>
                      </div>
                      <div class="flex items-center gap-3">
                          <div class="w-6 h-6 bg-green-100 rounded-full flex items-center justify-center">
                              <svg class="w-4 h-4 text-green-600" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                              </svg>
                          </div>
                          <span class="text-slate-700">Sesi tanya jawab interaktif</span>
                      </div>
                  </div>
              </div>
              
              <!-- Right Form -->
              <div>
                  <div class="bg-slate-50 rounded-3xl p-8 border border-slate-200">
                      <h3 class="text-xl font-bold text-slate-800 mb-6 text-center">Form Pendaftaran</h3>
                      
                      <form @submit.prevent="handleRegister" class="space-y-4">
                          <div>
                              <label class="block text-sm font-medium text-slate-700 mb-2">Nama Lengkap</label>
                              <input 
                                v-model="form.fullName" 
                                type="text" 
                                placeholder="Masukkan nama lengkap"
                                class="w-full bg-white border border-slate-200 rounded-xl px-4 py-3.5 text-slate-800 placeholder-slate-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all"
                              />
                          </div>
                          <div>
                              <label class="block text-sm font-medium text-slate-700 mb-2">Email <span class="text-red-500">*</span></label>
                              <input 
                                v-model="form.email" 
                                type="email" 
                                placeholder="nama@email.com" 
                                required
                                class="w-full bg-white border border-slate-200 rounded-xl px-4 py-3.5 text-slate-800 placeholder-slate-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all"
                              />
                          </div>
                          <div>
                              <label class="block text-sm font-medium text-slate-700 mb-2">WhatsApp <span class="text-red-500">*</span></label>
                              <input 
                                v-model="form.phone" 
                                type="tel" 
                                placeholder="08xxxxxxxxxx" 
                                required
                                class="w-full bg-white border border-slate-200 rounded-xl px-4 py-3.5 text-slate-800 placeholder-slate-400 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition-all"
                              />
                          </div>
                          
                          <div v-if="submitError" class="bg-red-50 text-red-700 text-sm rounded-xl px-4 py-3 border border-red-100">{{ submitError }}</div>
                          <div v-if="submitSuccess" class="bg-green-50 text-green-700 text-sm rounded-xl px-4 py-3 border border-green-100">{{ submitSuccess }}</div>
                          
                          <button 
                            :disabled="submitting" 
                            class="w-full bg-blue-900 hover:bg-blue-800 text-white font-semibold py-4 rounded-xl shadow-lg shadow-blue-900/20 transition-all hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed mt-2"
                          >
                              <span v-if="submitting">Memproses...</span>
                              <span v-else>{{ block.data.button_text || 'Daftar Webinar Gratis' }}</span>
                          </button>
                      </form>
                      
                      <p class="text-center text-slate-500 text-sm mt-6">Dengan mendaftar, kamu setuju dengan syarat & ketentuan kami.</p>
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
