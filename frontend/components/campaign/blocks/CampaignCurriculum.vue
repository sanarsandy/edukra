<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  courseCurriculum?: any[] | null
  isMobileLayout?: boolean
  animationDelay?: number
}

const props = withDefaults(defineProps<Props>(), {
  animationDelay: 0
})

// Container class
const containerClass = computed(() => {
  return props.isMobileLayout 
    ? 'w-full max-w-[480px] mx-auto px-4' 
    : 'w-full max-w-4xl mx-auto px-4'
})

// Modules from block data or course curriculum
const modules = computed(() => {
  return props.block.data.modules || props.courseCurriculum || []
})

// Track expanded modules
const expandedModules = ref<number[]>([0])

const toggleModule = (idx: number) => {
  const index = expandedModules.value.indexOf(idx)
  if (index > -1) {
    expandedModules.value.splice(index, 1)
  } else {
    expandedModules.value.push(idx)
  }
}

// Total lessons count
const totalLessons = computed(() => {
  return modules.value.reduce((acc: number, mod: any) => {
    return acc + (mod.lessons_count || mod.lessons?.length || 0)
  }, 0)
})
</script>

<template>
  <section 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#ffffff',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <!-- Header -->
      <div class="text-center mb-8 sm:mb-10">
        <h2 
          class="text-2xl sm:text-3xl font-bold mb-2" 
          :style="{color: styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
        >
          {{ block.data.title || '📚 Materi Kursus' }}
        </h2>
        <p :style="{color: styles.textSecondaryColor}">
          {{ block.data.subtitle || `${modules.length} Modul • ${totalLessons} Pelajaran` }}
        </p>
      </div>
      
      <!-- Modules List -->
      <div class="space-y-3">
        <div 
          v-for="(module, idx) in modules" 
          :key="idx"
          class="border rounded-xl overflow-hidden transition-all"
          :class="expandedModules.includes(idx) ? 'border-neutral-300 shadow-sm' : 'border-neutral-200'"
        >
          <!-- Module Header -->
          <button 
            @click="toggleModule(idx)"
            class="w-full p-4 flex items-center gap-4 text-left transition-colors"
            :class="expandedModules.includes(idx) ? 'bg-neutral-50' : 'bg-white hover:bg-neutral-50'"
          >
            <!-- Module Number -->
            <div 
              class="w-8 h-8 rounded-lg flex items-center justify-center text-sm font-bold flex-shrink-0"
              :style="{backgroundColor: styles.primaryColor + '15', color: styles.primaryColor}"
            >
              {{ idx + 1 }}
            </div>
            
            <!-- Module Title -->
            <div class="flex-1 min-w-0">
              <h4 class="font-semibold truncate" :style="{color: styles.textPrimaryColor}">
                {{ module.title }}
              </h4>
              <p class="text-xs text-neutral-500 mt-0.5">
                {{ module.lessons_count || module.lessons?.length || 0 }} Pelajaran
              </p>
            </div>
            
            <!-- Toggle Icon -->
            <svg 
              class="w-5 h-5 flex-shrink-0 transition-transform duration-200" 
              :class="expandedModules.includes(idx) ? 'rotate-180' : ''"
              :style="{color: styles.textSecondaryColor}"
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
          </button>
          
          <!-- Lessons (Expanded) -->
          <div 
            class="overflow-hidden transition-all duration-200"
            :class="expandedModules.includes(idx) ? 'max-h-[500px]' : 'max-h-0'"
          >
            <div class="p-4 pt-0 space-y-2">
              <div 
                v-for="(lesson, lidx) in (module.lessons || [])" 
                :key="lidx"
                class="flex items-center gap-3 p-2 rounded-lg hover:bg-neutral-50"
              >
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <span class="text-sm" :style="{color: styles.textSecondaryColor}">
                  {{ lesson.title }}
                </span>
                <span v-if="lesson.duration" class="text-xs text-neutral-400 ml-auto">
                  {{ lesson.duration }}
                </span>
              </div>
              <div v-if="!module.lessons?.length" class="text-center text-neutral-400 text-sm py-2">
                {{ module.lessons_count || 0 }} pelajaran tersedia
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
