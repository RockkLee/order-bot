import { createRouter, createWebHistory } from 'vue-router'

import CSide from '../views/CSide.vue'
import BLogin from '../views/BLogin.vue'
import BSignup from '../views/BSignup.vue'
import BDashboard from '../views/BDashboard.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/c',
    },
    {
      path: '/c',
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
  ],
})

function isLoggedIn(): boolean {
  return Boolean(localStorage.getItem("access_token"));
}

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !isLoggedIn()) {
    return {
      name: "b-login", // name of a router
      query: { redirect: to.fullPath }, // so you can return after login
    };
  }
  return true;
});

export default router
