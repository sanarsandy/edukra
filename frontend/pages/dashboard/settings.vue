<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Pengaturan Akun</h1>
      <p class="text-neutral-500 mt-1">Kelola profil dan preferensi Anda</p>
    </div>

    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Sidebar Navigation -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
          <nav class="divide-y divide-neutral-100">
            <button 
              v-for="tab in tabs" 
              :key="tab.id"
              @click="activeTab = tab.id"
              class="w-full flex items-center gap-3 p-4 text-left hover:bg-neutral-50 transition-colors"
              :class="activeTab === tab.id ? 'bg-primary-50 border-r-2 border-primary-600' : ''"
            >
              <div 
                class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="activeTab === tab.id ? 'bg-primary-100 text-primary-600' : 'bg-neutral-100 text-neutral-500'"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="tab.icon"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium" :class="activeTab === tab.id ? 'text-primary-700' : 'text-neutral-700'">{{ tab.name }}</p>
                <p class="text-xs text-neutral-500">{{ tab.description }}</p>
              </div>
            </button>
          </nav>
        </div>
      </div>

      <!-- Content -->
      <div class="lg:col-span-2">
        <!-- Profile Section -->
        <div v-if="activeTab === 'profile'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Informasi Profil</h3>
          
          <!-- Avatar -->
          <div class="flex items-center gap-4 mb-6 pb-6 border-b border-neutral-100">
            <div class="w-20 h-20 bg-primary-600 rounded-full flex items-center justify-center text-white text-2xl font-bold">
              {{ profile.name.charAt(0) }}
            </div>
            <div>
              <input type="file" ref="avatarInput" accept="image/*" class="hidden" @change="handleAvatarChange">
              <button @click="$refs.avatarInput.click()" class="text-sm text-primary-600 font-medium hover:text-primary-700">Ubah Foto</button>
              <p class="text-xs text-neutral-500 mt-1">JPG, PNG. Maksimal 2MB</p>
            </div>
          </div>
          
          <!-- Form -->
          <form @submit.prevent="saveProfile" class="space-y-4">
            <div class="grid sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Lengkap <span class="text-red-500">*</span></label>
                <input 
                  v-model="profile.name"
                  type="text"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Email <span class="text-red-500">*</span></label>
                <input 
                  v-model="profile.email"
                  type="email"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">No. Telepon</label>
              <input 
                v-model="profile.phone"
                type="tel" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
                placeholder="+62 xxx xxxx xxxx"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Bio</label>
              <textarea 
                v-model="profile.bio"
                rows="3" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm resize-none"
                placeholder="Ceritakan sedikit tentang diri Anda..."
              ></textarea>
            </div>
            <div class="pt-4 flex gap-3">
              <button type="submit" :disabled="saving" class="btn-primary">
                <svg v-if="saving" class="w-4 h-4 mr-2 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                {{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}
              </button>
            </div>
          </form>
        </div>

        <!-- Password Section -->
        <div v-if="activeTab === 'password'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Ubah Password</h3>
          
          <form @submit.prevent="changePassword" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Password Saat Ini <span class="text-red-500">*</span></label>
              <div class="relative">
                <input 
                  v-model="passwordForm.currentPassword"
                  :type="showCurrentPassword ? 'text' : 'password'"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm pr-10"
                />
                <button type="button" @click="showCurrentPassword = !showCurrentPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600">
                  <svg v-if="showCurrentPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Password Baru <span class="text-red-500">*</span></label>
              <div class="relative">
                <input 
                  v-model="passwordForm.newPassword"
                  :type="showNewPassword ? 'text' : 'password'"
                  required
                  minlength="8"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm pr-10"
                />
                <button type="button" @click="showNewPassword = !showNewPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600">
                  <svg v-if="showNewPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                </button>
              </div>
              <p class="text-xs text-neutral-400 mt-1">Minimal 8 karakter</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Konfirmasi Password Baru <span class="text-red-500">*</span></label>
              <input 
                v-model="passwordForm.confirmPassword"
                type="password" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent text-sm"
              />
              <p v-if="passwordForm.confirmPassword && passwordForm.newPassword !== passwordForm.confirmPassword" class="text-xs text-red-500 mt-1">Password tidak cocok</p>
            </div>
            <div class="pt-4">
              <button 
                type="submit" 
                :disabled="savingPassword || (passwordForm.newPassword !== passwordForm.confirmPassword)" 
                class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg v-if="savingPassword" class="w-4 h-4 mr-2 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
                {{ savingPassword ? 'Mengubah...' : 'Ubah Password' }}
              </button>
            </div>
          </form>
        </div>

        <!-- Notifications Section -->
        <div v-if="activeTab === 'notifications'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Preferensi Notifikasi</h3>
          
          <div class="space-y-4">
            <div v-for="notif in notifications" :key="notif.id" class="flex items-center justify-between py-3 border-b border-neutral-100 last:border-0">
              <div>
                <p class="text-sm font-medium text-neutral-900">{{ notif.title }}</p>
                <p class="text-xs text-neutral-500">{{ notif.description }}</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="notif.enabled" @change="saveNotificationPreference(notif)" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

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
  layout: 'dashboard',
  middleware: 'auth'
})

useHead({
  title: 'Pengaturan - LearnHub'
})

const activeTab = ref('profile')
const saving = ref(false)
const savingPassword = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)

const toast = ref({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

const tabs = [
  { id: 'profile', name: 'Profil', description: 'Informasi dasar', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z' },
  { id: 'password', name: 'Password', description: 'Keamanan akun', icon: 'M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z' },
  { id: 'notifications', name: 'Notifikasi', description: 'Preferensi pesan', icon: 'M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9' }
]

const profile = ref({
  name: '',
  email: '',
  phone: '',
  bio: ''
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const notifications = ref([
  { id: 1, title: 'Email Kursus Baru', description: 'Dapatkan notifikasi saat ada kursus baru', enabled: true },
  { id: 2, title: 'Reminder Belajar', description: 'Pengingat harian untuk melanjutkan belajar', enabled: true },
  { id: 3, title: 'Promo & Diskon', description: 'Informasi promo dan diskon kursus', enabled: false },
  { id: 4, title: 'Update Kursus', description: 'Notifikasi saat kursus yang Anda ikuti di-update', enabled: true }
])

// API integration
const api = useApi()
const { user } = useAuth()

// Load user profile on mount
onMounted(async () => {
  if (user.value) {
    profile.value.name = user.value.full_name || ''
    profile.value.email = user.value.email || ''
  }
  
  // Fetch full profile from API
  try {
    const userData = await api.fetch<any>('/api/me')
    if (userData) {
      profile.value.name = userData.full_name || ''
      profile.value.email = userData.email || ''
      profile.value.bio = userData.bio || ''
      profile.value.phone = userData.phone || ''
    }
  } catch (err) {
    // Use cached user data
  }
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const handleAvatarChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    // Handle avatar upload
    showToast('Foto profil berhasil diubah')
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    const result = await api.fetch<any>('/api/me', {
      method: 'PUT',
      body: {
        full_name: profile.value.name,
        bio: profile.value.bio || null,
        phone: profile.value.phone || null
      }
    })
    
    if (result) {
      // Update user cookie so header updates immediately
      const userCookie = useCookie('user')
      if (userCookie.value) {
        const currentUser = typeof userCookie.value === 'string' 
          ? JSON.parse(userCookie.value) 
          : userCookie.value
        currentUser.full_name = profile.value.name
        userCookie.value = JSON.stringify(currentUser)
      }
      showToast('Profil berhasil disimpan')
    } else {
      showToast('Gagal menyimpan profil', 'error')
    }
  } catch (err: any) {
    showToast(err.message || 'Gagal menyimpan profil', 'error')
  } finally {
    saving.value = false
  }
}

const changePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    showToast('Password tidak cocok', 'error')
    return
  }
  
  savingPassword.value = true
  try {
    const result = await api.fetch<any>('/api/me/password', {
      method: 'PUT',
      body: {
        current_password: passwordForm.value.currentPassword,
        new_password: passwordForm.value.newPassword
      }
    })
    
    if (result?.message) {
      passwordForm.value = { currentPassword: '', newPassword: '', confirmPassword: '' }
      showToast('Password berhasil diubah')
    } else {
      showToast('Gagal mengubah password', 'error')
    }
  } catch (err: any) {
    showToast(err.message || 'Gagal mengubah password', 'error')
  } finally {
    savingPassword.value = false
  }
}

const saveNotificationPreference = (notif: any) => {
  // Notification preferences can be saved to localStorage or API
  showToast(`${notif.title} ${notif.enabled ? 'diaktifkan' : 'dinonaktifkan'}`)
}
</script>

<style scoped>
.slide-up-enter-active, .slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from, .slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
