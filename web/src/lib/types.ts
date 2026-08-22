export interface User {
  id: string
  name: string
  email: string
  role: 'admin' | 'viewer'
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface Wallet {
  id: string
  name: string
  currency: string
  balance: number
  createdAt: string
  updatedAt: string
}

export interface Provider {
  id: string
  name: string
  method: string // "provider|value"
  provider: string
  label: string
  creds: Record<string, string> // masked values
  enabled: boolean
  createdAt: string
  updatedAt: string
}

export interface MethodField {
  key: string
  label: string
  required: boolean
  placeholder?: string
}

export interface CatalogMethod {
  provider: string
  value: string
  label: string
  fields: MethodField[]
}

export interface AuthUser {
  id: string
  name: string
  email: string
  role: 'admin' | 'viewer'
  active: boolean
  createdAt: string
  updatedAt: string
}
