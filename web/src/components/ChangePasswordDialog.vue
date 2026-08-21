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
    <div class="absolute inset-0 bg-midnight-ink/20 backdrop-blur-sm" @click="open = false" />
    <div class="card-surface relative w-full max-w-md p-6 shadow-xl">
      <h3 class="flex items-center gap-2 text-base font-semibold text-midnight-ink">
        <ShieldAlert class="size-4.5 text-tangerine" />
        Change password
      </h3>
      <p class="mt-1 text-xs text-fog">
        Signed in as {{ auth.user?.email }}
      </p>

      <form class="mt-5 space-y-4" @submit.prevent="submit">
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-graphite">Current password</label>
          <input
            v-model="oldPassword"
            :type="show ? 'text' : 'password'"
            required
            autocomplete="current-password"
            class="input-surface"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-graphite">New password</label>
          <input
            v-model="newPassword"
            :type="show ? 'text' : 'password'"
            required
            minlength="8"
            autocomplete="new-password"
            placeholder="At least 8 characters"
            class="input-surface"
          />
        </div>
        <div>
          <label class="mb-1.5 block text-xs font-semibold text-graphite">Confirm new password</label>
          <input
            v-model="confirmPassword"
            :type="show ? 'text' : 'password'"
            required
            autocomplete="new-password"
            class="input-surface"
          />
        </div>

        <label class="flex items-center gap-2 text-xs text-fog">
          <input v-model="show" type="checkbox" class="size-3.5 accent-electric-blue" />
          Show passwords
        </label>

        <div v-if="error" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs font-medium text-tangerine">
          {{ error }}
        </div>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn-ghost" @click="open = false">Cancel</button>
          <button type="submit" class="btn-primary" :disabled="saving">
            <Loader2 v-if="saving" class="size-3.5 animate-spin" />
            Save
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
