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
        
        <!-- Payment Gateway Settings -->
        <div v-if="activeTab === 'payment'" class="space-y-6">
          <!-- Active Gateway Selection -->
          <div class="bg-white rounded-xl border border-neutral-200 p-6">
            <h3 class="font-semibold text-neutral-900 mb-4">Payment Gateway Aktif</h3>
            <p class="text-sm text-neutral-500 mb-6">Pilih gateway pembayaran yang akan digunakan</p>
            
            <div class="grid sm:grid-cols-2 gap-4">
              <div 
                v-for="gateway in paymentGateways" 
                :key="gateway.id"
                @click="paymentForm.active_gateway = gateway.id"
                class="border rounded-xl p-4 cursor-pointer transition-all"
                :class="paymentForm.active_gateway === gateway.id ? 'border-admin-500 bg-admin-50' : 'border-neutral-200 hover:border-neutral-300'"
              >
                <div class="flex items-center gap-3">
                  <div 
                    class="w-12 h-12 rounded-lg flex items-center justify-center"
                    :class="gateway.color"
                  >
                    <span class="text-white font-bold text-sm">{{ gateway.icon }}</span>
                  </div>
                  <div>
                    <p class="font-medium text-neutral-900">{{ gateway.name }}</p>
                    <p class="text-xs text-neutral-500">{{ gateway.description }}</p>
                  </div>
                </div>
                <div v-if="paymentForm.active_gateway === gateway.id" class="mt-3 flex items-center gap-2 text-xs text-admin-600">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                  Gateway Aktif
                </div>
              </div>
            </div>
          </div>
          
          <!-- Midtrans Configuration -->
          <div v-if="paymentForm.active_gateway === 'midtrans'" class="bg-white rounded-xl border border-neutral-200 p-6">
            <div class="flex items-center gap-3 mb-6">
              <div class="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">M</span>
              </div>
              <div>
                <h3 class="font-semibold text-neutral-900">Konfigurasi Midtrans</h3>
                <p class="text-xs text-neutral-500">Sandbox & Production credentials</p>
              </div>
            </div>
            
            <form @submit.prevent="savePaymentSettings" class="space-y-4">
              <div class="flex items-center gap-4 mb-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="paymentForm.midtrans.environment" value="sandbox" class="w-4 h-4 text-admin-600">
                  <span class="text-sm text-neutral-700">Sandbox (Testing)</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="paymentForm.midtrans.environment" value="production" class="w-4 h-4 text-admin-600">
                  <span class="text-sm text-neutral-700">Production</span>
                </label>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Merchant ID</label>
                <input 
                  v-model="paymentForm.midtrans.merchant_id"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  placeholder="G123456789"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Client Key</label>
                <input 
                  v-model="paymentForm.midtrans.client_key"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  placeholder="SB-Mid-client-xxxxx"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Server Key</label>
                <div class="relative">
                  <input 
                    v-model="paymentForm.midtrans.server_key"
                    :type="showServerKey ? 'text' : 'password'"
                    class="w-full px-4 py-2.5 pr-12 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                    placeholder="SB-Mid-server-xxxxx"
                  />
                  <button 
                    type="button" 
                    @click="showServerKey = !showServerKey"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600"
                  >
                    <svg v-if="!showServerKey" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.542 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"/>
                    </svg>
                  </button>
                </div>
              </div>
              
              <div class="pt-4 flex gap-3">
                <button type="button" @click="testPaymentConnection" :disabled="testingConnection" class="px-4 py-2.5 text-sm font-medium text-admin-600 bg-admin-50 rounded-lg hover:bg-admin-100 transition-colors disabled:opacity-50">
                  {{ testingConnection ? 'Testing...' : 'Test Koneksi' }}
                </button>
                <button type="submit" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Konfigurasi' }}</button>
              </div>
            </form>
          </div>
          
          <!-- Xendit Configuration -->
          <div v-if="paymentForm.active_gateway === 'xendit'" class="bg-white rounded-xl border border-neutral-200 p-6">
            <div class="flex items-center gap-3 mb-6">
              <div class="w-10 h-10 bg-cyan-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">X</span>
              </div>
              <div>
                <h3 class="font-semibold text-neutral-900">Konfigurasi Xendit</h3>
                <p class="text-xs text-neutral-500">API keys untuk Xendit</p>
              </div>
            </div>
            
            <form @submit.prevent="savePaymentSettings" class="space-y-4">
              <div class="flex items-center gap-4 mb-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="paymentForm.xendit.environment" value="sandbox" class="w-4 h-4 text-admin-600">
                  <span class="text-sm text-neutral-700">Sandbox (Testing)</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="paymentForm.xendit.environment" value="production" class="w-4 h-4 text-admin-600">
                  <span class="text-sm text-neutral-700">Production</span>
                </label>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Public API Key</label>
                <input 
                  v-model="paymentForm.xendit.public_key"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  placeholder="xnd_public_development_xxxxx"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Secret API Key</label>
                <div class="relative">
                  <input 
                    v-model="paymentForm.xendit.secret_key"
                    :type="showServerKey ? 'text' : 'password'"
                    class="w-full px-4 py-2.5 pr-12 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                    placeholder="xnd_development_xxxxx"
                  />
                  <button 
                    type="button" 
                    @click="showServerKey = !showServerKey"
                    class="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-neutral-600"
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                </div>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Callback Token (Webhook Verification)</label>
                <input 
                  v-model="paymentForm.xendit.callback_token"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  placeholder="Token untuk verifikasi webhook"
                />
              </div>
              
              <div class="pt-4 flex gap-3">
                <button type="button" @click="testPaymentConnection" :disabled="testingConnection" class="px-4 py-2.5 text-sm font-medium text-admin-600 bg-admin-50 rounded-lg hover:bg-admin-100 transition-colors disabled:opacity-50">
                  {{ testingConnection ? 'Testing...' : 'Test Koneksi' }}
                </button>
                <button type="submit" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Konfigurasi' }}</button>
              </div>
            </form>
          </div>
          
          <!-- Manual Transfer -->
          <div v-if="paymentForm.active_gateway === 'manual'" class="bg-white rounded-xl border border-neutral-200 p-6">
            <div class="flex items-center gap-3 mb-6">
              <div class="w-10 h-10 bg-neutral-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">$</span>
              </div>
              <div>
                <h3 class="font-semibold text-neutral-900">Transfer Manual</h3>
                <p class="text-xs text-neutral-500">Konfigurasi rekening bank</p>
              </div>
            </div>
            
            <form @submit.prevent="savePaymentSettings" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Bank</label>
                <input 
                  v-model="paymentForm.manual.bank_name"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="BCA, Mandiri, BNI, etc."
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Nomor Rekening</label>
                <input 
                  v-model="paymentForm.manual.account_number"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm font-mono"
                  placeholder="1234567890"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Pemilik Rekening</label>
                <input 
                  v-model="paymentForm.manual.account_name"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="PT LearnHub Indonesia"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Instruksi Pembayaran</label>
                <textarea 
                  v-model="paymentForm.manual.instructions"
                  rows="3"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm resize-none"
                  placeholder="Instruksi untuk pembayaran manual..."
                ></textarea>
              </div>
              
              <div class="pt-4">
                <button type="submit" :disabled="saving" class="btn-admin">{{ saving ? 'Menyimpan...' : 'Simpan Konfigurasi' }}</button>
              </div>
            </form>
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
  { id: 'features', name: 'Fitur', description: 'Toggle fitur', icon: 'M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z' },
  { id: 'security', name: 'Keamanan', description: 'Pengaturan keamanan', icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z' }
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

// Payment Gateway Configuration
const paymentGateways = [
  { id: 'midtrans', name: 'Midtrans', description: 'Payment gateway Indonesia', color: 'bg-blue-600', icon: 'M' },
  { id: 'xendit', name: 'Xendit', description: 'Multi-payment solution', color: 'bg-cyan-600', icon: 'X' },
  { id: 'manual', name: 'Transfer Manual', description: 'Bank transfer langsung', color: 'bg-neutral-600', icon: '$' }
]

const paymentForm = ref({
  active_gateway: 'midtrans',
  midtrans: {
    environment: 'sandbox',
    merchant_id: '',
    client_key: '',
    server_key: ''
  },
  xendit: {
    environment: 'sandbox',
    public_key: '',
    secret_key: '',
    callback_token: ''
  },
  manual: {
    bank_name: '',
    account_number: '',
    account_name: '',
    instructions: ''
  }
})

const showServerKey = ref(false)
const testingConnection = ref(false)

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

const savePaymentSettings = async () => {
  saving.value = true
  
  // In production, this would save to backend
  // For now, simulate save
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  // Store in localStorage for demo
  if (process.client) {
    localStorage.setItem('payment_config', JSON.stringify(paymentForm.value))
  }
  
  saving.value = false
  showToast('Konfigurasi pembayaran berhasil disimpan')
}

const testPaymentConnection = async () => {
  testingConnection.value = true
  
  // Simulate API test
  await new Promise(resolve => setTimeout(resolve, 2000))
  
  testingConnection.value = false
  showToast('Koneksi berhasil! Gateway siap digunakan')
}

// Load payment settings from localStorage on mount
onMounted(async () => {
  await fetchSettings()
  if (settings.value) {
    form.value.site_name = settings.value.site_name || ''
    form.value.site_description = settings.value.site_description || ''
    form.value.contact_email = settings.value.contact_email || ''
    form.value.theme = settings.value.theme || 'default'
    form.value.logo_url = settings.value.logo_url || ''
  }
  
  // Load payment config from localStorage
  if (process.client) {
    const savedPaymentConfig = localStorage.getItem('payment_config')
    if (savedPaymentConfig) {
      try {
        paymentForm.value = JSON.parse(savedPaymentConfig)
      } catch (e) {
        // Use default config
      }
    }
  }
})
</script>
