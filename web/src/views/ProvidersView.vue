<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { CreditCard, Loader2, Pencil, Plus, ShieldCheck, Trash2 } from 'lucide-vue-next'
import type { CatalogMethod, Provider } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { useMethodsStore, useProvidersStore } from '@/stores/resources'

const store = useProvidersStore()
const methodsStore = useMethodsStore()

const providers = computed(() => store.query.data ?? [])
const loading = computed(() => store.query.isLoading)
const error = computed(() => store.query.error?.message ?? '')

// ---- form state ----
const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref<{
  name: string
  method: string
  creds: Record<string, string>
}>({
  name: '',
  method: '',
  creds: {},
})
const saving = ref(false)
const formError = ref('')

// ---- delete state ----
const deleteOpen = ref(false)
const deleteTarget = ref<Provider | null>(null)
const deleteError = ref('')
const deleting = ref(false)

// ---- toggle state (switch on card) ----
const togglingId = ref<string | null>(null)

// ---- methods catalog ----
const catalog = computed<CatalogMethod[]>(() => methodsStore.query.data?.methods ?? [])
const methodsLoading = computed(() => methodsStore.query.isLoading)

// current method fields (dynamic) for the selected method
const selectedMethod = computed<CatalogMethod | undefined>(() =>
  catalog.value.find((m) => `${m.provider}|${m.value}` === form.value.method),
)
const methodFields = computed(() => selectedMethod.value?.fields ?? [])

function openCreate() {
  editingId.value = null
  form.value = { name: '', method: '', creds: {} }
  formError.value = ''
  dialogOpen.value = true
}

function openEdit(p: Provider) {
  editingId.value = p.id
  form.value = {
    name: p.name,
    method: p.method,
    creds: { ...p.creds },
  }
  formError.value = ''
  dialogOpen.value = true
}

// when method changes, reset creds (fields differ per method)
watch(
  () => form.value.method,
  () => {
    form.value.creds = {}
  },
)

async function submit() {
  formError.value = ''
  if (!form.value.name.trim()) {
    formError.value = 'Name is required'
    return
  }
  if (!form.value.method) {
    formError.value = 'Pick a payment method'
    return
  }
  // required credential fields must be filled on create
  const missing = methodFields.value.filter((f) => f.required && !form.value.creds[f.key]?.trim())
  if (!editingId.value && missing.length) {
    formError.value = `Required fields: ${missing.map((f) => f.label).join(', ')}`
    return
  }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      method: form.value.method,
      creds: form.value.creds,
    }
    if (editingId.value) {
      await store.update.mutate({ id: editingId.value, ...payload })
    } else {
      await store.create.mutate(payload)
    }
    dialogOpen.value = false
  } catch (e: any) {
    formError.value = e?.message ?? 'Failed to save provider'
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(p: Provider) {
  togglingId.value = p.id
  try {
    await store.update.mutate({ id: p.id, enabled: !p.enabled })
  } catch (e: any) {
    console.error('toggle failed', e)
  } finally {
    togglingId.value = null
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

    <!-- provider cards -->
    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-3">
      <div v-for="p in providers" :key="p.id" class="card-surface flex flex-col p-5">
        <div class="flex items-start justify-between gap-3">
          <div class="flex items-center gap-3">
            <div class="flex size-10 items-center justify-center rounded-xl bg-electric-blue/10">
              <CreditCard class="size-5 text-electric-blue" />
            </div>
            <div class="min-w-0">
              <div class="truncate font-medium text-charcoal">{{ p.name }}</div>
              <div class="text-xs text-fog">{{ p.provider }} · {{ p.label }}</div>
            </div>
          </div>
          <!-- enable/disable switch -->
          <button
            type="button"
            role="switch"
            :aria-checked="p.enabled"
            :title="p.enabled ? 'Disable provider' : 'Enable provider'"
            class="relative h-6 w-11 shrink-0 rounded-full transition-colors disabled:opacity-50"
            :class="p.enabled ? 'bg-electric-blue' : 'bg-ash'"
            :disabled="togglingId === p.id"
            @click="toggleEnabled(p)"
          >
            <span
              class="absolute top-0.5 size-5 rounded-full bg-white shadow transition-all"
              :class="p.enabled ? 'left-[22px]' : 'left-0.5'"
            />
          </button>
        </div>

        <div class="mt-4 space-y-1.5 border-t border-ash pt-3">
          <div v-for="(v, k) in p.creds" :key="k" class="flex items-center justify-between text-xs">
            <span class="capitalize text-fog">{{ k.replace(/([A-Z])/g, ' $1') }}</span>
            <span class="font-mono text-graphite">{{ v }}</span>
          </div>
          <div v-if="!Object.keys(p.creds ?? {}).length" class="text-xs text-fog">No credentials</div>
        </div>

        <div class="mt-auto flex items-center justify-between border-t border-ash pt-3 text-xs text-fog">
          <span>{{ formatDate(p.createdAt) }}</span>
          <div class="flex gap-1">
            <button type="button" title="Edit" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(p)">
              <Pencil class="size-4" />
            </button>
            <button type="button" title="Delete" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(p)">
              <Trash2 class="size-4" />
            </button>
          </div>
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
            <input v-model="form.name" type="text" required placeholder="e.g. Midtrans QRIS Production" class="input-surface" />
          </div>

          <!-- method dropdown (from go-payment-method catalog) -->
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Method</label>
            <div v-if="methodsLoading" class="flex items-center gap-2 py-2 text-xs text-fog">
              <Loader2 class="size-3.5 animate-spin" /> Loading methods…
            </div>
            <select v-else v-model="form.method" class="input-surface">
              <option value="" disabled>Select a payment method…</option>
              <optgroup v-for="g in catalog" :key="g.provider" :label="g.provider">
                <option :value="`${g.provider}|${g.value}`">{{ g.label }}</option>
              </optgroup>
            </select>
          </div>

          <!-- dynamic credential fields per method -->
          <div v-if="methodFields.length" class="space-y-4">
            <template v-for="f in methodFields" :key="f.key">
              <div>
                <label class="mb-1.5 block text-xs font-semibold text-graphite">
                  {{ f.label }}
                  <span v-if="f.required" class="text-tangerine">*</span>
                </label>
                <input
                  v-model="form.creds[f.key]"
                  type="password"
                  autocomplete="new-password"
                  :placeholder="editingId ? (f.required ? 'Leave blank to keep current' : 'Leave blank to keep current') : (f.placeholder ?? '')"
                  class="input-surface"
                />
              </div>
            </template>
          </div>
          <div v-else-if="form.method" class="rounded-lg border border-dashed border-ash px-3 py-2 text-xs text-fog">
            No credential fields for this method.
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
