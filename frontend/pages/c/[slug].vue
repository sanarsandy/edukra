<template>
  <!-- Main wrapper -->
  <div class="min-h-screen bg-neutral-900">
    <!-- SSR Placeholder - shows on server, replaced by ClientOnly on client -->
    <ClientOnly>
      <!-- Loading State -->
      <div v-if="pending" class="min-h-screen flex items-center justify-center bg-neutral-900">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
      </div>

      <!-- Not Found -->
      <div v-else-if="!campaign" class="min-h-screen flex items-center justify-center bg-neutral-900 text-white">
        <div class="text-center">
          <h1 class="text-4xl font-bold mb-4">Halaman Tidak Ditemukan</h1>
          <p class="text-neutral-400 mb-6">Campaign yang Anda cari tidak tersedia</p>
          <NuxtLink to="/" class="px-6 py-3 bg-primary-600 rounded-lg font-medium hover:bg-primary-700">Kembali ke Home</NuxtLink>
        </div>
      </div>

      <!-- Campaign Content -->
      <div v-else class="min-h-screen" :style="pageStyle">
        <!-- Background Image Layer -->
        <div v-if="campaignStyles.backgroundImage" class="fixed inset-0 z-0 bg-cover bg-center bg-no-repeat" 
             :style="{backgroundImage: campaignStyles.backgroundImage}"></div>
        
        <!-- Gradient Layer (if needed) -->
        <div v-if="campaignStyles.hasGradient && !campaignStyles.backgroundImage" class="fixed inset-0 z-0 bg-gradient-to-br from-primary-500/10 to-accent-500/10 pointer-events-none"></div>

        <!-- Main Content Container -->
        <div class="relative z-10">
          <!-- Global Custom CSS (page-wide) -->
          <component :is="'style'" v-if="campaign?.globalCustomCSS">
            {{ campaign.globalCustomCSS }}
          </component>

          <!-- Render all enabled blocks using component -->
          <CampaignBlockRenderer
            v-for="(block, blockIndex) in enabledBlocks"
            :key="block.id"
            :block="block"
            :block-index="blockIndex"
            :styles="campaignStyles"
            :button-style="buttonStyle"
            :countdown="countdown"
            :display-price="displayPrice"
            :course-price="coursePrice"
            :is-free="isFree"
            :is-mobile-layout="isMobileLayout"
            :instructor-data="instructorData"
            :campaign-id="campaign?.id"
            :end-date="campaign?.end_date"
            @buy="handleBuy"
            @scroll-to="scrollTo"
          />

          <!-- Social Proof Component -->
          <CampaignSocialProof
            v-if="hasSocialProof"
            :block="socialProofBlock!"
            :styles="campaignStyles"
          />

          <!-- Floating WhatsApp -->
          <CampaignFloatingChat
            v-if="campaign?.whatsapp_number"
            :whatsapp-number="campaign.whatsapp_number"
            :styles="campaignStyles"
            :message="`Halo, saya tertarik dengan ${campaign?.title || 'kursus ini'}!`"
          />
        </div>
      </div>

      <!-- SSR Fallback (shown while ClientOnly hydrates) -->
      <template #fallback>
        <div class="min-h-screen flex items-center justify-center bg-neutral-900">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-primary-500 border-t-transparent"></div>
        </div>
      </template>
    </ClientOnly>

    <!-- Mobile Sticky CTA -->
    <ClientOnly>
      <div v-if="campaign" class="fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-xl border-t border-neutral-200 p-3 sm:p-4 md:hidden z-50 safe-area-pb">
        <div class="flex items-center gap-3">
          <div class="flex-1 min-w-0">
            <div class="text-sm text-neutral-500 truncate">{{ campaign.course?.title || campaign.webinar?.title }}</div>
            <div v-if="isFree" class="font-bold text-lg text-green-600">GRATIS</div>
            <div v-else class="font-bold text-lg" :style="{color: campaignStyles.primaryColor}">Rp {{ formatPrice(displayPrice) }}</div>
          </div>
          <button @click="handleBuy" class="px-6 py-3 text-white font-bold rounded-xl flex-shrink-0 min-h-[48px]" :style="buttonStyle">
            {{ isFree ? 'Daftar Gratis' : 'Beli Sekarang' }}
          </button>
        </div>
      </div>
    </ClientOnly>

    <!-- Checkout Modal -->
    <CampaignCheckoutModal
      :show="showCheckoutModal"
      :campaign-id="campaign?.id || ''"
      :course-name="campaign?.course?.title || campaign?.webinar?.title || ''"
      :course-price="campaign?.course?.price || campaign?.webinar?.price || 0"
      :discount-price="campaign?.course?.discount_price"
      :is-free="isFree"
      @close="showCheckoutModal = false"
      @success="handleCheckoutSuccess"
    />

    <!-- Success Modal -->
    <Teleport to="body">
      <Transition name="fade">
        <div v-if="showSuccessModal" class="fixed inset-0 bg-black/60 flex items-center justify-center z-[60] p-4" @click.self="closeSuccessModal">
          <Transition name="scale">
            <div v-if="showSuccessModal" class="bg-white rounded-3xl w-full max-w-md p-8 text-center shadow-2xl">
              <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-green-100 flex items-center justify-center">
                <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"/>
                </svg>
              </div>
              <h3 class="text-2xl font-bold text-neutral-900 mb-2">🎉 Selamat!</h3>
              <p class="text-neutral-600 mb-6">{{ successMessage }}</p>
              <button @click="goToCourse" class="w-full py-4 bg-green-600 text-white font-bold rounded-xl hover:bg-green-700 transition-all transform hover:scale-[1.02] active:scale-[0.98]">
                Lihat Kursus Saya
              </button>
              <button @click="closeSuccessModal" class="mt-4 text-sm text-neutral-500 hover:text-neutral-700 transition-colors">
                Tutup
              </button>
            </div>
          </Transition>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import CampaignBlockRenderer from '~/components/campaign/CampaignBlockRenderer.vue'
import CampaignSocialProof from '~/components/campaign/blocks/CampaignSocialProof.vue'
import CampaignFloatingChat from '~/components/campaign/blocks/CampaignFloatingChat.vue'
import { useCampaignStyles, useCampaignCountdown, type CampaignStyles } from '~/composables/useCampaignStyles'
import type { Block } from '~/utils/templates'

interface Campaign {
  id: string
  slug: string
  title: string
  meta_description?: string
  course_id?: string
  webinar_id?: string
  end_date?: string
  is_free_webinar?: boolean | null
  blocks?: Block[]
  styles?: CampaignStyles | string
  globalCustomCSS?: string
  whatsapp_number?: string
  og_image_url?: string
  gtm_id?: string
  facebook_pixel_id?: string
  course?: {
    id: string
    slug: string
    title: string
    price: number
    discount_price?: number
    instructor?: {
      full_name: string
      bio?: string
      avatar_url?: string
    }
  }
  webinar?: {
    id: string
    title: string
    price: number
    cover_image_url?: string
  }
}

definePageMeta({ 
  layout: false,
  validate: async (route) => {
    // Explicitly exclude static pages within this directory to prevent conflict
    if (route.params.slug === 'kelas-kepabeanan') return false
    return true
  }
})

const route = useRoute()
const router = useRouter()
const slug = computed(() => route.params.slug as string)

const config = useRuntimeConfig()
const apiBase = import.meta.server 
  ? (config.apiInternal || 'http://api:8080')
  : (config.public.apiBase || 'http://localhost:8080')

const { data: campaign, pending, refresh } = await useFetch<Campaign>(`/api/c/${slug.value}`, {
  baseURL: apiBase,
  key: `campaign-${slug.value}`,
  watch: [slug],
  // Disable SSR caching to always get fresh data
  getCachedData: () => null
})

onMounted(() => {
  // Force refresh on client to ensure fresh data
  refresh()
})

// SEO Head
useHead({
  title: computed(() => campaign.value?.title || 'Campaign'),
  meta: [
    { name: 'description', content: computed(() => campaign.value?.meta_description || '') },
    { property: 'og:title', content: computed(() => campaign.value?.title || '') },
    { property: 'og:description', content: computed(() => campaign.value?.meta_description || '') },
    { property: 'og:image', content: computed(() => campaign.value?.og_image_url || campaign.value?.course?.thumbnail_url || '') },
    { name: 'twitter:card', content: 'summary_large_image' }
  ],
  script: [
    computed(() => campaign.value?.gtm_id ? {
      children: `(function(w,d,s,l,i){w[l]=w[l]||[];w[l].push({'gtm.start':
new Date().getTime(),event:'gtm.js'});var f=d.getElementsByTagName(s)[0],
j=d.createElement(s),dl=l!='dataLayer'?'&l='+l:'';j.async=true;j.src=
'https://www.googletagmanager.com/gtm.js?id='+i+dl;f.parentNode.insertBefore(j,f);
})(window,document,'script','dataLayer','${campaign.value.gtm_id}');`,
      type: 'text/javascript'
    } : undefined),
    computed(() => campaign.value?.facebook_pixel_id ? {
      children: `!function(f,b,e,v,n,t,s)
{if(f.fbq)return;n=f.fbq=function(){n.callMethod?
n.callMethod.apply(n,arguments):n.queue.push(arguments)};
if(!f._fbq)f._fbq=n;n.push=n;n.loaded=!0;n.version='2.0';
n.queue=[];t=b.createElement(e);t.async=!0;
t.src=v;s=b.getElementsByTagName(e)[0];
s.parentNode.insertBefore(t,s)}(window, document,'script',
'https://connect.facebook.net/en_US/fbevents.js');
fbq('init', '${campaign.value.facebook_pixel_id}');
fbq('track', 'PageView');`,
      type: 'text/javascript'
    } : undefined)
  ]
})

// Use shared composable for styles
const campaignRef = computed(() => campaign.value as any)
const { 
  campaignStyles, 
  isMobileLayout, 
  enabledBlocks, 
  pageStyle, 
  buttonStyle 
} = useCampaignStyles(campaignRef)

// Countdown
const discountEndDate = computed(() => campaign.value?.end_date)
const { countdown } = useCampaignCountdown(discountEndDate)

// Price calculations
const coursePrice = computed(() => campaign.value?.course?.price || campaign.value?.webinar?.price || 0)
const displayPrice = computed(() => campaign.value?.course?.discount_price || coursePrice.value)
const isFree = computed(() => {
  if (campaign.value?.is_free_webinar !== null && campaign.value?.is_free_webinar !== undefined) {
    return campaign.value.is_free_webinar
  }
  return coursePrice.value === 0
})

// Instructor data
const instructorData = computed(() => {
  return campaign.value?.course?.instructor || null
})

// Social proof block
const hasSocialProof = computed(() => enabledBlocks.value.some((b: Block) => b.type === 'social_proof'))
const socialProofBlock = computed(() => enabledBlocks.value.find((b: Block) => b.type === 'social_proof'))

// Format price
const formatPrice = (price: number) => new Intl.NumberFormat('id-ID').format(price || 0)

// Track view
const trackView = async () => {
  if (!campaign.value?.id) return
  try {
    await $fetch(`/api/c/${campaign.value.id}/track`, {
      method: 'POST',
      baseURL: config.public.apiBase || 'http://localhost:8080',
      body: { event_type: 'view' }
    })
  } catch (e) {}
}

onMounted(() => {
  trackView()
})

// Checkout modal
const showCheckoutModal = ref(false)
const showSuccessModal = ref(false)
const successMessage = ref('')

const handleBuy = async () => {
  if (!campaign.value?.course && !campaign.value?.webinar) {
    alert('Kursus/Webinar belum terhubung')
    return
  }
  
  try {
    await $fetch(`/api/c/${campaign.value.id}/track`, {
      method: 'POST',
      baseURL: config.public.apiBase || 'http://localhost:8080',
      body: { event_type: 'click' }
    })
  } catch (e) {}

  showCheckoutModal.value = true
}

const handleCheckoutSuccess = (result: { isFree: boolean, message?: string }) => {
  showCheckoutModal.value = false
  if (result.isFree) {
    successMessage.value = result.message || 'Anda berhasil terdaftar. Informasi webinar akan dikirim via WhatsApp.'
    showSuccessModal.value = true
  }
}

const closeSuccessModal = () => {
  showSuccessModal.value = false
}

const goToCourse = () => {
  showSuccessModal.value = false
  router.push(`/dashboard/courses/${campaign.value?.course?.id}`)
}

const scrollTo = (selector: string) => {
  const el = document.querySelector(selector)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<style>
/* Hide scrollbar for testimonials carousel */
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

/* Safe area padding for mobile bottom sticky */
.safe-area-pb {
  padding-bottom: max(env(safe-area-inset-bottom), 12px);
}

/* Smooth scroll snap */
.snap-x {
  scroll-snap-type: x mandatory;
}
.snap-center {
  scroll-snap-align: center;
}

/* Tabular nums for countdown */
.tabular-nums {
  font-variant-numeric: tabular-nums;
}

/* Modal transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.scale-enter-active,
.scale-leave-active {
  transition: all 0.3s ease;
}
.scale-enter-from,
.scale-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
