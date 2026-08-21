<script setup lang="ts">
import { computed } from 'vue'
import {
  CreditCard,
  LayoutDashboard,
  LogOut,
  Users,
  Wallet,
} from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()

const nav = [
  { name: 'Dashboard', path: '/app/dashboard', icon: LayoutDashboard },
  { name: 'Users', path: '/app/users', icon: Users },
  { name: 'Wallets', path: '/app/wallets', icon: Wallet },
  { name: 'Providers', path: '/app/providers', icon: CreditCard },
]

const drawerOpen = defineModel<boolean>('open', { default: false })

const initials = computed(() => {
  const n = auth.user?.name ?? '?'
  return n
    .split(' ')
    .map((p) => p[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
})

function isActive(path: string) {
  return router.currentRoute.value.path === path
}

function go(path: string) {
  drawerOpen.value = false
  router.push(path)
}

function logout() {
  auth.logout()
  router.replace('/')
}
</script>

<template>
  <!-- mobile backdrop -->
  <Transition
    enter-active-class="transition-opacity duration-200"
    enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-200"
    leave-to-class="opacity-0"
  >
    <div
      v-if="drawerOpen"
      class="fixed inset-0 z-40 bg-midnight-ink/20 backdrop-blur-sm lg:hidden"
      @click="drawerOpen = false"
    />
  </Transition>

  <!-- sidebar — Dub style: white canvas, hairline border -->
  <aside
    class="fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-ash bg-canvas-white transition-transform duration-300 lg:translate-x-0"
    :class="drawerOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex h-16 items-center justify-between border-b border-ash px-5">
      <div class="flex items-center gap-2.5">
        <div class="flex size-8 items-center justify-center rounded-lg bg-midnight-ink">
          <CreditCard class="size-4 text-white" />
        </div>
        <div>
          <div class="text-sm font-bold leading-tight tracking-tight text-midnight-ink">Payment Gateway</div>
          <div class="text-[10px] text-fog">payment.sawang.tech</div>
        </div>
      </div>
    </div>

    <div class="px-3 pt-4 text-[10px] font-semibold uppercase tracking-[0.08em] text-fog">Menu</div>

    <nav class="flex-1 space-y-0.5 px-3 py-2">
      <button
        v-for="item in nav"
        :key="item.path"
        type="button"
        class="flex w-full items-center gap-2.5 rounded-lg px-3 py-2 text-[13px] font-medium transition"
        :class="
          isActive(item.path)
            ? 'bg-paper-mist text-midnight-ink'
            : 'text-graphite hover:bg-paper-mist/60 hover:text-charcoal'
        "
        @click="go(item.path)"
      >
        <component :is="item.icon" class="size-4" :class="isActive(item.path) ? 'text-electric-blue' : 'text-steel'" />
        {{ item.name }}
      </button>
    </nav>

    <div class="border-t border-ash p-3">
      <div class="flex items-center gap-2.5 rounded-lg px-2 py-2">
        <div class="flex size-8 shrink-0 items-center justify-center rounded-full bg-electric-blue/10 text-xs font-semibold text-electric-blue">
          {{ initials }}
        </div>
        <div class="min-w-0 flex-1">
          <div class="truncate text-[13px] font-medium text-charcoal">{{ auth.user?.name }}</div>
          <div class="truncate text-[11px] text-fog">{{ auth.user?.email }}</div>
        </div>
      </div>
      <button
        type="button"
        class="mt-1 flex w-full items-center gap-2 rounded-lg px-3 py-2 text-[13px] font-medium text-steel transition hover:bg-paper-mist hover:text-charcoal"
        @click="logout"
      >
        <LogOut class="size-4" />
        Sign out
      </button>
    </div>
  </aside>
</template>
