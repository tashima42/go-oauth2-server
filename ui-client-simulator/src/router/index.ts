import { createRouter, createWebHistory } from 'vue-router'
import ClientSimulatorView from "@/views/ClientSimulatorView.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: ClientSimulatorView
    },
  ]
})

export default router
