<script setup lang="ts">
import { computed, ref } from 'vue'
import { Loader2, Plus, Pencil, Trash2, Wallet } from 'lucide-vue-next'
import { api } from '@/lib/api'
import type { Wallet as WalletType } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const wallets = ref<WalletType[]>([])
const loading = ref(true)
const error = ref('')

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref({ name: '', currency: 'IDR', balance: 0 })
const saving = ref(false)
const formError = ref('')

const deleteOpen = ref(false)
const deleteTarget = ref<WalletType | null>(null)
const deleteError = ref('')
const deleting = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('wallets').json<{ wallets: WalletType[] }>()
    wallets.value = res.wallets
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load wallets'
  } finally {
    loading.value = false
  }
}
load()

function openCreate() {
  editingId.value = null
  form.value = { name: '', currency: 'IDR', balance: 0 }
  formError.value = ''
  dialogOpen.value = true
}

function openEdit(w: WalletType) {
  editingId.value = w.id
  form.value = { name: w.name, currency: w.currency, balance: w.balance }
  formError.value = ''
  dialogOpen.value = true
}

async function submit() {
  formError.value = ''
  if (!form.value.name) {
    formError.value = 'Name is required'
    return
  }
  if (form.value.balance < 0) {
    formError.value = 'Balance cannot be negative'
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await api.patch(`wallets/${editingId.value}`, { json: form.value })
    } else {
      await api.post('wallets', { json: form.value })
    }
    dialogOpen.value = false
    await load()
  } catch (e: any) {
    formError.value = e?.message ?? 'Failed to save wallet'
  } finally {
    saving.value = false
  }
}

function requestDelete(w: WalletType) {
  deleteTarget.value = w
  deleteError.value = ''
  deleteOpen.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`wallets/${deleteTarget.value.id}`)
    deleteOpen.value = false
    deleteTarget.value = null
    await load()
  } catch (e: any) {
    deleteError.value = e?.message ?? 'Failed to delete wallet'
  } finally {
    deleting.value = false
  }
}

const totalBalance = computed(() => wallets.value.reduce((acc, w) => acc + w.balance, 0))

function fmt(n: number, cur: string) {
  try {
    return new Intl.NumberFormat('en-US', { style: 'currency', currency: cur }).format(n)
  } catch {
    return `${cur} ${n.toLocaleString()}`
  }
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-5">
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5">
        <div class="text-xs text-slate-500">Total wallets</div>
        <div class="mt-1 text-2xl font-extrabold tracking-tight text-slate-100">{{ wallets.length }}</div>
      </div>
      <div class="rounded-2xl border border-white/10 bg-white/[0.03] p-5 sm:col-span-2">
        <div class="text-xs text-slate-500">Combined balance</div>
        <div class="mt-1 truncate text-2xl font-extrabold tracking-tight text-slate-100">
          {{ fmt(totalBalance, 'IDR') }}
        </div>
      </div>
    </div>

    <div class="flex items-center justify-between gap-3">
      <div class="text-sm text-slate-500">{{ wallets.length }} wallet{{ wallets.length !== 1 ? 's' : '' }}</div>
      <button
        type="button"
        class="flex h-9 items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-violet-600 px-4 text-sm font-semibold text-white shadow-lg shadow-indigo-900/30 transition hover:from-indigo-400 hover:to-violet-500"
        @click="openCreate"
      >
        <Plus class="size-4" />
        Add wallet
      </button>
    </div>

    <div v-if="error" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-slate-500">
      <Loader2 class="size-4 animate-spin" />
      Loading wallets…
    </div>

    <!-- desktop table -->
    <div v-else-if="!error" class="hidden overflow-hidden rounded-2xl border border-white/10 bg-white/[0.03] md:block">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-white/10 text-left text-xs text-slate-500 uppercase">
            <th class="px-5 py-3.5 font-semibold">Name</th>
            <th class="px-5 py-3.5 font-semibold">Currency</th>
            <th class="px-5 py-3.5 font-semibold">Balance</th>
            <th class="px-5 py-3.5 font-semibold">Created</th>
            <th class="px-5 py-3.5 text-right font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="w in wallets"
            :key="w.id"
            class="border-b border-white/5 last:border-0 transition hover:bg-white/[0.03]"
          >
            <td class="px-5 py-3.5">
              <div class="flex items-center gap-3">
                <div class="flex size-8 items-center justify-center rounded-lg bg-gradient-to-br from-amber-500/20 to-orange-500/10">
                  <Wallet class="size-4 text-amber-400" />
                </div>
                <span class="font-medium text-slate-200">{{ w.name }}</span>
              </div>
            </td>
            <td class="px-5 py-3.5">
              <span class="rounded-full bg-slate-500/15 px-2.5 py-0.5 text-xs font-semibold text-slate-300">
                {{ w.currency }}
              </span>
            </td>
            <td class="px-5 py-3.5 font-semibold text-emerald-400">{{ fmt(w.balance, w.currency) }}</td>
            <td class="px-5 py-3.5 text-slate-500">{{ formatDate(w.createdAt) }}</td>
            <td class="px-5 py-3.5">
              <div class="flex justify-end gap-1">
                <button
                  type="button"
                  title="Edit"
                  class="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
                  @click="openEdit(w)"
                >
                  <Pencil class="size-4" />
                </button>
                <button
                  type="button"
                  title="Delete"
                  class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400"
                  @click="requestDelete(w)"
                >
                  <Trash2 class="size-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- mobile cards -->
    <div v-else class="space-y-3">
      <div v-for="w in wallets" :key="w.id" class="rounded-2xl border border-white/10 bg-white/[0.03] p-4">
        <div class="flex items-center gap-3">
          <div class="flex size-10 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-500/10">
            <Wallet class="size-5 text-amber-400" />
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium text-slate-200">{{ w.name }}</div>
            <div class="text-xs text-slate-500">{{ w.currency }}</div>
          </div>
          <button type="button" class="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-slate-200" @click="openEdit(w)">
            <Pencil class="size-4" />
          </button>
          <button type="button" class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400" @click="requestDelete(w)">
            <Trash2 class="size-4" />
          </button>
        </div>
        <div class="mt-3 flex items-center justify-between border-t border-white/5 pt-3">
          <span class="text-sm font-semibold text-emerald-400">{{ fmt(w.balance, w.currency) }}</span>
          <span class="text-xs text-slate-600">{{ formatDate(w.createdAt) }}</span>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="relative w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
        <h3 class="text-base font-bold text-slate-100">
          {{ editingId ? 'Edit wallet' : 'Add wallet' }}
        </h3>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Name</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="e.g. Main bank account"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Currency</label>
            <input
              v-model="form.currency"
              type="text"
              required
              maxlength="8"
              placeholder="IDR"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Balance</label>
            <input
              v-model.number="form.balance"
              type="number"
              min="0"
              required
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div v-if="formError" class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300">
            {{ formError }}
          </div>
          <div class="flex justify-end gap-2 pt-1">
            <button
              type="button"
              class="h-9 rounded-lg px-4 text-sm font-medium text-slate-400 transition hover:bg-white/5 hover:text-slate-200"
              @click="dialogOpen = false"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="flex h-9 items-center gap-2 rounded-lg bg-gradient-to-r from-indigo-500 to-violet-600 px-4 text-sm font-semibold text-white transition hover:from-indigo-400 hover:to-violet-500 disabled:opacity-60"
            >
              <Loader2 v-if="saving" class="size-3.5 animate-spin" />
              {{ editingId ? 'Save changes' : 'Create wallet' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <ConfirmDeleteDialog
      v-model:open="deleteOpen"
      title="Delete wallet"
      description="This will permanently remove the wallet from the dashboard."
      item-name="wallet"
      :confirm-text="deleteTarget?.name ?? ''"
      :error="deleteError"
      :saving="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>
