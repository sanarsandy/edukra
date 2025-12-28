<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Dashboard Admin</h1>
      <p class="text-neutral-500 mt-1">Selamat datang kembali, Admin.</p>
    </div>
    
    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-admin-600"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-6 mb-8">
      <p class="text-red-600">{{ error }}</p>
      <button @click="loadDashboard" class="mt-2 text-sm text-red-700 hover:underline">Coba lagi</button>
    </div>

    <template v-else>
      <!-- Stats Grid -->
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        <div class="bg-white rounded-xl p-5 border border-neutral-200">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 bg-admin-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-admin-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-2xl font-bold text-neutral-900">{{ formatNumber(adminStats?.total_users || 0) }}</p>
          <p class="text-sm text-neutral-500">Total Pengguna</p>
        </div>
        
        <div class="bg-white rounded-xl p-5 border border-neutral-200">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 bg-primary-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
              </svg>
            </div>
          </div>
          <p class="text-2xl font-bold text-neutral-900">{{ formatNumber(adminStats?.total_courses || 0) }}</p>
          <p class="text-sm text-neutral-500">Total Kursus</p>
        </div>
        
        <div class="bg-white rounded-xl p-5 border border-neutral-200">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 bg-warm-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-2xl font-bold text-neutral-900">{{ formatCurrency(adminStats?.total_revenue || 0) }}</p>
          <p class="text-sm text-neutral-500">Pendapatan</p>
        </div>
        
        <div class="bg-white rounded-xl p-5 border border-neutral-200">
          <div class="flex items-center justify-between mb-3">
            <div class="w-10 h-10 bg-accent-100 rounded-lg flex items-center justify-center">
              <svg class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
          </div>
          <p class="text-2xl font-bold text-neutral-900">{{ formatNumber(adminStats?.total_transactions || 0) }}</p>
          <p class="text-sm text-neutral-500">Transaksi</p>
        </div>
      </div>
      
      <!-- Charts Section -->
      <div class="grid lg:grid-cols-2 gap-6 mb-8">
        <!-- Revenue Chart -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="font-semibold text-neutral-900">Pendapatan</h3>
              <p class="text-sm text-neutral-500">6 bulan terakhir</p>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full bg-primary-500"></span>
              <span class="text-neutral-600">Revenue</span>
            </div>
          </div>
          <!-- Line Chart -->
          <div class="relative h-32">
            <svg class="w-full h-full" viewBox="0 0 300 100" preserveAspectRatio="none">
              <!-- Grid lines -->
              <line x1="0" y1="25" x2="300" y2="25" stroke="#f3f4f6" stroke-width="1"/>
              <line x1="0" y1="50" x2="300" y2="50" stroke="#f3f4f6" stroke-width="1"/>
              <line x1="0" y1="75" x2="300" y2="75" stroke="#f3f4f6" stroke-width="1"/>
              <!-- Area fill -->
              <path 
                :d="getAreaPath(chartData.revenue)" 
                fill="url(#revenueGradient)" 
                opacity="0.3"
              />
              <!-- Line -->
              <path 
                :d="getLinePath(chartData.revenue)" 
                fill="none" 
                stroke="#1E4976" 
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <!-- Data points -->
              <circle 
                v-for="(point, index) in getChartPoints(chartData.revenue)" 
                :key="'rev-point-'+index"
                :cx="point.x" 
                :cy="point.y" 
                r="4" 
                fill="#1E4976"
                class="transition-all duration-300 hover:r-6"
              />
              <defs>
                <linearGradient id="revenueGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#1E4976"/>
                  <stop offset="100%" stop-color="#1E4976" stop-opacity="0"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <!-- Labels -->
          <div class="flex justify-between mt-2">
            <span v-for="(label, index) in chartData.labels" :key="'label-rev-'+index" class="text-xs text-neutral-500">{{ label }}</span>
          </div>
          <div class="mt-4 pt-4 border-t border-neutral-100 flex justify-between text-sm">
            <span class="text-neutral-500">Total Periode</span>
            <span class="font-semibold text-neutral-900">{{ formatCurrency(chartData.revenue.reduce((a: number, b: number) => a + b, 0)) }}</span>
          </div>
        </div>
        
        <!-- Users Growth Chart -->
        <div class="bg-white rounded-xl border border-neutral-200 p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h3 class="font-semibold text-neutral-900">Pertumbuhan Pengguna</h3>
              <p class="text-sm text-neutral-500">6 bulan terakhir</p>
            </div>
            <div class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full bg-accent-500"></span>
              <span class="text-neutral-600">Users</span>
            </div>
          </div>
          <!-- Line Chart -->
          <div class="relative h-32">
            <svg class="w-full h-full" viewBox="0 0 300 100" preserveAspectRatio="none">
              <!-- Grid lines -->
              <line x1="0" y1="25" x2="300" y2="25" stroke="#f3f4f6" stroke-width="1"/>
              <line x1="0" y1="50" x2="300" y2="50" stroke="#f3f4f6" stroke-width="1"/>
              <line x1="0" y1="75" x2="300" y2="75" stroke="#f3f4f6" stroke-width="1"/>
              <!-- Area fill -->
              <path 
                :d="getAreaPath(chartData.users)" 
                fill="url(#usersGradient)" 
                opacity="0.3"
              />
              <!-- Line -->
              <path 
                :d="getLinePath(chartData.users)" 
                fill="none" 
                stroke="#7FA99B" 
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <!-- Data points -->
              <circle 
                v-for="(point, index) in getChartPoints(chartData.users)" 
                :key="'usr-point-'+index"
                :cx="point.x" 
                :cy="point.y" 
                r="4" 
                fill="#7FA99B"
                class="transition-all duration-300 hover:r-6"
              />
              <defs>
                <linearGradient id="usersGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                  <stop offset="0%" stop-color="#7FA99B"/>
                  <stop offset="100%" stop-color="#7FA99B" stop-opacity="0"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <!-- Labels -->
          <div class="flex justify-between mt-2">
            <span v-for="(label, index) in chartData.labels" :key="'label-usr-'+index" class="text-xs text-neutral-500">{{ label }}</span>
          </div>
          <div class="mt-4 pt-4 border-t border-neutral-100 flex justify-between text-sm">
            <span class="text-neutral-500">Total Pengguna Baru</span>
            <span class="font-semibold text-neutral-900">{{ chartData.users.reduce((a: number, b: number) => a + b, 0) }} user</span>
          </div>
        </div>
      </div>
      
      <!-- Two Column Layout -->
      <div class="grid lg:grid-cols-2 gap-6">
        <!-- Recent Users -->
        <div class="bg-white rounded-xl border border-neutral-200">
          <div class="flex items-center justify-between p-5 border-b border-neutral-100">
            <h3 class="font-semibold text-neutral-900">Pengguna Baru</h3>
            <NuxtLink to="/admin/users" class="text-sm text-admin-600 hover:text-admin-700 font-medium">Lihat Semua →</NuxtLink>
          </div>
          <div class="divide-y divide-neutral-100">
            <div v-if="recentUsers.length === 0" class="p-6 text-center text-neutral-500">
              Belum ada pengguna
            </div>
            <div v-for="user in recentUsers" :key="user.id" class="flex items-center gap-4 p-4 hover:bg-neutral-50 transition-colors">
              <div class="w-10 h-10 bg-primary-500 rounded-full flex items-center justify-center text-white font-semibold text-sm flex-shrink-0">
                {{ getInitials(user.full_name || user.email) }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-neutral-900 truncate">{{ user.full_name || 'Tanpa Nama' }}</p>
                <p class="text-xs text-neutral-500 truncate">{{ user.email }}</p>
              </div>
              <span class="text-xs px-2 py-1 rounded-full" :class="getRoleBadge(user.role)">
                {{ user.role }}
              </span>
            </div>
          </div>
        </div>
        
        <!-- Recent Transactions -->
        <div class="bg-white rounded-xl border border-neutral-200">
          <div class="flex items-center justify-between p-5 border-b border-neutral-100">
            <h3 class="font-semibold text-neutral-900">Transaksi Terbaru</h3>
            <NuxtLink to="/admin/transactions" class="text-sm text-admin-600 hover:text-admin-700 font-medium">Lihat Semua →</NuxtLink>
          </div>
          <div class="divide-y divide-neutral-100">
            <div v-if="transactions.length === 0" class="p-6 text-center text-neutral-500">
              Belum ada transaksi
            </div>
            <div v-for="tx in transactions" :key="tx.id" class="flex items-center justify-between p-4 hover:bg-neutral-50 transition-colors">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0" :class="getStatusBg(tx.status)">
                  <svg v-if="tx.status === 'success'" class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 13l4 4L19 7"/>
                  </svg>
                  <svg v-else class="w-5 h-5 text-warm-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-medium text-neutral-900">{{ tx.payment_gateway || 'Payment' }}</p>
                  <p class="text-xs text-neutral-500">{{ formatDate(tx.created_at) }}</p>
                </div>
              </div>
              <div class="text-right">
                <span class="text-sm font-semibold" :class="getStatusColor(tx.status)">{{ formatCurrency(tx.amount) }}</span>
                <p class="text-xs capitalize" :class="getStatusColor(tx.status)">{{ tx.status }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Activity Log -->
      <div class="mt-8 bg-white rounded-xl border border-neutral-200">
        <div class="flex items-center justify-between p-5 border-b border-neutral-100">
          <h3 class="font-semibold text-neutral-900">Aktivitas Terbaru</h3>
          <span class="text-xs text-neutral-400">Live updates</span>
        </div>
        <div class="p-5">
          <div class="relative">
            <!-- Empty state -->
            <div v-if="activityLog.length === 0" class="text-center py-8 text-neutral-500">
              <p>Belum ada aktivitas terbaru</p>
            </div>
            
            <template v-else>
              <!-- Timeline line -->
              <div class="absolute left-5 top-0 bottom-0 w-px bg-neutral-200"></div>
              
              <!-- Activity items -->
              <div class="space-y-6">
                <div v-for="activity in activityLog" :key="activity.id" class="relative flex gap-4">
                  <div 
                    class="relative z-10 w-10 h-10 rounded-full flex items-center justify-center"
                    :class="getActivityColor(activity.type)"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="getActivityIcon(activity.type)"/>
                    </svg>
                  </div>
                  <div class="flex-1 min-w-0 pt-1.5">
                    <p class="text-sm text-neutral-900">
                      <span class="font-medium">{{ activity.user }}</span>{{ ' ' }}
                      <span class="text-neutral-500">{{ activity.action }}</span>{{ ' ' }}
                      <span v-if="activity.target" class="font-medium">{{ activity.target }}</span>
                    </p>
                    <p class="text-xs text-neutral-400 mt-0.5">{{ activity.time }}</p>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

useHead({
  title: 'Admin Dashboard - LearnHub'
})

const api = useApi()

const { 
  loading, 
  error, 
  adminStats, 
  recentUsers, 
  transactions, 
  fetchAdminDashboard 
} = useDashboard()

// Chart data interface
interface ChartData {
  labels: string[]
  revenue: number[]
  users: number[]
}

// Chart data
const chartData = ref<ChartData>({
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun'],
  revenue: [0, 0, 0, 0, 0, 0],
  users: [0, 0, 0, 0, 0, 0]
})

// Fetch chart data
const fetchChartData = async () => {
  try {
    const response = await api.fetch<{
      labels: string[]
      revenue: number[]
      users: number[]
    }>('/api/admin/dashboard/charts')
    
    console.log('[Dashboard] Chart data response:', response)
    
    if (response && response.labels && response.revenue && response.users) {
      chartData.value = {
        labels: response.labels,
        revenue: response.revenue,
        users: response.users
      }
      console.log('[Dashboard] Chart data updated:', chartData.value)
    }
  } catch (err) {
    console.error('[Dashboard] Error fetching chart data:', err)
    // Use default sample data if API fails
  }
}

// Line chart helper functions
const getChartPoints = (data: number[]) => {
  const max = Math.max(...data, 1) // Ensure at least 1 to avoid division by zero
  const width = 300
  const height = 100
  const padding = 10
  
  return data.map((value, index) => ({
    x: padding + (index * (width - 2 * padding)) / (data.length - 1 || 1),
    y: height - padding - ((value / max) * (height - 2 * padding))
  }))
}

const getLinePath = (data: number[]) => {
  const points = getChartPoints(data)
  if (points.length === 0) return ''
  
  return points.map((point, index) => 
    `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`
  ).join(' ')
}

const getAreaPath = (data: number[]) => {
  const points = getChartPoints(data)
  if (points.length === 0) return ''
  
  const width = 300
  const height = 100
  const padding = 10
  
  // Start from bottom left
  let path = `M ${padding} ${height - padding}`
  
  // Draw line to first point then through all points
  points.forEach(point => {
    path += ` L ${point.x} ${point.y}`
  })
  
  // Close path to bottom right and back
  path += ` L ${width - padding} ${height - padding} Z`
  
  return path
}

// Load dashboard data
const loadDashboard = async () => {
  await fetchAdminDashboard()
  await fetchChartData()
  await fetchActivityLog()
}

onMounted(() => {
  loadDashboard()
})

// Activity log types
interface Activity {
  id: string
  type: string
  user: string
  action: string
  target: string
  time: string
}

// Activity log
const activityLog = ref<Activity[]>([])

// Fetch activity log from API
const fetchActivityLog = async () => {
  try {
    const response = await api.fetch<Activity[]>('/api/admin/dashboard/activities?limit=5')
    
    if (response && Array.isArray(response)) {
      activityLog.value = response
    }
  } catch (err) {
    console.error('[Dashboard] Error fetching activity log:', err)
    // Fallback to empty array
    activityLog.value = []
  }
}

const getActivityColor = (type: string) => {
  const colors: Record<string, string> = {
    user: 'bg-primary-100 text-primary-600',
    course: 'bg-accent-100 text-accent-600',
    payment: 'bg-warm-100 text-warm-600',
    edit: 'bg-admin-100 text-admin-600',
    delete: 'bg-red-100 text-red-600'
  }
  return colors[type] || 'bg-neutral-100 text-neutral-600'
}

const getActivityIcon = (type: string) => {
  const icons: Record<string, string> = {
    user: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
    course: 'M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253',
    payment: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z',
    edit: 'M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z',
    delete: 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16'
  }
  return icons[type] || icons.edit
}

// Helpers
const formatNumber = (num: number) => {
  return new Intl.NumberFormat('id-ID').format(num)
}

const formatCurrency = (amount: number) => {
  if (amount >= 1000000) {
    return `Rp ${(amount / 1000000).toFixed(0)}M`
  } else if (amount >= 1000) {
    return `Rp ${(amount / 1000).toFixed(0)}K`
  }
  return `Rp ${amount}`
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const getInitials = (name: string) => {
  if (!name) return '?'
  return name.split(' ').map(n => n[0]).join('').toUpperCase().slice(0, 2)
}

const getRoleBadge = (role: string) => {
  const badges: Record<string, string> = {
    admin: 'bg-admin-100 text-admin-700',
    instructor: 'bg-primary-100 text-primary-700',
    student: 'bg-accent-100 text-accent-700'
  }
  return badges[role] || 'bg-neutral-100 text-neutral-700'
}

const getStatusBg = (status: string) => {
  return status === 'success' ? 'bg-accent-100' : 'bg-warm-100'
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    success: 'text-accent-600',
    pending: 'text-warm-600',
    failed: 'text-red-600',
    refunded: 'text-neutral-600'
  }
  return colors[status] || 'text-neutral-600'
}
</script>
