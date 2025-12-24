<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Transaksi</h1>
        <p class="text-neutral-500 mt-1">Riwayat transaksi dan pembayaran</p>
      </div>
      <button @click="handleExport" class="px-4 py-2.5 text-sm font-medium text-neutral-700 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors flex items-center">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
        </svg>
        Export CSV
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-accent-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <span class="text-xs font-medium text-accent-600 bg-accent-50 px-2 py-1 rounded-full">+18%</span>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ formatCurrency(totalRevenue) }}</p>
        <p class="text-sm text-neutral-500">Total Pendapatan</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ transactions.length }}</p>
        <p class="text-sm text-neutral-500">Total Transaksi</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-accent-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-accent-600">{{ successCount }}</p>
        <p class="text-sm text-neutral-500">Berhasil</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-warm-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-warm-600">{{ pendingCount }}</p>
        <p class="text-sm text-neutral-500">Pending</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari transaksi..." 
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
        />
      </div>
      <select v-model="statusFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Status</option>
        <option value="success">Berhasil</option>
        <option value="pending">Pending</option>
        <option value="failed">Gagal</option>
      </select>
    </div>

    <!-- Transactions Table -->
    <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">ID</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Pengguna</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kursus / Gateway</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Jumlah</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Tanggal</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-for="tx in filteredTransactions" :key="tx.id" class="hover:bg-neutral-50 transition-colors">
              <td class="px-6 py-4">
                <span class="text-sm font-mono text-neutral-600">{{ tx.id.substring(0, 12) }}...</span>
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div :class="['w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-semibold', getAvatarColor(tx.id)]">
                    {{ tx.user?.full_name ? getInitials(tx.user.full_name) : 'U' }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-neutral-900">{{ tx.user?.full_name || 'Unknown' }}</p>
                    <p class="text-xs text-neutral-500">{{ tx.user?.email || tx.user_id?.substring(0, 8) || '-' }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <p class="text-sm text-neutral-900">{{ tx.course?.title || '-' }}</p>
                <p class="text-xs text-neutral-500">{{ tx.payment_gateway || 'Direct' }}</p>
              </td>
              <td class="px-6 py-4 text-sm font-medium text-neutral-900">{{ formatCurrency(tx.amount) }}</td>
              <td class="px-6 py-4">
                <span 
                  class="px-2.5 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-accent-100 text-accent-700': tx.status === 'success',
                    'bg-warm-100 text-warm-700': tx.status === 'pending',
                    'bg-red-100 text-red-700': tx.status === 'failed'
                  }"
                >
                  {{ tx.status === 'success' ? 'Berhasil' : tx.status === 'pending' ? 'Pending' : 'Gagal' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-500">{{ formatDate(tx.created_at) }}</td>
              <td class="px-6 py-4 text-right">
                <button @click="viewTransaction(tx)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

useHead({
  title: 'Transaksi - Admin'
})

// API composable
const { 
  loading, 
  error, 
  transactions, 
  transactionsTotal, 
  fetchTransactions, 
  updateTransactionStatus 
} = useDashboard()

const searchQuery = ref('')
const statusFilter = ref('')

// Load transactions on mount
onMounted(async () => {
  await loadTransactions()
})

const loadTransactions = async () => {
  await fetchTransactions({ limit: 100 })
}

// Computed stats
const totalRevenue = computed(() => {
  return transactions.value
    .filter(tx => tx.status === 'success')
    .reduce((sum, tx) => sum + tx.amount, 0)
})

const successCount = computed(() => transactions.value.filter(tx => tx.status === 'success').length)
const pendingCount = computed(() => transactions.value.filter(tx => tx.status === 'pending').length)

const filteredTransactions = computed(() => {
  return transactions.value.filter(tx => {
    const matchesSearch = (tx.id || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          (tx.payment_gateway || '').toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !statusFilter.value || tx.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

// Helpers
const formatCurrency = (amount: number) => {
  if (amount >= 1000000) {
    return `Rp ${(amount / 1000000).toFixed(1)}M`
  } else if (amount >= 1000) {
    return `Rp ${(amount / 1000).toFixed(0)}K`
  }
  return `Rp ${amount.toLocaleString()}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    success: 'Berhasil',
    pending: 'Pending',
    failed: 'Gagal',
    refunded: 'Refund'
  }
  return labels[status] || status
}

const getAvatarColor = (id: string) => {
  const colors = ['bg-primary-500', 'bg-accent-500', 'bg-warm-500', 'bg-rose-500', 'bg-cyan-500', 'bg-indigo-500']
  const hash = (id || '').split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[hash % colors.length]
}

const getInitials = (name: string) => {
  if (!name) return '?'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const viewTransaction = (tx: any) => {
  // TODO: Open detail modal
  console.log('View transaction:', tx)
}

// Export composable
const { exportTransactions } = useExport()

const handleExport = () => {
  if (filteredTransactions.value.length === 0) {
    return
  }
  exportTransactions(filteredTransactions.value)
}
</script>
