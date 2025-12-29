<template>
  <div v-if="show" class="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-2xl w-full max-w-md max-h-[90vh] overflow-y-auto">
      <!-- Header -->
      <div class="flex items-center justify-between p-5 border-b">
        <div>
          <h3 class="text-xl font-bold text-neutral-900">Checkout</h3>
          <p class="text-sm text-neutral-500">{{ props.courseName }}</p>
        </div>
        <button @click="close" class="text-neutral-400 hover:text-neutral-600 text-2xl">&times;</button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="p-5 space-y-5">
        <!-- Email -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-1.5">Email *</label>
          <input 
            v-model="form.email" 
            type="email" 
            required
            placeholder="email@gmail.com"
            class="w-full px-4 py-3 border border-neutral-300 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          />
          <p class="text-xs text-neutral-500 mt-1">Email akan digunakan untuk akses akun LMS Anda</p>
        </div>

        <!-- Phone (WhatsApp) -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-1.5">Nomor WhatsApp *</label>
          <input 
            v-model="form.phone" 
            type="tel" 
            required
            placeholder="08123456789"
            class="w-full px-4 py-3 border border-neutral-300 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          />
        </div>

        <!-- Name (Optional) -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-1.5">Nama Lengkap (opsional)</label>
          <input 
            v-model="form.fullName" 
            type="text" 
            placeholder="John Doe"
            class="w-full px-4 py-3 border border-neutral-300 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          />
        </div>

        <!-- Payment Method -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-2">Metode Pembayaran *</label>
          
          <div v-if="loadingMethods" class="flex items-center justify-center py-8">
            <div class="animate-spin rounded-full h-6 w-6 border-2 border-primary-500 border-t-transparent"></div>
          </div>

          <div v-else-if="paymentMethods.length" class="space-y-2 max-h-48 overflow-y-auto">
            <label 
              v-for="method in paymentMethods" 
              :key="method.paymentMethod"
              class="flex items-center gap-3 p-3 border rounded-xl cursor-pointer transition-all"
              :class="form.paymentMethod === method.paymentMethod ? 'border-primary-500 bg-primary-50' : 'border-neutral-200 hover:border-neutral-300'"
            >
              <input 
                type="radio" 
                :value="method.paymentMethod" 
                v-model="form.paymentMethod"
                class="text-primary-600"
              />
              <img v-if="method.paymentImage" :src="method.paymentImage" class="h-6 w-auto" :alt="method.paymentName"/>
              <div class="flex-1">
                <div class="font-medium text-sm">{{ method.paymentName }}</div>
                <div v-if="method.totalFee && method.totalFee !== '0'" class="text-xs text-neutral-500">
                  Fee: Rp {{ formatNumber(parseInt(method.totalFee)) }}
                </div>
              </div>
            </label>
          </div>

          <div v-else class="text-center py-4 text-neutral-500 text-sm">
            Tidak ada metode pembayaran tersedia
          </div>
        </div>

        <!-- Coupon Code -->
        <div>
          <label class="block text-sm font-medium text-neutral-700 mb-1.5">Kode Kupon (opsional)</label>
          <input 
            v-model="form.couponCode" 
            type="text" 
            placeholder="Masukkan kode kupon"
            class="w-full px-4 py-3 border border-neutral-300 rounded-xl focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          />
        </div>

        <!-- Price Summary -->
        <div class="bg-neutral-50 rounded-xl p-4 space-y-2">
          <div class="flex justify-between text-sm">
            <span class="text-neutral-600">Harga</span>
            <span>Rp {{ formatNumber(props.coursePrice) }}</span>
          </div>
          <div v-if="props.discountPrice && props.discountPrice < props.coursePrice" class="flex justify-between text-sm">
            <span class="text-neutral-600">Diskon</span>
            <span class="text-green-600">-Rp {{ formatNumber(props.coursePrice - props.discountPrice) }}</span>
          </div>
          <div class="flex justify-between font-bold text-lg pt-2 border-t border-neutral-200">
            <span>Total</span>
            <span class="text-primary-600">Rp {{ formatNumber(props.discountPrice || props.coursePrice) }}</span>
          </div>
        </div>

        <!-- Error Message -->
        <div v-if="error" class="bg-red-50 text-red-600 p-3 rounded-lg text-sm">
          {{ error }}
        </div>

        <!-- Submit Button -->
        <button 
          type="submit"
          :disabled="submitting || !form.paymentMethod"
          class="w-full py-4 bg-primary-600 text-white font-bold rounded-xl hover:bg-primary-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
        >
          <span v-if="submitting" class="flex items-center justify-center gap-2">
            <div class="animate-spin rounded-full h-5 w-5 border-2 border-white border-t-transparent"></div>
            Memproses...
          </span>
          <span v-else>Bayar Sekarang</span>
        </button>

        <!-- Terms -->
        <p class="text-xs text-neutral-500 text-center">
          Dengan melanjutkan, Anda menyetujui 
          <a href="/terms" class="text-primary-600 hover:underline">Syarat & Ketentuan</a>
        </p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
interface PaymentMethod {
  paymentMethod: string
  paymentName: string
  paymentImage: string
  totalFee: string
}

interface Props {
  show: boolean
  campaignId: string
  courseName: string
  coursePrice: number
  discountPrice?: number
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'success'])

const config = useRuntimeConfig()
const apiBase = config.public.apiBase || 'http://localhost:8080'

const form = ref({
  email: '',
  phone: '',
  fullName: '',
  paymentMethod: '',
  couponCode: ''
})

const paymentMethods = ref<PaymentMethod[]>([])
const loadingMethods = ref(false)
const submitting = ref(false)
const error = ref('')

const formatNumber = (num: number) => new Intl.NumberFormat('id-ID').format(num)

const close = () => {
  error.value = ''
  emit('close')
}

const fetchPaymentMethods = async () => {
  loadingMethods.value = true
  try {
    const amount = props.discountPrice || props.coursePrice
    const response = await $fetch<{ payment_methods: PaymentMethod[] }>(`/api/public/payment-methods?amount=${amount}`, {
      baseURL: apiBase
    })
    paymentMethods.value = response.payment_methods || []
  } catch (err: any) {
    console.error('Failed to fetch payment methods:', err)
    error.value = 'Gagal memuat metode pembayaran'
  } finally {
    loadingMethods.value = false
  }
}

const handleSubmit = async () => {
  error.value = ''
  
  if (!form.value.email || !form.value.phone || !form.value.paymentMethod) {
    error.value = 'Mohon lengkapi semua field yang wajib diisi'
    return
  }

  submitting.value = true
  try {
    const response = await $fetch<{
      payment_url?: string
      is_free?: boolean
      message?: string
      order_id?: string
    }>('/api/campaign-checkout', {
      method: 'POST',
      baseURL: apiBase,
      body: {
        campaign_id: props.campaignId,
        email: form.value.email,
        phone: form.value.phone,
        full_name: form.value.fullName,
        payment_method: form.value.paymentMethod,
        coupon_code: form.value.couponCode || undefined
      }
    })

    if (response.is_free) {
      // Free enrollment - show success
      emit('success', { isFree: true, message: response.message })
      close()
    } else if (response.payment_url) {
      // Redirect to payment page
      window.location.href = response.payment_url
    } else {
      error.value = 'Tidak dapat memproses pembayaran'
    }
  } catch (err: any) {
    console.error('Checkout failed:', err)
    error.value = err.data?.error || err.message || 'Gagal memproses checkout'
  } finally {
    submitting.value = false
  }
}

// Fetch payment methods when modal opens
watch(() => props.show, (newVal) => {
  if (newVal) {
    fetchPaymentMethods()
  }
})
</script>
