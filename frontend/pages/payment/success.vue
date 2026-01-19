<template>
  <div class="min-h-screen bg-gradient-to-b from-neutral-50 to-neutral-100 flex items-center justify-center p-4">
    <div class="bg-white rounded-2xl shadow-xl max-w-lg w-full p-8 text-center">
      <!-- Loading -->
      <div v-if="pending" class="py-12">
        <div class="animate-spin rounded-full h-16 w-16 border-4 border-primary-500 border-t-transparent mx-auto mb-4"></div>
        <p class="text-neutral-600">Memuat status pembayaran...</p>
      </div>

      <!-- Success -->
      <div v-else-if="transaction?.status === 'success' || transaction?.status === 'settlement'" class="space-y-6">
        <div class="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center mx-auto">
          <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        
        <div>
          <h1 class="text-2xl font-bold text-neutral-900 mb-2">Pembayaran Berhasil! 🎉</h1>
          <p class="text-neutral-600">Selamat! Anda berhasil terdaftar di kursus.</p>
        </div>

        <div class="bg-neutral-50 rounded-xl p-4 text-left space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Order ID</span>
            <span class="font-medium">{{ orderId }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Kursus</span>
            <span class="font-medium">{{ transaction?.course_name }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Total</span>
            <span class="font-medium text-primary-600">Rp {{ formatNumber(transaction?.amount || 0) }}</span>
          </div>
        </div>

        <div class="bg-blue-50 rounded-xl p-4 text-left">
          <h3 class="font-semibold text-blue-900 mb-2">📧 Cek Email Anda</h3>
          <p class="text-sm text-blue-700">
            Kami telah mengirimkan instruksi login ke email Anda. 
            Gunakan email tersebut untuk masuk ke platform LMS.
          </p>
        </div>

        <div class="space-y-3">
          <NuxtLink 
            v-if="transaction?.course_slug"
            :to="`/courses/${transaction.course_slug}`"
            class="block w-full py-3 bg-primary-600 text-white font-semibold rounded-xl hover:bg-primary-700 transition-all"
          >
            Lihat Kursus
          </NuxtLink>
          <NuxtLink 
            to="/login"
            class="block w-full py-3 border border-neutral-300 text-neutral-700 font-medium rounded-xl hover:bg-neutral-50 transition-all"
          >
            Login ke Akun
          </NuxtLink>
        </div>
      </div>

      <!-- Pending -->
      <div v-else-if="transaction?.status === 'pending'" class="space-y-6">
        <div class="w-20 h-20 bg-yellow-100 rounded-full flex items-center justify-center mx-auto">
          <svg class="w-10 h-10 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        
        <div>
          <h1 class="text-2xl font-bold text-neutral-900 mb-2">Menunggu Pembayaran</h1>
          <p class="text-neutral-600">Silakan selesaikan pembayaran Anda.</p>
        </div>

        <div class="bg-neutral-50 rounded-xl p-4 text-left space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Order ID</span>
            <span class="font-medium">{{ orderId }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Kursus</span>
            <span class="font-medium">{{ transaction?.course_name }}</span>
          </div>
          <div class="flex justify-between text-sm">
            <span class="text-neutral-500">Total</span>
            <span class="font-medium text-primary-600">Rp {{ formatNumber(transaction?.amount || 0) }}</span>
          </div>
        </div>

        <button 
          @click="refreshStatus"
          class="w-full py-3 bg-primary-600 text-white font-semibold rounded-xl hover:bg-primary-700 transition-all"
        >
          Refresh Status
        </button>
      </div>

      <!-- Failed / Not Found -->
      <div v-else class="space-y-6">
        <div class="w-20 h-20 bg-red-100 rounded-full flex items-center justify-center mx-auto">
          <svg class="w-10 h-10 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </div>
        
        <div>
          <h1 class="text-2xl font-bold text-neutral-900 mb-2">Pembayaran Gagal</h1>
          <p class="text-neutral-600">{{ transaction?.status || 'Transaksi tidak ditemukan' }}</p>
        </div>

        <NuxtLink 
          to="/"
          class="block w-full py-3 bg-primary-600 text-white font-semibold rounded-xl hover:bg-primary-700 transition-all"
        >
          Kembali ke Home
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: false })

const route = useRoute()
const config = useRuntimeConfig()

const orderId = computed(() => route.query.order_id as string)
const apiBase = import.meta.server 
  ? (config.apiInternal as string) || 'http://api:8080'
  : (config.public.apiBase as string) || 'http://localhost:8080'

interface Transaction {
  order_id: string
  status: string
  amount: number
  course_name: string
  course_slug: string
  created_at: string
}

const { data: transaction, pending, refresh } = await useFetch<Transaction>(
  () => `/api/transaction-status/${orderId.value}`,
  {
    baseURL: apiBase,
    watch: [orderId],
    key: `transaction-${orderId.value}`
  }
)

const formatNumber = (num: number) => new Intl.NumberFormat('id-ID').format(num)

const refreshStatus = () => {
  refresh()
}

// Auto-refresh for pending status
let refreshInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  if (transaction.value?.status === 'pending') {
    refreshInterval = setInterval(() => {
      refresh()
    }, 5000)
  }
})

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval)
})

// Watch for status changes to stop polling
watch(() => transaction.value?.status, (newStatus) => {
  if (newStatus !== 'pending' && refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
})
</script>
