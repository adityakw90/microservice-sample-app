import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import UsersView from '@/views/UsersView.vue'
import LoginView from '@/views/LoginView.vue'
import ProfileView from '@/views/ProfileView.vue'
import UserDetailView from '@/views/UserDetailView.vue'
import ProtoTestView from '@/views/ProtoTestView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { public: true }
    },
    {
      path: '/',
      component: DashboardLayout,
      children: [
        {
          path: '',
          redirect: '/dashboard'
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: DashboardView
        },
        {
          path: 'users',
          name: 'users',
          component: UsersView
        },
        {
          path: 'users/:uid',
          name: 'user-detail',
          component: UserDetailView
        },
        {
          path: 'profile',
          name: 'profile',
          component: ProfileView,
          meta: { requiresAuth: true }
        },
        {
          path: 'activity',
          name: 'activity',
          component: DashboardView
        },
        {
          path: 'security',
          name: 'security',
          component: DashboardView
        },
        {
          path: 'settings',
          name: 'settings',
          component: DashboardView
        },
        {
          path: 'proto-test',
          name: 'proto-test',
          component: ProtoTestView
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
