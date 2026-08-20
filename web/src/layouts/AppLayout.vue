<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppSidebar from '@/components/AppSidebar.vue'
import AppHeader from '@/components/AppHeader.vue'
import ChangePasswordDialog from '@/components/ChangePasswordDialog.vue'
import { KeyRound } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const route = useRoute()
const drawerOpen = ref(false)
const passwordOpen = ref(false)

const titles: Record<string, { title: string; subtitle: string }> = {
  dashboard: { title: 'Dashboard', subtitle: 'Payment infrastructure overview' },
  users: { title: 'Users', subtitle: 'Manage team access to this dashboard' },
  wallets: { title: 'Wallets', subtitle: 'Track balances across currencies' },
  providers: { title: 'Providers', subtitle: 'Payment gateway credentials' },
}

const current = ref<{ title: string; subtitle: string }>(titles.dashboard!)
watch(
  () => route.name,
  (name) => {
    const key = (name as string) ?? 'dashboard'
    current.value = titles[key] ?? titles.dashboard!
  },
  { immediate: true },
)
</script>

<template>
  <div class="min-h-screen bg-slate-950 font-sans text-slate-100 antialiased">
    <div class="pointer-events-none fixed inset-0">
      <div class="absolute top-0 left-1/4 h-96 w-96 rounded-full bg-indigo-600/10 blur-[120px]" />
      <div class="absolute right-0 bottom-0 h-96 w-96 rounded-full bg-violet-600/10 blur-[120px]" />
    </div>

    <AppSidebar v-model:open="drawerOpen" />

    <div class="relative lg:pl-64">
      <AppHeader
        :title="current.title"
        :subtitle="current.subtitle"
        v-model:open="drawerOpen"
      >
        <template #actions>
          <button
            type="button"
            class="flex h-9 items-center gap-2 rounded-xl border border-white/10 bg-white/5 px-3 text-xs font-medium text-slate-300 transition hover:bg-white/10 hover:text-slate-100"
            @click="passwordOpen = true"
          >
            <KeyRound class="size-3.5" />
            <span class="hidden sm:inline">Change password</span>
          </button>
        </template>
      </AppHeader>

      <main class="mx-auto max-w-7xl px-4 py-6 sm:px-6">
        <RouterView />
      </main>
    </div>

    <ChangePasswordDialog v-model:open="passwordOpen" />
  </div>
</template>
