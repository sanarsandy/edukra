<template>
  <div class="immersive-study-page min-h-screen bg-neutral-50 flex flex-col">
    <!-- Top Bar -->
    <header class="h-14 bg-white border-b border-neutral-200 flex items-center justify-between px-4 flex-shrink-0 z-30">
      <div class="flex items-center gap-4">
        <NuxtLink 
          :to="`/dashboard/courses/${courseId}`" 
          class="flex items-center gap-2 text-neutral-600 hover:text-neutral-900 transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
          </svg>
          <span class="text-sm font-medium hidden sm:inline">Kembali</span>
        </NuxtLink>
        
        <div class="h-6 w-px bg-neutral-200 hidden sm:block"></div>
        
        <div class="text-sm text-neutral-900 font-medium truncate max-w-[200px] sm:max-w-md">
          {{ currentLesson?.title || 'Pilih Materi' }}
        </div>
      </div>
      
      <div class="flex items-center gap-3">
        <!-- Progress Badge -->
        <div class="hidden sm:flex items-center gap-2 px-3 py-1.5 bg-neutral-100 rounded-full">
          <div class="w-2 h-2 rounded-full" :class="progressPercentage >= 100 ? 'bg-accent-500' : 'bg-primary-500'"></div>
          <span class="text-xs font-medium text-neutral-700">{{ progressPercentage }}% selesai</span>
        </div>
        
        <!-- Toggle Sidebar (Mobile) -->
        <button 
          @click="sidebarOpen = !sidebarOpen"
          class="lg:hidden p-2 text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100 rounded-lg"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
        
        <!-- Close Button -->
        <NuxtLink 
          :to="`/dashboard/courses/${courseId}`"
          class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </NuxtLink>
      </div>
    </header>
    
    <!-- Main Layout -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Sidebar Overlay (Mobile) -->
      <div 
        v-if="sidebarOpen"
        class="fixed inset-0 bg-black/40 z-40 lg:hidden"
        @click="sidebarOpen = false"
      ></div>
      
      <!-- Left Sidebar -->
      <aside 
        class="fixed lg:relative inset-y-14 lg:inset-y-0 left-0 z-40 w-80 bg-white border-r border-neutral-200 flex flex-col transform transition-transform duration-200 lg:transform-none"
        :class="sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
      >
        <!-- Course Info -->
        <div class="p-4 border-b border-neutral-100">
          <h2 class="font-semibold text-neutral-900 text-sm line-clamp-2">{{ course?.title }}</h2>
          <p class="text-xs text-neutral-500 mt-1">{{ course?.instructor?.full_name }}</p>
          
          <!-- Progress Bar -->
          <div class="mt-3">
            <div class="flex justify-between text-xs text-neutral-500 mb-1">
              <span>Progress</span>
              <span>{{ completedLessonIds.length }}/{{ actualLessons.length }}</span>
            </div>
            <div class="h-1.5 bg-neutral-100 rounded-full overflow-hidden">
              <div 
                class="h-full rounded-full transition-all duration-300"
                :class="progressPercentage >= 100 ? 'bg-accent-500' : 'bg-primary-500'"
                :style="{ width: progressPercentage + '%' }"
              ></div>
            </div>
          </div>
        </div>
        
        <!-- Lessons List -->
        <div class="flex-1 overflow-y-auto p-2">
          <template v-if="lessonsTree.length > 0">
            <StudentLessonTree
              v-for="(lesson, idx) in lessonsTree"
              :key="lesson.id"
              :lesson="lesson"
              :lesson-number="idx + 1"
              :selected-id="currentLessonId"
              :completed-ids="completedLessonIds"
              @select="handleLessonSelect"
            />
          </template>
          <div v-else class="p-8 text-center text-neutral-400 text-sm">
            Belum ada materi
          </div>
        </div>
        
        <!-- Footer Actions -->
        <div class="p-4 border-t border-neutral-100">
          <button 
            v-if="currentLesson && !completedLessonIds.includes(currentLesson.id)"
            @click="markLessonComplete"
            class="w-full flex items-center justify-center gap-2 px-4 py-2.5 bg-accent-500 text-white rounded-lg hover:bg-accent-600 transition-colors text-sm font-medium"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            Tandai Selesai
          </button>
          <div v-else-if="currentLesson" class="flex items-center justify-center gap-2 px-4 py-2.5 bg-accent-50 text-accent-600 rounded-lg text-sm font-medium">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            Sudah Selesai
          </div>
          
          <!-- Next Lesson Button -->
          <button 
            v-if="nextLesson"
            @click="handleLessonSelect(nextLesson)"
            class="w-full mt-2 flex items-center justify-center gap-2 px-4 py-2 bg-neutral-100 text-neutral-700 rounded-lg hover:bg-neutral-200 transition-colors text-sm"
          >
            Materi Selanjutnya
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
      </aside>
      
      <!-- Main Content Area -->
      <main class="flex-1 overflow-y-auto bg-neutral-100">
        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center h-full">
          <div class="animate-spin w-10 h-10 border-4 border-primary-500 border-t-transparent rounded-full"></div>
        </div>
        
        <!-- No Lesson Selected -->
        <div v-else-if="!currentLesson" class="flex items-center justify-center h-full">
          <div class="text-center p-8">
            <div class="w-20 h-20 mx-auto mb-4 bg-neutral-200 rounded-full flex items-center justify-center">
              <svg class="w-10 h-10 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
              </svg>
            </div>
            <p class="text-neutral-500">Pilih materi dari sidebar untuk memulai belajar</p>
          </div>
        </div>
        
        <!-- Content Display -->
        <div v-else class="h-full flex flex-col">
          <!-- Video Content -->
          <div 
            v-if="currentLesson.type === 'video'" 
            class="flex-1 bg-neutral-900 flex items-center justify-center protected-content"
            @contextmenu.prevent
          >
            <!-- Embed Video -->
            <template v-if="isEmbedVideo(currentLesson.videoUrl)">
              <SecureEmbedPlayer 
                v-if="isYouTubeVideo(currentLesson.videoUrl) && getYouTubeVideoId(currentLesson.videoUrl)"
                :video-id="getYouTubeVideoId(currentLesson.videoUrl)!"
                :autoplay="true"
                @ended="onVideoEnded"
                class="w-full h-full"
              />
              <iframe 
                v-else
                :src="getEmbedUrl(currentLesson.videoUrl)"
                class="w-full h-full"
                frameborder="0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; fullscreen"
                allowfullscreen
              ></iframe>
            </template>
            <!-- Direct Video -->
            <template v-else>
              <div class="relative w-full h-full">
                <video 
                  ref="videoPlayer"
                  :src="secureVideoUrl || getFileUrl(currentLesson.videoUrl)"
                  controls
                  controlsList="nodownload noplaybackrate"
                  disablepictureinpicture
                  class="w-full h-full object-contain"
                  @ended="onVideoEnded"
                  @contextmenu.prevent
                  @playing="videoIsPlaying = true"
                  @pause="videoIsPlaying = false"
                ></video>
                <!-- Watermark -->
                <VideoWatermark 
                  v-if="videoIsPlaying && currentUserEmail"
                  :user-email="currentUserEmail"
                  :opacity="0.12"
                  :rotate-interval="25"
                />
              </div>
            </template>
          </div>
          
          <!-- PDF Content -->
          <div 
            v-else-if="currentLesson.type === 'pdf' || currentLesson.type === 'document'" 
            class="flex-1 bg-neutral-50"
          >
            <ClientOnly>
              <SecurePDFViewer
                :is-open="true"
                :pdf-url="pdfUrl"
                :title="currentLesson.title"
                :user-email="currentUserEmail"
                :is-completed="completedLessonIds.includes(currentLesson.id)"
                :inline-mode="true"
                @complete="markLessonComplete"
              />
              <template #fallback>
                <div class="flex items-center justify-center h-full">
                  <div class="animate-spin w-8 h-8 border-4 border-primary-500 border-t-transparent rounded-full"></div>
                </div>
              </template>
            </ClientOnly>
          </div>
          
          <!-- Quiz Content -->
          <div 
            v-else-if="currentLesson.type === 'quiz'" 
            class="flex-1 flex items-center justify-center bg-gradient-to-br from-accent-50 to-accent-100"
          >
            <div class="text-center p-8">
              <div class="w-24 h-24 mx-auto mb-6 bg-accent-500 rounded-full flex items-center justify-center shadow-lg">
                <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-accent-800 mb-2">{{ currentLesson.title }}</h3>
              <p class="text-accent-600 mb-6">Uji pemahaman Anda dengan mengerjakan kuis ini</p>
              <button 
                @click="openQuiz"
                class="px-8 py-3 bg-accent-600 text-white rounded-xl font-medium hover:bg-accent-700 transition-colors shadow-lg"
              >
                Mulai Kuis
              </button>
            </div>
          </div>
          
          <!-- Text Content -->
          <div 
            v-else-if="currentLesson.type === 'text'" 
            class="flex-1 overflow-y-auto bg-white"
          >
            <div class="max-w-4xl mx-auto p-6 lg:p-10">
              <h1 class="text-2xl font-bold text-neutral-900 mb-6">{{ currentLesson.title }}</h1>
              <div 
                class="prose prose-neutral max-w-none protected-content relative"
                @contextmenu.prevent
                @copy.prevent
                v-html="currentLesson.content || '<p class=text-neutral-500>Tidak ada konten</p>'"
              ></div>
              <!-- Watermark -->
              <VideoWatermark 
                v-if="currentUserEmail"
                :user-email="currentUserEmail"
                :opacity="0.06"
                :rotate-interval="30"
              />
            </div>
          </div>
          
          <!-- External Link Content -->
          <div 
            v-else-if="currentLesson.type === 'link'" 
            class="flex-1 flex items-center justify-center bg-neutral-50"
          >
            <div class="text-center p-8">
              <div class="w-20 h-20 mx-auto mb-4 bg-blue-100 rounded-full flex items-center justify-center">
                <svg class="w-10 h-10 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-neutral-900 mb-2">{{ currentLesson.title }}</h3>
              <p class="text-neutral-500 mb-4 text-sm">Konten ini akan dibuka di tab baru</p>
              <a 
                :href="currentLesson.videoUrl" 
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 px-6 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                Buka Tautan
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
                </svg>
              </a>
            </div>
          </div>
          
          <!-- Unknown Content Type -->
          <div v-else class="flex-1 flex items-center justify-center">
            <div class="text-center p-8">
              <div class="w-16 h-16 mx-auto mb-4 bg-neutral-200 rounded-full flex items-center justify-center">
                <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <p class="text-neutral-500">Tipe konten tidak didukung</p>
            </div>
          </div>
        </div>
      </main>
    </div>
    
    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3" :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'">
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
import { ref, computed, onMounted } from 'vue'

const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()
const courseId = computed(() => route.params.id as string)
const apiBase = config.public.apiBase || 'http://localhost:8080'

// Composables
const { course, fetchCourse, loading: loadingCourse } = useCourses()
const { checkEnrollment } = useEnrollments()
const { 
  completedLessonIds, 
  fetchCourseProgress, 
  markLessonComplete: markLessonCompleteApi 
} = useLessonProgress()
const { 
  getContentBlobUrl,
  isMinioObject 
} = useSecureContent()

// State
const sidebarOpen = ref(false)
const currentLessonId = ref<string | null>(null)
const secureVideoUrl = ref<string | null>(null)
const videoIsPlaying = ref(false)
const videoPlayer = ref<HTMLVideoElement | null>(null)
const pdfUrl = ref('')
const isEnrolled = ref(false)
const enrollmentData = ref<any>(null)
const currentUserEmail = ref('')

// Toast
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

definePageMeta({
  layout: false // Fullscreen layout
})

// Building lessons tree
const lessonsTree = computed(() => {
  if (!course.value?.lessons) return []
  
  const allLessons = course.value.lessons
  const lessonMap = new Map<string, any>()
  const roots: any[] = []
  
  allLessons.forEach((l: any) => {
    lessonMap.set(l.id, { ...l, children: [] })
  })
  
  allLessons.forEach((l: any) => {
    const node = lessonMap.get(l.id)!
    if (l.parent_id && lessonMap.has(l.parent_id)) {
      lessonMap.get(l.parent_id)!.children.push(node)
    } else {
      roots.push(node)
    }
  })
  
  const sortChildren = (items: any[]) => {
    items.sort((a, b) => (a.order_index || 0) - (b.order_index || 0))
    items.forEach(item => {
      if (item.children?.length) sortChildren(item.children)
    })
  }
  sortChildren(roots)
  
  return roots
})

// Flatten lessons for navigation
const flattenLessons = (items: any[]): any[] => {
  let result: any[] = []
  items.forEach(item => {
    if (!item.is_container) result.push(item)
    if (item.children?.length) {
      result = result.concat(flattenLessons(item.children))
    }
  })
  return result
}

const actualLessons = computed(() => flattenLessons(lessonsTree.value))

const currentLesson = computed(() => {
  if (!currentLessonId.value || !course.value?.lessons) return null
  const lesson = course.value.lessons.find((l: any) => l.id === currentLessonId.value)
  if (!lesson) return null
  return {
    id: lesson.id,
    title: lesson.title,
    description: lesson.description,
    type: lesson.content_type,
    videoUrl: lesson.video_url,
    content: lesson.content
  }
})

const progressPercentage = computed(() => {
  if (!actualLessons.value.length) return 0
  return Math.round((completedLessonIds.value.length / actualLessons.value.length) * 100)
})

const nextLesson = computed(() => {
  if (!currentLessonId.value) return actualLessons.value[0]
  const currentIndex = actualLessons.value.findIndex((l: any) => l.id === currentLessonId.value)
  if (currentIndex >= 0 && currentIndex < actualLessons.value.length - 1) {
    return actualLessons.value[currentIndex + 1]
  }
  return null
})

const loading = computed(() => loadingCourse.value)

// Video helpers
type VideoType = 'youtube' | 'vimeo' | 'gdrive' | 'zoom' | 'msstream' | 'sharepoint' | 'direct' | 'embed' | 'unknown'

const detectVideoType = (url: string | null | undefined): VideoType => {
  if (!url) return 'unknown'
  if (url.includes('youtube.com') || url.includes('youtu.be')) return 'youtube'
  if (url.includes('vimeo.com')) return 'vimeo'
  if (url.includes('drive.google.com')) return 'gdrive'
  if (url.includes('zoom.us/rec')) return 'zoom'
  if (url.includes('microsoftstream.com')) return 'msstream'
  if (url.includes('sharepoint.com') || url.includes('onedrive.live.com')) return 'sharepoint'
  if (/\.(mp4|webm|ogg|mov|m4v)(\?.*)?$/i.test(url)) return 'direct'
  if (url.includes('/embed') || url.includes('/player')) return 'embed'
  return 'direct'
}

const isEmbedVideo = (url: string | null | undefined): boolean => {
  const type = detectVideoType(url)
  return ['youtube', 'vimeo', 'gdrive', 'zoom', 'msstream', 'sharepoint', 'embed'].includes(type)
}

const isYouTubeVideo = (url: string | null | undefined): boolean => detectVideoType(url) === 'youtube'

const getYouTubeVideoId = (url: string): string | null => {
  const watchMatch = url.match(/youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)/)
  if (watchMatch) return watchMatch[1]
  const shortMatch = url.match(/youtu\.be\/([a-zA-Z0-9_-]+)/)
  if (shortMatch) return shortMatch[1]
  const embedMatch = url.match(/youtube\.com\/embed\/([a-zA-Z0-9_-]+)/)
  if (embedMatch) return embedMatch[1]
  return null
}

const getEmbedUrl = (url: string): string => {
  const type = detectVideoType(url)
  switch (type) {
    case 'youtube': {
      const ytParams = new URLSearchParams({
        rel: '0', modestbranding: '1', disablekb: '1', playsinline: '1'
      }).toString()
      const match = url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]+)/)
      if (match) return `https://www.youtube-nocookie.com/embed/${match[1]}?${ytParams}`
      return url
    }
    case 'vimeo': {
      const match = url.match(/vimeo\.com\/(\d+)/)
      if (match) return `https://player.vimeo.com/video/${match[1]}`
      return url
    }
    case 'gdrive': {
      const match = url.match(/drive\.google\.com\/file\/d\/([a-zA-Z0-9_-]+)/)
      if (match) return `https://drive.google.com/file/d/${match[1]}/preview`
      return url
    }
    default:
      return url
  }
}

const getFileUrl = (path: string | null | undefined): string => {
  if (!path) return ''
  if (path.startsWith('http://') || path.startsWith('https://')) return path
  if (path.startsWith('/uploads')) return `${apiBase}${path}`
  return ''
}

// Handlers
const handleLessonSelect = async (lesson: any) => {
  // Don't select containers
  if (lesson.is_container) return
  
  currentLessonId.value = lesson.id
  sidebarOpen.value = false // Close mobile sidebar
  secureVideoUrl.value = null
  
  // Load content based on type
  if (lesson.content_type === 'video' && lesson.video_url && isMinioObject(lesson.video_url)) {
    const url = await getContentBlobUrl(lesson.id)
    if (url) secureVideoUrl.value = url
  }
  
  if ((lesson.content_type === 'pdf' || lesson.content_type === 'document') && lesson.video_url) {
    if (isMinioObject(lesson.video_url)) {
      const url = await getContentBlobUrl(lesson.id)
      if (url) pdfUrl.value = url
    } else {
      pdfUrl.value = getFileUrl(lesson.video_url)
    }
  }
}

const onVideoEnded = () => {
  if (currentLesson.value && !completedLessonIds.value.includes(currentLesson.value.id)) {
    markLessonComplete()
  }
}

const markLessonComplete = async () => {
  if (!currentLesson.value || completedLessonIds.value.includes(currentLesson.value.id)) return
  
  await markLessonCompleteApi(currentLesson.value.id)
  showToast('Materi ditandai selesai!')
}

const openQuiz = () => {
  if (currentLesson.value) {
    localStorage.setItem('lastCourseId', courseId.value)
    router.push(`/dashboard/quiz/${currentLesson.value.id}?lessonId=${currentLesson.value.id}`)
  }
}

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

// Initialize
onMounted(async () => {
  // Get user email from localStorage (client-side only)
  try {
    const userData = localStorage.getItem('user')
    if (userData) {
      const user = JSON.parse(userData)
      currentUserEmail.value = user.email || ''
    }
  } catch (e) {
    // Ignore errors
  }
  
  // Fetch course data
  await fetchCourse(courseId.value)
  
  // Check enrollment
  const check = await checkEnrollment(courseId.value)
  isEnrolled.value = check.enrolled
  enrollmentData.value = check.enrollment
  
  // Redirect if not enrolled
  if (!check.enrolled) {
    showToast('Anda belum terdaftar di kursus ini', 'error')
    setTimeout(() => router.push(`/dashboard/courses/${courseId.value}`), 1500)
    return
  }
  
  // Fetch progress
  await fetchCourseProgress(courseId.value)
  
  // Select lesson from query or first uncompleted
  const lessonIdFromQuery = route.query.lesson as string
  if (lessonIdFromQuery) {
    const lesson = course.value?.lessons?.find((l: any) => l.id === lessonIdFromQuery)
    if (lesson) {
      handleLessonSelect(lesson)
      return
    }
  }
  
  // Select first uncompleted lesson
  const firstUncompleted = actualLessons.value.find(
    (l: any) => !completedLessonIds.value.includes(l.id)
  )
  if (firstUncompleted) {
    handleLessonSelect(firstUncompleted)
  } else if (actualLessons.value.length > 0) {
    handleLessonSelect(actualLessons.value[0])
  }
})

useHead({
  title: computed(() => course.value?.title ? `Belajar: ${course.value.title}` : 'Loading...')
})
</script>

<style scoped>
.immersive-study-page {
  height: 100vh;
  height: 100dvh;
}

.protected-content {
  -webkit-user-select: none;
  -moz-user-select: none;
  user-select: none;
}

/* Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}
.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(20px);
  opacity: 0;
}

/* Line clamp */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Prose styling for text content */
.prose h1, .prose h2, .prose h3 {
  color: #1f2937;
  margin-top: 1.5em;
  margin-bottom: 0.5em;
}
.prose p {
  margin-bottom: 1em;
  line-height: 1.75;
}
.prose ul, .prose ol {
  margin-left: 1.5em;
  margin-bottom: 1em;
}
.prose li {
  margin-bottom: 0.25em;
}
</style>
