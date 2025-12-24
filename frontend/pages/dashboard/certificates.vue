<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Sertifikat Saya</h1>
      <p class="text-neutral-500 mt-1">Koleksi sertifikat yang telah Anda peroleh</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-16">
      <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
    </div>

    <template v-else>
      <!-- Stats -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div class="bg-white rounded-xl p-5 border border-neutral-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-neutral-500 mb-1">Total Sertifikat</p>
              <p class="text-2xl font-bold text-neutral-900">{{ certificates.length }}</p>
            </div>
            <div class="w-11 h-11 bg-primary-100 rounded-xl flex items-center justify-center">
              <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Certificates Grid -->
      <div v-if="certificates.length > 0" class="grid md:grid-cols-2 gap-6">
        <div 
          v-for="cert in certificates" 
          :key="cert.id" 
          class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all"
        >
          <!-- Certificate Preview -->
          <div class="relative bg-gradient-to-br from-primary-600 to-primary-800 p-8 text-center text-white">
            <div class="absolute inset-0 opacity-10">
              <svg class="w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
                <pattern id="grid" width="10" height="10" patternUnits="userSpaceOnUse">
                  <path d="M 10 0 L 0 0 0 10" fill="none" stroke="white" stroke-width="0.5"/>
                </pattern>
                <rect width="100" height="100" fill="url(#grid)"/>
              </svg>
            </div>
            <div class="relative">
              <div class="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
                </svg>
              </div>
              <p class="text-xs uppercase tracking-widest opacity-80 mb-2">Sertifikat Penyelesaian</p>
              <h3 class="text-lg font-bold">{{ cert.course_name }}</h3>
            </div>
          </div>
          
          <!-- Certificate Info -->
          <div class="p-5">
            <div class="flex items-center justify-between mb-4">
              <div>
                <p class="text-sm text-neutral-500">Diterbitkan</p>
                <p class="font-medium text-neutral-900">{{ formatDate(cert.issued_at) }}</p>
              </div>
              <div class="text-right">
                <p class="text-sm text-neutral-500">ID Sertifikat</p>
                <p class="font-medium text-neutral-900 font-mono text-sm">{{ cert.certificate_number }}</p>
              </div>
            </div>
            
            <div class="flex gap-3">
              <button 
                @click="viewCert(cert)"
                class="flex-1 py-2.5 text-sm font-medium text-primary-600 bg-primary-50 rounded-lg hover:bg-primary-100 transition-colors flex items-center justify-center"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                </svg>
                Lihat
              </button>
              <button 
                @click="downloadCert(cert)"
                class="flex-1 py-2.5 text-sm font-medium text-white bg-primary-600 rounded-lg hover:bg-primary-700 transition-colors flex items-center justify-center"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
                </svg>
                Unduh PDF
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16">
        <div class="w-20 h-20 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-10 h-10 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z"/>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada sertifikat</h3>
        <p class="text-neutral-500 mb-6">Selesaikan kursus untuk mendapatkan sertifikat</p>
        <NuxtLink to="/dashboard/courses" class="btn-primary">Lihat Kursus Saya</NuxtLink>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
const { certificates, fetchCertificates, formatDate, loading } = useCertificates()
const config = useRuntimeConfig()

definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Sertifikat Saya - LearnHub'
})

onMounted(async () => {
  await fetchCertificates()
})

const viewCert = async (cert: any) => {
  try {
    const baseUrl = config.public.apiBase || ''
    // Get token from cookie
    const token = useCookie('token').value
    
    const response = await fetch(`${baseUrl}/api/certificates/${cert.id}/download`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })
    
    if (!response.ok) {
      throw new Error('Failed to fetch certificate')
    }
    
    const html = await response.text()
    
    // Open new window and write HTML
    const newWindow = window.open('', '_blank')
    if (newWindow) {
      newWindow.document.write(html)
      newWindow.document.close()
    }
  } catch (err) {
    console.error('Error viewing certificate:', err)
    alert('Gagal membuka sertifikat')
  }
}

const downloadCert = async (cert: any) => {
  // Same as view - user can print to PDF from browser
  await viewCert(cert)
}
</script>

