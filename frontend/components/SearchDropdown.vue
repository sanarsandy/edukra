<template>
  <div class="relative" ref="searchRef">
    <!-- Search Input -->
    <div class="relative">
      <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
      </svg>
      <input 
        v-model="query"
        type="text" 
        :placeholder="placeholder"
        class="w-full pl-10 pr-4 py-2.5 border-0 rounded-lg focus:outline-none focus:ring-2 transition-all text-sm"
        :class="variant === 'light' ? 'bg-neutral-100 focus:ring-primary-500 focus:bg-white placeholder-neutral-500' : 'bg-neutral-800 focus:ring-admin-500 focus:bg-neutral-700 placeholder-neutral-400 text-white'"
        @input="handleSearch"
        @focus="showDropdown = true"
      />
      <!-- Loading Spinner -->
      <svg v-if="loading" class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400 animate-spin" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
      </svg>
    </div>
    
    <!-- Dropdown Results -->
    <Transition name="dropdown">
      <div v-if="showDropdown && query.length >= 2" class="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-lg border border-neutral-200 py-2 z-50 max-h-80 overflow-y-auto">
        <!-- Loading State -->
        <div v-if="loading" class="px-4 py-8 text-center">
          <svg class="w-6 h-6 mx-auto text-neutral-400 animate-spin mb-2" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          <p class="text-sm text-neutral-500">Mencari...</p>
        </div>
        
        <!-- Results -->
        <template v-else-if="results.length > 0">
          <div v-for="result in results" :key="result.id" class="px-4 py-3 hover:bg-neutral-50 cursor-pointer transition-colors" @click="selectResult(result)">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0" :class="result.type === 'course' ? 'bg-primary-100' : 'bg-accent-100'">
                <svg v-if="result.type === 'course'" class="w-5 h-5 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                </svg>
                <svg v-else class="w-5 h-5 text-accent-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-neutral-900 truncate">{{ result.title }}</p>
                <p class="text-xs text-neutral-500 truncate">{{ result.subtitle }}</p>
              </div>
            </div>
          </div>
        </template>
        
        <!-- No Results -->
        <div v-else class="px-4 py-8 text-center">
          <svg class="w-12 h-12 mx-auto text-neutral-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <p class="text-sm text-neutral-500">Tidak ditemukan hasil untuk "{{ query }}"</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
interface SearchResult {
  id: string | number
  title: string
  subtitle: string
  type: 'course' | 'lesson' | 'user'
  url?: string
}

const props = defineProps<{
  placeholder?: string
  variant?: 'light' | 'dark'
}>()

const emit = defineEmits(['select', 'search'])

const searchRef = ref<HTMLElement | null>(null)
const query = ref('')
const showDropdown = ref(false)
const loading = ref(false)
const results = ref<SearchResult[]>([])
let debounceTimer: ReturnType<typeof setTimeout>

const handleSearch = () => {
  clearTimeout(debounceTimer)
  if (query.value.length < 2) {
    results.value = []
    return
  }
  loading.value = true
  debounceTimer = setTimeout(() => {
    // Simulated search - replace with actual API call
    results.value = [
      { id: 1, title: 'Fullstack JavaScript', subtitle: 'React, Node.js, MongoDB', type: 'course' },
      { id: 2, title: 'UI/UX Design Fundamentals', subtitle: 'Figma, Design System', type: 'course' },
      { id: 3, title: 'Pengenalan HTML & CSS', subtitle: 'Modul Dasar Web', type: 'lesson' },
    ].filter(r => r.title.toLowerCase().includes(query.value.toLowerCase()))
    loading.value = false
    emit('search', query.value)
  }, 300)
}

const selectResult = (result: SearchResult) => {
  emit('select', result)
  showDropdown.value = false
  query.value = ''
}

onMounted(() => {
  document.addEventListener('click', (e) => {
    if (searchRef.value && !searchRef.value.contains(e.target as Node)) {
      showDropdown.value = false
    }
  })
})
</script>
