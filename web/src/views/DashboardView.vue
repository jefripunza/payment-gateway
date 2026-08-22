<script setup lang="ts">
import { computed } from 'vue'
import { CreditCard, KeyRound, Users, Wallet } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'
import { useProvidersStore, useUsersStore, useWalletsStore } from '@/stores/resources'

const auth = useAuthStore()
const users = useUsersStore()
const wallets = useWalletsStore()
const providers = useProvidersStore()

// Pinia Colada — reactive queries, no manual fetch in components
users.query
wallets.query
providers.query

const userCount = computed(() => users.query.data?.length ?? 0)
const walletCount = computed(() => wallets.query.data?.length ?? 0)
const providerCount = computed(() => providers.query.data?.length ?? 0)
const activeCount = computed(() => providers.query.data?.filter((p) => p.enabled).length ?? 0)
const combinedBalance = computed(() => wallets.query.data?.reduce((acc, w) => acc + (w.balance ?? 0), 0) ?? 0)

const fmtIDR = new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 })

const stats = computed(() => [
  { label: 'Team members', value: String(userCount.value), sub: `${userCount.value} active`, icon: Users, color: 'blue' },
  { label: 'Wallets', value: String(walletCount.value), sub: 'across currencies', icon: Wallet, color: 'violet' },
  { label: 'Providers', value: String(providerCount.value), sub: `${activeCount.value} active`, icon: CreditCard, color: 'green' },
  { label: 'Combined balance', value: fmtIDR.format(combinedBalance.value), sub: 'all wallets', icon: KeyRound, color: 'orange' },
])
</script>

<template>
  <div class="space-y-6">
    <!-- hero -->
    <div class="alt-surface flex flex-col gap-1 px-6 py-6">
      <h2 class="display-heading text-2xl text-midnight-ink">Welcome back, {{ auth.user?.name }} 👋</h2>
      <p class="text-sm text-fog">Here's the state of your payment infrastructure.</p>
    </div>

    <!-- stat cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div
        v-for="s in stats"
        :key="s.label"
        class="card-surface flex items-start gap-3 p-4"
      >
        <div
          class="flex size-9 shrink-0 items-center justify-center rounded-lg"
          :class="{
            'bg-electric-blue/10 text-electric-blue': s.color === 'blue',
            'bg-lavender/10 text-lavender': s.color === 'violet',
            'bg-vivid-green/10 text-vivid-green': s.color === 'green',
            'bg-tangerine/10 text-tangerine': s.color === 'orange',
          }"
        >
          <component :is="s.icon" class="size-4.5" />
        </div>
        <div class="min-w-0">
          <div class="text-xs font-medium text-fog">{{ s.label }}</div>
          <div class="mt-0.5 truncate text-xl font-bold tracking-tight text-midnight-ink">{{ s.value }}</div>
          <div class="text-[11px] text-silver">{{ s.sub }}</div>
        </div>
      </div>
    </div>

    <!-- content grid -->
    <div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
      <!-- providers -->
      <section class="card-surface p-5">
        <div class="mb-3 flex items-center justify-between">
          <h3 class="text-sm font-semibold text-midnight-ink">Payment providers</h3>
          <RouterLink to="/app/providers" class="text-xs font-medium text-electric-blue hover:underline">View all →</RouterLink>
        </div>
        <div v-if="providers.query.isLoading" class="py-8 text-center text-sm text-fog">Loading…</div>
        <div v-else-if="!providerCount" class="py-8 text-center text-sm text-fog">
          No providers yet — add one from the Providers menu.
        </div>
        <ul v-else class="divide-y divide-ash">
          <li v-for="p in providers.query.data" :key="p.id" class="flex items-center justify-between py-2.5">
            <div class="flex items-center gap-2.5">
              <span class="pill" :class="p.enabled ? 'pill-green' : 'pill-blue'">{{ p.provider }}</span>
              <span class="text-sm font-medium text-charcoal">{{ p.name }}</span>
            </div>
            <span class="text-[11px] text-silver">{{ p.label }}</span>
          </li>
        </ul>
      </section>

      <!-- wallets -->
      <section class="card-surface p-5">
        <div class="mb-3 flex items-center justify-between">
          <h3 class="text-sm font-semibold text-midnight-ink">Wallets</h3>
          <RouterLink to="/app/wallets" class="text-xs font-medium text-electric-blue hover:underline">View all →</RouterLink>
        </div>
        <div v-if="wallets.query.isLoading" class="py-8 text-center text-sm text-fog">Loading…</div>
        <div v-else-if="!walletCount" class="py-8 text-center text-sm text-fog">
          No wallets yet — add one from the Wallets menu.
        </div>
        <ul v-else class="divide-y divide-ash">
          <li v-for="w in wallets.query.data" :key="w.id" class="flex items-center justify-between py-2.5">
            <div class="flex items-center gap-2.5">
              <Wallet class="size-4 text-steel" />
              <span class="text-sm font-medium text-charcoal">{{ w.name }}</span>
            </div>
            <span class="text-sm font-semibold text-midnight-ink">
              {{ fmtIDR.format(w.balance ?? 0) }}
              <span class="text-xs font-normal text-fog">{{ w.currency }}</span>
            </span>
          </li>
        </ul>
      </section>
    </div>

    <!-- team -->
    <section class="card-surface p-5">
      <div class="mb-3 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-midnight-ink">Team</h3>
        <RouterLink to="/app/users" class="text-xs font-medium text-electric-blue hover:underline">Manage →</RouterLink>
      </div>
      <div v-if="!userCount" class="py-6 text-center text-sm text-fog">No users yet.</div>
      <ul v-else class="divide-y divide-ash">
        <li v-for="u in users.query.data" :key="u.id" class="flex items-center gap-3 py-2.5">
          <div class="flex size-8 items-center justify-center rounded-full bg-electric-blue/10 text-xs font-semibold text-electric-blue">
            {{ u.name.slice(0, 1).toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate text-sm font-medium text-charcoal">{{ u.name }}</div>
            <div class="truncate text-xs text-fog">{{ u.email }}</div>
          </div>
          <span class="pill pill-blue">{{ u.role }}</span>
        </li>
      </ul>
    </section>
  </div>
</template>
