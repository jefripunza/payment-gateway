<script setup lang="ts">
import { ref } from 'vue'
import {
  CreditCard,
  Eye,
  EyeOff,
  Loader2,
  Pencil,
  Plus,
  ShieldCheck,
  Trash2,
} from 'lucide-vue-next'
import { api } from '@/lib/api'
import type { Provider } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const PROVIDER_TYPES = ['midtrans', 'xendit', 'tripay', 'duitku', 'paypal', 'stripe', 'other']

const providers = ref<Provider[]>([])
const loading = ref(true)
const error = ref('')

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref({
  name: '',
  type: 'midtrans',
  isProduction: false,
  merchantId: '',
  apiKey: '',
  apiSecret: '',
  webhookKey: '',
  enabled: true,
})
const saving = ref(false)
const formError = ref('')

const deleteOpen = ref(false)
const deleteTarget = ref<Provider | null>(null)
const deleteError = ref('')
const deleting = ref(false)

// masked value visibility
const reveal = ref<Record<string, boolean>>({})

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('providers').json<{ providers: Provider[] }>()
    providers.value = res.providers
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load providers'
  } finally {
    loading.value = false
  }
}
load()

function openCreate() {
  editingId.value = null
  form.value = {
    name: '',
    type: 'midtrans',
    isProduction: false,
    merchantId: '',
    apiKey: '',
    apiSecret: '',
    webhookKey: '',
    enabled: true,
  }
  formError.value = ''
  dialogOpen.value = true
}

function openEdit(p: Provider) {
  editingId.value = p.id
  form.value = {
    name: p.name,
    type: p.type,
    isProduction: p.isProduction,
    merchantId: p.merchantId,
    apiKey: p.apiKey,
    apiSecret: p.apiSecret,
    webhookKey: p.webhookKey,
    enabled: p.enabled,
  }
  formError.value = ''
  dialogOpen.value = true
}

async function submit() {
  formError.value = ''
  if (!form.value.name) {
    formError.value = 'Name is required'
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await api.patch(`providers/${editingId.value}`, { json: form.value })
    } else {
      await api.post('providers', { json: form.value })
    }
    dialogOpen.value = false
    await load()
  } catch (e: any) {
    formError.value = e?.message ?? 'Failed to save provider'
    if (e?.response?.status === 409) formError.value = 'Provider with this name already exists'
  } finally {
    saving.value = false
  }
}

function requestDelete(p: Provider) {
  deleteTarget.value = p
  deleteError.value = ''
  deleteOpen.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`providers/${deleteTarget.value.id}`)
    deleteOpen.value = false
    deleteTarget.value = null
    await load()
  } catch (e: any) {
    deleteError.value = e?.message ?? 'Failed to delete provider'
  } finally {
    deleting.value = false
  }
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-5">
    <div class="flex items-center justify-between gap-3">
      <div class="text-sm text-slate-500">{{ providers.length }} provider{{ providers.length !== 1 ? 's' : '' }}</div>
      <button
        type="button"
        class="flex h-9 items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-violet-600 px-4 text-sm font-semibold text-white shadow-lg shadow-indigo-900/30 transition hover:from-indigo-400 hover:to-violet-500"
        @click="openCreate"
      >
        <Plus class="size-4" />
        Add provider
      </button>
    </div>

    <div v-if="error" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-slate-500">
      <Loader2 class="size-4 animate-spin" />
      Loading providers…
    </div>

    <div v-else-if="providers.length === 0" class="rounded-2xl border border-dashed border-white/15 py-16 text-center">
      <CreditCard class="mx-auto size-10 text-slate-600" />
      <div class="mt-3 text-sm font-medium text-slate-400">No payment providers yet</div>
      <p class="mt-1 text-xs text-slate-600">
        Add a provider to store its API credentials — e.g. Midtrans, Xendit, Tripay.
      </p>
    </div>

    <!-- provider cards -->
    <div v-else class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-3">
      <div
        v-for="p in providers"
        :key="p.id"
        class="group rounded-2xl border border-white/10 bg-white/[0.03] p-5 transition hover:border-indigo-400/30 hover:bg-white/[0.05]"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="flex items-center gap-3">
            <div class="flex size-10 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500/25 to-violet-500/15">
              <CreditCard class="size-5 text-indigo-300" />
            </div>
            <div>
              <div class="font-semibold text-slate-100">{{ p.name }}</div>
              <div class="text-xs text-slate-500">{{ p.type }}</div>
            </div>
          </div>
          <div class="flex gap-1">
            <button
              type="button"
              title="Edit"
              class="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
              @click="openEdit(p)"
            >
              <Pencil class="size-4" />
            </button>
            <button
              type="button"
              title="Delete"
              class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400"
              @click="requestDelete(p)"
            >
              <Trash2 class="size-4" />
            </button>
          </div>
        </div>

        <div class="mt-4 flex flex-wrap gap-2">
          <span
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
            :class="p.enabled ? 'bg-emerald-500/15 text-emerald-300' : 'bg-slate-500/15 text-slate-400'"
          >
            <span class="size-1.5 rounded-full" :class="p.enabled ? 'bg-emerald-400' : 'bg-slate-500'" />
            {{ p.enabled ? 'Enabled' : 'Disabled' }}
          </span>
          <span
            class="inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-semibold"
            :class="
              p.isProduction
                ? 'bg-amber-500/15 text-amber-300'
                : 'bg-sky-500/15 text-sky-300'
            "
          >
            <ShieldCheck class="size-3" />
            {{ p.isProduction ? 'Production' : 'Sandbox' }}
          </span>
        </div>

        <div class="mt-4 space-y-1.5 border-t border-white/5 pt-4 text-xs">
          <div class="flex items-center justify-between">
            <span class="text-slate-500">Merchant ID</span>
            <button
              type="button"
              class="flex items-center gap-1.5 font-mono text-slate-300 transition hover:text-slate-100"
              @click="reveal[p.id] = !reveal[p.id]"
            >
              <Eye v-if="reveal[p.id]" class="size-3 text-slate-500" />
              <EyeOff v-else class="size-3 text-slate-500" />
              {{ reveal[p.id] ? p.merchantId : '••••••••' }}
            </button>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-500">API Key</span>
            <span class="font-mono text-slate-400">{{ p.apiKey }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-500">API Secret</span>
            <span class="font-mono text-slate-400">{{ p.apiSecret }}</span>
          </div>
          <div class="flex items-center justify-between">
            <span class="text-slate-500">Webhook Key</span>
            <span class="font-mono text-slate-400">{{ p.webhookKey }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center overflow-y-auto p-4">
      <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="relative my-auto w-full max-w-lg rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
        <h3 class="text-base font-bold text-slate-100">
          {{ editingId ? 'Edit provider' : 'Add provider' }}
        </h3>
        <p class="mt-1 text-xs text-slate-500">
          Credentials are encrypted at rest and never shown in full again.
        </p>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-300">Name</label>
              <input
                v-model="form.name"
                type="text"
                required
                placeholder="e.g. Midtrans Production"
                class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
              />
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-300">Type</label>
              <select
                v-model="form.type"
                class="h-10 w-full rounded-lg border border-white/10 bg-slate-900 px-3 text-sm text-slate-100 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
              >
                <option v-for="t in PROVIDER_TYPES" :key="t" :value="t">{{ t }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <label
              class="flex cursor-pointer items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-2.5 text-xs font-medium text-slate-300"
            >
              <input v-model="form.enabled" type="checkbox" class="size-3.5 accent-emerald-500" />
              Enabled
            </label>
            <label
              class="flex cursor-pointer items-center gap-2 rounded-lg border border-white/10 bg-white/5 px-3 py-2.5 text-xs font-medium text-slate-300"
            >
              <input v-model="form.isProduction" type="checkbox" class="size-3.5 accent-amber-500" />
              Production
            </label>
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Merchant ID / Client ID</label>
            <input
              v-model="form.merchantId"
              type="text"
              autocomplete="off"
              :placeholder="editingId ? 'Leave blank to keep current' : 'e.g. G123456789'"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 font-mono text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">API Key</label>
            <input
              v-model="form.apiKey"
              type="password"
              autocomplete="new-password"
              :placeholder="editingId ? 'Leave blank to keep current' : 'Server key / API key'"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 font-mono text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">API Secret</label>
            <input
              v-model="form.apiSecret"
              type="password"
              autocomplete="new-password"
              :placeholder="editingId ? 'Leave blank to keep current' : 'Secret key'"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 font-mono text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Webhook Key</label>
            <input
              v-model="form.webhookKey"
              type="password"
              autocomplete="new-password"
              :placeholder="editingId ? 'Leave blank to keep current' : 'Webhook verification key (optional)'"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 font-mono text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
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
              {{ editingId ? 'Save changes' : 'Add provider' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <ConfirmDeleteDialog
      v-model:open="deleteOpen"
      title="Delete provider"
      description="The provider credentials will be permanently removed."
      item-name="provider"
      :confirm-text="deleteTarget?.name ?? ''"
      :error="deleteError"
      :saving="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>
