<template>
  <div class="animate-pulse">
    <!-- Stats Skeleton -->
    <div v-if="type === 'stats'" class="bg-white rounded-xl p-5 border border-neutral-200">
      <div class="flex items-center justify-between">
        <div>
          <div class="h-3 w-20 bg-neutral-200 rounded mb-2"></div>
          <div class="h-7 w-12 bg-neutral-200 rounded"></div>
        </div>
        <div class="w-11 h-11 bg-neutral-200 rounded-xl"></div>
      </div>
    </div>
    
    <!-- Card Skeleton -->
    <div v-else-if="type === 'card'" class="bg-white rounded-xl border border-neutral-200 overflow-hidden">
      <div class="h-36 bg-neutral-200"></div>
      <div class="p-5">
        <div class="h-5 w-3/4 bg-neutral-200 rounded mb-2"></div>
        <div class="h-4 w-1/2 bg-neutral-200 rounded mb-4"></div>
        <div class="h-3 bg-neutral-200 rounded mb-4"></div>
        <div class="h-10 bg-neutral-200 rounded"></div>
      </div>
    </div>
    
    <!-- Table Row Skeleton -->
    <div v-else-if="type === 'table-row'" class="flex items-center gap-4 p-4 border-b border-neutral-100">
      <div class="w-10 h-10 bg-neutral-200 rounded-full flex-shrink-0"></div>
      <div class="flex-1 space-y-2">
        <div class="h-4 w-1/3 bg-neutral-200 rounded"></div>
        <div class="h-3 w-1/4 bg-neutral-200 rounded"></div>
      </div>
      <div class="h-6 w-16 bg-neutral-200 rounded-full"></div>
      <div class="h-8 w-20 bg-neutral-200 rounded"></div>
    </div>
    
    <!-- Text Lines Skeleton -->
    <div v-else-if="type === 'text'" class="space-y-3">
      <div class="h-4 w-full bg-neutral-200 rounded"></div>
      <div class="h-4 w-5/6 bg-neutral-200 rounded"></div>
      <div class="h-4 w-4/6 bg-neutral-200 rounded"></div>
    </div>
    
    <!-- Avatar Skeleton -->
    <div v-else-if="type === 'avatar'" class="w-10 h-10 bg-neutral-200 rounded-full"></div>
    
    <!-- Button Skeleton -->
    <div v-else-if="type === 'button'" class="h-10 w-24 bg-neutral-200 rounded-lg"></div>
    
    <!-- Default Box Skeleton -->
    <div v-else class="h-20 bg-neutral-200 rounded-lg"></div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  type?: 'stats' | 'card' | 'table-row' | 'text' | 'avatar' | 'button' | 'box'
}>()
</script>
