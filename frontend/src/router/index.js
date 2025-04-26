import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('../layouts/DefaultLayout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/Home.vue')
      },
      {
        path: 'novels',
        name: 'Novels',
        component: () => import('../views/NovelSquare.vue')
      },
      {
        path: 'novels/:id',
        name: 'Novel',
        component: () => import('../views/Novel.vue')
      },
      {
        path: 'novels/:id/chapter/:chapterId',
        name: 'Chapter',
        component: () => import('../views/Chapter.vue')
      },
      {
        path: 'write',
        name: 'Write',
        component: () => import('../views/WriteCenter.vue')
      },
      {
        path: 'login',
        name: 'Login',
        component: () => import('../views/Login.vue')
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('../views/Register.vue')
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('../views/UserProfile.vue'),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/write',
    name: 'WriteCenter',
    component: () => import('../views/WriteCenter.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/write/editor/:id',
    name: 'NovelEditor',
    component: () => import('../views/write/NovelEditor.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/write/chapters/:id',
    name: 'ChapterManager',
    component: () => import('../views/write/ChapterManager.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/write/settings/:id',
    name: 'NovelSettings',
    component: () => import('../views/write/NovelSettings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/write/outline/:id',
    name: 'NovelOutline',
    component: () => import('../views/write/NovelOutline.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/write/stats/:id',
    name: 'NovelStats',
    component: () => import('../views/write/NovelStats.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/library',
    name: 'Library',
    component: () => import('../views/Library.vue'),
    meta: {
      title: '书库'
    }
  },
  {
    path: '/ranking',
    name: 'Ranking',
    component: () => import('../views/Ranking.vue'),
    meta: {
      title: '排行榜'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router 