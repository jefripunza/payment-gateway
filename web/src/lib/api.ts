import ky from 'ky'
import type { AuthUser, ProviderMethod } from './types'

const TOKEN_KEY = 'payment:token'
const USER_KEY = 'payment:user'

export const api = ky.create({
  prefix: '/api/v1/',
  timeout: 15000,
  hooks: {
    beforeRequest: [
      ({ request }) => {
        const token = getToken()
        if (token) request.headers.set('Authorization', `Bearer ${token}`)
      },
    ],
  },
})

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getStoredUser(): AuthUser | null {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? (JSON.parse(raw) as AuthUser) : null
  } catch {
    return null
  }
}

export function setSession(token: string, user: AuthUser) {
  localStorage.setItem(TOKEN_KEY, token)
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

export async function login(email: string, password: string) {
  const res = await api
    .post('auth/login', { json: { email, password } })
    .json<{ token: string; user: AuthUser }>()
  setSession(res.token, res.user)
  return res
}

export async function fetchMe(): Promise<AuthUser> {
  const user = await api.get('auth/me').json<AuthUser>()
  localStorage.setItem(USER_KEY, JSON.stringify(user))
  return user
}

export async function changePassword(oldPassword: string, newPassword: string) {
  return api.patch('auth/password', { json: { oldPassword, newPassword } }).json<{ ok: boolean }>()
}

// ---------- typed resource helpers ----------

export async function fetchList<T>(path: string, key: string): Promise<T[]> {
  const res = await api.get(path).json<Record<string, T[]>>()
  return res[key] ?? []
}

export async function createItem<T>(path: string, body: unknown): Promise<T> {
  return api.post(path, { json: body }).json<T>()
}

export async function updateItem<T>(path: string, body: unknown): Promise<T> {
  return api.patch(path, { json: body }).json<T>()
}

export async function deleteItem(path: string): Promise<void> {
  await api.delete(path)
}

// ---------- provider methods catalog ----------

export interface MethodsCatalog {
  providers: Record<string, ProviderMethod[]>
}

export async function fetchMethods(): Promise<MethodsCatalog> {
  return api.get('methods').json<MethodsCatalog>()
}
