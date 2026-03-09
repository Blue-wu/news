import axios from 'axios';

// 创建axios实例
const api = axios.create({
  // 在开发环境中，基础URL将通过Vite代理转发到后端
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器 - 添加token
api.interceptors.request.use(
  config => {
    // 从localStorage获取token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器 - 处理401错误
api.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    // 如果是401错误，只对需要认证的API路径进行重定向
    // 公开路径如文章列表、文章详情、分类等不应触发重定向
    if (error.response && error.response.status === 401) {
      // 检查请求的URL是否包含需要认证的路径
      const url = error.config.url || '';
      const requiresAuthPaths = ['/admin', '/user', '/auth/me'];
      const isAuthRequired = requiresAuthPaths.some(path => url.includes(path));
      
      if (isAuthRequired) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
      } else {
        // 对于公开API的401错误，仅打印错误，不重定向
        console.warn('公开API返回401错误，但不进行重定向:', url);
      }
    }
    console.error('API请求错误:', error);
    return Promise.reject(error);
  }
);

// 文章相关API
export const articleApi = {
  // 获取文章列表，添加分页参数
  getArticles: (page = 1, pageSize = 10, additionalParams = {}) => api.get('/articles', {
    params: { page, page_size: pageSize, ...additionalParams }
  }),
  
  // 获取分类下的文章，添加分页参数
  getArticlesByCategory: (slug, page = 1, pageSize = 10) => api.get(`/categories/${slug}/articles`, {
    params: { page, page_size: pageSize }
  }),
   // 获取随机文章
  getRandomArticles: (count = 6) => api.get('/articles/random', {
    params: { limit: count }
  }),
  	// 获取热门文章
	getPopularArticles: (count = 8) => api.get('/articles/popular', {
		params: { limit: count }
	}),
  // 获取单篇文章
  getArticle: (id) => api.get(`/articles/${id}`),
    // 添加增加阅读计数的方法
  incrementViewCount: (id) => api.post(`/articles/${id}/view`),
  // 搜索文章
	searchArticles: (keyword, page = 1, pageSize = 10) => api.get('/articles/search', {
  params: { q: keyword, page, page_size: pageSize }
  
}),
  // 创建文章
  createArticle: (articleData) => api.post('/articles', articleData),
  
  // 更新文章
  updateArticle: (id, articleData) => api.put(`/articles/${id}`, articleData),
  
  // 删除文章
  deleteArticle: (id) => api.delete(`/articles/${id}`)
};

// 添加认证相关API
export const authApi = {
  // 用户登录
  login: (credentials) => api.post('/auth/login', credentials),
  
  // 用户注册
  register: (userData) => api.post('/auth/register', userData),
  
  // 获取用户信息
  getProfile: () => api.get('/user/profile')
};

// 添加管理员API
// 添加管理员API
export const adminApi = {
  // 用户管理
  getUsers: () => api.get('/admin/users'),
  updateUserStatus: (id, status) => api.put(`/admin/users/${id}/status`, { status }),
  deleteUser: (id) => api.delete(`/admin/users/${id}`),
  
  // 文章管理
  getDraftArticles: () => api.get('/admin/articles/draft'),
  updateArticleStatus: (id, status) => api.put(`/admin/articles/${id}/status`, { status }),
   // 添加获取文章详情的方法
  getArticleDetail: (id) => api.get(`/admin/articles/${id}`),
  // 分类管理
  getCategories: () => api.get('/admin/categories'),
  createCategory: (category) => api.post('/admin/categories', category),
  updateCategory: (id, category) => api.put(`/admin/categories/${id}`, category),
  deleteCategory: (id) => api.delete(`/admin/categories/${id}`),
  // 系统配置管理
  getSystemSettings: () => api.get('/admin/settings'),
  updateSystemSettings: (settings) => api.put('/admin/settings', settings)
};

// 添加站点地图相关API
export const sitemapAPI = {
  // 获取站点地图设置
  getSettings() {
    return api.get('/admin/settings/sitemap');
  },
  
  // 更新站点地图设置
  updateSettings(settings) {
    return api.put('/admin/settings/sitemap', settings);
  },
  
  // 重新生成站点地图
  regenerate() {
    return api.post('/admin/sitemap/generate');
  },
  
  // 清除站点地图缓存
  clearCache() {
    return api.post('/admin/sitemap/clear-cache');
  },
  
  // 获取站点地图统计
  getStats() {
    return api.get('/admin/sitemap/stats');
  }
};
// 导出
export default api;