import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/LoginView.vue'
import DashboardView from '@/views/DashboardView.vue'
import UsersView from '@/views/UsersView.vue'
import WalletsView from '@/views/WalletsView.vue'
import ProvidersView from '@/views/ProvidersView.vue'
import { getToken } from '@/lib/api'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'login', component: LoginView },
    { path: '/app', redirect: '/app/dashboard' },
    {
      path: '/app/dashboard',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true },
    },
    {
      path: '/app/users',
      name: 'users',
      component: UsersView,
      meta: { requiresAuth: true },
    },
    {
      path: '/app/wallets',
      name: 'wallets',
      component: WalletsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/app/providers',
      name: 'providers',
      component: ProvidersView,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to) => {
  const authed = Boolean(getToken())
  if (to.meta.requiresAuth && !authed) {
    return { path: '/' }
  }
  if (to.path === '/' && authed) {
    return { path: '/app/dashboard' }
  }
})

export default router
