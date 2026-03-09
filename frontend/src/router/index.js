import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/login/Login.vue'
import HomeLayout from '../views/HomeLayout.vue';
import SearchResults from '../views/SearchResults.vue';
const routes = [
  // 根路径路由，确保直接匹配'/'
  {
    path: '/',
    name: 'Home',
    component: HomeLayout,
    meta: { title: '文章列表', requiresAuth: false },
    // 这里不需要嵌套路由，直接在HomeLayout中处理内容显示
  },
  
  // 其他前台页面路由
  {
    path: '/articles',
    name: 'articleList',
    component: HomeLayout,
    meta: { title: '文章列表', requiresAuth: false }
  },
  {
    path: '/articles/:id',
    name: 'articleDetail',
    component: () => import('../views/ArticleDetail.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/category/:slug',
    name: 'categoryArticles',
    component: HomeLayout,
    props: true,
    meta: { requiresAuth: false }
  },
  
  // 管理员页面使用AdminLayout布局
  {
    path: '/admin',
    component: () => import('../views/admin/Layout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'adminDashboard',
        component: () => import('../views/admin/Dashboard.vue')
      },
      {
        path: 'articles',
        name: 'adminArticleManagement',
        component: () => import('../views/admin/ArticleManagement.vue')
      },
      {
        path: 'articles/edit/:id?',
        name: 'adminArticleEdit',
        component: () => import('../views/ArticleEdit.vue')
      },
      {
        path: 'users',
        name: 'userManagement',
        component: () => import('../views/admin/UserManagement.vue')
      },
      {
        path: 'settings',
        name: 'adminSettings',
        component: () => import('../views/admin/Setting.vue')
      },
      {
        path: 'categories',
        name: 'categoryManagement',
        component: () => import('../views/admin/CategoryManagement.vue')
      },
       // 站点地图管理路由
      {
        path: 'sitemap',
        name: 'SitemapManagement',
        component: () => import('../views/admin/SitemapManagement.vue'),
        meta: {
          title: '站点地图管理'
        }
      }
    ]
  },
  
  // 登录页面不使用布局 - 移到最后，确保默认不会优先匹配
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const userStr = localStorage.getItem('user');
  let user = null;
  
  if (userStr) {
    user = JSON.parse(userStr);
  }
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - 博客系统` : '博客系统'
  
  // 添加详细调试日志，帮助跟踪路由跳转行为
  console.log('===== 路由跳转开始 =====');
  console.log('访问路径:', to.path);
  console.log('需要认证:', to.meta.requiresAuth);
  console.log('是否已登录:', !!token);
  console.log('从哪个路径跳转:', from.path);
  
  // 特殊处理根路径，确保直接允许访问
  if (to.path === '/') {
    console.log('访问根路径，直接允许访问');
    return next();
  }
  
  // 只有当requiresAuth严格等于true时才检查登录状态
  if (to.meta.requiresAuth === true) {
    if (!token) {
      console.log('未登录，重定向到登录页');
      return next('/login');
    }
    
    // 检查是否需要管理员权限
    if (to.meta.requiresAdmin && user && user.role !== 'admin') {
      alert('您没有管理员权限');
      return next('/');
    }
  } else {
    console.log('允许访问，无需登录验证');
  }
  
  // 如果已登录访问登录页，重定向到首页
  if (to.path === '/login' && token) {
    console.log('已登录，重定向到首页');
    return next('/');
  }
  
  console.log('允许继续访问当前路径');
  next();
});

export default router