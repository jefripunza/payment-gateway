<script setup lang="ts">
import { ref } from 'vue'
import { Loader2, ShieldAlert } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const show = ref(false)
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const error = ref('')
const saving = ref(false)

const open = defineModel<boolean>('open', { default: false })

function reset() {
  show.value = false
  oldPassword.value = ''
  newPassword.value = ''
  confirmPassword.value = ''
  error.value = ''
}

async function submit() {
  error.value = ''
  if (newPassword.value.length < 8) {
    error.value = 'New password must be at least 8 characters'
    return
  }
  if (newPassword.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }
  saving.value = true
  try {
    const { changePassword } = await import('@/lib/api')
    await changePassword(oldPassword.value, newPassword.value)
    open.value = false
    reset()
  } catch (e: any) {
    error.value = e?.message ?? 'Failed to change password'
    if (e?.response?.status === 401) error.value = 'Current password is incorrect'
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div v-if="open" class="fixed inset-0 z-[60] flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-slate-950/70 backdrop-blur-sm" @click="open = false" />
    <div class="relative w-full max-w-md rounded-2xl border border-white/10 bg-slate-900 p-6 shadow-2xl">
      <h3 class="flex items-center gap-2 text-base font-bold text-slate-100">
        <ShieldAlert class="size-4.5 text-amber-400" />
        Change password
      </h3>
      <p class="mt-1 text-xs text-slate-500">
        Signed in as {{ auth.user?.email }}
      </p>

      <form class="mt-5 space-y-4" @submit.prevent="submit">
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-300">Current password</label>
          <input
            v-model="oldPassword"
            :type="show ? 'text' : 'password'"
            required
            autocomplete="current-password"
            class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-300">New password</label>
          <input
            v-model="newPassword"
            :type="show ? 'text' : 'password'"
            required
            minlength="8"
            autocomplete="new-password"
            placeholder="At least 8 characters"
            class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-medium text-slate-300">Confirm new password</label>
          <input
            v-model="confirmPassword"
            :type="show ? 'text' : 'password'"
            required
            autocomplete="new-password"
            class="h-10 w-full rounded-lg border border-white/10 bg-white/5 px-3 text-sm text-slate-100 placeholder:text-slate-600 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
          />
        </div>

        <label class="flex items-center gap-2 text-xs text-slate-400">
          <input v-model="show" type="checkbox" class="size-3.5 accent-indigo-500" />
          Show passwords
        </label>

        <div
          v-if="error"
          class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-xs text-red-300"
        >
          {{ error }}
        </div>

        <div class="flex justify-end gap-2 pt-1">
          <button
            type="button"
            class="h-9 rounded-lg px-4 text-sm font-medium text-slate-400 transition hover:bg-white/5 hover:text-slate-200"
            @click="open = false"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="flex h-9 items-center gap-2 rounded-lg bg-gradient-to-r from-indigo-500 to-violet-600 px-4 text-sm font-semibold text-white transition hover:from-indigo-400 hover:to-violet-500 disabled:opacity-60"
          >
            <Loader2 v-if="saving" class="size-3.5 animate-spin" />
            Save
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
