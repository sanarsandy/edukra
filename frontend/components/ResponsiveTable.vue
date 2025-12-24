<template>
  <div class="w-full">
    <!-- Desktop Table -->
    <div class="hidden md:block overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-neutral-200">
            <slot name="header" />
          </tr>
        </thead>
        <tbody>
          <slot name="body" />
        </tbody>
      </table>
    </div>
    
    <!-- Mobile Card View -->
    <div class="md:hidden space-y-4">
      <slot name="mobile" />
    </div>
  </div>
</template>

<script setup lang="ts">
// This component provides a responsive table wrapper
// On desktop: renders as table
// On mobile: renders slot "mobile" which should contain card-style views
</script>
