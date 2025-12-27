<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Header with Breadcrumb -->
    <header class="bg-white border-b border-neutral-200">
      <div class="max-w-4xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <button 
            @click="goBack" 
            class="flex items-center gap-2 text-admin-600 hover:text-admin-700 transition-colors"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
            </svg>
            <span>Kembali ke Kursus</span>
          </button>

          <!-- Security Badge -->
          <div class="flex items-center gap-1 text-xs text-neutral-500">
            <svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
            </svg>
            <span>Transaksi Aman</span>
          </div>
        </div>
      </div>
    </header>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-24">
      <div class="text-center">
        <div class="animate-spin w-12 h-12 border-4 border-admin-500 border-t-transparent rounded-full mx-auto mb-4"></div>
        <p class="text-neutral-500">Memuat data kursus...</p>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="max-w-4xl mx-auto px-4 py-12">
      <div class="bg-red-50 border border-red-200 rounded-xl p-6 text-center">
        <svg class="w-12 h-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <h2 class="text-lg font-semibold text-red-800 mb-2">{{ error }}</h2>
        <button @click="goBack" class="text-red-600 hover:text-red-700 underline">Kembali</button>
      </div>
    </div>

    <!-- Already Enrolled -->
    <div v-else-if="isEnrolled" class="max-w-4xl mx-auto px-4 py-12">
      <div class="bg-accent-50 border border-accent-200 rounded-xl p-8 text-center">
        <div class="w-16 h-16 bg-accent-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <h2 class="text-xl font-bold text-accent-800 mb-2">Anda Sudah Terdaftar!</h2>
        <p class="text-accent-600 mb-6">Anda sudah terdaftar di kursus ini. Silakan lanjutkan belajar.</p>
        <NuxtLink 
          :to="`/dashboard/courses/${courseId}`"
          class="inline-block px-6 py-3 bg-accent-600 text-white font-semibold rounded-lg hover:bg-accent-700 transition-colors"
        >
          Mulai Belajar
        </NuxtLink>
      </div>
    </div>

    <!-- Checkout Content -->
    <div v-else-if="course" class="max-w-4xl mx-auto px-4 py-8">
      <!-- Checkout Steps (Visual Progress) -->
      <div class="mb-8">
        <div class="flex items-center justify-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-admin-600 text-white rounded-full flex items-center justify-center text-sm font-bold">1</div>
            <span class="text-sm font-medium text-neutral-900">Detail Pesanan</span>
          </div>
          <div class="w-12 h-0.5 bg-neutral-200"></div>
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-neutral-200 text-neutral-500 rounded-full flex items-center justify-center text-sm font-bold">2</div>
            <span class="text-sm text-neutral-500">Pembayaran</span>
          </div>
          <div class="w-12 h-0.5 bg-neutral-200"></div>
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 bg-neutral-200 text-neutral-500 rounded-full flex items-center justify-center text-sm font-bold">3</div>
            <span class="text-sm text-neutral-500">Selesai</span>
          </div>
        </div>
      </div>

      <div class="grid md:grid-cols-5 gap-8">
        <!-- Course Info -->
        <div class="md:col-span-3">
          <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
            <!-- Course Thumbnail -->
            <div class="aspect-video bg-neutral-100">
              <img 
                v-if="course.thumbnail_url" 
                :src="getThumbnailUrl(course.thumbnail_url)" 
                :alt="course.title"
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
                <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
              </div>
            </div>

            <!-- Course Details -->
            <div class="p-6">
              <!-- Category Badge -->
              <div v-if="course.category" class="mb-3">
                <span class="px-2.5 py-1 bg-admin-100 text-admin-700 text-xs font-medium rounded-full">
                  {{ course.category.name }}
                </span>
              </div>

              <h1 class="text-2xl font-bold text-neutral-900 mb-2">{{ course.title }}</h1>
              <p class="text-neutral-600 mb-4 line-clamp-3">{{ course.description }}</p>
              
              <!-- Instructor Info -->
              <div v-if="course.instructor" class="flex items-center gap-3 mb-4 p-3 bg-neutral-50 rounded-lg">
                <div class="w-10 h-10 bg-admin-100 rounded-full flex items-center justify-center">
                  <span class="text-sm font-bold text-admin-600">
                    {{ course.instructor.full_name?.charAt(0) || 'I' }}
                  </span>
                </div>
                <div>
                  <p class="text-sm font-medium text-neutral-900">{{ course.instructor.full_name }}</p>
                  <p class="text-xs text-neutral-500">Instruktur</p>
                </div>
              </div>

              <div class="flex items-center gap-4 text-sm text-neutral-500">
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                  </svg>
                  {{ course.lessons_count || course.lessons?.length || 0 }} Materi
                </span>
                <span v-if="course.duration" class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                  {{ course.duration }}
                </span>
                <span class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
                  </svg>
                  {{ course.enrolled_count || 0 }} Siswa
                </span>
              </div>
            </div>
          </div>

          <!-- What You'll Get -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6 mt-6">
            <h3 class="font-semibold text-neutral-900 mb-4">Yang Anda Dapatkan</h3>
            <ul class="space-y-3">
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-accent-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span class="text-neutral-600">Akses seumur hidup ke semua materi</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-accent-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span class="text-neutral-600">Sertifikat penyelesaian</span>
              </li>
              <li class="flex items-start gap-3">
                <svg class="w-5 h-5 text-accent-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                </svg>
                <span class="text-neutral-600">Update materi gratis</span>
              </li>
            </ul>
          </div>
        </div>

        <!-- Order Summary -->
        <div class="md:col-span-2">
          <div class="bg-white rounded-xl border border-neutral-200 p-6 sticky top-4">
            <h2 class="text-lg font-semibold text-neutral-900 mb-4">Ringkasan Pesanan</h2>
            
            <div class="space-y-3 mb-6">
              <div class="flex justify-between text-sm">
                <span class="text-neutral-600">Harga Kursus</span>
                <span class="font-medium text-neutral-900">{{ formatPrice(course.price) }}</span>
              </div>
              <div class="border-t border-neutral-100 pt-3 flex justify-between">
                <span class="font-semibold text-neutral-900">Total</span>
                <span class="font-bold text-xl" :class="course.price === 0 ? 'text-accent-600' : 'text-admin-600'">
                  {{ course.price === 0 ? 'GRATIS' : formatPrice(course.price) }}
                </span>
              </div>
            </div>

            <!-- Free Course - Simplified UI -->
            <div v-if="course.price === 0">
              <!-- Free Badge -->
              <div class="bg-gradient-to-r from-accent-500 to-emerald-500 rounded-lg p-4 mb-4">
                <div class="flex items-center justify-center gap-2 text-white">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7"/>
                  </svg>
                  <span class="text-lg font-bold">KURSUS GRATIS!</span>
                </div>
              </div>

              <!-- Quick Benefits for Free Course -->
              <ul class="space-y-2 mb-4">
                <li class="flex items-center gap-2 text-sm">
                  <svg class="w-4 h-4 text-accent-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span class="text-neutral-600">Tanpa biaya pendaftaran</span>
                </li>
                <li class="flex items-center gap-2 text-sm">
                  <svg class="w-4 h-4 text-accent-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span class="text-neutral-600">Akses langsung ke semua materi</span>
                </li>
                <li class="flex items-center gap-2 text-sm">
                  <svg class="w-4 h-4 text-accent-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  <span class="text-neutral-600">Sertifikat penyelesaian</span>
                </li>
              </ul>

              <button 
                @click="handleEnroll"
                :disabled="processing"
                class="w-full py-4 bg-accent-600 text-white font-bold text-lg rounded-xl hover:bg-accent-700 transition-all transform hover:scale-[1.02] disabled:opacity-50 disabled:transform-none flex items-center justify-center gap-2 shadow-lg shadow-accent-500/30"
              >
                <svg v-if="processing" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ processing ? 'Mendaftarkan...' : 'Mulai Belajar Sekarang' }}</span>
                <svg v-if="!processing" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                </svg>
              </button>

              <p class="text-center mt-4 text-sm text-neutral-500">
                <span class="flex items-center justify-center gap-1">
                  <svg class="w-4 h-4 text-accent-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  {{ course.enrolled_count || 0 }}+ sudah bergabung
                </span>
              </p>
            </div>

            <!-- Paid Course -->
            <div v-else>
              <!-- Payment Method Selection for Duitku -->
              <div v-if="paymentConfig?.provider === 'duitku'" class="mb-4">
                <label class="block text-sm font-medium text-neutral-700 mb-2">Pilih Metode Pembayaran</label>
                
                <!-- Loading state -->
                <div v-if="loadingPaymentMethods" class="flex items-center justify-center py-4">
                  <svg class="w-6 h-6 animate-spin text-admin-600" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
                
                <!-- Payment Methods List -->
                <div v-else-if="paymentMethods.length > 0" class="space-y-2 max-h-64 overflow-y-auto">
                  <label
                    v-for="method in paymentMethods"
                    :key="method.paymentMethod"
                    class="relative flex items-center p-3 border rounded-lg cursor-pointer transition-all hover:border-admin-300"
                    :class="selectedPaymentMethod === method.paymentMethod ? 'border-admin-500 bg-admin-50' : 'border-neutral-200'"
                  >
                    <input
                      type="radio"
                      name="paymentMethod"
                      :value="method.paymentMethod"
                      v-model="selectedPaymentMethod"
                      class="hidden"
                    />
                    <img 
                      :src="method.paymentImage" 
                      :alt="method.paymentName"
                      class="w-12 h-8 object-contain mr-3"
                    />
                    <div class="flex-1">
                      <p class="text-sm font-medium text-neutral-900">{{ method.paymentName }}</p>
                      <p v-if="method.totalFee !== '0'" class="text-xs text-neutral-500">
                        Biaya: Rp {{ parseInt(method.totalFee).toLocaleString('id-ID') }}
                      </p>
                    </div>
                    <div 
                      v-if="selectedPaymentMethod === method.paymentMethod"
                      class="w-5 h-5 rounded-full bg-admin-600 flex items-center justify-center"
                    >
                      <svg class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                      </svg>
                    </div>
                  </label>
                </div>
                
                <!-- No payment methods available -->
                <div v-else class="text-center py-4 text-neutral-500 text-sm">
                  Tidak ada metode pembayaran tersedia
                </div>
              </div>

              <button 
                @click="handleCheckout"
                :disabled="processing || !paymentConfig?.enabled || (paymentConfig?.provider === 'duitku' && !selectedPaymentMethod)"
                class="w-full py-3 bg-admin-600 text-white font-semibold rounded-lg hover:bg-admin-700 transition-colors disabled:opacity-50 flex items-center justify-center gap-2"
              >
                <svg v-if="processing" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span>{{ processing ? 'Memproses...' : 'Bayar Sekarang' }}</span>
              </button>

              <!-- Payment not enabled warning -->
              <p v-if="!paymentConfig?.enabled" class="text-xs text-amber-600 text-center mt-3">
                Pembayaran belum diaktifkan oleh admin
              </p>

              <p v-else class="text-xs text-neutral-500 text-center mt-3">
                Pembayaran aman melalui {{ providerDisplayName }}
              </p>

              <!-- Trust Badges -->
              <div class="mt-4 pt-4 border-t border-neutral-100">
                <div class="flex items-center justify-center gap-4 text-xs text-neutral-500">
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                    </svg>
                    SSL Secure
                  </span>
                  <span class="flex items-center gap-1">
                    <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"/>
                    </svg>
                    Verified
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div 
          class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3"
          :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'"
        >
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const { createCheckout, loadSnapScript, openSnapPayment, formatPrice } = usePayment()

const courseId = route.params.id as string

const loading = ref(true)
const processing = ref(false)
const error = ref('')
const course = ref<any>(null)
const paymentConfig = ref<any>(null)
const isEnrolled = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

// Dynamic payment methods from Duitku
const paymentMethods = ref<any[]>([])
const selectedPaymentMethod = ref('')
const loadingPaymentMethods = ref(false)

// Computed display name for payment provider
const providerDisplayName = computed(() => {
  if (!paymentConfig.value) return 'Payment Gateway'
  switch (paymentConfig.value.provider) {
    case 'midtrans': return 'Midtrans'
    case 'duitku': return 'Duitku'
    case 'xendit': return 'Xendit'
    default: return 'Payment Gateway'
  }
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

// Go back to course detail page
const goBack = () => {
  // Check if we have history to go back
  if (window.history.length > 1) {
    router.back()
  } else {
    // Fallback to course detail page
    router.push(`/dashboard/courses/${courseId}`)
  }
}

const getThumbnailUrl = (url: string) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

// Check enrollment status
const checkEnrollment = async () => {
  try {
    const token = useCookie('token')
    if (!token.value) return
    
    const data = await $fetch<any>(`${apiBase}/api/enrollments/check/${courseId}`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    isEnrolled.value = data.enrolled
  } catch (err) {
    // Not enrolled, continue with checkout
  }
}

// Fetch course details
const fetchCourse = async () => {
  try {
    const token = useCookie('token')
    const data = await $fetch(`${apiBase}/api/courses/${courseId}`, {
      headers: token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
    })
    course.value = data
  } catch (err: any) {
    console.error('Failed to fetch course:', err)
    error.value = err.data?.error || 'Kursus tidak ditemukan'
  }
}

// Fetch payment config dan load snap script
const initPayment = async () => {
  try {
    const token = useCookie('token')
    const config = await $fetch<any>(`${apiBase}/api/checkout/config`, {
      headers: token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
    })
    paymentConfig.value = config

    // Only load Snap script for Midtrans
    if (config.enabled && config.provider === 'midtrans' && config.client_key) {
      await loadSnapScript(config.client_key, config.is_production)
    }
  } catch (err) {
    console.error('Failed to init payment:', err)
  }
}

// Fetch available payment methods for Duitku
const fetchPaymentMethods = async () => {
  if (paymentConfig.value?.provider !== 'duitku' || !course.value) return
  
  loadingPaymentMethods.value = true
  try {
    const token = useCookie('token')
    const amount = course.value.price || 10000
    const data = await $fetch<any>(`${apiBase}/api/checkout/payment-methods?amount=${amount}`, {
      headers: token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
    })
    
    if (data.payment_methods && data.payment_methods.length > 0) {
      paymentMethods.value = data.payment_methods
      // Auto-select first payment method
      selectedPaymentMethod.value = data.payment_methods[0].paymentMethod
    }
  } catch (err) {
    console.error('Failed to fetch payment methods:', err)
  } finally {
    loadingPaymentMethods.value = false
  }
}

// Handle free course enrollment
const handleEnroll = async () => {
  processing.value = true
  try {
    const result = await createCheckout(courseId)
    if (result.is_free) {
      showToast('Berhasil mendaftar kursus!')
      setTimeout(() => {
        router.push(`/dashboard/courses/${courseId}`)
      }, 1500)
    }
  } catch (err: any) {
    showToast(err.data?.error || err.message || 'Gagal mendaftar', 'error')
  } finally {
    processing.value = false
  }
}

// Handle paid course checkout
const handleCheckout = async () => {
  processing.value = true
  try {
    const returnUrl = window.location.origin + `/dashboard/courses/${courseId}`
    // Pass selected payment method for Duitku
    const paymentMethod = paymentConfig.value?.provider === 'duitku' ? selectedPaymentMethod.value : undefined
    const result = await createCheckout(courseId, returnUrl, paymentMethod)

    if (result.is_free) {
      showToast('Berhasil mendaftar kursus!')
      setTimeout(() => {
        router.push(`/dashboard/courses/${courseId}`)
      }, 1500)
      return
    }

    // For Midtrans - use Snap popup
    if (result.snap_token && paymentConfig.value?.provider === 'midtrans') {
      openSnapPayment(result.snap_token, {
        onSuccess: (res) => {
          showToast('Pembayaran berhasil!')
          router.push(`/dashboard/courses/${courseId}`)
        },
        onPending: (res) => {
          showToast('Menunggu pembayaran...', 'success')
          router.push('/dashboard/transactions')
        },
        onError: (res) => {
          showToast('Pembayaran gagal', 'error')
        },
        onClose: () => {
          processing.value = false
        }
      })
    } 
    // For Duitku and other redirect-based providers
    else if (result.payment_url) {
      window.location.href = result.payment_url
    } else {
      showToast('Konfigurasi pembayaran tidak valid', 'error')
      processing.value = false
    }
  } catch (err: any) {
    showToast(err.data?.error || err.message || 'Gagal memproses pembayaran', 'error')
    processing.value = false
  }
}

onMounted(async () => {
  await Promise.all([fetchCourse(), checkEnrollment(), initPayment()])
  // Fetch payment methods after course and payment config are ready
  await fetchPaymentMethods()
  loading.value = false
})
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}
</style>
