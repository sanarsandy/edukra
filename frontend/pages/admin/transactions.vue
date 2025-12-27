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
        <option value="cancelled">Dibatalkan</option>
      </select>
    </div>

    <!-- Bulk Actions Bar -->
    <Transition name="slide-down">
      <div v-if="selectedTransactions.length > 0" class="bg-admin-600 text-white px-6 py-3 rounded-xl mb-4 flex items-center justify-between">
        <span class="text-sm font-medium">{{ selectedTransactions.length }} transaksi dipilih</span>
        <div class="flex items-center gap-3">
          <button @click="openBulkDeleteModal" class="px-3 py-1.5 text-sm font-medium bg-red-500 hover:bg-red-600 rounded-lg transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            Hapus
          </button>
          <button @click="clearSelection" class="p-1.5 hover:bg-white/20 rounded-lg transition-colors">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
    </Transition>

    <!-- Transactions Table -->
    <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-4 py-4 text-left">
                <input 
                  type="checkbox" 
                  :checked="isAllSelected"
                  @change="toggleSelectAll"
                  class="w-4 h-4 text-admin-600 rounded border-neutral-300 focus:ring-admin-500"
                />
              </th>
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
            <tr v-for="tx in filteredTransactions" :key="tx.id" class="hover:bg-neutral-50 transition-colors" :class="selectedTransactions.includes(tx.id) ? 'bg-admin-50' : ''">
              <td class="px-4 py-4">
                <input 
                  type="checkbox" 
                  :checked="selectedTransactions.includes(tx.id)"
                  @change="toggleSelectTransaction(tx.id)"
                  class="w-4 h-4 text-admin-600 rounded border-neutral-300 focus:ring-admin-500"
                  :disabled="tx.status === 'success'"
                />
              </td>
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
                    'bg-red-100 text-red-700': tx.status === 'failed',
                    'bg-neutral-100 text-neutral-700': tx.status === 'cancelled'
                  }"
                >
                  {{ getStatusLabel(tx.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-500">{{ formatDate(tx.created_at) }}</td>
              <td class="px-6 py-4 text-right flex items-center justify-end gap-1">
                <button @click="viewTransaction(tx)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors" title="Lihat Detail">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
                <button 
                  v-if="tx.status !== 'success'" 
                  @click="confirmDelete(tx)" 
                  class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                  title="Hapus Transaksi"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Transaction Detail Modal -->
    <Teleport to="body">
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showDetailModal = false">
        <div class="bg-white rounded-xl shadow-xl max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto">
          <div class="p-6 border-b border-neutral-200 flex items-center justify-between">
            <h3 class="text-lg font-semibold text-neutral-900">Detail Transaksi</h3>
            <button @click="showDetailModal = false" class="text-neutral-400 hover:text-neutral-600">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div v-if="selectedTransaction" class="p-6 space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-xs text-neutral-500 mb-1">ID Transaksi</p>
                <p class="text-sm font-mono text-neutral-900 break-all">{{ selectedTransaction.id }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Order ID</p>
                <p class="text-sm font-mono text-neutral-900">{{ selectedTransaction.order_id || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Status</p>
                <span 
                  class="px-2.5 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-accent-100 text-accent-700': selectedTransaction.status === 'success',
                    'bg-warm-100 text-warm-700': selectedTransaction.status === 'pending',
                    'bg-red-100 text-red-700': selectedTransaction.status === 'failed',
                    'bg-neutral-100 text-neutral-700': selectedTransaction.status === 'cancelled'
                  }"
                >
                  {{ getStatusLabel(selectedTransaction.status) }}
                </span>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Jumlah</p>
                <p class="text-sm font-semibold text-neutral-900">Rp {{ Number(selectedTransaction.amount).toLocaleString('id-ID') }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Payment Gateway</p>
                <p class="text-sm text-neutral-900">{{ selectedTransaction.payment_gateway || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Payment Type</p>
                <p class="text-sm text-neutral-900">{{ selectedTransaction.payment_type || '-' }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-neutral-500 mb-1">Pengguna</p>
                <p class="text-sm text-neutral-900">{{ selectedTransaction.user?.full_name || 'Unknown' }}</p>
                <p class="text-xs text-neutral-500">{{ selectedTransaction.user?.email || '-' }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-neutral-500 mb-1">Kursus</p>
                <p class="text-sm text-neutral-900">{{ selectedTransaction.course?.title || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Tanggal Dibuat</p>
                <p class="text-sm text-neutral-900">{{ formatDateTime(selectedTransaction.created_at) }}</p>
              </div>
              <div>
                <p class="text-xs text-neutral-500 mb-1">Terakhir Update</p>
                <p class="text-sm text-neutral-900">{{ formatDateTime(selectedTransaction.updated_at) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showDeleteModal = false">
        <div class="bg-white rounded-xl shadow-xl max-w-md w-full mx-4">
          <div class="p-6">
            <div class="flex items-center justify-center w-12 h-12 mx-auto bg-red-100 rounded-full mb-4">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 text-center mb-2">Hapus Transaksi?</h3>
            <p class="text-sm text-neutral-600 text-center mb-6">
              Apakah Anda yakin ingin menghapus transaksi ini? Tindakan ini tidak dapat dibatalkan.
              <span v-if="transactionToDelete" class="block mt-2 font-mono text-xs text-neutral-500">
                ID: {{ transactionToDelete.id?.substring(0, 12) }}...
              </span>
            </p>
            <div class="flex gap-3">
              <button 
                @click="showDeleteModal = false" 
                class="flex-1 px-4 py-2.5 text-sm font-medium text-neutral-700 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors"
                :disabled="deleting"
              >
                Batal
              </button>
              <button 
                @click="deleteTransaction" 
                class="flex-1 px-4 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50"
                :disabled="deleting"
              >
                {{ deleting ? 'Menghapus...' : 'Hapus' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Bulk Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showBulkDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click.self="showBulkDeleteModal = false">
        <div class="bg-white rounded-xl shadow-xl max-w-md w-full mx-4">
          <div class="p-6">
            <div class="flex items-center justify-center w-12 h-12 mx-auto bg-red-100 rounded-full mb-4">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 text-center mb-2">Hapus {{ selectedTransactions.length }} Transaksi?</h3>
            <p class="text-sm text-neutral-600 text-center mb-6">
              Apakah Anda yakin ingin menghapus <strong>{{ selectedTransactions.length }} transaksi</strong> yang dipilih? Tindakan ini tidak dapat dibatalkan.
            </p>
            <div class="flex gap-3">
              <button 
                @click="showBulkDeleteModal = false" 
                class="flex-1 px-4 py-2.5 text-sm font-medium text-neutral-700 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors"
                :disabled="deleting"
              >
                Batal
              </button>
              <button 
                @click="bulkDeleteTransactions" 
                class="flex-1 px-4 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50"
                :disabled="deleting"
              >
                {{ deleting ? 'Menghapus...' : 'Hapus Semua' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Toast -->
    <Teleport to="body">
      <Transition name="slide-up">
        <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
          <div 
            class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3"
            :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'"
          >
            <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            <span class="text-sm font-medium">{{ toast.message }}</span>
          </div>
        </div>
      </Transition>
    </Teleport>
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

const { fetch: apiFetch } = useApi()

const searchQuery = ref('')
const statusFilter = ref('')
const showDetailModal = ref(false)
const selectedTransaction = ref<any>(null)
const showDeleteModal = ref(false)
const transactionToDelete = ref<any>(null)
const deleting = ref(false)
const showBulkDeleteModal = ref(false)

// Bulk selection
const selectedTransactions = ref<string[]>([])

const toast = ref({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

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
    .filter((tx: any) => tx.status === 'success')
    .reduce((sum: number, tx: any) => sum + tx.amount, 0)
})

const successCount = computed(() => transactions.value.filter((tx: any) => tx.status === 'success').length)
const pendingCount = computed(() => transactions.value.filter((tx: any) => tx.status === 'pending').length)

const filteredTransactions = computed(() => {
  return transactions.value.filter((tx: any) => {
    const matchesSearch = (tx.id || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          (tx.payment_gateway || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          (tx.order_id || '').toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !statusFilter.value || tx.status === statusFilter.value
    return matchesSearch && matchesStatus
  })
})

// Bulk selection - only select non-success transactions
const selectableTransactions = computed(() => {
  return filteredTransactions.value.filter((tx: any) => tx.status !== 'success')
})

const isAllSelected = computed(() => {
  return selectableTransactions.value.length > 0 && 
         selectedTransactions.value.length === selectableTransactions.value.length
})

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedTransactions.value = []
  } else {
    selectedTransactions.value = selectableTransactions.value.map((tx: any) => tx.id)
  }
}

const toggleSelectTransaction = (txId: string) => {
  const index = selectedTransactions.value.indexOf(txId)
  if (index === -1) {
    selectedTransactions.value.push(txId)
  } else {
    selectedTransactions.value.splice(index, 1)
  }
}

const clearSelection = () => {
  selectedTransactions.value = []
}

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

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    success: 'Berhasil',
    pending: 'Pending',
    failed: 'Gagal',
    cancelled: 'Dibatalkan',
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
  selectedTransaction.value = tx
  showDetailModal.value = true
}

// Single delete
const confirmDelete = (tx: any) => {
  transactionToDelete.value = tx
  showDeleteModal.value = true
}

const deleteTransaction = async () => {
  if (!transactionToDelete.value) return
  
  deleting.value = true
  try {
    await apiFetch(`/api/admin/transactions/${transactionToDelete.value.id}`, {
      method: 'DELETE'
    })
    
    await loadTransactions()
    showDeleteModal.value = false
    transactionToDelete.value = null
    showToast('Transaksi berhasil dihapus')
  } catch (error: any) {
    console.error('Delete failed:', error)
    showToast(error.data?.error || 'Gagal menghapus transaksi', 'error')
  } finally {
    deleting.value = false
  }
}

// Bulk delete
const openBulkDeleteModal = () => {
  showBulkDeleteModal.value = true
}

const bulkDeleteTransactions = async () => {
  deleting.value = true
  let successCount = 0
  let failCount = 0
  
  for (const txId of selectedTransactions.value) {
    try {
      await apiFetch(`/api/admin/transactions/${txId}`, {
        method: 'DELETE'
      })
      successCount++
    } catch (error) {
      failCount++
    }
  }
  
  deleting.value = false
  showBulkDeleteModal.value = false
  
  if (successCount > 0) {
    showToast(`${successCount} transaksi berhasil dihapus`)
    clearSelection()
    await loadTransactions()
  }
  
  if (failCount > 0) {
    showToast(`${failCount} transaksi gagal dihapus`, 'error')
  }
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

<style scoped>
.slide-down-enter-active, .slide-down-leave-active {
  transition: all 0.3s ease;
}
.slide-down-enter-from, .slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
