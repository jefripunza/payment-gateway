<script setup lang="ts">
import { ref } from 'vue'
import { Loader2, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import { api } from '@/lib/api'
import type { User } from '@/lib/types'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const users = ref<User[]>([])
const loading = ref(true)
const error = ref('')

// dialog state
const dialogOpen = ref(false)
const editingId = ref<string | null>(null)
const form = ref({ name: '', email: '', password: '', role: 'viewer', active: true })
const saving = ref(false)
const formError = ref('')

const deleteOpen = ref(false)
const deleteTarget = ref<User | null>(null)
const deleteError = ref('')
const deleting = ref(false)

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await api.get('users').json<{ users: User[] }>()
    users.value = res.users
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to load users'
  } finally {
    loading.value = false
  }
}
load()

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
  if (!editingId.value && form.value.password.length < 8) {
    formError.value = 'Password must be at least 8 characters'
    return
  }
  if (form.value.password && form.value.password.length < 8) {
    formError.value = 'Password must be at least 8 characters'
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await api.patch(`users/${editingId.value}`, {
        json: {
          name: form.value.name,
          email: form.value.email,
          role: form.value.role,
          active: form.value.active,
          ...(form.value.password ? { password: form.value.password } : {}),
        },
      })
    } else {
      await api.post('users', { json: form.value })
    }
    dialogOpen.value = false
    await load()
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
    await api.delete(`users/${deleteTarget.value.id}`)
    deleteOpen.value = false
    deleteTarget.value = null
    await load()
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
      <div class="text-sm text-slate-500">{{ users.length }} user{{ users.length !== 1 ? 's' : '' }}</div>
      <button
        type="button"
        class="flex h-9 items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-violet-600 px-4 text-sm font-semibold text-white shadow-lg shadow-indigo-900/30 transition hover:from-indigo-400 hover:to-violet-500"
        @click="openCreate"
      >
        <Plus class="size-4" />
        Add user
      </button>
    </div>

    <div v-if="error" class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-3 text-sm text-red-300">
      {{ error }}
    </div>

    <div v-if="loading" class="flex items-center justify-center gap-2 py-16 text-sm text-slate-500">
      <Loader2 class="size-4 animate-spin" />
      Loading users…
    </div>

    <!-- desktop table -->
    <div v-else-if="!error" class="hidden overflow-hidden rounded-2xl border border-white/10 bg-white/[0.03] md:block">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-white/10 text-left text-xs text-slate-500 uppercase">
            <th class="px-5 py-3.5 font-semibold">Name</th>
            <th class="px-5 py-3.5 font-semibold">Email</th>
            <th class="px-5 py-3.5 font-semibold">Role</th>
            <th class="px-5 py-3.5 font-semibold">Status</th>
            <th class="px-5 py-3.5 font-semibold">Created</th>
            <th class="px-5 py-3.5 text-right font-semibold">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="u in users"
            :key="u.id"
            class="border-b border-white/5 last:border-0 transition hover:bg-white/[0.03]"
          >
            <td class="px-5 py-3.5">
              <div class="flex items-center gap-3">
                <div class="flex size-8 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 text-[10px] font-bold text-white">
                  {{ u.name.split(' ').map((p) => p[0]).slice(0, 2).join('').toUpperCase() }}
                </div>
                <span class="font-medium text-slate-200">{{ u.name }}</span>
              </div>
            </td>
            <td class="px-5 py-3.5 text-slate-400">{{ u.email }}</td>
            <td class="px-5 py-3.5">
              <span
                class="inline-flex rounded-full px-2.5 py-0.5 text-xs font-semibold"
                :class="u.role === 'admin' ? 'bg-indigo-500/15 text-indigo-300' : 'bg-slate-500/15 text-slate-300'"
              >
                {{ u.role }}
              </span>
            </td>
            <td class="px-5 py-3.5">
              <span
                class="inline-flex items-center gap-1.5 text-xs font-medium"
                :class="u.active ? 'text-emerald-400' : 'text-slate-500'"
              >
                <span class="size-1.5 rounded-full" :class="u.active ? 'bg-emerald-400' : 'bg-slate-600'" />
                {{ u.active ? 'Active' : 'Disabled' }}
              </span>
            </td>
            <td class="px-5 py-3.5 text-slate-500">{{ formatDate(u.createdAt) }}</td>
            <td class="px-5 py-3.5">
              <div class="flex justify-end gap-1">
                <button
                  type="button"
                  title="Edit"
                  class="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
                  @click="openEdit(u)"
                >
                  <Pencil class="size-4" />
                </button>
                <button
                  type="button"
                  title="Delete"
                  class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400"
                  @click="requestDelete(u)"
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
      <div
        v-for="u in users"
        :key="u.id"
        class="rounded-2xl border border-white/10 bg-white/[0.03] p-4"
      >
        <div class="flex items-center gap-3">
          <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-gradient-to-br from-indigo-500 to-violet-600 text-xs font-bold text-white">
            {{ u.name.split(' ').map((p) => p[0]).slice(0, 2).join('').toUpperCase() }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="truncate font-medium text-slate-200">{{ u.name }}</div>
            <div class="truncate text-xs text-slate-500">{{ u.email }}</div>
          </div>
          <button
            type="button"
            class="rounded-lg p-2 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
            @click="openEdit(u)"
          >
            <Pencil class="size-4" />
          </button>
          <button
            type="button"
            class="rounded-lg p-2 text-slate-500 transition hover:bg-red-500/10 hover:text-red-400"
            @click="requestDelete(u)"
          >
            <Trash2 class="size-4" />
          </button>
        </div>
        <div class="mt-3 flex items-center gap-2 border-t border-white/5 pt-3 text-xs">
          <span
            class="rounded-full px-2 py-0.5 font-semibold"
            :class="u.role === 'admin' ? 'bg-indigo-500/15 text-indigo-300' : 'bg-slate-500/15 text-slate-300'"
          >
            {{ u.role }}
          </span>
          <span class="ml-auto flex items-center gap-1.5" :class="u.active ? 'text-emerald-400' : 'text-slate-500'">
            <span class="size-1.5 rounded-full" :class="u.active ? 'bg-emerald-400' : 'bg-slate-600'" />
            {{ u.active ? 'Active' : 'Disabled' }}
          </span>
        </div>
      </div>
    </div>

    <!-- create / edit dialog -->
    <div v-if="dialogOpen" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" @click="dialogOpen = false" />
      <div class="relative w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
        <h3 class="text-base font-bold text-slate-100">
          {{ editingId ? 'Edit user' : 'Add user' }}
        </h3>
        <form class="mt-5 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Name</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">Email</label>
            <input
              v-model="form.email"
              type="email"
              required
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div>
            <label class="mb-1.5 block text-xs font-medium text-slate-300">
              {{ editingId ? 'New password' : 'Password' }}
            </label>
            <input
              v-model="form.password"
              type="password"
              :required="!editingId"
              minlength="8"
              :placeholder="editingId ? 'Leave blank to keep current' : 'At least 8 characters'"
              class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
            />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-300">Role</label>
              <select
                v-model="form.role"
                class="h-10 w-full rounded-lg border border-white/10 bg-slate-900 px-3 text-sm text-slate-100 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
              >
                <option value="viewer">Viewer</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div>
              <label class="mb-1.5 block text-xs font-medium text-slate-300">Status</label>
              <select
                v-model="form.active"
                class="h-10 w-full rounded-lg border border-white/10 bg-slate-900 px-3 text-sm text-slate-100 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
              >
                <option :value="true">Active</option>
                <option :value="false">Disabled</option>
              </select>
            </div>
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
