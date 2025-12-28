<template>
  <div>
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div>
        <h1 class="text-2xl font-bold text-neutral-900">Manajemen Kursus</h1>
        <p class="text-neutral-500 mt-1">Kelola semua kursus platform</p>
      </div>
      <div class="flex gap-3">
        <button @click="handleExport" class="px-4 py-2.5 text-sm font-medium text-neutral-600 bg-white border border-neutral-200 rounded-lg hover:bg-neutral-50 transition-colors flex items-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          Export CSV
        </button>
        <button @click="openAddModal" class="btn-admin">
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          Tambah Kursus
        </button>
      </div>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total Kursus</p>
        <p class="text-2xl font-bold text-neutral-900">{{ courses.length }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Dipublikasi</p>
        <p class="text-2xl font-bold text-accent-600">{{ publishedCount }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Draft</p>
        <p class="text-2xl font-bold text-warm-600">{{ draftCount }}</p>
      </div>
      <div class="bg-white rounded-xl p-5 border border-neutral-200">
        <p class="text-sm text-neutral-500 mb-1">Total</p>
        <p class="text-2xl font-bold text-neutral-900">{{ total }}</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <div class="relative flex-1">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
        </svg>
        <input 
          v-model="searchQuery"
          type="text" 
          placeholder="Cari kursus..." 
          class="w-full pl-10 pr-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
        />
      </div>
      <select v-model="categoryFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Kategori</option>
        <option v-for="cat in categoriesList" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
      <select v-model="statusFilter" class="px-4 py-2.5 bg-white border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm">
        <option value="">Semua Status</option>
        <option value="published">Published</option>
        <option value="draft">Draft</option>
      </select>
    </div>

    <!-- Courses Table -->
    <div class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-neutral-50 border-b border-neutral-200">
            <tr>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kursus</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Kategori</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Instruktur</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Siswa</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Rating</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Harga</th>
              <th class="px-6 py-4 text-left text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
              <th class="px-6 py-4 text-right text-xs font-semibold text-neutral-600 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-neutral-100">
            <tr v-for="course in filteredCourses" :key="course.id" class="hover:bg-neutral-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div 
                    v-if="course.thumbnail_url" 
                    class="w-12 h-12 rounded-lg overflow-hidden flex-shrink-0"
                  >
                    <img :src="getThumbnailUrl(course.thumbnail_url)" :alt="course.title" class="w-full h-full object-cover" />
                  </div>
                  <div v-else :class="['w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0', getCourseColor(course.id)]">
                    <svg class="w-6 h-6 text-white/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                    </svg>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-neutral-900">{{ course.title }}</p>
                    <p class="text-xs text-neutral-500">{{ course.description?.substring(0, 50) || 'No description' }}...</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-neutral-600">{{ course.category?.name || '-' }}</td>
              <td class="px-6 py-4 text-sm text-neutral-600">{{ course.instructor?.full_name || '-' }}</td>
              <td class="px-6 py-4 text-sm text-neutral-900 font-medium">-</td>
              <td class="px-6 py-4">
                <div v-if="getCourseRating(course.id)" class="flex items-center gap-1">
                  <svg class="w-4 h-4 text-yellow-400 fill-current" viewBox="0 0 24 24">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  <span class="text-sm font-medium text-neutral-900">{{ getCourseRating(course.id).toFixed(1) }}</span>
                  <span class="text-xs text-neutral-500">({{ getCourseRatingCount(course.id) }})</span>
                </div>
                <span v-else class="text-sm text-neutral-400">-</span>
              </td>
              <td class="px-6 py-4">
                <div v-if="course.discount_price && isDiscountActive(course)" class="flex flex-col">
                  <span class="text-xs text-neutral-400 line-through">{{ formatCurrency(course.price) }}</span>
                  <span class="text-sm font-bold text-red-600">{{ formatCurrency(course.discount_price) }}</span>
                </div>
                <span v-else class="text-sm text-neutral-900 font-medium">{{ formatCurrency(course.price) }}</span>
              </td>
              <td class="px-6 py-4">
                <button 
                  @click="handleToggleStatus(course)"
                  class="px-2.5 py-1 text-xs font-medium rounded-full cursor-pointer transition-colors"
                  :class="course.is_published ? 'bg-accent-100 text-accent-700 hover:bg-accent-200' : 'bg-warm-100 text-warm-700 hover:bg-warm-200'"
                >
                  {{ course.is_published ? 'Published' : 'Draft' }}
                </button>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-2">
                  <NuxtLink :to="`/admin/courses/${course.id}/materials`" class="p-2 text-neutral-400 hover:text-accent-600 hover:bg-accent-50 rounded-lg transition-colors" title="Kelola Materi">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
                    </svg>
                  </NuxtLink>
                  <button @click="viewCourse(course)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                  <button @click="openEditModal(course)" class="p-2 text-neutral-400 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                    </svg>
                  </button>
                  <button @click="openDeleteModal(course)" class="p-2 text-neutral-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Add/Edit Course Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-lg m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-neutral-900">{{ isEditing ? 'Edit Kursus' : 'Tambah Kursus' }}</h3>
            <button @click="closeModal" class="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveCourse" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Kursus <span class="text-red-500">*</span></label>
              <input 
                v-model="form.title"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Masukkan judul kursus"
              />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Kategori <span class="text-red-500">*</span></label>
                <select 
                  v-model="form.category_id"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                >
                  <option value="">Pilih kategori</option>
                  <option v-for="cat in categoriesList" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Instruktur <span class="text-red-500">*</span></label>
                <select 
                  v-model="form.instructor_id"
                  required
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                >
                  <option value="">Pilih instruktur</option>
                  <option v-for="inst in instructors" :key="inst.id" :value="inst.id">{{ inst.full_name }}</option>
                </select>
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Jumlah Modul</label>
                <input 
                  v-model.number="form.lessons"
                  type="number" 
                  min="1"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="0"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Durasi</label>
                <input 
                  v-model="form.duration"
                  type="text" 
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="Contoh: 12 jam"
                />
              </div>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Harga (Rp) <span class="text-red-500">*</span></label>
                <input 
                  v-model.number="form.price"
                  type="number" 
                  required
                  min="0"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  placeholder="Contoh: 499000"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-2">Harga Diskon (Rp)</label>
                <input 
                  v-model.number="form.discount_price"
                  type="number" 
                  min="0"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500 text-sm"
                  placeholder="Kosongkan jika tidak ada"
                />
              </div>
            </div>
            <div v-if="form.discount_price">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Diskon Berakhir</label>
              <input 
                v-model="form.discount_valid_until"
                type="datetime-local"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
              />
              <p class="text-xs text-neutral-500 mt-1">Kosongkan jika diskon berlaku selamanya</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Thumbnail Kursus</label>
              <!-- Thumbnail Preview -->
              <div class="mb-3">
                <div 
                  class="w-full h-40 rounded-lg border-2 border-dashed border-neutral-300 flex flex-col items-center justify-center overflow-hidden relative"
                  :class="{'border-admin-400': isDragging}"
                  @dragover.prevent="isDragging = true"
                  @dragleave="isDragging = false"
                  @drop.prevent="handleFileDrop"
                >
                  <template v-if="thumbnailPreview || form.thumbnail_url">
                    <img 
                      :src="thumbnailPreview || getThumbnailUrl(form.thumbnail_url)" 
                      alt="Thumbnail Preview" 
                      class="w-full h-full object-cover"
                    />
                    <button 
                      type="button"
                      @click="clearThumbnail"
                      class="absolute top-2 right-2 p-1.5 bg-red-500 text-white rounded-full hover:bg-red-600 transition-colors"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                      </svg>
                    </button>
                  </template>
                  <template v-else>
                    <div :class="['w-16 h-16 rounded-lg flex items-center justify-center mb-2', form.color || 'bg-neutral-200']">
                      <svg class="w-8 h-8 text-white/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                    <p class="text-sm text-neutral-500">Drag & drop atau klik untuk upload</p>
                    <p class="text-xs text-neutral-400">PNG, JPG maksimal 2MB</p>
                  </template>
                  <input 
                    type="file" 
                    ref="thumbnailInput"
                    accept="image/png,image/jpeg,image/webp"
                    @change="handleFileSelect"
                    class="absolute inset-0 opacity-0 cursor-pointer"
                  />
                </div>
              </div>
              <!-- URL input alternative -->
              <div class="flex items-center gap-2 text-xs text-neutral-500 mb-2">
                <span>atau masukkan URL gambar:</span>
              </div>
              <input 
                v-model="form.thumbnail_url"
                type="url"
                placeholder="https://example.com/image.jpg"
                class="w-full px-3 py-2 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                @input="onThumbnailUrlChange"
              />
              <!-- Fallback color picker -->
              <div class="mt-3">
                <p class="text-xs text-neutral-500 mb-2">Warna fallback (jika tidak ada gambar):</p>
                <div class="flex gap-2">
                  <button 
                    v-for="color in colorOptions" 
                    :key="color.value"
                    type="button"
                    @click="form.color = color.value"
                    class="w-7 h-7 rounded-md transition-all"
                    :class="[color.class, form.color === color.value ? 'ring-2 ring-offset-1 ring-neutral-900' : '']"
                  ></button>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Status</label>
              <div class="flex items-center gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="form.status" value="draft" class="w-4 h-4 text-admin-600 focus:ring-admin-500">
                  <span class="text-sm text-neutral-700">Draft</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="form.status" value="published" class="w-4 h-4 text-admin-600 focus:ring-admin-500">
                  <span class="text-sm text-neutral-700">Published</span>
                </label>
              </div>
            </div>
            
            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors">
                {{ isEditing ? 'Simpan Perubahan' : 'Tambah Kursus' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- View Course Modal -->
    <Transition name="fade">
      <div v-if="showViewModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showViewModal = false">
        <div class="bg-white rounded-xl w-full max-w-lg m-4 overflow-hidden" @click.stop>
          <div :class="['h-32 flex items-center justify-center', selectedCourse?.color]">
            <svg class="w-16 h-16 text-white/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
            </svg>
          </div>
          <div class="p-6">
            <h3 class="text-xl font-bold text-neutral-900 mb-2">{{ selectedCourse?.title }}</h3>
            <div class="grid grid-cols-2 gap-4 text-sm mb-4">
              <div>
                <p class="text-neutral-500">Kategori</p>
                <p class="font-medium text-neutral-900">{{ selectedCourse?.category?.name || '-' }}</p>
              </div>
              <div>
                <p class="text-neutral-500">Instruktur</p>
                <p class="font-medium text-neutral-900">{{ selectedCourse?.instructor?.full_name || '-' }}</p>
              </div>
              <div>
                <p class="text-neutral-500">Modul</p>
                <p class="font-medium text-neutral-900">{{ selectedCourse?.lessons_count || 0 }} Modul</p>
              </div>
              <div>
                <p class="text-neutral-500">Durasi</p>
                <p class="font-medium text-neutral-900">{{ selectedCourse?.duration || '-' }}</p>
              </div>
              <div>
                <p class="text-neutral-500">Siswa</p>
                <p class="font-medium text-neutral-900">{{ (selectedCourse?.students || 0).toLocaleString() }}</p>
              </div>
              <div>
                <p class="text-neutral-500">Harga</p>
                <p class="font-medium text-neutral-900">Rp {{ (selectedCourse?.price || 0).toLocaleString() }}</p>
              </div>
            </div>
            <button @click="showViewModal = false" class="w-full py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Tutup
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Delete Confirmation Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-sm m-4" @click.stop>
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-8 h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-neutral-900 mb-2">Hapus Kursus</h3>
            <p class="text-sm text-neutral-500 mb-6">Apakah Anda yakin ingin menghapus <strong>{{ selectedCourse?.title }}</strong>? Semua data kursus akan dihapus permanen.</p>
          </div>
          <div class="flex gap-3">
            <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
              Batal
            </button>
            <button @click="confirmDeleteCourse" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50">
              Hapus
            </button>
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
  title: 'Manajemen Kursus - Admin'
})

// API composable
const { 
  loading, 
  error, 
  courses, 
  total, 
  fetchAdminCourses, 
  createCourse: createCourseApi, 
  updateCourse: updateCourseApi, 
  deleteCourse: deleteCourseApi,
  togglePublish
} = useCourses()

const { instructors, fetchInstructors } = useInstructors()
const { categories: categoriesList, fetchCategories } = useCategories()
const { exportCourses } = useExport()

// Rating stats for courses
const api = useApi()
const ratingStats = ref<Record<string, { average_rating: number, total_ratings: number }>>({})

const fetchCourseRatingStats = async () => {
  if (courses.value.length === 0) return
  
  try {
    // Fetch stats for each course
    for (const course of courses.value) {
      const stats = await api.fetch<{ average_rating: number, total_ratings: number }>(`/api/courses/${course.id}/ratings/stats`)
      if (stats) {
        ratingStats.value[course.id] = stats
      }
    }
  } catch (err) {
    console.error('Failed to fetch rating stats:', err)
  }
}

const getCourseRating = (courseId: string): number | null => {
  const stats = ratingStats.value[courseId]
  return stats?.average_rating > 0 ? stats.average_rating : null
}

const getCourseRatingCount = (courseId: string): number => {
  return ratingStats.value[courseId]?.total_ratings || 0
}

const handleExport = () => {
  if (filteredCourses.value.length === 0) {
    showToast('Tidak ada data untuk di-export', 'error')
    return
  }
  exportCourses(filteredCourses.value)
  showToast('Data berhasil di-export ke CSV')
}

const searchQuery = ref('')
const categoryFilter = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const showViewModal = ref(false)
const showDeleteModal = ref(false)
const isEditing = ref(false)
const selectedCourse = ref<any>(null)
const saving = ref(false)

const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

// const categories = ['Web Development', 'Design', 'Data Science', 'Mobile', 'Business', 'Marketing']
const colorOptions = [
  { value: 'bg-primary-600', class: 'bg-primary-600' },
  { value: 'bg-accent-600', class: 'bg-accent-600' },
  { value: 'bg-warm-500', class: 'bg-warm-500' },
  { value: 'bg-cyan-600', class: 'bg-cyan-600' },
  { value: 'bg-rose-500', class: 'bg-rose-500' },
  { value: 'bg-indigo-600', class: 'bg-indigo-600' }
]

const form = ref({
  title: '',
  description: '',
  price: 0,
  discount_price: null as number | null,
  discount_valid_until: '',
  status: 'draft' as 'draft' | 'published',
  instructor_id: '',
  category_id: '',
  lessons: 0,
  duration: '',
  thumbnail_url: '',
  color: 'bg-primary-600'
})

// Helper to check if discount is active
const isDiscountActive = (course: any): boolean => {
  if (!course.discount_price) return false
  if (!course.discount_valid_until) return true
  return new Date(course.discount_valid_until) > new Date()
}

// Thumbnail handling
const thumbnailInput = ref<HTMLInputElement | null>(null)
const thumbnailPreview = ref<string | null>(null)
const thumbnailFile = ref<File | null>(null)
const isDragging = ref(false)

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  if (input.files && input.files[0]) {
    processFile(input.files[0])
  }
}

const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false
  if (event.dataTransfer?.files && event.dataTransfer.files[0]) {
    processFile(event.dataTransfer.files[0])
  }
}

const processFile = (file: File) => {
  // Validate file type
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type)) {
    showToast('Format file tidak didukung. Gunakan PNG, JPG, atau WebP', 'error')
    return
  }
  
  // Validate file size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    showToast('Ukuran file terlalu besar. Maksimal 2MB', 'error')
    return
  }
  
  thumbnailFile.value = file
  
  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    thumbnailPreview.value = e.target?.result as string
    form.value.thumbnail_url = '' // Clear URL if file is selected
  }
  reader.readAsDataURL(file)
}

const clearThumbnail = () => {
  thumbnailPreview.value = null
  thumbnailFile.value = null
  form.value.thumbnail_url = ''
  if (thumbnailInput.value) {
    thumbnailInput.value.value = ''
  }
}

const onThumbnailUrlChange = () => {
  // Clear file preview if URL is entered
  if (form.value.thumbnail_url) {
    thumbnailPreview.value = null
    thumbnailFile.value = null
  }
}

// Load courses on mount
onMounted(async () => {
  await loadCourses()
})

const loadCourses = async () => {
  await Promise.all([
    fetchAdminCourses({ limit: 100 }),
    fetchInstructors(),
    fetchCategories()
  ])
  // Fetch rating stats after courses are loaded
  await fetchCourseRatingStats()
}

const publishedCount = computed(() => courses.value.filter(c => c.is_published).length)
const draftCount = computed(() => courses.value.filter(c => !c.is_published).length)

const filteredCourses = computed(() => {
  return courses.value.filter(course => {
    const matchesSearch = course.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = !statusFilter.value || 
      (statusFilter.value === 'published' ? course.is_published : !course.is_published)
    const matchesCategory = !categoryFilter.value || course.category_id === categoryFilter.value
    return matchesSearch && matchesStatus && matchesCategory
  })
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => { toast.value.show = false }, 3000)
}

const formatCurrency = (amount: number) => {
  return `Rp ${amount.toLocaleString('id-ID')}`
}

const getCourseColor = (id: string) => {
  const colors = ['bg-primary-600', 'bg-accent-600', 'bg-warm-500', 'bg-cyan-600', 'bg-rose-500', 'bg-indigo-600']
  const hash = (id || '').split('').reduce((acc, char) => acc + char.charCodeAt(0), 0)
  return colors[hash % colors.length]
}

// Get full thumbnail URL
const getThumbnailUrl = (url: string | null | undefined) => {
  if (!url) return ''
  // Full URL - return as-is
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  // Legacy /uploads path - prepend apiBase
  if (url.startsWith('/uploads')) {
    const config = useRuntimeConfig()
    return `${config.public.apiBase}${url}`
  }
  // MinIO object key - use public images endpoint
  const config = useRuntimeConfig()
  return `${config.public.apiBase}/api/images/${url}`
}

const openAddModal = () => {
  isEditing.value = false
  form.value = { 
    title: '', 
    description: '', 
    price: 0,
    discount_price: null,
    discount_valid_until: '',
    status: 'draft',
    instructor_id: '',
    category_id: '',
    lessons: 0,
    duration: '',
    thumbnail_url: '',
    color: 'bg-primary-600'
  }
  // Clear thumbnail preview
  thumbnailPreview.value = null
  thumbnailFile.value = null
  showModal.value = true
}

const openEditModal = (course: any) => {
  isEditing.value = true
  selectedCourse.value = course
  form.value = { 
    title: course.title,
    description: course.description || '',
    price: course.price,
    discount_price: course.discount_price || null,
    discount_valid_until: course.discount_valid_until ? new Date(course.discount_valid_until).toISOString().slice(0, 16) : '',
    status: course.is_published ? 'published' : 'draft',
    instructor_id: course.instructor_id || '',
    category_id: course.category_id || '',
    lessons: course.lessons_count || 0,
    duration: course.duration || '',
    thumbnail_url: course.thumbnail_url || '',
    color: 'bg-primary-600'
  }
  // Clear file preview when editing
  thumbnailPreview.value = null
  thumbnailFile.value = null
  showModal.value = true
}

const viewCourse = (course: any) => {
  selectedCourse.value = course
  showViewModal.value = true
}

const openDeleteModal = (course: any) => {
  selectedCourse.value = course
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
}

const handleToggleStatus = async (course: any) => {
  const result = await togglePublish(course.id)
  if (result) {
    showToast(`Status kursus berhasil diubah`)
    await loadCourses()
  } else {
    showToast('Gagal mengubah status', 'error')
  }
}

const saveCourse = async () => {
  saving.value = true
  try {
    let thumbnailUrl = form.value.thumbnail_url || null
    
    // Upload thumbnail file if selected
    if (thumbnailFile.value) {
      const formData = new FormData()
      formData.append('file', thumbnailFile.value)
      
      const config = useRuntimeConfig()
      const token = useCookie('token')
      
      try {
        const uploadResponse = await $fetch<{ url: string }>(`${config.public.apiBase}/api/admin/upload`, {
          method: 'POST',
          body: formData,
          headers: {
            'Authorization': `Bearer ${token.value}`
          }
        })
        
        if (uploadResponse?.url) {
          thumbnailUrl = uploadResponse.url
        }
      } catch (uploadErr) {
        console.error('Thumbnail upload error:', uploadErr)
        showToast('Gagal mengupload thumbnail', 'error')
        saving.value = false
        return
      }
    }
    
    // Convert form status to is_published boolean for API
    const courseData = {
      title: form.value.title,
      description: form.value.description,
      price: form.value.price,
      discount_price: form.value.discount_price || null,
      discount_valid_until: form.value.discount_valid_until ? new Date(form.value.discount_valid_until).toISOString() : null,
      is_published: form.value.status === 'published',
      instructor_id: form.value.instructor_id,
      category_id: form.value.category_id,
      lessons_count: form.value.lessons,
      duration: form.value.duration,
      thumbnail_url: thumbnailUrl,
    }

    if (isEditing.value && selectedCourse.value) {
      const result = await updateCourseApi(selectedCourse.value.id, courseData)
      if (result) {
        showToast('Kursus berhasil diperbarui')
        await loadCourses()
      } else {
        showToast('Gagal memperbarui kursus', 'error')
      }
    } else {
      const result = await createCourseApi(courseData)
      if (result) {
        showToast('Kursus berhasil ditambahkan')
        await loadCourses()
      } else {
        showToast('Gagal menambahkan kursus', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}

const confirmDeleteCourse = async () => {
  if (!selectedCourse.value) return
  
  saving.value = true
  try {
    const success = await deleteCourseApi(selectedCourse.value.id)
    if (success) {
      showToast('Kursus berhasil dihapus')
      await loadCourses()
    } else {
      showToast('Gagal menghapus kursus', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedCourse.value = null
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.slide-up-enter-active, .slide-up-leave-active { transition: all 0.3s ease; }
.slide-up-enter-from, .slide-up-leave-to { opacity: 0; transform: translateY(20px); }
</style>
