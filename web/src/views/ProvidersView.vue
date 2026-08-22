<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  CreditCard,
  Eye,
  EyeOff,
  Loader2,
  Pencil,
  Plus,
  ShieldCheck,
  Trash2,
  X,
} from 'lucide-vue-next'
import type { Provider, ProviderMethod } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { useMethodsStore, useProvidersStore } from '@/stores/resources'

const PROVIDER_TYPES = ['midtrans', 'xendit', 'tripay', 'duitku', 'paypal', 'stripe', 'other'] as const

const store = useProvidersStore()
const methodsStore = useMethodsStore()

const providers = computed(() => store.query.data ?? [])
const loading = computed(() => store.query.isLoading)
const error = computed(() => store.query.error?.message ?? '')

const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref<{
  name: string
  type: string
  methods: string[]
  isProduction: boolean
  merchantId: string
  apiKey: string
  apiSecret: string
  webhookKey: string
  enabled: boolean
}>({
  name: '',
  type: 'midtrans',
  methods: [],
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

// methods available for the currently selected type
const availableMethods = computed<ProviderMethod[]>(() => {
  const cat = methodsStore.query.data?.providers ?? {}
  return cat[form.value.type] ?? []
})

function openCreate() {
  editingId.value = null
  form.value = {
    name: '',
    type: 'midtrans',
    methods: [],
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
    methods: Array.isArray(p.methods) ? [...p.methods] : [],
    isProduction: p.isProduction,
    merchantId: p.merchantId ?? '',
    apiKey: p.apiKey ?? '',
    apiSecret: p.apiSecret ?? '',
    webhookKey: p.webhookKey ?? '',
    enabled: p.enabled,
  }
  formError.value = ''
  dialogOpen.value = true
}

function toggleMethod(m: string) {
  const idx = form.value.methods.indexOf(m)
  if (idx >= 0) {
    form.value.methods.splice(idx, 1)
  } else {
    form.value.methods.push(m)
  }
}

function isMethodSelected(m: string) {
  return form.value.methods.includes(m)
}

// reset selected methods when provider type changes
watch(
  () => form.value.type,
  () => {
    form.value.methods = []
  },
)

async function submit() {
  formError.value = ''
  if (!form.value.name) {
    formError.value = 'Name is required'
    return
  }
  if (!editingId.value && !form.value.apiKey && !form.value.apiSecret) {
    formError.value = 'Add at least an API key or secret for the provider'
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
    formError.value = e?.message ?? 'Failed to save provider'
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
    await store.remove.mutate(deleteTarget.value.id)
    deleteOpen.value = false
    deleteTarget.value = null
  } catch (e: any) {
    deleteError.value = e?.message ?? 'Failed to delete provider'
  } finally {
    deleting.value = false
  }
}

function mask(v: string) {
  if (!v) return '—'
  if (v.length <= 8) return '••••'
  return `${v.slice(0, 4)}••••••${v.slice(-4)}`
}

function toggleReveal(id: string) {
  reveal.value[id] = !reveal.value[id]
}

function isRevealed(id: string) {
  return !!reveal.value[id]
}

function formatDate(s: string) {
  return new Date(s).toLocaleDateString('en-GB', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="space-y-5">
    <!-- toolbar -->
    <div class="flex items-center justify-between gap-3">
      <div class="text-sm text-fog">{{ providers.length }} provider{{ providers.length !== 1 ? 's' : '' }}</div>
      <button type="button" class="btn-primary" @click="openCreate">
        <Plus class="size-4" />
        Add provider
      </button>
    </div>

    <div v-if="error" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-4 py-3 text-sm text-tangerine">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-fog">
      <Loader2 class="size-4 animate-spin" />
      Loading providers…
    </div>

    <!-- empty state -->
    <div v-else-if="!error && !providers.length" class="card-surface flex flex-col items-center justify-center gap-3 px-6 py-16 text-center">
      <div class="flex size-12 items-center justify-center rounded-xl bg-paper-mist">
        <CreditCard class="size-6 text-silver" />
      </div>
      <div>
        <div class="text-sm font-semibold text-charcoal">No payment providers yet</div>
        <p class="mt-1 max-w-sm text-sm text-fog">
          Add a payment gateway (Midtrans, Xendit, Tripay, …) to store its credentials securely — encrypted at rest.
        </p>
      </div>
      <button type="button" class="btn-primary mt-2" @click="openCreate">
        <Plus class="size-4" />
        Add your first provider
      </button>
    </div>

    <!-- desktop table -->
    <div v-else-if="!error && providers.length" class="card-surface hidden overflow-hidden md:block">
      <table class="table-surface">
        <thead>
          <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Methods</th>
            <th>Credentials</th>
            <th>Environment</th>
            <th>Status</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="p in providers" :key="p.id">
            <td>
              <div class="flex items-center gap-3">
                <div class="flex size-8 items-center justify-center rounded-lg bg-electric-blue/10">
                  <CreditCard class="size-4 text-electric-blue" />
                </div>
                <span class="font-medium text-charcoal">{{ p.name }}</span>
              </div>
            </td>
            <td><span class="pill pill-blue">{{ p.type }}</span></td>
            <td>
              <div v-if="p.methods?.length" class="flex max-w-[220px] flex-wrap gap-1">
                <span v-for="m in p.methods" :key="m" class="pill pill-ash">{{ m }}</span>
              </div>
              <span v-else class="text-xs text-fog">—</span>
            </td>
            <td>
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs text-fog">{{ mask(p.apiKey) }}</span>
                <button
                  type="button"
                  class="rounded p-1 text-silver transition hover:text-charcoal"
                  :title="isRevealed(p.id) ? 'Hide credential' : 'Show credential'"
                  @click="toggleReveal(p.id)"
                >
                  <EyeOff v-if="isRevealed(p.id)" class="size-3.5" />
                  <Eye v-else class="size-3.5" />
                </button>
              </div>
            </td>
            <td>
              <span class="pill" :class="p.isProduction ? 'pill-orange' : 'pill-violet'">
                {{ p.isProduction ? 'Production' : 'Sandbox' }}
              </span>
            </td>
            <td>
              <span class="inline-flex items-center gap-1.5 text-xs font-medium" :class="p.enabled ? 'text-vivid-green' : 'text-silver'">
                <span class="size-1.5 rounded-full" :class="p.enabled ? 'bg-vivid-green' : 'bg-silver'" />
                {{ p.enabled ? 'Enabled' : 'Disabled' }}
              </span>
            </td>
            <td>
              <div class="flex justify-end gap-1">
                <button type="button" title="Edit" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(p)">
                  <Pencil class="size-4" />
                </button>
                <button type="button" title="Delete" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(p)">
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
      <div v-for="p in providers" :key="p.id" class="card-surface p-4">
        <div class="flex items-center gap-3">
          <div class="flex size-10 items-center justify-center rounded-xl bg-electric-blue/10">
            <CreditCard class="size-5 text-electric-blue" />
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium text-charcoal">{{ p.name }}</div>
            <div class="text-xs text-fog">{{ p.type }}</div>
          </div>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(p)">
            <Pencil class="size-4" />
          </button>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(p)">
            <Trash2 class="size-4" />
          </button>
        </div>
        <div v-if="p.methods?.length" class="mt-3 flex flex-wrap gap-1">
          <span v-for="m in p.methods" :key="m" class="pill pill-ash">{{ m }}</span>
        </div>
        <div class="mt-3 flex items-center justify-between border-t border-ash pt-3">
          <span class="inline-flex items-center gap-1.5 text-xs font-medium" :class="p.enabled ? 'text-vivid-green' : 'text-silver'">
            <span class="size-1.5 rounded-full" :class="p.enabled ? 'bg-vivid-green' : 'bg-silver'" />
            {{ p.enabled ? 'Enabled' : 'Disabled' }}
          </span>
          <span class="pill" :class="p.isProduction ? 'pill-orange' : 'pill-violet'">
            {{ p.isProduction ? 'Production' : 'Sandbox' }}
          </span>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-midnight-ink/20 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="card-surface relative max-h-[90vh] w-full max-w-lg overflow-y-auto p-6 shadow-xl">
        <h3 class="text-base font-semibold text-midnight-ink">
          {{ editingId ? 'Edit provider' : 'Add provider' }}
        </h3>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Name</label>
            <input v-model="form.name" type="text" required placeholder="e.g. Midtrans Production" class="input-surface" />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1.5 block text-xs font-semibold text-graphite">Type</label>
              <select v-model="form.type" class="input-surface">
                <option v-for="t in PROVIDER_TYPES" :key="t" :value="t">{{ t }}</option>
              </select>
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-semibold text-graphite">Environment</label>
              <select v-model="form.isProduction" class="input-surface">
                <option :value="false">Sandbox</option>
                <option :value="true">Production</option>
              </select>
            </div>
          </div>

          <!-- payment methods (from go-payment-method library) -->
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Payment Methods</label>
            <div v-if="availableMethods.length" class="grid max-h-44 grid-cols-2 gap-1 overflow-y-auto rounded-lg border border-ash p-2">
              <label
                v-for="m in availableMethods"
                :key="m.value"
                class="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm text-graphite transition hover:bg-paper-mist"
              >
                <input
                  type="checkbox"
                  class="size-4 rounded border-smoke text-electric-blue focus:ring-electric-blue/30"
                  :checked="isMethodSelected(m.value)"
                  @change="toggleMethod(m.value)"
                />
                <span class="truncate">{{ m.label }}</span>
              </label>
            </div>
            <div v-else class="rounded-lg border border-dashed border-ash px-3 py-2 text-xs text-fog">
              No method list available for {{ form.type }}.
            </div>
            <div v-if="form.methods.length" class="mt-2 flex flex-wrap gap-1">
              <span v-for="m in form.methods" :key="m" class="pill pill-blue">
                {{ m }}
                <button type="button" class="ml-1 rounded hover:text-charcoal" @click="toggleMethod(m)">
                  <X class="size-3" />
                </button>
              </span>
            </div>
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Merchant ID</label>
            <input v-model="form.merchantId" type="text" placeholder="Optional" class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">API Key</label>
            <input v-model="form.apiKey" type="text" :placeholder="editingId ? 'Leave blank to keep current' : 'Server key / API key'" class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">API Secret</label>
            <input v-model="form.apiSecret" type="text" :placeholder="editingId ? 'Leave blank to keep current' : 'Secret key / client secret'" class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Webhook Key</label>
            <input v-model="form.webhookKey" type="text" :placeholder="editingId ? 'Leave blank to keep current' : 'Optional' " class="input-surface" />
          </div>
          <div class="flex items-center gap-2">
            <input v-model="form.enabled" type="checkbox" id="provider-enabled" class="size-4 rounded border-smoke text-electric-blue focus:ring-electric-blue/30" />
            <label for="provider-enabled" class="text-sm text-graphite">Enabled</label>
          </div>
          <div v-if="formError" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs font-medium text-tangerine">
            {{ formError }}
          </div>
          <div class="flex items-center justify-between gap-2 pt-1">
            <span class="flex items-center gap-1.5 text-xs text-fog">
              <ShieldCheck class="size-3.5 text-vivid-green" />
              Encrypted at rest
            </span>
            <div class="flex gap-2">
              <button type="button" class="btn-ghost" @click="dialogOpen = false">Cancel</button>
              <button type="submit" class="btn-primary" :disabled="saving">
                <Loader2 v-if="saving" class="size-3.5 animate-spin" />
                {{ editingId ? 'Save changes' : 'Create provider' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <ConfirmDeleteDialog
      v-model:open="deleteOpen"
      title="Delete provider"
      description="This will permanently remove the provider and its credentials."
      item-name="provider"
      :confirm-text="deleteTarget?.name ?? ''"
      :error="deleteError"
      :saving="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>
