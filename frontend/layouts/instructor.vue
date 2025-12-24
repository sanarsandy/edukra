<template>
  <div class="min-h-screen bg-neutral-100">
    <!-- Sidebar - Same as Admin Theme (Dark) -->
    <aside 
      class="fixed top-0 left-0 z-40 h-screen bg-neutral-900 hidden lg:block transition-all duration-300 overflow-hidden"
      :class="sidebarCollapsed ? 'w-[72px]' : 'w-72'"
    >
      <div class="h-full flex flex-col">
        <!-- Logo -->
        <div class="h-16 flex items-center border-b border-neutral-800" :class="sidebarCollapsed ? 'justify-center' : 'justify-between px-4'">
          <div class="flex items-center gap-3 min-w-0" :class="sidebarCollapsed ? 'justify-center' : ''">
            <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
            <div v-if="!sidebarCollapsed"><span class="font-display font-bold text-lg text-white">EDUKRA</span><span class="block text-xs text-admin-400 -mt-0.5">Instructor Panel</span></div>
          </div>
          <button v-if="!sidebarCollapsed" @click="toggleSidebar" class="p-2 text-neutral-500 hover:text-white hover:bg-neutral-800 rounded-lg transition-colors flex-shrink-0" title="Tutup Sidebar">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
        
        <!-- Navigation -->
        <nav class="flex-1 py-6 overflow-y-auto overflow-x-hidden" :class="sidebarCollapsed ? 'px-3' : 'px-4'">
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Overview</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/instructor" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/instructor', true) ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"/>
                </svg>
                <span v-if="!sidebarCollapsed">Dashboard</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Dashboard
                </div>
              </NuxtLink>
            </div>
          </div>
          
          <div class="mb-8">
            <p v-if="!sidebarCollapsed" class="px-3 mb-3 text-xs font-semibold text-neutral-500 uppercase tracking-wider">Manajemen</p>
            <div class="space-y-1">
              <NuxtLink 
                to="/instructor/courses" 
                class="flex items-center text-sm font-medium rounded-lg transition-all group relative"
                :class="[
                  isActive('/instructor/courses') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white',
                  sidebarCollapsed ? 'justify-center p-3' : 'px-3 py-2.5'
                ]"
              >
                <svg class="w-5 h-5 flex-shrink-0" :class="sidebarCollapsed ? '' : 'mr-3'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <span v-if="!sidebarCollapsed">Kursus Saya</span>
                <span v-if="!sidebarCollapsed && stats.total_courses > 0" class="ml-auto bg-neutral-700 text-neutral-300 text-xs font-medium px-2 py-0.5 rounded-full">{{ stats.total_courses }}</span>
                <div v-if="sidebarCollapsed" class="absolute left-full ml-2 px-2 py-1 bg-white text-neutral-900 text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 pointer-events-none whitespace-nowrap transition-opacity z-50">
                  Kursus Saya
                </div>
              </NuxtLink>
            </div>
          </div>
        </nav>
        
        <!-- Toggle Button (collapsed state) -->
        <div v-if="sidebarCollapsed" class="p-3 border-t border-neutral-800">
          <button @click="toggleSidebar" class="w-full flex items-center justify-center p-2.5 text-neutral-500 hover:text-white hover:bg-neutral-800 rounded-lg transition-colors" title="Buka Sidebar">
            <svg class="w-5 h-5 rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 19l-7-7 7-7m8 14l-7-7 7-7"/></svg>
          </button>
        </div>
      </div>
    </aside>
    
    <!-- Mobile Header -->
    <header class="lg:hidden fixed top-0 left-0 right-0 z-50 h-16 bg-neutral-900 flex items-center justify-between px-4">
      <div class="flex items-center space-x-3">
        <button @click="mobileMenuOpen = true" class="p-2 -ml-2 text-neutral-400 hover:text-white">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        <img src="/logo.png" alt="EDUKRA Logo" class="w-8 h-8 object-contain" />
        <span class="text-xs text-admin-400 font-medium">Instructor</span>
      </div>
    </header>
    
    <!-- Main Content -->
    <div class="transition-all duration-300" :class="sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-72'">
      <!-- Top Bar -->
      <header class="hidden lg:flex h-16 bg-white border-b border-neutral-200 items-center justify-between px-8">
        <div class="flex items-center gap-4">
          <span class="px-3 py-1.5 bg-admin-50 text-admin-700 text-xs font-semibold rounded-full">Instructor Panel</span>
          <div class="relative">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/></svg>
            <input type="text" placeholder="Cari..." class="w-64 pl-9 pr-4 py-2 bg-neutral-100 border-0 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 focus:bg-white text-sm placeholder-neutral-500 transition-all"/>
          </div>
        </div>
        <div class="flex items-center gap-3">
          <!-- Notification Bell -->
          <div class="relative" ref="notificationRef">
            <button @click="notificationOpen = !notificationOpen" class="relative p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>
              <span v-if="unreadCount > 0" class="absolute top-1 right-1 w-4 h-4 bg-red-500 rounded-full text-[10px] text-white font-bold flex items-center justify-center">{{ unreadCount > 9 ? '9+' : unreadCount }}</span>
            </button>
            
            <!-- Notification Dropdown -->
            <Transition name="dropdown">
              <div v-if="notificationOpen" class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-lg border border-neutral-200 overflow-hidden z-50">
                <div class="p-4 border-b border-neutral-100 flex items-center justify-between">
                  <h3 class="font-semibold text-neutral-900">Notifikasi</h3>
                </div>
                <div class="max-h-80 overflow-y-auto">
                  <div v-if="notifications.length === 0" class="p-6 text-center text-neutral-500 text-sm">
                    Tidak ada notifikasi
                  </div>
                  <div 
                    v-for="notif in notifications" 
                    :key="notif.id" 
                    @click="handleNotificationClick(notif)"
                    class="p-4 hover:bg-neutral-50 border-b border-neutral-100 last:border-0 cursor-pointer transition-colors"
                    :class="!notif.is_read ? 'bg-admin-50/50' : ''"
                  >
                    <div class="flex gap-3">
                      <div 
                        class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
                        :class="getNotificationColor(notif.type)"
                      >
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="getNotificationIcon(notif.type)"/>
                        </svg>
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="text-sm text-neutral-900" :class="!notif.is_read ? 'font-medium' : ''">{{ notif.title }}</p>
                        <p class="text-xs text-neutral-500 mt-0.5 truncate">{{ notif.message }}</p>
                        <p class="text-xs text-neutral-400 mt-1">{{ formatTimeAgo(notif.created_at) }}</p>
                      </div>
                      <div v-if="!notif.is_read" class="w-2 h-2 bg-admin-500 rounded-full flex-shrink-0 mt-2"></div>
                    </div>
                  </div>
                </div>
              </div>
            </Transition>
          </div>
          <div class="relative" ref="userDropdownRef">
            <button @click="userDropdownOpen = !userDropdownOpen" class="flex items-center gap-3 p-2 hover:bg-neutral-100 rounded-lg transition-colors">
              <div class="w-8 h-8 bg-admin-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">{{ userInitial }}</div>
              <div class="text-left hidden sm:block"><p class="text-sm font-medium text-neutral-900">{{ userName }}</p></div>
              <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 9l-7 7-7-7"/></svg>
            </button>
            <Transition name="dropdown">
              <div v-if="userDropdownOpen" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-50">
                <div class="border-t border-neutral-100 my-2"></div>
                <button @click="handleLogout" class="flex items-center w-full px-4 py-2 text-sm text-red-600 hover:bg-red-50">Keluar</button>
              </div>
            </Transition>
          </div>
        </div>
      </header>
      
      <!-- Page Content -->
      <main class="p-6 lg:p-8 pt-20 lg:pt-8">
        <slot />
      </main>
    </div>
    
    <!-- Mobile Menu Overlay -->
    <Transition name="fade">
      <div v-if="mobileMenuOpen" class="fixed inset-0 z-50 bg-black/60 lg:hidden" @click="mobileMenuOpen = false"></div>
    </Transition>
    
    <!-- Mobile Menu -->
    <Transition name="slide">
      <aside v-if="mobileMenuOpen" class="fixed top-0 left-0 z-50 w-72 h-screen bg-neutral-900 lg:hidden">
        <div class="h-full flex flex-col">
          <div class="h-16 flex items-center justify-between px-6 border-b border-neutral-800">
            <div class="flex items-center space-x-3">
              <img src="/logo.png" alt="EDUKRA Logo" class="w-10 h-10 object-contain" />
              <span class="font-display font-bold text-lg text-white">EDUKRA Instructor</span>
            </div>
            <button @click="mobileMenuOpen = false" class="p-2 text-neutral-500 hover:text-white">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <nav class="flex-1 px-4 py-6 overflow-y-auto">
            <div class="space-y-1">
              <NuxtLink to="/instructor" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/instructor', true) ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z"/>
                </svg>
                Dashboard
              </NuxtLink>
              <NuxtLink to="/instructor/courses" @click="mobileMenuOpen = false" class="flex items-center px-3 py-3 text-sm font-medium rounded-lg" :class="isActive('/instructor/courses') ? 'bg-admin-600 text-white' : 'text-neutral-400 hover:bg-neutral-800 hover:text-white'">
                <svg class="w-5 h-5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                Kursus Saya
              </NuxtLink>
            </div>
          </nav>
        </div>
      </aside>
    </Transition>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const instructorPanel = useInstructorPanel()

const mobileMenuOpen = ref(false)
const sidebarCollapsed = ref(false)
const userDropdownOpen = ref(false)
const userDropdownRef = ref<HTMLElement | null>(null)
const notificationOpen = ref(false)
const notificationRef = ref<HTMLElement | null>(null)

const userCookie = useCookie<{ full_name?: string } | null>('user')
const userName = computed(() => userCookie.value?.full_name || 'Instructor')
const userInitial = computed(() => userName.value.charAt(0).toUpperCase())

const stats = ref({
  total_courses: 0,
  published_courses: 0,
  draft_courses: 0,
  pending_review: 0
})

const notifications = ref<any[]>([])
const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

const getNotificationColor = (type: string) => {
  const colors: Record<string, string> = {
    course_approved: 'bg-accent-100 text-accent-600',
    course_rejected: 'bg-red-100 text-red-600',
    course_published: 'bg-primary-100 text-primary-600',
    course_unpublished: 'bg-warm-100 text-warm-600'
  }
  return colors[type] || 'bg-neutral-100 text-neutral-600'
}

const getNotificationIcon = (type: string) => {
  const icons: Record<string, string> = {
    course_approved: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    course_rejected: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
    course_published: 'M5 13l4 4L19 7',
    course_unpublished: 'M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636'
  }
  return icons[type] || 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
}

const formatTimeAgo = (dateStr: string) => {
  const date = new Date(dateStr)
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  if (seconds < 60) return 'Baru saja'
  if (seconds < 3600) return `${Math.floor(seconds / 60)} menit lalu`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} jam lalu`
  return `${Math.floor(seconds / 86400)} hari lalu`
}

const handleNotificationClick = async (notif: any) => {
  if (!notif.is_read) {
    await instructorPanel.markNotificationRead(notif.id)
    notif.is_read = true
  }
  notificationOpen.value = false
  if (notif.reference_type === 'course' && notif.reference_id) {
    navigateTo(`/instructor/courses/${notif.reference_id}`)
  }
}

const isActive = (path: string, exact = false) => {
  if (exact) return route.path === path
  return route.path === path || route.path.startsWith(path + '/')
}

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  if (process.client) {
    localStorage.setItem('instructor-sidebar-collapsed', String(sidebarCollapsed.value))
  }
}

const fetchDashboardStats = async () => {
  const data = await instructorPanel.fetchDashboard()
  if (data?.stats) {
    stats.value = data.stats
  }
}

const fetchNotifications = async () => {
  const data = await instructorPanel.fetchNotifications()
  if (data?.notifications) {
    notifications.value = data.notifications
  }
}

const handleLogout = () => {
  instructorPanel.instructorLogout()
}

onMounted(() => {
  if (process.client) {
    const saved = localStorage.getItem('instructor-sidebar-collapsed')
    if (saved === 'true') {
      sidebarCollapsed.value = true
    }
    document.addEventListener('click', (e) => {
      if (userDropdownRef.value && !userDropdownRef.value.contains(e.target as Node)) {
        userDropdownOpen.value = false
      }
      if (notificationRef.value && !notificationRef.value.contains(e.target as Node)) {
        notificationOpen.value = false
      }
    })
    
    fetchDashboardStats()
    fetchNotifications()
  }
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-enter-active, .slide-leave-active { transition: transform 0.3s ease; }
.slide-enter-from, .slide-leave-to { transform: translateX(-100%); }
.dropdown-enter-active, .dropdown-leave-active { transition: all 0.2s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-10px); }
</style>
