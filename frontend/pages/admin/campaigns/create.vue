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
      <!-- Left Sidebar -->
      <div class="w-80 bg-white border-r border-neutral-200 overflow-y-auto">
        <div class="p-5 space-y-6">
          <!-- Basic Info -->
          <div class="space-y-4">
            <h2 class="font-semibold text-neutral-900">📋 Pengaturan Dasar</h2>
            
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Judul Campaign *</label>
              <input 
                v-model="form.title"
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
              <label class="block text-sm font-medium text-neutral-700 mb-1">Link ke Kursus</label>
              <select v-model="form.course_id" class="w-full px-3 py-2 border border-neutral-300 rounded-lg">
                <option :value="null">-- Pilih Kursus --</option>
                <option v-for="course in courses" :key="course.id" :value="course.id">
                  {{ course.title }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Countdown End Date</label>
              <input v-model="form.end_date" type="datetime-local" class="w-full px-3 py-2 border border-neutral-300 rounded-lg"/>
            </div>
          </div>

          <hr class="border-neutral-200">

          <!-- Style Customization -->
          <div class="space-y-4">
            <h2 class="font-semibold text-neutral-900">🎨 Kustomisasi Warna</h2>
            
            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Primary</label>
                <input v-model="form.styles.primaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Accent</label>
                <input v-model="form.styles.accentColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Background</label>
                <input v-model="form.styles.backgroundColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Button</label>
                <input v-model="form.styles.buttonColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Text Primary</label>
                <input v-model="form.styles.textPrimaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Text Secondary</label>
                <input v-model="form.styles.textSecondaryColor" type="color" class="w-full h-10 rounded border cursor-pointer"/>
              </div>
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Button Style</label>
                <select v-model="form.styles.buttonStyle" class="w-full px-3 py-2 border border-neutral-300 rounded-lg text-sm">
                  <option value="solid">Solid</option>
                  <option value="outline">Outline</option>
                  <option value="gradient">Gradient</option>
                </select>
              </div>
              <div>
                <label class="block text-xs font-medium text-neutral-600 mb-1">Corner Radius</label>
                <select v-model="form.styles.borderRadius" class="w-full px-3 py-2 border border-neutral-300 rounded-lg text-sm">
                  <option value="sharp">Sharp (0px)</option>
                  <option value="rounded">Rounded (8px)</option>
                  <option value="pill">Pill (999px)</option>
                </select>
              </div>
            </div>

            <div class="flex items-center gap-2">
              <input type="checkbox" v-model="form.styles.hasGradient" id="useGradient" class="rounded border-neutral-300 text-primary-600 focus:ring-primary-500"/>
              <label for="useGradient" class="text-sm text-neutral-700">Gunakan Gradient Background</label>
            </div>

            <div>
              <label class="block text-xs font-medium text-neutral-600 mb-1">Font</label>
              <select v-model="form.styles.fontFamily" class="w-full px-3 py-2 border border-neutral-300 rounded-lg text-sm">
                <option value="Inter">Inter (Modern)</option>
                <option value="Poppins">Poppins (Friendly)</option>
                <option value="Playfair Display">Playfair Display (Elegant)</option>
                <option value="Roboto">Roboto (Clean)</option>
              </select>
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
                    <button @click="editBlock(block)" class="p-1 text-neutral-500 hover:text-primary-600">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                    </button>
                  </div>
                </template>
              </draggable>
            </ClientOnly>
          </div>
        </div>
      </div>

      <!-- Right - Preview Area -->
      <div class="flex-1 overflow-y-auto bg-neutral-100 p-6">
        <div class="max-w-4xl mx-auto">
          <div class="bg-white rounded-xl shadow-lg overflow-hidden" :style="previewStyles">
            <!-- Live Preview -->
            <template v-for="block in enabledBlocks" :key="block.id">
              <!-- Hero Preview -->
              <div v-if="block.type === 'hero'" class="relative min-h-[300px] bg-gradient-to-br from-neutral-800 to-neutral-900 text-white p-8 flex flex-col justify-center items-center text-center">
                <h1 class="text-3xl md:text-4xl font-bold mb-4">{{ block.data.headline }}</h1>
                <p class="text-lg text-neutral-300 mb-6 max-w-xl">{{ block.data.subheadline }}</p>
                <button class="px-6 py-3 rounded-lg font-semibold" :style="{backgroundColor: form.styles.buttonColor}">
                  {{ block.data.cta_text }}
                </button>
              </div>

              <!-- Countdown Preview -->
              <div v-else-if="block.type === 'countdown'" class="bg-neutral-900 text-white py-8 text-center">
                <p class="text-neutral-400 mb-4">{{ block.data.label }}</p>
                <div class="flex justify-center gap-4">
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: form.styles.primaryColor}">00</div>
                    <div class="text-xs text-neutral-400">Hari</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: form.styles.primaryColor}">00</div>
                    <div class="text-xs text-neutral-400">Jam</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: form.styles.primaryColor}">00</div>
                    <div class="text-xs text-neutral-400">Menit</div>
                  </div>
                  <div class="bg-neutral-800 px-4 py-3 rounded-lg">
                    <div class="text-2xl font-bold" :style="{color: form.styles.primaryColor}">00</div>
                    <div class="text-xs text-neutral-400">Detik</div>
                  </div>
                </div>
              </div>

              <!-- Benefits Preview -->
              <div v-else-if="block.type === 'benefits'" class="py-10 px-8 bg-white">
                <h2 class="text-2xl font-bold text-center mb-6">{{ block.data.title }}</h2>
                <div class="grid md:grid-cols-2 gap-4">
                  <div v-for="(item, idx) in block.data.items" :key="idx" class="flex items-center gap-3 p-3 bg-neutral-50 rounded-lg">
                    <div class="w-8 h-8 rounded-full flex items-center justify-center" :style="{backgroundColor: form.styles.primaryColor + '20', color: form.styles.primaryColor}">
                      ✓
                    </div>
                    <span class="text-neutral-700">{{ item.text }}</span>
                  </div>
                </div>
              </div>

              <!-- Pricing Preview -->
              <div v-else-if="block.type === 'pricing'" class="py-10 px-8 bg-gradient-to-b from-neutral-900 to-neutral-800 text-white text-center">
                <h2 class="text-2xl font-bold mb-6">Penawaran Spesial</h2>
                <div class="text-neutral-400 line-through text-lg">Rp {{ formatPrice(block.data.original_price) }}</div>
                <div class="text-4xl font-bold my-2" :style="{color: form.styles.primaryColor}">
                  Rp {{ formatPrice(block.data.discount_price) }}
                </div>
                <button class="mt-4 px-8 py-3 rounded-lg font-bold" :style="{backgroundColor: form.styles.buttonColor}">
                  {{ block.data.cta_text }}
                </button>
              </div>

              <!-- Testimonials Preview -->
              <div v-else-if="block.type === 'testimonials'" class="py-10 px-8 bg-neutral-100">
                <h2 class="text-2xl font-bold text-center mb-6">{{ block.data.title }}</h2>
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
              <div v-else-if="block.type === 'faq'" class="py-10 px-8 bg-white">
                <h2 class="text-2xl font-bold text-center mb-6">{{ block.data.title }}</h2>
                <div class="space-y-3 max-w-2xl mx-auto">
                  <div v-for="(item, idx) in block.data.items" :key="idx" class="border border-neutral-200 rounded-lg p-4">
                    <div class="font-medium text-neutral-900">{{ item.question }}</div>
                    <div class="text-neutral-600 text-sm mt-2">{{ item.answer }}</div>
                  </div>
                  <div v-if="!block.data.items?.length" class="text-center text-neutral-400 py-8">
                    Belum ada FAQ
                  </div>
                </div>
              </div>

              <!-- Instructor Preview -->
              <div v-else-if="block.type === 'instructor'" class="py-10 px-8 bg-neutral-50 text-center">
                <h2 class="text-2xl font-bold mb-6">👨‍🏫 Tentang Instruktur</h2>
                <div class="w-20 h-20 bg-neutral-300 rounded-full mx-auto mb-4 flex items-center justify-center text-3xl">
                  👨‍🏫
                </div>
                <div class="font-bold text-lg">{{ block.data.name || 'Nama Instruktur' }}</div>
                <p class="text-neutral-600 max-w-md mx-auto mt-2">{{ block.data.bio || 'Bio instruktur' }}</p>
              </div>

              <!-- Video Preview (NEW) -->
              <div v-else-if="block.type === 'video'" class="py-10 px-8 bg-neutral-900 text-center">
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
              <div v-else-if="block.type === 'trust'" class="py-8 px-6 bg-white border-y border-neutral-100">
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
              <div v-else-if="block.type === 'cta_banner'" class="py-10 px-8 text-white text-center" :style="{backgroundColor: form.styles.buttonColor}">
                <h2 class="text-2xl font-bold mb-3">{{ block.data.headline || '🚀 Siap Untuk Memulai?' }}</h2>
                <p class="opacity-80 mb-6">{{ block.data.subheadline || 'Jangan lewatkan kesempatan untuk meningkatkan skill Anda' }}</p>
                <button class="px-6 py-3 bg-white rounded-lg font-semibold" :style="{color: form.styles.buttonColor}">
                  {{ block.data.cta_text || 'Daftar Sekarang' }}
                </button>
              </div>

              <!-- Statistics Preview (NEW) -->
              <div v-else-if="block.type === 'statistics'" class="py-10 px-8 bg-white">
                <h2 v-if="block.data.title" class="text-xl font-bold text-center mb-6">{{ block.data.title }}</h2>
                <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
                  <div v-for="(stat, idx) in block.data.items" :key="idx" class="text-center p-4 bg-neutral-50 rounded-xl">
                    <div class="text-2xl font-black mb-1" :style="{color: form.styles.primaryColor}">
                      {{ stat.value }}{{ stat.suffix || '' }}
                    </div>
                    <div class="text-neutral-600 text-xs font-medium">{{ stat.label }}</div>
                  </div>
                </div>
              </div>

              <!-- Bonus Preview (NEW) -->
              <div v-else-if="block.type === 'bonus'" class="py-10 px-8 bg-amber-50">
                <div class="text-center mb-6">
                  <span class="inline-block bg-yellow-400 text-yellow-900 text-[10px] font-bold px-2 py-0.5 rounded-full mb-2">🎁 BONUS</span>
                  <h2 class="text-xl font-bold text-neutral-900">{{ block.data.title || 'Bonus Eksklusif' }}</h2>
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
              <div v-else-if="block.type === 'curriculum'" class="py-10 px-8 bg-white">
                <div class="text-center mb-8">
                  <h2 class="text-xl font-bold text-neutral-900">{{ block.data.title || 'Materi Kursus' }}</h2>
                  <p v-if="block.data.subtitle" class="text-neutral-500 text-sm">{{ block.data.subtitle }}</p>
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
              <div v-else-if="block.type === 'gallery'" class="py-10 bg-white text-center">
                 <h2 v-if="block.data.title" class="font-bold mb-4">{{ block.data.title }}</h2>
                 <div class="flex gap-2 overflow-x-auto px-4 pb-2 snap-x">
                    <div v-for="(img, idx) in block.data.items" :key="idx" class="w-48 h-32 flex-shrink-0 bg-neutral-200 rounded-lg overflow-hidden relative snap-center">
                       <img v-if="img.image" :src="img.image" class="w-full h-full object-cover"/>
                       <div v-else class="w-full h-full flex items-center justify-center text-4xl">🖼️</div>
                    </div>
                 </div>
              </div>

              <!-- Achievement Preview -->
              <div v-else-if="block.type === 'achievement'" class="py-6 bg-neutral-50 border-y border-neutral-100">
                 <div class="flex flex-wrap justify-center gap-4 px-4">
                    <div v-for="(badge, idx) in block.data.items" :key="idx" class="flex items-center gap-2 bg-white px-3 py-2 rounded-full shadow-sm border border-neutral-100">
                       <span>{{ badge.emoji }}</span>
                       <span class="text-sm font-medium">{{ badge.text }}</span>
                    </div>
                 </div>
              </div>
              <!-- Comparison Preview -->
              <div v-else-if="block.type === 'comparison'" class="py-10 px-4 bg-white">
                 <h2 v-if="block.data.title" class="text-center font-bold mb-4">{{ block.data.title }}</h2>
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
              <div v-else-if="block.type === 'floating_chat'" class="py-10 bg-neutral-100 flex items-center justify-center relative h-40">
                 <div class="absolute bottom-4 right-4 bg-green-500 text-white px-4 py-3 rounded-full shadow-lg flex items-center gap-2 cursor-pointer hover:bg-green-600 transition-colors">
                    <span class="text-xl">💬</span>
                    <span v-if="block.data.label" class="font-bold text-sm">{{ block.data.label }}</span>
                 </div>
                 <p class="text-neutral-400 text-sm">Tombol akan melayang di pojok kanan bawah halaman.</p>
              </div>

              <!-- Social Proof Preview -->
              <div v-else-if="block.type === 'social_proof'" class="py-10 bg-neutral-100 flex items-center justify-center relative h-32">
                 <div class="absolute bottom-4 left-4 bg-white/90 backdrop-blur border border-neutral-200 px-4 py-2 rounded-lg shadow-lg flex items-center gap-3 animate-pulse">
                    <div class="w-2 h-2 rounded-full bg-red-500 animate-ping"></div>
                    <span class="text-sm font-medium text-neutral-800">
                      <span class="font-bold">24</span> {{ block.data.text || 'orang sedang melihat ini' }}
                    </span>
                 </div>
                 <p class="text-neutral-400 text-sm">Notifikasi Live Viewer di pojok kiri bawah.</p>
              </div>
            </template>

            <div v-if="!enabledBlocks.length" class="py-20 text-center text-neutral-400">
              <p class="text-lg">Aktifkan block untuk melihat preview</p>
            </div>
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
          <!-- Hero Editor -->
          <template v-if="editingBlock.type === 'hero'">
            <div>
              <label class="block text-sm font-medium mb-1">Headline</label>
              <input v-model="editingBlock.data.headline" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Subheadline</label>
              <textarea v-model="editingBlock.data.subheadline" rows="2" class="w-full px-3 py-2 border rounded-lg"></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Background Image URL</label>
              <input v-model="editingBlock.data.background_image" type="url" class="w-full px-3 py-2 border rounded-lg" placeholder="https://..."/>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium mb-1">CTA Text</label>
                <input v-model="editingBlock.data.cta_text" type="text" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
              <div>
                <label class="block text-sm font-medium mb-1">CTA Link</label>
                <input v-model="editingBlock.data.cta_link" type="text" class="w-full px-3 py-2 border rounded-lg"/>
              </div>
            </div>
          </template>

          <!-- Countdown Editor -->
          <template v-else-if="editingBlock.type === 'countdown'">
            <div>
              <label class="block text-sm font-medium mb-1">Label</label>
              <input v-model="editingBlock.data.label" type="text" class="w-full px-3 py-2 border rounded-lg"/>
            </div>
            <p class="text-sm text-neutral-500">* Waktu countdown diambil dari "Countdown End Date"</p>
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
            <div class="flex items-center gap-2">
              <input type="checkbox" v-model="editingBlock.data.show_timer" class="rounded border-neutral-300"/>
              <label class="text-sm">Tampilkan timer</label>
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
                <input v-model="editingBlock.data.avatar" type="url" class="w-full px-3 py-2 border rounded-lg"/>
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

          <!-- Trust Badges - No editor needed, just toggle -->
          <template v-else-if="editingBlock.type === 'trust'">
            <div class="text-center text-neutral-500 py-4">
              <p class="text-4xl mb-2">🛡️</p>
              <p>Block ini menampilkan trust badges otomatis:</p>
              <div class="flex flex-wrap gap-3 justify-center mt-4 text-sm">
                <span>🔒 Pembayaran Aman</span>
                <span>💯 Garansi 30 Hari</span>
                <span>🎓 Sertifikat</span>
                <span>♾️ Akses Selamanya</span>
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
        </div>

        <div class="flex justify-end gap-3 p-4 border-t">
          <button @click="editingBlock = null" class="px-4 py-2 border rounded-lg hover:bg-neutral-50">Tutup</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import draggable from 'vuedraggable'

interface Course {
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
const courses = ref<Course[]>([])
const editingBlock = ref<Block | null>(null)

const form = ref({
  title: '',
  slug: '',
  course_id: null as string | null,
  meta_description: '',
  end_date: '',
  is_active: false,
  styles: {
    primaryColor: '#6366f1',
    accentColor: '#f59e0b',
    backgroundColor: '#1f2937',
    buttonColor: '#6366f1',
    textPrimaryColor: '#111827',
    textSecondaryColor: '#4b5563',
    buttonStyle: 'solid',
    borderRadius: 'rounded',
    hasGradient: false,
    fontFamily: 'Inter'
  },
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
  fontFamily: form.value.styles.fontFamily + ', sans-serif',
  '--primary-color': form.value.styles.primaryColor,
  '--accent-color': form.value.styles.accentColor,
  '--bg-color': form.value.styles.backgroundColor,
  '--btn-color': form.value.styles.buttonColor,
  '--text-primary': form.value.styles.textPrimaryColor,
  '--text-secondary': form.value.styles.textSecondaryColor,
  '--radius': form.value.styles.borderRadius === 'pill' ? '999px' : form.value.styles.borderRadius === 'sharp' ? '0px' : '0.75rem',
}))

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

const fetchCourses = async () => {
  try {
    const response = await api.fetch<{ courses: Course[] }>('/api/courses?limit=100')
    courses.value = response?.courses || []
  } catch (err) {
    console.error('Failed to fetch courses:', err)
  }
}

const fetchCampaign = async () => {
  if (!isEdit.value) return
  
  try {
    const data = await api.fetch<any>(`/api/admin/campaigns/${campaignId.value}`)
    form.value.title = data.title
    form.value.slug = data.slug
    form.value.course_id = data.course_id || null
    form.value.meta_description = data.meta_description || ''
    form.value.is_active = data.is_active
    
    if (data.end_date) {
      const date = new Date(data.end_date)
      form.value.end_date = date.toISOString().slice(0, 16)
    }
    
    if (data.blocks && Array.isArray(data.blocks)) {
      form.value.blocks = data.blocks
    }
    
    if (data.styles) {
      form.value.styles = { ...form.value.styles, ...data.styles }
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
      meta_description: form.value.meta_description || null,
      end_date: form.value.end_date ? new Date(form.value.end_date).toISOString() : null,
      is_active: publish,
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

onMounted(() => {
  fetchCourses()
  if (isEdit.value) {
    fetchCampaign()
  }
})
</script>
