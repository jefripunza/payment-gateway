<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { CreditCard, KeyRound, Loader2, Users, Wallet } from 'lucide-vue-next'
import { api } from '@/lib/api'
import type { Provider, User, Wallet as WalletType } from '@/lib/types'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

const users = ref<User[]>([])
const wallets = ref<WalletType[]>([])
const providers = ref<Provider[]>([])
const loading = ref(true)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [u, w, p] = await Promise.all([
      api.get('users').json<{ users: User[] }>(),
      api.get('wallets').json<{ wallets: WalletType[] }>(),
      api.get('providers').json<{ providers: Provider[] }>(),
    ])
    users.value = u.users
    wallets.value = w.wallets
    providers.value = p.providers
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load dashboard'
  } finally {
    loading.value = false
  }
}
onMounted(load)

const activeProviders = computed(() => providers.value.filter((p) => p.enabled))
const totalBalance = computed(() => wallets.value.reduce((acc, w) => acc + w.balance, 0))

function fmt(n: number, cur = 'IDR') {
  try {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: cur, maximumFractionDigits: 0 }).format(n)
  } catch {
    return `${cur} ${n.toLocaleString()}`
  }
}

function initials(name: string) {
  return name
    .split(' ')
    .map((p) => p[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()
}
</script>

<template>
  <div class="space-y-6">
    <!-- greeting -->
    <div class="rounded-2xl border border-white/10 bg-gradient-to-r from-indigo-500/10 via-violet-500/5 to-transparent p-6">
      <h2 class="text-xl font-bold tracking-tight text-slate-100">
        Welcome back, {{ auth.user?.name?.split(' ')[0] }} 👋
      </h2>
      <p class="mt-1 text-sm text-slate-400">
        Here's the state of your payment infrastructure.
      </p>
    </div>

    <div v-if="error" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-20 text-sm text-slate-500">
      <Loader2 class="size-4 animate-spin" />
      Loading dashboard…
    </div>

    <template v-else>
      <!-- stat cards -->
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5">
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <Users class="size-3.5" />
            Team members
          </div>
          <div class="mt-2 text-3xl font-extrabold tracking-tight text-slate-100">{{ users.length }}</div>
        </div>
        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5">
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <Wallet class="size-3.5" />
            Wallets
          </div>
          <div class="mt-2 text-3xl font-extrabold tracking-tight text-slate-100">{{ wallets.length }}</div>
        </div>
        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5">
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <CreditCard class="size-3.5" />
            Providers
          </div>
          <div class="mt-2 text-3xl font-extrabold tracking-tight text-slate-100">{{ providers.length }}</div>
          <div class="mt-1 text-xs text-emerald-400">{{ activeProviders.length }} active</div>
        </div>
        <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5">
          <div class="flex items-center gap-2 text-xs text-slate-500">
            <KeyRound class="size-3.5" />
            Combined balance
          </div>
          <div class="mt-2 truncate text-3xl font-extrabold tracking-tight text-emerald-400">
            {{ fmt(totalBalance) }}
          </div>
        </div>
      </div>

      <!-- two-column content -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <!-- recent providers -->
        <div class="rounded-2xl border border-white/10 bg-white/[0.03]">
          <div class="flex items-center justify-between border-b border-white/10 px-5 py-4">
            <h3 class="text-sm font-bold text-slate-100">Payment providers</h3>
            <RouterLink to="/app/providers" class="text-xs font-medium text-indigo-400 hover:text-indigo-300">
              View all →
            </RouterLink>
          </div>
          <div class="divide-y divide-white/5">
            <div v-for="p in providers.slice(0, 4)" :key="p.id" class="flex items-center gap-3 px-5 py-3.5">
              <div class="flex size-9 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500/25 to-violet-500/15">
                <CreditCard class="size-4.5 text-indigo-300" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="truncate text-sm font-medium text-slate-200">{{ p.name }}</div>
                <div class="text-xs text-slate-500">{{ p.type }}</div>
              </div>
              <span
                class="rounded-full px-2 py-0.5 text-[10px] font-semibold"
                :class="p.enabled ? 'bg-emerald-500/15 text-emerald-300' : 'bg-slate-500/15 text-slate-400'"
              >
                {{ p.enabled ? 'Enabled' : 'Disabled' }}
              </span>
            </div>
            <div v-if="providers.length === 0" class="px-5 py-8 text-center text-xs text-slate-600">
              No providers yet — add one from the Providers menu.
            </div>
          </div>
        </div>

        <!-- wallets summary -->
        <div class="rounded-2xl border border-white/10 bg-white/[0.03]">
          <div class="flex items-center justify-between border-b border-white/10 px-5 py-4">
            <h3 class="text-sm font-bold text-slate-100">Wallets</h3>
            <RouterLink to="/app/wallets" class="text-xs font-medium text-indigo-400 hover:text-indigo-300">
              View all →
            </RouterLink>
          </div>
          <div class="divide-y divide-white/5">
            <div v-for="w in wallets.slice(0, 4)" :key="w.id" class="flex items-center gap-3 px-5 py-3.5">
              <div class="flex size-9 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-500/10">
                <Wallet class="size-4.5 text-amber-400" />
              </div>
              <div class="min-w-0 flex-1">
                <div class="truncate text-sm font-medium text-slate-200">{{ w.name }}</div>
                <div class="text-xs text-slate-500">{{ w.currency }}</div>
              </div>
              <span class="text-sm font-semibold text-emerald-400">{{ fmt(w.balance, w.currency) }}</span>
            </div>
            <div v-if="wallets.length === 0" class="px-5 py-8 text-center text-xs text-slate-600">
              No wallets yet — add one from the Wallets menu.
            </div>
          </div>
        </div>
      </div>

      <!-- team -->
      <div class="rounded-2xl border border-white/10 bg-white/[0.03]">
        <div class="flex items-center justify-between border-b border-white/10 px-5 py-4">
          <h3 class="text-sm font-bold text-slate-100">Team</h3>
          <RouterLink to="/app/users" class="text-xs font-medium text-indigo-400 hover:text-indigo-300">
            Manage →
          </RouterLink>
        </div>
        <div class="grid grid-cols-1 divide-y divide-white/5 sm:grid-cols-2 sm:divide-x">
          <div v-for="u in users.slice(0, 4)" :key="u.id" class="flex items-center gap-3 px-5 py-3.5">
            <div class="flex size-9 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 text-[10px] font-bold text-white">
              {{ initials(u.name) }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="truncate text-sm font-medium text-slate-200">{{ u.name }}</div>
              <div class="truncate text-xs text-slate-500">{{ u.email }}</div>
            </div>
            <span
              class="rounded-full px-2 py-0.5 text-[10px] font-semibold"
              :class="u.role === 'admin' ? 'bg-indigo-500/15 text-indigo-300' : 'bg-slate-500/15 text-slate-300'"
            >
              {{ u.role }}
            </span>
          </div>
          <div v-if="users.length === 0" class="px-5 py-8 text-center text-xs text-slate-600">
            No users yet.
          </div>
        </div>
      </div>
    </template>
  </div>
</template>
