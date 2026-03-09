<script setup>
import { ref, onMounted, watch } from 'vue'; // 添加watch导入
import { useRoute, useRouter  } from 'vue-router';
import axios from 'axios';
// 导入我们的api实例
import { articleApi } from '@/utils/api.js';

const route = useRoute();
const router = useRouter(); // 添加router变量
const articles = ref([]);
const loading = ref(false);
const error = ref('');
// 添加分页相关变量
const currentPage = ref(1);
const pageSize = ref(10);
const totalCount = ref(0);
const totalPages = ref(0);
// 添加跳转到文章详情页的函数
const goToArticleDetail = (articleId) => {
  router.push(`/articles/${articleId}`);
};

// 获取文章列表
// 获取文章列表 - 添加分页参数
const fetchArticles = async () => {
  loading.value = true;
  error.value = '';
  try {
    let response;
    if (route.params.slug) {
      // 如果有分类参数，获取该分类下的文章
      response = await articleApi.getArticlesByCategory(route.params.slug, currentPage.value, pageSize.value);
    } else {
      // 获取所有文章
      response = await articleApi.getArticles(currentPage.value, pageSize.value);
    }
    
    // 假设后端返回的数据格式为 { data: [...], total: 1693 }
    articles.value = response.data || response; // 兼容不同的响应格式
    
    // 设置总页数信息
    totalCount.value = response.total || 1693; // 从响应中获取总数，或者使用默认值
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } catch (err) {
    error.value = '获取文章列表失败：' + (err.response?.data?.message || err.message);
    // 如果API调用失败，使用模拟数据
    articles.value = [
      // ... 模拟数据保持不变 ...
    ];
    // 模拟分页数据
    totalCount.value = 1693;
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } finally {
    loading.value = false;
  }
};

// 分页处理函数
const handlePageChange = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
    fetchArticles();
  }
};

const goToNextPage = () => {
  if (currentPage.value < totalPages.value) {
    handlePageChange(currentPage.value + 1);
  }
};

const goToPrevPage = () => {
  if (currentPage.value > 1) {
    handlePageChange(currentPage.value - 1);
  }
};

const goToFirstPage = () => {
  handlePageChange(1);
};

const goToLastPage = () => {
  handlePageChange(totalPages.value);
};

// 生成页码数组
const getPageNumbers = () => {
  const pages = [];
  const maxPagesToShow = 10;
  let startPage = Math.max(1, currentPage.value - Math.floor(maxPagesToShow / 2));
  let endPage = Math.min(totalPages.value, startPage + maxPagesToShow - 1);
  
  // 调整起始页码，确保显示的页码数量一致
  if (endPage - startPage + 1 < maxPagesToShow) {
    startPage = Math.max(1, endPage - maxPagesToShow + 1);
  }
  
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }
  return pages;
};

// 删除文章
const deleteArticle = async (id) => {
  if (!confirm('确定要删除这篇文章吗？')) return;
  
  try {
    await axios.delete(`http://localhost:8000/api/articles/${id}`);
    // 重新获取文章列表
    fetchArticles();
  } catch (err) {
    error.value = '删除文章失败：' + (err.response?.data?.message || err.message);
  }
};

// 监听路由参数变化
route.params.slug && fetchArticles();
// 监听路由参数变化
watch(() => route.params.slug, (newSlug) => {
     currentPage.value = 1; // 切换分类时重置到第一页
  fetchArticles(); // 当路由参数变化时重新获取数据
});
onMounted(() => {
  fetchArticles();
});
</script>

<template>
  <div class="article-list-page">
    <h1 class="section-title">{{ route.params.slug ? route.params.slug + '分类' : '最新文章' }}</h1>
    
    <!-- 错误提示 -->
    <div v-if="error" class="error">{{ error }}</div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading">加载中...</div>
    
    <!-- 文章列表 -->
    <div v-else-if="articles.length > 0" class="articles-grid">
<div v-for="article in articles" :key="article.id" class="article-item">
  <div class="article-image-container">
    <router-link :to="'/articles/' + article.id">
      <img :src="article.cover_image || 'https://picsum.photos/800/600?random=' + article.id" alt="文章封面" class="article-image">
    </router-link>
  </div>
  <div class="article-content">
    <div class="article-header">
      <span v-if="article.categories && article.categories.length > 0" class="article-tag">{{ article.categories[0].name }}</span>
      <h3 class="article-title">
        <router-link :to="'/articles/' + article.id">{{ article.title }}</router-link>
      </h3>
    </div>
      <p class="article-summary" @click="goToArticleDetail(article.id)">{{ article.summary }}</p>
    <div class="article-meta">
      <span>{{ article.user?.nickname || article.user?.username || article.user_name || '未知用户' }}</span>
      <span class="dot">·</span>
      <span>{{ new Date(article.created_at).toLocaleDateString() }}</span>
      <span class="dot">·</span>
      <span>{{ article.view_count }} 阅读</span>
      <span class="dot">·</span>
      <span>{{ article.comment_count }} 评论</span>
    </div>
  </div>
</div>
    </div>
    
    <!-- 空状态 -->
    <div v-else class="empty">暂无文章</div>
  <!-- 分页组件 -->
    <div class="pagination">
      <button 
        @click="goToFirstPage" 
        :disabled="currentPage === 1"
        class="pagination-btn"
      >
        首页
      </button>
      <button 
        @click="goToPrevPage" 
        :disabled="currentPage === 1"
        class="pagination-btn"
      >
        上一页
      </button>
      
      <button 
        v-for="page in getPageNumbers()" 
        :key="page"
        @click="handlePageChange(page)"
        :class="['pagination-btn', { active: currentPage === page }]"
      >
        {{ page }}
      </button>
      
      <button 
        @click="goToNextPage" 
        :disabled="currentPage === totalPages"
        class="pagination-btn"
      >
        下一页
      </button>
      <button 
        @click="goToLastPage" 
        :disabled="currentPage === totalPages"
        class="pagination-btn"
      >
        末页
      </button>
      
      <span class="page-info">共 {{ totalPages }} 页</span>
    </div>
  </div>
</template>

<style scoped>
/* 分页样式 */
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-top: 30px;
  flex-wrap: wrap;
}

.pagination-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  background-color: #fff;
  color: #333;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.pagination-btn:hover:not(:disabled) {
  background-color: #f5f5f5;
  border-color: #999;
}

.pagination-btn.active {
  background-color: #ff6600;
  color: #fff;
  border-color: #ff6600;
}

.pagination-btn:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.page-info {
  margin-left: 10px;
  color: #666;
}

@media (max-width: 768px) {
  .pagination {
    gap: 4px;
  }
  
  .pagination-btn {
    padding: 4px 8px;
    font-size: 14px;
  }
}

.article-list-page {
  padding: 20px;
}

.section-title {
  font-size: 20px;
  font-weight: bold;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 10px;
  border-bottom: 2px solid #f0f0f0;
}

.article-image-container a {
  display: block;
  width: 100%;
  height: 100%;
}

.article-image-container a:hover .article-image {
  transform: scale(1.05);
}

.article-summary {
  color: #666;
  line-height: 1.6;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  cursor: pointer;
  transition: color 0.3s ease;
}

.article-summary:hover {
  color: #ff6600;
}

.article-item:hover .article-title a {
  color: #ff6600;
}

.articles-grid {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.article-item {
  display: flex;
  gap: 20px;
  padding: 20px;
  background-color: #fff;
  border: 1px solid #eaeaea;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.article-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.article-image-container {
  flex-shrink: 0;
  width: 250px;
  height: 180px;
  overflow: hidden;
  border-radius: 6px;
}

.article-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.article-item:hover .article-image {
  transform: scale(1.05);
}

.article-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.article-header {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  margin-bottom: 12px;
}

.article-tag {
  background-color: #ff6600;
  color: white;
  padding: 2px 8px;
  border-radius: 3px;
  font-size: 12px;
  font-weight: bold;
  flex-shrink: 0;
}

.article-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin: 0;
  line-height: 1.4;
  flex: 1;
}

.article-title a {
  color: #333;
  text-decoration: none;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-title a:hover {
  color: #ff6600;
}

.article-summary {
  color: #666;
  line-height: 1.6;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.article-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #999;
}

.article-meta .dot {
  color: #ccc;
}

.error {
  color: #ff4757;
  padding: 10px;
  background-color: #ffebee;
  border-radius: 4px;
  margin-bottom: 20px;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .article-item {
    flex-direction: column;
  }
  
  .article-image-container {
    width: 100%;
    height: 200px;
  }
  
  .article-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .article-title {
    font-size: 16px;
  }
}
</style>