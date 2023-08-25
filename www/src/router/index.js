import { createRouter, createWebHistory } from 'vue-router'
import IndexView from '../views/IndexView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: IndexView
    },
    // {
    //   path: '/about',
    //   name: 'about',
    //   // route level code-splitting
    //   // this generates a separate chunk (About.[hash].js) for this route
    //   // which is lazy-loaded when the route is visited.
    //   component: () => import('../views/AboutView.vue')
    // }
    {
      path: '/admin',
      name: "admin",
      component: () => import("../views/AdminView.vue"),
    },
    {
      path: '/streamer/data/member',
      name: "member",
      component: () => import("../views/StreamerMemberView.vue"),
    },
    {
      path: '/streamer/data/sc',
      name: "sc",
      component: () => import("../views/StreamerSuperchatView.vue"),
    },
    {
      path: '/streamer/data/gift',
      name: "gift",
      component: () => import("../views/StreamerGiftView.vue"),
    },
    {
      path: '/streamer/events',
      name: "events",
      component: () => import("../views/StreamerEventListView.vue"),
    },
    {
      path: '/streamer/events/new',
      name: "events_add",
      component: () => import("../views/StreamerEventAddView.vue"),
    },
    {
      path: '/streamer/event/:id',
      name: "event_detail",
      component: () => import("../views/StreamerEventDetailView.vue"),
    },
  ]
})

export default router
