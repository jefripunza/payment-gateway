<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CreditCard, Eye, EyeOff, Loader2, Lock, Mail, ShieldCheck } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const show = ref(false)
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  if (!email.value || !password.value) {
    error.value = 'Enter your email and password.'
    return
  }
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.replace('/app/dashboard')
  } catch {
    error.value = 'Invalid email or password.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-paper-mist px-4 py-10">
    <!-- frosted card on rice paper -->
    <div class="w-full max-w-md">
      <div class="mb-8 flex flex-col items-center text-center">
        <div class="mb-4 flex size-11 items-center justify-center rounded-xl bg-midnight-ink">
          <CreditCard class="size-5 text-white" />
        </div>
        <h1 class="display-heading text-[28px] text-midnight-ink">Payment Gateway</h1>
        <p class="mt-1 text-sm text-fog">payment.sawang.tech</p>
      </div>

      <div class="card-surface p-6 shadow-sm sm:p-8">
        <h2 class="display-heading text-xl text-midnight-ink">Welcome back</h2>
        <p class="mt-1 text-sm text-fog">Sign in to your payment dashboard</p>

        <form class="mt-6 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite" for="email">Email</label>
            <div class="relative">
              <Mail class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-silver" />
              <input
                id="email"
                v-model="email"
                type="email"
                autocomplete="email"
                placeholder="you@example.com"
                class="input-surface pl-9"
              />
            </div>
          </div>

          <div>
            <label class="mb-1.5 block text-xs font-semibold text-graphite" for="password">Password</label>
            <div class="relative">
              <Lock class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-silver" />
              <input
                id="password"
                v-model="password"
                :type="show ? 'text' : 'password'"
                autocomplete="current-password"
                placeholder="••••••••"
                class="input-surface pr-10 pl-9"
              />
              <button
                type="button"
                class="absolute top-1/2 right-3 -translate-y-1/2 text-silver transition hover:text-graphite"
                @click="show = !show"
              >
                <EyeOff v-if="show" class="size-4" />
                <Eye v-else class="size-4" />
              </button>
            </div>
          </div>

          <p v-if="error" class="rounded-lg border border-tangerine/20 bg-tangerine/5 px-3 py-2 text-xs font-medium text-tangerine">
            {{ error }}
          </p>

          <button type="submit" class="btn-primary w-full justify-center" :disabled="loading">
            <Loader2 v-if="loading" class="size-4 animate-spin" />
            {{ loading ? 'Signing in…' : 'Sign in' }}
          </button>
        </form>
      </div>

      <p class="mt-6 flex items-center justify-center gap-1.5 text-xs text-fog">
        <ShieldCheck class="size-3.5 text-vivid-green" />
        Secured connection · Credentials are encrypted at rest
      </p>
    </div>
  </div>
</template>
