<template>
  <div class="student-lesson-tree">
    <!-- Module Row -->
    <div 
      v-if="lesson.is_container"
      class="module-header"
    >
      <button 
        @click="toggleExpand"
        class="flex items-center gap-3 w-full p-3 hover:bg-neutral-50 rounded-lg transition-colors"
        :style="{ paddingLeft: `${depth * 16 + 12}px` }"
      >
        <!-- Expand Icon -->
        <svg 
          class="w-4 h-4 text-neutral-400 transition-transform flex-shrink-0" 
          :class="{ 'rotate-90': expanded }"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>
        
        <!-- Folder Icon -->
        <div class="w-8 h-8 bg-amber-100 text-amber-600 rounded-lg flex items-center justify-center flex-shrink-0">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
          </svg>
        </div>
        
        <!-- Title -->
        <div class="flex-1 text-left">
          <h4 class="font-medium text-neutral-800 text-sm">{{ lesson.title }}</h4>
          <p v-if="hasChildren" class="text-xs text-neutral-400">{{ childCount }} item</p>
        </div>
        
        <!-- Module Progress -->
        <div v-if="moduleProgress > 0" class="flex items-center gap-2">
          <div class="w-16 h-1.5 bg-neutral-100 rounded-full overflow-hidden">
            <div class="h-full bg-accent-500 rounded-full" :style="{ width: moduleProgress + '%' }"></div>
          </div>
          <span class="text-xs text-neutral-500">{{ moduleProgress }}%</span>
        </div>
      </button>
    </div>

    <!-- Lesson Row (non-container) -->
    <div 
      v-else
      @click="$emit('select', lesson)"
      class="flex items-center gap-3 p-3 hover:bg-neutral-50 cursor-pointer transition-colors rounded-lg"
      :class="{ 'bg-primary-50': isSelected }"
      :style="{ paddingLeft: `${depth * 16 + 12}px` }"
    >
      <!-- Status Icon / Number -->
      <div 
        class="w-8 h-8 rounded-full flex items-center justify-center flex-shrink-0 text-xs font-semibold"
        :class="getStatusClass()"
      >
        <svg v-if="isCompleted" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
        </svg>
        <span v-else>{{ lessonNumber }}</span>
      </div>
      
      <!-- Content -->
      <div class="flex-1 min-w-0">
        <h4 class="text-sm font-medium text-neutral-800 truncate">{{ lesson.title }}</h4>
        <div class="flex items-center gap-2 mt-0.5">
          <span 
            class="px-1.5 py-0.5 text-xs font-medium rounded"
            :class="getContentTypeBadgeClass()"
          >
            {{ getContentTypeLabel() }}
          </span>
          <span v-if="lesson.video_duration" class="text-xs text-neutral-400">
            {{ formatDuration(lesson.video_duration) }}
          </span>
        </div>
      </div>
      
      <!-- Type Icon -->
      <div class="flex-shrink-0">
        <svg v-if="lesson.content_type === 'video'" class="w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-else-if="lesson.content_type === 'pdf'" class="w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
        </svg>
        <svg v-else-if="lesson.content_type === 'quiz'" class="w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-else class="w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
      </div>
    </div>

    <!-- Children (recursive) -->
    <Transition name="expand">
      <div v-if="expanded && hasChildren" class="children">
        <StudentLessonTree
          v-for="(child, idx) in lesson.children"
          :key="child.id"
          :lesson="child"
          :depth="depth + 1"
          :lesson-number="idx + 1"
          :selected-id="selectedId"
          :completed-ids="completedIds"
          @select="$emit('select', $event)"
        />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
interface Lesson {
  id: string
  title: string
  description?: string
  content_type?: string
  video_url?: string
  video_duration?: number
  is_container?: boolean
  children?: Lesson[]
}

const props = withDefaults(defineProps<{
  lesson: Lesson
  depth?: number
  lessonNumber?: number
  selectedId?: string
  completedIds?: string[]
}>(), {
  depth: 0,
  lessonNumber: 1,
  completedIds: () => []
})

defineEmits<{
  select: [lesson: Lesson]
}>()

const expanded = ref(true)
const hasChildren = computed(() => props.lesson.children && props.lesson.children.length > 0)
const childCount = computed(() => countAllLessons(props.lesson))
const isSelected = computed(() => props.selectedId === props.lesson.id)
const isCompleted = computed(() => props.completedIds.includes(props.lesson.id))

// Calculate module progress
const moduleProgress = computed(() => {
  if (!props.lesson.is_container || !hasChildren.value) return 0
  const allLessonIds = getAllLessonIds(props.lesson)
  if (allLessonIds.length === 0) return 0
  const completedCount = allLessonIds.filter(id => props.completedIds.includes(id)).length
  return Math.round((completedCount / allLessonIds.length) * 100)
})

const toggleExpand = () => {
  expanded.value = !expanded.value
}

// Count all non-container lessons (including nested)
const countAllLessons = (lesson: Lesson): number => {
  let count = 0
  if (lesson.children) {
    for (const child of lesson.children) {
      if (child.is_container) {
        count += countAllLessons(child)
      } else {
        count++
      }
    }
  }
  return count
}

// Get all lesson IDs (non-container)
const getAllLessonIds = (lesson: Lesson): string[] => {
  const ids: string[] = []
  if (lesson.children) {
    for (const child of lesson.children) {
      if (child.is_container) {
        ids.push(...getAllLessonIds(child))
      } else {
        ids.push(child.id)
      }
    }
  }
  return ids
}

const getStatusClass = () => {
  if (isCompleted.value) {
    return 'bg-accent-100 text-accent-600'
  }
  if (isSelected.value) {
    return 'bg-primary-600 text-white'
  }
  return 'bg-neutral-100 text-neutral-600'
}

const getContentTypeBadgeClass = () => {
  switch (props.lesson.content_type) {
    case 'video': return 'bg-blue-100 text-blue-700'
    case 'pdf': return 'bg-red-100 text-red-700'
    case 'quiz': return 'bg-purple-100 text-purple-700'
    default: return 'bg-neutral-100 text-neutral-700'
  }
}

const getContentTypeLabel = () => {
  switch (props.lesson.content_type) {
    case 'video': return 'Video'
    case 'pdf': return 'PDF'
    case 'quiz': return 'Kuis'
    case 'text': return 'Teks'
    default: return 'Materi'
  }
}

const formatDuration = (seconds: number): string => {
  if (!seconds) return ''
  const mins = Math.floor(seconds / 60)
  return `${mins} menit`
}
</script>

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
  overflow: hidden;
}
.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
