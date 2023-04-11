import { createRouter, createWebHistory, type RouteLocationNormalized } from 'vue-router'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginView
    },
    {
      path: '/authorize',
      name: 'authorize',
      component: () => import('../views/AuthorizeView.vue')
    },
    {
      path: '/loggedin',
      name: 'loggedin',
      component: () => import('../views/LoggedInView.vue')
    },
  ]
})

router.beforeEach((to, from, next) => {
  if (!hasQueryParams(to) && hasQueryParams(from)) {
    next({ path: to.path, query: from.query, name: to.name ? to.name : undefined });
  } else {
    next()
  }
})


function hasQueryParams(route: RouteLocationNormalized) {
  return !!Object.keys(route.query).length
}


export default router
