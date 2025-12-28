<template>
  <div class="relative" ref="notifRef">
    <!-- Bell Button -->
    <button @click="open = !open" class="relative p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
      </svg>
      <span v-if="unreadCount > 0" class="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
    </button>
    
    <!-- Dropdown -->
    <Transition name="dropdown">
      <div v-if="open" class="absolute right-0 mt-2 w-80 bg-white rounded-xl shadow-lg border border-neutral-200 z-50 overflow-hidden">
        <!-- Header -->
        <div class="px-4 py-3 border-b border-neutral-100 flex items-center justify-between">
          <h3 class="font-semibold text-neutral-900">Notifikasi</h3>
          <button v-if="unreadCount > 0" @click="markAllRead" class="text-xs text-primary-600 hover:text-primary-700">Tandai sudah dibaca</button>
        </div>
        
        <div class="max-h-80 overflow-y-auto">
          <template v-if="notifications.length > 0">
            <div v-for="notif in notifications" :key="notif.id" class="px-4 py-3 hover:bg-neutral-50 cursor-pointer border-b border-neutral-50 last:border-0 transition-colors" :class="{ 'bg-primary-50/50': !notif.is_read }">
              <div class="flex gap-3">
                <div class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0" :class="notifIconClass(notif.type)">
                  <svg v-if="notif.type === 'course'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                  </svg>
                  <svg v-else-if="notif.type === 'achievement'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/>
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-neutral-900" :class="{ 'font-medium': !notif.is_read }">{{ notif.title }}</p>
                  <p class="text-xs text-neutral-500 mt-0.5">{{ notif.time }}</p>
                </div>
              </div>
            </div>
          </template>
          
          <!-- Empty State -->
          <div v-else class="px-4 py-8 text-center">
            <svg class="w-12 h-12 mx-auto text-neutral-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
            </svg>
            <p class="text-sm text-neutral-500">Belum ada notifikasi</p>
          </div>
        </div>
        
        <!-- Footer -->
        <div class="px-4 py-2 border-t border-neutral-100 bg-neutral-50">
          <a href="#" class="block text-center text-sm text-primary-600 hover:text-primary-700 font-medium">Lihat Semua Notifikasi</a>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
interface Notification {
  id: string
  type: string
  title: string
  message?: string
  reference_id?: string
  reference_type?: string
  is_read: boolean
  created_at: string
  time: string
}

const api = useApi()
const notifRef = ref<HTMLElement | null>(null)
const open = ref(false)

// Notifications from API
const notifications = ref<Notification[]>([])

const unreadCount = computed(() => notifications.value.filter(n => !n.is_read).length)

// Fetch notifications
const fetchNotifications = async () => {
  try {
    const response = await api.fetch<{ notifications: Notification[], unread_count: number }>('/api/notifications')
    if (response?.notifications) {
      notifications.value = response.notifications
    }
  } catch (err) {
    console.error('[Notifications] Error fetching:', err)
  }
}

const notifIconClass = (type: string) => {
  switch (type) {
    case 'course': return 'bg-primary-100 text-primary-600'
    case 'achievement': return 'bg-warm-100 text-warm-600'
    default: return 'bg-neutral-100 text-neutral-600'
  }
}

const markAllRead = async () => {
  try {
    await api.fetch('/api/notifications/read-all', { method: 'PUT' })
    notifications.value.forEach(n => n.is_read = true)
  } catch (err) {
    console.error('[Notifications] Error marking all read:', err)
  }
}

onMounted(() => {
  document.addEventListener('click', (e) => {
    if (notifRef.value && !notifRef.value.contains(e.target as Node)) {
      open.value = false
    }
  })
  
  // Fetch notifications on mount
  fetchNotifications()
})
</script>
