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
    whatsapp: '',
    businessType: ''
})
const loading = ref(false)

const submitForm = async () => {
    loading.value = true
    // Simulate API call or emit event
    await new Promise(r => setTimeout(r, 1500))
    alert('Pendaftaran Berhasil! (Demo)')
    loading.value = false
}
</script>

<template>
  <section id="cta_clean" class="py-20 bg-white" :style="{ fontFamily: block.data.font_family || 'Inter, sans-serif' }">
      <div class="max-w-2xl mx-auto px-4 text-center">
          <h2 class="text-3xl font-bold text-neutral-900 mb-2">{{ block.data.headline }}</h2>
          <p class="text-neutral-500 mb-10 text-sm">{{ block.data.subheadline }}</p>

          <form @submit.prevent="submitForm" class="bg-[#F9FAFB] p-8 rounded-2xl border border-neutral-100 text-left space-y-5">
               <div>
                   <label class="block text-xs font-semibold text-neutral-600 mb-2">Nama Lengkap</label>
                   <input v-model="form.name" type="text" placeholder="Masukkan nama lengkap Anda" class="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-neutral-400 focus:outline-none text-sm" required />
               </div>

               <div>
                   <label class="block text-xs font-semibold text-neutral-600 mb-2">Email</label>
                   <input v-model="form.email" type="email" placeholder="Masukkan email aktif Anda" class="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-neutral-400 focus:outline-none text-sm" required />
               </div>
               
               <div>
                   <label class="block text-xs font-semibold text-neutral-600 mb-2">Nomor WhatsApp</label>
                   <input v-model="form.whatsapp" type="tel" placeholder="Contoh: 081234567890" class="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-neutral-400 focus:outline-none text-sm" required />
               </div>

                <div>
                   <label class="block text-xs font-semibold text-neutral-600 mb-2">Jenis Usaha</label>
                   <div class="relative">
                       <select v-model="form.businessType" class="w-full px-4 py-3 rounded-lg border border-neutral-200 focus:border-neutral-400 focus:outline-none text-sm appearance-none bg-white">
                          <option value="" disabled selected>Pilih jenis usaha Anda</option>
                          <option value="fashion">Fashion & Busana</option>
                          <option value="fnb">Kuliner (F&B)</option>
                          <option value="beauty">Kecantikan & Skincare</option>
                          <option value="elektronik">Gadget & Elektronik</option>
                          <option value="other">Lainnya</option>
                       </select>
                       <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-neutral-400">▼</div>
                   </div>
               </div>

               <button type="submit" :disabled="loading" 
                       class="w-full text-white font-bold py-3.5 rounded-lg shadow-md hover:shadow-lg transition-all mt-4 disabled:opacity-70 disabled:cursor-not-allowed"
                       :class="{ 'bg-gradient-to-r from-blue-900 to-blue-600': !block.data.button_color }"
                       :style="{ background: block.data.button_color }">
                   {{ loading ? 'MEMPROSES...' : 'DAFTAR SEKARANG' }}
               </button>

          </form>
      </div>
  </section>
</template>

<style scoped>
.font-inter {
    font-family: 'Inter', sans-serif;
}
</style>
