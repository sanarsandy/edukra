<template>
  <div class="blog-editor">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 p-3 bg-neutral-50 border border-neutral-200 rounded-t-lg">
      <!-- Text Formatting -->
      <div class="flex items-center gap-0.5 pr-2 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBold().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('bold') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Bold (Ctrl+B)"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleItalic().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('italic') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Italic (Ctrl+I)"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 4h4m-2 0v16m4-16h-4"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleUnderline().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('underline') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Underline (Ctrl+U)"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4v6a5 5 0 0010 0V4M5 20h14"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleStrike().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('strike') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Strikethrough"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 12H7m10-5a4 4 0 00-8 0v10a4 4 0 008 0"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleSubscript().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('subscript') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors text-xs font-bold"
          title="Subscript"
        >
          X₂
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleSuperscript().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('superscript') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors text-xs font-bold"
          title="Superscript"
        >
          X²
        </button>
      </div>

      <!-- Text Color & Highlight -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <div class="relative">
          <button 
            type="button"
            @click="showColorPicker = !showColorPicker"
            class="p-1.5 rounded hover:bg-neutral-200 transition-colors flex items-center gap-1"
            title="Text Color"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v18m-4-8l4-10 4 10"/>
            </svg>
            <div class="w-3 h-1 rounded" :style="{ backgroundColor: currentColor }"></div>
          </button>
          <div v-if="showColorPicker" class="absolute top-full left-0 mt-1 p-2 bg-white rounded-lg shadow-lg border z-20 grid grid-cols-6 gap-1">
            <button 
              v-for="color in textColors" 
              :key="color"
              @click="setColor(color)" 
              class="w-5 h-5 rounded border border-neutral-200 hover:scale-110 transition-transform"
              :style="{ backgroundColor: color }"
            ></button>
          </div>
        </div>
        <div class="relative">
          <button 
            type="button"
            @click="showHighlightPicker = !showHighlightPicker"
            :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('highlight') }"
            class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
            title="Highlight"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
              <path d="M4 19h16v2H4v-2zM19.4 4.6l-2.8-2.8c-.8-.8-2-.8-2.8 0l-10 10c-.4.4-.6.9-.6 1.4V17h3.8c.5 0 1-.2 1.4-.6l10-10c.8-.8.8-2 0-2.8zM7.6 15H6v-1.6l6-6 1.6 1.6-6 6z"/>
            </svg>
          </button>
          <div v-if="showHighlightPicker" class="absolute top-full left-0 mt-1 p-2 bg-white rounded-lg shadow-lg border z-20 grid grid-cols-5 gap-1">
            <button 
              v-for="color in highlightColors" 
              :key="color"
              @click="setHighlight(color)" 
              class="w-5 h-5 rounded border border-neutral-200 hover:scale-110 transition-transform"
              :style="{ backgroundColor: color }"
            ></button>
            <button 
              @click="removeHighlight()" 
              class="w-5 h-5 rounded border border-neutral-200 hover:scale-110 transition-transform flex items-center justify-center text-red-500"
              title="Remove highlight"
            >
              ✕
            </button>
          </div>
        </div>
      </div>

      <!-- Headings & Font Size -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <select 
          @change="setHeading($event.target.value)"
          class="px-2 py-1 text-xs border border-neutral-200 rounded hover:bg-neutral-100"
        >
          <option value="p">Paragraph</option>
          <option value="1">Heading 1</option>
          <option value="2">Heading 2</option>
          <option value="3">Heading 3</option>
        </select>
      </div>

      <!-- Text Align -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().setTextAlign('left').run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive({ textAlign: 'left' }) }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Align Left"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h10M4 18h14"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().setTextAlign('center').run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive({ textAlign: 'center' }) }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Align Center"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M7 12h10M5 18h14"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().setTextAlign('right').run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive({ textAlign: 'right' }) }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Align Right"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M10 12h10M6 18h14"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().setTextAlign('justify').run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive({ textAlign: 'justify' }) }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Justify"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
          </svg>
        </button>
      </div>

      <!-- Lists -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBulletList().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('bulletList') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Bullet List"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 6h13M8 12h13M8 18h13M3 6h.01M3 12h.01M3 18h.01"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleOrderedList().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('orderedList') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Numbered List"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 6h13M8 12h13M8 18h13"/>
            <text x="1" y="8" font-size="6" fill="currentColor">1</text>
            <text x="1" y="14" font-size="6" fill="currentColor">2</text>
            <text x="1" y="20" font-size="6" fill="currentColor">3</text>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleTaskList().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('taskList') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Task List"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"/>
          </svg>
        </button>
      </div>

      <!-- Block Elements -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBlockquote().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('blockquote') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Quote"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('codeBlock') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Code Block"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().setHorizontalRule().run()"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Horizontal Line"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 12h16"/>
          </svg>
        </button>
      </div>

      <!-- Table -->
      <div class="flex items-center gap-0.5 px-2 border-r border-neutral-200">
        <div class="relative">
          <button 
            type="button"
            @click="showTableMenu = !showTableMenu"
            class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
            title="Table"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18M10 3v18M14 3v18M3 6a3 3 0 013-3h12a3 3 0 013 3v12a3 3 0 01-3 3H6a3 3 0 01-3-3V6z"/>
            </svg>
          </button>
          <div v-if="showTableMenu" class="absolute top-full left-0 mt-1 p-2 bg-white rounded-lg shadow-lg border z-20 min-w-[140px]">
            <button @click="insertTable" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded">Insert Table</button>
            <hr class="my-1">
            <button @click="editor?.chain().focus().addColumnBefore().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded">Add Col Before</button>
            <button @click="editor?.chain().focus().addColumnAfter().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded">Add Col After</button>
            <button @click="editor?.chain().focus().deleteColumn().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded text-red-600">Delete Column</button>
            <hr class="my-1">
            <button @click="editor?.chain().focus().addRowBefore().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded">Add Row Before</button>
            <button @click="editor?.chain().focus().addRowAfter().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded">Add Row After</button>
            <button @click="editor?.chain().focus().deleteRow().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded text-red-600">Delete Row</button>
            <hr class="my-1">
            <button @click="editor?.chain().focus().deleteTable().run()" class="w-full text-left px-2 py-1 text-sm hover:bg-neutral-100 rounded text-red-600">Delete Table</button>
          </div>
        </div>
      </div>

      <!-- Media -->
      <div class="flex items-center gap-0.5 px-2">
        <button 
          type="button"
          @click="setLink"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('link') }"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Add Link"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
          </svg>
        </button>
        <label class="p-1.5 rounded hover:bg-neutral-200 transition-colors cursor-pointer" title="Upload Image">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
          <input type="file" accept="image/*" @change="uploadImage" class="hidden" :disabled="uploadingImage" />
        </label>
        <button 
          type="button"
          @click="addYoutube"
          class="p-1.5 rounded hover:bg-neutral-200 transition-colors"
          title="Embed YouTube"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M19.615 3.184c-3.604-.246-11.631-.245-15.23 0-3.897.266-4.356 2.62-4.385 8.816.029 6.185.484 8.549 4.385 8.816 3.6.245 11.626.246 15.23 0 3.897-.266 4.356-2.62 4.385-8.816-.029-6.185-.484-8.549-4.385-8.816zm-10.615 12.816v-8l8 3.993-8 4.007z"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Editor Content -->
    <div class="border border-t-0 border-neutral-200 rounded-b-lg bg-white relative">
      <div v-if="uploadingImage" class="absolute inset-0 bg-white/80 flex items-center justify-center z-10">
        <div class="flex items-center gap-2">
          <div class="animate-spin rounded-full h-5 w-5 border-b-2 border-admin-600"></div>
          <span class="text-sm text-neutral-600">Uploading image...</span>
        </div>
      </div>
      <EditorContent :editor="editor" class="prose prose-neutral max-w-none min-h-[400px] p-6" />
    </div>

    <!-- Character Count -->
    <div class="flex items-center justify-between px-3 py-2 bg-neutral-50 border border-t-0 border-neutral-200 rounded-b-lg text-xs text-neutral-500">
      <div>
        {{ characterCount }} karakter · {{ wordCount }} kata
      </div>
      <div v-if="uploadingImage" class="text-admin-600">
        Uploading...
      </div>
    </div>
  </div>
</template>

<script setup>
import { useEditor, EditorContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import Link from '@tiptap/extension-link'
import Image from '@tiptap/extension-image'
import Underline from '@tiptap/extension-underline'
import Placeholder from '@tiptap/extension-placeholder'
import TextAlign from '@tiptap/extension-text-align'
import Highlight from '@tiptap/extension-highlight'
import Color from '@tiptap/extension-color'
import TextStyle from '@tiptap/extension-text-style'
import Subscript from '@tiptap/extension-subscript'
import Superscript from '@tiptap/extension-superscript'
import Table from '@tiptap/extension-table'
import TableRow from '@tiptap/extension-table-row'
import TableHeader from '@tiptap/extension-table-header'
import TableCell from '@tiptap/extension-table-cell'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import CharacterCount from '@tiptap/extension-character-count'
import Youtube from '@tiptap/extension-youtube'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Tulis konten artikel Anda di sini...'
  }
})

const emit = defineEmits(['update:modelValue'])
const config = useRuntimeConfig()
const apiBase = config.public.apiBase

// UI States
const showColorPicker = ref(false)
const showHighlightPicker = ref(false)
const showTableMenu = ref(false)
const uploadingImage = ref(false)
const currentColor = ref('#000000')

// Color palettes
const textColors = [
  '#000000', '#374151', '#6B7280', '#DC2626', '#EA580C', '#D97706',
  '#16A34A', '#059669', '#0891B2', '#2563EB', '#7C3AED', '#DB2777'
]
const highlightColors = [
  '#FEF08A', '#FDE68A', '#BBF7D0', '#A5F3FC', '#C7D2FE', '#FBCFE8'
]

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] }
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: { class: 'text-primary-600 underline hover:text-primary-700' }
    }),
    Image.configure({
      HTMLAttributes: { class: 'rounded-lg max-w-full mx-auto my-4' }
    }),
    Underline,
    Placeholder.configure({ placeholder: props.placeholder }),
    TextAlign.configure({ types: ['heading', 'paragraph'] }),
    Highlight.configure({ multicolor: true }),
    TextStyle,
    Color,
    Subscript,
    Superscript,
    Table.configure({ resizable: true }),
    TableRow,
    TableHeader,
    TableCell,
    TaskList,
    TaskItem.configure({ nested: true }),
    CharacterCount,
    Youtube.configure({
      width: 640,
      height: 360,
      HTMLAttributes: { class: 'rounded-lg overflow-hidden my-4' }
    })
  ],
  editorProps: {
    attributes: { class: 'focus:outline-none' }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

// Computed
const characterCount = computed(() => editor.value?.storage.characterCount.characters() || 0)
const wordCount = computed(() => editor.value?.storage.characterCount.words() || 0)

// Watch for external changes
watch(() => props.modelValue, (newVal) => {
  if (editor.value && newVal !== editor.value.getHTML()) {
    editor.value.commands.setContent(newVal, false)
  }
})

// Close dropdowns when clicking outside
const closeDropdowns = () => {
  showColorPicker.value = false
  showHighlightPicker.value = false
  showTableMenu.value = false
}

// Methods
const setHeading = (level) => {
  if (level === 'p') {
    editor.value?.chain().focus().setParagraph().run()
  } else {
    editor.value?.chain().focus().toggleHeading({ level: parseInt(level) }).run()
  }
}

const setColor = (color) => {
  currentColor.value = color
  editor.value?.chain().focus().setColor(color).run()
  showColorPicker.value = false
}

const setHighlight = (color) => {
  editor.value?.chain().focus().toggleHighlight({ color }).run()
  showHighlightPicker.value = false
}

const removeHighlight = () => {
  editor.value?.chain().focus().unsetHighlight().run()
  showHighlightPicker.value = false
}

const setLink = () => {
  const previousUrl = editor.value?.getAttributes('link').href
  const url = window.prompt('Masukkan URL:', previousUrl)
  
  if (url === null) return
  
  if (url === '') {
    editor.value?.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }
  
  editor.value?.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}

const uploadImage = async (event) => {
  const file = event.target?.files?.[0]
  if (!file) return
  
  uploadingImage.value = true
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const token = useCookie('token')
    const response = await $fetch(`${apiBase}/api/admin/upload`, {
      method: 'POST',
      body: formData,
      headers: { 'Authorization': `Bearer ${token.value}` }
    })
    
    if (response?.url || response?.object_key) {
      const imageUrl = `${apiBase}/api/images/${encodeURIComponent(response.object_key || response.url)}`
      editor.value?.chain().focus().setImage({ src: imageUrl }).run()
    }
  } catch (error) {
    console.error('Upload error:', error)
    alert('Gagal mengupload gambar')
  } finally {
    uploadingImage.value = false
    event.target.value = ''
  }
}

const addYoutube = () => {
  const url = window.prompt('Masukkan URL YouTube:')
  if (url) {
    editor.value?.chain().focus().setYoutubeVideo({ src: url }).run()
  }
}

const insertTable = () => {
  editor.value?.chain().focus().insertTable({ rows: 3, cols: 3, withHeaderRow: true }).run()
  showTableMenu.value = false
}

onMounted(() => {
  document.addEventListener('click', closeDropdowns)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', closeDropdowns)
  editor.value?.destroy()
})
</script>

<style scoped>
:deep(.ProseMirror) {
  min-height: 400px;
}

:deep(.ProseMirror p.is-editor-empty:first-child::before) {
  content: attr(data-placeholder);
  color: #9ca3af;
  pointer-events: none;
  float: left;
  height: 0;
}

:deep(.ProseMirror h1) {
  @apply text-3xl font-bold text-neutral-900 mt-6 mb-4;
}

:deep(.ProseMirror h2) {
  @apply text-2xl font-bold text-neutral-900 mt-5 mb-3;
}

:deep(.ProseMirror h3) {
  @apply text-xl font-semibold text-neutral-900 mt-4 mb-2;
}

:deep(.ProseMirror p) {
  @apply text-neutral-700 leading-relaxed mb-4;
}

:deep(.ProseMirror ul) {
  @apply list-disc list-inside mb-4 space-y-1;
}

:deep(.ProseMirror ol) {
  @apply list-decimal list-inside mb-4 space-y-1;
}

:deep(.ProseMirror blockquote) {
  @apply border-l-4 border-primary-500 pl-4 py-2 my-4 text-neutral-600 italic bg-neutral-50 rounded-r;
}

:deep(.ProseMirror pre) {
  @apply bg-neutral-900 text-neutral-100 p-4 rounded-lg my-4 overflow-x-auto font-mono text-sm;
}

:deep(.ProseMirror code) {
  @apply bg-neutral-100 text-admin-600 px-1.5 py-0.5 rounded text-sm font-mono;
}

:deep(.ProseMirror hr) {
  @apply my-8 border-neutral-200;
}

:deep(.ProseMirror img) {
  @apply rounded-lg max-w-full mx-auto my-4;
}

/* Table styles */
:deep(.ProseMirror table) {
  @apply border-collapse w-full my-4;
}

:deep(.ProseMirror th) {
  @apply border border-neutral-300 bg-neutral-100 p-2 font-semibold text-left;
}

:deep(.ProseMirror td) {
  @apply border border-neutral-300 p-2;
}

/* Task list styles */
:deep(.ProseMirror ul[data-type="taskList"]) {
  @apply list-none p-0;
}

:deep(.ProseMirror ul[data-type="taskList"] li) {
  @apply flex items-start gap-2;
}

:deep(.ProseMirror ul[data-type="taskList"] li input) {
  @apply mt-1;
}

/* YouTube embed */
:deep(.ProseMirror iframe) {
  @apply w-full aspect-video rounded-lg my-4;
}
</style>
