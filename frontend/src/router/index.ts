import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    redirect: '/dashboard/tasks',
    children: [
      {
        path: 'tasks',
        name: 'TaskManagement',
        component: () => import('../views/TaskManagement.vue')
      },
      {
        path: 'evidences',
        name: 'EvidenceManagement',
        component: () => import('../views/EvidenceManagement.vue')
      },
      {
        path: 'results',
        name: 'ResultManagement',
        component: () => import('../views/ResultManagement.vue')
      },
      {
        path: 'reports',
        name: 'ReportManagement',
        component: () => import('../views/ReportManagement.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && to.path !== '/register' && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router