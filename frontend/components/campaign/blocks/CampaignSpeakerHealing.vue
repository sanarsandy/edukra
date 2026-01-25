<script setup lang="ts">
const props = defineProps<{
  block: any
  styles?: any
}>()

const defaultCredentials = [
  'Certified Hypnotherapist (NGH & IACT)',
  'Coaching Level 1 ICF',
  'Ericksonian Corporate Hypnotherapist',
  'Master Neuro-Semantics'
]

const credentials = props.block.data?.credentials || defaultCredentials
</script>

<template>
  <section class="relative py-16 md:py-24 bg-gradient-to-b from-[#f5f0e8] to-[#e8e4dc]">
    <div class="max-w-5xl mx-auto px-4">
      <!-- Header -->
      <div class="text-center mb-12">
        <p class="text-sm font-medium text-amber-700 uppercase tracking-widest mb-3">
          {{ block.data?.eyebrow || 'Dipandu Oleh' }}
        </p>
        <h2 class="text-3xl md:text-4xl font-serif text-[#4a4540]">
          {{ block.data?.headline || 'Tentang Instruktur' }}
        </h2>
      </div>

      <!-- Speaker Card -->
      <div class="bg-white rounded-3xl shadow-xl overflow-hidden md:flex">
        <!-- Photo -->
        <div class="md:w-2/5 relative">
          <div class="aspect-square md:aspect-auto md:h-full bg-gradient-to-br from-[#d4cfc4] to-[#c4bcb0]">
            <img 
              :src="block.data?.image_url || 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?auto=format&fit=crop&q=80&w=600'" 
              :alt="block.data?.name || 'Speaker'"
              class="w-full h-full object-cover"
            />
          </div>
          <!-- Decorative Glow -->
          <div class="absolute -bottom-4 -right-4 w-32 h-32 bg-amber-300/30 rounded-full blur-2xl"></div>
        </div>

        <!-- Info -->
        <div class="md:w-3/5 p-8 md:p-10 flex flex-col justify-center">
          <h3 class="text-2xl md:text-3xl font-serif text-[#4a4540] mb-2">
            {{ block.data?.name || 'Fransiskus Indra Pratama' }}
          </h3>
          <p class="text-amber-700 font-medium mb-6">
            {{ block.data?.title || 'Coach, Hypnotherapist & Trainer' }}
          </p>
          
          <p class="text-[#6b635a] leading-relaxed mb-6">
            {{ block.data?.bio || 'Seorang Coach, Hypnotherapist dan Trainer dengan ratusan jam pembelajaran dan praktik. Pendekatan kerja mengintegrasikan hipnosis Ericksonian, coaching, NLP dan Neuro-Semantics untuk membantu proses kejelasan, perubahan, dan pengambilan keputusan yang bertanggung jawab.' }}
          </p>

          <!-- Credentials -->
          <div class="space-y-2">
            <p class="text-sm font-semibold text-[#4a4540] uppercase tracking-wide mb-3">Sertifikasi:</p>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="(cred, index) in credentials" 
                :key="index"
                class="inline-flex items-center gap-1.5 bg-amber-50 text-amber-800 text-xs font-medium px-3 py-1.5 rounded-full border border-amber-200"
              >
                <svg class="w-3.5 h-3.5 text-amber-600" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
                </svg>
                {{ cred }}
              </span>
            </div>
          </div>
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
