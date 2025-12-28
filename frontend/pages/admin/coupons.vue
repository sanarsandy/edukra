<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Kelola Kupon</h1>
        <p class="text-neutral-500 mt-1">Buat dan kelola kupon diskon untuk kursus.</p>
      </div>
      <button 
        @click="openModal('create')"
        class="inline-flex items-center px-4 py-2.5 bg-admin-600 hover:bg-admin-700 text-white font-medium rounded-lg transition-colors"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
        </svg>
        Buat Kupon
      </button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ coupons?.length || 0 }}</p>
        <p class="text-sm text-neutral-500">Total Kupon</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-success-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ activeCoupons }}</p>
        <p class="text-sm text-neutral-500">Aktif</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-warm-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ totalUsage }}</p>
        <p class="text-sm text-neutral-500">Total Penggunaan</p>
      </div>
      
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 bg-accent-100 rounded-lg flex items-center justify-center">
            <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-2xl font-bold text-neutral-900">{{ expiredCoupons }}</p>
        <p class="text-sm text-neutral-500">Kadaluarsa</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl border border-neutral-200 p-4 mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex-1 min-w-[200px]">
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Cari kode kupon..." 
              class="w-full pl-10 pr-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
            />
          </div>
        </div>
        <select 
          v-model="filterStatus"
          class="px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
        >
          <option value="all">Semua Status</option>
          <option value="active">Aktif</option>
          <option value="inactive">Tidak Aktif</option>
          <option value="expired">Kadaluarsa</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-admin-600"></div>
    </div>

    <!-- Coupons Table -->
    <div v-else class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kode</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Diskon</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Penggunaan</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Berlaku</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-for="coupon in filteredCoupons" :key="coupon.id" class="hover:bg-neutral-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center">
                  <code class="px-3 py-1.5 bg-primary-50 text-primary-700 font-mono font-bold rounded-lg text-sm">
                    {{ coupon.code }}
                  </code>
                </div>
              </td>
              <td class="px-6 py-4">
                <span v-if="coupon.discount_type === 'percentage'" class="font-semibold text-neutral-900">
                  {{ coupon.discount_value }}%
                </span>
                <span v-else class="font-semibold text-neutral-900">
                  {{ formatCurrency(coupon.discount_value) }}
                </span>
                <span v-if="coupon.max_discount" class="text-xs text-neutral-500 block">
                  Max: {{ formatCurrency(coupon.max_discount) }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="text-neutral-900">{{ coupon.usage_count }}</span>
                <span class="text-neutral-500">/ {{ coupon.usage_limit || '∞' }}</span>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-600">
                <div v-if="coupon.valid_until">
                  {{ formatDate(coupon.valid_from) }} - {{ formatDate(coupon.valid_until) }}
                </div>
                <div v-else class="text-success-600">Tidak ada batas</div>
              </td>
              <td class="px-6 py-4">
                <span 
                  :class="[
                    'px-2.5 py-1 rounded-full text-xs font-medium',
                    getCouponStatus(coupon) === 'active' ? 'bg-success-100 text-success-700' :
                    getCouponStatus(coupon) === 'expired' ? 'bg-neutral-100 text-neutral-600' :
                    'bg-red-100 text-red-700'
                  ]"
                >
                  {{ getCouponStatus(coupon) === 'active' ? 'Aktif' : getCouponStatus(coupon) === 'expired' ? 'Kadaluarsa' : 'Tidak Aktif' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    @click="openModal('edit', coupon)"
                    class="p-2 text-neutral-500 hover:text-admin-600 hover:bg-admin-50 rounded-lg transition-colors"
                    title="Edit"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                  </button>
                  <button 
                    @click="toggleCouponStatus(coupon)"
                    class="p-2 text-neutral-500 hover:text-warm-600 hover:bg-warm-50 rounded-lg transition-colors"
                    :title="coupon.is_active ? 'Nonaktifkan' : 'Aktifkan'"
                  >
                    <svg v-if="coupon.is_active" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"/>
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                  </button>
                  <button 
                    @click="confirmDelete(coupon)"
                    class="p-2 text-neutral-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title="Hapus"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredCoupons.length === 0">
              <td colspan="6" class="px-6 py-12 text-center text-neutral-500">
                <svg class="w-12 h-12 mx-auto text-neutral-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A2 2 0 013 12V7a4 4 0 014-4z"/>
                </svg>
                <p>Belum ada kupon. Buat kupon pertama Anda!</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="closeModal"></div>
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <div class="sticky top-0 bg-white border-b border-neutral-200 px-6 py-4 rounded-t-2xl">
          <h2 class="text-lg font-bold text-neutral-900">
            {{ modalMode === 'create' ? 'Buat Kupon Baru' : 'Edit Kupon' }}
          </h2>
          <button 
            @click="closeModal"
            class="absolute top-4 right-4 p-2 text-neutral-400 hover:text-neutral-600 rounded-lg hover:bg-neutral-100"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        
        <form @submit.prevent="submitForm" class="p-6 space-y-5">
          <!-- Code -->
          <div>
            <label class="block text-sm font-medium text-neutral-700 mb-2">Kode Kupon *</label>
            <input 
              v-model="form.code"
              type="text" 
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 uppercase"
              placeholder="DISKON20"
              required
            />
          </div>
          
          <!-- Discount Type & Value -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Tipe Diskon *</label>
              <select 
                v-model="form.discount_type"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                required
              >
                <option value="percentage">Persentase (%)</option>
                <option value="fixed">Nominal (Rp)</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Nilai Diskon *</label>
              <input 
                v-model.number="form.discount_value"
                type="number" 
                min="0"
                :max="form.discount_type === 'percentage' ? 100 : undefined"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                :placeholder="form.discount_type === 'percentage' ? '20' : '50000'"
                required
              />
            </div>
          </div>
          
          <!-- Max Discount (for percentage) -->
          <div v-if="form.discount_type === 'percentage'">
            <label class="block text-sm font-medium text-neutral-700 mb-2">Maksimal Potongan (Opsional)</label>
            <input 
              v-model.number="form.max_discount"
              type="number" 
              min="0"
              class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              placeholder="100000"
            />
            <p class="text-xs text-neutral-500 mt-1">Batas maksimal potongan harga dalam Rupiah</p>
          </div>
          
          <!-- Usage Limits -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Batas Penggunaan</label>
              <input 
                v-model.number="form.usage_limit"
                type="number" 
                min="1"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                placeholder="∞ Unlimited"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Per User</label>
              <input 
                v-model.number="form.per_user_limit"
                type="number" 
                min="1"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
                placeholder="1"
              />
            </div>
          </div>
          
          <!-- Validity Period -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Berlaku Dari</label>
              <input 
                v-model="form.valid_from"
                type="datetime-local"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Berlaku Sampai</label>
              <input 
                v-model="form.valid_until"
                type="datetime-local"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500"
              />
            </div>
          </div>
          
          <!-- Actions -->
          <div class="flex items-center justify-end gap-3 pt-4 border-t border-neutral-200">
            <button 
              type="button"
              @click="closeModal"
              class="px-5 py-2.5 text-neutral-600 hover:text-neutral-900 font-medium rounded-lg hover:bg-neutral-100 transition-colors"
            >
              Batal
            </button>
            <button 
              type="submit"
              :disabled="submitting"
              class="px-5 py-2.5 bg-admin-600 hover:bg-admin-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
            >
              {{ submitting ? 'Menyimpan...' : (modalMode === 'create' ? 'Buat Kupon' : 'Simpan Perubahan') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="showDeleteModal = false"></div>
      <div class="relative bg-white rounded-2xl shadow-xl w-full max-w-md p-6">
        <div class="text-center">
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
            </svg>
          </div>
          <h3 class="text-lg font-bold text-neutral-900 mb-2">Hapus Kupon?</h3>
          <p class="text-neutral-600 mb-6">
            Kupon <strong>{{ couponToDelete?.code }}</strong> akan dihapus permanen. Tindakan ini tidak dapat dibatalkan.
          </p>
          <div class="flex items-center justify-center gap-3">
            <button 
              @click="showDeleteModal = false"
              class="px-5 py-2.5 text-neutral-600 hover:text-neutral-900 font-medium rounded-lg hover:bg-neutral-100 transition-colors"
            >
              Batal
            </button>
            <button 
              @click="deleteCoupon"
              :disabled="deleting"
              class="px-5 py-2.5 bg-red-600 hover:bg-red-700 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
            >
              {{ deleting ? 'Menghapus...' : 'Ya, Hapus' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
definePageMeta({
  layout: 'admin'
})

const config = useRuntimeConfig()
const apiBase = config.public.apiBase

const loading = ref(true)
const coupons = ref([])
const searchQuery = ref('')
const filterStatus = ref('all')

// Modal state
const showModal = ref(false)
const modalMode = ref('create')
const submitting = ref(false)
const form = ref({
  code: '',
  discount_type: 'percentage',
  discount_value: null,
  max_discount: null,
  usage_limit: null,
  per_user_limit: 1,
  valid_from: '',
  valid_until: ''
})

// Delete modal state
const showDeleteModal = ref(false)
const couponToDelete = ref(null)
const deleting = ref(false)

// Computed
const activeCoupons = computed(() => {
  return coupons.value.filter(c => c.is_active && !isExpired(c)).length
})

const expiredCoupons = computed(() => {
  return coupons.value.filter(c => isExpired(c)).length
})

const totalUsage = computed(() => {
  return coupons.value.reduce((sum, c) => sum + (c.usage_count || 0), 0)
})

const filteredCoupons = computed(() => {
  let result = coupons.value
  
  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(c => c.code.toLowerCase().includes(query))
  }
  
  // Status filter
  if (filterStatus.value !== 'all') {
    result = result.filter(c => {
      const status = getCouponStatus(c)
      return status === filterStatus.value
    })
  }
  
  return result
})

// Methods
const loadCoupons = async () => {
  loading.value = true
  try {
    const token = useCookie('token')
    const response = await $fetch(`${apiBase}/api/admin/coupons`, {
      headers: { Authorization: `Bearer ${token.value}` }
    })
    coupons.value = response.coupons || []
  } catch (error) {
    console.error('Failed to load coupons:', error)
  } finally {
    loading.value = false
  }
}

const isExpired = (coupon) => {
  if (!coupon.valid_until) return false
  return new Date(coupon.valid_until) < new Date()
}

const getCouponStatus = (coupon) => {
  if (isExpired(coupon)) return 'expired'
  if (!coupon.is_active) return 'inactive'
  return 'active'
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount || 0)
}

const openModal = (mode, coupon = null) => {
  modalMode.value = mode
  if (mode === 'edit' && coupon) {
    form.value = {
      id: coupon.id,
      code: coupon.code,
      discount_type: coupon.discount_type,
      discount_value: coupon.discount_value,
      max_discount: coupon.max_discount,
      usage_limit: coupon.usage_limit,
      per_user_limit: coupon.per_user_limit || 1,
      valid_from: coupon.valid_from ? new Date(coupon.valid_from).toISOString().slice(0, 16) : '',
      valid_until: coupon.valid_until ? new Date(coupon.valid_until).toISOString().slice(0, 16) : ''
    }
  } else {
    form.value = {
      code: '',
      discount_type: 'percentage',
      discount_value: null,
      max_discount: null,
      usage_limit: null,
      per_user_limit: 1,
      valid_from: '',
      valid_until: ''
    }
  }
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const submitForm = async () => {
  submitting.value = true
  try {
    const token = useCookie('token')
    const payload = {
      code: form.value.code.toUpperCase(),
      discount_type: form.value.discount_type,
      discount_value: form.value.discount_value,
      max_discount: form.value.max_discount || null,
      usage_limit: form.value.usage_limit || null,
      per_user_limit: form.value.per_user_limit || 1,
      valid_from: form.value.valid_from ? new Date(form.value.valid_from).toISOString() : null,
      valid_until: form.value.valid_until ? new Date(form.value.valid_until).toISOString() : null
    }
    
    if (modalMode.value === 'create') {
      await $fetch(`${apiBase}/api/admin/coupons`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${token.value}` },
        body: payload
      })
    } else {
      await $fetch(`${apiBase}/api/admin/coupons/${form.value.id}`, {
        method: 'PUT',
        headers: { Authorization: `Bearer ${token.value}` },
        body: payload
      })
    }
    
    closeModal()
    await loadCoupons()
  } catch (error) {
    console.error('Failed to save coupon:', error)
    alert(error.data?.error || 'Gagal menyimpan kupon')
  } finally {
    submitting.value = false
  }
}

const toggleCouponStatus = async (coupon) => {
  try {
    const token = useCookie('token')
    await $fetch(`${apiBase}/api/admin/coupons/${coupon.id}`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${token.value}` },
      body: { is_active: !coupon.is_active }
    })
    await loadCoupons()
  } catch (error) {
    console.error('Failed to toggle status:', error)
  }
}

const confirmDelete = (coupon) => {
  couponToDelete.value = coupon
  showDeleteModal.value = true
}

const deleteCoupon = async () => {
  deleting.value = true
  try {
    const token = useCookie('token')
    await $fetch(`${apiBase}/api/admin/coupons/${couponToDelete.value.id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token.value}` }
    })
    showDeleteModal.value = false
    await loadCoupons()
  } catch (error) {
    console.error('Failed to delete coupon:', error)
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadCoupons()
})
</script>
