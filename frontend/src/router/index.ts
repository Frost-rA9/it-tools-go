import { createRouter, createWebHashHistory } from 'vue-router'
import BaseLayout from '../layouts/BaseLayout.vue'
import ToolLayout from '../layouts/ToolLayout.vue'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      component: BaseLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomeView,
        },
        {
          path: 'tool/:id',
          component: ToolLayout,
          children: [
            {
              path: '',
              name: 'tool',
              component: () => import('../views/ToolView.vue'),
            },
          ],
        },
      ],
    },
  ],
})

export default router
