<template>
  <nav class="flex items-center space-x-2 text-sm">
    <NuxtLink to="/" class="text-neutral-500 hover:text-neutral-700 transition-colors">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
      </svg>
    </NuxtLink>
    
    <template v-for="(item, index) in items" :key="index">
      <svg class="w-4 h-4 text-neutral-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5l7 7-7 7"/>
      </svg>
      
      <NuxtLink 
        v-if="index < items.length - 1" 
        :to="item.path" 
        class="text-neutral-500 hover:text-neutral-700 transition-colors"
      >
        {{ item.label }}
      </NuxtLink>
      <span v-else class="text-neutral-900 font-medium">{{ item.label }}</span>
    </template>
  </nav>
</template>

<script setup lang="ts">
interface BreadcrumbItem {
  label: string
  path?: string
}

defineProps<{
  items: BreadcrumbItem[]
}>()
</script>
