<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

interface Props {
  block: Block
  styles: CampaignStyles
  instructorData?: {
    name?: string
    avatar_url?: string
    bio?: string
  } | null
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

// Get instructor name
const instructorName = computed(() => {
  return props.block.data.name || props.instructorData?.name || 'Nama Instruktur'
})

// Get instructor image
const instructorImage = computed(() => {
  return props.block.data.profile_image || props.instructorData?.avatar_url
})

// Get instructor bio
const instructorBio = computed(() => {
  return props.block.data.bio || props.instructorData?.bio || ''
})

// Credentials list
const credentials = computed(() => {
  return props.block.data.credentials || []
})
</script>

<template>
  <section 
    class="py-12 sm:py-20 animate-fade-in-up" 
    :style="{
      backgroundColor: block.data.backgroundColor || '#fafafa',
      animationDelay: `${animationDelay}ms`
    }"
  >
    <div :class="containerClass">
      <h2 
        class="text-2xl sm:text-3xl md:text-4xl font-bold text-center mb-8 sm:mb-12" 
        :style="{color: block.data.titleColor || styles.textPrimaryColor, fontFamily: styles.fontFamilyHeading}"
      >
        {{ block.data.title || 'Tentang Instruktur' }}
      </h2>
      
      <div class="flex flex-col md:flex-row items-center md:items-start gap-8">
        <!-- Profile Image -->
        <div class="flex-shrink-0">
          <div 
            v-if="instructorImage"
            class="w-32 h-32 sm:w-40 sm:h-40 rounded-2xl overflow-hidden shadow-lg"
          >
            <img 
              :src="instructorImage" 
              :alt="instructorName"
              class="w-full h-full object-cover"
            />
          </div>
          <div 
            v-else 
            class="w-32 h-32 sm:w-40 sm:h-40 rounded-2xl bg-neutral-100 flex items-center justify-center text-neutral-400"
          >
            <svg class="w-16 h-16 sm:w-20 sm:h-20" fill="currentColor" viewBox="0 0 24 24">
              <path fill-rule="evenodd" d="M12 2a5 5 0 100 10 5 5 0 000-10zM6 14a6 6 0 016-6h.06a6 6 0 015.94 6v.203a6.002 6.002 0 01-5.308 5.966L12 20.203l-.693-.035A6.002 6.002 0 016 14.204V14z" clip-rule="evenodd" />
            </svg>
          </div>
        </div>
        
        <!-- Info -->
        <div class="flex-1 text-center md:text-left">
          <h3 
            class="text-xl sm:text-2xl font-bold mb-3" 
            :style="{color: block.data.titleColor || styles.textPrimaryColor}"
          >
            {{ instructorName }}
          </h3>
          
          <!-- Credentials/Tags -->
          <div v-if="credentials.length" class="flex flex-wrap gap-2 mb-4 justify-center md:justify-start">
            <span 
              v-for="(cred, idx) in credentials" 
              :key="idx"
              class="px-3 py-1 text-xs font-medium rounded-full"
              :style="{backgroundColor: (block.data.accentColor || styles.primaryColor) + '15', color: block.data.accentColor || styles.primaryColor}"
            >
              {{ cred }}
            </span>
          </div>
          
          <!-- Bio -->
          <p 
            class="text-sm sm:text-base leading-relaxed whitespace-pre-line" 
            :style="{color: block.data.textColor || styles.textSecondaryColor}"
          >
            {{ instructorBio }}
          </p>
          
          <!-- Social Links (if any) -->
          <div v-if="block.data.social_links" class="flex gap-3 mt-4 justify-center md:justify-start">
            <a 
              v-for="(link, platform) in block.data.social_links" 
              :key="platform"
              :href="link"
              target="_blank"
              class="w-10 h-10 rounded-full bg-neutral-100 hover:bg-neutral-200 flex items-center justify-center transition-colors"
            >
              <span v-if="platform === 'instagram'">📷</span>
              <span v-else-if="platform === 'youtube'">▶️</span>
              <span v-else-if="platform === 'linkedin'">💼</span>
              <span v-else>🔗</span>
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
