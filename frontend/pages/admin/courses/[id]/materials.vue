<template>
  <div>
    <!-- Header with Back Button -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-8">
      <div class="flex items-center gap-4">
        <NuxtLink :to="'/admin/courses'" class="w-10 h-10 flex items-center justify-center bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
          <svg class="w-5 h-5 text-neutral-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7"/>
          </svg>
        </NuxtLink>
        <div>
          <h1 class="text-2xl font-bold text-neutral-900">Materi Kursus</h1>
          <p class="text-neutral-500 mt-1">{{ course?.title || 'Loading...' }}</p>
        </div>
      </div>
      <div class="flex gap-3">
        <button @click="openAddModuleModal" class="btn-secondary flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
          </svg>
          Tambah Modul
        </button>
        <button @click="openQuizEditor" class="btn-secondary flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          Tambah Kuis
        </button>
        <button @click="openAddModal()" class="btn-admin flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          Tambah Materi
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-4 border-admin-600 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="lessonsTree.length === 0" class="bg-white rounded-xl border border-neutral-200 p-12 text-center">
      <div class="w-16 h-16 bg-neutral-100 rounded-full flex items-center justify-center mx-auto mb-4">
        <svg class="w-8 h-8 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/>
        </svg>
      </div>
      <h3 class="text-lg font-semibold text-neutral-900 mb-2">Belum ada materi</h3>
      <p class="text-neutral-500 mb-6">Tambahkan modul atau materi pertama untuk kursus ini</p>
      <div class="flex gap-3 justify-center">
        <button @click="openAddModuleModal" class="btn-secondary">Tambah Modul</button>
        <button @click="openAddModal()" class="btn-admin">Tambah Materi</button>
      </div>
    </div>

    <!-- Materials Tree View -->
    <div v-else class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="p-2">
        <AdminLessonTreeItem
          v-for="lesson in lessonsTree"
          :key="lesson.id"
          :lesson="lesson"
          @view="openViewModal"
          @edit="openEditModal"
          @delete="openDeleteModal"
          @add-child="openAddModal"
          @add-submodule="openAddSubmoduleModal"
          @add-quiz="openQuizEditorWithParent"
          @manage-quiz="openQuizForLesson"
        />
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Transition name="fade">
      <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="closeModal">
        <div class="bg-white rounded-xl p-6 w-full max-w-lg m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex items-center justify-between mb-6">
            <h2 class="text-xl font-bold text-neutral-900">{{ isEditing ? 'Edit Materi' : 'Tambah Materi' }}</h2>
            <button @click="closeModal" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors">
              <svg class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
              </svg>
            </button>
          </div>
          
          <form @submit.prevent="saveLesson" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Judul Materi <span class="text-red-500">*</span></label>
              <input 
                v-model="form.title"
                type="text" 
                required
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                placeholder="Contoh: Pengenalan HTML"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Deskripsi</label>
              <textarea 
                v-model="form.description"
                rows="3"
                class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm resize-none"
                placeholder="Deskripsi singkat materi"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-2">Tipe Konten <span class="text-red-500">*</span></label>
              <div class="grid grid-cols-4 gap-2">
                <button 
                  v-for="type in contentTypes" 
                  :key="type.value"
                  type="button"
                  @click="form.content_type = type.value"
                  class="flex flex-col items-center gap-1 p-3 border rounded-lg transition-all"
                  :class="form.content_type === type.value ? 'border-admin-500 bg-admin-50 text-admin-600' : 'border-neutral-200 hover:border-neutral-300'"
                >
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" :d="type.icon"/>
                  </svg>
                  <span class="text-xs font-medium">{{ type.label }}</span>
                </button>
              </div>
            </div>

            <!-- Video URL or File Upload -->
            <div v-if="form.content_type === 'video' || form.content_type === 'pdf'">
              <label class="block text-sm font-medium text-neutral-700 mb-2">
                {{ form.content_type === 'video' ? 'Video' : 'Dokumen' }}
              </label>
              
              <!-- Tab: URL or Upload -->
              <div class="flex gap-2 mb-3">
                <button 
                  type="button" 
                  @click="uploadMode = 'url'" 
                  class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
                  :class="uploadMode === 'url' ? 'bg-admin-100 text-admin-700' : 'bg-neutral-100 text-neutral-600 hover:bg-neutral-200'"
                >
                  URL / Link
                </button>
                <button 
                  type="button" 
                  @click="uploadMode = 'file'" 
                  class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors"
                  :class="uploadMode === 'file' ? 'bg-admin-100 text-admin-700' : 'bg-neutral-100 text-neutral-600 hover:bg-neutral-200'"
                >
                  Upload File
                </button>
              </div>

              <!-- URL Input -->
              <div v-if="uploadMode === 'url'">
                <input 
                  v-model="form.video_url"
                  type="url"
                  class="w-full px-4 py-2.5 border border-neutral-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-admin-500 text-sm"
                  :placeholder="form.content_type === 'video' ? 'https://youtube.com/watch?v=...' : 'https://example.com/dokumen.pdf'"
                />
                <!-- Detected Platform Badge -->
                <div v-if="form.content_type === 'video' && form.video_url" class="mt-2 flex items-center gap-2">
                  <span class="text-xs text-neutral-500">Terdeteksi:</span>
                  <span class="px-2 py-0.5 rounded text-xs font-medium" :class="getVideoProviderBadgeClass(form.video_url)">
                    {{ getVideoProviderName(form.video_url) }}
                  </span>
                  <svg class="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                </div>
                <div v-if="form.content_type === 'video' && !form.video_url" class="mt-2 p-2 bg-neutral-50 rounded-lg">
                  <p class="text-xs font-medium text-neutral-600 mb-1">Platform yang didukung:</p>
                  <div class="flex flex-wrap gap-1">
                    <span class="px-2 py-0.5 bg-red-100 text-red-700 rounded text-xs">YouTube</span>
                    <span class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">Vimeo</span>
                    <span class="px-2 py-0.5 bg-green-100 text-green-700 rounded text-xs">Google Drive</span>
                    <span class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">Zoom</span>
                    <span class="px-2 py-0.5 bg-purple-100 text-purple-700 rounded text-xs">MS Stream</span>
                    <span class="px-2 py-0.5 bg-teal-100 text-teal-700 rounded text-xs">Teams/SharePoint</span>
                    <span class="px-2 py-0.5 bg-neutral-100 text-neutral-700 rounded text-xs">Direct URL (.mp4)</span>
                  </div>
                </div>
                <p v-else-if="form.content_type !== 'video'" class="text-xs text-neutral-500 mt-1">Masukkan link dokumen PDF</p>
              </div>

              <!-- File Upload -->
              <div v-else>
                <div 
                  class="border-2 border-dashed border-neutral-300 rounded-lg p-6 text-center hover:border-admin-400 transition-colors cursor-pointer"
                  @click="triggerFileInput"
                  @dragover.prevent
                  @drop.prevent="handleFileDrop"
                >
                  <input 
                    ref="fileInputRef"
                    type="file" 
                    class="hidden" 
                    :accept="form.content_type === 'video' ? 'video/*' : '.pdf,.doc,.docx,.ppt,.pptx'"
                    @change="handleFileSelect"
                  />
                  <div v-if="uploading" class="flex flex-col items-center">
                    <div class="w-8 h-8 border-4 border-admin-600 border-t-transparent rounded-full animate-spin mb-2"></div>
                    <p class="text-sm text-neutral-600">Mengupload {{ uploadProgress }}%</p>
                  </div>
                  <div v-else-if="form.video_url && uploadMode === 'file'" class="flex flex-col items-center">
                    <svg class="w-10 h-10 text-accent-500 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M5 13l4 4L19 7"/>
                    </svg>
                    <p class="text-sm text-neutral-600 font-medium">File berhasil diupload</p>
                    <p class="text-xs text-neutral-500 mt-1 truncate max-w-full">{{ form.video_url }}</p>
                  </div>
                  <div v-else class="flex flex-col items-center">
                    <svg class="w-10 h-10 text-neutral-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                    </svg>
                    <p class="text-sm text-neutral-600">Klik atau drag file ke sini</p>
                    <p class="text-xs text-neutral-500 mt-1">
                      {{ form.content_type === 'video' ? 'MP4, WebM, MOV (Maks. 500MB)' : 'PDF, DOC, DOCX, PPT (Maks. 50MB)' }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Rich Text Editor for Text Content -->
            <div v-if="form.content_type === 'text'">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Konten Teks</label>
              <div class="border border-neutral-200 rounded-lg overflow-hidden">
                <!-- Simple toolbar -->
                <div class="flex flex-wrap gap-1 p-2 bg-neutral-50 border-b border-neutral-200">
                  <button type="button" @click="formatText('bold')" class="p-2 rounded hover:bg-neutral-200 transition-colors" title="Bold">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
                      <path d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z"/>
                      <path d="M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z"/>
                    </svg>
                  </button>
                  <button type="button" @click="formatText('italic')" class="p-2 rounded hover:bg-neutral-200 transition-colors" title="Italic">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                      <path d="M19 4h-9m4 16H5M15 4L9 20"/>
                    </svg>
                  </button>
                  <button type="button" @click="formatText('underline')" class="p-2 rounded hover:bg-neutral-200 transition-colors" title="Underline">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                      <path d="M6 3v7a6 6 0 006 6 6 6 0 006-6V3M4 21h16"/>
                    </svg>
                  </button>
                  <div class="w-px h-6 bg-neutral-300 mx-1 self-center"></div>
                  <button type="button" @click="formatText('insertUnorderedList')" class="p-2 rounded hover:bg-neutral-200 transition-colors" title="Bullet List">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                      <path d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/>
                    </svg>
                  </button>
                  <button type="button" @click="formatText('insertOrderedList')" class="p-2 rounded hover:bg-neutral-200 transition-colors" title="Numbered List">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2">
                      <path d="M10 6h11M10 12h11M10 18h11M4 6h1v4M4 10h2M6 18H4c0-1 2-2 2-3s-1-1.5-2-1"/>
                    </svg>
                  </button>
                  <div class="w-px h-6 bg-neutral-300 mx-1 self-center"></div>
                  <button type="button" @click="formatText('formatBlock', 'h2')" class="p-2 rounded hover:bg-neutral-200 transition-colors text-xs font-bold" title="Heading">
                    H2
                  </button>
                  <button type="button" @click="formatText('formatBlock', 'h3')" class="p-2 rounded hover:bg-neutral-200 transition-colors text-xs font-bold" title="Subheading">
                    H3
                  </button>
                </div>
                <!-- Editable content area -->
                <div 
                  ref="textEditorRef"
                  contenteditable="true"
                  class="min-h-[200px] p-4 prose prose-sm max-w-none focus:outline-none"
                  @input="onTextEditorInput"
                  @paste="onTextEditorPaste"
                ></div>
              </div>
              <p class="text-xs text-neutral-500 mt-1">Gunakan toolbar untuk memformat teks</p>
            </div>

            <div class="flex items-center gap-3">
              <label class="relative inline-flex items-center cursor-pointer">
                <input type="checkbox" v-model="form.is_preview" class="sr-only peer">
                <div class="w-11 h-6 bg-neutral-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-admin-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-neutral-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-admin-600"></div>
              </label>
              <div>
                <p class="text-sm font-medium text-neutral-900">Materi Preview</p>
                <p class="text-xs text-neutral-500">Bisa dilihat tanpa beli kursus</p>
              </div>
            </div>

            <div class="flex gap-3 pt-4">
              <button type="button" @click="closeModal" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button type="submit" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-admin-600 rounded-lg hover:bg-admin-700 transition-colors disabled:opacity-50">
                {{ saving ? 'Menyimpan...' : (isEditing ? 'Simpan Perubahan' : 'Tambah Materi') }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- Delete Confirmation Modal -->
    <Transition name="fade">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showDeleteModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-md m-4" @click.stop>
          <div class="text-center">
            <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
              </svg>
            </div>
            <h3 class="text-lg font-bold text-neutral-900 mb-2">Hapus Materi?</h3>
            <p class="text-neutral-500 mb-6">Materi "{{ selectedLesson?.title }}" akan dihapus permanen.</p>
            <div class="flex gap-3">
              <button @click="showDeleteModal = false" class="flex-1 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
                Batal
              </button>
              <button @click="confirmDeleteLesson" :disabled="saving" class="flex-1 py-2.5 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50">
                {{ saving ? 'Menghapus...' : 'Hapus' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- View Material Modal -->
    <Transition name="fade">
      <div v-if="showViewModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showViewModal = false">
        <div 
          :class="[
            'bg-white rounded-xl p-6 overflow-y-auto transition-all duration-300',
            isFullscreen ? 'w-full h-full max-w-none m-0 rounded-none' : 'w-full max-w-3xl m-4 max-h-[90vh]'
          ]" 
          @click.stop
        >
          <div class="flex items-center justify-between mb-6">
            <div>
              <span class="px-2 py-0.5 rounded-full text-xs font-medium mr-2" :class="getContentTypeBadgeClass(viewLesson?.content_type || 'text')">
                {{ getContentTypeLabel(viewLesson?.content_type || 'text') }}
              </span>
              <h2 class="text-xl font-bold text-neutral-900 mt-2">{{ viewLesson?.title }}</h2>
            </div>
            <div class="flex items-center gap-2">
              <!-- Fullscreen Toggle -->
              <button @click="toggleFullscreen" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors" :title="isFullscreen ? 'Keluar Fullscreen' : 'Fullscreen'">
                <svg v-if="!isFullscreen" class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"/>
                </svg>
                <svg v-else class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25"/>
                </svg>
              </button>
              <!-- Close Button -->
              <button @click="closeViewModal" class="p-2 hover:bg-neutral-100 rounded-lg transition-colors">
                <svg class="w-5 h-5 text-neutral-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
          </div>
          
          <!-- Content based on type -->
          <div class="mb-4">
            <!-- Video Content -->
            <div v-if="viewLesson?.content_type === 'video'" class="space-y-4">
              <div v-if="viewLesson?.video_url">
                <!-- Video Provider Badge -->
                <div class="mb-2">
                  <span class="px-2 py-1 rounded-full text-xs font-medium" :class="getVideoProviderBadgeClass(viewLesson.video_url)">
                    {{ getVideoProviderName(viewLesson.video_url) }}
                  </span>
                </div>
                
                <div class="aspect-video rounded-lg overflow-hidden bg-neutral-900">
                  <!-- Embed Video (YouTube, Vimeo, GDrive, Zoom, Teams, etc) -->
                  <iframe 
                    v-if="isEmbedVideo(viewLesson.video_url)" 
                    :src="getEmbedUrl(viewLesson.video_url)" 
                    class="w-full h-full"
                    frameborder="0" 
                    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; fullscreen" 
                    allowfullscreen
                  ></iframe>
                  <!-- Direct Video File -->
                  <video v-else :src="getFullUrl(viewLesson.video_url)" controls class="w-full h-full"></video>
                </div>
                
                <!-- Video URL Info -->
                <div class="text-xs text-neutral-500 truncate">
                  <span class="font-medium">URL:</span> {{ viewLesson.video_url }}
                </div>
              </div>
              <p v-else class="text-neutral-500 text-center py-8">Tidak ada video</p>
            </div>

            <!-- PDF/Document Content -->
            <div v-else-if="viewLesson?.content_type === 'pdf'" class="space-y-4">
              <div v-if="viewLesson?.video_url">
                <div class="rounded-lg border border-neutral-200 overflow-hidden">
                  <iframe :src="getFullUrl(viewLesson.video_url)" class="w-full h-[500px]" frameborder="0"></iframe>
                </div>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-xs text-neutral-500 truncate flex-1">{{ viewLesson.video_url }}</span>
                  <a :href="getFullUrl(viewLesson.video_url)" target="_blank" class="ml-2 px-3 py-1 bg-neutral-100 hover:bg-neutral-200 rounded text-xs font-medium text-neutral-700 transition-colors">
                    Buka di Tab Baru ↗
                  </a>
                </div>
              </div>
              <div v-else class="text-neutral-500 text-center py-8">
                <svg class="w-12 h-12 mx-auto mb-2 text-neutral-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
                </svg>
                <p>Tidak ada dokumen</p>
              </div>
            </div>

            <!-- Text Content -->
            <div v-else-if="viewLesson?.content_type === 'text'" class="prose prose-neutral max-w-none">
              <div v-if="viewLesson?.content" class="p-4 bg-neutral-50 rounded-lg" v-html="viewLesson.content"></div>
              <p v-else class="text-neutral-500 text-center py-8">Tidak ada konten teks</p>
            </div>

            <!-- Quiz placeholder -->
            <div v-else-if="viewLesson?.content_type === 'quiz'" class="text-center py-8">
              <svg class="w-12 h-12 mx-auto mb-2 text-accent-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <p class="text-neutral-600">Klik tombol "Kelola Kuis" untuk mengelola pertanyaan</p>
            </div>
          </div>

          <!-- Description -->
          <div v-if="viewLesson?.description" class="mt-4 pt-4 border-t border-neutral-200">
            <h4 class="text-sm font-medium text-neutral-700 mb-2">Deskripsi</h4>
            <p class="text-neutral-600 text-sm">{{ viewLesson.description }}</p>
          </div>

          <button @click="showViewModal = false" class="w-full mt-6 py-2.5 text-sm font-medium text-neutral-600 bg-neutral-100 rounded-lg hover:bg-neutral-200 transition-colors">
            Tutup
          </button>
        </div>
      </div>
    </Transition>

    <!-- Quiz Modal -->
    <Transition name="fade">
      <div v-if="showQuizModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50" @click="showQuizModal = false">
        <div class="bg-white rounded-xl w-full max-w-3xl m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <!-- Header -->
          <div class="p-6 border-b border-neutral-200">
            <h2 class="text-xl font-bold text-neutral-900">{{ quiz ? 'Kelola Kuis' : 'Buat Kuis Baru' }}</h2>
          </div>
          
          <!-- Content -->
          <div class="p-6">
            <!-- Quiz Settings Form -->
            <div v-if="!quiz" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Judul Kuis *</label>
                <input v-model="quizForm.title" type="text" required class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Contoh: Kuis Bab 1" />
              </div>
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Deskripsi</label>
                <textarea v-model="quizForm.description" rows="2" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Deskripsi singkat kuis"></textarea>
              </div>
              <div class="grid grid-cols-3 gap-4">
                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Batas Waktu (menit)</label>
                  <input v-model.number="quizForm.time_limit" type="number" min="0" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="0 = tidak terbatas" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Nilai Kelulusan (%)</label>
                  <input v-model.number="quizForm.passing_score" type="number" min="0" max="100" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" />
                </div>
                <div>
                  <label class="block text-sm font-medium text-neutral-700 mb-1">Maks Percobaan</label>
                  <input v-model.number="quizForm.max_attempts" type="number" min="0" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="0 = unlimited" />
                </div>
              </div>
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="quizForm.show_correct_answers" type="checkbox" class="w-4 h-4 text-admin-500 rounded" />
                <span class="text-sm text-neutral-700">Tampilkan jawaban benar setelah submit</span>
              </label>
            </div>
            
            <!-- Quiz Questions List (when quiz exists) -->
            <div v-else class="space-y-4">
              <div class="flex items-center justify-between">
                <h3 class="font-semibold text-neutral-900">Pertanyaan ({{ questions.length }})</h3>
                <button @click="openQuestionForm" class="px-4 py-2 bg-admin-500 text-white rounded-lg hover:bg-admin-600 transition-colors text-sm font-medium">
                  + Tambah Pertanyaan
                </button>
              </div>
              
              <div v-if="questions.length === 0" class="p-8 text-center text-neutral-500 bg-neutral-50 rounded-lg">
                Belum ada pertanyaan. Klik tombol di atas untuk menambahkan.
              </div>
              
              <div v-else class="space-y-3">
                <div v-for="(q, index) in questions" :key="q.id" class="p-4 border border-neutral-200 rounded-lg hover:border-neutral-300 transition-colors">
                  <div class="flex items-start gap-3">
                    <div class="w-8 h-8 bg-neutral-100 rounded-full flex items-center justify-center text-sm font-medium text-neutral-600">
                      {{ index + 1 }}
                    </div>
                    <div class="flex-1">
                      <div class="flex items-center gap-2 mb-1">
                        <span class="px-2 py-0.5 rounded text-xs font-medium" :class="getQuestionTypeBadge(q.question_type)">
                          {{ getQuestionTypeLabel(q.question_type) }}
                        </span>
                        <span class="text-xs text-neutral-500">{{ q.points }} poin</span>
                      </div>
                      <p class="text-neutral-900">{{ q.question_text }}</p>
                      <div v-if="q.options && q.options.length > 0" class="mt-2 space-y-1">
                        <div v-for="opt in q.options" :key="opt.id" class="flex items-center gap-2 text-sm">
                          <span :class="opt.is_correct ? 'text-green-500' : 'text-neutral-400'">{{ opt.is_correct ? '✓' : '○' }}</span>
                          <span :class="opt.is_correct ? 'text-green-700 font-medium' : 'text-neutral-600'">{{ opt.option_text }}</span>
                        </div>
                      </div>
                    </div>
                    <div class="flex items-center gap-1">
                      <button @click="editQuestion(q)" class="p-2 hover:bg-neutral-100 rounded-lg text-neutral-600">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                      </button>
                      <button @click="handleDeleteQuestion(q)" class="p-2 hover:bg-red-50 rounded-lg text-red-500">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Footer -->
          <div class="p-6 border-t border-neutral-200 flex justify-end gap-3">
            <button @click="showQuizModal = false" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">Tutup</button>
            <button v-if="!quiz" @click="handleCreateQuiz" :disabled="saving" class="px-6 py-2 bg-admin-500 text-white rounded-lg hover:bg-admin-600 transition-colors disabled:opacity-50">
              {{ saving ? 'Menyimpan...' : 'Buat Kuis' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Question Modal -->
    <Transition name="fade">
      <div v-if="showQuestionModal" class="fixed inset-0 z-[60] flex items-center justify-center bg-black/50" @click="showQuestionModal = false">
        <div class="bg-white rounded-xl p-6 w-full max-w-2xl m-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <h2 class="text-xl font-bold text-neutral-900 mb-6">{{ editingQuestion ? 'Edit Pertanyaan' : 'Tambah Pertanyaan' }}</h2>
          <form @submit.prevent="handleSaveQuestion" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Tipe Pertanyaan *</label>
              <select v-model="questionForm.question_type" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent">
                <option value="multiple_choice">Pilihan Ganda (Satu Jawaban)</option>
                <option value="multiple_answer">Pilihan Ganda (Multi Jawaban)</option>
                <option value="true_false">Benar/Salah</option>
                <option value="short_answer">Jawaban Singkat</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Pertanyaan *</label>
              <textarea v-model="questionForm.question_text" rows="3" required class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Tulis pertanyaan di sini..."></textarea>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-neutral-700 mb-1">Poin</label>
                <input v-model.number="questionForm.points" type="number" min="1" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" />
              </div>
              <div class="flex items-end">
                <label class="flex items-center gap-2 cursor-pointer pb-2">
                  <input v-model="questionForm.required" type="checkbox" class="w-4 h-4 text-admin-500 rounded" />
                  <span class="text-sm text-neutral-700">Wajib dijawab</span>
                </label>
              </div>
            </div>

            <!-- Options for Multiple Choice / True False -->
            <div v-if="questionForm.question_type !== 'short_answer'">
              <label class="block text-sm font-medium text-neutral-700 mb-2">Pilihan Jawaban</label>
              <div class="space-y-2">
                <div v-for="(option, index) in questionForm.options" :key="index" class="flex items-center gap-2">
                  <input 
                    v-if="questionForm.question_type === 'multiple_choice' || questionForm.question_type === 'true_false'"
                    type="radio" 
                    :name="'correct-option'" 
                    :checked="option.is_correct"
                    @change="setCorrectOption(index)"
                    class="w-4 h-4 text-green-500"
                  />
                  <input 
                    v-else
                    type="checkbox" 
                    v-model="option.is_correct"
                    class="w-4 h-4 text-green-500 rounded"
                  />
                  <input v-model="option.option_text" type="text" class="flex-1 px-3 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Teks pilihan" />
                  <button v-if="questionForm.options.length > 2 && questionForm.question_type !== 'true_false'" type="button" @click="removeOptionField(index)" class="p-2 text-red-500 hover:bg-red-50 rounded-lg">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
                  </button>
                </div>
              </div>
              <button v-if="questionForm.question_type !== 'true_false'" type="button" @click="addOptionField" class="mt-2 text-sm text-admin-500 hover:text-admin-600">+ Tambah pilihan</button>
            </div>

            <!-- Short Answer -->
            <div v-else>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Jawaban Benar</label>
              <input v-model="questionForm.options[0].option_text" type="text" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Masukkan jawaban yang benar" />
              <p class="text-xs text-neutral-500 mt-1">Jawaban akan dicocokkan secara case-insensitive</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-neutral-700 mb-1">Penjelasan (opsional)</label>
              <textarea v-model="questionForm.explanation" rows="2" class="w-full px-4 py-2 border border-neutral-200 rounded-lg focus:ring-2 focus:ring-admin-500 focus:border-transparent" placeholder="Penjelasan yang ditampilkan setelah menjawab"></textarea>
            </div>

            <div class="flex justify-end gap-3 pt-4">
              <button type="button" @click="showQuestionModal = false" class="px-4 py-2 text-neutral-700 hover:bg-neutral-100 rounded-lg transition-colors">Batal</button>
              <button type="submit" :disabled="saving" class="px-6 py-2 bg-admin-500 text-white rounded-lg hover:bg-admin-600 transition-colors disabled:opacity-50">
                {{ saving ? 'Menyimpan...' : 'Simpan' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

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

const route = useRoute()
const courseId = computed(() => route.params.id as string)

useHead({
  title: 'Materi Kursus - Admin'
})

// API Composables
const { loading, error, lessons, lessonsTree, fetchLessons, fetchLessonsTree, createLesson, updateLesson, deleteLesson: deleteLessonApi } = useLessons()
const selectedParent = ref<any>(null)
const { course, fetchCourse } = useCourses()
const { 
  quiz, 
  questions, 
  loading: quizLoading, 
  getQuizByLesson, 
  createQuiz, 
  updateQuiz, 
  deleteQuiz, 
  addQuestion, 
  updateQuestion, 
  deleteQuestion 
} = useQuiz()

const showModal = ref(false)
const showDeleteModal = ref(false)
const showViewModal = ref(false)
const showQuizModal = ref(false)
const showQuestionModal = ref(false)
const isFullscreen = ref(false)
const isEditing = ref(false)
const selectedLesson = ref<any>(null)
const viewLesson = ref<any>(null)
const editingQuestion = ref<any>(null)
const selectedLessonForQuiz = ref<any>(null)
const saving = ref(false)
const toast = ref({ show: false, message: '', type: 'success' as 'success' | 'error' })

// Quiz form state
const quizForm = ref({
  title: '',
  description: '',
  time_limit: 0,
  passing_score: 70,
  max_attempts: 0,
  show_correct_answers: true
})

const questionForm = ref({
  question_type: 'multiple_choice' as string,
  question_text: '',
  explanation: '',
  points: 10,
  required: true,
  options: [
    { option_text: '', is_correct: true },
    { option_text: '', is_correct: false },
    { option_text: '', is_correct: false },
    { option_text: '', is_correct: false }
  ]
})

// Upload state
const uploadMode = ref<'url' | 'file'>('url')
const uploading = ref(false)
const uploadProgress = ref(0)
const fileInputRef = ref<HTMLInputElement | null>(null)

const contentTypes = [
  { value: 'video', label: 'Video', icon: 'M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664zM21 12a9 9 0 11-18 0 9 9 0 0118 0z' },
  { value: 'pdf', label: 'PDF', icon: 'M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z' },
  { value: 'text', label: 'Teks', icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' }
]

const form = ref({
  title: '',
  description: '',
  content_type: 'video',
  video_url: '',
  content: '', // Rich text content for text type
  is_preview: false,
  parent_id: null as string | null,
  is_container: false
})

// Text editor ref
const textEditorRef = ref<HTMLDivElement | null>(null)

// Load data on mount
onMounted(async () => {
  await Promise.all([
    fetchLessonsTree(courseId.value),
    fetchCourse(courseId.value)
  ])
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

// Quiz functions
const openQuizEditor = () => {
  // Reset state for new quiz
  quiz.value = null
  questions.value = []
  selectedLessonForQuiz.value = null
  
  quizForm.value = {
    title: '',
    description: '',
    time_limit: 0,
    passing_score: 70,
    max_attempts: 0,
    show_correct_answers: true
  }
  showQuizModal.value = true
}

const handleCreateQuiz = async () => {
  if (!quizForm.value.title) {
    showToast('Judul kuis wajib diisi', 'error')
    return
  }
  
  saving.value = true
  
  // First create a quiz-type lesson
  const lessonData = {
    title: quizForm.value.title,
    description: quizForm.value.description || 'Kuis',
    content_type: 'quiz',
    is_preview: false,
    parent_id: selectedParent.value?.id || null
  }
  
  const lesson = await createLesson(courseId.value, lessonData)
  
  if (lesson) {
    // Then create the quiz for this lesson
    const result = await createQuiz(lesson.id, quizForm.value)
    if (result) {
      showToast('Kuis berhasil dibuat! Tambahkan pertanyaan sekarang.')
      // Keep modal open and switch to questions view
      // quiz.value is already set by createQuiz in useQuiz composable
      await fetchLessonsTree(courseId.value)
      // Fetch the quiz to ensure we have the full data
      await getQuizByLesson(lesson.id)
    } else {
      showToast('Gagal membuat kuis', 'error')
    }
  } else {
    showToast('Gagal membuat lesson kuis', 'error')
  }
  
  saving.value = false
}

// Open quiz editor for existing quiz
const openQuizForLesson = async (lessonItem: any) => {
  if (lessonItem.content_type !== 'quiz') return
  
  // Reset quiz state
  quiz.value = null
  questions.value = []
  
  // Fetch quiz data
  await getQuizByLesson(lessonItem.id)
  
  if (quiz.value) {
    showQuizModal.value = true
  } else {
    // No quiz yet for this lesson, show create form
    quizForm.value = {
      title: lessonItem.title,
      description: lessonItem.description || '',
      time_limit: 0,
      passing_score: 70,
      max_attempts: 0,
      show_correct_answers: true
    }
    // Store the lesson id to create quiz for this existing lesson
    selectedLessonForQuiz.value = lessonItem
    showQuizModal.value = true
  }
}

const openQuestionForm = () => {
  editingQuestion.value = null
  questionForm.value = {
    question_type: 'multiple_choice',
    question_text: '',
    explanation: '',
    points: 10,
    required: true,
    options: [
      { option_text: '', is_correct: true },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false }
    ]
  }
  showQuestionModal.value = true
}

const handleSaveQuestion = async () => {
  if (!quiz.value) {
    showToast('Kuis belum dibuat, buat kuis terlebih dahulu', 'error')
    return
  }
  if (!questionForm.value.question_text) {
    showToast('Pertanyaan wajib diisi', 'error')
    return
  }
  
  saving.value = true
  
  let options = questionForm.value.options
  if (questionForm.value.question_type === 'true_false') {
    options = [
      { option_text: 'Benar', is_correct: questionForm.value.options[0]?.is_correct || false },
      { option_text: 'Salah', is_correct: !questionForm.value.options[0]?.is_correct }
    ]
  } else if (questionForm.value.question_type === 'short_answer') {
    options = [{ option_text: questionForm.value.options[0]?.option_text || '', is_correct: true }]
  }
  
  const data = {
    question_type: questionForm.value.question_type,
    question_text: questionForm.value.question_text,
    explanation: questionForm.value.explanation,
    points: questionForm.value.points,
    required: questionForm.value.required,
    options
  }
  
  let result
  if (editingQuestion.value) {
    result = await updateQuestion(editingQuestion.value.id, data)
  } else {
    result = await addQuestion(quiz.value.id, data)
  }
  
  saving.value = false
  
  if (result) {
    showQuestionModal.value = false
    showToast(editingQuestion.value ? 'Pertanyaan berhasil diperbarui' : 'Pertanyaan berhasil ditambahkan')
    // Refresh quiz data to get updated questions list
    if (quiz.value) {
      await getQuizByLesson(quiz.value.lesson_id)
    }
  } else {
    showToast('Gagal menyimpan pertanyaan', 'error')
  }
}

const editQuestion = (q: any) => {
  editingQuestion.value = q
  questionForm.value = {
    question_type: q.question_type,
    question_text: q.question_text,
    explanation: q.explanation || '',
    points: q.points,
    required: q.required,
    options: q.options?.map((o: any) => ({ option_text: o.option_text, is_correct: o.is_correct })) || []
  }
  showQuestionModal.value = true
}

const handleDeleteQuestion = async (q: any) => {
  if (!confirm('Hapus pertanyaan ini?')) return
  
  const result = await deleteQuestion(q.id)
  if (result) {
    showToast('Pertanyaan berhasil dihapus')
    // Refresh
    if (quiz.value) await getQuizByLesson(quiz.value.lesson_id)
  } else {
    showToast('Gagal menghapus pertanyaan', 'error')
  }
}

const setCorrectOption = (index: number) => {
  questionForm.value.options.forEach((opt, i) => {
    opt.is_correct = i === index
  })
}

const addOptionField = () => {
  questionForm.value.options.push({ option_text: '', is_correct: false })
}

const removeOptionField = (index: number) => {
  questionForm.value.options.splice(index, 1)
}

const getQuestionTypeBadge = (type: string) => {
  const badges: Record<string, string> = {
    'multiple_choice': 'bg-blue-100 text-blue-700',
    'multiple_answer': 'bg-purple-100 text-purple-700',
    'true_false': 'bg-green-100 text-green-700',
    'short_answer': 'bg-orange-100 text-orange-700'
  }
  return badges[type] || 'bg-neutral-100 text-neutral-700'
}

const getQuestionTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    'multiple_choice': 'Pilihan Ganda',
    'multiple_answer': 'Multi Jawaban',
    'true_false': 'Benar/Salah',
    'short_answer': 'Jawaban Singkat'
  }
  return labels[type] || type
}

// Watch for question type changes
watch(() => questionForm.value.question_type, (newType) => {
  if (newType === 'true_false') {
    questionForm.value.options = [
      { option_text: 'Benar', is_correct: true },
      { option_text: 'Salah', is_correct: false }
    ]
  } else if (newType === 'short_answer') {
    questionForm.value.options = [{ option_text: '', is_correct: true }]
  } else if (questionForm.value.options.length < 2) {
    questionForm.value.options = [
      { option_text: '', is_correct: true },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false },
      { option_text: '', is_correct: false }
    ]
  }
})

const openAddModal = (parent?: any) => {
  isEditing.value = false
  selectedLesson.value = null
  selectedParent.value = parent || null
  form.value = {
    title: '',
    description: '',
    content_type: 'video',
    video_url: '',
    content: '',
    is_preview: false,
    parent_id: parent?.id || null,
    is_container: false
  }
  showModal.value = true
}

const openAddModuleModal = () => {
  isEditing.value = false
  selectedLesson.value = null
  selectedParent.value = null
  form.value = {
    title: '',
    description: '',
    content_type: '',
    video_url: '',
    content: '',
    is_preview: false,
    parent_id: null,
    is_container: true
  }
  showModal.value = true
}

// Create submodule inside a parent module
const openAddSubmoduleModal = (parent: any) => {
  isEditing.value = false
  selectedLesson.value = null
  selectedParent.value = parent
  form.value = {
    title: '',
    description: '',
    content_type: '',
    video_url: '',
    content: '',
    is_preview: false,
    parent_id: parent.id,
    is_container: true
  }
  showModal.value = true
}

// Create quiz inside a parent module
const openQuizEditorWithParent = (parent: any) => {
  selectedParent.value = parent
  quiz.value = null
  questions.value = []
  selectedLessonForQuiz.value = null
  
  quizForm.value = {
    title: '',
    description: '',
    time_limit: 0,
    passing_score: 70,
    max_attempts: 0,
    show_correct_answers: true
  }
  showQuizModal.value = true
}


const openEditModal = (lesson: any) => {
  isEditing.value = true
  selectedLesson.value = lesson
  form.value = {
    title: lesson.title,
    description: lesson.description || '',
    content_type: lesson.content_type,
    video_url: lesson.video_url || '',
    content: lesson.content || '',
    is_preview: lesson.is_preview
  }
  // Populate text editor if content type is text
  if (lesson.content_type === 'text' && textEditorRef.value) {
    nextTick(() => {
      if (textEditorRef.value) {
        textEditorRef.value.innerHTML = lesson.content || ''
      }
    })
  }
  showModal.value = true
}

const openDeleteModal = (lesson: any) => {
  selectedLesson.value = lesson
  showDeleteModal.value = true
}

const closeModal = () => {
  showModal.value = false
  selectedLesson.value = null
}

const saveLesson = async () => {
  saving.value = true
  try {
    if (isEditing.value && selectedLesson.value) {
      const result = await updateLesson(selectedLesson.value.id, {
        title: form.value.title,
        description: form.value.description,
        content_type: form.value.content_type || undefined,
        video_url: form.value.video_url || undefined,
        content: form.value.content || undefined,
        is_preview: form.value.is_preview
      })
      if (result) {
        showToast('Materi berhasil diperbarui')
        await fetchLessonsTree(courseId.value)
      } else {
        showToast('Gagal memperbarui materi', 'error')
      }
    } else {
      const result = await createLesson(courseId.value, {
        title: form.value.title,
        description: form.value.description,
        content_type: form.value.content_type || undefined,
        video_url: form.value.video_url || undefined,
        content: form.value.content || undefined,
        is_preview: form.value.is_preview,
        parent_id: form.value.parent_id,
        is_container: form.value.is_container
      })
      if (result) {
        showToast(form.value.is_container ? 'Modul berhasil ditambahkan' : 'Materi berhasil ditambahkan')
        await fetchLessonsTree(courseId.value)
      } else {
        showToast('Gagal menambahkan materi', 'error')
      }
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    closeModal()
  }
}


const confirmDeleteLesson = async () => {
  if (!selectedLesson.value) return
  
  saving.value = true
  try {
    const success = await deleteLessonApi(selectedLesson.value.id)
    if (success) {
      showToast('Materi berhasil dihapus')
      await fetchLessonsTree(courseId.value)
    } else {
      showToast('Gagal menghapus materi', 'error')
    }
  } catch (err) {
    showToast('Terjadi kesalahan', 'error')
  } finally {
    saving.value = false
    showDeleteModal.value = false
    selectedLesson.value = null
  }
}

const getContentTypeClass = (type: string) => {
  const classes: Record<string, string> = {
    video: 'bg-primary-100 text-primary-600',
    pdf: 'bg-red-100 text-red-600',
    quiz: 'bg-accent-100 text-accent-600',
    text: 'bg-neutral-100 text-neutral-600'
  }
  return classes[type] || 'bg-neutral-100 text-neutral-600'
}

const getContentTypeBadgeClass = (type: string) => {
  const classes: Record<string, string> = {
    video: 'bg-primary-100 text-primary-700',
    pdf: 'bg-red-100 text-red-700',
    quiz: 'bg-accent-100 text-accent-700',
    text: 'bg-neutral-100 text-neutral-700'
  }
  return classes[type] || 'bg-neutral-100 text-neutral-600'
}

const getContentTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    video: 'Video',
    pdf: 'Dokumen',
    quiz: 'Kuis',
    text: 'Teks'
  }
  return labels[type] || type
}

// Upload functions
const triggerFileInput = () => {
  fileInputRef.value?.click()
}

const handleFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    await uploadFile(file)
  }
}

const handleFileDrop = async (event: DragEvent) => {
  const file = event.dataTransfer?.files?.[0]
  if (file) {
    await uploadFile(file)
  }
}

const uploadFile = async (file: File) => {
  uploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('file', file)

    // Get auth token and API URL
    const token = useCookie('token')
    const config = useRuntimeConfig()
    const apiUrl = config.public.apiBase || 'http://localhost:8080'
    
    const response = await fetch(`${apiUrl}/api/admin/upload`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token.value}`
      },
      body: formData
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Upload failed')
    }

    const result = await response.json()
    form.value.video_url = result.url
    uploadProgress.value = 100
    showToast('File berhasil diupload')
  } catch (err: any) {
    showToast(err.message || 'Gagal mengupload file', 'error')
  } finally {
    uploading.value = false
  }
}

// Reorder lessons
const moveLesson = async (index: number, direction: 'up' | 'down') => {
  const newIndex = direction === 'up' ? index - 1 : index + 1
  
  if (newIndex < 0 || newIndex >= lessons.value.length) return
  
  // Swap lessons in array
  const lessonsCopy = [...lessons.value]
  const temp = lessonsCopy[index]
  lessonsCopy[index] = lessonsCopy[newIndex]
  lessonsCopy[newIndex] = temp
  
  // Update order_index for both lessons
  const lesson1 = lessonsCopy[index]
  const lesson2 = lessonsCopy[newIndex]
  
  try {
    await Promise.all([
      updateLesson(lesson1.id, { order_index: index }),
      updateLesson(lesson2.id, { order_index: newIndex })
    ])
    
    // Refresh the list
    await fetchLessonsTree(courseId.value)
    showToast('Urutan materi berhasil diubah')
  } catch (err) {
    showToast('Gagal mengubah urutan', 'error')
  }
}

// Rich text editor functions
const formatText = (command: string, value?: string) => {
  if (command === 'formatBlock' && value) {
    document.execCommand(command, false, `<${value}>`)
  } else {
    document.execCommand(command, false)
  }
}

const onTextEditorInput = () => {
  if (textEditorRef.value) {
    form.value.content = textEditorRef.value.innerHTML
  }
}

const onTextEditorPaste = (event: ClipboardEvent) => {
  event.preventDefault()
  const text = event.clipboardData?.getData('text/plain') || ''
  document.execCommand('insertText', false, text)
}

// View modal
const openViewModal = (lesson: any) => {
  viewLesson.value = lesson
  isFullscreen.value = false
  showViewModal.value = true
}

const closeViewModal = () => {
  showViewModal.value = false
  isFullscreen.value = false
}

const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// URL helpers - Video platform detection
type VideoType = 'youtube' | 'vimeo' | 'gdrive' | 'zoom' | 'msstream' | 'sharepoint' | 'direct' | 'embed' | 'unknown'

const detectVideoType = (url: string | null | undefined): VideoType => {
  if (!url) return 'unknown'
  
  // YouTube
  if (url.includes('youtube.com') || url.includes('youtu.be')) return 'youtube'
  // Vimeo
  if (url.includes('vimeo.com')) return 'vimeo'
  // Google Drive
  if (url.includes('drive.google.com')) return 'gdrive'
  // Zoom Cloud Recording
  if (url.includes('zoom.us/rec')) return 'zoom'
  // Microsoft Stream
  if (url.includes('microsoftstream.com') || url.includes('web.microsoftstream.com')) return 'msstream'
  // SharePoint / OneDrive (Teams recordings)
  if (url.includes('sharepoint.com') || url.includes('onedrive.live.com') || url.includes('1drv.ms')) return 'sharepoint'
  // Direct video file
  if (/\.(mp4|webm|ogg|mov|m4v)(\?.*)?$/i.test(url)) return 'direct'
  // Generic embed
  if (url.includes('/embed') || url.includes('/player')) return 'embed'
  
  return 'direct'
}

const isEmbedVideo = (url: string | null | undefined): boolean => {
  const type = detectVideoType(url)
  return ['youtube', 'vimeo', 'gdrive', 'zoom', 'msstream', 'sharepoint', 'embed'].includes(type)
}

const getEmbedUrl = (url: string): string => {
  const type = detectVideoType(url)
  
  switch (type) {
    case 'youtube': {
      const watchMatch = url.match(/youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)/)
      if (watchMatch) return `https://www.youtube.com/embed/${watchMatch[1]}`
      const shortMatch = url.match(/youtu\.be\/([a-zA-Z0-9_-]+)/)
      if (shortMatch) return `https://www.youtube.com/embed/${shortMatch[1]}`
      if (url.includes('/embed/')) return url
      return url
    }
    case 'vimeo': {
      const vimeoMatch = url.match(/vimeo\.com\/(\d+)/)
      if (vimeoMatch) return `https://player.vimeo.com/video/${vimeoMatch[1]}`
      if (url.includes('player.vimeo.com')) return url
      return url
    }
    case 'gdrive': {
      const driveMatch = url.match(/drive\.google\.com\/file\/d\/([a-zA-Z0-9_-]+)/)
      if (driveMatch) return `https://drive.google.com/file/d/${driveMatch[1]}/preview`
      return url
    }
    case 'zoom': {
      if (url.includes('/rec/share/')) return url.replace('/rec/share/', '/rec/play/')
      return url
    }
    case 'msstream': {
      const streamMatch = url.match(/microsoftstream\.com\/video\/([a-zA-Z0-9-]+)/)
      if (streamMatch) return `https://web.microsoftstream.com/embed/video/${streamMatch[1]}`
      if (url.includes('/embed/')) return url
      return url
    }
    case 'sharepoint': {
      if (url.includes('sharepoint.com') && !url.includes('embed=1')) {
        const separator = url.includes('?') ? '&' : '?'
        return `${url}${separator}embed=1`
      }
      return url
    }
    default:
      return url
  }
}

const getVideoProviderName = (url: string | null | undefined): string => {
  const type = detectVideoType(url)
  const names: Record<string, string> = {
    youtube: 'YouTube', vimeo: 'Vimeo', gdrive: 'Google Drive',
    zoom: 'Zoom', msstream: 'MS Stream', sharepoint: 'Teams/SharePoint'
  }
  return names[type] || 'Video'
}

const getVideoProviderBadgeClass = (url: string | null | undefined): string => {
  const type = detectVideoType(url)
  const classes: Record<string, string> = {
    youtube: 'bg-red-100 text-red-700',
    vimeo: 'bg-blue-100 text-blue-700',
    gdrive: 'bg-green-100 text-green-700',
    zoom: 'bg-blue-100 text-blue-700',
    msstream: 'bg-purple-100 text-purple-700',
    sharepoint: 'bg-teal-100 text-teal-700'
  }
  return classes[type] || 'bg-neutral-100 text-neutral-700'
}

const getFullUrl = (url: string) => {
  if (!url) return ''
  // If it's already a full URL, return it
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  // Otherwise, prepend the API base URL
  const config = useRuntimeConfig()
  const apiUrl = config.public.apiBase || 'http://localhost:8080'
  return `${apiUrl}${url}`
}
</script>
