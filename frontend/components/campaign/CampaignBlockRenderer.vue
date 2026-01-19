<script setup lang="ts">
import type { Block } from '~/utils/templates'
import type { CampaignStyles } from '~/composables/useCampaignStyles'

// Import all block components
import CampaignHero from './blocks/CampaignHero.vue'
import CampaignPricing from './blocks/CampaignPricing.vue'
import CampaignBenefits from './blocks/CampaignBenefits.vue'
import CampaignTestimonials from './blocks/CampaignTestimonials.vue'
import CampaignCountdown from './blocks/CampaignCountdown.vue'
import CampaignFaq from './blocks/CampaignFaq.vue'
import CampaignInstructor from './blocks/CampaignInstructor.vue'
import CampaignVideo from './blocks/CampaignVideo.vue'
import CampaignTrust from './blocks/CampaignTrust.vue'
import CampaignCtaBanner from './blocks/CampaignCtaBanner.vue'
import CampaignStatistics from './blocks/CampaignStatistics.vue'
import CampaignBonus from './blocks/CampaignBonus.vue'
import CampaignCurriculum from './blocks/CampaignCurriculum.vue'
import CampaignGallery from './blocks/CampaignGallery.vue'
import CampaignComparison from './blocks/CampaignComparison.vue'

interface Props {
  block: Block
  blockIndex: number
  styles: CampaignStyles
  buttonStyle: Record<string, any>
  countdown?: { days: number; hours: number; minutes: number; seconds: number } | null
  displayPrice?: number
  coursePrice?: number
  isFree?: boolean
  isMobileLayout?: boolean
  instructorData?: any
  courseCurriculum?: any[]
  courseImages?: string[]
  preview?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  preview: false
})

const emit = defineEmits<{
  (e: 'buy'): void
  (e: 'scrollTo', target: string): void
}>()

// Map block types to components
const blockComponentMap: Record<string, any> = {
  hero: CampaignHero,
  pricing: CampaignPricing,
  benefits: CampaignBenefits,
  testimonials: CampaignTestimonials,
  countdown: CampaignCountdown,
  faq: CampaignFaq,
  instructor: CampaignInstructor,
  video: CampaignVideo,
  trust: CampaignTrust,
  cta_banner: CampaignCtaBanner,
  statistics: CampaignStatistics,
  bonus: CampaignBonus,
  curriculum: CampaignCurriculum,
  gallery: CampaignGallery,
  comparison: CampaignComparison
}

// Get component for block type
const blockComponent = computed(() => {
  return blockComponentMap[props.block.type] || null
})

// Animation delay based on block index
const animationDelay = computed(() => props.blockIndex * 100)

// Handle buy event
const handleBuy = () => {
  emit('buy')
}

// Handle scroll event
const handleScrollTo = (target: string) => {
  emit('scrollTo', target)
}
</script>

<template>
  <div :id="`block-${block.id}`">
    <!-- Custom CSS for this block -->
    <component 
      v-if="block.data.customCSS" 
      :is="'style'" 
      v-text="`#block-${block.id} { ${block.data.customCSS} }`"
    />
    
    <!-- Hero Block -->
    <CampaignHero 
      v-if="block.type === 'hero'"
      :block="block"
      :styles="styles"
      :button-style="buttonStyle"
      :countdown="countdown"
      :display-price="displayPrice"
      :course-price="coursePrice"
      :is-free="isFree"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
      :preview="preview"
      @buy="handleBuy"
      @scroll-to="handleScrollTo"
    />
    
    <!-- Countdown Block -->
    <CampaignCountdown 
      v-else-if="block.type === 'countdown'"
      :block="block"
      :countdown="countdown || { days: 0, hours: 0, minutes: 0, seconds: 0 }"
      :styles="styles"
      :label="block.data.label"
      :animation-delay="animationDelay"
    />
    
    <!-- Benefits Block -->
    <CampaignBenefits 
      v-else-if="block.type === 'benefits'"
      :block="block"
      :styles="styles"
      :block-index="blockIndex"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Pricing Block -->
    <CampaignPricing 
      v-else-if="block.type === 'pricing'"
      :block="block"
      :styles="styles"
      :button-style="buttonStyle"
      :display-price="displayPrice"
      :course-price="coursePrice"
      :is-free="isFree"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
      :preview="preview"
      @buy="handleBuy"
    />
    
    <!-- Testimonials Block -->
    <CampaignTestimonials 
      v-else-if="block.type === 'testimonials'"
      :block="block"
      :styles="styles"
      :block-index="blockIndex"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- FAQ Block -->
    <CampaignFaq 
      v-else-if="block.type === 'faq'"
      :block="block"
      :styles="styles"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Instructor Block -->
    <CampaignInstructor 
      v-else-if="block.type === 'instructor'"
      :block="block"
      :styles="styles"
      :instructor-data="instructorData"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Video Block -->
    <CampaignVideo 
      v-else-if="block.type === 'video'"
      :block="block"
      :styles="styles"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Trust Block -->
    <CampaignTrust 
      v-else-if="block.type === 'trust'"
      :block="block"
      :styles="styles"
      :animation-delay="animationDelay"
    />
    
    <!-- CTA Banner Block -->
    <CampaignCtaBanner 
      v-else-if="block.type === 'cta_banner'"
      :block="block"
      :styles="styles"
      :button-style="buttonStyle"
      :animation-delay="animationDelay"
      @buy="handleBuy"
    />
    
    <!-- Statistics Block -->
    <CampaignStatistics 
      v-else-if="block.type === 'statistics'"
      :block="block"
      :styles="styles"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Bonus Block -->
    <CampaignBonus 
      v-else-if="block.type === 'bonus'"
      :block="block"
      :styles="styles"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Curriculum Block -->
    <CampaignCurriculum 
      v-else-if="block.type === 'curriculum'"
      :block="block"
      :styles="styles"
      :course-curriculum="courseCurriculum"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Gallery Block -->
    <CampaignGallery 
      v-else-if="block.type === 'gallery'"
      :block="block"
      :styles="styles"
      :course-images="courseImages"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Comparison Block -->
    <CampaignComparison 
      v-else-if="block.type === 'comparison'"
      :block="block"
      :styles="styles"
      :is-mobile-layout="isMobileLayout"
      :animation-delay="animationDelay"
    />
    
    <!-- Fallback for unknown block types -->
    <div 
      v-else-if="!['social_proof', 'floating_chat'].includes(block.type)"
      class="py-8 px-4 bg-neutral-100 text-center text-neutral-500"
    >
      <p class="text-sm">Block type "{{ block.type }}" tidak dikenali</p>
    </div>
  </div>
</template>
