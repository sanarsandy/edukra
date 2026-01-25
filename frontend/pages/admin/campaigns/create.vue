<template>
  <div class="min-h-screen bg-neutral-50">
    <!-- Header -->
    <div class="bg-white border-b border-neutral-200 sticky top-0 z-50">
      <div class="flex items-center justify-between px-6 py-4">
        <div class="flex items-center gap-4">
          <NuxtLink to="/admin/campaigns" class="text-neutral-500 hover:text-neutral-700">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/></svg>
          </NuxtLink>
          <div>
            <h1 class="text-xl font-bold text-neutral-900">{{ isEdit ? 'Edit Campaign' : 'Buat Campaign Baru' }}</h1>
            <p class="text-sm text-neutral-500">Landing Page Builder</p>
          </div>
        </div>
        
        <div class="flex items-center gap-3">
          <button 
            @click="previewPage"
            class="px-4 py-2 border border-neutral-300 rounded-lg hover:bg-neutral-50 flex items-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
            Preview
          </button>

          <button 
            @click="saveDraft"
            :disabled="saving"
            class="px-4 py-2 border border-neutral-300 rounded-lg hover:bg-neutral-50 flex items-center gap-2"
          >
            {{ saving ? 'Menyimpan...' : 'Simpan Draft' }}
          </button>

          <button 
            @click="publishCampaign"
            :disabled="saving"
            class="px-5 py-2.5 bg-primary-600 text-white rounded-lg hover:bg-primary-700 font-medium"
          >
            Publish
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex h-[calc(100vh-73px)]">
      <!-- Left - Editor -->
      <div class="w-[450px] bg-white border-r border-neutral-200 flex flex-col z-20 shadow-xl overflow-hidden">
        
        <!-- Top Navigation Tabs (Content vs Appearance) -->
        <div class="flex border-b border-neutral-200 bg-neutral-50/50 backdrop-blur">
          <button 
            @click="activeTab = 'content'"
            class="flex-1 py-4 text-sm font-semibold border-b-2 transition-all flex items-center justify-center gap-2"
            :class="activeTab === 'content' ? 'border-primary-600 text-primary-700 bg-white' : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100'"
          >
            <span class="text-lg">🧱</span> Content
          </button>
          <button 
            @click="activeTab = 'appearance'"
            class="flex-1 py-4 text-sm font-semibold border-b-2 transition-all flex items-center justify-center gap-2"
            :class="activeTab === 'appearance' ? 'border-primary-600 text-primary-700 bg-white' : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:bg-neutral-100'"
            v-if="!isRigidTemplate"
          >
            <span class="text-lg">🎨</span> Appearance
          </button>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar">
          
          <!-- TAB: CONTENT -->
          <div v-show="activeTab === 'content'" class="p-6">
            <div class="mb-6">
              <label class="block text-xs font-bold text-neutral-500 uppercase tracking-wider mb-2">Campaign Info</label>
              <div class="space-y-4">
                 <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Judul Campaign *</label>
                  <input v-model="form.title"
                         @input="slugify"
                type="text" 
                class="w-full px-3 py-2 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-primary-500"
                placeholder="Flash Sale Akhir Tahun"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Slug URL *</label>
              <div class="flex items-center">
                <span class="text-sm text-neutral-500 mr-1">/c/</span>
                <input 
                  v-model="form.slug"
                  type="text" 
                  class="flex-1 px-3 py-2 border border-neutral-300 rounded-lg"
                  placeholder="flash-sale"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Tipe Campaign</label>
              <select v-model="form.campaign_type" class="w-full px-3 py-2 border border-neutral-300 rounded-lg">
                <option value="ecourse_only">📚 Hanya eCourse (akses kursus)</option>
                <option value="webinar_only">🎥 Hanya Webinar (registrasi webinar)</option>
                <option value="webinar_ecourse">📚🎥 Webinar + eCourse (keduanya)</option>
              </select>
              <p class="text-xs text-neutral-500 mt-1">
                {{ form.campaign_type === 'webinar_only' ? 'User hanya terdaftar webinar tanpa akses kursus' :
                   form.campaign_type === 'webinar_ecourse' ? 'User dapat akses kursus + terdaftar webinar' :
                   'User langsung dapat akses kursus (default)' }}
              </p>
            </div>

            <div v-if="['ecourse_only', 'webinar_ecourse'].includes(form.campaign_type)">
              <label class="block text-sm font-medium text-neutral-700 mb-1">Link ke Kursus</label>
              <select v-model="form.course_id" class="w-full px-3 py-2 border border-neutral-300 rounded-lg">
                <option :value="null">-- Pilih Kursus --</option>
                <option v-for="course in courses" :key="course.id" :value="course.id">
                  {{ course.title }}
                </option>
              </select>
            </div>

            <div v-if="['webinar_only', 'webinar_ecourse'].includes(form.campaign_type)">
              <label class="block text-sm font-medium text-neutral-700 mb-1">Link ke Webinar</label>
              <select v-model="form.webinar_id" class="w-full px-3 py-2 border border-neutral-300 rounded-lg">
                <option :value="null">-- Pilih Webinar (Opsional jika ambil dari Kursus) --</option>
                <option v-for="webinar in webinars" :key="webinar.id" :value="webinar.id">
                  {{ webinar.title }} ({{ new Date(webinar.scheduled_at).toLocaleDateString() }})
                </option>
              </select>
              <p class="text-xs text-neutral-500 mt-1">
                Jika kosong, sistem akan mencari webinar dari kursus yang dipilih.
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Countdown End Date</label>
              <input v-model="form.end_date" type="datetime-local" class="w-full px-3 py-2 border border-neutral-300 rounded-lg"/>
            </div>
          </div>

          <hr class="border-neutral-200">

          <!-- Style Customization -->
          <div v-if="!isRigidTemplate" class="space-y-4">
            <h2 class="font-semibold text-neutral-900">🎨 Kustomisasi Warna</h2>
            
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Primary</label>
                <input v-model="form.styles.primaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Harga, icons, angka statistik</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Accent</label>
                <input v-model="form.styles.accentColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Badge, gradient button</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Background</label>
                <input v-model="form.styles.backgroundColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Hero & Pricing bg</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Button</label>
                <input v-model="form.styles.buttonColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Tombol CTA</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Text Primary</label>
                <input v-model="form.styles.textPrimaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Judul, heading</p>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Text Secondary</label>
                <input v-model="form.styles.textSecondaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
                <p class="text-[10px] text-neutral-400 mt-0.5">Deskripsi, subtitle</p>
              </div>
            </div>

            <!-- End of Basic Info -->
            </div>
          </div>

          <hr class="border-neutral-200 my-6">

          <!-- Global Custom CSS Section (NEW) -->
          <details class="group">
            <summary class="flex items-center gap-2 font-semibold text-neutral-900 cursor-pointer">
              <span>⚡ Advanced: Global CSS</span>
              <svg class="w-4 h-4 group-open:rotate-180 transition-transform ml-auto text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
            </summary>
            <div class="mt-3 space-y-3">
              <p class="text-xs text-neutral-500">CSS kustom untuk seluruh landing page. Akan diterapkan di atas styling default.</p>
              <div class="p-3 bg-neutral-900 rounded-lg">
                <textarea 
                  v-model="form.globalCustomCSS"
                  rows="10"
                  placeholder="/* Contoh CSS kustom */
.hero-section h1 {
  text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
}

.pricing-card {
  transform: scale(1.05);
  box-shadow: 0 20px 40px rgba(0,0,0,0.2);
}

/* Animasi custom */
@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.cta-button {
  animation: bounce 2s infinite;
}"
                  class="w-full px-3 py-2 bg-neutral-800 border border-neutral-700 rounded font-mono text-xs text-cyan-400 placeholder-neutral-500 focus:border-orange-500 focus:ring-1 focus:ring-orange-500 resize-y"
                  spellcheck="false"
                ></textarea>
              </div>
              <div class="flex items-center gap-2 text-xs">
                <span class="px-2 py-0.5 bg-green-100 text-green-700 rounded">💡 Tips</span>
                <span class="text-neutral-500">Gunakan class pada block default untuk targeting spesifik</span>
              </div>
            </div>
          </details>

          <hr class="border-neutral-200">

          <!-- Analytics Section -->
          <div class="space-y-4">
             <h3 class="font-semibold text-neutral-900">📊 Tracking & Analytics</h3>
             <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Google Tag Manager ID</label>
                <input v-model="form.gtm_id" type="text" placeholder="GTM-XXXXXX" class="w-full px-3 py-2 border border-neutral-300 rounded-lg text-sm" />
                <p class="text-xs text-neutral-500 mt-1">Masukkan ID container GTM Anda</p>
             </div>
             <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Facebook Pixel ID</label>
                <input v-model="form.facebook_pixel_id" type="text" placeholder="1234567890" class="w-full px-3 py-2 border border-neutral-300 rounded-lg text-sm" />
             </div>
          </div>

          <hr class="border-neutral-200">

          <!-- Blocks -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-neutral-900">📦 Content Blocks</h2>
              <span class="text-xs text-neutral-500">Drag untuk urutan</span>
            </div>
            
            <ClientOnly>
              <draggable 
                v-model="form.blocks" 
                item-key="id"
                handle=".drag-handle"
                ghost-class="opacity-50"
                class="space-y-2"
              >
                <template #item="{ element: block }">
                  <div 
                    class="flex items-center gap-2 p-3 bg-neutral-50 border border-neutral-200 rounded-lg hover:border-primary-300 transition-all"
                  >
                    <div class="drag-handle cursor-grab text-neutral-400 hover:text-neutral-600">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16"/></svg>
                    </div>
                    <span class="text-lg">{{ getBlockEmoji(block.type) }}</span>
                    <div class="flex-1 min-w-0">
                      <p class="font-medium text-sm truncate">{{ getBlockTitle(block.type) }}</p>
                    </div>
                    <label class="flex items-center">
                      <input type="checkbox" v-model="block.enabled" class="rounded border-neutral-300 text-primary-600"/>
                    </label>
                    <button v-if="!isRigidTemplate" @click="editBlock(block)" class="p-1 text-neutral-500 hover:text-primary-600">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                    </button>
                  </div>
                </template>
              </draggable>
            </ClientOnly>
            
            <!-- Add Block Button -->
            <button 
              v-if="!isRigidTemplate"
              @click="showAddBlockModal = true"
              class="w-full py-4 mt-6 border-2 border-dashed border-primary-300 rounded-xl flex items-center justify-center gap-2 text-primary-600 font-bold hover:bg-primary-50 hover:border-primary-500 transition-all group"
            >
              <span class="w-8 h-8 rounded-full bg-primary-100 flex items-center justify-center group-hover:bg-primary-200 group-hover:scale-110 transition-transform">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              </span>
              <span>Tambah Block Baru</span>
            </button>
          </div>
          </div>

          <!-- TAB: APPEARANCE (New Design) -->
          <div v-show="activeTab === 'appearance'" class="flex flex-col h-full bg-neutral-50">
             
             <!-- Sub-Navigation for Appearance (Sidebar style or Horizontal Pill) -->
             <div class="px-6 py-4 bg-white border-b border-neutral-200 overflow-x-auto whitespace-nowrap hide-scrollbar flex gap-2">
                <button v-for="tab in ['themes', 'background', 'buttons', 'fonts']" :key="tab"
                        @click="appearanceTab = tab"
                        class="px-4 py-2 rounded-full text-xs font-bold transition-all border capitalization"
                        :class="appearanceTab === tab ? 'bg-neutral-900 text-white border-neutral-900' : 'bg-white text-neutral-600 border-neutral-200 hover:bg-neutral-50'">
                    {{ tab }}
                </button>
             </div>

             <div class="p-6 flex-1 overflow-y-auto">
                <!-- 1. THEMES -->
                <div v-if="appearanceTab === 'themes'">
                   <h3 class="text-sm font-bold text-neutral-900 mb-3">Preset Themes</h3>
                   
                   <!-- Category Tabs -->
                   <div class="flex flex-wrap gap-1.5 mb-4">
                      <button v-for="cat in themeCategories" :key="cat.id"
                              @click="selectedThemeCategory = cat.id"
                              class="px-2.5 py-1 text-[10px] font-bold rounded-full transition-all"
                              :class="selectedThemeCategory === cat.id 
                                ? 'bg-neutral-900 text-white' 
                                : 'bg-neutral-100 text-neutral-500 hover:bg-neutral-200'">
                         {{ cat.emoji }} {{ cat.name }}
                      </button>
                   </div>

                   <!-- Themes Grid -->
                   <div class="grid grid-cols-2 gap-3">
                     <button v-for="theme in filteredThemes" :key="theme.id" 
                             @click="applyTheme(theme)"
                             class="group relative aspect-[3/4] rounded-xl border-2 overflow-hidden transition-all text-left"
                             :class="currentThemeId === theme.id ? 'border-primary-500 ring-2 ring-primary-200' : 'border-neutral-200 hover:border-neutral-300'">
                        <div class="absolute inset-0" :style="{background: theme.thumbnail}"></div>
                        <div class="absolute inset-x-0 bottom-0 p-3 bg-gradient-to-t from-black/80 to-transparent pt-8">
                           <div class="text-white font-bold text-sm truncate">{{ theme.name }}</div>
                           <div v-if="currentThemeId === theme.id" class="text-[10px] text-green-300 font-medium flex items-center gap-1">
                              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                              Active
                           </div>
                        </div>
                     </button>
                   </div>
                </div>

                <!-- 2. BACKGROUND -->
                <div v-else-if="appearanceTab === 'background'" class="space-y-6">
                   <!-- Layout Mode -->
                   <div>
                     <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Page Layout</label>
                     <div class="flex p-1 bg-white border border-neutral-200 rounded-lg">
                        <button @click="form.styles.layoutMode = 'wide'" 
                                class="flex-1 py-2 text-xs font-medium rounded-md transition-all flex items-center justify-center gap-2"
                                :class="!form.styles.layoutMode || form.styles.layoutMode === 'wide' ? 'bg-neutral-100 text-neutral-900 shadow-sm' : 'text-neutral-500 hover:text-neutral-700'">
                           <span>🖥️</span> Wide
                        </button>
                        <button @click="form.styles.layoutMode = 'mobile'" 
                                class="flex-1 py-2 text-xs font-medium rounded-md transition-all flex items-center justify-center gap-2"
                                :class="form.styles.layoutMode === 'mobile' ? 'bg-neutral-100 text-neutral-900 shadow-sm' : 'text-neutral-500 hover:text-neutral-700'">
                           <span>📱</span> Bio Link
                        </button>
                     </div>
                   </div>

                   <!-- Background Type -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Page Background</label>
                      <div class="bg-white p-4 rounded-xl border border-neutral-200 shadow-sm space-y-4">
                         
                         <!-- Solid Color -->
                         <div>
                            <label class="text-xs text-neutral-600 mb-1.5 block">Solid Color</label>
                            <div class="flex gap-2">
                               <input type="color" v-model="form.styles.backgroundColor" class="h-10 w-12 rounded cursor-pointer border p-0.5 bg-white"/>
                               <input type="text" v-model="form.styles.backgroundColor" class="flex-1 px-3 py-2 border border-neutral-200 rounded-lg text-sm bg-neutral-50"/>
                            </div>
                         </div>

                         <!-- Gradient Toggle -->
                         <div class="flex items-center justify-between">
                            <label class="text-xs text-neutral-600">Gradient Overlay</label>
                            <button @click="form.styles.hasGradient = !form.styles.hasGradient" 
                                    class="w-10 h-5 rounded-full transition-colors relative"
                                    :class="form.styles.hasGradient ? 'bg-primary-600' : 'bg-neutral-200'">
                               <div class="w-4 h-4 bg-white rounded-full absolute top-0.5 transition-transform shadow-sm"
                                    :class="form.styles.hasGradient ? 'left-5.5' : 'left-0.5'"></div>
                            </button>
                         </div>

                         <!-- Image URL -->
                         <div>
                            <label class="text-xs text-neutral-600 mb-1.5 block">Background Image URL</label>
                            <input type="text" v-model="form.styles.backgroundImage" placeholder="https://..." class="w-full px-3 py-2 border border-neutral-200 rounded-lg text-sm bg-neutral-50"/>
                            <p class="text-[10px] text-neutral-400 mt-1">Recommended: 1920x1080px (optimized)</p>
                         </div>
                         
                         <!-- Background Overlay -->
                         <div>
                            <label class="text-xs text-neutral-600 mb-1.5 block">Overlay Opacity</label>
                            <input type="text" v-model="form.styles.backgroundOverlay" placeholder="rgba(0,0,0,0.5)" class="w-full px-3 py-2 border border-neutral-200 rounded-lg text-sm bg-neutral-50"/>
                         </div>
                      </div>
                   </div>
                </div>

                <!-- 3. BUTTONS -->
                <div v-else-if="appearanceTab === 'buttons'" class="space-y-6">
                   <!-- Button Shape -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Button Shape</label>
                      <div class="grid grid-cols-2 gap-3">
                         <button @click="form.styles.buttonShape = 'sharp'" class="h-10 border border-neutral-300 bg-white hover:border-neutral-400 transition-colors flex items-center justify-center font-medium text-xs text-neutral-600" :class="{'ring-2 ring-primary-500 border-primary-500 text-primary-700': form.styles.buttonShape === 'sharp'}">Sharp</button>
                         <button @click="form.styles.buttonShape = 'rounded'" class="h-10 border border-neutral-300 bg-white hover:border-neutral-400 rounded-lg transition-colors flex items-center justify-center font-medium text-xs text-neutral-600" :class="{'ring-2 ring-primary-500 border-primary-500 text-primary-700': form.styles.buttonShape === 'rounded'}">Rounded</button>
                         <button @click="form.styles.buttonShape = 'pill'" class="h-10 border border-neutral-300 bg-white hover:border-neutral-400 rounded-full transition-colors flex items-center justify-center font-medium text-xs text-neutral-600" :class="{'ring-2 ring-primary-500 border-primary-500 text-primary-700': form.styles.buttonShape === 'pill'}">Pill</button>
                         <button @click="form.styles.buttonShape = 'leaf'" class="h-10 border border-neutral-300 bg-white hover:border-neutral-400 rounded-tl-2xl rounded-br-2xl transition-colors flex items-center justify-center font-medium text-xs text-neutral-600" :class="{'ring-2 ring-primary-500 border-primary-500 text-primary-700': form.styles.buttonShape === 'leaf'}">Leaf</button>
                      </div>
                   </div>

                   <!-- Button Style -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Button Style</label>
                      <div class="grid grid-cols-2 gap-3">
                         <button v-for="style in ['solid', 'outline', 'gradient', 'glass', 'neo-brutalism', 'soft-shadow']" :key="style"
                                 @click="form.styles.buttonStyle = style"
                                 class="h-10 rounded-lg flex items-center justify-center text-xs font-bold transition-all capitalize border"
                                 :class="form.styles.buttonStyle === style ? 'border-primary-500 bg-primary-50 text-primary-700 ring-1 ring-primary-500' : 'border-neutral-200 bg-white text-neutral-600 hover:bg-neutral-50'">
                            {{ style.replace('-', ' ') }}
                         </button>
                      </div>
                   </div>

                   <!-- Button Colors -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Colors</label>
                      <div class="bg-white p-4 rounded-xl border border-neutral-200 space-y-3">
                         <div>
                            <label class="text-xs text-neutral-500 mb-1 block">Button Background</label>
                            <div class="flex gap-2">
                               <input type="color" v-model="form.styles.buttonColor" class="h-8 w-10 rounded cursor-pointer border p-0.5"/>
                               <input type="text" v-model="form.styles.buttonColor" class="flex-1 px-2 py-1 border border-neutral-200 rounded text-xs"/>
                            </div>
                         </div>
                         <div>
                            <label class="text-xs text-neutral-500 mb-1 block">Primary Accent</label>
                            <div class="flex gap-2">
                               <input type="color" v-model="form.styles.primaryColor" class="h-8 w-10 rounded cursor-pointer border p-0.5"/>
                               <input type="text" v-model="form.styles.primaryColor" class="flex-1 px-2 py-1 border border-neutral-200 rounded text-xs"/>
                            </div>
                         </div>
                      </div>
                   </div>
                </div>

                  <!-- 4. FONTS -->
                <div v-else-if="appearanceTab === 'fonts'" class="space-y-6">
                   <!-- Heading Font -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Heading Font</label>
                      <div class="space-y-2">
                         <button v-for="font in ['Inter', 'Poppins', 'Playfair Display', 'Roboto', 'Montserrat', 'Orbitron']" :key="font"
                                 @click="form.styles.fontFamilyHeading = font"
                                 class="w-full p-3 border rounded-xl text-left hover:border-neutral-300 transition-all flex items-center justify-between group"
                                 :class="form.styles.fontFamilyHeading === font ? 'border-primary-500 bg-primary-50' : 'border-neutral-200 bg-white'">
                            <span class="text-lg" :style="{fontFamily: font}">{{ font }}</span>
                            <span v-if="form.styles.fontFamilyHeading === font" class="text-primary-600">
                               <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
                            </span>
                         </button>
                      </div>
                   </div>

                   <!-- Body Font -->
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Body Font</label>
                      <select v-model="form.styles.fontFamilyBody" class="w-full p-3 border border-neutral-200 rounded-xl bg-white">
                         <option v-for="font in ['Inter', 'Roboto', 'Lato', 'Quicksand', 'Space Grotesk']" :key="font" :value="font">{{ font }}</option>
                      </select>
                   </div>
                   
                   <div>
                      <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Text Colors</label>
                      <div class="bg-white p-4 rounded-xl border border-neutral-200 space-y-3">
                         <div>
                            <label class="text-xs text-neutral-500 mb-1 block">Primary Text (Headings)</label>
                            <div class="flex gap-2">
                               <input type="color" v-model="form.styles.textPrimaryColor" class="h-8 w-10 rounded cursor-pointer border p-0.5"/>
                               <input type="text" v-model="form.styles.textPrimaryColor" class="flex-1 px-2 py-1 border border-neutral-200 rounded text-xs"/>
                            </div>
                         </div>
                       </div>
                    </div>
                </div>

             </div>
          </div>
        </div>
      </div>

      <!-- Right - Preview Area -->
      <div class="flex-1 overflow-y-auto bg-neutral-100 p-6">
        <div class="max-w-4xl mx-auto">
          <!-- Global Custom CSS -->
          <component :is="'style'" v-if="form.globalCustomCSS">
            {{ form.globalCustomCSS }}
          </component>
          <!-- Device Selector (Mobile Mode Only) -->
          <div v-if="form.styles.layoutMode === 'mobile'" class="flex justify-center mb-6">
             <div class="bg-white rounded-lg shadow-sm border border-neutral-200 p-1 flex gap-1">
               <button 
                 v-for="(preset, key) in devicePresets" :key="key"
                 @click="selectedDevice = key"
                 class="px-3 py-1.5 text-xs font-medium rounded-md transition-colors"
                 :class="selectedDevice === key ? 'bg-neutral-900 text-white' : 'text-neutral-500 hover:bg-neutral-50'"
               >
                 {{ preset.name }}
               </button>
             </div>
          </div>

          <!-- Live Preview Container -->
          <div class="transition-all duration-500 mx-auto"
               :class="{
                 'bg-white border-neutral-900 shadow-2xl overflow-hidden relative': form.styles.layoutMode === 'mobile',
                 'w-full bg-white rounded-xl shadow-lg overflow-hidden': !form.styles.layoutMode || form.styles.layoutMode === 'wide',
                 'border-[8px]': form.styles.layoutMode === 'mobile' && devicePresets[selectedDevice].radius === '4px', // Basic border
                 'border-[12px] shadow-[0_0_0_2px_#333,0_20px_40px_-10px_rgba(0,0,0,0.5)]': form.styles.layoutMode === 'mobile' && devicePresets[selectedDevice].radius !== '4px' // Premium border
               }"
               :style="form.styles.layoutMode === 'mobile' ? {
                  width: devicePresets[selectedDevice].width,
                  height: devicePresets[selectedDevice].height,
                  borderRadius: devicePresets[selectedDevice].radius
               } : previewStyles">
            
            <!-- Dynamic Island / Notch -->
            <div v-if="form.styles.layoutMode === 'mobile'">
                <!-- Island -->
                <div v-if="devicePresets[selectedDevice].notch === 'island'" 
                     class="absolute top-3 left-1/2 -translate-x-1/2 h-7 w-[120px] bg-black rounded-full z-50 transition-all"></div>
                <!-- Standard Notch -->
                <div v-else-if="devicePresets[selectedDevice].notch === 'notch'" 
                     class="absolute top-0 left-1/2 -translate-x-1/2 h-6 w-36 bg-neutral-900 rounded-b-xl z-50"></div>
                <!-- Punch Hole -->
                <div v-else-if="devicePresets[selectedDevice].notch === 'hole'" 
                     class="absolute top-4 left-1/2 -translate-x-1/2 h-4 w-4 bg-black rounded-full z-50"></div>
            </div>
            
            <!-- Status Bar (Mobile Only) -->
            <div v-if="form.styles.layoutMode === 'mobile'" class="absolute top-0 inset-x-0 h-10 px-6 flex justify-between items-center z-50 pointer-events-none text-black mix-blend-difference">
               <span class="text-[10px] font-bold">9:41</span>
               <div class="flex gap-1.5 opacity-80">
                  <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 14H9v-2h2v2zm0-4H9V7h2v5z"/></svg>
                  <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 24 24"><path d="M15.67 4H14V2h-4v2H8.33C7.6 4 7 4.6 7 5.33v15.33C7 21.4 7.6 22 8.33 22h7.33c.74 0 1.34-.6 1.34-1.33V5.33C17 4.6 16.4 4 15.67 4z"/></svg>
               </div>
            </div>
            
            <div class="h-full overflow-y-auto scrollbar-hide" :class="{'pt-8 pb-4': form.styles.layoutMode === 'mobile'}">
            <!-- Global Custom CSS (injected inside frame) -->
            <component :is="'style'" v-if="form.globalCustomCSS">
              {{ form.globalCustomCSS }}
            </component>
            <template v-for="block in enabledBlocks" :key="block.id">
              <!-- Block Wrapper with Resize Handle -->
              <div class="relative group" :style="getBlockPreviewStyle(block)">
                <!-- Height Indicator (on hover) -->
                <div v-if="block.data.minHeight" 
                     class="absolute -left-12 top-1/2 -translate-y-1/2 bg-neutral-800 text-white text-[10px] px-1.5 py-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap z-10">
                  H: {{ block.data.minHeight }}px
                </div>
                
                <!-- Resize Handle (bottom) -->
                <div 
                  class="absolute bottom-0 left-0 right-0 h-3 cursor-ns-resize bg-gradient-to-t from-primary-500/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                  @mousedown="startResize(block, $event)"
                >
                  <div class="w-12 h-1 bg-primary-500 rounded-full"></div>
                </div>
                
                <!-- Quick Edit Button -->
                <button 
                  @click="editBlock(block)"
                  class="absolute top-2 right-2 p-1.5 bg-white/90 rounded-lg shadow opacity-0 group-hover:opacity-100 transition-opacity z-10 hover:bg-primary-50"
                >
                  <svg class="w-4 h-4 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                </button>


              <!-- Custom CSS Injection -->
              <component :is="'style'" v-if="block.data.customCSS">
                {{ getScopedCSS(block) }}
              </component>

              <!-- Block Templates -->
              <div class="block-content">
              <!-- Hero Preview -->
              <div v-if="block.type === 'hero'" 
                   class="relative min-h-[300px] text-white p-8 flex flex-col justify-center items-center text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : '300px', background: block.data.backgroundColor || 'linear-gradient(to bottom right, #1f2937, #111827)'}">
                <h1 class="text-3xl md:text-4xl font-bold mb-4" :style="{fontFamily: previewStyles['--font-heading'], color: block.data.titleColor || '#ffffff'}">{{ block.data.headline }}</h1>
                <p class="text-lg mb-6 max-w-xl" :style="{color: block.data.textColor || 'rgba(255,255,255,0.7)'}">{{ block.data.subheadline }}</p>
                <button class="px-6 py-3 font-semibold" :style="{...previewButtonStyle, backgroundColor: block.data.buttonColor || previewButtonStyle.backgroundColor, color: block.data.buttonTextColor || previewButtonStyle.color}">
                  {{ block.data.cta_text || 'Daftar Sekarang' }}
                </button>
              </div>

              <!-- Countdown Preview -->
              <div v-else-if="block.type === 'countdown'" 
                   class="text-white py-8 text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#111827'}">
                <p class="mb-4" :style="{color: block.data.labelColor || 'rgba(255,255,255,0.6)'}">{{ block.data.label }}</p>
                <div class="flex justify-center gap-4">
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: block.data.numberColor || block.data.accentColor || form.styles.primaryColor}">00</div>
                    <div class="text-xs" :style="{color: block.data.labelColor || 'rgba(255,255,255,0.5)'}">Hari</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: block.data.numberColor || block.data.accentColor || form.styles.primaryColor}">00</div>
                    <div class="text-xs" :style="{color: block.data.labelColor || 'rgba(255,255,255,0.5)'}">Jam</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: block.data.numberColor || block.data.accentColor || form.styles.primaryColor}">00</div>
                    <div class="text-xs" :style="{color: block.data.labelColor || 'rgba(255,255,255,0.5)'}">Menit</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: block.data.numberColor || block.data.accentColor || form.styles.primaryColor}">00</div>
                    <div class="text-xs" :style="{color: block.data.labelColor || 'rgba(255,255,255,0.5)'}">Detik</div>
                  </div>
                </div>
              </div>

              <!-- Benefits Preview -->
              <div v-else-if="block.type === 'benefits'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                <h2 class="text-2xl font-bold text-center mb-6" :style="{color: block.data.titleColor || form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title }}</h2>
                <div class="grid md:grid-cols-2 gap-4">
                  <div v-for="(item, idx) in block.data.items" :key="idx" class="flex items-center gap-3 p-3 bg-neutral-50 rounded-lg">
                    <div class="w-8 h-8 rounded-full flex items-center justify-center" :style="{backgroundColor: (block.data.iconColor || block.data.accentColor || form.styles.primaryColor) + '20', color: block.data.iconColor || block.data.accentColor || form.styles.primaryColor}">
                      ✓
                    </div>
                    <span :style="{color: block.data.textColor || form.styles.textSecondaryColor}">{{ item.text }}</span>
                  </div>
                </div>
              </div>

              <!-- Pricing Preview -->
              <div v-else-if="block.type === 'pricing'" 
                   class="py-10 px-8 text-white text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', background: block.data.backgroundColor || 'linear-gradient(to bottom, #111827, #1f2937)'}">
                <h2 class="text-2xl font-bold mb-6" :style="{fontFamily: previewStyles['--font-heading'], color: block.data.titleColor || '#ffffff'}">Penawaran Spesial</h2>
                <div class="text-neutral-400 line-through text-lg">Rp {{ formatPrice(block.data.original_price) }}</div>
                <div class="text-4xl font-bold my-2" :style="{color: block.data.priceColor || block.data.accentColor || form.styles.primaryColor}">
                  Rp {{ formatPrice(block.data.discount_price) }}
                </div>
                <button class="mt-4 px-8 py-3 font-bold" :style="{...previewButtonStyle, backgroundColor: block.data.buttonColor || previewButtonStyle.backgroundColor, color: block.data.buttonTextColor || previewButtonStyle.color}">
                  {{ block.data.cta_text || 'Beli Sekarang' }}
                </button>
              </div>

              <!-- Testimonials Preview -->
              <div v-else-if="block.type === 'testimonials'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#f5f5f5'}">
                <h2 class="text-2xl font-bold text-center mb-6" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title }}</h2>
                <div class="grid md:grid-cols-2 gap-4">
                  <div v-for="(item, idx) in block.data.items" :key="idx" class="bg-white p-4 rounded-lg shadow-sm">
                    <p class="text-neutral-600 italic mb-3">"{{ item.text }}"</p>
                    <div class="font-medium text-neutral-900">- {{ item.name }}</div>
                  </div>
                  <div v-if="!block.data.items?.length" class="col-span-2 text-center text-neutral-400 py-8">
                    Belum ada testimonial
                  </div>
                </div>
              </div>

              <!-- FAQ Preview -->
              <div v-else-if="block.type === 'faq'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                <h2 class="text-2xl font-bold text-center mb-6" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title }}</h2>
                <div class="space-y-3 max-w-2xl mx-auto">
                  <div v-for="(item, idx) in block.data.items" :key="idx" class="border border-neutral-200 rounded-lg p-4">
                    <div class="font-medium" :style="{color: form.styles.textPrimaryColor}">{{ item.question }}</div>
                    <div class="text-sm mt-2" :style="{color: form.styles.textSecondaryColor}">{{ item.answer }}</div>
                  </div>
                  <div v-if="!block.data.items?.length" class="text-center text-neutral-400 py-8">
                    Belum ada FAQ
                  </div>
                </div>
              </div>

              <!-- Instructor Preview -->
              <div v-else-if="block.type === 'instructor'" 
                   class="py-10 px-8 text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#fafafa'}">
                <h2 class="text-2xl font-bold mb-6" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">👨‍🏫 Tentang Instruktur</h2>
                <div class="w-20 h-20 bg-neutral-300 rounded-full mx-auto mb-4 flex items-center justify-center text-3xl">
                  👨‍🏫
                </div>
                <div class="font-bold text-lg" :style="{color: form.styles.textPrimaryColor}">{{ block.data.name || 'Nama Instruktur' }}</div>
                <p class="max-w-md mx-auto mt-2" :style="{color: form.styles.textSecondaryColor}">{{ block.data.bio || 'Bio instruktur' }}</p>
              </div>

              <!-- Video Preview (NEW) -->
              <div v-else-if="block.type === 'video'" 
                   class="py-10 px-8 text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#111827'}">
                <h2 v-if="block.data.title" class="text-2xl font-bold text-white mb-2">{{ block.data.title }}</h2>
                <p v-if="block.data.subtitle" class="text-neutral-400 mb-6">{{ block.data.subtitle }}</p>
                <div class="max-w-2xl mx-auto aspect-video bg-neutral-800 rounded-xl flex items-center justify-center">
                  <div v-if="block.data.youtube_url" class="text-white text-center">
                    <span class="text-4xl">▶️</span>
                    <p class="text-sm text-neutral-400 mt-2">YouTube Video</p>
                  </div>
                  <div v-else class="text-neutral-500 text-center">
                    <span class="text-4xl">🎬</span>
                    <p class="text-sm mt-2">Tambahkan YouTube URL</p>
                  </div>
                </div>
              </div>

              <!-- Trust Badges Preview (NEW) -->
              <div v-else-if="block.type === 'trust'" 
                   class="py-8 px-6 border-y border-neutral-100"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                <div class="flex flex-wrap items-center justify-center gap-6">
                  <div class="flex items-center gap-2 text-neutral-600">
                    <span class="text-xl">🔒</span>
                    <span class="text-sm font-medium">Pembayaran Aman</span>
                  </div>
                  <div class="flex items-center gap-2 text-neutral-600">
                    <span class="text-xl">💯</span>
                    <span class="text-sm font-medium">Garansi 30 Hari</span>
                  </div>
                  <div class="flex items-center gap-2 text-neutral-600">
                    <span class="text-xl">🎓</span>
                    <span class="text-sm font-medium">Sertifikat</span>
                  </div>
                  <div class="flex items-center gap-2 text-neutral-600">
                    <span class="text-xl">♾️</span>
                    <span class="text-sm font-medium">Akses Selamanya</span>
                  </div>
                </div>
              </div>

              <!-- CTA Banner Preview (NEW) -->
              <div v-else-if="block.type === 'cta_banner'" 
                   class="py-10 px-8 text-white text-center" 
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || form.styles.buttonColor}">
                <h2 class="text-2xl font-bold mb-3" :style="{fontFamily: previewStyles['--font-heading']}">{{ block.data.headline || '🚀 Siap Untuk Memulai?' }}</h2>
                <p class="opacity-80 mb-6">{{ block.data.subheadline || 'Jangan lewatkan kesempatan untuk meningkatkan skill Anda' }}</p>
                <button class="px-6 py-3 bg-white font-semibold" :style="{color: form.styles.buttonColor, borderRadius: previewButtonStyle.borderRadius}">
                  {{ block.data.cta_text || 'Daftar Sekarang' }}
                </button>
              </div>

              <!-- Statistics Preview (NEW) -->
              <div v-else-if="block.type === 'statistics'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                <h2 v-if="block.data.title" class="text-xl font-bold text-center mb-6" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title }}</h2>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
                  <div v-for="(stat, idx) in block.data.items" :key="idx" class="text-center p-4 bg-neutral-50 rounded-xl">
                    <div class="text-2xl font-black mb-1" :style="{color: block.data.accentColor || form.styles.primaryColor}">
                      {{ stat.value }}{{ stat.suffix || '' }}
                    </div>
                    <div class="text-xs font-medium" :style="{color: form.styles.textSecondaryColor}">{{ stat.label }}</div>
                  </div>
                </div>
              </div>

              <!-- Bonus Preview (NEW) -->
              <div v-else-if="block.type === 'bonus'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#fffbeb'}">
                <div class="text-center mb-6">
                  <span class="inline-block text-[10px] font-bold px-2 py-0.5 rounded-full mb-2" :style="{backgroundColor: block.data.accentColor || '#facc15', color: '#78350f'}">🎁 BONUS</span>
                  <h2 class="text-xl font-bold" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title || 'Bonus Eksklusif' }}</h2>
                </div>
                <div class="space-y-3">
                  <div v-for="(bonus, idx) in block.data.items" :key="idx" class="flex items-center gap-3 p-3 bg-white rounded-xl shadow-sm border border-orange-100">
                    <div class="w-10 h-10 rounded-lg flex items-center justify-center text-xl bg-orange-100">
                      {{ bonus.emoji || '🎁' }}
                    </div>
                    <div class="flex-1 min-w-0">
                      <h4 class="font-bold text-neutral-900 text-sm">{{ bonus.title }}</h4>
                      <p v-if="bonus.description" class="text-neutral-500 text-xs">{{ bonus.description }}</p>
                    </div>
                    <div class="text-green-600 font-bold text-xs">GRATIS</div>
                  </div>
                </div>
                <div v-if="block.data.total_value" class="mt-6 p-4 bg-gradient-to-r from-green-500 to-emerald-600 rounded-xl text-white text-center">
                  <p class="text-xs opacity-80">Total Nilai Bonus</p>
                  <p class="text-xl font-black">Rp {{ formatPrice(block.data.total_value) }}</p>
                </div>
              </div>

              <!-- Curriculum Preview (NEW) -->
              <div v-else-if="block.type === 'curriculum'" 
                   class="py-10 px-8"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                <div class="text-center mb-8">
                  <h2 class="text-xl font-bold" :style="{color: form.styles.textPrimaryColor, fontFamily: previewStyles['--font-heading']}">{{ block.data.title || 'Materi Kursus' }}</h2>
                  <p v-if="block.data.subtitle" class="text-sm" :style="{color: form.styles.textSecondaryColor}">{{ block.data.subtitle }}</p>
                </div>
                <div class="space-y-2">
                  <div v-for="(module, idx) in block.data.modules" :key="idx" class="border rounded-lg overflow-hidden">
                    <div class="p-3 bg-neutral-50 flex items-center gap-3">
                      <span class="w-6 h-6 bg-neutral-200 rounded flex items-center justify-center text-xs font-bold">{{ idx + 1 }}</span>
                      <div class="flex-1 font-semibold text-sm">{{ module.title }}</div>
                      <span class="text-xs text-neutral-500">{{ module.lessons_count }} Lessons</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Gallery Preview -->
              <div v-else-if="block.type === 'gallery'" 
                   class="py-10 text-center"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                 <h2 v-if="block.data.title" class="font-bold mb-4" :style="{color: form.styles.textPrimaryColor}">{{ block.data.title }}</h2>
                 <div class="flex gap-2 overflow-x-auto px-4 pb-2 snap-x">
                    <div v-for="(img, idx) in block.data.items" :key="idx" class="w-48 h-32 flex-shrink-0 bg-neutral-200 rounded-lg overflow-hidden relative snap-center">
                       <img v-if="img.image" :src="img.image" class="w-full h-full object-cover"/>
                       <div v-else class="w-full h-full flex items-center justify-center text-4xl">🖼️</div>
                    </div>
                 </div>
              </div>

              <!-- Achievement Preview -->
              <div v-else-if="block.type === 'achievement'" 
                   class="py-6 border-y border-neutral-100"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#fafafa'}">
                 <div class="flex flex-wrap justify-center gap-4 px-4">
                    <div v-for="(badge, idx) in block.data.items" :key="idx" class="flex items-center gap-2 bg-white px-3 py-2 rounded-full shadow-sm border border-neutral-100">
                       <span>{{ badge.emoji }}</span>
                       <span class="text-sm font-medium">{{ badge.text }}</span>
                    </div>
                 </div>
              </div>
              <!-- Comparison Preview -->
              <div v-else-if="block.type === 'comparison'" 
                   class="py-10 px-4"
                   :style="{minHeight: block.data.minHeight ? block.data.minHeight + 'px' : 'auto', backgroundColor: block.data.backgroundColor || '#ffffff'}">
                 <h2 v-if="block.data.title" class="text-center font-bold mb-4" :style="{color: form.styles.textPrimaryColor}">{{ block.data.title }}</h2>
                 <div class="overflow-x-auto">
                    <table class="w-full text-sm border-collapse">
                       <thead>
                          <tr class="bg-neutral-50">
                             <th class="p-2 border text-left">{{ block.data.headers?.[0] || 'Fitur' }}</th>
                             <th class="p-2 border text-center font-bold text-neutral-500">{{ block.data.headers?.[1] || 'A' }}</th>
                             <th class="p-2 border text-center font-bold text-primary-600 bg-primary-50">{{ block.data.headers?.[2] || 'B' }}</th>
                          </tr>
                       </thead>
                       <tbody>
                          <tr v-for="(row, idx) in block.data.rows" :key="idx" class="border-b">
                             <td class="p-2 border-r">{{ row.feature }}</td>
                             <td class="p-2 text-center border-r bg-neutral-50/30">
                                <span v-if="row.val_a === true">✅</span>
                                <span v-else-if="row.val_a === false" class="text-neutral-300">❌</span>
                                <span v-else>{{ row.text_a }}</span>
                             </td>
                             <td class="p-2 text-center bg-primary-50/10">
                                <span v-if="row.val_b === true">✅</span>
                                <span v-else-if="row.val_b === false" class="text-neutral-300">❌</span>
                                <span v-else>{{ row.text_b }}</span>
                             </td>
                          </tr>
                       </tbody>
                    </table>
                 </div>
              </div>
              <!-- Floating Chat Preview -->
              <div v-else-if="block.type === 'floating_chat'" 
                   class="py-10 flex items-center justify-center relative h-40"
                   :style="{backgroundColor: block.data.backgroundColor || '#f5f5f5'}">
                 <div class="absolute bottom-4 right-4 bg-green-500 text-white px-4 py-3 rounded-full shadow-lg flex items-center gap-2 cursor-pointer hover:bg-green-600 transition-colors">
                    <span class="text-xl">💬</span>
                    <span v-if="block.data.label" class="font-bold text-sm">{{ block.data.label }}</span>
                 </div>
                 <p class="text-neutral-400 text-sm">Tombol akan melayang di pojok kanan bawah halaman.</p>
              </div>

              <!-- Social Proof Preview -->
              <div v-else-if="block.type === 'social_proof'" 
                   class="py-10 flex items-center justify-center relative h-32"
                   :style="{backgroundColor: block.data.backgroundColor || '#f5f5f5'}">
                 <div class="absolute bottom-4 left-4 bg-white/90 backdrop-blur border border-neutral-200 px-4 py-2 rounded-lg shadow-lg flex items-center gap-3 animate-pulse">
                    <div class="w-2 h-2 rounded-full bg-red-500 animate-ping"></div>
                    <span class="text-sm font-medium text-neutral-800">
                      <span class="font-bold">24</span> {{ block.data.text || 'orang sedang melihat ini' }}
                    </span>
                 </div>
                 <p class="text-neutral-400 text-sm">Notifikasi Live Viewer di pojok kiri bawah.</p>
              </div>

               <!-- Gen Z Blocks -->
               <div v-else-if="block.type === 'hero_gen_z'">
                 <CampaignHeroGenZ :block="block" :styles="form.styles" :global-end-date="form.end_date" />
               </div>
               <div v-else-if="block.type === 'features_gen_z'">
                 <CampaignFeaturesGenZ :block="block" :styles="form.styles" />
               </div>
               <div v-else-if="block.type === 'cta_gen_z'">
                 <CampaignCtaGenZ :block="block" :styles="form.styles" :campaign-id="campaignId" />
               </div>
               <!-- Clean/TikTok Blocks -->
               <div v-else-if="block.type === 'hero_clean'">
                  <CampaignHeroClean :block="block" :styles="form.styles" :global-end-date="form.end_date" />
               </div>
                <div v-else-if="block.type === 'speaker_clean'">
                  <CampaignSpeakerClean :block="block" :styles="form.styles" />
               </div>
                <div v-else-if="block.type === 'features_clean'">
                  <CampaignFeaturesClean :block="block" :styles="form.styles" />
               </div>
                <div v-else-if="block.type === 'cta_clean'">
                  <CampaignCtaClean :block="block" :styles="form.styles" :campaign-id="campaignId" />
               </div>

                <!-- Pro Blocks (Professional Minimalist) -->
                <div v-else-if="block.type === 'hero_pro'">
                  <CampaignHeroPro :block="block" :styles="form.styles" :global-end-date="form.end_date" />
               </div>
                <div v-else-if="block.type === 'speaker_pro'">
                  <CampaignSpeakerPro :block="block" :styles="form.styles" />
               </div>
                <div v-else-if="block.type === 'features_pro'">
                  <CampaignFeaturesPro :block="block" :styles="form.styles" />
               </div>
                <div v-else-if="block.type === 'cta_pro'">
                  <CampaignCtaPro :block="block" :styles="form.styles" :campaign-id="campaignId" />
                </div>

                <!-- Healing Blocks (P.U.L.I.H Theme) -->
                <div v-else-if="block.type === 'hero_healing'">
                  <CampaignHeroHealing :block="block" :styles="form.styles" :global-end-date="form.end_date" />
                </div>
                <div v-else-if="block.type === 'problem_healing'">
                  <CampaignProblemHealing :block="block" :styles="form.styles" />
                </div>
                <div v-else-if="block.type === 'solution_healing'">
                  <CampaignSolutionHealing :block="block" :styles="form.styles" />
                </div>
                <div v-else-if="block.type === 'speaker_healing'">
                  <CampaignSpeakerHealing :block="block" :styles="form.styles" />
                </div>
                <div v-else-if="block.type === 'cta_healing'">
                  <CampaignCtaHealing :block="block" :styles="form.styles" :campaign-id="campaignId" />
                </div>

              </div><!-- Close Default Block Templates -->
              </div><!-- Close Block Wrapper -->
            </template>

            <div v-if="!enabledBlocks.length" class="py-20 text-center text-neutral-400">
              <p class="text-lg">Aktifkan block untuk melihat preview</p>
            </div>
          </div>
      </div>
    </div>

    <!-- Edit Block Modal -->
    <div v-if="editingBlock" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div class="bg-white rounded-xl w-full max-w-lg max-h-[80vh] overflow-y-auto">
        <div class="flex items-center justify-between p-4 border-b">
          <h3 class="font-semibold text-lg">{{ getBlockEmoji(editingBlock.type) }} Edit {{ getBlockTitle(editingBlock.type) }}</h3>
          <button @click="editingBlock = null" class="text-neutral-500 hover:text-neutral-700">✕</button>
        </div>
        
        <div class="p-4 space-y-4">
          <!-- Variant Selector (Generic) -->
          <div class="p-3 bg-neutral-50 rounded-lg border border-neutral-100">
             <label class="block text-xs font-bold text-neutral-500 uppercase mb-1">Layout Variant</label>
             <select v-model="editingBlock.variant" class="w-full px-3 py-2 border rounded-lg text-sm bg-white">
                <option value="default">Default / Standard</option>
                <template v-if="editingBlock.type === 'hero'">
                  <option value="centered">Centered</option>
                  <option value="split">Split (Text Left - Image Right)</option>
                  <option value="bio_profile">Bio Link Profile (Avatar &amp; Header)</option>
                </template>
                <template v-if="editingBlock.type === 'pricing'">
                  <option value="highlight">Highlight Center</option>
                  <option value="minimal">Minimal List</option>
                </template>
                <template v-if="editingBlock.type === 'testimonials'">
                  <option value="grid">Grid View</option>
                  <option value="carousel">Carousel View</option>
                </template>
                <template v-if="editingBlock.type === 'gallery'">
                  <option value="grid">Grid Masonry</option>
                  <option value="carousel">Carousel/Slider</option>
                </template>
             </select>
          </div>

          <!-- Per-Block Styling Section (ENHANCED) -->
          <details class="group" open>
             <summary v-if="!isRigidTemplate" class="flex items-center gap-2 text-xs font-bold text-purple-700 cursor-pointer p-3 bg-gradient-to-br from-purple-50 to-indigo-50 rounded-lg border border-purple-200 hover:bg-purple-100 transition-colors">
               <span>🎨 Warna & Styling Blok Ini</span>
               <svg class="w-3 h-3 group-open:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
             </summary>
             
             <div class="mt-3 space-y-3">
               <!-- Min Height (Universal) -->
               <div class="p-3 bg-neutral-50 rounded-lg">
                 <label class="block text-xs font-medium text-neutral-600 mb-1">📐 Min Height (px)</label>
                 <div class="flex items-center gap-2">
                   <input type="range" v-model.number="editingBlock.data.minHeight" min="0" max="600" step="20" class="flex-1"/>
                   <input type="number" v-model.number="editingBlock.data.minHeight" min="0" max="600" step="10" 
                          class="w-16 px-2 py-1 border rounded text-sm text-center"/>
                 </div>
               </div>

               <!-- Background Colors -->
               <details class="group/bg">
                 <summary class="flex items-center gap-2 text-xs font-medium text-neutral-700 cursor-pointer p-2 bg-neutral-100 rounded hover:bg-neutral-200 transition-colors">
                   <span>🖼️ Background</span>
                   <div class="w-4 h-4 rounded border" :style="{backgroundColor: editingBlock.data.backgroundColor || '#transparent'}"></div>
                   <svg class="w-3 h-3 group-open/bg:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                 </summary>
                 <div class="mt-2 p-3 bg-white rounded border space-y-2">
                   <div class="flex items-center gap-2">
                     <input type="color" v-model="editingBlock.data.backgroundColor" class="h-8 w-10 rounded cursor-pointer border-0 p-0"/>
                     <input type="text" v-model="editingBlock.data.backgroundColor" placeholder="#ffffff" class="flex-1 px-2 py-1 border rounded text-sm uppercase"/>
                     <button @click="editingBlock.data.backgroundColor = undefined" class="text-xs text-red-500 hover:underline">Reset</button>
                   </div>
                 </div>
               </details>

               <!-- Text Colors (Show for blocks with text) -->
               <details v-if="['hero', 'benefits', 'pricing', 'faq', 'instructor', 'cta_banner', 'bonus', 'curriculum'].includes(editingBlock.type)" class="group/text">
                 <summary class="flex items-center gap-2 text-xs font-medium text-neutral-700 cursor-pointer p-2 bg-neutral-100 rounded hover:bg-neutral-200 transition-colors">
                   <span>✏️ Warna Teks</span>
                   <div class="flex gap-1">
                     <div class="w-4 h-4 rounded border" :style="{backgroundColor: editingBlock.data.titleColor || '#111827'}"></div>
                     <div class="w-4 h-4 rounded border" :style="{backgroundColor: editingBlock.data.textColor || '#6b7280'}"></div>
                   </div>
                   <svg class="w-3 h-3 group-open/text:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                 </summary>
                 <div class="mt-2 p-3 bg-white rounded border space-y-3">
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Title / Heading</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.titleColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.titleColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.titleColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Body / Description</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.textColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.textColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.textColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                 </div>
               </details>

               <!-- Button Colors (Show for blocks with buttons) -->
               <details v-if="['hero', 'pricing', 'cta_banner'].includes(editingBlock.type)" class="group/btn">
                 <summary class="flex items-center gap-2 text-xs font-medium text-neutral-700 cursor-pointer p-2 bg-neutral-100 rounded hover:bg-neutral-200 transition-colors">
                   <span>🔘 Warna Tombol</span>
                   <div class="w-4 h-4 rounded border" :style="{backgroundColor: editingBlock.data.buttonColor || '#6366f1'}"></div>
                   <svg class="w-3 h-3 group-open/btn:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                 </summary>
                 <div class="mt-2 p-3 bg-white rounded border space-y-3">
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Button Background</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.buttonColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.buttonColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.buttonColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Button Text</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.buttonTextColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.buttonTextColor" placeholder="#ffffff" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.buttonTextColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                 </div>
               </details>

               <!-- Accent/Icon Colors (Show for relevant blocks) -->
               <details v-if="['benefits', 'countdown', 'statistics', 'pricing', 'faq', 'comparison'].includes(editingBlock.type)" class="group/accent">
                 <summary class="flex items-center gap-2 text-xs font-medium text-neutral-700 cursor-pointer p-2 bg-neutral-100 rounded hover:bg-neutral-200 transition-colors">
                   <span>✨ Warna Aksen</span>
                   <div class="w-4 h-4 rounded border" :style="{backgroundColor: editingBlock.data.accentColor || '#f59e0b'}"></div>
                   <svg class="w-3 h-3 group-open/accent:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                 </summary>
                 <div class="mt-2 p-3 bg-white rounded border space-y-3">
                   <!-- Countdown specific -->
                   <template v-if="editingBlock.type === 'countdown'">
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Angka</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.numberColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.numberColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.numberColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Label</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.labelColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.labelColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.labelColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                   <!-- Statistics specific -->
                   <template v-else-if="editingBlock.type === 'statistics'">
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Angka Statistik</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.numberColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.numberColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.numberColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                   <!-- Benefits/FAQ - icon color -->
                   <template v-else-if="['benefits', 'faq'].includes(editingBlock.type)">
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Icon (Check/Arrow)</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.iconColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.iconColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.iconColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                   <!-- Comparison specific -->
                   <template v-else-if="editingBlock.type === 'comparison'">
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Check ✓</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.checkColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.checkColor" placeholder="#22c55e" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.checkColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Cross ✗</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.crossColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.crossColor" placeholder="#ef4444" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.crossColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                   <!-- Pricing badge -->
                   <template v-else-if="editingBlock.type === 'pricing'">
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Harga</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.priceColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.priceColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.priceColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Warna Badge Diskon</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.badgeColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.badgeColor" placeholder="#ef4444" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.badgeColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                   <!-- Generic accent -->
                   <template v-else>
                     <div>
                       <label class="block text-[10px] font-medium text-neutral-500 mb-1">Accent Color</label>
                       <div class="flex items-center gap-2">
                         <input type="color" v-model="editingBlock.data.accentColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                         <input type="text" v-model="editingBlock.data.accentColor" placeholder="Default" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                         <button @click="editingBlock.data.accentColor = undefined" class="text-[10px] text-red-500">×</button>
                       </div>
                     </div>
                   </template>
                 </div>
               </details>

               <!-- Testimonials specific -->
               <details v-if="editingBlock.type === 'testimonials'" class="group/test">
                 <summary class="flex items-center gap-2 text-xs font-medium text-neutral-700 cursor-pointer p-2 bg-neutral-100 rounded hover:bg-neutral-200 transition-colors">
                   <span>💬 Warna Testimonial</span>
                   <svg class="w-3 h-3 group-open/test:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
                 </summary>
                 <div class="mt-2 p-3 bg-white rounded border space-y-3">
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Card Background</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.cardBgColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.cardBgColor" placeholder="#f3f4f6" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.cardBgColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                   <div>
                     <label class="block text-[10px] font-medium text-neutral-500 mb-1">Star/Rating Color</label>
                     <div class="flex items-center gap-2">
                       <input type="color" v-model="editingBlock.data.starColor" class="h-7 w-8 rounded cursor-pointer border-0 p-0"/>
                       <input type="text" v-model="editingBlock.data.starColor" placeholder="#fbbf24" class="flex-1 px-2 py-1 border rounded text-xs uppercase"/>
                       <button @click="editingBlock.data.starColor = undefined" class="text-[10px] text-red-500">×</button>
                     </div>
                   </div>
                 </div>
               </details>
             </div>
          </details>

          <!-- Advanced: Custom CSS Override (NEW) -->
          <details class="group">
            <summary class="flex items-center gap-2 text-xs font-bold text-orange-600 cursor-pointer p-3 bg-gradient-to-br from-orange-50 to-amber-50 rounded-lg border border-orange-200 hover:bg-orange-100 transition-colors">
              <span>⚡ Advanced: Custom CSS Override</span>
              <svg class="w-3 h-3 group-open:rotate-180 transition-transform ml-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/></svg>
            </summary>
            
            <div class="mt-3 space-y-4 p-3 bg-neutral-900 rounded-lg">
              <!-- Info -->
              <div class="p-2 bg-blue-500/10 border border-blue-500/30 rounded text-blue-300 text-[10px]">
                ℹ️ CSS ini akan ditambahkan ke block ini saja (scoped). Gunakan untuk override styling default.
              </div>

              <!-- Custom CSS Editor -->
              <div>
                <label class="block text-xs font-medium text-neutral-300 mb-1">CSS Override</label>
                <textarea 
                  v-model="editingBlock.data.customCSS"
                  rows="8"
                  placeholder="/* Contoh Override */
h1 {
  text-transform: uppercase;
  letter-spacing: 2px;
}

.btn-primary {
  box-shadow: 0 10px 20px rgba(0,0,0,0.2);
}

/* Animasi */
img {
  transition: transform 0.3s;
}
img:hover {
  transform: scale(1.1);
}"
                  class="w-full px-3 py-2 bg-neutral-800 border border-neutral-700 rounded font-mono text-xs text-cyan-400 placeholder-neutral-500 focus:border-orange-500 focus:ring-1 focus:ring-orange-500"
                  spellcheck="false"
                ></textarea>
                <p class="text-[10px] text-neutral-500 mt-1">CSS akan di-scope otomatis dengan ID block (<code>#block-{{ editingBlock.id }}</code>).</p>
              </div>
              
              <!-- Action Buttons -->
              <div class="flex items-center gap-2 pt-2 border-t border-neutral-700">
                <button 
                  @click="editingBlock.data.customCSS = ''"
                  class="px-3 py-2 text-xs font-medium bg-red-500/20 text-red-400 rounded hover:bg-red-500/30"
                >
                  🗑️ Reset CSS
                </button>
              </div>
            </div>
          </details>

          <template v-if="editingBlock.type === 'hero'">
            <div v-if="editingBlock.variant === 'bio_profile'" class="bg-primary-50 p-3 rounded-lg border border-primary-100 mb-4">
               <label class="block text-xs font-bold text-primary-700 uppercase mb-2">👤 Profile Settings</label>
               <div class="space-y-3">
                 <div>
                   <label class="block text-sm font-medium mb-1">Profile Image URL</label>
                   <input v-model="editingBlock.data.profile_image" type="url" class="w-full px-3 py-2 border rounded-lg" placeholder="https://... (Avatar)"/>
                 </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Handle / Username</label>
                   <input v-model="editingBlock.data.handle" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="@username"/>
                 </div>
               </div>
            </div>

            <div>
              <label class="block text-sm font-medium mb-1">Headline / Name</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline / Bio</label>
              <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg"></textarea>
            </div>
            <div v-if="editingBlock.variant !== 'bio_profile'">
              <label class="block text-sm font-medium mb-1">Promo Text (Floating Badge)</label>
              <input v-model="editingBlock.data.promo_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Promo Berakhir dalam"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">{{ editingBlock.variant === 'bio_profile' ? 'Header Cover Image URL' : 'Background Image URL' }}</label>
              <input v-model="editingBlock.data.background_image" type="url" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
            </div>
            <div class="grid grid-cols-2 gap-4" v-if="editingBlock.variant !== 'bio_profile'">
              <div>
                <label class="block text-sm font-medium mb-1">CTA Text</label>
                <input v-model="editingBlock.data.cta_text" type="text" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">CTA Link</label>
                <input v-model="editingBlock.data.cta_link" type="text" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
            </div>
            <div v-if="editingBlock.variant !== 'bio_profile'">
              <label class="block text-sm font-medium mb-1">Trust Text (Bawah Tombol)</label>
              <input v-model="editingBlock.data.trust_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="✓ Akses Seumur Hidup • ✓ Garansi 30 Hari"/>
            </div>
          </template>

          <!-- Countdown Editor -->
          <template v-else-if="editingBlock.type === 'countdown'">
            <div>
              <label class="block text-sm font-medium mb-1">Label</label>
              <input v-model="editingBlock.data.label" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="⏰ Penawaran berakhir dalam:"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Countdown Target (Waktu Berakhir)</label>
              <input v-model="editingBlock.data.end_date" type="datetime-local" class="w-full px-3 py-2 border rounded-lg"/>
              <p class="text-xs text-neutral-500 mt-1">Jika tidak diisi, akan menggunakan "Countdown End Date" dari pengaturan campaign global.</p>
            </div>
          </template>

          <!-- Benefits Editor -->
          <template v-else-if="editingBlock.type === 'benefits'">
            <div>
              <label class="block text-sm font-medium mb-1">Title</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Items</label>
              <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="flex gap-2 mb-2">
                <input v-model="item.text" type="text" class="flex-1 px-3 py-2 border rounded-lg" placeholder="Benefit..."/>
                <button @click="editingBlock.data.items.splice(idx, 1)" class="px-3 py-2 text-red-500 hover:bg-red-50 rounded-lg">✕</button>
              </div>
              <button @click="editingBlock.data.items.push({icon: 'check', text: ''})" class="text-primary-600 text-sm font-medium">+ Tambah Item</button>
            </div>
          </template>

          <!-- Pricing Editor -->
          <template v-else-if="editingBlock.type === 'pricing'">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1">Harga Asli</label>
                <input v-model.number="editingBlock.data.original_price" type="number" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Harga Diskon</label>
                <input v-model.number="editingBlock.data.discount_price" type="number" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">CTA Text</label>
              <input v-model="editingBlock.data.cta_text" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div class="flex items-center gap-2 mb-3">
              <input type="checkbox" v-model="editingBlock.data.show_timer" class="rounded border-neutral-300"/>
              <label class="text-sm">Tampilkan timer</label>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Trust Text (Bawah Tombol)</label>
              <input v-model="editingBlock.data.trust_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Garansi 30 Hari Uang Kembali"/>
            </div>
          </template>

          <!-- Testimonials Editor -->
          <template v-else-if="editingBlock.type === 'testimonials'">
            <div>
              <label class="block text-sm font-medium mb-1">Title</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Testimonials</label>
              <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2">
                <div class="flex justify-between">
                  <input v-model="item.name" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Nama"/>
                  <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-500 text-sm ml-2">Hapus</button>
                </div>
                <textarea v-model="item.text" rows="2" class="w-full px-2 py-1 border rounded text-sm" placeholder="Testimonial..."></textarea>
              </div>
              <button @click="editingBlock.data.items.push({name: '', text: ''})" class="text-primary-600 text-sm font-medium">+ Tambah</button>
            </div>
          </template>

          <!-- FAQ Editor -->
          <template v-else-if="editingBlock.type === 'faq'">
            <div>
              <label class="block text-sm font-medium mb-1">Title</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">FAQ Items</label>
              <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2">
                <div class="flex justify-between">
                  <input v-model="item.question" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Pertanyaan"/>
                  <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-500 text-sm ml-2">Hapus</button>
                </div>
                <textarea v-model="item.answer" rows="2" class="w-full px-2 py-1 border rounded text-sm" placeholder="Jawaban..."></textarea>
              </div>
              <button @click="editingBlock.data.items.push({question: '', answer: ''})" class="text-primary-600 text-sm font-medium">+ Tambah</button>
            </div>
          </template>

          <!-- Instructor Editor -->
          <template v-else-if="editingBlock.type === 'instructor'">
            <div class="flex items-center gap-2 mb-4">
              <input type="checkbox" v-model="editingBlock.data.auto_fill" class="rounded border-neutral-300"/>
              <label class="text-sm">Ambil dari data instruktur kursus</label>
            </div>
            <template v-if="!editingBlock.data.auto_fill">
              <div>
                <label class="block text-sm font-medium mb-1">Nama</label>
                <input v-model="editingBlock.data.name" type="text" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Bio</label>
                <textarea v-model="editingBlock.data.bio" rows="3" class="w-full px-3 py-2 border rounded-lg"></textarea>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Avatar URL</label>
                <input v-model="editingBlock.data.profile_image" type="url" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
              </div>
            </template>
          </template>

          <!-- Video Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'video'">
            <div>
              <label class="block text-sm font-medium mb-1">Title (Opsional)</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Lihat Apa yang Akan Anda Pelajari"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subtitle (Opsional)</label>
              <input v-model="editingBlock.data.subtitle" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Preview materi kursus"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">YouTube URL *</label>
              <input v-model="editingBlock.data.youtube_url" type="url" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="https://www.youtube.com/watch?v=..."/>
              <p class="text-xs text-neutral-500 mt-1">Dukung format: youtube.com/watch?v=..., youtu.be/...</p>
            </div>
          </template>

            <!-- Hero Gen Z Editor -->
            <template v-else-if="editingBlock.type === 'hero_gen_z'">
              <div class="p-3 bg-purple-50 rounded-lg border border-purple-100 mb-4">
                <label class="block text-xs font-bold text-purple-700 uppercase mb-2">⚡ Gen Z Settings</label>
                <div>
                   <label class="block text-sm font-medium mb-1">Badge Text (Neon)</label>
                   <input v-model="editingBlock.data.badge" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="WEBINAR EKSKLUSIF"/>
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Headline</label>
                <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="BELAJAR KEPABEANAN & CUKAI"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Subheadline</label>
                <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg" placeholder="Deskripsi singkat..."></textarea>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Promo Text (Above Countdown)</label>
                <input v-model="editingBlock.data.promo_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Promo Early Bird Berakhir Dalam:"/>
              </div>
              <div>
                 <!-- Link to Webinar removed as per request to rely on global -->
              </div>
            </template>

            <!-- CLEAN / TIKTOK EDITORS -->
            <template v-else-if="editingBlock.type === 'hero_clean'">
                <div>
                   <label class="block text-sm font-medium mb-1">Headline</label>
                   <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Subheadline</label>
                   <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg"></textarea>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Badge Text</label>
                   <input v-model="editingBlock.data.badge" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                <div>
                   <input v-model="editingBlock.data.image_url" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
                   <p class="text-xs text-neutral-400 mt-1">Recommended size: 800x600px</p>
                </div>
                <!-- Style Overrides -->
                <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t">
                   <div>
                      <label class="block text-sm font-medium mb-1">Button Color</label>
                      <input v-model="editingBlock.data.button_color" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="bg-blue-600 or hex"/>
                   </div>
                   <div>
                       <label class="block text-sm font-medium mb-1">Font Family</label>
                       <select v-model="editingBlock.data.font_family" class="w-full px-3 py-2 border rounded-lg">
                          <option value="Inter">Inter</option>
                          <option value="Poppins">Poppins</option>
                          <option value="Roboto">Roboto</option>
                           <option value="Outfit">Outfit</option>
                       </select>
                   </div>
                </div>
            </template>

            <template v-else-if="editingBlock.type === 'hero_pro'">
                <div>
                   <label class="block text-sm font-medium mb-1">Headline</label>
                   <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Subheadline</label>
                   <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg"></textarea>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Badge Text</label>
                   <input v-model="editingBlock.data.badge" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                <div class="p-3 bg-amber-50 rounded-lg border border-amber-200 mt-4">
                   <h4 class="font-semibold text-amber-800 mb-3">📷 Foto Pembicara (Hero)</h4>
                   <div>
                      <label class="block text-sm font-medium mb-1">Nama Pembicara</label>
                      <input v-model="editingBlock.data.speaker_name" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                   </div>
                   <div class="mt-2">
                      <label class="block text-sm font-medium mb-1">Jabatan Pembicara</label>
                      <input v-model="editingBlock.data.speaker_title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                   </div>
                   <div class="mt-2">
                      <label class="block text-sm font-medium mb-1">URL Foto Pembicara</label>
                      <input v-model="editingBlock.data.speaker_image" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
                      <p class="text-xs text-neutral-400 mt-1">Recommended: 400x500px (Portrait)</p>
                   </div>
                   <div v-if="editingBlock.data.speaker_image" class="mt-3">
                      <p class="text-xs text-neutral-500 mb-1">Preview:</p>
                      <img :src="editingBlock.data.speaker_image" class="w-24 h-24 object-cover rounded-lg border"/>
                   </div>
                </div>
            </template>

             <template v-else-if="editingBlock.type === 'speaker_clean'">
                <div>
                   <label class="block text-sm font-medium mb-1">Nama Pembicara</label>
                   <input v-model="editingBlock.data.name" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Title / Jabatan</label>
                   <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                <div>
                   <label class="block text-sm font-medium mb-1">Bio</label>
                   <textarea v-model="editingBlock.data.bio" rows="4" class="w-full px-3 py-2 border rounded-lg"></textarea>
                </div>
                <div>
                   <label class="block text-sm font-medium mb-1">Speaker Photo URL</label>
                   <input v-model="editingBlock.data.image_url" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
                   <p class="text-xs text-neutral-400 mt-1">Recommended size: 800x1000px (Portrait)</p>
                </div>
                 <!-- Simple Stats Editor -->
                 <label class="block text-sm font-medium mt-4 mb-2">Stats (2 Items)</label>
                 <div v-for="(stat, idx) in editingBlock.data.stats" :key="idx" class="flex gap-2 mb-2">
                     <input v-model="stat.icon" class="w-12 px-2 py-1 border rounded" placeholder="Icon"/>
                     <input v-model="stat.label" class="flex-1 px-2 py-1 border rounded" placeholder="Label"/>
                 </div>
            </template>

            <template v-else-if="editingBlock.type === 'features_clean'">
                 <div>
                   <label class="block text-sm font-medium mb-1">Section Title</label>
                   <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                <div class="mt-4">
                     <label class="block text-sm font-medium mb-2">Features Items</label>
                     <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-3 bg-neutral-50">
                        <input v-model="item.title" class="w-full mb-2 px-2 py-1 border rounded bg-white text-sm font-bold" placeholder="Title"/>
                        <textarea v-model="item.text" rows="2" class="w-full mb-2 px-2 py-1 border rounded bg-white text-sm" placeholder="Description"></textarea>
                        <input v-model="item.icon" class="w-12 px-2 py-1 border rounded bg-white text-center" placeholder="Icon"/>
                     </div>
                </div>
            </template>

            <template v-else-if="editingBlock.type === 'cta_clean'">
                <div>
                   <label class="block text-sm font-medium mb-1">Headline</label>
                   <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                 <div>
                   <label class="block text-sm font-medium mb-1">Subheadline</label>
                   <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg"></textarea>
                </div>
                <!-- Style Overrides -->
                <div class="grid grid-cols-2 gap-4 mt-4 pt-4 border-t">
                   <div>
                      <label class="block text-sm font-medium mb-1">Button Color</label>
                      <input v-model="editingBlock.data.button_color" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="bg-blue-600 or hex"/>
                   </div>
                   <div>
                       <label class="block text-sm font-medium mb-1">Font Family</label>
                       <select v-model="editingBlock.data.font_family" class="w-full px-3 py-2 border rounded-lg">
                          <option value="Inter">Inter</option>
                          <option value="Poppins">Poppins</option>
                          <option value="Roboto">Roboto</option>
                           <option value="Outfit">Outfit</option>
                       </select>
                   </div>
                </div>
            </template>

            <!-- Features Gen Z Editor -->
            <template v-else-if="editingBlock.type === 'features_gen_z'">
              <div>
                <label class="block text-sm font-medium mb-1">Section Title</label>
                <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Yang Bakal Kamu DAPETIN"/>
              </div>
              <div class="mt-4">
                <label class="block text-sm font-medium mb-2">Grid Items (4 Cards)</label>
                <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-3 bg-neutral-900 text-white border-neutral-700">
                   <div class="flex gap-2 mb-2">
                      <div class="w-10">
                         <label class="text-[10px] text-neutral-400">Icon</label>
                         <input v-model="item.icon" type="text" class="w-full px-1 py-1 bg-neutral-800 border border-neutral-600 rounded text-center" placeholder="📘"/>
                      </div>
                      <div class="flex-1">
                         <label class="text-[10px] text-neutral-400">Title</label>
                         <input v-model="item.title" type="text" class="w-full px-2 py-1 bg-neutral-800 border border-neutral-600 rounded text-white" placeholder="Judul Fitur"/>
                      </div>
                      <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-400 hover:text-red-300">✕</button>
                   </div>
                   <div>
                      <label class="text-[10px] text-neutral-400">Description</label>
                      <textarea v-model="item.desc" rows="2" class="w-full px-2 py-1 bg-neutral-800 border border-neutral-600 rounded text-xs text-neutral-300"></textarea>
                   </div>
                   <div class="grid grid-cols-2 gap-2 mt-2">
                      <div>
                        <label class="text-[10px] text-neutral-400">Text Color Class</label>
                        <input v-model="item.color" type="text" class="w-full px-2 py-1 bg-neutral-800 border border-neutral-600 rounded text-xs text-green-400 font-mono" placeholder="text-purple-400"/>
                      </div>
                      <div>
                        <label class="text-[10px] text-neutral-400">Bg Color Class</label>
                        <input v-model="item.bg" type="text" class="w-full px-2 py-1 bg-neutral-800 border border-neutral-600 rounded text-xs text-green-400 font-mono" placeholder="bg-purple-500/10"/>
                      </div>
                   </div>
                </div>
                <button @click="editingBlock.data.items = editingBlock.data.items || []; editingBlock.data.items.push({title:'New Feature', desc:'', icon:'✨', color:'text-white', bg:'bg-white/10'})" class="text-purple-600 text-sm font-bold border border-purple-200 w-full py-2 rounded-lg hover:bg-purple-50">+ Add Gen Z Card</button>
              </div>
            </template>

            <!-- CTA Gen Z Editor -->
            <template v-else-if="editingBlock.type === 'cta_gen_z'">
              <div>
                <label class="block text-sm font-medium mb-1">Headline</label>
                <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="SIAP JADI AHLI KEPABEANAN?"/>
              </div>
               <div>
                <label class="block text-sm font-medium mb-1">Subheadline</label>
                <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Daftar sekarang dan amankan slot early bird!"/>
              </div>
              <div class="grid grid-cols-2 gap-4">
                 <div>
                    <label class="block text-sm font-medium mb-1">Price (Rp)</label>
                    <input v-model.number="editingBlock.data.price" type="number" class="w-full px-3 py-2 border rounded-lg"/>
                 </div>
                 <div>
                    <label class="block text-sm font-medium mb-1">Original Price (Rp)</label>
                    <input v-model.number="editingBlock.data.original_price" type="number" class="w-full px-3 py-2 border rounded-lg"/>
                 </div>
              </div>
              <div>
                 <label class="block text-sm font-medium mb-1">Button Text</label>
                 <input v-model="editingBlock.data.button_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="DAFTAR SEKARANG 🚀"/>
              </div>
              <div class="p-3 bg-neutral-50 rounded border mt-2">
                 <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Form Integration</label>
                 <div class="flex items-center justify-between text-sm">
                    <span class="text-neutral-600">Whatsapp Notification</span>
                    <span class="px-2 py-0.5 bg-green-100 text-green-700 rounded text-xs font-bold">ACTIVE</span>
                 </div>
              </div>
            </template>

          <!-- CTA Banner Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'cta_banner'">
            <div>
              <label class="block text-sm font-medium mb-1">Headline</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="🚀 Siap Untuk Memulai?"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="Jangan lewatkan kesempatan untuk meningkatkan skill Anda"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">CTA Text</label>
              <input v-model="editingBlock.data.cta_text" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="Daftar Sekarang"/>
            </div>
          </template>

          <!-- Statistics Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'statistics'">
            <div>
              <label class="block text-sm font-medium mb-1">Title (Opsional)</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="Mengapa Memilih Kami?"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Statistics</label>
              <div v-for="(stat, idx) in editingBlock.data.items" :key="idx" class="flex gap-2 mb-2">
                <input v-model="stat.value" type="text" class="w-20 px-2 py-1 border rounded text-sm" placeholder="5000"/>
                <input v-model="stat.suffix" type="text" class="w-12 px-2 py-1 border rounded text-sm" placeholder="+"/>
                <input v-model="stat.label" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Siswa"/>
                <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-500 text-sm">✕</button>
              </div>
              <button @click="editingBlock.data.items.push({value: '', suffix: '', label: ''})" class="text-primary-600 text-sm font-medium">+ Tambah Statistik</button>
            </div>
          </template>

          <!-- Bonus Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'bonus'">
            <div>
              <label class="block text-sm font-medium mb-1">Title</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="Bonus Eksklusif untuk Anda!"/>
            </div>
            <div>
               <label class="block text-sm font-medium mb-1">Tampilan</label>
               <div class="flex gap-2 p-1 bg-neutral-100 rounded-lg">
                  <button @click="editingBlock.data.displayMode = 'list'" class="flex-1 py-1.5 rounded-md text-sm transition-all" :class="!editingBlock.data.displayMode || editingBlock.data.displayMode === 'list' ? 'bg-white shadow text-neutral-900 font-medium' : 'text-neutral-500 hover:text-neutral-700'">List</button>
                  <button @click="editingBlock.data.displayMode = 'card'" class="flex-1 py-1.5 rounded-md text-sm transition-all" :class="editingBlock.data.displayMode === 'card' ? 'bg-white shadow text-neutral-900 font-medium' : 'text-neutral-500 hover:text-neutral-700'">Cards</button>
               </div>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Bonus Items</label>
              <div v-for="(bonus, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2">
                <div class="flex gap-2">
                  <input v-model="bonus.emoji" type="text" class="w-12 px-2 py-1 border rounded text-center" placeholder="🎁"/>
                  <input v-model="bonus.title" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Nama Bonus"/>
                  <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-500 text-sm">Hapus</button>
                </div>
                <div class="flex gap-2">
                  <input v-model.number="bonus.value" type="number" class="w-32 px-2 py-1 border rounded text-sm" placeholder="Nilai (Rp)"/>
                  <input v-model="bonus.description" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Deskripsi singkat"/>
                </div>
              </div>
              <button @click="editingBlock.data.items.push({emoji: '🎁', title: '', value: 0, description: ''})" class="text-primary-600 text-sm font-medium">+ Tambah Bonus</button>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Total Nilai Bonus</label>
              <input v-model.number="editingBlock.data.total_value" type="number" class="w-full px-3 py-2 border rounded-lg" placeholder="250000"/>
            </div>
          </template>

          <!-- Curriculum Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'curriculum'">
            <div>
              <label class="block text-sm font-medium mb-1">Title</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="📚 Materi yang Akan Anda Pelajari"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subtitle</label>
              <input v-model="editingBlock.data.subtitle" type="text" class="w-full px-3 py-2 border rounded-lg" 
                     placeholder="Kurikulum lengkap dan terstruktur"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Modules</label>
              <div v-for="(mod, idx) in editingBlock.data.modules" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2">
                <div class="flex gap-2 items-center">
                  <span class="w-6 h-6 bg-primary-100 text-primary-700 rounded text-sm flex items-center justify-center font-bold">{{ idx + 1 }}</span>
                  <input v-model="mod.title" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Judul Modul"/>
                  <button @click="editingBlock.data.modules.splice(idx, 1)" class="text-red-500 text-sm">Hapus</button>
                </div>
                <input v-model.number="mod.lessons_count" type="number" class="w-32 px-2 py-1 border rounded text-sm" placeholder="Jumlah pelajaran"/>
              </div>
              <button @click="editingBlock.data.modules.push({title: '', lessons_count: 0, lessons: []})" class="text-primary-600 text-sm font-medium">+ Tambah Modul</button>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium mb-1">Total Pelajaran</label>
                <input v-model.number="editingBlock.data.total_lessons" type="number" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">Total Jam</label>
                <input v-model.number="editingBlock.data.total_hours" type="number" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
            </div>
          </template>

          <!-- Gallery Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'gallery'">
            <div>
              <label class="block text-sm font-medium mb-1">Title (Opsional)</label>
              <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Galeri Kegiatan"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-2">Images</label>
              <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2 bg-neutral-50/50">
                <div class="flex items-center gap-2">
                    <span class="text-xs bg-neutral-200 px-2 py-1 rounded">Image {{ idx + 1 }}</span>
                    <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-500 text-xs ml-auto">Hapus</button>
                </div>
                <input v-model="item.image" type="url" class="w-full px-2 py-1 border rounded text-sm" placeholder="URL Gambar (https://...)"/>
                <input v-model="item.caption" type="text" class="w-full px-2 py-1 border rounded text-sm" placeholder="Caption (optional)"/>
              </div>
              <button @click="editingBlock.data.items.push({image: '', caption: ''})" class="text-primary-600 text-sm font-medium border border-dashed border-primary-300 w-full py-2 rounded-lg hover:bg-primary-50">+ Tambah Gambar</button>
            </div>
          </template>

          <!-- Achievement Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'achievement'">
            <div>
              <label class="block text-sm font-medium mb-2">Badges</label>
              <div v-for="(item, idx) in editingBlock.data.items" :key="idx" class="flex gap-2 mb-2 items-center">
                <input v-model="item.emoji" type="text" class="w-12 px-2 py-2 border rounded-lg text-center" placeholder="🏆"/>
                <input v-model="item.text" type="text" class="flex-1 px-3 py-2 border rounded-lg" placeholder="Teks (mis: 100% Puas)"/>
                <button @click="editingBlock.data.items.splice(idx, 1)" class="text-red-400 hover:text-red-600 px-2">✕</button>
              </div>
              <button @click="editingBlock.data.items.push({emoji: '🏅', text: ''})" class="text-primary-600 text-sm font-medium mt-1">+ Tambah Badge</button>
            </div>
          </template>
          <!-- Comparison Table Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'comparison'">
            <div>
               <label class="block text-sm font-medium mb-1">Title</label>
               <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Bandingkan Paket"/>
            </div>
            <div>
               <label class="block text-sm font-medium mb-1">Header Label (Kolom 1)</label>
               <input v-model="editingBlock.data.headers[0]" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Fitur"/>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                 <label class="block text-sm font-medium mb-1">Label Opsi A</label>
                 <input v-model="editingBlock.data.headers[1]" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Gratis"/>
              </div>
              <div>
                 <label class="block text-sm font-medium mb-1">Label Opsi B</label>
                 <input v-model="editingBlock.data.headers[2]" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Premium"/>
              </div>
            </div>
            
            <div>
               <label class="block text-sm font-medium mb-2">Rows</label>
               <div v-for="(row, idx) in editingBlock.data.rows" :key="idx" class="p-3 border rounded-lg mb-2 space-y-2 bg-neutral-50/50">
                  <input v-model="row.feature" type="text" class="w-full px-2 py-1 border rounded text-sm font-medium" placeholder="Nama Fitur"/>
                  <div class="grid grid-cols-2 gap-4">
                     <select v-model="row.val_a" class="w-full px-2 py-1 border rounded text-sm">
                        <option :value="true">✅ Ya</option>
                        <option :value="false">❌ Tidak</option>
                        <option value="text">Teks...</option>
                     </select>
                     <select v-model="row.val_b" class="w-full px-2 py-1 border rounded text-sm">
                        <option :value="true">✅ Ya</option>
                        <option :value="false">❌ Tidak</option>
                        <option value="text">Teks...</option>
                     </select>
                  </div>
                  <div class="grid grid-cols-2 gap-4" v-if="row.val_a === 'text' || row.val_b === 'text'">
                     <input v-if="row.val_a === 'text'" v-model="row.text_a" type="text" class="w-full px-2 py-1 border rounded text-sm" placeholder="Teks A"/>
                     <input v-if="row.val_b === 'text'" v-model="row.text_b" type="text" class="w-full px-2 py-1 border rounded text-sm" placeholder="Teks B"/>
                  </div>
                  <button @click="editingBlock.data.rows.splice(idx, 1)" class="text-red-500 text-xs ml-auto block">Hapus Baris</button>
               </div>
               <button @click="editingBlock.data.rows.push({feature: '', val_a: false, val_b: true})" class="text-primary-600 text-sm font-medium border border-dashed border-primary-300 w-full py-2 rounded-lg hover:bg-primary-50">+ Tambah Baris Perbandingan</button>
            </div>
          </template>

          <!-- Floating Chat Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'floating_chat'">
             <div>
                <label class="block text-sm font-medium mb-1">WhatsApp Number</label>
                <input v-model="editingBlock.data.phone" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="628123456789"/>
                <p class="text-xs text-neutral-500 mt-1">Gunakan format internasional tanpa '+' (contoh: 628...)</p>
             </div>
             <div>
                <label class="block text-sm font-medium mb-1">Pesan Otomatis</label>
                <textarea v-model="editingBlock.data.message" rows="2" class="w-full px-3 py-2 border rounded-lg" placeholder="Halo, saya mau tanya tentang kursus ini..."></textarea>
             </div>
             <div>
                <label class="block text-sm font-medium mb-1">Label Tombol (Opsional)</label>
                <input v-model="editingBlock.data.label" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Chat Kami"/>
             </div>
          </template>

          <!-- Social Proof Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'social_proof'">
             <div>
                <label class="block text-sm font-medium mb-1">Template Teks</label>
                <input v-model="editingBlock.data.text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="orang sedang melihat halaman ini"/>
                <p class="text-xs text-neutral-500 mt-1">Gunakan kata-kata yang menarik. Angka akan muncul sebelum teks ini.</p>
             </div>
             <div class="grid grid-cols-2 gap-4 mt-3">
                <div>
                   <label class="block text-sm font-medium mb-1">Min Count</label>
                   <input v-model.number="editingBlock.data.min_count" type="number" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
                <div>
                   <label class="block text-sm font-medium mb-1">Max Count</label>
                   <input v-model.number="editingBlock.data.max_count" type="number" class="w-full px-3 py-2 border rounded-lg"/>
                </div>
             </div>
             <div class="mt-3">
                <label class="block text-sm font-medium mb-1">Durasi Update (detik)</label>
                <input v-model.number="editingBlock.data.interval" type="number" class="w-full px-3 py-2 border rounded-lg" placeholder="5"/>
                <p class="text-xs text-neutral-500 mt-1">Seberapa sering angka berubah (simulasi live).</p>
             </div>
          </template>

          <!-- Link Button Editor (NEW) -->
          <template v-else-if="editingBlock.type === 'link_button'">
             <div class="space-y-4">
               <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Label Tombol</label>
                  <input v-model="editingBlock.data.label" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Klik Disini"/>
               </div>
               <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Link URL</label>
                  <input v-model="editingBlock.data.url" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
               </div>
               <div class="grid grid-cols-2 gap-4">
                 <div>
                    <label class="block text-sm font-medium text-neutral-700 mb-1">Icon (Emoji)</label>
                    <input v-model="editingBlock.data.icon" type="text" class="w-full px-3 py-2 border rounded-lg text-center" placeholder="🔗"/>
                 </div>
                 <div>
                     <label class="block text-sm font-medium text-neutral-700 mb-1">Style</label>
                     <select v-model="editingBlock.data.style" class="w-full px-3 py-2 border rounded-lg">
                        <option value="solid">Solid Color</option>
                        <option value="outline">Outline</option>
                        <option value="glass">Glassmorphism</option>
                     </select>
                 </div>
               </div>
             </div>
          </template>

          <!-- HEALING BLOCKS EDITORS (P.U.L.I.H Theme) -->
          
          <!-- Hero Healing Editor -->
          <template v-else-if="editingBlock.type === 'hero_healing'">
            <div>
              <label class="block text-sm font-medium mb-1">Headline Utama</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="P.U.L.I.H"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Seni Berdamai dengan Luka Batin"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Badge Text</label>
              <input v-model="editingBlock.data.badge" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="FREE WEBINAR"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Platform</label>
              <input v-model="editingBlock.data.platform" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Online via Zoom"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Button Text</label>
              <input v-model="editingBlock.data.button_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="DAFTAR GRATIS SEKARANG"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">CTA Note</label>
              <input v-model="editingBlock.data.cta_note" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Kuota terbatas, amankan tempatmu sekarang!"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Background Image URL</label>
              <input v-model="editingBlock.data.image_url" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
            </div>
          </template>

          <!-- Problem Healing Editor -->
          <template v-else-if="editingBlock.type === 'problem_healing'">
            <div>
              <label class="block text-sm font-medium mb-1">Eyebrow Text</label>
              <input v-model="editingBlock.data.eyebrow" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Apakah kamu merasakan ini?"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Headline</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Bukan Masalah Hari Ini..."/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Reflection Quote</label>
              <textarea v-model="editingBlock.data.reflection" rows="3" class="w-full px-3 py-2 border rounded-lg"></textarea>
            </div>
            <div class="mt-4">
              <label class="block text-sm font-medium mb-2">Problem Cards</label>
              <div v-for="(problem, idx) in editingBlock.data.problems" :key="idx" class="p-3 border rounded-lg mb-3 bg-amber-50">
                <div class="flex gap-2 mb-2">
                  <input v-model="problem.emoji" type="text" class="w-12 px-2 py-1 border rounded text-center" placeholder="😤"/>
                  <input v-model="problem.title" type="text" class="flex-1 px-2 py-1 border rounded" placeholder="Title"/>
                  <button @click="editingBlock.data.problems.splice(idx, 1)" class="text-red-500">✕</button>
                </div>
                <textarea v-model="problem.description" rows="2" class="w-full px-2 py-1 border rounded text-sm" placeholder="Description"></textarea>
              </div>
              <button @click="editingBlock.data.problems = editingBlock.data.problems || []; editingBlock.data.problems.push({emoji:'❓', title:'', description:''})" class="text-amber-700 text-sm font-medium border border-dashed border-amber-300 w-full py-2 rounded-lg hover:bg-amber-50">+ Add Problem Card</button>
            </div>
          </template>

          <!-- Solution Healing Editor -->
          <template v-else-if="editingBlock.type === 'solution_healing'">
            <div>
              <label class="block text-sm font-medium mb-1">Eyebrow Text</label>
              <input v-model="editingBlock.data.eyebrow" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Di Webinar P.U.L.I.H"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Headline</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Yang Akan Kamu Pelajari"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Key Message</label>
              <textarea v-model="editingBlock.data.key_message" rows="3" class="w-full px-3 py-2 border rounded-lg"></textarea>
            </div>
            <div class="mt-4">
              <label class="block text-sm font-medium mb-2">Learning Points</label>
              <div v-for="(item, idx) in editingBlock.data.learnings" :key="idx" class="p-3 border rounded-lg mb-3 bg-green-50">
                <div class="flex gap-2 mb-2">
                  <input v-model="item.icon" type="text" class="w-12 px-2 py-1 border rounded text-center" placeholder="🧠"/>
                  <input v-model="item.title" type="text" class="flex-1 px-2 py-1 border rounded" placeholder="Title"/>
                  <button @click="editingBlock.data.learnings.splice(idx, 1)" class="text-red-500">✕</button>
                </div>
                <textarea v-model="item.description" rows="2" class="w-full px-2 py-1 border rounded text-sm" placeholder="Description"></textarea>
              </div>
              <button @click="editingBlock.data.learnings = editingBlock.data.learnings || []; editingBlock.data.learnings.push({icon:'✨', title:'', description:''})" class="text-green-700 text-sm font-medium border border-dashed border-green-300 w-full py-2 rounded-lg hover:bg-green-50">+ Add Learning Point</button>
            </div>
          </template>

          <!-- Speaker Healing Editor -->
          <template v-else-if="editingBlock.type === 'speaker_healing'">
            <div class="p-3 bg-amber-50 rounded-lg border border-amber-200">
              <h4 class="font-semibold text-amber-800 mb-3">👤 Profil Instruktur</h4>
              <div>
                <label class="block text-sm font-medium mb-1">Eyebrow Text</label>
                <input v-model="editingBlock.data.eyebrow" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Dipandu Oleh"/>
              </div>
              <div class="mt-2">
                <label class="block text-sm font-medium mb-1">Headline</label>
                <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Tentang Instruktur"/>
              </div>
              <div class="mt-2">
                <label class="block text-sm font-medium mb-1">Nama Instruktur</label>
                <input v-model="editingBlock.data.name" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Fransiskus Indra Pratama"/>
              </div>
              <div class="mt-2">
                <label class="block text-sm font-medium mb-1">Jabatan / Title</label>
                <input v-model="editingBlock.data.title" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Coach, Hypnotherapist & Trainer"/>
              </div>
              <div class="mt-2">
                <label class="block text-sm font-medium mb-1">Bio</label>
                <textarea v-model="editingBlock.data.bio" rows="4" class="w-full px-3 py-2 border rounded-lg"></textarea>
              </div>
              <div class="mt-2">
                <label class="block text-sm font-medium mb-1">Foto Instruktur (URL)</label>
                <input v-model="editingBlock.data.image_url" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
                <p class="text-xs text-neutral-400 mt-1">Recommended: 600x600px (Square/Portrait)</p>
              </div>
              <div v-if="editingBlock.data.image_url" class="mt-3">
                <p class="text-xs text-neutral-500 mb-1">Preview:</p>
                <img :src="editingBlock.data.image_url" class="w-24 h-24 object-cover rounded-lg border"/>
              </div>
            </div>
            <div class="mt-4">
              <label class="block text-sm font-medium mb-2">Sertifikasi / Credentials</label>
              <div v-for="(cred, idx) in editingBlock.data.credentials" :key="idx" class="flex gap-2 mb-2">
                <input v-model="editingBlock.data.credentials[idx]" type="text" class="flex-1 px-2 py-1 border rounded text-sm" placeholder="Certified Hypnotherapist"/>
                <button @click="editingBlock.data.credentials.splice(idx, 1)" class="text-red-500">✕</button>
              </div>
              <button @click="editingBlock.data.credentials = editingBlock.data.credentials || []; editingBlock.data.credentials.push('')" class="text-amber-700 text-sm font-medium border border-dashed border-amber-300 w-full py-2 rounded-lg hover:bg-amber-50">+ Add Credential</button>
            </div>
          </template>

          <!-- CTA Healing Editor -->
          <template v-else-if="editingBlock.type === 'cta_healing'">
            <div>
              <label class="block text-sm font-medium mb-1">Headline</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Siap Memulai Perjalanan Pemulihan?"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <input v-model="editingBlock.data.subheadline" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Daftar sekarang untuk mengamankan tempatmu..."/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Button Text</label>
              <input v-model="editingBlock.data.button_text" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="DAFTAR GRATIS SEKARANG"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Button Color</label>
              <input v-model="editingBlock.data.button_color" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="#4a4540 or gradient"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Privacy Note</label>
              <input v-model="editingBlock.data.privacy_note" type="text" class="w-full px-3 py-2 border rounded-lg" placeholder="Data kamu aman dan tidak akan dibagikan ke pihak lain."/>
            </div>
            <div class="p-3 bg-neutral-50 rounded border mt-3">
              <label class="block text-xs font-bold text-neutral-500 uppercase mb-2">Form Integration</label>
              <div class="flex items-center justify-between text-sm">
                <span class="text-neutral-600">WhatsApp Notification</span>
                <span class="px-2 py-0.5 bg-green-100 text-green-700 rounded text-xs font-bold">ACTIVE</span>
              </div>
            </div>
          </template>

        </div>

        <div class="flex justify-end gap-3 p-4 border-t">
          <button @click="editingBlock = null" class="px-4 py-2 border rounded-lg hover:bg-neutral-50">Tutup</button>
        </div>
      </div>
    </div>



    <!-- Template Selection Modal (or Initial Step) -->
    <div v-if="showTemplateSelector" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
       <div class="bg-white rounded-2xl w-full max-w-5xl h-[80vh] flex flex-col shadow-2xl">
          <div class="p-6 border-b flex justify-between items-center">
             <div>
               <h2 class="text-xl font-bold text-neutral-900">Pilih Template Campaign</h2>
               <p class="text-sm text-neutral-500">Mulai dengan desain profesional yang siap pakai.</p>
             </div>
             <button @click="showTemplateSelector = false" class="p-2 hover:bg-neutral-100 rounded-full">
               <svg class="w-6 h-6 text-neutral-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
               </svg>
             </button>
          </div>
          
          <div class="flex-1 overflow-y-auto p-6 bg-neutral-50">
             <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <div 
                  v-for="tpl in templates" 
                  :key="tpl.id"
                  @click="applyTemplate(tpl)"
                  class="group bg-white rounded-xl border border-neutral-200 overflow-hidden cursor-pointer hover:shadow-lg hover:border-primary-500 transition-all relative"
                >
                   <div class="h-40 bg-neutral-100 flex items-center justify-center text-4xl group-hover:scale-105 transition-transform">
                      {{ tpl.thumbnail }}
                   </div>
                   <div class="p-5">
                      <div class="flex justify-between items-start mb-2">
                         <span class="text-xs font-bold px-2 py-1 rounded-full uppercase tracking-wider" 
                           :class="{
                             'bg-blue-100 text-blue-700': tpl.category === 'ecourse',
                             'bg-red-100 text-red-700': tpl.category === 'webinar',
                             'bg-gray-100 text-gray-700': tpl.category === 'bio',
                             'bg-green-100 text-green-700': tpl.category === 'product',
                             'bg-purple-100 text-purple-700': tpl.category === 'service'
                           }">
                           {{ tpl.category }}
                         </span>
                      </div>
                      <h3 class="font-bold text-neutral-900 mb-1 group-hover:text-primary-600">{{ tpl.name }}</h3>
                      <p class="text-sm text-neutral-500 line-clamp-2">{{ tpl.description }}</p>
                   </div>
                   
                   <!-- Hover Overlay -->
                   <div class="absolute inset-x-0 bottom-0 p-4 bg-gradient-to-t from-white via-white to-transparent translate-y-full group-hover:translate-y-0 transition-transform duration-300">
                      <button class="w-full py-2 bg-primary-600 text-white font-semibold rounded-lg shadow-lg">Gunakan Template Ini</button>
                   </div>
                </div>
             </div>
          </div>
       </div>
    </div>
    <!-- Add Block Modal (Lynk Style) -->
    <div v-if="showAddBlockModal" class="fixed inset-0 bg-black/50 z-[60] flex items-center justify-center p-4">
       <div class="bg-white rounded-2xl w-full max-w-4xl h-[80vh] flex flex-col shadow-2xl overflow-hidden">
          <div class="p-4 border-b flex justify-between items-center bg-white">
             <h3 class="font-bold text-lg">Add New Block</h3>
             <button @click="showAddBlockModal = false" class="p-2 hover:bg-neutral-100 rounded-full">✕</button>
          </div>
          
          <div class="flex-1 overflow-y-auto bg-neutral-50 p-6">
             <div class="mb-8">
                <h4 class="text-sm font-bold text-neutral-500 uppercase tracking-wider mb-4">Basic Blocks</h4>
                <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                   <button v-for="item in blockCategories.basic" :key="item.type" @click="addBlock(item.type)"
                           class="bg-white p-4 rounded-xl border border-neutral-200 shadow-sm hover:shadow-md hover:border-primary-500 transition-all text-left group">
                      <div class="w-10 h-10 rounded-lg bg-neutral-100 flex items-center justify-center text-2xl mb-3 group-hover:bg-primary-50 group-hover:scale-110 transition-transform">
                         {{ item.icon }}
                      </div>
                      <div class="font-bold text-neutral-900 group-hover:text-primary-600">{{ item.label }}</div>
                      <div class="text-xs text-neutral-500 mt-1">{{ item.desc }}</div>
                   </button>
                </div>
             </div>
             
             <div>
                <h4 class="text-sm font-bold text-neutral-500 uppercase tracking-wider mb-4">Monetization & Marketing</h4>
                <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                   <button v-for="item in blockCategories.monetization" :key="item.type" @click="addBlock(item.type)"
                           class="bg-white p-4 rounded-xl border border-neutral-200 shadow-sm hover:shadow-md hover:border-green-500 transition-all text-left group">
                      <div class="w-10 h-10 rounded-lg bg-green-50 flex items-center justify-center text-2xl mb-3 group-hover:scale-110 transition-transform">
                         {{ item.icon }}
                      </div>
                      <div class="font-bold text-neutral-900 group-hover:text-green-600">{{ item.label }}</div>
                      <div class="text-xs text-neutral-500 mt-1">{{ item.desc }}</div>
                   </button>
                </div>
             </div>
          </div>
       </div>
    </div>
  </div>
  </div>
  </div>
</template>

<script setup lang="ts">
import draggable from 'vuedraggable'
import { CAMPAIGN_TEMPLATES, type CampaignTemplate, type Block as TemplateBlock } from '~/utils/templates'
import { THEME_PRESETS, type ThemePreset } from '~/utils/themes'
import CampaignHeroGenZ from '~/components/campaign/blocks/CampaignHeroGenZ.vue'
import CampaignFeaturesGenZ from '~/components/campaign/blocks/CampaignFeaturesGenZ.vue'
import CampaignCtaGenZ from '~/components/campaign/blocks/CampaignCtaGenZ.vue'
import CampaignHeroClean from '~/components/campaign/blocks/CampaignHeroClean.vue'
import CampaignSpeakerClean from '~/components/campaign/blocks/CampaignSpeakerClean.vue'
import CampaignFeaturesClean from '~/components/campaign/blocks/CampaignFeaturesClean.vue'
import CampaignCtaClean from '~/components/campaign/blocks/CampaignCtaClean.vue'
import CampaignSpeakerGenZ from '~/components/campaign/blocks/CampaignSpeakerGenZ.vue'
import CampaignHeroPro from '~/components/campaign/blocks/CampaignHeroPro.vue'
import CampaignSpeakerPro from '~/components/campaign/blocks/CampaignSpeakerPro.vue'
import CampaignFeaturesPro from '~/components/campaign/blocks/CampaignFeaturesPro.vue'
import CampaignCtaPro from '~/components/campaign/blocks/CampaignCtaPro.vue'
import CampaignHeroHealing from '~/components/campaign/blocks/CampaignHeroHealing.vue'
import CampaignProblemHealing from '~/components/campaign/blocks/CampaignProblemHealing.vue'
import CampaignSolutionHealing from '~/components/campaign/blocks/CampaignSolutionHealing.vue'
import CampaignSpeakerHealing from '~/components/campaign/blocks/CampaignSpeakerHealing.vue'
import CampaignCtaHealing from '~/components/campaign/blocks/CampaignCtaHealing.vue'

interface Course {
  id: string
  title: string
}

interface Webinar {
  id: string
  title: string
}

interface Block {
  id: string
  type: string
  enabled: boolean
  order: number
  data: any
}

definePageMeta({
  layout: 'admin',
  middleware: ['admin']
})

const route = useRoute()
const router = useRouter()
const api = useApi()

const isEdit = computed(() => !!route.params.id)
const campaignId = computed(() => route.params.id as string)

const saving = ref(false)
const showTemplateSelector = ref(false)
const showThemeSelector = ref(false)
const templates = CAMPAIGN_TEMPLATES
const themePresets = THEME_PRESETS

// Theme category filtering
const themeCategories = [
  { id: 'all', name: 'All', emoji: '✨' },
  { id: 'sales', name: 'Sales', emoji: '🔥' },
  { id: 'premium', name: 'Premium', emoji: '👑' },
  { id: 'creative', name: 'Creative', emoji: '🎨' },
  { id: 'dark', name: 'Dark', emoji: '🌑' },
  { id: 'professional', name: 'Pro', emoji: '💼' },
  { id: 'minimal', name: 'Minimal', emoji: '✨' },
]
const selectedThemeCategory = ref('all')
const filteredThemes = computed(() => {
  if (selectedThemeCategory.value === 'all') return themePresets
  return themePresets.filter(t => t.category === selectedThemeCategory.value)
})

const courses = ref<Course[]>([])
const webinars = ref<Webinar[]>([])
const editingBlock = ref<Block | null>(null)
const selectedDevice = ref('iphone_14_pro')
const activeTab = ref<'content' | 'appearance'>('content') // NEW: Tab State
const appearanceTab = ref<'themes' | 'background' | 'buttons' | 'fonts'>('themes') // Sub-tabs for Appearance


const devicePresets = {
  iphone_14_pro: { 
     name: 'iPhone 14 Pro', 
     width: '393px', 
     height: '852px', 
     radius: '50px', 
     notch: 'island' 
  },
  iphone_se: { 
     name: 'iPhone SE', 
     width: '375px', 
     height: '667px', 
     radius: '4px', // Standard phone border 
     notch: 'notch' 
  },
  pixel_7: { 
     name: 'Pixel 7', 
     width: '412px', 
     height: '915px', 
     radius: '24px', 
     notch: 'hole' 
  }
}

const form = ref({
  title: '',
  slug: '',
  course_id: null as string | null,
  webinar_id: null as string | null,
  gtm_id: '',
  facebook_pixel_id: '',
  meta_description: '',
  end_date: '',
  is_active: false,
  campaign_type: 'ecourse_only',
  styles: {
    primaryColor: '#6366f1',
    accentColor: '#f59e0b',
    backgroundColor: '#ffffff',
    buttonColor: '#6366f1',
    buttonTextColor: '#ffffff',
    textPrimaryColor: '#111827',
    textSecondaryColor: '#4b5563',
    buttonStyle: 'solid',
    buttonShape: 'rounded',
    cardStyle: 'flat',
    hasGradient: false,
    fontFamilyHeading: 'Inter',
    fontFamilyBody: 'Inter',
    
    // Background Fields
    backgroundType: 'solid', 
    backgroundImage: '',
    backgroundGradient: '',
    backgroundOverlay: '',
    
    layoutMode: 'wide'
  },
  globalCustomCSS: '', // NEW: Global CSS for entire landing page
  blocks: [
    {id: 'hero_1', type: 'hero', enabled: true, order: 1, data: {headline: 'Judul Kursus Anda', subheadline: 'Deskripsi singkat yang menarik', cta_text: '🚀 Daftar Sekarang', cta_link: '#pricing', background_image: null, badge: ''}},
    {id: 'countdown_1', type: 'countdown', enabled: true, order: 2, data: {label: '⏰ Promo Berakhir Dalam'}},
    {id: 'benefits_1', type: 'benefits', enabled: true, order: 3, data: {title: 'Yang Akan Anda Dapatkan', subtitle: 'Semua yang Anda butuhkan untuk sukses', items: [{icon: 'check', text: 'Akses materi selamanya'}, {icon: 'video', text: 'Video berkualitas tinggi'}, {icon: 'certificate', text: 'Sertifikat kelulusan'}]}},
    {id: 'pricing_1', type: 'pricing', enabled: true, order: 4, data: {original_price: 500000, discount_price: 199000, currency: 'IDR', cta_text: '🔥 Beli Sekarang', show_timer: true}},
    {id: 'testimonials_1', type: 'testimonials', enabled: false, order: 5, data: {title: 'Apa Kata Mereka?', items: []}},
    {id: 'video_1', type: 'video', enabled: false, order: 6, data: {title: 'Preview Materi', subtitle: 'Lihat apa yang akan Anda pelajari', youtube_url: ''}},
    {id: 'trust_1', type: 'trust', enabled: true, order: 7, data: {}},
    {id: 'faq_1', type: 'faq', enabled: false, order: 8, data: {title: 'Pertanyaan Umum', items: []}},
    {id: 'instructor_1', type: 'instructor', enabled: false, order: 9, data: {auto_fill: true}},
    {id: 'cta_1', type: 'cta_banner', enabled: false, order: 10, data: {headline: '🚀 Siap Untuk Memulai?', subheadline: 'Jangan lewatkan kesempatan untuk meningkatkan skill Anda', cta_text: 'Daftar Sekarang'}},
    {id: 'statistics_1', type: 'statistics', enabled: false, order: 11, data: {title: '', items: [{value: '5000', suffix: '+', label: 'Siswa'}, {value: '98', suffix: '%', label: 'Rating Puas'}, {value: '50', suffix: '+', label: 'Video'}, {value: '24/7', suffix: '', label: 'Support'}]}},
    {id: 'bonus_1', type: 'bonus', enabled: false, order: 12, data: {title: 'Bonus Eksklusif untuk Anda!', items: [{emoji: '📚', title: 'Ebook Panduan Lengkap', value: 150000, description: 'PDF 50+ halaman'}, {emoji: '📋', title: 'Template Excel', value: 100000, description: 'Siap pakai'}], total_value: 250000}},
    {id: 'curriculum_1', type: 'curriculum', enabled: false, order: 13, data: {title: '📚 Materi yang Akan Anda Pelajari', subtitle: 'Kurikulum lengkap dan terstruktur', modules: [{title: 'Modul 1: Pengenalan', lessons_count: 5, lessons: ['Apa itu...', 'Mengapa penting...', 'Persiapan awal']}, {title: 'Modul 2: Praktik Dasar', lessons_count: 8, lessons: ['Langkah pertama', 'Latihan 1', 'Latihan 2']}], total_lessons: 25, total_hours: 10}},
    {id: 'gallery_1', type: 'gallery', enabled: false, order: 14, data: {title: 'Galeri Kegiatan', items: [{image: 'https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&w=800&q=80', caption: 'Suasana Kelas'}, {image: 'https://images.unsplash.com/photo-1524178232363-1fb2b075b655?auto=format&fit=crop&w=800&q=80', caption: 'Mentoring Session'}]}},
    {id: 'achievement_1', type: 'achievement', enabled: false, order: 15, data: {items: [{emoji: '🏆', text: 'Tersertifikasi Nasional'}, {emoji: '⭐', text: 'Best Seller 2024'}, {emoji: '👥', text: '10.000+ Alumni'}]}},
    {id: 'comparison_1', type: 'comparison', enabled: false, order: 16, data: {title: 'Bandingkan Paket', headers: ['Fitur', 'Basic', 'PRO'], rows: [{feature: 'Akses Materi', val_a: true, val_b: true}, {feature: 'Sertifikat', val_a: false, val_b: true}, {feature: 'Konsultasi Mentor', val_a: false, val_b: true}]}},
    {id: 'floating_chat_1', type: 'floating_chat', enabled: false, order: 17, data: {phone: '', message: 'Halo, saya tertarik dengan kursus ini', label: 'Chat Kami'}},
    {id: 'social_proof_1', type: 'social_proof', enabled: false, order: 18, data: {text: 'orang sedang melihat penawaran ini', min_count: 15, max_count: 45, interval: 5}},
  ] as Block[]
})

const enabledBlocks = computed(() => form.value.blocks.filter(b => b.enabled))

const previewStyles = computed(() => ({
  fontFamily: (form.value.styles.fontFamilyBody || form.value.styles.fontFamily || 'Inter') + ', sans-serif',
  backgroundColor: form.value.styles.backgroundColor || '#ffffff',
  '--primary-color': form.value.styles.primaryColor,
  '--accent-color': form.value.styles.accentColor,
  '--bg-color': form.value.styles.backgroundColor,
  '--btn-color': form.value.styles.buttonColor,
  '--text-primary': form.value.styles.textPrimaryColor,
  '--text-secondary': form.value.styles.textSecondaryColor,
  '--font-heading': (form.value.styles.fontFamilyHeading || 'Inter') + ', sans-serif',
  '--font-body': (form.value.styles.fontFamilyBody || 'Inter') + ', sans-serif',
  '--radius': form.value.styles.borderRadius === 'pill' ? '999px' : form.value.styles.borderRadius === 'sharp' ? '0px' : '0.75rem',
  color: form.value.styles.textPrimaryColor || '#111827'
}))

// Computed button style for preview - mirrors live page logic
const previewButtonStyle = computed(() => {
  const s = form.value.styles
  const shape = s.buttonShape || 'rounded'
  const style = s.buttonStyle || 'solid'
  const bgColor = s.buttonColor || '#6366f1'
  const textColor = s.buttonTextColor || '#ffffff'
  const primary = s.primaryColor || '#6366f1'
  const accent = s.accentColor || '#f59e0b'

  const base = {
    fontFamily: (s.fontFamilyHeading || 'Inter') + ', sans-serif',
    borderRadius: shape === 'pill' ? '999px' : shape === 'sharp' ? '0px' : shape === 'leaf' ? '1.5rem 0.25rem 1.5rem 0.25rem' : '0.75rem',
    transition: 'all 0.2s ease-in-out',
    fontWeight: '600'
  }

  const computedTextColor = s.buttonTextColor || (style === 'outline' || style === 'glass' ? bgColor : '#ffffff')

  switch (style) {
    case 'outline':
      return { ...base, backgroundColor: 'transparent', border: `2px solid ${bgColor}`, color: bgColor }
    case 'gradient':
      return { ...base, backgroundImage: `linear-gradient(to right, ${primary}, ${accent})`, color: '#ffffff', border: 'none', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }
    case 'glass':
      return { ...base, backgroundColor: 'rgba(255, 255, 255, 0.15)', backdropFilter: 'blur(10px)', border: '1px solid rgba(255, 255, 255, 0.3)', color: s.buttonTextColor || '#ffffff' }
    case 'neo-brutalism':
      return { ...base, backgroundColor: bgColor, color: computedTextColor, border: '3px solid #000000', boxShadow: '5px 5px 0px #000000', fontWeight: '800', textTransform: 'uppercase' as const }
    case 'soft-shadow':
      return { ...base, backgroundColor: bgColor, color: computedTextColor, boxShadow: `0 10px 25px -5px ${bgColor}66`, border: 'none' }
    case 'solid':
    default:
      return { ...base, backgroundColor: bgColor, color: computedTextColor, border: 'none', boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1)' }
  }
})

// NEW: Identify current theme based on styles to highlight it
// NEW: Identify current theme based on colors (simplified match)
const currentThemeId = computed(() => {
  // Simple check based on primary color match to avoid full object comparison issues
  const match = themePresets.find(t => t.styles.primaryColor === form.value.styles.primaryColor && t.styles.backgroundColor === form.value.styles.backgroundColor)
  return match ? match.id : null
})

const isRigidTemplate = computed(() => {
  // Identify if we are using the "Rigid" Gen Z template based on its unique blocks
  return form.value.blocks.some(b => ['hero_gen_z', 'features_gen_z', 'cta_gen_z'].includes(b.type))
})

const applyTheme = (theme: any) => {
  form.value.styles = { ...theme.styles }
}

const getBlockEmoji = (type: string) => {
  const emojis: Record<string, string> = {
    hero: '🎯', countdown: '⏰', benefits: '✨', pricing: '💰',
    testimonials: '💬', faq: '❓', instructor: '👨‍🏫',
    video: '🎬', trust: '🛡️', cta_banner: '📢',
    statistics: '📊', bonus: '🎁', curriculum: '📚',
    gallery: '📷', achievement: '🏆', comparison: '🆚',
    floating_chat: '💬', social_proof: '👀'
  }
  return emojis[type] || '📄'
}

const getBlockTitle = (type: string) => {
  const titles: Record<string, string> = {
    hero: 'Hero Section', countdown: 'Countdown Timer', benefits: 'Benefits List',
    pricing: 'Pricing Section', testimonials: 'Testimonials', faq: 'FAQ Section', 
    instructor: 'Instructor', video: 'Video Section', trust: 'Trust Badges', cta_banner: 'CTA Banner',
    statistics: 'Statistics', bonus: 'Bonus Section', curriculum: 'Curriculum',
    gallery: 'Gallery', achievement: 'Achievement', comparison: 'Comparison Table',
    floating_chat: 'Floating Chat Button', social_proof: 'Social Proof / Live'
  }
  return titles[type] || type
}

const formatPrice = (price: number) => new Intl.NumberFormat('id-ID').format(price || 0)

const slugify = () => {
  form.value.slug = form.value.title
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

const editBlock = (block: Block) => {
  editingBlock.value = block
}

// Drag-to-resize state
const resizing = ref<{ block: Block; startY: number; startHeight: number } | null>(null)

const startResize = (block: Block, event: MouseEvent) => {
  event.preventDefault()
  resizing.value = {
    block,
    startY: event.clientY,
    startHeight: block.data.minHeight || 0
  }
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
}

const handleResize = (event: MouseEvent) => {
  if (!resizing.value) return
  const { block, startY, startHeight } = resizing.value
  const delta = event.clientY - startY
  const newHeight = Math.max(0, Math.min(600, startHeight + delta))
  block.data.minHeight = Math.round(newHeight / 10) * 10 // Snap to 10px increments
}

const stopResize = () => {
  resizing.value = null
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
}

// Helper to get preview styles for a block
const getBlockPreviewStyle = (block: Block) => {
  return {
    ...(block.data.minHeight && { minHeight: `${block.data.minHeight}px` }),
    ...(block.data.backgroundColor && { backgroundColor: block.data.backgroundColor })
  }
}


// Scope CSS to specific block by prepending block ID
const getScopedCSS = (block: Block): string => {
  const css = block.data.customCSS || ''
  const blockId = `#block-${block.id}`
  
  // Simple CSS scoping: prepend block ID to all selectors
  return css.split('}').map((rule: string) => {
    if (!rule.trim()) return ''
    const [selector, ...rest] = rule.split('{')
    if (!selector || !rest.length) return rule + '}'
    
    // Handle multiple selectors separated by comma
    const scopedSelector = selector.split(',').map((s: string) => {
      const trimmed = s.trim()
      if (!trimmed) return ''
      // Don't scope @keyframes or @media
      if (trimmed.startsWith('@')) return trimmed
      return `${blockId} ${trimmed}`
    }).join(', ')
    
    return `${scopedSelector} {${rest.join('{')}}` 
  }).join('\n')
}

// Code Snippet Templates
const codeSnippets: Record<string, { html: string; css: string }> = {
  hero_gradient: {
    html: `<div class="hero-gradient">
  <div class="hero-content">
    <h1>Your Headline Here</h1>
    <p>Your subheadline goes here. Make it compelling!</p>
    <button class="cta-btn">Get Started →</button>
  </div>
</div>`,
    css: `.hero-gradient {
  min-height: 500px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 60px 20px;
}
.hero-content h1 {
  font-size: 3rem;
  font-weight: 800;
  color: white;
  margin-bottom: 1rem;
}
.hero-content p {
  font-size: 1.25rem;
  color: rgba(255,255,255,0.8);
  margin-bottom: 2rem;
}
.cta-btn {
  padding: 16px 32px;
  background: white;
  color: #667eea;
  font-weight: 700;
  border-radius: 50px;
  border: none;
  cursor: pointer;
  transition: transform 0.3s;
}
.cta-btn:hover { transform: scale(1.05); }`
  },
  hero_video_bg: {
    html: `<div class="hero-video">
  <video autoplay muted loop playsinline class="bg-video">
    <source src="YOUR_VIDEO_URL.mp4" type="video/mp4">
  </video>
  <div class="overlay"></div>
  <div class="hero-text">
    <h1>Powerful Video Hero</h1>
    <p>Background video creates immersive experience</p>
  </div>
</div>`,
    css: `.hero-video {
  position: relative;
  min-height: 500px;
  overflow: hidden;
}
.bg-video {
  position: absolute;
  top: 50%; left: 50%;
  transform: translate(-50%, -50%);
  min-width: 100%; min-height: 100%;
  object-fit: cover;
}
.overlay {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.5);
}
.hero-text {
  position: relative;
  z-index: 10;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 500px;
  text-align: center;
  color: white;
  padding: 20px;
}
.hero-text h1 { font-size: 3rem; font-weight: 800; }
.hero-text p { font-size: 1.25rem; opacity: 0.8; }`
  },
  card_glass: {
    html: `<div class="glass-container">
  <div class="glass-card">
    <h3>✨ Glassmorphism</h3>
    <p>Modern frosted glass effect that looks premium and elegant.</p>
  </div>
</div>`,
    css: `.glass-container {
  min-height: 300px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.glass-card {
  background: rgba(255,255,255,0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255,255,255,0.2);
  border-radius: 20px;
  padding: 40px;
  max-width: 400px;
  text-align: center;
  color: white;
}
.glass-card h3 { font-size: 1.5rem; margin-bottom: 1rem; }
.glass-card p { opacity: 0.8; }`
  },
  card_neon: {
    html: `<div class="neon-container">
  <div class="neon-card">
    <h3>⚡ Neon Glow</h3>
    <p>Cyberpunk-style glowing effect for maximum impact.</p>
  </div>
</div>`,
    css: `.neon-container {
  min-height: 300px;
  background: #0a0a0a;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.neon-card {
  background: #111;
  border: 2px solid #00f2ff;
  border-radius: 15px;
  padding: 40px;
  max-width: 400px;
  text-align: center;
  color: #00f2ff;
  box-shadow: 0 0 20px rgba(0,242,255,0.3), inset 0 0 20px rgba(0,242,255,0.1);
  animation: neon-pulse 2s infinite;
}
@keyframes neon-pulse {
  0%, 100% { box-shadow: 0 0 20px rgba(0,242,255,0.3); }
  50% { box-shadow: 0 0 40px rgba(0,242,255,0.6); }
}
.neon-card h3 { font-size: 1.5rem; margin-bottom: 1rem; }
.neon-card p { color: #94a3b8; }`
  },
  anim_fade_in: {
    html: `<div class="fade-in-section">
  <h2>Fade In Content</h2>
  <p>This content fades in beautifully when it appears.</p>
</div>`,
    css: `.fade-in-section {
  padding: 60px 20px;
  text-align: center;
  animation: fadeIn 1s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
.fade-in-section h2 { font-size: 2rem; font-weight: 700; margin-bottom: 1rem; }
.fade-in-section p { color: #666; }`
  },
  anim_slide_up: {
    html: `<div class="slide-up-section">
  <h2>Slide Up Animation</h2>
  <p>Content slides up smoothly for a dynamic feel.</p>
</div>`,
    css: `.slide-up-section {
  padding: 60px 20px;
  text-align: center;
  animation: slideUp 0.8s ease-out;
}
@keyframes slideUp {
  from { opacity: 0; transform: translateY(40px); }
  to { opacity: 1; transform: translateY(0); }
}
.slide-up-section h2 { font-size: 2rem; font-weight: 700; margin-bottom: 1rem; }
.slide-up-section p { color: #666; }`
  }
}




const showAddBlockModal = ref(false)

const blockCategories = {
  basic: [
    { type: 'hero', label: 'Hero / Profile Header', icon: '🎯', desc: 'Judul utama atau profil' },
    { type: 'video', label: 'Video Embed', icon: '🎬', desc: 'YouTube/Vimeo' },
    { type: 'link_button', label: 'Link Button', icon: '🔗', desc: 'Tombol dengan Icon' },
    { type: 'cta_banner', label: 'CTA Banner', icon: '📢', desc: 'Banner Promosi' },
    { type: 'floating_chat', label: 'Floating Chat', icon: '💬', desc: 'WhatsApp Button' },
    { type: 'faq', label: 'FAQ Accordion', icon: '❓', desc: 'Tanya Jawab' },
    { type: 'gallery', label: 'Gallery', icon: '📷', desc: 'Slider Gambar' },
    { type: 'trust', label: 'Trust Badges', icon: '🛡️', desc: 'Logos & Keamanan' },
    { type: 'achievement', label: 'Achievement', icon: '🏆', desc: 'Pencapaian' },
    { type: 'social_proof', label: 'Live Viewer', icon: '👀', desc: 'Notifikasi Fake' }
  ],
  monetization: [
    { type: 'curriculum', label: 'Course Content', icon: '📚', desc: 'Daftar Materi' },
    { type: 'bonus', label: 'Product / Bonus', icon: '🎁', desc: 'Produk Digital' },
    { type: 'pricing', label: 'Pricing Card', icon: '💰', desc: 'Harga & Diskon' },
    { type: 'instructor', label: 'Instructor Profile', icon: '👨‍🏫', desc: 'Bio Pengajar' },
    { type: 'statistics', label: 'Statistics', icon: '📊', desc: 'Angka Statistik' },
    { type: 'testimonials', label: 'Testimonials', icon: '💬', desc: 'Review Siswa' },
    { type: 'countdown', label: 'Countdown', icon: '⏰', desc: 'Timer Promo' },
    { type: 'benefits', label: 'Benefits', icon: '✨', desc: 'Keuntungan' }
  ],
  special: [
    { type: 'hero_gen_z', label: 'Hero Gen Z', icon: '🚀', desc: 'Floating Emojis & Gradient (Lovable)' },
    { type: 'features_gen_z', label: 'Features Gen Z', icon: '✨', desc: 'Glassmorphism Cards (Lovable)' },
    { type: 'cta_gen_z', label: 'CTA Gen Z', icon: '🔥', desc: 'Neon Form Registration (Lovable)' }
  ]
}

const addBlock = (type: string) => {
  const existingBlock = form.value.blocks.find(b => b.type === type && !b.enabled)
  
  // Helper to get default styled data
  const getStyleDefaults = () => ({
     backgroundColor: form.value.styles.backgroundColor || '#ffffff',
     titleColor: form.value.styles.textPrimaryColor || '#000000',
     textColor: form.value.styles.textSecondaryColor || '#4b5563',
     buttonColor: form.value.styles.buttonColor || '#2563eb'
  })

  // Specific defaults per block type
  const getBlockDefaults = (type: string) => {
      const styles = form.value.styles
      const baseStyles = getStyleDefaults()
      
      switch(type) {
          case 'faq':
             return {
                title: 'Pertanyaan Umum',
                items: [{ question: 'Apakah ini cocok untuk pemula?', answer: 'Tentu saja! Materi disusun dari dasar.' }],
                ...baseStyles
             }
          case 'pricing':
             return {
                original_price: 1000000,
                discount_price: 499000,
                cta_text: 'Beli Sekarang',
                trust_text: 'Garansi 30 Hari',
                // Special case for Gen Z/Dark themes: ensure contrast if bg is dark
                cardBgColor: styles.backgroundColor === '#050505' ? '#171717' : '#ffffff',
                ...baseStyles
             }
          case 'testimonials':
             return {
                title: 'Apa Kata Mereka',
                items: [{ name: 'Budi Santoso', text: 'Materinya sangat daging!' }],
                ...baseStyles
             }
           case 'hero_gen_z':
              return {
                 headline: 'NEW HERO',
                 subheadline: 'Subtitle here...',
                 badge: 'NEW',
                 promo_text: 'Limited Time',
                 ...baseStyles
              }
           case 'features_gen_z':
              return {
                 title: 'Key Features',
                 items: [{ title: 'Feature 1', desc: 'Description', icon: '✨', color: 'text-white', bg: 'bg-white/10' }],
                 ...baseStyles
              }
           case 'cta_gen_z':
              return {
                  headline: 'Ready to Join?',
                  subheadline: 'Secure your spot',
                  button_text: 'Join Now',
                  price: 0,
                  original_price: 0,
                  ...baseStyles
              }
           case 'hero_clean':
               return { headline: 'Headline', subheadline: 'Subheadline', badge: 'WEBINAR', image_url: 'https://images.unsplash.com/photo-1573164713988-8665fc963095?auto=format&fit=crop&q=80&w=800', button_color: 'linear-gradient(to right, #1e3a8a, #3b82f6)', font_family: 'Inter', ...baseStyles }
           case 'speaker_gen_z':
               return { name: 'Joni', title: 'Expert', bio: 'Bio...', image_url: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=300', stats: [], ...baseStyles }
           case 'speaker_clean':
               return { name: 'Name', title: 'Title', bio: 'Bio...', stats: [{ icon: 'A', label: 'Stat 1' }, { icon: 'B', label: 'Stat 2' }], image_url: 'https://images.unsplash.com/photo-1573496359142-b8d87734a5a2?auto=format&fit=crop&q=80&w=800', font_family: 'Inter', ...baseStyles }
           case 'features_clean':
               return { title: 'Section Title', items: [{ title: 'Feature 1', text: 'Description', icon: '✨' }], font_family: 'Inter', ...baseStyles }
           case 'cta_clean':
               return { headline: 'Headline', subheadline: 'Subheadline', button_color: 'linear-gradient(to right, #1e3a8a, #3b82f6)', font_family: 'Inter', ...baseStyles }
           case 'link_button':
              return { label: 'Klik disini', url: '#', icon: '', style: 'solid', thumbnail: '' }
           case 'bonus':
              return { title: 'Produk Digital', items: [], total_value: 0, displayMode: 'list' }
           default:
              return { ...baseStyles }
      }
  }

  if (existingBlock) {
    existingBlock.enabled = true
    const maxOrder = Math.max(...form.value.blocks.map(b => b.order || 0))
    existingBlock.order = maxOrder + 1
    // Update data with current styles to ensure it matches theme if re-enabled
    existingBlock.data = { ...existingBlock.data, ...getStyleDefaults() }
  } else {
    const baseBlock = form.value.blocks.find(b => b.type === type)
    if (baseBlock) {
       // Clone existing block structure but refresh styles
       const newBlock = JSON.parse(JSON.stringify(baseBlock))
       newBlock.id = `${type}_${Date.now()}`
       newBlock.enabled = true
       newBlock.order = Math.max(...form.value.blocks.map(b => b.order || 0)) + 1
       newBlock.data = { ...newBlock.data, ...getStyleDefaults() }
       form.value.blocks.push(newBlock)
    } else {
       // Initialize completely new block
       const newBlock = {
          id: `${type}_${Date.now()}`,
          type: type,
          enabled: true,
          order: Math.max(0, ...form.value.blocks.map((b: any) => b.order || 0)) + 1,
          data: getBlockDefaults(type)
       }
       form.value.blocks.push(newBlock)
    }
  }
  showAddBlockModal.value = false
}

const removeBlock = (blockId: string) => {
   const idx = form.value.blocks.findIndex(b => b.id === blockId)
   if (idx !== -1) {
      form.value.blocks[idx].enabled = false
   }
}

const fetchCourses = async () => {
  try {
    const data = await api.fetch<{ courses: Course[] }>('/api/admin/courses')
    courses.value = data.courses || []
  } catch (err) {
    console.error('Failed to fetch courses:', err)
  }
}

const fetchWebinars = async () => {
  try {
    const data = await api.fetch<{ webinars: Webinar[] }>('/api/admin/webinars')
    webinars.value = data.webinars || []
  } catch (err) {
    console.error('Failed to fetch webinars:', err)
  }
}

const fetchCampaign = async () => {
  if (!isEdit.value) return
  
  try {
    const data = await api.fetch<any>(`/api/admin/campaigns/${campaignId.value}`)
    form.value.title = data.title
    form.value.slug = data.slug
    form.value.course_id = data.course_id || null
    form.value.webinar_id = data.webinar_id || null
    form.value.meta_description = data.meta_description || ''
    form.value.is_active = data.is_active
    form.value.campaign_type = data.campaign_type || 'ecourse_only'
    
    if (data.end_date) {
      const date = new Date(data.end_date)
      form.value.end_date = date.toISOString().slice(0, 16)
    }
    
    // Load blocks from server
    if (data.blocks) {
      // Handle if blocks is a string (needs parsing) or already an array
      const blocksData = typeof data.blocks === 'string' ? JSON.parse(data.blocks) : data.blocks
      if (Array.isArray(blocksData) && blocksData.length > 0) {
        form.value.blocks = blocksData
      }
    }
    
    // Load styles from server
    if (data.styles) {
      const stylesData = typeof data.styles === 'string' ? JSON.parse(data.styles) : data.styles
      form.value.styles = { ...form.value.styles, ...stylesData }
    }
  } catch (err) {
    console.error('Failed to fetch campaign:', err)
    alert('Gagal memuat campaign')
    router.push('/admin/campaigns')
  }
}

const previewPage = () => {
  if (form.value.slug) {
    window.open(`/c/${form.value.slug}`, '_blank')
  } else {
    alert('Silakan isi slug terlebih dahulu')
  }
}

const saveDraft = async () => await saveCampaign(false)
const publishCampaign = async () => await saveCampaign(true)

const saveCampaign = async (publish: boolean) => {
  if (!form.value.title || !form.value.slug) {
    alert('Judul dan slug wajib diisi')
    return
  }

  saving.value = true
  try {
    // Update block orders based on current array position
    const blocksWithOrder = form.value.blocks.map((block, index) => ({
      ...block,
      order: Number(index + 1)
    }))

    const payload = {
      title: form.value.title,
      slug: form.value.slug,
      course_id: form.value.course_id,
      webinar_id: form.value.webinar_id,
      meta_description: form.value.meta_description || null,
      end_date: form.value.end_date ? new Date(form.value.end_date).toISOString() : null,
      is_active: publish,
      campaign_type: form.value.campaign_type,
      blocks: blocksWithOrder,
      styles: form.value.styles
    }

    if (isEdit.value) {
      await api.fetch(`/api/admin/campaigns/${campaignId.value}`, {
        method: 'PUT',
        body: JSON.stringify(payload)
      })
    } else {
      await api.fetch('/api/admin/campaigns', {
        method: 'POST',
        body: JSON.stringify(payload)
      })
    }

    router.push('/admin/campaigns')
  } catch (err: any) {
    console.error('Failed to save campaign:', err)
    alert(err.data?.error || 'Gagal menyimpan campaign')
  } finally {
    saving.value = false
  }
}

const applyTemplate = (tpl: CampaignTemplate) => {
  if (!confirm('Apply template? This will replace current blocks.')) return

  // Apply blocks (deep copy to avoid reference issues)
  form.value.blocks = JSON.parse(JSON.stringify(tpl.blocks))
  
  // Deep merge to ensure all properties are updated
  form.value.styles = {
    ...form.value.styles,
    ...tpl.styles, // Assuming tpl.styles is the 'theme' being applied
    // Ensure backwards compatibility or specific mapping if needed
    borderRadius: tpl.styles?.buttonShape // Map new token to old property if used elsewhere, or just use buttonShape
  }
  
  // Auto-set campaign type if needed based on category
  if (tpl.category === 'webinar') {
     form.value.campaign_type = 'webinar_only'
  } else if (tpl.category === 'ecourse') {
     form.value.campaign_type = 'ecourse_only'
  }
  
  showTemplateSelector.value = false
}

onMounted(() => {
  fetchCourses()
  fetchWebinars()
  if (isEdit.value) {
    fetchCampaign()
  } else {
    // Show template selector for new campaigns
    showTemplateSelector.value = true
  }
})
</script>
