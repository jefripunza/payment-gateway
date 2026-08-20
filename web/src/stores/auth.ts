import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  api,
  clearSession,
  fetchMe,
  getStoredUser,
  getToken,
  login as apiLogin,
  setSession,
} from '@/lib/api'
import type { AuthUser } from '@/lib/types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(getStoredUser())
  const initialized = ref(false)

  async function login(email: string, password: string) {
    const res = await apiLogin(email, password)
    user.value = res.user
    return res.user
  }

  async function init() {
    if (initialized.value) return
    const token = getToken()
    if (token) {
      try {
        user.value = await fetchMe()
      } catch {
        clearSession()
        user.value = null
      }
    }
    initialized.value = true
  }

  function logout() {
    clearSession()
    user.value = null
  }

  async function refreshMe() {
    if (!getToken()) return
    user.value = await fetchMe()
  }

  return { user, initialized, login, init, logout, refreshMe }
})
