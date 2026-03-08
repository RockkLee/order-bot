import { createRouter, createWebHistory } from 'vue-router'

import CSide from '../views/CSide.vue'
import BLogin from '../views/BLogin.vue'
import BSignup from '../views/BSignup.vue'
import BDashboard from '../views/BDashboard.vue'
import ClosedNotice from '../views/ClosedNotice.vue'

import { isBusinessOpenUtc8 } from '../utils/businessHours'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/c/:botId/:menuId',
      name: 'c-side',
      component: CSide,
    },
    {
      path: '/b',
      redirect: '/b/login',
    },
    {
      path: '/b/login',
      name: 'b-login',
      component: BLogin,
    },
    {
      path: '/b/signup',
      name: 'b-signup',
      component: BSignup,
    },
    {
      path: '/b/app',
      name: 'b-dashboard',
      component: BDashboard,
      meta: { requiresAuth: true },
    },
    {
      path: '/closed',
      name: 'closed',
      component: ClosedNotice,
    },
  ],
})

function isLoggedIn(): boolean {
  return Boolean(localStorage.getItem('access_token'))
}

router.beforeEach((to) => {
  if (!isBusinessOpenUtc8() && to.name !== 'closed') {
    return {
      name: 'closed',
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.requiresAuth && !isLoggedIn()) {
    return {
      name: 'b-login', // name of a router
      query: { redirect: to.fullPath }, // so you can return after login
    }
  }
  return true
})

export default router
