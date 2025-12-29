<template>
  <div>
    <!-- Header -->
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-neutral-900">Pengaturan Sistem</h1>
      <p class="text-neutral-500 mt-1">Konfigurasi platform dan preferensi</p>
    </div>

    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Sidebar Navigation -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
          <nav class="divide-y divide-neutral-100">
            <button 
              v-for="tab in tabs" 
              :key="tab.id"
              @click="activeTab = tab.id"
              class="w-full flex items-center gap-3 p-4 text-left hover:bg-neutral-50 transition-colors"
              :class="activeTab === tab.id ? 'bg-admin-50 border-r-2 border-admin-600' : ''"
            >
              <div 
                class="w-10 h-10 rounded-lg flex items-center justify-center"
                :class="activeTab === tab.id ? 'bg-admin-100 text-admin-600' : 'bg-neutral-100 text-neutral-500'"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="tab.icon"/>
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium" :class="activeTab === tab.id ? 'text-admin-700' : 'text-neutral-700'">{{ tab.name }}</p>
                <p class="text-xs text-neutral-500">{{ tab.description }}</p>
              </div>
            </button>
          </nav>
        </div>
      </div>

      <!-- Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- General Settings -->
        <div v-if="activeTab === 'general'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Pengaturan Umum</h3>
          
          <form class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Platform</label>
              <input 
                v-model="form.site_name"
                type="text" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Deskripsi</label>
              <input 
                v-model="form.site_description"
                type="text" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Email Kontak</label>
              <input 
                v-model="form.contact_email"
                type="email" 
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
            <div class="pt-4">
              <button type="button" @click="saveGeneralSettings" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}</button>
            </div>
          </form>
        </div>

        <!-- Theme Settings -->
        <div v-if="activeTab === 'theme'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Pengaturan Tema</h3>
          
          <!-- Primary Color -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-neutral-700 mb-3">Warna Primer</label>
            <div class="flex gap-3">
              <button 
                v-for="color in primaryColors" 
                :key="color.value"
                @click="form.theme = color.value"
                :class="['w-10 h-10 rounded-lg transition-all', color.class, form.theme === color.value ? 'ring-2 ring-offset-2 ring-neutral-900' : '']"
              ></button>
            </div>
          </div>
          
          <!-- Logo -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-neutral-700 mb-3">Logo</label>
            <div class="flex items-center gap-4">
              <div class="w-20 h-20 bg-neutral-100 rounded-xl flex items-center justify-center overflow-hidden">
                <img 
                  v-if="form.logo_url" 
                  :src="getLogoUrl(form.logo_url)" 
                  alt="Logo" 
                  class="w-full h-full object-contain"
                />
                <span v-else class="text-2xl font-bold text-neutral-400">L</span>
              </div>
              <div>
                <input 
                  ref="logoInput"
                  type="file" 
                  accept="image/png,image/jpeg,image/svg+xml,image/webp"
                  class="hidden"
                  @change="uploadLogo"
                />
                <button 
                  type="button"
                  @click="($refs.logoInput as HTMLInputElement)?.click()"
                  :disabled="uploadingLogo"
                  class="text-sm text-admin-600 font-medium hover:text-admin-700 disabled:opacity-50"
                >
                  {{ uploadingLogo ? 'Mengupload...' : 'Ubah Logo' }}
                </button>
                <p class="text-xs text-neutral-500 mt-1">PNG, JPG, WebP. Maksimal 10MB</p>
              </div>
            </div>
          </div>

          <button type="button" @click="saveThemeSettings" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Perubahan' }}</button>
        </div>

        <!-- Feature Flags -->
        <div v-if="activeTab === 'features'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Feature Toggles</h3>
          <p class="text-sm text-neutral-500 mb-6">Aktifkan atau nonaktifkan fitur platform</p>
          
          <div class="space-y-4">
            <div v-for="feature in features" :key="feature.id" class="flex items-center justify-between py-3 border-b border-neutral-100 last:border-0">
              <div>
                <p class="text-sm font-medium text-neutral-900">{{ feature.name }}</p>
                <p class="text-xs text-neutral-500">{{ feature.description }}</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="feature.enabled" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-admin-600"></div>
              </label>
            </div>
          </div>
        </div>

        <!-- Security Settings -->
        <div v-if="activeTab === 'security'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Keamanan</h3>
          
          <div class="space-y-6">
            <div class="flex items-center justify-between py-3">
              <div>
                <p class="text-sm font-medium text-neutral-900">Two-Factor Authentication</p>
                <p class="text-xs text-neutral-500">Wajibkan 2FA untuk semua admin</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="settings.require2FA" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-admin-600"></div>
              </label>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Session Timeout (menit)</label>
              <input 
                v-model="form.sessionTimeout"
                type="number" 
                class="w-32 px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
          </div>
        </div>
        
        <!-- Announcement Settings -->
        <div v-if="activeTab === 'announcement'" class="bg-white rounded-xl border border-neutral-200 p-6">
          <h3 class="font-semibold text-neutral-900 mb-6">Banner Pengumuman</h3>
          <p class="text-sm text-neutral-500 mb-6">Tampilkan informasi penting atau promo di bagian atas halaman utama.</p>
          
          <div class="space-y-6">
            <div class="flex items-center justify-between py-3 border-b border-neutral-100">
              <div>
                <p class="text-sm font-medium text-neutral-900">Aktifkan Banner</p>
                <p class="text-xs text-neutral-500">Banner akan muncul di halaman landing page</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="bannerSettings.enabled" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-admin-600"></div>
              </label>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Teks Pengumuman</label>
              <input 
                v-model="bannerSettings.text"
                type="text" 
                placeholder="Contoh: Diskon Flash Sale 50% Akhir Tahun! Gunakan kode: FLASH50"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Link Tujuan (Opsional)</label>
              <input 
                v-model="bannerSettings.link"
                type="text" 
                placeholder="https://example.com/promo atau /register"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
            </div>
            
            <!-- Color Settings -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Warna Latar Belakang</label>
                <div class="flex items-center gap-2">
                  <input 
                    v-model="bannerSettings.bgColor"
                    type="color" 
                    class="w-12 h-10 rounded border border-neutral-200 cursor-pointer"
                  />
                  <input 
                    v-model="bannerSettings.bgColor"
                    type="text" 
                    class="flex-1 px-3 py-2 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  />
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Warna Teks</label>
                <div class="flex items-center gap-2">
                  <input 
                    v-model="bannerSettings.textColor"
                    type="color" 
                    class="w-12 h-10 rounded border border-neutral-200 cursor-pointer"
                  />
                  <input 
                    v-model="bannerSettings.textColor"
                    type="text" 
                    class="flex-1 px-3 py-2 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  />
                </div>
              </div>
            </div>
            
            <!-- Preview -->
            <div class="mt-8 p-4 bg-neutral-50 rounded-xl border border-neutral-200">
              <p class="text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-2">Preview</p>
              <div v-if="bannerSettings.enabled" class="w-full px-4 py-3 text-sm font-medium text-center rounded-lg shadow-sm" :style="{ backgroundColor: bannerSettings.bgColor, color: bannerSettings.textColor }">
                 {{ bannerSettings.text || 'Teks pengumuman akan muncul di sini' }}
                 <span v-if="bannerSettings.link" class="ml-2 underline opacity-80 cursor-pointer">Cek Sekarang &rarr;</span>
              </div>
              <div v-else class="text-center text-neutral-400 text-sm py-2">
                Banner tidak aktif
              </div>
            </div>

            <div class="pt-4">
              <button type="button" @click="saveBannerSettings" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Pengaturan' }}</button>
            </div>
          </div>
        </div>
        
        <!-- Payment Settings -->
        <div v-if="activeTab === 'payment'" class="space-y-6">
          <!-- Loading Payment -->
          <div v-if="paymentLoading" class="flex items-center justify-center py-8">
            <div class="animate-spin w-8 h-8 border-4 border-admin-500 border-t-transparent rounded-full"></div>
          </div>

          <template v-else>
            <!-- Enable Payment Toggle -->
            <div class="bg-white rounded-xl border border-neutral-200 p-6">
              <div class="flex items-center justify-between">
                <div>
                  <h3 class="font-semibold text-neutral-900">Aktifkan Modul Pembayaran</h3>
                  <p class="text-sm text-neutral-500 mt-1">Aktifkan untuk menggunakan fitur pembayaran kursus berbayar</p>
                </div>
                <label class="relative inline-flex items-center cursor-pointer">
                  <input type="checkbox" v-model="paymentSettings.enabled" @change="savePaymentSettings" class="sr-only peer">
                  <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-admin-600"></div>
                </label>
              </div>
            </div>

            <!-- Provider Selection -->
            <div class="bg-white rounded-xl border border-neutral-200 p-6">
              <h3 class="font-semibold text-neutral-900 mb-4">Payment Provider</h3>
              
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <!-- Midtrans -->
                <div 
                  @click="paymentSettings.provider = 'midtrans'"
                  :class="[
                    'border-2 rounded-xl p-4 cursor-pointer transition-all',
                    paymentSettings.provider === 'midtrans' 
                      ? 'border-admin-500 bg-admin-50' 
                      : 'border-neutral-200 hover:border-neutral-300'
                  ]"
                >
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                      <span class="text-lg font-bold text-blue-600">M</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-neutral-900 text-sm">Midtrans</h4>
                      <p class="text-xs text-neutral-500">Popular</p>
                    </div>
                  </div>
                </div>

                <!-- Duitku -->
                <div 
                  @click="paymentSettings.provider = 'duitku'"
                  :class="[
                    'border-2 rounded-xl p-4 cursor-pointer transition-all',
                    paymentSettings.provider === 'duitku' 
                      ? 'border-admin-500 bg-admin-50' 
                      : 'border-neutral-200 hover:border-neutral-300'
                  ]"
                >
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                      <span class="text-lg font-bold text-green-600">D</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-neutral-900 text-sm">Duitku</h4>
                      <p class="text-xs text-neutral-500">Alternatif</p>
                    </div>
                  </div>
                </div>

                <!-- Xendit (Coming Soon) -->
                <div class="border-2 border-neutral-200 rounded-xl p-4 opacity-50 cursor-not-allowed">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
                      <span class="text-lg font-bold text-purple-600">X</span>
                    </div>
                    <div>
                      <h4 class="font-medium text-neutral-900 text-sm">Xendit</h4>
                      <p class="text-xs text-neutral-500">Coming Soon</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Midtrans Configuration -->
            <div v-if="paymentSettings.provider === 'midtrans'" class="bg-white rounded-xl border border-neutral-200 p-6">
              <h3 class="font-semibold text-neutral-900 mb-4">Konfigurasi Midtrans</h3>
              
              <div class="space-y-4">
                <!-- Environment -->
                <div class="flex items-center gap-4">
                  <label class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" v-model="paymentSettings.midtrans_is_production" class="sr-only peer">
                    <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-amber-500"></div>
                  </label>
                  <div>
                    <span class="font-medium text-neutral-900 text-sm">
                      {{ paymentSettings.midtrans_is_production ? 'Production Mode' : 'Sandbox Mode' }}
                    </span>
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Server Key</label>
                  <input 
                    type="password"
                    v-model="paymentSettings.midtrans_server_key"
                    placeholder="SB-Mid-server-xxxxx atau Mid-server-xxxxx"
                    class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 text-sm"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Client Key</label>
                  <input 
                    type="text"
                    v-model="paymentSettings.midtrans_client_key"
                    placeholder="SB-Mid-client-xxxxx atau Mid-client-xxxxx"
                    class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 text-sm"
                  />
                </div>

                <div class="pt-2">
                  <button @click="savePaymentSettings" :disabled="savingPayment" class="btn-admin">
                    {{ savingPayment ? 'Menyimpan...' : 'Simpan Pengaturan' }}
                  </button>
                </div>
              </div>
            </div>

            <!-- Duitku Configuration -->
            <div v-if="paymentSettings.provider === 'duitku'" class="bg-white rounded-xl border border-neutral-200 p-6">
              <h3 class="font-semibold text-neutral-900 mb-4">Konfigurasi Duitku</h3>
              
              <div class="space-y-4">
                <!-- Environment -->
                <div class="flex items-center gap-4">
                  <label class="relative inline-flex items-center cursor-pointer">
                    <input type="checkbox" v-model="paymentSettings.duitku_is_production" class="sr-only peer">
                    <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-amber-500"></div>
                  </label>
                  <div>
                    <span class="font-medium text-neutral-900 text-sm">
                      {{ paymentSettings.duitku_is_production ? 'Production Mode' : 'Sandbox Mode' }}
                    </span>
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Merchant Code</label>
                  <input 
                    type="text"
                    v-model="paymentSettings.duitku_merchant_code"
                    placeholder="DXXXXX"
                    class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 text-sm"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Merchant Key (API Key)</label>
                  <input 
                    type="password"
                    v-model="paymentSettings.duitku_merchant_key"
                    placeholder="API Key dari Duitku"
                    class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-admin-500 text-sm"
                  />
                </div>

                <div class="pt-2">
                  <button @click="savePaymentSettings" :disabled="savingPayment" class="btn-admin">
                    {{ savingPayment ? 'Menyimpan...' : 'Simpan Pengaturan' }}
                  </button>
                </div>
              </div>
            </div>
          </template>
        </div>
        <!-- AI Tutor Settings -->
        <div v-if="activeTab === 'ai'" class="space-y-6">
          <!-- AI Enable Toggle -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="font-semibold text-neutral-900">AI Tutor</h3>
                <p class="text-sm text-neutral-500 mt-1">Aktifkan fitur AI Tutor untuk membantu siswa memahami materi</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="aiSettings.enabled" @change="saveAISettings" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-accent-600"></div>
              </label>
            </div>
          </div>

          <!-- AI Provider Selection -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6">
            <h3 class="font-semibold text-neutral-900 mb-4">Provider AI</h3>
            <p class="text-sm text-neutral-500 mb-6">Pilih penyedia layanan AI yang akan digunakan</p>
            
            <div class="grid sm:grid-cols-2 gap-4 mb-6">
              <div 
                v-for="provider in aiProviders" 
                :key="provider.id"
                @click="aiSettings.provider = provider.id"
                class="border rounded-xl p-4 cursor-pointer transition-all"
                :class="aiSettings.provider === provider.id ? 'border-admin-500 bg-admin-50' : 'border-neutral-200 hover:border-neutral-300'"
              >
                <div class="flex items-center gap-3">
                  <div 
                    class="w-12 h-12 rounded-lg flex items-center justify-center"
                    :class="provider.color"
                  >
                    <span class="text-white font-bold text-lg">{{ provider.icon }}</span>
                  </div>
                  <div>
                    <p class="font-medium text-neutral-900">{{ provider.name }}</p>
                    <p class="text-xs text-neutral-500">{{ provider.description }}</p>
                  </div>
                </div>
                <div v-if="aiSettings.provider === provider.id" class="mt-3 flex items-center gap-2 text-xs text-admin-600">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  Provider Aktif
                </div>
              </div>
            </div>

            <!-- API Key Input -->
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">API Key {{ selectedProviderName }}</label>
                <div class="relative">
                  <input 
                    v-model="aiSettings.apiKey"
                    :type="showAIKey ? 'text' : 'password'"
                    class="w-full px-4 py-2.5 pr-12 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                    :placeholder="getAPIKeyPlaceholder()"
                  />
                  <button 
                    type="button" 
                    @click="showAIKey = !showAIKey"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                </div>
              </div>

              <!-- Model Selection -->
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Model Chat</label>
                <select 
                  v-model="aiSettings.model"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                >
                  <option v-for="model in getModelsForProvider()" :key="model" :value="model">{{ model }}</option>
                </select>
              </div>
              
              <!-- Embedding Model Selection -->
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Model Embedding</label>
                <select 
                  v-model="aiSettings.embeddingModel"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  :disabled="!supportsEmbedding()"
                >
                  <option v-if="!supportsEmbedding()" value="">Provider tidak mendukung embedding</option>
                  <option v-for="model in getEmbeddingModelsForProvider()" :key="model" :value="model">{{ model }}</option>
                </select>
                <p class="text-xs text-neutral-500 mt-1">Model untuk generate embeddings konten kursus (RAG)</p>
              </div>

              <div class="pt-4 flex gap-3">
                <button type="button" @click="validateAIKey" :disabled="validatingKey" class="px-4 py-2.5 text-sm font-medium text-admin-600 bg-admin-50 rounded-lg hover:bg-admin-100 transition-colors disabled:opacity-50">
                  {{ validatingKey ? 'Validating...' : 'Test API Key' }}
                </button>
                <button type="button" @click="saveAISettings" :disabled="saving" class="btn-admin">
                  {{ saving ? 'Menyimpan...' : 'Simpan Konfigurasi' }}
                </button>
              </div>
            </div>
          </div>

          <!-- AI Parameters -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6">
            <h3 class="font-semibold text-neutral-900 mb-4">Parameter AI</h3>

            <div class="grid sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Max Tokens</label>
                <input 
                  v-model.number="aiSettings.maxTokens"
                  type="number" 
                  min="100"
                  max="8000"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Temperature</label>
                <input 
                  v-model.number="aiSettings.temperature"
                  type="number" 
                  step="0.1"
                  min="0"
                  max="2"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Batas Harian per User</label>
                <input 
                  v-model.number="aiSettings.rateLimitPerDay"
                  type="number" 
                  min="1"
                  max="1000"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                />
              </div>
            </div>
          </div>

          <!-- System Prompt -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6">
            <h3 class="font-semibold text-neutral-900 mb-4">System Prompt</h3>
            <p class="text-sm text-neutral-500 mb-4">Instruksi dan persona untuk AI Tutor</p>
            
            <textarea 
              v-model="aiSettings.systemPrompt"
              rows="8"
              class="w-full px-4 py-3 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm resize-none"
              placeholder="Contoh: Kamu adalah AI Tutor yang membantu siswa memahami materi kursus..."
            ></textarea>

            <div class="pt-4">
              <button type="button" @click="saveAISettings" :disabled="saving" class="btn-admin">
                {{ saving ? 'Menyimpan...' : 'Simpan Semua' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div 
          class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3"
          :class="toast.type === 'success' ? 'bg-accent-600 text-white' : 'bg-red-600 text-white'"
        >
          <svg v-if="toast.type === 'success'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
          <span class="text-sm font-medium">{{ toast.message }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'admin',
  middleware: 'admin'
})

useHead({
  title: 'Pengaturan Sistem - Admin'
})

// API Composable
const { loading, error, settings, fetchSettings, updateSettings } = useSettings()

const activeTab = ref('general')
const saving = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const tabs = [
  { id: 'general', name: 'Umum', description: 'Informasi dasar', icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z' },
  { id: 'theme', name: 'Tema', description: 'Tampilan platform', icon: 'M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zm0 0h12a2 2 0 002-2v-4a2 2 0 00-2-2h-2.343M11 7.343l1.657-1.657a2 2 0 012.828 0l2.829 2.829a2 2 0 010 2.828l-8.486 8.485M7 17h.01' },
  { id: 'payment', name: 'Pembayaran', description: 'Gateway pembayaran', icon: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z' },
  { id: 'ai', name: 'AI Tutor', description: 'Konfigurasi AI', icon: 'M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714c0 .597.237 1.17.659 1.591L19.8 15.3M14.25 3.104c.251.023.501.05.75.082M19.8 15.3l-1.57.393A9.065 9.065 0 0112 15a9.065 9.065 0 00-6.23-.693L5 14.5m14.8.8l1.402 1.402c1.232 1.232.65 3.318-1.067 3.611l-.932.156a2.25 2.25 0 01-2.585-2.586l.156-.931a2.25 2.25 0 00.598-1.652' },
  { id: 'features', name: 'Fitur', description: 'Toggle fitur', icon: 'M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z' },
  { id: 'security', name: 'Keamanan', description: 'Pengaturan keamanan', icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z' },
  { id: 'announcement', name: 'Pengumuman', description: 'Banner informasi', icon: 'M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z' }
]

// Local form state (syncs from API)
const form = ref({
  site_name: '',
  site_description: '',
  contact_email: '',
  theme: 'default',
  logo_url: '',
  require2FA: false,
  sessionTimeout: 30
})

// Logo upload state
const uploadingLogo = ref(false)
const logoInput = ref<HTMLInputElement | null>(null)

// Payment settings state
const paymentLoading = ref(true)
const savingPayment = ref(false)
const paymentSettings = ref({
  enabled: false,
  provider: 'midtrans',
  midtrans_server_key: '',
  midtrans_client_key: '',
  midtrans_is_production: false,
  duitku_merchant_code: '',
  duitku_merchant_key: '',
  duitku_is_production: false
})

// Banner settings state
const bannerSettings = ref({
  enabled: false,
  text: '',
  link: '',
  bgColor: '#1E3A5F',
  textColor: '#FFFFFF'
})

// Fetch banner settings
const fetchBannerSettings = async () => {
    // Sync from global settings if available
    if (settings.value) {
        bannerSettings.value = {
            enabled: settings.value.banner_enabled || false,
            text: settings.value.banner_text || '',
            link: settings.value.banner_link || '',
            bgColor: settings.value.banner_bg_color || '#1E3A5F',
            textColor: settings.value.banner_text_color || '#FFFFFF'
        }
    }
}

// Watch global settings to update banner local state
watch(settings, (newSettings: any) => {
    if (newSettings) {
        bannerSettings.value = {
            enabled: newSettings.banner_enabled || false,
            text: newSettings.banner_text || '',
            link: newSettings.banner_link || '',
            bgColor: newSettings.banner_bg_color || '#1E3A5F',
            textColor: newSettings.banner_text_color || '#FFFFFF'
        }
    }
}, { immediate: true })

const saveBannerSettings = async () => {
  saving.value = true
  const result = await updateSettings({
    banner_enabled: bannerSettings.value.enabled,
    banner_text: bannerSettings.value.text,
    banner_link: bannerSettings.value.link,
    banner_bg_color: bannerSettings.value.bgColor,
    banner_text_color: bannerSettings.value.textColor
  })
  saving.value = false
  
  if (result) {
    showToast('Pengumuman berhasil disimpan')
  } else {
    showToast('Gagal menyimpan pengumuman', 'error')
  }
}

// Fetch payment settings
const fetchPaymentSettings = async () => {
  try {
    const config = useRuntimeConfig()
    const token = useCookie('token')
    const data = await $fetch<any>(`${config.public.apiBase}/api/admin/payment/settings`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    paymentSettings.value = { ...paymentSettings.value, ...data }
  } catch (err) {
    console.error('Failed to fetch payment settings:', err)
  } finally {
    paymentLoading.value = false
  }
}

// Save payment settings
const savePaymentSettings = async () => {
  savingPayment.value = true
  try {
    const config = useRuntimeConfig()
    const token = useCookie('token')
    await $fetch(`${config.public.apiBase}/api/admin/payment/settings`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token.value}` },
      body: paymentSettings.value
    })
    showToast('Pengaturan pembayaran berhasil disimpan')
  } catch (err: any) {
    console.error('Failed to save payment settings:', err)
    showToast(err.data?.error || 'Gagal menyimpan pengaturan', 'error')
  } finally {
    savingPayment.value = false
  }
}

// Get full logo URL
const getLogoUrl = (url: string) => {
  if (!url) return ''
  if (url.startsWith('http')) return url
  const config = useRuntimeConfig()
  return `${config.public.apiBase}${url}`
}

// Upload logo function
const uploadLogo = async (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  
  uploadingLogo.value = true
  
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const config = useRuntimeConfig()
    const token = useCookie('token')
    
    const response = await $fetch<{ url: string }>(`${config.public.apiBase}/api/admin/upload`, {
      method: 'POST',
      body: formData,
      headers: {
        'Authorization': `Bearer ${token.value}`
      }
    })
    
    if (response?.url) {
      form.value.logo_url = response.url
      // Save logo_url to settings immediately
      await updateSettings({ logo_url: response.url })
      showToast('Logo berhasil diupload')
    }
  } catch (err: any) {
    console.error('Upload error:', err)
    showToast('Gagal mengupload logo', 'error')
  } finally {
    uploadingLogo.value = false
    // Reset input
    if (input) input.value = ''
  }
}

const primaryColors = [
  { value: 'blue', class: 'bg-primary-600' },
  { value: 'teal', class: 'bg-admin-600' },
  { value: 'green', class: 'bg-accent-600' },
  { value: 'orange', class: 'bg-warm-500' },
  { value: 'rose', class: 'bg-rose-500' },
  { value: 'indigo', class: 'bg-indigo-600' }
]

const features = ref([
  { id: 1, name: 'Sistem Kuis', description: 'Ujian dan kuis untuk siswa', enabled: true },
  { id: 2, name: 'E-Sertifikat', description: 'Sertifikat otomatis setelah selesai kursus', enabled: true },
  { id: 3, name: 'Forum Diskusi', description: 'Forum tanya jawab antar siswa', enabled: false },
  { id: 4, name: 'Live Chat', description: 'Chat langsung dengan instruktur', enabled: false },
  { id: 5, name: 'Drip Content', description: 'Buka materi berdasarkan jadwal', enabled: true },
  { id: 6, name: 'Sistem Afiliasi', description: 'Program referral untuk pengguna', enabled: false }
])

// AI Settings
const showAIKey = ref(false)
const validatingKey = ref(false)
const aiSettings = ref({
  enabled: false,
  provider: 'openai',
  model: 'gpt-4-turbo',
  apiKey: '',
  maxTokens: 2048,
  temperature: 0.7,
  rateLimitPerDay: 50,
  systemPrompt: '',
  embeddingModel: ''
})

const aiProviders = [
  { id: 'openai', name: 'OpenAI', description: 'GPT-4, GPT-3.5', color: 'bg-emerald-600', icon: '✦' },
  { id: 'claude', name: 'Anthropic Claude', description: 'Claude 3.5, Claude 3', color: 'bg-amber-600', icon: 'A' },
  { id: 'groq', name: 'Groq', description: 'LLaMA 3, Mixtral', color: 'bg-blue-600', icon: 'G' },
  { id: 'gemini', name: 'Google Gemini', description: 'Gemini Pro, Gemini Flash', color: 'bg-purple-600', icon: '◆' }
]

const aiModels: Record<string, string[]> = {
  openai: ['gpt-4-turbo', 'gpt-4', 'gpt-4o', 'gpt-4o-mini', 'gpt-3.5-turbo'],
  claude: ['claude-3-5-sonnet-20241022', 'claude-3-opus-20240229', 'claude-3-sonnet-20240229', 'claude-3-haiku-20240307'],
  groq: ['llama-3.3-70b-versatile', 'llama-3.1-70b-versatile', 'llama-3.1-8b-instant', 'mixtral-8x7b-32768'],
  gemini: ['gemini-3-flash', 'gemini-2.5-flash', 'gemini-2.5-flash-lite']
}

// Embedding models per provider (only some providers support embedding)
const aiEmbeddingModels: Record<string, string[]> = {
  openai: ['text-embedding-ada-002', 'text-embedding-3-small', 'text-embedding-3-large'],
  gemini: ['gemini-embedding-001', 'text-embedding-004']
}

const selectedProviderName = computed(() => {
  const provider = aiProviders.find(p => p.id === aiSettings.value.provider)
  return provider?.name || ''
})

const getModelsForProvider = () => {
  return aiModels[aiSettings.value.provider] || []
}

// Check if current provider supports embedding
const supportsEmbedding = () => {
  return Object.keys(aiEmbeddingModels).includes(aiSettings.value.provider)
}

// Get embedding models for current provider
const getEmbeddingModelsForProvider = () => {
  return aiEmbeddingModels[aiSettings.value.provider] || []
}

const getAPIKeyPlaceholder = () => {
  switch (aiSettings.value.provider) {
    case 'openai': return 'sk-...'
    case 'claude': return 'sk-ant-...'
    case 'groq': return 'gsk_...'
    case 'gemini': return 'AI...'
    default: return 'API Key'
  }
}

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

const saveGeneralSettings = async () => {
  saving.value = true
  const result = await updateSettings({
    site_name: form.value.site_name,
    site_description: form.value.site_description,
    contact_email: form.value.contact_email
  })
  saving.value = false
  
  if (result) {
    showToast('Pengaturan berhasil disimpan')
  } else {
    showToast('Gagal menyimpan pengaturan', 'error')
  }
}

const saveThemeSettings = async () => {
  saving.value = true
  const result = await updateSettings({
    theme: form.value.theme
  })
  saving.value = false
  
  if (result) {
    showToast('Tema berhasil disimpan')
  } else {
    showToast('Gagal menyimpan tema', 'error')
  }
}

// AI Settings Functions
const fetchAISettings = async () => {
  try {
    const config = useRuntimeConfig()
    const token = useCookie('token')
    const data = await $fetch<any>(`${config.public.apiBase}/api/admin/ai/settings`, {
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    
    // Determine which provider's key to show based on active provider
    let currentApiKey = ''
    const provider = data.provider || 'openai'
    
    if (provider === 'openai' && data.openai_configured) {
      currentApiKey = data.api_key_openai || ''
    } else if (provider === 'claude' && data.claude_configured) {
      currentApiKey = data.api_key_claude || ''
    } else if (provider === 'groq' && data.groq_configured) {
      currentApiKey = data.api_key_groq || ''
    } else if (provider === 'gemini' && data.gemini_configured) {
      currentApiKey = data.api_key_gemini || ''
    }
    
    aiSettings.value = {
      enabled: data.enabled,
      provider: provider,
      model: data.model || 'gpt-4-turbo',
      apiKey: currentApiKey,
      maxTokens: data.max_tokens || 2048,
      temperature: data.temperature || 0.7,
      rateLimitPerDay: data.rate_limit_per_day || 50,
      systemPrompt: data.system_prompt || '',
      embeddingModel: data.embedding_model || ''
    }
  } catch (err) {
    console.error('Failed to fetch AI settings:', err)
  }
}

const saveAISettings = async () => {
  saving.value = true
  try {
    const config = useRuntimeConfig()
    const token = useCookie('token')
    
    const payload: any = {
      enabled: aiSettings.value.enabled,
      provider: aiSettings.value.provider,
      model: aiSettings.value.model,
      embedding_model: aiSettings.value.embeddingModel, // Separate embedding model
      max_tokens: aiSettings.value.maxTokens,
      temperature: aiSettings.value.temperature,
      rate_limit_per_day: aiSettings.value.rateLimitPerDay,
      system_prompt: aiSettings.value.systemPrompt
    }
    
    // Add API key based on provider
    const apiKeyField = `api_key_${aiSettings.value.provider}`
    if (aiSettings.value.apiKey && !aiSettings.value.apiKey.startsWith('****')) {
      payload[apiKeyField] = aiSettings.value.apiKey
    }
    
    await $fetch(`${config.public.apiBase}/api/admin/ai/settings`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${token.value}` },
      body: payload
    })
    
    showToast('Pengaturan AI berhasil disimpan')
  } catch (err) {
    console.error('Failed to save AI settings:', err)
    showToast('Gagal menyimpan pengaturan AI', 'error')
  } finally {
    saving.value = false
  }
}

const validateAIKey = async () => {
  if (!aiSettings.value.apiKey) {
    showToast('Masukkan API Key terlebih dahulu', 'error')
    return
  }
  
  validatingKey.value = true
  try {
    const config = useRuntimeConfig()
    const token = useCookie('token')
    
    const result = await $fetch<{ valid: boolean; message: string }>(`${config.public.apiBase}/api/admin/ai/validate-key`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token.value}` },
      body: {
        provider: aiSettings.value.provider,
        api_key: aiSettings.value.apiKey
      }
    })
    
    if (result.valid) {
      showToast('API Key valid! ✓')
    } else {
      showToast(result.message || 'API Key tidak valid', 'error')
    }
  } catch (err: any) {
    showToast('Gagal memvalidasi API Key', 'error')
  } finally {
    validatingKey.value = false
  }
}

// Load settings on mount
onMounted(async () => {
  // Load AI settings
  await fetchAISettings()
  await fetchSettings()
  await fetchPaymentSettings()
  if (settings.value) {
    form.value.site_name = settings.value.site_name || ''
    form.value.site_description = settings.value.site_description || ''
    form.value.contact_email = settings.value.contact_email || ''
    form.value.theme = settings.value.theme || 'default'
    form.value.logo_url = settings.value.logo_url || ''
  }
})
</script>
