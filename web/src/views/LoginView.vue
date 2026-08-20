<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CreditCard, Eye, EyeOff, Loader2, Lock, Mail, ShieldCheck, Zap } from 'lucide-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref('')
const loading = ref(false)

async function submit() {
  if (!email.value || !password.value) {
    error.value = 'Email and password are required'
    return
  }
  error.value = ''
  loading.value = true
  try {
    await auth.login(email.value, password.value)
    router.replace('/app/dashboard')
  } catch (e: any) {
    error.value = e?.message ?? 'Login failed'
    if (e?.response?.status === 401) error.value = 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div
    class="relative min-h-screen overflow-hidden bg-slate-950 font-sans text-slate-100 antialiased"
  >
    <!-- ambient background -->
    <div class="pointer-events-none absolute inset-0">
      <div
        class="absolute -top-40 -left-40 h-[32rem] w-[32rem] rounded-full bg-indigo-600/30 blur-[120px]"
      />
      <div
        class="absolute top-1/3 -right-40 h-[28rem] w-[28rem] rounded-full bg-violet-600/25 blur-[120px]"
      />
      <div
        class="absolute bottom-0 left-1/3 h-[24rem] w-[24rem] rounded-full bg-sky-600/20 blur-[120px]"
      />
      <div
        class="absolute inset-0 opacity-[0.04]"
        style="background-image: radial-gradient(circle at 1px 1px, white 1px, transparent 0)"
        style2="background-size: 32px 32px"
      />
    </div>

    <div class="relative z-10 mx-auto flex min-h-screen w-full max-w-6xl items-center justify-center px-4 py-10">
      <div class="grid w-full items-center gap-10 lg:grid-cols-2">
        <!-- left brand panel -->
        <div class="hidden lg:block">
          <div class="mb-8 flex items-center gap-3">
            <div class="flex size-11 items-center justify-center rounded-2xl bg-gradient-to-br from-indigo-500 to-violet-600 shadow-lg shadow-indigo-900/50">
              <CreditCard class="size-6 text-white" />
            </div>
            <div>
              <div class="text-xl font-bold tracking-tight">Payment Gateway</div>
              <div class="text-sm text-slate-400">payment.sawang.tech</div>
            </div>
          </div>

          <h1 class="max-w-md text-4xl font-extrabold leading-tight tracking-tight">
            Manage every payment provider
            <span class="bg-gradient-to-r from-indigo-400 via-violet-400 to-sky-400 bg-clip-text text-transparent">
              from one dashboard
            </span>
          </h1>
          <p class="mt-4 max-w-md text-slate-400">
            Store credentials for multiple payment gateways, track wallets, and manage team access —
            securely, in one place.
          </p>

          <div class="mt-10 grid max-w-md grid-cols-3 gap-4">
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <ShieldCheck class="mb-2 size-5 text-emerald-400" />
              <div class="text-sm font-semibold">Encrypted</div>
              <div class="text-xs text-slate-400">AES-256 at rest</div>
            </div>
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <Zap class="mb-2 size-5 text-amber-400" />
              <div class="text-sm font-semibold">Multi-gateway</div>
              <div class="text-xs text-slate-400">One dashboard</div>
            </div>
            <div class="rounded-2xl border border-white/10 bg-white/5 p-4 backdrop-blur">
              <Lock class="mb-2 size-5 text-sky-400" />
              <div class="text-sm font-semibold">Role-based</div>
              <div class="text-xs text-slate-400">Admin & viewer</div>
            </div>
          </div>
        </div>

        <!-- login card -->
        <div class="mx-auto w-full max-w-md">
          <div class="rounded-3xl border border-white/10 bg-white/[0.06] p-8 shadow-2xl shadow-black/40 backdrop-blur-xl">
            <div class="mb-8 flex items-center gap-3 lg:hidden">
              <div class="flex size-11 items-center justify-center rounded-2xl bg-gradient-to-br from-indigo-500 to-violet-600">
                <CreditCard class="size-6 text-white" />
              </div>
              <div>
                <div class="text-lg font-bold tracking-tight">Payment Gateway</div>
                <div class="text-xs text-slate-400">payment.sawang.tech</div>
              </div>
            </div>

            <h2 class="text-2xl font-bold tracking-tight">Welcome back</h2>
            <p class="mt-1 text-sm text-slate-400">Sign in to your payment dashboard</p>

            <form class="mt-8 space-y-5" @submit.prevent="submit">
              <div>
                <label for="email" class="mb-1.5 block text-sm font-medium text-slate-300">Email</label>
                <div class="relative">
                  <Mail class="pointer-events-none absolute top-1/2 left-3.5 size-4.5 -translate-y-1/2 text-slate-500" />
                  <input
                    id="email"
                    v-model="email"
                    type="email"
                    required
                    autocomplete="email"
                    placeholder="you@example.com"
                    class="h-11 w-full rounded-xl border border-white/10 bg-white/5 pr-4 pl-10 text-sm text-slate-100 placeholder:text-slate-500 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label for="password" class="mb-1.5 block text-sm font-medium text-slate-300">Password</label>
                <div class="relative">
                  <Lock class="pointer-events-none absolute top-1/2 left-3.5 size-4.5 -translate-y-1/2 text-slate-500" />
                  <input
                    id="password"
                    v-model="password"
                    :type="showPassword ? 'text' : 'password'"
                    required
                    autocomplete="current-password"
                    placeholder="••••••••"
                    class="h-11 w-full rounded-xl border border-white/10 bg-white/5 pr-11 pl-10 text-sm text-slate-100 placeholder:text-slate-500 focus:border-indigo-400/60 focus:ring-2 focus:ring-indigo-500/30 focus:outline-none"
                  />
                  <button
                    type="button"
                    class="absolute top-1/2 right-3 -translate-y-1/2 text-slate-500 transition hover:text-slate-300"
                    @click="showPassword = !showPassword"
                  >
                    <Eye v-if="!showPassword" class="size-4.5" />
                    <EyeOff v-else class="size-4.5" />
                  </button>
                </div>
              </div>

              <div
                v-if="error"
                class="rounded-xl border border-red-500/30 bg-red-500/10 px-4 py-2.5 text-sm text-red-300"
              >
                {{ error }}
              </div>

              <button
                type="submit"
                :disabled="loading"
                class="flex h-11 w-full items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-indigo-500 to-violet-600 text-sm font-semibold text-white shadow-lg shadow-indigo-900/40 transition hover:from-indigo-400 hover:to-violet-500 disabled:cursor-not-allowed disabled:opacity-60"
              >
                <Loader2 v-if="loading" class="size-4 animate-spin" />
                <span>{{ loading ? 'Signing in…' : 'Sign in' }}</span>
              </button>
            </form>
          </div>
          <p class="mt-6 text-center text-xs text-slate-500">
            Secured connection · Credentials are encrypted at rest
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
