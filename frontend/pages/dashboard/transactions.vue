<template>
  <div>
    <!-- Page Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Riwayat Transaksi</h1>
      <p class="text-neutral-500 mt-1">Lihat semua riwayat pembayaran kursus Anda</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin w-8 h-8 border-4 border-admin-500 border-t-transparent rounded-full"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!transactions.length" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <svg class="w-16 h-16 text-neutral-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
      </svg>
      <h3 class="text-lg font-medium text-neutral-900 mb-2">Belum Ada Transaksi</h3>
      <p class="text-neutral-500 mb-4">Anda belum melakukan transaksi pembelian kursus</p>
      <NuxtLink to="/" class="inline-flex items-center gap-2 px-4 py-2 bg-admin-600 text-white rounded-lg hover:bg-admin-700 transition-colors">
        <span>Jelajahi Kursus</span>
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
        </svg>
      </NuxtLink>
    </div>

    <!-- Transactions List -->
    <div v-else class="space-y-4">
      <div 
        v-for="tx in transactions" 
        :key="tx.id"
        class="bg-white rounded-xl border border-neutral-200 p-4 hover:border-neutral-300 transition-colors"
      >
        <div class="flex flex-col md:flex-row md:items-center gap-4">
          <!-- Course Thumbnail -->
          <div v-if="tx.course" class="w-full md:w-24 h-16 bg-neutral-100 rounded-lg overflow-hidden flex-shrink-0">
            <img 
              v-if="tx.course.thumbnail_url"
              :src="getThumbnailUrl(tx.course.thumbnail_url)" 
              :alt="tx.course.title"
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
              </svg>
            </div>
          </div>

          <!-- Transaction Info -->
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-neutral-900 truncate">
              {{ tx.course?.title || 'Unknown Course' }}
            </h3>
            <p class="text-sm text-neutral-500">
              {{ formatDate(tx.created_at) }}
              <span v-if="tx.order_id" class="text-neutral-400"> • {{ tx.order_id }}</span>
            </p>
          </div>

          <!-- Amount & Status -->
          <div class="flex items-center gap-4 md:text-right">
            <div>
              <p class="font-semibold text-neutral-900">{{ formatPrice(tx.amount) }}</p>
              <p v-if="tx.payment_type" class="text-xs text-neutral-500">{{ tx.payment_type }}</p>
            </div>
            <span 
              :class="[getStatusLabel(tx.status).color, 'px-3 py-1 rounded-full text-xs font-medium']"
            >
              {{ getStatusLabel(tx.status).label }}
            </span>
          </div>
        </div>

        <!-- Action Button for Pending -->
        <div v-if="tx.status === 'pending' && tx.snap_token" class="mt-4 pt-4 border-t border-neutral-100">
          <button 
            @click="continuePayment(tx)"
            class="text-sm text-admin-600 hover:text-admin-700 font-medium flex items-center gap-1"
          >
            <span>Lanjutkan Pembayaran</span>
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Load More -->
      <div v-if="hasMore" class="text-center py-4">
        <button 
          @click="loadMore"
          :disabled="loadingMore"
          class="px-4 py-2 text-admin-600 hover:text-admin-700 font-medium disabled:opacity-50"
        >
          {{ loadingMore ? 'Memuat...' : 'Muat Lebih Banyak' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Riwayat Transaksi'
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase
const router = useRouter()

const { fetchMyTransactions, loadSnapScript, openSnapPayment, formatPrice, getStatusLabel } = usePayment()

const loading = ref(true)
const loadingMore = ref(false)
const transactions = ref<any[]>([])
const total = ref(0)
const offset = ref(0)
const limit = 10

const hasMore = computed(() => transactions.value.length < total.value)

const getThumbnailUrl = (url: string) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads')) return `${apiBase}${url}`
  return `${apiBase}/api/images/${url}`
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const fetchTransactions = async () => {
  try {
    const data = await fetchMyTransactions(limit, offset.value)
    transactions.value = data.transactions || []
    total.value = data.total || 0
  } catch (err) {
    console.error('Failed to fetch transactions:', err)
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  loadingMore.value = true
  offset.value += limit
  try {
    const data = await fetchMyTransactions(limit, offset.value)
    transactions.value.push(...(data.transactions || []))
  } catch (err) {
    console.error('Failed to load more:', err)
    offset.value -= limit
  } finally {
    loadingMore.value = false
  }
}

const continuePayment = async (tx: any) => {
  if (tx.snap_token) {
    // Get payment config for client key
    try {
      const token = useCookie('token')
      const config = await $fetch<any>(`${apiBase}/api/checkout/config`, {
        headers: token.value ? { 'Authorization': `Bearer ${token.value}` } : {}
      })
      
      if (config.client_key) {
        await loadSnapScript(config.client_key, config.is_production)
        openSnapPayment(tx.snap_token, {
          onSuccess: () => {
            if (tx.course_id) {
              router.push(`/dashboard/courses/${tx.course_id}`)
            } else {
              fetchTransactions()
            }
          },
          onPending: () => {
            fetchTransactions()
          },
          onError: () => {
            fetchTransactions()
          }
        })
      }
    } catch (err) {
      console.error('Failed to continue payment:', err)
    }
  } else if (tx.payment_url) {
    window.location.href = tx.payment_url
  }
}

onMounted(() => {
  fetchTransactions()
})
</script>
