import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '平台登录' }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
    // 添加路由元信息，标记这个页面需要登录才能进入
    meta: { requiresAuth: true, title: '数据大屏' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫：阻止未登录用户强行访问系统内部页面
router.beforeEach((to, from, next) => {
  // 设置浏览器标题
  if (to.meta.title) {
    document.title = to.meta.title + ' - 杭州餐饮数据分析'
  }

  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    // 需要鉴权但没 token，强制跳转登录页
    next('/login')
  } else if (to.path === '/login' && token) {
    // 已经登录的用户如果还要去登录页，直接带入大屏
    next('/dashboard')
  } else {
    next()
  }
})

export default router
