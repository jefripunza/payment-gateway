<script setup lang="ts">
import { computed } from 'vue'
import {
  CreditCard,
  LayoutDashboard,
  LogOut,
  Menu,
  Users,
  Wallet,
  X,
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
      class="fixed inset-0 z-40 bg-slate-950/70 backdrop-blur-sm lg:hidden"
      @click="drawerOpen = false"
    />
  </Transition>

  <!-- sidebar -->
  <aside
    class="fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-white/5 bg-slate-900/80 backdrop-blur-xl transition-transform duration-300 lg:translate-x-0"
    :class="drawerOpen ? 'translate-x-0' : '-translate-x-full'"
  >
    <div class="flex h-16 items-center justify-between border-b border-white/5 px-5">
      <div class="flex items-center gap-2.5">
        <div class="flex size-9 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-violet-600">
          <CreditCard class="size-5 text-white" />
        </div>
        <div>
          <div class="text-sm font-bold leading-tight tracking-tight">Payment Gateway</div>
          <div class="text-[10px] text-slate-500">payment.sawang.tech</div>
        </div>
      </div>
      <button
        type="button"
        class="rounded-lg p-1.5 text-slate-400 hover:bg-white/5 hover:text-slate-200 lg:hidden"
        @click="drawerOpen = false"
      >
        <X class="size-5" />
      </button>
    </div>

    <nav class="flex-1 space-y-1 overflow-y-auto px-3 py-4">
      <div class="px-3 pb-2 text-[10px] font-semibold tracking-widest text-slate-600 uppercase">
        Menu
      </div>
      <button
        v-for="item in nav"
        :key="item.path"
        type="button"
        class="flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition"
        :class="
          isActive(item.path)
            ? 'bg-gradient-to-r from-indigo-500/20 to-violet-500/10 text-white ring-1 ring-indigo-400/30'
            : 'text-slate-400 hover:bg-white/5 hover:text-slate-200'
        "
        @click="go(item.path)"
      >
        <component :is="item.icon" class="size-4.5" />
        {{ item.name }}
      </button>
    </nav>

    <div class="border-t border-white/5 p-4">
      <div class="flex items-center gap-3">
        <div class="flex size-9 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 text-xs font-bold text-white">
          {{ initials }}
        </div>
        <div class="min-w-0 flex-1">
          <div class="truncate text-sm font-semibold text-slate-200">{{ auth.user?.name }}</div>
          <div class="truncate text-xs text-slate-500">{{ auth.user?.email }}</div>
        </div>
        <button
          type="button"
          title="Sign out"
          class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400"
          @click="logout"
        >
          <LogOut class="size-4.5" />
        </button>
      </div>
    </div>
  </aside>
</template>
