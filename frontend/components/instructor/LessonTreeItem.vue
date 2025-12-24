<template>
  <div class="lesson-tree-item">
    <!-- Lesson Item Row -->
    <div 
      class="flex items-center gap-3 p-3 hover:bg-neutral-50 transition-colors rounded-lg group"
      :style="{ paddingLeft: `${depth * 24 + 12}px` }"
    >
      <!-- Expand/Collapse Button (for containers) -->
      <button 
        v-if="lesson.is_container && hasChildren"
        @click="toggleExpand"
        class="w-6 h-6 flex items-center justify-center text-neutral-400 hover:text-neutral-600 transition-colors"
      >
        <svg 
          class="w-4 h-4 transition-transform" 
          :class="{ 'rotate-90': expanded }"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
        </svg>
      </button>
      <div v-else class="w-6"></div>

      <!-- Icon -->
      <div 
        class="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0"
        :class="getIconClass()"
      >
        <!-- Container/Folder Icon -->
        <svg v-if="lesson.is_container" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
        </svg>
        <!-- Video Icon -->
        <svg v-else-if="lesson.content_type === 'video'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <!-- PDF Icon -->
        <svg v-else-if="lesson.content_type === 'pdf'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
        </svg>
        <!-- Quiz Icon -->
        <svg v-else-if="lesson.content_type === 'quiz'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <!-- Text Icon -->
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
      </div>

      <!-- Title & Badge -->
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-2">
          <h4 class="font-medium text-neutral-900 truncate text-sm">{{ lesson.title }}</h4>
          <span 
            v-if="lesson.is_container" 
            class="px-1.5 py-0.5 bg-amber-100 text-amber-700 rounded text-xs font-medium flex-shrink-0"
          >
            Modul
          </span>
          <span 
            v-else-if="lesson.content_type"
            class="px-1.5 py-0.5 rounded text-xs font-medium flex-shrink-0"
            :class="getContentTypeBadgeClass()"
          >
            {{ getContentTypeLabel() }}
          </span>
          <span 
            v-if="lesson.is_preview" 
            class="px-1.5 py-0.5 bg-green-100 text-green-700 rounded text-xs font-medium flex-shrink-0"
          >
            Preview
          </span>
        </div>
        <p v-if="lesson.description" class="text-xs text-neutral-500 truncate mt-0.5">
          {{ lesson.description }}
        </p>
      </div>

      <!-- Children Count (for containers) -->
      <span v-if="lesson.is_container && hasChildren" class="text-xs text-neutral-400">
        {{ lesson.children?.length }} item
      </span>

      <!-- Actions -->
      <div class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity">
        <!-- Add Submodule (for containers) -->
        <button 
          v-if="lesson.is_container"
          @click="$emit('add-submodule', lesson)"
          class="p-1.5 text-neutral-400 hover:text-amber-600 hover:bg-amber-50 rounded-lg transition-colors"
          title="Tambah submodul"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
          </svg>
        </button>
        <!-- Add Quiz (for containers) -->
        <button 
          v-if="lesson.is_container"
          @click="$emit('add-quiz', lesson)"
          class="p-1.5 text-neutral-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-colors"
          title="Tambah kuis"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </button>
        <!-- Add Material (for containers) -->
        <button 
          v-if="lesson.is_container"
          @click="$emit('add-child', lesson)"
          class="p-1.5 text-neutral-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-colors"
          title="Tambah materi"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
        </button>
        <!-- View -->
        <button 
          v-if="!lesson.is_container"
          @click="$emit('view', lesson)"
          class="p-1.5 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors"
          title="Lihat"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
          </svg>
        </button>
        <!-- Manage Quiz (for quiz items) -->
        <button 
          v-if="lesson.content_type === 'quiz'"
          @click="$emit('manage-quiz', lesson)"
          class="p-1.5 text-neutral-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-colors"
          title="Kelola Kuis"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
        </button>
        <!-- Edit -->
        <button 
          @click="$emit('edit', lesson)"
          class="p-1.5 text-neutral-400 hover:text-admin-600 hover:bg-admin-50 rounded-lg transition-colors"
          title="Edit"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
          </svg>
        </button>
        <!-- Delete -->
        <button 
          @click="$emit('delete', lesson)"
          class="p-1.5 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
          title="Hapus"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Children (recursive) -->
    <Transition name="expand">
      <div v-if="expanded && hasChildren" class="children">
        <InstructorLessonTreeItem
          v-for="child in lesson.children"
          :key="child.id"
          :lesson="child"
          :depth="depth + 1"
          @view="$emit('view', $event)"
          @edit="$emit('edit', $event)"
          @delete="$emit('delete', $event)"
          @add-child="$emit('add-child', $event)"
          @add-submodule="$emit('add-submodule', $event)"
          @add-quiz="$emit('add-quiz', $event)"
          @manage-quiz="$emit('manage-quiz', $event)"
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
  is_container?: boolean
  is_preview?: boolean
  children?: Lesson[]
}

const props = defineProps<{
  lesson: Lesson
  depth?: number
}>()

defineEmits<{
  view: [lesson: Lesson]
  edit: [lesson: Lesson]
  delete: [lesson: Lesson]
  'add-child': [parent: Lesson]
  'add-submodule': [parent: Lesson]
  'add-quiz': [parent: Lesson]
  'manage-quiz': [lesson: Lesson]
}>()

const depth = computed(() => props.depth ?? 0)
const expanded = ref(true)
const hasChildren = computed(() => props.lesson.children && props.lesson.children.length > 0)

const toggleExpand = () => {
  expanded.value = !expanded.value
}

const getIconClass = () => {
  if (props.lesson.is_container) {
    return 'bg-amber-100 text-amber-600'
  }
  switch (props.lesson.content_type) {
    case 'video': return 'bg-blue-100 text-blue-600'
    case 'pdf': return 'bg-red-100 text-red-600'
    case 'quiz': return 'bg-purple-100 text-purple-600'
    default: return 'bg-neutral-100 text-neutral-600'
  }
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

