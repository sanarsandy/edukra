<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Manajemen Pengguna</h1>
        <p class="text-neutral-500 mt-1">Kelola semua pengguna platform</p>
      </div>
      <div class="flex gap-3">
        <button @click="handleExport" class="px-4 py-2.5 text-sm font-medium text-neutral-600 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          Export CSV
        </button>
        <button @click="openAddModal" class="btn-admin">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
          </svg>
          Tambah Pengguna
        </button>
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
          placeholder="Cari pengguna..." 
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 focus:border-transparent text-sm"
        />
      </div>
      <select v-model="roleFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Role</option>
        <option value="admin">Admin</option>
        <option value="instructor">Instructor</option>
        <option value="student">Student</option>
      </select>
      <select v-model="statusFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Status</option>
        <option value="active">Aktif</option>
        <option value="inactive">Nonaktif</option>
      </select>
    </div>

    <!-- Bulk Actions Bar -->
    <Transition name="slide-down">
      <div v-if="selectedUsers.length > 0" class="bg-admin-600 text-white px-6 py-3 rounded-xl mb-4 flex items-center justify-between">
        <span class="text-sm font-medium">{{ selectedUsers.length }} pengguna dipilih</span>
        <div class="flex items-center gap-3">
          <button @click="exportSelectedUsers" class="px-3 py-1.5 text-sm font-medium bg-white/20 hover:bg-white/30 rounded-lg transition-colors flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            Export
          </button>
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

    <!-- Users Table -->
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
              <th class="px-4 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Pengguna</th>
              <th class="px-4 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Role</th>
              <th class="px-4 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-4 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Bergabung</th>
              <th class="px-4 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kursus</th>
              <th class="px-4 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-for="user in filteredUsers" :key="user.id" class="hover:bg-neutral-50 transition-colors" :class="selectedUsers.includes(user.id) ? 'bg-admin-50' : ''">
              <td class="px-4 py-4">
                <input 
                  type="checkbox" 
                  :checked="selectedUsers.includes(user.id)"
                  @change="toggleSelectUser(user.id)"
                  class="w-4 h-4 text-admin-600 rounded border-neutral-300 focus:ring-admin-500"
                />
              </td>
              <td class="px-4 py-4">
                <div class="flex items-center gap-3">
                  <div 
                    class="w-10 h-10 rounded-full flex items-center justify-center text-white font-semibold text-sm"
                    :class="getAvatarColor(user.email)"
                  >
                    {{ getInitials(user.full_name || user.email) }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-neutral-900">{{ user.full_name || 'Tanpa Nama' }}</p>
                    <p class="text-xs text-neutral-500">{{ user.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-4">
                <span 
                  class="px-2.5 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-admin-100 text-admin-700': user.role === 'admin',
                    'bg-primary-100 text-primary-700': user.role === 'instructor',
                    'bg-neutral-100 text-neutral-700': user.role === 'student'
                  }"
                >
                  {{ user.role }}
                </span>
              </td>
              <td class="px-4 py-4">
                <span 
                  class="px-2.5 py-1 text-xs font-medium rounded-full"
                  :class="user.status === 'active' ? 'bg-accent-100 text-accent-700' : 'bg-neutral-100 text-neutral-500'"
                >
                  {{ user.status === 'active' ? 'Aktif' : 'Nonaktif' }}
                </span>
              </td>
              <td class="px-4 py-4 text-sm text-neutral-500">{{ formatDate(user.created_at) }}</td>
              <td class="px-4 py-4">
                <span v-if="user.enrollment_count > 0" class="text-sm font-medium text-primary-600">{{ user.enrollment_count }} kursus</span>
                <span v-else class="text-sm text-neutral-400">-</span>
              </td>
              <td class="px-4 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openEditModal(user)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                    </svg>
                  </button>
                  <button @click="openDeleteModal(user)" class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
    <!-- Pagination -->
    <div class="px-6 py-4 border-t border-neutral-100 flex items-center justify-between">
      <p class="text-sm text-neutral-500">
        Menampilkan {{ (page - 1) * limit + 1 }}-{{ Math.min(page * limit, total) }} dari {{ total }} pengguna
      </p>
      <div class="flex gap-2">
        <button 
          @click="prevPage" 
          :disabled="page === 1"
          class="px-3 py-1.5 text-sm text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Prev
        </button>
        
        <button 
          v-for="p in displayedPages" 
          :key="p"
          @click="setPage(p)"
          class="px-3 py-1.5 text-sm rounded-lg transition-colors"
          :class="p === page ? 'text-white bg-admin-600' : 'text-neutral-600 bg-neutral-100 hover:bg-neutral-200'"
        >
          {{ p }}
        </button>
        
        <button 
          @click="nextPage" 
          :disabled="page === totalPages"
          class="px-3 py-1.5 text-sm text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>
  </div>

  <!-- Add/Edit User Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-neutral-900">{{ isEditing ? 'Edit Pengguna' : 'Tambah Pengguna' }}</h3>
            <button @click="closeModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveUser" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Lengkap <span class="text-red-500">*</span></label>
              <input 
                v-model="form.name"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan nama lengkap"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Email <span class="text-red-500">*</span></label>
              <input 
                v-model="form.email"
                type="email" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="email@example.com"
              />
            </div>
            <div v-if="!isEditing">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Password <span class="text-red-500">*</span></label>
              <input 
                v-model="form.password"
                type="password" 
                :required="!isEditing"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Minimal 8 karakter"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Role <span class="text-red-500">*</span></label>
              <select 
                v-model="form.role"
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              >
                <option value="">Pilih role</option>
                <option value="student">Student</option>
                <option value="instructor">Instructor</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Status</label>
              <div class="flex items-center gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="form.status" value="active" class="w-4 h-4 text-admin-600 focus:ring-admin-500">
                  <span class="text-sm text-neutral-700">Aktif</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="form.status" value="inactive" class="w-4 h-4 text-admin-600 focus:ring-admin-500">
                  <span class="text-sm text-neutral-700">Nonaktif</span>
                </label>
              </div>
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors">
                {{ isEditing ? 'Simpan Perubahan' : 'Tambah Pengguna' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Delete Confirmation Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Pengguna</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus <strong>{{ selectedUser?.full_name || selectedUser?.email }}</strong>? Tindakan ini tidak dapat dibatalkan.</p>
          </div>
          <div class="flex gap-3">
            <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Batal
            </button>
            <button @click="confirmDeleteUser" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50">
              Hapus
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Bulk Delete Confirmation Modal -->
    <Transition name="fade">
      <div v-if="showBulkDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showBulkDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus {{ selectedUsers.length }} Pengguna</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus <strong>{{ selectedUsers.length }} pengguna</strong> yang dipilih? Tindakan ini tidak dapat dibatalkan.</p>
          </div>
          <div class="flex gap-3">
            <button @click="showBulkDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Batal
            </button>
            <button @click="confirmBulkDelete" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50">
              {{ saving ? 'Menghapus...' : 'Hapus Semua' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast Notification -->
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
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

useHead({
  title: 'Manajemen Pengguna - Admin'
})

// API Composable
const { 
  loading: apiLoading, 
  error: apiError, 
  users, 
  total, 
  page, 
  limit, 
  totalPages, 
  fetchUsers, 
  createUser, 
  updateUser, 
  deleteUser: deleteUserApi,
  nextPage,
  prevPage,
  setPage
} = useUsers()

// Export composable
const { exportUsers } = useExport()

const handleExport = () => {
  if (filteredUsers.value.length === 0) {
    showToast('Tidak ada data untuk di-export', 'error')
    return
  }
  exportUsers(filteredUsers.value)
  showToast('Data berhasil di-export ke CSV')
}

const searchQuery = ref('')
const roleFilter = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const showDeleteModal = ref(false)
const showBulkDeleteModal = ref(false)
const isEditing = ref(false)
const selectedUser = ref<any>(null)
const saving = ref(false)

// Bulk selection
const selectedUsers = ref<string[]>([])

const isAllSelected = computed(() => {
  return filteredUsers.value.length > 0 && selectedUsers.value.length === filteredUsers.value.length
})

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedUsers.value = []
  } else {
    selectedUsers.value = filteredUsers.value.map(u => u.id)
  }
}

const toggleSelectUser = (userId: string) => {
  const index = selectedUsers.value.indexOf(userId)
  if (index === -1) {
    selectedUsers.value.push(userId)
  } else {
    selectedUsers.value.splice(index, 1)
  }
}

const clearSelection = () => {
  selectedUsers.value = []
}

const exportSelectedUsers = () => {
  const selectedData = filteredUsers.value.filter(u => selectedUsers.value.includes(u.id))
  exportUsers(selectedData)
  showToast(`${selectedData.length} pengguna di-export ke CSV`)
}

const openBulkDeleteModal = () => {
  showBulkDeleteModal.value = true
}

const confirmBulkDelete = async () => {
  saving.value = true
  let successCount = 0
  
  for (const userId of selectedUsers.value) {
    const success = await deleteUserApi(userId)
    if (success) successCount++
  }
  
  saving.value = false
  showBulkDeleteModal.value = false
  
  if (successCount > 0) {
    showToast(`${successCount} pengguna berhasil dihapus`)
    clearSelection()
    await fetchUsers()
  } else {
    showToast('Gagal menghapus pengguna', 'error')
  }
}

const toast = ref({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

const form = ref({
  name: '',
  email: '',
  password: '',
  role: '',
  status: 'active'
})

const avatarColors = ['bg-primary-500', 'bg-accent-500', 'bg-warm-500', 'bg-rose-500', 'bg-cyan-500', 'bg-indigo-500']

// Load users on mount
onMounted(async () => {
  await fetchUsers()
})

// Watch search and filters to reset pagination
watch([searchQuery, roleFilter, statusFilter], () => {
  setPage(1)
})

const displayedPages = computed(() => {
  const delta = 2
  const range = []
  for (let i = Math.max(2, page.value - delta); i <= Math.min(totalPages.value - 1, page.value + delta); i++) {
    range.push(i)
  }

  if (page.value - delta > 2) {
    range.unshift('...')
  }
  if (page.value + delta < totalPages.value - 1) {
    range.push('...')
  }

  range.unshift(1)
  if (totalPages.value !== 1) {
    range.push(totalPages.value)
  }

  return range.filter(x => typeof x === 'number') // Simplified for now, handling ellipsis requires more logic
})

// Simplified displayedPages for now to avoid complexity with ellipsis
const simpleDisplayedPages = computed(() => {
  const pages = []
  const maxButtons = 5
  let start = Math.max(1, page.value - 2)
  let end = Math.min(totalPages.value, start + maxButtons - 1)
  
  if (end - start < maxButtons - 1) {
    start = Math.max(1, end - maxButtons + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})


const filteredUsers = computed(() => {
  if (!users.value) return []
  return users.value.filter(user => {
    const matchesSearch = (user.full_name || '').toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                          user.email.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesRole = !roleFilter.value || user.role === roleFilter.value
    // Status filter using is_active from API
    const userStatus = user.is_active ? 'active' : 'inactive'
    const matchesStatus = !statusFilter.value || userStatus === statusFilter.value
    return matchesSearch && matchesRole && matchesStatus
  }).map(user => ({
    ...user,
    status: user.is_active ? 'active' : 'inactive'
  }))
})

const getInitials = (name: string) => {
  if (!name) return '?'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const getAvatarColor = (email: string) => {
  const hash = email.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return avatarColors[hash % avatarColors.length]
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const openAddModal = () => {
  isEditing.value = false
  form.value = { name: '', email: '', password: '', role: 'student', status: 'active' }
  showModal.value = true
}

const openEditModal = (user: any) => {
  isEditing.value = true
  selectedUser.value = user
  form.value = { 
    name: user.full_name || '', 
    email: user.email, 
    password: '', 
    role: user.role, 
    status: user.is_active ? 'active' : 'inactive'
  }
  showModal.value = true
}

const openDeleteModal = (user: any) => {
  selectedUser.value = user
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
  form.value = { name: '', email: '', password: '', role: '', status: 'active' }
}

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const saveUser = async () => {
  saving.value = true
  
  try {
    if (isEditing.value && selectedUser.value) {
      // Update existing user via API
      const result = await updateUser(selectedUser.value.id, {
        full_name: form.value.name,
        role: form.value.role as any,
        is_active: form.value.status === 'active'
      })
      
      if (result) {
        showToast('Pengguna berhasil diperbarui')
        await fetchUsers() // Refresh list
      } else {
        showToast('Gagal memperbarui pengguna', 'error')
      }
    } else {
      // Create new user via API
      if (!form.value.password) {
        showToast('Password wajib diisi', 'error')
        saving.value = false
        return
      }
      
      const result = await createUser({
        email: form.value.email,
        password: form.value.password,
        full_name: form.value.name,
        role: form.value.role || 'student'
      })
      
      if (result) {
        showToast('Pengguna berhasil ditambahkan')
        await fetchUsers() // Refresh list
      } else {
        showToast('Gagal menambahkan pengguna', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}

const confirmDeleteUser = async () => {
  if (!selectedUser.value) return
  
  saving.value = true
  try {
    const success = await deleteUserApi(selectedUser.value.id)
    if (success) {
      showToast('Pengguna berhasil dihapus')
      await fetchUsers() // Refresh list
    } else {
      showToast('Gagal menghapus pengguna', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedUser.value = null
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
.slide-down-enter-active, .slide-down-leave-active {
  transition: all 0.3s ease;
}
.slide-down-enter-from, .slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
