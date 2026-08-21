import { defineStore } from 'pinia'
import { useQuery, useMutation, useQueryCache } from '@pinia/colada'
import { createItem, deleteItem, fetchList, updateItem } from '@/lib/api'
import type { Provider, User, Wallet } from '@/lib/types'

export const useUsersStore = defineStore('users', () => {
  const cache = useQueryCache()

  const query = useQuery({
    key: ['users'],
    query: () => fetchList<User>('users', 'users'),
  })

  const create = useMutation({
    mutation: (body: Partial<User> & { password: string }) => createItem<User>('users', body),
    onSuccess: () => cache.invalidateQueries({ key: ['users'] }, true),
  })
  const update = useMutation({
    mutation: (body: Partial<User> & { id: string }) => updateItem<User>(`users/${body.id}`, body),
    onSuccess: () => cache.invalidateQueries({ key: ['users'] }, true),
  })
  const remove = useMutation({
    mutation: (id: string) => deleteItem(`users/${id}`),
    onSuccess: () => cache.invalidateQueries({ key: ['users'] }, true),
  })

  return { query, create, update, remove }
})

export const useWalletsStore = defineStore('wallets', () => {
  const cache = useQueryCache()

  const query = useQuery({
    key: ['wallets'],
    query: () => fetchList<Wallet>('wallets', 'wallets'),
  })

  const create = useMutation({
    mutation: (body: Partial<Wallet>) => createItem<Wallet>('wallets', body),
    onSuccess: () => cache.invalidateQueries({ key: ['wallets'] }, true),
  })
  const update = useMutation({
    mutation: (body: Partial<Wallet> & { id: string }) => updateItem<Wallet>(`wallets/${body.id}`, body),
    onSuccess: () => cache.invalidateQueries({ key: ['wallets'] }, true),
  })
  const remove = useMutation({
    mutation: (id: string) => deleteItem(`wallets/${id}`),
    onSuccess: () => cache.invalidateQueries({ key: ['wallets'] }, true),
  })

  return { query, create, update, remove }
})

export const useProvidersStore = defineStore('providers', () => {
  const cache = useQueryCache()

  const query = useQuery({
    key: ['providers'],
    query: () => fetchList<Provider>('providers', 'providers'),
  })

  const create = useMutation({
    mutation: (body: Partial<Provider>) => createItem<Provider>('providers', body),
    onSuccess: () => cache.invalidateQueries({ key: ['providers'] }, true),
  })
  const update = useMutation({
    mutation: (body: Partial<Provider> & { id: string }) => updateItem<Provider>(`providers/${body.id}`, body),
    onSuccess: () => cache.invalidateQueries({ key: ['providers'] }, true),
  })
  const remove = useMutation({
    mutation: (id: string) => deleteItem(`providers/${id}`),
    onSuccess: () => cache.invalidateQueries({ key: ['providers'] }, true),
  })

  return { query, create, update, remove }
})
