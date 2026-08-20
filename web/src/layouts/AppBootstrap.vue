<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/layouts/AppLayout.vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()

onMounted(async () => {
  await auth.init()
  if (!auth.user) {
    router.replace('/')
  }
})
</script>

<template>
  <AppLayout v-if="auth.user" />
  <div v-else class="flex min-h-screen items-center justify-center bg-slate-950 text-sm text-slate-500">
    Checking session…
  </div>
</template>
