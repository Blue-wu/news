<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import { articleApi } from '@/utils/api.js'; // 导入API

// 原有状态
const user = ref(null);
const router = useRouter();
const route = useRoute();
const currentPath = ref('首页');
const categories = ref([]);
const loadingCategories = ref(false);
const searchKeyword = ref('');

// 文章列表相关状态
const articles = ref([]);
const loading = ref(false);
const error = ref('');
const currentPage = ref(1);
const pageSize = ref(10);
const totalArticles = ref(0);
const totalPages = ref(0);

// 新增搜索结果相关状态
const showSearchResults = ref(false);
const searchArticles = ref([]);
const searchLoading = ref(false);
const searchError = ref('');
const searchCurrentPage = ref(1);
const searchPageSize = ref(5); // 小页面显示更少结果
const searchTotalCount = ref(0);
const searchTotalPages = ref(0);
// 新增随便看看相关状态
const randomArticles = ref([]);
const randomLoading = ref(false);
const randomError = ref('');

const popularArticles = ref([]) // 热门文章数据
const popularLoading = ref(false)
const popularError = ref('')

// 获取主要文章列表
const fetchArticles = async (page = 1, additionalParams = {}) => {
  loading.value = true;
  error.value = '';
  try {
    const response = await articleApi.getArticles(page, pageSize.value, additionalParams);
    
    // 处理API响应，考虑不同的响应格式
    if (response && response.data && Array.isArray(response.data)) {
      articles.value = response.data;
      // 尝试从不同字段获取总数
      totalArticles.value = response.total || response.totalCount || response.count || 0;
    } else if (Array.isArray(response)) {
      // 如果响应直接是文章数组
      articles.value = response;
      totalArticles.value = response.length;
    } else {
      throw new Error('Invalid API response format');
    }
    
    // 计算总页数
    totalPages.value = Math.ceil(totalArticles.value / pageSize.value);
    currentPage.value = page;
    
    console.log('文章列表数据:', articles.value);
    console.log('总文章数:', totalArticles.value);
    console.log('总页数:', totalPages.value);
  } catch (err) {
    error.value = '获取文章列表失败，请稍后重试';
    console.error('Error fetching articles:', err);
    // 提供一些默认文章数据
    articles.value = [
      {
        id: 1,
        title: '示例文章1：欢迎使用我们的网站',
        content: '这是一篇示例文章的内容...',
        created_at: new Date().toISOString(),
        user: { username: '管理员' },
        view_count: 123,
        cover_image: 'https://picsum.photos/800/400?random=1'
      },
      {
        id: 2,
        title: '示例文章2：网站使用指南',
        content: '本文将指导您如何使用我们的网站...',
        created_at: new Date(Date.now() - 86400000).toISOString(),
        user: { username: '编辑' },
        view_count: 45,
        cover_image: 'https://picsum.photos/800/400?random=2'
      }
    ];
    totalArticles.value = articles.value.length;
    totalPages.value = 1;
  } finally {
    loading.value = false;
  }
};

// 分页方法
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    fetchArticles(page);
  }
};

const goToPrevPage = () => {
  if (currentPage.value > 1) {
    goToPage(currentPage.value - 1);
  }
};

const goToNextPage = () => {
  if (currentPage.value < totalPages.value) {
    goToPage(currentPage.value + 1);
  }
};

// 获取热门文章
const fetchPopularArticles = async () => {
	popularLoading.value = true
	popularError.value = ''
	try {
		const response = await articleApi.getPopularArticles(8)
		popularArticles.value = Array.isArray(response) ? response : []
	} catch (err) {
		popularError.value = '获取热门文章失败'
		console.error('Error fetching popular articles:', err)
		// 提供默认数据
		popularArticles.value = [
			{ id: 1, title: '热门文章1', view_count: 1234 },
			{ id: 2, title: '热门文章2', view_count: 987 },
			{ id: 3, title: '热门文章3', view_count: 567 },
			{ id: 4, title: '热门文章4', view_count: 456 },
			{ id: 5, title: '热门文章5', view_count: 345 },
			{ id: 6, title: '热门文章6', view_count: 234 },
			{ id: 7, title: '热门文章7', view_count: 123 },
			{ id: 8, title: '热门文章8', view_count: 99 }
		]
	} finally {
		popularLoading.value = false
	}
}
function formatDate(dateString) {
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) {
    return '今天';
  } else if (diffDays === 1) {
    return '昨天';
  } else if (diffDays < 30) {
    return `${diffDays}天前`;
  } else {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
  }
}

// 获取随机文章
const fetchRandomArticles = async () => {
  randomLoading.value = true;
  randomError.value = '';
  try {
    const response = await articleApi.getRandomArticles(6);
    // 假设响应格式与其他文章API一致
    randomArticles.value = response || [];
  } catch (error) {
    randomError.value = '获取随机文章失败';
    console.error('获取随机文章失败:', error);
    // 错误时使用默认数据
    randomArticles.value = [];
  } finally {
    randomLoading.value = false;
  }
};

// 刷新随机文章（"换一换"功能）
const refreshRandomArticles = () => {
  fetchRandomArticles();
};

// 跳转到文章详情
const goToArticleDetail = (articleId) => {
  closeSearchResults();
  router.push(`/articles/${articleId}`);
};

// 搜索结果获取方法
const fetchSearchResults = async (q, page = 1) => {
  if (!q) return;
  
  searchLoading.value = true;
  searchError.value = '';
  try {
    const response = await articleApi.searchArticles(q, page, searchPageSize.value);
    
    searchArticles.value = response.data || [];
    searchTotalCount.value = response.total || 0;
    searchTotalPages.value = Math.ceil(searchTotalCount.value / searchPageSize.value);
  } catch (err) {
    searchError.value = '搜索失败，请稍后重试';
    console.error('Search error:', err);
  } finally {
    searchLoading.value = false;
  }
};

// 修改处理搜索逻辑，不再跳转页面
const handleSearch = () => {
  const keyword = searchKeyword.value.trim();
  if (keyword) {
    searchCurrentPage.value = 1;
    showSearchResults.value = true;
    fetchSearchResults(keyword, 1);
  }
};

// 处理回车键搜索
const handleKeyPress = (event) => {
  if (event.key === 'Enter') {
    handleSearch();
  }
};

// 关闭搜索结果
const closeSearchResults = () => {
  showSearchResults.value = false;
  searchArticles.value = [];
};

// 点击其他地方关闭搜索结果
const handleClickOutside = (event) => {
  const searchBox = document.querySelector('.search-result-container');
  const searchInput = document.querySelector('.search-box input');
  const searchBtn = document.querySelector('.search-btn');
  
  if (searchBox && !searchBox.contains(event.target) && 
      event.target !== searchInput && event.target !== searchBtn) {
    closeSearchResults();
  }
};

// 分页方法
const goToSearchPage = (page) => {
  if (page >= 1 && page <= searchTotalPages.value) {
    searchCurrentPage.value = page;
    fetchSearchResults(searchKeyword.value.trim(), page);
  }
};

const goToSearchPrevPage = () => {
  if (searchCurrentPage.value > 1) {
    goToSearchPage(searchCurrentPage.value - 1);
  }
};

const goToSearchNextPage = () => {
  if (searchCurrentPage.value < searchTotalPages.value) {
    goToSearchPage(searchCurrentPage.value + 1);
  }
};

// 根据当前路由获取文章数据
const fetchArticlesByRoute = async () => {
  // 重置页码到第一页
  currentPage.value = 1;
  
  // 获取当前路径
  const path = window.location.pathname;
  
  // 检查是否是分类页面
  const categoryMatch = path.match(/^\/category\/([^/]+)$/);
  
  if (categoryMatch) {
    const slug = categoryMatch[1];
    console.log('获取分类文章，slug:', slug);
    fetchCategoryArticles(slug);
  } else {
    // 否则获取全部文章
    console.log('获取全部文章');
    fetchArticles(1);
  }
};

// 获取分类下的文章
const fetchCategoryArticles = async (slug, page = 1) => {
  loading.value = true;
  error.value = '';
  try {
    const response = await articleApi.getArticlesByCategory(slug, page, pageSize.value);
    
    // 处理API响应，考虑不同的响应格式
    if (response && response.data && Array.isArray(response.data)) {
      articles.value = response.data;
      totalArticles.value = response.total || response.totalCount || response.count || 0;
    } else if (Array.isArray(response)) {
      // 如果响应直接是文章数组
      articles.value = response;
      totalArticles.value = response.length;
    } else {
      throw new Error('Invalid API response format');
    }
    
    // 计算总页数
    totalPages.value = Math.ceil(totalArticles.value / pageSize.value);
    currentPage.value = page;
    
    console.log('分类文章列表数据:', articles.value);
    console.log('总文章数:', totalArticles.value);
    console.log('总页数:', totalPages.value);
  } catch (err) {
    error.value = '获取分类文章列表失败，请稍后重试';
    console.error('Error fetching category articles:', err);
    // 提供一些默认文章数据
    articles.value = [];
    totalArticles.value = 0;
    totalPages.value = 0;
  } finally {
    loading.value = false;
  }
};

// 监听路由变化
watch(() => route.path, () => {
  console.log('路由变化，重新加载数据，新路径:', route.path);
  updateCurrentPath();
  fetchArticlesByRoute();
});

onMounted(() => {
  // 原有初始化逻辑
  const userStr = localStorage.getItem('user');
  if (userStr) {
    user.value = JSON.parse(userStr);
  }
  
  updateCurrentPath();
  fetchCategories();
  fetchRandomArticles(); // 添加这行，初始化随机文章
  fetchPopularArticles(); // 添加热门文章获取
  fetchArticlesByRoute(); // 使用新函数替代直接调用fetchArticles
  
  // 添加点击事件监听器
  document.addEventListener('click', handleClickOutside);
});

// 清理事件监听器
const cleanup = () => {
  document.removeEventListener('click', handleClickOutside);
};

// 原有方法保持不变
const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  user.value = null;
  router.push('/login');
};

const updateCurrentPath = () => {
  const path = window.location.pathname;
  if (path === '/') {
    currentPath.value = [{ text: '首页', url: '/' }];
  } else if (path === '/articles') {
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: '文章列表', url: '/articles' }
    ];
  } else if (path.match(/^\/articles\/\d+$/)) {
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: '文章列表', url: '/articles' },
      { text: '文章详情', url: path }
    ];
  } else if (path.includes('/login')) {
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: '登录', url: '/login' }
    ];
  } else if (path.includes('/golden-shovel')) {
    if (path.includes('/lineup')) {
      currentPath.value = [
        { text: '首页', url: '/' },
        { text: '金铲铲阵容', url: '/golden-shovel/lineup' }
      ];
    } else if (path.includes('/heroes')) {
      currentPath.value = [
        { text: '首页', url: '/' },
        { text: '金铲铲海克斯', url: '/golden-shovel/heroes' }
      ];
    } else if (path.includes('/synergies')) {
      currentPath.value = [
        { text: '首页', url: '/' },
        { text: '金铲铲羁绊', url: '/golden-shovel/synergies' }
      ];
    }
  } else if (path.includes('/strategies')) {
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: '攻略', url: '/strategies' }
    ];
  } else if (path.match(/^\/category\/[^/]+$/)) {
    const categoryName = decodeURIComponent(path.split('/').pop());
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: categoryName, url: path }
    ];
  } else {
    currentPath.value = [
      { text: '首页', url: '/' },
      { text: '关于软件的文章', url: path }
    ];
  }
};

const fetchCategories = async () => {
  loadingCategories.value = true;
  try {
    const response = await axios.get('/api/categories');
    
    if (response.data && response.data.success && Array.isArray(response.data.data)) {
      categories.value = response.data.data.filter(cat => cat.parent_id === 0);
    } else {
      console.warn('分类数据格式不正确:', response.data);
      categories.value = [];
    }
  } catch (error) {
    console.error('获取分类失败:', error);
    categories.value = [
      { id: 1, name: '金铲铲阵容', slug: 'golden-shovel-lineup', parent_id: 0 },
      { id: 2, name: '金铲铲海克斯', slug: 'golden-shovel-heroes', parent_id: 0 },
      { id: 3, name: '金铲铲羁绊', slug: 'golden-shovel-synergies', parent_id: 0 },
      { id: 4, name: '攻略', slug: 'strategies', parent_id: 0 }
    ];
  } finally {
    loadingCategories.value = false;
  }
};
</script>

<template>
  <div class="home-layout">
    <!-- 顶部区域 -->
    <header class="header">
      <div class="header-content">
        <!-- 移动设备登录按钮占位符 -->
        <div class="mobile-login-placeholder"></div>
        <div class="logo">
          <router-link to="/">优品飞游戏</router-link>
        </div>
        
        <!-- 搜索区域 -->
        <div class="search-container">
          <div class="search-box">
            <input 
              type="text" 
              placeholder="搜索感兴趣的知识和文章" 
              v-model="searchKeyword"
              @keyup.enter="handleKeyPress"
              @focus="showSearchResults = searchArticles.length > 0"
            />
            <button class="search-btn" @click="handleSearch">搜索一下</button>
          </div>
          
          <!-- 搜索结果下拉面板 -->
          <div v-if="showSearchResults" class="search-result-container">
            <div class="search-result-header">
              <h3>搜索结果: "{{ searchKeyword }}"</h3>
              <button class="close-btn" @click="closeSearchResults">×</button>
            </div>
            
            <div class="search-result-content">
              <div v-if="searchLoading" class="search-loading">加载中...</div>
              <div v-else-if="searchError" class="search-error">{{ searchError }}</div>
              <div v-else-if="searchArticles.length === 0" class="search-empty">
                没有找到相关内容，请尝试其他关键词
              </div>
              <div v-else class="search-article-list">
                <div 
                  v-for="article in searchArticles" 
                  :key="article.id" 
                  class="search-article-item"
                  @click="goToArticleDetail(article.id)"
                >
                  <div class="article-title">{{ article.title }}</div>
                  <div class="article-meta">
                    <span>{{ article.user?.username || '未知用户' }}</span>
                    <span class="dot">·</span>
                    <span>{{ new Date(article.created_at).toLocaleDateString() }}</span>
                  </div>
                </div>
              </div>
            </div>
            
            <!-- 分页控制 -->
            <div v-if="searchTotalPages > 1" class="search-pagination">
              <button 
                @click="goToSearchPrevPage" 
                :disabled="searchCurrentPage === 1" 
                class="page-btn"
              >
                上一页
              </button>
              <span class="page-info">
                第 {{ searchCurrentPage }} / {{ searchTotalPages }} 页
              </span>
              <button 
                @click="goToSearchNextPage" 
                :disabled="searchCurrentPage === searchTotalPages" 
                class="page-btn"
              >
                下一页
              </button>
            </div>
          </div>
        </div>
        
        <!-- 用户操作区域 - 移动设备上会移到右上角 -->
        <div class="user-actions">
          <template v-if="user">
            <span class="welcome-text">欢迎，{{ user.username }}</span>
            <button @click="handleLogout" class="logout-btn">登出</button>
          </template>
          <template v-else>
            <router-link to="/login" class="login-register-btn mobile-login-btn">注册/登录</router-link>
          </template>
        </div>
      </div>
    </header>

    <!-- 导航栏 -->
    <nav class="main-nav">
      <div class="nav-content">
        <ul class="nav-links">
          <li><router-link to="/" class="nav-link">首页</router-link></li>
          <li v-for="category in categories" :key="category.id">
            <router-link 
              :to="category.slug && category.slug.includes('golden-shovel') 
                  ? `/golden-shovel/${category.slug.replace('golden-shovel-', '')}` 
                  : `/category/${category.slug || 'default'}`" 
              class="nav-link"
            >
              {{ category.name || '未命名分类' }}
            </router-link>
          </li>
        </ul>
      </div>
    </nav>

    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <span>当前位置: </span>
      <span class="path">
        <template v-for="(item, index) in currentPath" :key="index">
          <template v-if="index === currentPath.length - 1">
            <span>{{ item.text }}</span>
          </template>
          <template v-else>
            <router-link :to="item.url" class="breadcrumb-link">{{ item.text }}</router-link>
            <span class="breadcrumb-separator"> &gt; </span>
          </template>
        </template>
      </span>
    </div>

    <!-- 主内容区域 -->
    <div class="content-wrapper">
      <main class="main-content">
        <!-- 文章列表内容 -->
        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p>加载文章中...</p>
        </div>
        
        <div v-else-if="error" class="error-container">
          <p class="error-message">{{ error }}</p>
        </div>
        
        <div v-else-if="articles.length === 0" class="empty-container">
          <p>暂无文章</p>
        </div>
        
        <div v-else class="article-list">
          <div 
            v-for="article in articles" 
            :key="article.id" 
            class="article-item"
            @click="goToArticleDetail(article.id)"
          >
            <!-- 文章封面图 -->
            <div class="article-cover" v-if="article.cover_image">
              <img :src="article.cover_image" :alt="article.title" class="cover-image" />
            </div>
            
            <!-- 文章内容 -->
            <div class="article-content">
              <h2 class="article-title">{{ article.title }}</h2>
              <div class="article-meta">
                <span class="author">{{ article.user?.username || '未知用户' }}</span>
                <span class="date">{{ formatDate(article.created_at) }}</span>
                <span class="views">{{ article.view_count || 0 }} 阅读</span>
              </div>
              <div class="article-excerpt">
                {{ article.content ? article.content.replace(/<[^>]*>/g, '').substring(0, 150) + '...' : '暂无内容' }}
              </div>
              <div class="article-actions">
                <button class="read-more-btn">阅读全文</button>
              </div>
            </div>
          </div>
          
          <!-- 分页控制 -->
          <div class="pagination" v-if="totalPages > 1">
            <button 
              @click="goToPrevPage" 
              :disabled="currentPage === 1" 
              class="page-btn"
            >
              上一页
            </button>
            <span class="page-info">
              第 {{ currentPage }} / {{ totalPages }} 页
            </span>
            <button 
              @click="goToNextPage" 
              :disabled="currentPage === totalPages" 
              class="page-btn"
            >
              下一页
            </button>
          </div>
        </div>
      </main>
      
      <aside class="sidebar">
        <!-- 随便看看 -->
        <!-- 随便看看 -->
<div class="sidebar-section">
  <div class="sidebar-title">
    <span>随便看看</span>
    <span class="more-link" @click="refreshRandomArticles">换一换</span>
  </div>
  <div class="sidebar-content">
    <div v-if="randomLoading" class="loading">加载中...</div>
    <div v-else-if="randomError" class="error-message">{{ randomError }}</div>
    <div v-else class="random-articles-grid">
      <div 
        v-for="article in randomArticles" 
        :key="article.id"
        class="random-article-item"
        @click="router.push(`/articles/${article.id}`)"
      >
        <img 
          :src="article.cover_image || `https://picsum.photos/200/150?random=${article.id}`" 
          alt="文章图片" 
          class="random-article-img"
        >
        <div class="random-article-title">{{ article.title }}</div>
      </div>
    </div>
  </div>
</div>
        
        <!-- 热门文章 -->
        <!-- 热门文章 -->
<div class="sidebar-section">
  <h3 class="sidebar-title">
    <span>热门文章</span>
    <div class="title-dots">
      <span class="dot red"></span>
      <span class="dot yellow"></span>
      <span class="dot blue"></span>
    </div>
  </h3>
  <div class="sidebar-content">
    <div v-if="popularLoading" class="loading">加载中...</div>
    <div v-else-if="popularError" class="error-message">{{ popularError }}</div>
    <div v-else-if="popularArticles.length > 0" class="hot-articles">
      <router-link 
        v-for="(article, index) in popularArticles" 
        :key="article.id"
        :to="`/articles/${article.id}`"
        class="hot-article-item"
        target="_blank"
      >
        <span 
          class="article-rank" 
          :class="{
            'rank-1': index === 0,
            'rank-2': index === 1,
            'rank-3': index === 2
          }"
        >
          {{ index + 1 }}
        </span>
        <img 
          :src="article.cover_image || `https://picsum.photos/80/60?random=${article.id}`" 
          alt="文章缩略图" 
          class="article-thumb"
        >
        <div class="article-info">
          <div class="article-title">{{ article.title }}</div>
          <div class="article-meta">
            <span class="article-date">{{ formatDate(article.created_at) }}</span>
            <span class="meta-separator">|</span>
            <span class="article-views">{{ article.view_count || 0 }} 阅读</span>
          </div>
        </div>
      </router-link>
    </div>
    <div v-else class="no-data">暂无热门文章</div>
  </div>
</div>
      </aside>
    </div>
  </div>

  <!-- 侧边悬浮关注功能 -->
<div class="floating-sidebar">
  <div class="floating-trigger">
    <span>关注我们</span>
    <div class="floating-content">
      <div class="qrcode-wrapper">
        <img src="https://picsum.photos/200/200?random=qrcode" alt="公众号二维码" class="qrcode-image">
        <div class="qrcode-text">微信扫一扫关注我们</div>
        <button class="follow-button">关注微信</button>
        <div class="social-share">
          <div class="social-icon wechat">微信</div>
          <div class="social-icon weibo">微博</div>
          <div class="social-icon github">GitHub</div>
        </div>
      </div>
    </div>
  </div>
</div>
</template>

<style scoped>
/* 搜索相关样式 */
.search-container {
  position: relative;
  flex: 1;
  max-width: 500px;
  margin: 0 20px;
}

.search-result-container {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: white;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1000;
  max-height: 400px;
  display: flex;
  flex-direction: column;
}

.search-result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  border-bottom: 1px solid #eee;
  background: #f8f8f8;
}

.search-result-header h3 {
  margin: 0;
  font-size: 14px;
  color: #666;
  font-weight: normal;
}

.close-btn {
  background: none;
  border: none;
  font-size: 18px;
  color: #999;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  color: #333;
}

.search-result-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.search-loading, .search-error, .search-empty {
  padding: 20px;
  text-align: center;
  color: #666;
}

.search-article-list {
  display: flex;
  flex-direction: column;
}

.search-article-item {
  padding: 10px;
  border-bottom: 1px solid #f0f0f0;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-article-item:hover {
  background-color: #f8f8f8;
}

.search-article-item:last-child {
  border-bottom: none;
}

.search-article-item .article-title {
  font-size: 14px;
  color: #333;
  margin-bottom: 5px;
  line-height: 1.4;
}

.search-article-item .article-meta {
  font-size: 12px;
  color: #999;
}

.search-article-item .dot {
  margin: 0 5px;
}

.search-pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 10px;
  border-top: 1px solid #eee;
  background: #f8f8f8;
  gap: 10px;
}

.search-pagination .page-btn {
  padding: 4px 8px;
  border: 1px solid #ddd;
  background: white;
  color: #333;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
}

.search-pagination .page-btn:hover:not(:disabled) {
  background: #f5f5f5;
}

.search-pagination .page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.search-pagination .page-info {
  font-size: 12px;
  color: #666;
}

/* 重置默认样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Arial', 'Microsoft YaHei', sans-serif;
  color: #333;
  background-color: #f5f5f5;
  line-height: 1.6;
}

/* 面包屑导航样式 */
.breadcrumb {
  padding: 10px 20px;
  background-color: #f9f9f9;
  font-size: 14px;
  color: #666;
  border-bottom: 1px solid #eaeaea;
}

.path {
  color: #333;
}

/* 面包屑链接样式 */
.breadcrumb-link {
  color: #0066cc;
  text-decoration: none;
  transition: color 0.2s;
}

.breadcrumb-link:hover {
  color: #ff6600;
  text-decoration: underline;
}

/* 分隔符样式 */
.breadcrumb-separator {
  color: #999;
  margin: 0 4px;
}

/* 主容器样式 */
.home-layout {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  background-color: #fff;
  min-height: 100vh;
}

/* 顶部区域样式 */
.header {
  padding: 15px 20px;
  background-color: #f8f8f8;
  border-bottom: 1px solid #eaeaea;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
}
.logo {
  font-size: 24px;
  font-weight: bold;
  color: #ff6600;
  display: flex;
  align-items: center;
}

.logo::before {
  content: "🎮";
  margin-right: 8px;
  font-size: 28px;
}

.logo a {
  text-decoration: none;
  color: #ff6600;
}

.search-box {
  display: flex;
  width: 100%;
}

.search-box input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
  font-size: 14px;
}

.search-btn {
  background-color: #ff6600;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
}

.user-actions {
  display: flex;
}

.login-register-btn {
  color: #ff6600;
  text-decoration: none;
  font-size: 14px;
}

.logout-btn {
  background: none;
  border: 1px solid #ff6600;
  color: #ff6600;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}

/* 主导航栏样式 */
.main-nav {
  background-color: #fff;
  border-bottom: 2px solid #ff6600;
}

.nav-content {
  padding: 0 20px;
}

.nav-links {
  display: flex;
  list-style: none;
}

.nav-link {
  display: block;
  padding: 15px 20px;
  color: #333;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s;
}

.nav-link:hover {
  background-color: #f8f8f8;
  color: #ff6600;
}

/* 内容区域样式 */
.content-wrapper {
  display: flex;
  padding: 20px;
}

.main-content {
  flex: 1;
  margin-right: 20px;
}

/* 侧边栏样式 */
.sidebar {
  width: 300px;
}

.sidebar-section {
  margin-bottom: 20px;
  border: 1px solid #eaeaea;
  border-radius: 4px;
  overflow: hidden;
}

.sidebar-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #f8f8f8;
  border-bottom: 1px solid #eaeaea;
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.more-link {
  color: #999;
  font-size: 12px;
  cursor: pointer;
}

.more-link:hover {
  color: #ff6600;
}

.title-dots {
  display: flex;
  gap: 5px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #ddd;
}

/* 为热门文章标题的彩色点添加样式 */
.dot.red {
  background-color: #ff6b6b;
}

.dot.yellow {
  background-color: #feca57;
}

.dot.blue {
  background-color: #48dbfb;
}

.dot.active {
  background-color: #ff6600;
}

.sidebar-content {
  padding: 15px;
}

/* 随便看看网格布局 */
.random-articles-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.random-article-item {
  display: flex;
  flex-direction: column;
  cursor: pointer;
}

.random-article-img {
  width: 100%;
  height: 100px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 5px;
}

.random-article-title {
  font-size: 12px;
  color: #333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.random-article-item:hover .random-article-title {
  color: #ff6600;
}

/* 热门文章完整样式 - 移到媒体查询外部 */
.hot-articles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hot-article-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 8px;
  border-radius: 4px;
  transition: background-color 0.2s;
  text-decoration: none;
  margin-bottom: 0; /* 移除默认的margin-bottom */
}

.hot-article-item:hover {
  background-color: #f8f8f8;
}

.article-rank {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  font-size: 12px;
  font-weight: bold;
  color: #666;
  background-color: #eee;
  border-radius: 2px;
  flex-shrink: 0;
}

.article-rank.rank-1 {
  background-color: #ff6b6b;
  color: white;
}

.article-rank.rank-2 {
  background-color: #feca57;
  color: white;
}

.article-rank.rank-3 {
  background-color: #48dbfb;
  color: white;
}

.article-thumb {
  width: 80px;
  height: 60px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.article-title {
  font-size: 14px;
  color: #333;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.2s;
}

.hot-article-item:hover .article-title {
  color: #ff6600;
}

.article-meta {
  font-size: 12px;
  color: #999;
  display: flex;
  align-items: center;
  gap: 6px;
}

.article-date {
  color: #666;
}

.meta-separator {
  color: #ddd;
}

.article-views {
  color: #666;
}

/* 加载和错误状态样式 */
.loading {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}

.error-message {
  text-align: center;
  padding: 20px;
  color: #ff4757;
  font-size: 14px;
}

.no-data {
  text-align: center;
  padding: 20px;
  color: #999;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    margin-top: 20px;
  }
  
  .search-result-container {
    max-width: 100%;
  }
}

/* 移动设备登录按钮样式 */
.mobile-login-btn {
  display: none; /* 默认隐藏 */
}

.mobile-login-placeholder {
  display: none; /* 默认隐藏 */
}

@media (max-width: 768px) {
  .header-content {
    flex-direction: row;
    padding: 0 15px;
  }
  
  .logo a {
    font-size: 18px;
  }
  
  .search-container {
    flex: 1;
    margin: 0 10px;
  }
  
  .search-box input {
    font-size: 14px;
    padding: 8px 12px;
  }
  
  .search-btn {
    padding: 8px 12px;
    font-size: 14px;
  }
  
  /* 移动设备上登录按钮移到右上角 */
  .user-actions {
    display: flex;
    align-items: center;
  }
  
  /* 正常的登录按钮在移动设备上隐藏 */
  .login-register-btn {
    display: none;
  }
  
  /* 移动设备专用登录按钮显示并定位在右上角 */
  .mobile-login-btn {
    display: block;
    position: absolute;
    top: 10px;
    right: 15px;
    z-index: 100;
    background-color: #ff6b6b;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 8px 16px;
    font-size: 14px;
    text-decoration: none;
  }
  
  /* 占位符确保搜索框不被遮挡 */
  .mobile-login-placeholder {
    display: block;
    width: 80px; /* 与登录按钮宽度匹配 */
  }
  
  .welcome-text {
    display: none;
  }
  
  .logout-btn {
    font-size: 14px;
    padding: 6px 12px;
  }
  
  .nav-links {
    flex-wrap: wrap;
  }
  
  .nav-link {
    padding: 10px;
    font-size: 14px;
  }
  
  .random-articles-grid {
    grid-template-columns: 1fr;
  }
  
  /* 移动端热门文章调整 */
  .hot-article-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .article-thumb {
    width: 100%;
    height: 120px;
  }
  
  .article-rank {
    margin-bottom: 8px;
  }
}

/* 针对小屏幕手机的额外优化 */
@media (max-width: 480px) {
  .search-container {
    margin: 0 5px;
  }
  
  .search-box input {
    font-size: 13px;
  }
  
  .search-btn {
    font-size: 13px;
    padding: 6px 10px;
  }
  
  /* 调整小屏幕上的登录按钮样式 */
  .mobile-login-btn {
    top: 8px;
    right: 10px;
    padding: 6px 12px;
    font-size: 13px;
  }
  
  .mobile-login-placeholder {
    width: 70px; /* 小屏幕上稍窄 */
  }
  
  .main-nav {
    padding: 10px 0;
  }
  
  .nav-link {
    padding: 8px 12px;
    font-size: 14px;
  }
}
/* 文章列表样式 */
.article-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.article-item {
  display: flex;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.article-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
}

.article-cover {
  width: 300px;
  height: 100%;
  flex-shrink: 0;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.article-content {
  flex: 1;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.article-title {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #333;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  color: #666;
  font-size: 14px;
}

.article-excerpt {
  color: #666;
  line-height: 1.6;
  margin-bottom: 15px;
  flex: 1;
}

.article-actions {
  margin-top: auto;
}

.read-more-btn {
  background: #1a237e;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.read-more-btn:hover {
  background: #283593;
}

/* 加载、错误和空状态样式 */
.loading-container,
.error-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  text-align: center;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #1a237e;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 15px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-message {
  color: #d32f2f;
  margin-bottom: 10px;
}

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  margin-top: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.page-btn {
  background: #1a237e;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: #283593;
}

.page-btn:disabled {
  background: #cccccc;
  cursor: not-allowed;
}

.page-info {
  color: #666;
  font-size: 14px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-item {
    flex-direction: column;
  }
  
  .article-cover {
    width: 100%;
    height: 200px;
  }
  
  .article-content {
    padding: 15px;
  }
  
  .article-title {
    font-size: 18px;
  }
  
  .article-meta {
    flex-wrap: wrap;
    gap: 10px;
  }
}

</style>
<style>
/* 侧边悬浮关注功能 - 全局样式 */
.floating-sidebar {
  position: fixed;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
  z-index: 999;
  transition: all 0.3s ease;
}

.floating-trigger {
  background-color: #ff6600;
  color: white;
  width: 40px;
  height: 120px;
  border-radius: 6px 0 0 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  box-shadow: -2px 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.floating-trigger:hover {
  background-color: #ff5500;
}

.floating-trigger span {
  writing-mode: vertical-rl;
  text-orientation: mixed;
  font-size: 14px;
  font-weight: bold;
  letter-spacing: 2px;
}

.floating-content {
  position: absolute;
  right: 100%;
  top: 50%;
  transform: translateY(-50%);
  background: white;
  border-radius: 6px;
  box-shadow: -2px 2px 12px rgba(0, 0, 0, 0.15);
  padding: 20px;
  min-width: 200px;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
  margin-right: -10px;
}

.floating-sidebar:hover .floating-content {
  opacity: 1;
  visibility: visible;
  margin-right: 0;
}

.qrcode-wrapper {
  text-align: center;
}

.qrcode-image {
  width: 160px;
  height: 160px;
  border: 1px solid #eaeaea;
  padding: 10px;
  border-radius: 4px;
  background: white;
  margin-bottom: 10px;
}

.qrcode-text {
  font-size: 14px;
  color: #666;
  margin-bottom: 15px;
}

.follow-button {
  width: 100%;
  padding: 8px 0;
  background-color: #ff6600;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-bottom: 15px;
}

.follow-button:hover {
  background-color: #ff5500;
}

.social-share {
  display: flex;
  justify-content: center;
  gap: 15px;
}

.social-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
  font-size: 16px;
}

.social-icon:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.social-icon.wechat {
  background-color: #07C160;
}

.social-icon.weibo {
  background-color: #E6162D;
}

.social-icon.github {
  background-color: #333;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .floating-sidebar {
    right: 10px;
  }
  
  .floating-trigger {
    width: 36px;
    height: 100px;
  }
  
  .floating-content {
    min-width: 180px;
    padding: 15px;
  }
  
  .qrcode-image {
    width: 140px;
    height: 140px;
  }
}
</style>