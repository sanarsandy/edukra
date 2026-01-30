<template>
  <div class="blog-editor">
    <!-- Toolbar -->
    <div class="flex flex-wrap items-center gap-1 p-3 bg-neutral-50 border border-neutral-200 rounded-t-lg">
      <!-- Text Formatting -->
      <div class="flex items-center gap-0.5 pr-3 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBold().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('bold') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
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
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
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
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
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
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Strikethrough"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 12H7m10-5a4 4 0 00-8 0v10a4 4 0 008 0"/>
          </svg>
        </button>
      </div>

      <!-- Headings -->
      <div class="flex items-center gap-0.5 px-3 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleHeading({ level: 1 }).run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('heading', { level: 1 }) }"
          class="px-2 py-1 rounded hover:bg-neutral-200 transition-colors text-sm font-bold"
          title="Heading 1"
        >
          H1
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleHeading({ level: 2 }).run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('heading', { level: 2 }) }"
          class="px-2 py-1 rounded hover:bg-neutral-200 transition-colors text-sm font-bold"
          title="Heading 2"
        >
          H2
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleHeading({ level: 3 }).run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('heading', { level: 3 }) }"
          class="px-2 py-1 rounded hover:bg-neutral-200 transition-colors text-sm font-bold"
          title="Heading 3"
        >
          H3
        </button>
      </div>

      <!-- Lists -->
      <div class="flex items-center gap-0.5 px-3 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBulletList().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('bulletList') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Bullet List"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
            <circle cx="2" cy="6" r="1" fill="currentColor"/>
            <circle cx="2" cy="12" r="1" fill="currentColor"/>
            <circle cx="2" cy="18" r="1" fill="currentColor"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleOrderedList().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('orderedList') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Numbered List"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 6h13M8 12h13M8 18h13"/>
            <text x="1" y="8" font-size="6" fill="currentColor">1</text>
            <text x="1" y="14" font-size="6" fill="currentColor">2</text>
            <text x="1" y="20" font-size="6" fill="currentColor">3</text>
          </svg>
        </button>
      </div>

      <!-- Block Elements -->
      <div class="flex items-center gap-0.5 px-3 border-r border-neutral-200">
        <button 
          type="button"
          @click="editor?.chain().focus().toggleBlockquote().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('blockquote') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Quote"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10.5c0-1.657 1.343-3 3-3s3 1.343 3 3c0 2.5-3 4.5-3 4.5m0-9c0-1.657 1.343-3 3-3s3 1.343 3 3c0 2.5-3 4.5-3 4.5"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().toggleCodeBlock().run()"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('codeBlock') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Code Block"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="editor?.chain().focus().setHorizontalRule().run()"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Horizontal Line"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 12h16"/>
          </svg>
        </button>
      </div>

      <!-- Link & Image -->
      <div class="flex items-center gap-0.5 pl-2">
        <button 
          type="button"
          @click="setLink"
          :class="{ 'bg-admin-100 text-admin-700': editor?.isActive('link') }"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Add Link"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1"/>
          </svg>
        </button>
        <button 
          type="button"
          @click="addImage"
          class="p-2 rounded hover:bg-neutral-200 transition-colors"
          title="Add Image"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Editor Content -->
    <div class="border border-t-0 border-neutral-200 rounded-b-lg bg-white">
      <EditorContent :editor="editor" class="prose prose-neutral max-w-none min-h-[400px] p-6" />
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

const editor = useEditor({
  content: props.modelValue,
  extensions: [
    StarterKit.configure({
      heading: {
        levels: [1, 2, 3]
      }
    }),
    Link.configure({
      openOnClick: false,
      HTMLAttributes: {
        class: 'text-primary-600 underline hover:text-primary-700'
      }
    }),
    Image.configure({
      HTMLAttributes: {
        class: 'rounded-lg max-w-full mx-auto my-4'
      }
    }),
    Underline,
    Placeholder.configure({
      placeholder: props.placeholder
    })
  ],
  editorProps: {
    attributes: {
      class: 'focus:outline-none'
    }
  },
  onUpdate: ({ editor }) => {
    emit('update:modelValue', editor.getHTML())
  }
})

// Watch for external changes
watch(() => props.modelValue, (newVal) => {
  if (editor.value && newVal !== editor.value.getHTML()) {
    editor.value.commands.setContent(newVal, false)
  }
})

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

const addImage = () => {
  const url = window.prompt('Masukkan URL gambar:')
  
  if (url) {
    editor.value?.chain().focus().setImage({ src: url }).run()
  }
}

onBeforeUnmount(() => {
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
</style>
