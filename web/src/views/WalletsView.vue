<script setup lang="ts">
import { computed, ref } from 'vue'
import { Loader2, Plus, Pencil, Trash2, Wallet } from 'lucide-vue-next'
import type { Wallet as WalletType } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { useWalletsStore } from '@/stores/resources'

const store = useWalletsStore()

const wallets = computed(() => store.query.data ?? [])
const loading = computed(() => store.query.isLoading)
const error = computed(() => store.query.error?.message ?? '')

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref({ name: '', currency: 'IDR', balance: 0 })
const saving = ref(false)
const formError = ref('')

const deleteOpen = ref(false)
const deleteTarget = ref<WalletType | null>(null)
const deleteError = ref('')
const deleting = ref(false)

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
      await store.update.mutate({ id: editingId.value, ...form.value })
    } else {
      await store.create.mutate({ ...form.value })
    }
    dialogOpen.value = false
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
    await store.remove.mutate(deleteTarget.value.id)
    deleteOpen.value = false
    deleteTarget.value = null
  } catch (e: any) {
    deleteError.value = e?.message ?? 'Failed to delete wallet'
  } finally {
    deleting.value = false
  }
}

const totalBalance = computed(() => wallets.value.reduce((acc: number, w: WalletType) => acc + w.balance, 0))

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
      <div class="alt-surface p-5">
        <div class="text-xs font-medium text-fog">Total wallets</div>
        <div class="mt-1 text-2xl font-bold tracking-tight text-midnight-ink">{{ wallets.length }}</div>
      </div>
      <div class="alt-surface p-5 sm:col-span-2">
        <div class="text-xs font-medium text-fog">Combined balance</div>
        <div class="mt-1 truncate text-2xl font-bold tracking-tight text-midnight-ink">
          {{ fmt(totalBalance, 'IDR') }}
        </div>
      </div>
    </div>

    <div class="flex items-center justify-between gap-3">
      <div class="text-sm text-fog">{{ wallets.length }} wallet{{ wallets.length !== 1 ? 's' : '' }}</div>
      <button type="button" class="btn-primary" @click="openCreate">
        <Plus class="size-4" />
        Add wallet
      </button>
    </div>

    <div v-if="error" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-4 py-3 text-sm text-tangerine">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-fog">
      <Loader2 class="size-4 animate-spin" />
      Loading wallets…
    </div>

    <!-- desktop table -->
    <div v-else-if="!error" class="card-surface hidden overflow-hidden md:block">
      <table class="table-surface">
        <thead>
          <tr>
            <th>Name</th>
            <th>Currency</th>
            <th>Balance</th>
            <th>Created</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="w in wallets" :key="w.id">
            <td>
              <div class="flex items-center gap-3">
                <div class="flex size-8 items-center justify-center rounded-lg bg-lavender/10">
                  <Wallet class="size-4 text-lavender" />
                </div>
                <span class="font-medium text-charcoal">{{ w.name }}</span>
              </div>
            </td>
            <td>
              <span class="pill pill-blue">{{ w.currency }}</span>
            </td>
            <td class="font-semibold text-vivid-green">{{ fmt(w.balance, w.currency) }}</td>
            <td class="text-fog">{{ formatDate(w.createdAt) }}</td>
            <td>
              <div class="flex justify-end gap-1">
                <button type="button" title="Edit" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(w)">
                  <Pencil class="size-4" />
                </button>
                <button type="button" title="Delete" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(w)">
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
      <div v-for="w in wallets" :key="w.id" class="card-surface p-4">
        <div class="flex items-center gap-3">
          <div class="flex size-10 items-center justify-center rounded-xl bg-lavender/10">
            <Wallet class="size-5 text-lavender" />
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium text-charcoal">{{ w.name }}</div>
            <div class="text-xs text-fog">{{ w.currency }}</div>
          </div>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(w)">
            <Pencil class="size-4" />
          </button>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(w)">
            <Trash2 class="size-4" />
          </button>
        </div>
        <div class="mt-3 flex items-center justify-between border-t border-ash pt-3">
          <span class="text-sm font-semibold text-vivid-green">{{ fmt(w.balance, w.currency) }}</span>
          <span class="text-xs text-silver">{{ formatDate(w.createdAt) }}</span>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-midnight-ink/20 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="card-surface relative w-full max-w-md p-6 shadow-xl">
        <h3 class="text-base font-semibold text-midnight-ink">
          {{ editingId ? 'Edit wallet' : 'Add wallet' }}
        </h3>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Name</label>
            <input v-model="form.name" type="text" required placeholder="e.g. Main bank account" class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Currency</label>
            <input v-model="form.currency" type="text" required maxlength="8" placeholder="IDR" class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Balance</label>
            <input v-model.number="form.balance" type="number" min="0" required class="input-surface" />
          </div>
          <div v-if="formError" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs font-medium text-tangerine">
            {{ formError }}
          </div>
          <div class="flex justify-end gap-2 pt-1">
            <button type="button" class="btn-ghost" @click="dialogOpen = false">Cancel</button>
            <button type="submit" class="btn-primary" :disabled="saving">
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
