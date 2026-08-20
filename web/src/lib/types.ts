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
  type: string
  isProduction: boolean
  merchantId: string
  apiKey: string
  apiSecret: string
  webhookKey: string
  enabled: boolean
  createdAt: string
  updatedAt: string
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
