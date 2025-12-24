<template>
  <div class="mb-4">
    <!-- Label -->
    <label v-if="label" :for="id" class="block text-sm font-medium text-neutral-700 mb-1.5">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <!-- Input Container -->
    <div class="relative">
      <input
        :id="id"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        class="w-full px-4 py-2.5 border rounded-lg transition-all focus:outline-none focus:ring-2"
        :class="[
          error 
            ? 'border-red-300 focus:ring-red-500 focus:border-red-500 bg-red-50' 
            : 'border-neutral-300 focus:ring-primary-500 focus:border-primary-500',
          disabled ? 'bg-neutral-100 cursor-not-allowed' : ''
        ]"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @blur="$emit('blur')"
      />
      
      <!-- Error Icon -->
      <div v-if="error" class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
        <svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
        </svg>
      </div>
    </div>
    
    <!-- Error Message -->
    <p v-if="error" class="mt-1.5 text-sm text-red-600">{{ error }}</p>
    
    <!-- Helper Text -->
    <p v-else-if="helper" class="mt-1.5 text-sm text-neutral-500">{{ helper }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  id?: string
  label?: string
  type?: string
  modelValue: string
  placeholder?: string
  error?: string
  helper?: string
  disabled?: boolean
  required?: boolean
}>()

defineEmits(['update:modelValue', 'blur'])

withDefaults(defineProps<{
  type?: string
}>(), {
  type: 'text'
})
</script>
