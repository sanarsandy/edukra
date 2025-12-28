<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Manajemen Instruktur</h1>
        <p class="text-neutral-500 mt-1">Kelola instruktur dan pengajar platform</p>
      </div>
      <button @click="openAddModal" class="btn-admin w-full sm:w-auto">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
        </svg>
        Tambah Instruktur
      </button>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total Instruktur</p>
        <p class="text-2xl font-bold text-neutral-900">{{ instructors.length }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total Kursus</p>
        <p class="text-2xl font-bold text-primary-600">{{ instructors.reduce((acc, i) => acc + (i.course_count || 0), 0) }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total Siswa</p>
        <p class="text-2xl font-bold text-accent-600">{{ instructors.reduce((acc, i) => acc + (i.students || 0), 0).toLocaleString() }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Rata-rata Rating</p>
        <p class="text-2xl font-bold text-warm-600">{{ averageRating }}</p>
      </div>
    </div>

    <!-- Instructors Grid -->
    <div class="grid md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div v-for="instructor in instructors" :key="instructor.id" class="bg-white rounded-xl border border-neutral-200 overflow-hidden hover:shadow-lg transition-all">
        <div class="p-6">
          <div class="flex items-start gap-4 mb-4">
            <div :class="['w-16 h-16 rounded-full flex items-center justify-center text-white text-xl font-bold', instructor.avatarColor]">
              {{ instructor.full_name.split(' ').map(n => n[0]).join('').substring(0, 2) }}
            </div>
            <div class="flex-1">
              <h3 class="font-semibold text-neutral-900">{{ instructor.full_name }}</h3>
              <p class="text-sm text-neutral-500">{{ instructor.email }}</p>
              <p class="text-xs text-neutral-400 mt-1">{{ instructor.specialty }}</p>
            </div>
          </div>
          
          <div class="grid grid-cols-3 gap-4 mb-4 py-4 border-y border-neutral-100">
            <div class="text-center">
              <p class="text-lg font-bold text-neutral-900">{{ instructor.course_count || 0 }}</p>
              <p class="text-xs text-neutral-500">Kursus</p>
            </div>
            <div class="text-center">
              <p class="text-lg font-bold text-neutral-900">{{ (instructor.students || 0).toLocaleString() }}</p>
              <p class="text-xs text-neutral-500">Siswa</p>
            </div>
            <div class="text-center">
              <div class="flex items-center justify-center gap-1">
                <svg class="w-4 h-4 text-warm-500 fill-current" viewBox="0 0 20 20">
                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                </svg>
                <span class="text-lg font-bold text-neutral-900">{{ instructor.rating || '0.0' }}</span>
              </div>
              <p class="text-xs text-neutral-500">Rating</p>
            </div>
          </div>
          
          <div class="flex gap-2">
            <button @click="openViewModal(instructor)" class="flex-1 py-2 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Lihat Profil
            </button>
            <button @click="openEditModal(instructor)" class="flex-1 py-2 text-sm font-medium text-admin-600 bg-admin-50 rounded-lg hover:bg-admin-100 transition-colors">
              Kelola
            </button>
            <button @click="openDeleteModal(instructor)" class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Add Instructor Card -->
      <div 
        @click="openAddModal"
        class="bg-neutral-50 rounded-xl border-2 border-dashed border-neutral-300 flex flex-col items-center justify-center p-8 cursor-pointer hover:border-admin-500 hover:bg-admin-50 transition-all group min-h-[280px]"
      >
        <div class="w-16 h-16 bg-neutral-200 group-hover:bg-admin-100 rounded-full flex items-center justify-center mb-4 transition-colors">
          <svg class="w-8 h-8 text-neutral-400 group-hover:text-admin-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
          </svg>
        </div>
        <p class="text-sm font-medium text-neutral-500 group-hover:text-admin-600 transition-colors">Tambah Instruktur Baru</p>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-neutral-900">{{ isEditing ? 'Edit Instruktur' : 'Tambah Instruktur' }}</h3>
            <button @click="closeModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveInstructor" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Nama Lengkap <span class="text-red-500">*</span></label>
              <input 
                v-model="form.name"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan nama lengkap"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Email <span class="text-red-500">*</span></label>
              <input 
                v-model="form.email"
                type="email" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="email@learnhub.id"
              />
            </div>
            <div v-if="!isEditing">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Password <span class="text-red-500">*</span></label>
              <input 
                v-model="form.password"
                type="password" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan password"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Spesialisasi <span class="text-red-500">*</span></label>
              <select 
                v-model="form.specialty"
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              >
                <option value="">Pilih spesialisasi</option>
                <option v-for="spec in specialties" :key="spec" :value="spec">{{ spec }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Bio</label>
              <textarea 
                v-model="form.bio"
                rows="3"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm resize-none"
                placeholder="Deskripsi singkat tentang instruktur"
              ></textarea>
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors">
                {{ isEditing ? 'Simpan Perubahan' : 'Tambah Instruktur' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- View Profile Modal -->
    <Transition name="fade">
      <div v-if="showViewModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showViewModal = false">
        <div class="bg-white rounded-xl w-full max-w-md m-4 overflow-hidden" @click.stop>
          <div :class="['h-24 relative', selectedInstructor?.avatarColor]"></div>
          <div class="px-6 pb-6 -mt-12">
            <div :class="['w-24 h-24 rounded-full flex items-center justify-center text-white text-3xl font-bold border-4 border-white shadow-lg', selectedInstructor?.avatarColor]">
              {{ selectedInstructor?.full_name.split(' ').map(n => n[0]).join('').substring(0, 2) }}
            </div>
            <h3 class="text-xl font-bold text-neutral-900 mt-4">{{ selectedInstructor?.full_name }}</h3>
            <p class="text-sm text-neutral-500">{{ selectedInstructor?.email }}</p>
            <span class="inline-block mt-2 px-3 py-1 bg-admin-100 text-admin-700 text-xs font-medium rounded-full">{{ selectedInstructor?.specialty }}</span>
            
            <p v-if="selectedInstructor?.bio" class="text-sm text-neutral-600 mt-4">{{ selectedInstructor?.bio }}</p>
            
            <div class="grid grid-cols-3 gap-4 mt-6 py-4 border-y border-neutral-100">
              <div class="text-center">
                <p class="text-2xl font-bold text-neutral-900">{{ selectedInstructor?.course_count || 0 }}</p>
                <p class="text-xs text-neutral-500">Kursus</p>
              </div>
              <div class="text-center">
                <p class="text-2xl font-bold text-neutral-900">{{ (selectedInstructor?.students || 0).toLocaleString() }}</p>
                <p class="text-xs text-neutral-500">Siswa</p>
              </div>
              <div class="text-center">
                <div class="flex items-center justify-center gap-1">
                  <svg class="w-5 h-5 text-warm-500 fill-current" viewBox="0 0 20 20">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                  </svg>
                  <span class="text-2xl font-bold text-neutral-900">{{ selectedInstructor?.rating || '0.0' }}</span>
                </div>
                <p class="text-xs text-neutral-500">Rating</p>
              </div>
            </div>
            
            <button @click="showViewModal = false" class="w-full mt-4 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Delete Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Instruktur</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus <strong>{{ selectedInstructor?.full_name }}</strong>?</p>
          </div>
          <div class="flex gap-3">
            <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">Batal</button>
            <button @click="deleteInstructor" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors">Hapus</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast -->
    <Transition name="slide-up">
      <div v-if="toast.show" class="fixed bottom-6 right-6 z-50">
        <div class="px-4 py-3 rounded-lg shadow-lg flex items-center gap-3 bg-accent-600 text-white">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
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
  title: 'Manajemen Instruktur - Admin'
})

// API Composable
const { 
  loading, 
  error, 
  instructors, 
  total, 
  fetchInstructors, 
  createInstructor: createInstructorApi, 
  updateInstructor: updateInstructorApi, 
  deleteInstructor: deleteInstructorApi 
} = useInstructors()

const showModal = ref(false)
const showViewModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const selectedInstructor = ref<any>(null)
const saving = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

const specialties = ['Web Development', 'UI/UX Design', 'Data Science', 'Mobile Development', 'Business', 'Marketing']
const avatarColors = ['bg-primary-600', 'bg-accent-600', 'bg-warm-600', 'bg-rose-500', 'bg-cyan-600', 'bg-indigo-600']

const form = ref({
  name: '',
  email: '',
  password: '',
  specialty: '',
  bio: ''
})

// Add avatar color to each instructor 
const enrichInstructors = () => {
  instructors.value = instructors.value.map((instructor, index) => ({
    ...instructor,
    avatarColor: avatarColors[index % avatarColors.length],
    // Ensure specialty is displayed from backend if empty
    specialty: instructor.specialty || 'Belum Ditentukan'
  }))
}

// Load instructors on mount
onMounted(async () => {
  await fetchInstructors()
  enrichInstructors()
})

const averageRating = computed(() => {
  if (instructors.value.length === 0) return '0.0'
  const ratingsWithData = instructors.value.filter((i: any) => i.rating && i.rating > 0)
  if (ratingsWithData.length === 0) return '0.0'
  const total = ratingsWithData.reduce((acc: number, i: any) => acc + (i.rating || 0), 0)
  return (total / ratingsWithData.length).toFixed(1)
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const openAddModal = () => {
  isEditing.value = false
  form.value = { name: '', email: '', password: '', specialty: '', bio: '' }
  showModal.value = true
}

const openEditModal = (instructor: any) => {
  isEditing.value = true
  selectedInstructor.value = instructor
  form.value = { 
    name: instructor.full_name, 
    email: instructor.email, 
    password: '', // Password not editable directly here usually
    specialty: instructor.specialty || '', 
    bio: instructor.bio || '' 
  }
  showModal.value = true
}

const openViewModal = (instructor: any) => {
  selectedInstructor.value = instructor
  showViewModal.value = true
}

const openDeleteModal = (instructor: any) => {
  selectedInstructor.value = instructor
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const saveInstructor = async () => {
  saving.value = true
  try {
    if (isEditing.value && selectedInstructor.value) {
      const result = await updateInstructorApi(selectedInstructor.value.id, {
        full_name: form.value.name,
        specialty: form.value.specialty,
        bio: form.value.bio
      })
      if (result) {
        showToast('Instruktur berhasil diperbarui')
        await fetchInstructors()
        enrichInstructors()
      } else {
        showToast('Gagal memperbarui instruktur', 'error')
      }
    } else {
      if (!form.value.password) {
        showToast('Password wajib diisi', 'error')
        saving.value = false
        return
      }
      
      const result = await createInstructorApi({
        full_name: form.value.name,
        email: form.value.email,
        password: form.value.password,
        specialty: form.value.specialty,
        bio: form.value.bio
      })
      
      if (result) {
        showToast('Instruktur berhasil ditambahkan')
        await fetchInstructors()
        enrichInstructors()
      } else {
        showToast('Gagal menambahkan instruktur', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}

const deleteInstructor = async () => {
  if (!selectedInstructor.value) return
  
  saving.value = true
  try {
    const success = await deleteInstructorApi(selectedInstructor.value.id)
    if (success) {
      showToast('Instruktur berhasil dihapus')
      await fetchInstructors()
      enrichInstructors()
    } else {
      showToast('Gagal menghapus instruktur', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedInstructor.value = null
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
