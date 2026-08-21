<script setup lang="ts">
import { computed, ref } from 'vue'
import { Loader2, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import type { User } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { useUsersStore } from '@/stores/resources'

const store = useUsersStore()

const users = computed(() => store.query.data ?? [])
const loading = computed(() => store.query.isLoading)
const error = computed(() => store.query.error?.message ?? '')

// dialog state
const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref<{ name: string; email: string; password: string; role: 'admin' | 'viewer'; active: boolean }>({
  name: '',
  email: '',
  password: '',
  role: 'viewer',
  active: true,
})
const saving = ref(false)
const formError = ref('')

const deleteOpen = ref(false)
const deleteTarget = ref<User | null>(null)
const deleteError = ref('')
const deleting = ref(false)

function openCreate() {
  editingId.value = null
  form.value = { name: '', email: '', password: '', role: 'viewer', active: true }
  formError.value = ''
  dialogOpen.value = true
}

function openEdit(u: User) {
  editingId.value = u.id
  form.value = { name: u.name, email: u.email, password: '', role: u.role, active: u.active }
  formError.value = ''
  dialogOpen.value = true
}

async function submit() {
  formError.value = ''
  if (!form.value.name || !form.value.email) {
    formError.value = 'Name and email are required'
    return
  }
  if (form.value.password && form.value.password.length < 8) {
    formError.value = 'Password must be at least 8 characters'
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await store.update.mutate({
        id: editingId.value,
        name: form.value.name,
        email: form.value.email,
        role: form.value.role,
        active: form.value.active,
        ...(form.value.password ? { password: form.value.password } : {}),
      })
    } else {
      await store.create.mutate({ ...form.value } as any)
    }
    dialogOpen.value = false
  } catch (e: any) {
    formError.value = e?.message ?? 'Failed to save user'
    if (e?.response?.status === 409) formError.value = 'User with this email already exists'
  } finally {
    saving.value = false
  }
}

function requestDelete(u: User) {
  deleteTarget.value = u
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
    deleteError.value = e?.message ?? 'Failed to delete user'
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
      <div class="text-sm text-fog">{{ users.length }} user{{ users.length !== 1 ? 's' : '' }}</div>
      <button type="button" class="btn-primary" @click="openCreate">
        <Plus class="size-4" />
        Add user
      </button>
    </div>

    <div v-if="error" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-4 py-3 text-sm text-tangerine">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-fog">
      <Loader2 class="size-4 animate-spin" />
      Loading users…
    </div>

    <!-- desktop table -->
    <div v-else-if="!error" class="card-surface hidden overflow-hidden md:block">
      <table class="table-surface">
        <thead>
          <tr>
            <th>Name</th>
            <th>Email</th>
            <th>Role</th>
            <th>Status</th>
            <th>Created</th>
            <th class="text-right">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="u in users" :key="u.id">
            <td>
              <div class="flex items-center gap-3">
                <div class="flex size-8 shrink-0 items-center justify-center rounded-full bg-electric-blue/10 text-[10px] font-bold text-electric-blue">
                  {{ u.name.split(' ').map((p: string) => p[0]).slice(0, 2).join('').toUpperCase() }}
                </div>
                <span class="font-medium text-charcoal">{{ u.name }}</span>
              </div>
            </td>
            <td class="text-fog">{{ u.email }}</td>
            <td>
              <span class="pill" :class="u.role === 'admin' ? 'pill-blue' : 'pill-violet'">{{ u.role }}</span>
            </td>
            <td>
              <span class="inline-flex items-center gap-1.5 text-xs font-medium" :class="u.active ? 'text-vivid-green' : 'text-silver'">
                <span class="size-1.5 rounded-full" :class="u.active ? 'bg-vivid-green' : 'bg-silver'" />
                {{ u.active ? 'Active' : 'Disabled' }}
              </span>
            </td>
            <td class="text-fog">{{ formatDate(u.createdAt) }}</td>
            <td>
              <div class="flex justify-end gap-1">
                <button type="button" title="Edit" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(u)">
                  <Pencil class="size-4" />
                </button>
                <button type="button" title="Delete" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(u)">
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
      <div v-for="u in users" :key="u.id" class="card-surface p-4">
        <div class="flex items-center gap-3">
          <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-electric-blue/10 text-xs font-bold text-electric-blue">
            {{ u.name.split(' ').map((p: string) => p[0]).slice(0, 2).join('').toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium text-charcoal">{{ u.name }}</div>
            <div class="truncate text-xs text-fog">{{ u.email }}</div>
          </div>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-paper-mist hover:text-charcoal" @click="openEdit(u)">
            <Pencil class="size-4" />
          </button>
          <button type="button" class="rounded-lg p-2 text-fog transition hover:bg-tangerine/10 hover:text-tangerine" @click="requestDelete(u)">
            <Trash2 class="size-4" />
          </button>
        </div>
        <div class="mt-3 flex items-center gap-2 border-t border-ash pt-3 text-xs">
          <span class="pill" :class="u.role === 'admin' ? 'pill-blue' : 'pill-violet'">{{ u.role }}</span>
          <span class="ml-auto flex items-center gap-1.5" :class="u.active ? 'text-vivid-green' : 'text-silver'">
            <span class="size-1.5 rounded-full" :class="u.active ? 'bg-vivid-green' : 'bg-silver'" />
            {{ u.active ? 'Active' : 'Disabled' }}
          </span>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-midnight-ink/20 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="card-surface relative w-full max-w-md p-6 shadow-xl">
        <h3 class="text-base font-semibold text-midnight-ink">
          {{ editingId ? 'Edit user' : 'Add user' }}
        </h3>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Name</label>
            <input v-model="form.name" type="text" required class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">Email</label>
            <input v-model="form.email" type="email" required class="input-surface" />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite">
              {{ editingId ? 'New password' : 'Password' }}
            </label>
            <input
              v-model="form.password"
              type="password"
              :required="!editingId"
              minlength="8"
              :placeholder="editingId ? 'Leave blank to keep current' : 'At least 8 characters'"
              class="input-surface"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1.5 block text-xs font-semibold text-graphite">Role</label>
              <select v-model="form.role" class="input-surface">
                <option value="viewer">Viewer</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-semibold text-graphite">Status</label>
              <select v-model="form.active" class="input-surface">
                <option :value="true">Active</option>
                <option :value="false">Disabled</option>
              </select>
            </div>
          </div>
          <div v-if="formError" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs font-medium text-tangerine">
            {{ formError }}
          </div>
          <div class="flex justify-end gap-2 pt-1">
            <button type="button" class="btn-ghost" @click="dialogOpen = false">Cancel</button>
            <button type="submit" class="btn-primary" :disabled="saving">
              <Loader2 v-if="saving" class="size-3.5 animate-spin" />
              {{ editingId ? 'Save changes' : 'Create user' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <ConfirmDeleteDialog
      v-model:open="deleteOpen"
      title="Delete user"
      description="This will permanently remove the user from the dashboard."
      item-name="user"
      :confirm-text="deleteTarget?.name ?? ''"
      :error="deleteError"
      :saving="deleting"
      @confirm="confirmDelete"
    />
  </div>
</template>
