<template>
  <!-- Edit page reuses create.vue logic via route.params.id -->
  <CampaignCreate />
</template>

<script setup lang="ts">
import CampaignCreate from '../create.vue'

definePageMeta({
  layout: 'admin',
  middleware: ['admin']
})

useHead({ title: 'Edit Campaign - Admin' })
</script>
