<template>
  <div>
    <div class="mb-8">
      <h1 class="text-2xl font-display font-bold text-slate-900">Dashboard</h1>
      <p class="text-slate-600 mt-1">Selamat datang kembali! Berikut ringkasan kursus Anda.</p>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- Total Courses -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/50 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-gradient-to-br from-teal-400 to-cyan-500 rounded-xl flex items-center justify-center shadow-lg shadow-teal-500/20">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
          </div>
        </div>
        <p class="text-3xl font-bold text-slate-900">{{ stats.total_courses }}</p>
        <p class="text-sm text-slate-500 mt-1">Total Kursus</p>
      </div>

      <!-- Published -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/50 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-gradient-to-br from-emerald-400 to-green-500 rounded-xl flex items-center justify-center shadow-lg shadow-emerald-500/20">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-3xl font-bold text-slate-900">{{ stats.published_courses }}</p>
        <p class="text-sm text-slate-500 mt-1">Dipublish</p>
      </div>

      <!-- Pending Review -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/50 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-gradient-to-br from-amber-400 to-orange-500 rounded-xl flex items-center justify-center shadow-lg shadow-amber-500/20">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-3xl font-bold text-slate-900">{{ stats.pending_review }}</p>
        <p class="text-sm text-slate-500 mt-1">Menunggu Review</p>
      </div>

      <!-- Total Students -->
      <div class="bg-white rounded-2xl p-6 border border-slate-200/50 shadow-sm hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-gradient-to-br from-violet-400 to-purple-500 rounded-xl flex items-center justify-center shadow-lg shadow-violet-500/20">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
            </svg>
          </div>
        </div>
        <p class="text-3xl font-bold text-slate-900">{{ stats.total_students }}</p>
        <p class="text-sm text-slate-500 mt-1">Total Siswa</p>
      </div>
    </div>

    <!-- Quick Actions & Recent Courses -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Quick Actions -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-2xl border border-slate-200/50 shadow-sm overflow-hidden">
          <div class="p-6 border-b border-slate-100">
            <h2 class="font-semibold text-slate-900">Aksi Cepat</h2>
          </div>
          <div class="p-4 space-y-2">
            <NuxtLink to="/instructor/courses/create" class="flex items-center gap-4 p-4 hover:bg-slate-50 rounded-xl transition-colors group">
              <div class="w-10 h-10 bg-gradient-to-br from-teal-400 to-cyan-500 rounded-lg flex items-center justify-center shadow-md shadow-teal-500/20 group-hover:scale-105 transition-transform">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </div>
              <div>
                <p class="font-medium text-slate-900">Buat Kursus Baru</p>
                <p class="text-sm text-slate-500">Mulai membuat kursus</p>
              </div>
            </NuxtLink>
            
            <NuxtLink to="/instructor/courses" class="flex items-center gap-4 p-4 hover:bg-slate-50 rounded-xl transition-colors group">
              <div class="w-10 h-10 bg-slate-100 rounded-lg flex items-center justify-center group-hover:scale-105 transition-transform">
                <svg class="w-5 h-5 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
                </svg>
              </div>
              <div>
                <p class="font-medium text-slate-900">Kelola Kursus</p>
                <p class="text-sm text-slate-500">Lihat semua kursus Anda</p>
              </div>
            </NuxtLink>
          </div>
        </div>

        <!-- Rating Overview -->
        <div class="bg-white rounded-2xl border border-slate-200/50 shadow-sm overflow-hidden mt-6">
          <div class="p-6 border-b border-slate-100">
            <h2 class="font-semibold text-slate-900">Rating Rata-rata</h2>
          </div>
          <div class="p-6 text-center">
            <div class="text-4xl font-bold text-slate-900 mb-2">
              {{ stats.average_rating?.toFixed(1) || '0.0' }}
            </div>
            <div class="flex justify-center gap-1 mb-2">
              <svg v-for="i in 5" :key="i" class="w-5 h-5" :class="i <= Math.round(stats.average_rating || 0) ? 'text-amber-400' : 'text-slate-200'" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
              </svg>
            </div>
            <p class="text-sm text-slate-500">dari semua kursus</p>
          </div>
        </div>
      </div>

      <!-- Recent Courses -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-2xl border border-slate-200/50 shadow-sm overflow-hidden">
          <div class="p-6 border-b border-slate-100 flex items-center justify-between">
            <h2 class="font-semibold text-slate-900">Kursus Terbaru</h2>
            <NuxtLink to="/instructor/courses" class="text-sm text-teal-600 hover:text-teal-700 font-medium">Lihat Semua</NuxtLink>
          </div>
          <div v-if="loading" class="p-8 text-center">
            <div class="inline-flex items-center gap-2 text-slate-500">
              <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Memuat...
            </div>
          </div>
          <div v-else-if="recentCourses.length === 0" class="p-8 text-center">
            <svg class="w-12 h-12 mx-auto text-slate-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
            <p class="text-slate-500 mb-4">Belum ada kursus</p>
            <NuxtLink to="/instructor/courses/create" class="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-teal-500 to-cyan-500 text-white font-medium rounded-lg hover:shadow-lg hover:shadow-teal-500/25 transition-all">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              Buat Kursus Pertama
            </NuxtLink>
          </div>
          <div v-else class="divide-y divide-slate-100">
            <NuxtLink 
              v-for="course in recentCourses" 
              :key="course.id" 
              :to="`/instructor/courses/${course.id}`"
              class="flex items-center gap-4 p-4 hover:bg-slate-50 transition-colors"
            >
              <div class="w-16 h-12 bg-slate-100 rounded-lg overflow-hidden flex-shrink-0">
                <img v-if="course.thumbnail_url" :src="course.thumbnail_url" class="w-full h-full object-cover" :alt="course.title" />
                <div v-else class="w-full h-full flex items-center justify-center">
                  <svg class="w-6 h-6 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                  </svg>
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-medium text-slate-900 truncate">{{ course.title }}</h3>
                <div class="flex items-center gap-3 mt-1">
                  <span class="text-sm text-slate-500">{{ course.lessons_count || 0 }} materi</span>
                  <span :class="getStatusClass(course.status)" class="px-2 py-0.5 text-xs font-medium rounded-full">{{ getStatusLabel(course.status) }}</span>
                </div>
              </div>
              <svg class="w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5l7 7-7 7"/>
              </svg>
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'instructor',
  middleware: ['instructor']
})

const instructorPanel = useInstructorPanel()

const loading = ref(true)
const stats = ref({
  total_courses: 0,
  published_courses: 0,
  draft_courses: 0,
  pending_review: 0,
  total_students: 0,
  average_rating: 0
})
const recentCourses = ref<any[]>([])

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    draft: 'bg-slate-100 text-slate-700',
    pending_review: 'bg-amber-100 text-amber-700',
    approved: 'bg-blue-100 text-blue-700',
    rejected: 'bg-rose-100 text-rose-700',
    published: 'bg-emerald-100 text-emerald-700'
  }
  return classes[status] || 'bg-slate-100 text-slate-700'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    draft: 'Draft',
    pending_review: 'Menunggu Review',
    approved: 'Disetujui',
    rejected: 'Ditolak',
    published: 'Dipublish'
  }
  return labels[status] || status
}

onMounted(async () => {
  const data = await instructorPanel.fetchDashboard()
  if (data) {
    stats.value = data.stats
    recentCourses.value = data.recent_courses || []
  }
  loading.value = false
})
</script>


