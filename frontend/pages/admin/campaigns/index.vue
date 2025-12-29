<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Campaign Landing Pages</h1>
        <p class="text-neutral-500 mt-1">Buat halaman promo dengan drag & drop untuk promosi kursus Anda</p>
      </div>
      <NuxtLink 
        to="/admin/campaigns/create" 
        class="inline-flex items-center justify-center px-4 py-2.5 bg-admin-600 hover:bg-admin-700 text-white text-sm font-medium rounded-lg transition-colors shadow-sm"
      >
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 4.5v15m7.5-7.5h-15"/></svg>
        Buat Campaign Baru
      </NuxtLink>
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="w-8 h-8 border-4 border-admin-200 border-t-admin-600 rounded-full animate-spin"></div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg">
      <p class="font-medium">Gagal memuat data campaigns</p>
      <button @click="fetchCampaigns" class="text-sm underline mt-1">Coba lagi</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="campaigns.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <div class="w-20 h-20 bg-gradient-to-br from-admin-100 to-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-10 h-10 text-admin-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z"/></svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada Campaign</h3>
      <p class="text-neutral-500 max-w-sm mx-auto mb-6">Mulai promosikan kursus Anda dengan membuat halaman landing page yang menarik.</p>
      <NuxtLink to="/admin/campaigns/create" class="inline-flex items-center px-4 py-2 bg-admin-600 text-white rounded-lg hover:bg-admin-700 text-sm font-medium">
        Buat Campaign Pertama →
      </NuxtLink>
    </div>

    <!-- Campaign List -->
    <div v-else class="grid gap-4">
      <div 
        v-for="campaign in campaigns" 
        :key="campaign.id"
        class="bg-white border border-neutral-200 rounded-xl p-5 hover:shadow-md transition-all"
      >
        <div class="flex items-start gap-4">
          <!-- Thumbnail -->
          <div class="w-24 h-24 rounded-lg bg-neutral-100 overflow-hidden flex-shrink-0">
            <img 
              v-if="campaign.course?.thumbnail_url" 
              :src="campaign.course.thumbnail_url" 
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-neutral-400">
              <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
            </div>
          </div>

          <!-- Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <h3 class="text-lg font-semibold text-neutral-900 truncate">{{ campaign.title }}</h3>
              <span 
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="campaign.is_active ? 'bg-green-100 text-green-700' : 'bg-neutral-100 text-neutral-600'"
              >
                <span class="w-1.5 h-1.5 rounded-full mr-1" :class="campaign.is_active ? 'bg-green-500' : 'bg-neutral-400'"></span>
                {{ campaign.is_active ? 'Active' : 'Draft' }}
              </span>
            </div>

            <!-- URL -->
            <div class="flex items-center gap-2 text-sm text-neutral-500 mb-3">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/></svg>
              <code class="bg-neutral-100 px-2 py-0.5 rounded text-xs">/c/{{ campaign.slug }}</code>
              <a 
                v-if="campaign.is_active"
                :href="`/c/${campaign.slug}`" 
                target="_blank" 
                class="text-admin-600 hover:text-admin-700 hover:underline text-xs font-medium"
              >
                Lihat Live →
              </a>
            </div>

            <!-- Stats -->
            <div class="flex items-center gap-4 text-sm">
              <div class="flex items-center gap-1.5 text-neutral-600">
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                <span>{{ campaign.view_count || 0 }} views</span>
              </div>
              <div class="flex items-center gap-1.5 text-neutral-600">
                <svg class="w-4 h-4 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122"/></svg>
                <span>{{ campaign.click_count || 0 }} clicks</span>
              </div>
              <div class="flex items-center gap-1.5 text-green-600 font-medium">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                <span>{{ campaign.conversion_count || 0 }} sales</span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2">
            <NuxtLink 
              :to="`/admin/campaigns/${campaign.id}/edit`"
              class="p-2 text-neutral-500 hover:text-admin-600 hover:bg-admin-50 rounded-lg transition-colors"
              title="Edit"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
            </NuxtLink>
            <button 
              @click="confirmDelete(campaign)"
              class="p-2 text-neutral-500 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
              title="Hapus"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="showDeleteModal = false"></div>
      <div class="bg-white rounded-xl shadow-xl max-w-md w-full p-6 relative z-10">
        <h3 class="text-lg font-bold text-neutral-900 mb-2">Hapus Campaign?</h3>
        <p class="text-neutral-600 mb-6">
          Apakah Anda yakin ingin menghapus campaign "<span class="font-medium">{{ selectedCampaign?.title }}</span>"? 
          Tindakan ini tidak dapat dibatalkan.
        </p>
        <div class="flex justify-end gap-3">
          <button 
            @click="showDeleteModal = false"
            class="px-4 py-2 border border-neutral-300 text-neutral-700 font-medium rounded-lg hover:bg-neutral-50"
          >
            Batal
          </button>
          <button 
            @click="deleteCampaign"
            class="px-4 py-2 bg-red-600 text-white font-medium rounded-lg hover:bg-red-700"
            :disabled="isDeleting"
          >
            {{ isDeleting ? 'Menghapus...' : 'Ya, Hapus' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Campaign {
  id: string
  slug: string
  title: string
  is_active: boolean
  course_id?: string
  course?: {
    id: string
    title: string
    thumbnail_url?: string
  }
  view_count: number
  click_count: number
  conversion_count: number
  created_at: string
}

definePageMeta({
  layout: 'admin',
  middleware: ['admin']
})

useHead({ title: 'Campaigns - Admin' })

const api = useApi()
const campaigns = ref<Campaign[]>([])
const pending = ref(true)
const error = ref<Error | null>(null)

const showDeleteModal = ref(false)
const selectedCampaign = ref<Campaign | null>(null)
const isDeleting = ref(false)

const fetchCampaigns = async () => {
  pending.value = true
  error.value = null
  try {
    const data = await api.fetch<Campaign[]>('/api/admin/campaigns')
    campaigns.value = data || []
  } catch (err) {
    error.value = err as Error
    console.error('[Campaigns] Failed to fetch:', err)
  } finally {
    pending.value = false
  }
}

const confirmDelete = (campaign: Campaign) => {
  selectedCampaign.value = campaign
  showDeleteModal.value = true
}

const deleteCampaign = async () => {
  if (!selectedCampaign.value) return
  
  isDeleting.value = true
  try {
    await api.fetch(`/api/admin/campaigns/${selectedCampaign.value.id}`, { method: 'DELETE' })
    campaigns.value = campaigns.value.filter(c => c.id !== selectedCampaign.value!.id)
    showDeleteModal.value = false
    selectedCampaign.value = null
  } catch (err) {
    alert('Gagal menghapus campaign')
  } finally {
    isDeleting.value = false
  }
}

onMounted(() => {
  fetchCampaigns()
})
</script>
