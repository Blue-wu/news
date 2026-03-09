<template>
  <div class="article-management">
    <div class="card">
      <div class="card-header">
        <h2>文章管理</h2>
        <router-link to="/admin/articles/edit" class="create-btn">创建新文章</router-link>
      </div>
      
      <div class="card-body">
        <!-- 搜索和筛选 -->
        <div class="filters">
          <div class="search-box">
            <input 
              type="text" 
              v-model="searchQuery" 
              placeholder="搜索文章标题..." 
              class="search-input"
              @keyup.enter="handleSearch"
            />
            <button class="search-btn" @click="handleSearch">搜索</button>
          </div>
          <select v-model="statusFilter" class="status-filter" @change="handleSearch">
            <option value="">全部状态</option>
            <option value="1">已发布</option>
            <option value="0">草稿</option>
          </select>
        </div>
        
        <!-- 文章列表 -->
        <div v-if="loading" class="loading">加载中...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="filteredArticles.length === 0" class="error">暂无文章数据</div>
        <div v-else class="table-container">
          <table class="article-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>标题</th>
                <th>作者</th>
                <th>发布日期</th>
                <th>浏览量</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(article, index) in filteredArticles" :key="index">
                <td>{{ article.id || 'N/A' }}</td>
                <td class="title-cell">{{ article.title || '无标题' }}</td>
                <td>{{ article.user?.nickname || article.user?.username || article.user_name || article.author || '未知' }}</td>
                <td>{{ formatDate(article.created_at) }}</td>
                <td>{{ article.view_count || article.viewCount || 0 }}</td>
                <td>
                  <span class="status-badge" :class="article.status === 1 ? 'published' : 'draft'">
                    {{ article.status === 1 ? '已发布' : '草稿' }}
                  </span>
                </td>
                <td class="actions" v-if="article.id">
                  <!-- 修改查看按钮，使用方法调用而不是router-link -->
                  <button @click="viewArticle(article.id)" class="view-btn">
                    👁️ 查看
                  </button>
                  <router-link :to="`/admin/articles/edit/${article.id}`" class="edit-btn">
                    ✏️ 编辑
                  </router-link>
                  <button @click="deleteArticle(article.id)" class="delete-btn">
                    🗑️ 删除
                  </button>
                </td>
                <td class="actions" v-else>
                  <span>无效文章</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <!-- 分页 -->
        <div v-if="!loading && filteredArticles.length > 0" class="pagination">
          <button @click="currentPage--" :disabled="currentPage === 1" class="page-btn">上一页</button>
          <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
          <button @click="currentPage++" :disabled="currentPage === totalPages" class="page-btn">下一页</button>
        </div>
      </div>
    </div>
    
    <!-- 文章详情模态框 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ selectedArticle.title }}</h3>
          <button class="close-btn" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div v-if="loadingArticle" class="loading">加载文章中...</div>
          <div v-else-if="articleError" class="error">{{ articleError }}</div>
          <div v-else-if="selectedArticle" class="article-detail">
            <div class="article-meta">
              <span>作者: {{ selectedArticle.user?.nickname || selectedArticle.user?.username || selectedArticle.user_name || selectedArticle.author || '未知' }}</span>
              <span>发布日期: {{ formatDate(selectedArticle.created_at) }}</span>
              <span>浏览量: {{ selectedArticle.view_count || 0 }}</span>
              <span :class="selectedArticle.status === 1 ? 'published' : 'draft'">
                {{ selectedArticle.status === 1 ? '已发布' : '草稿' }}
              </span>
            </div>
            <div class="article-content" v-html="selectedArticle.content"></div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="close-button">关闭</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { articleApi, adminApi } from '@/utils/api.js';

const articles = ref([]);
const loading = ref(false);
const error = ref('');
const searchQuery = ref('');
const statusFilter = ref('');
const currentPage = ref(1);
const pageSize = 10;

// 模态框相关状态
const showModal = ref(false);
const selectedArticle = ref(null);
const loadingArticle = ref(false);
const articleError = ref('');

onMounted(() => {
  fetchArticles();
});

const fetchArticles = async () => {
  loading.value = true;
  error.value = '';
  try {
    // 传递分页参数给API
    const response = await articleApi.getArticles(currentPage.value, pageSize);
    console.log('API响应原始数据:', response);
    
    // 正确处理API响应格式，特别关注total字段
    if (response) {
      // 检查是否有total字段（从截图中看到API返回了total: 37）
      if (response.total !== undefined) {
        totalArticles.value = response.total;
      }
      
      if (response.data) {
        // 确保data是数组格式
        if (Array.isArray(response.data)) {
          articles.value = response.data;
        } else {
          // 如果data不是数组，检查是否包含items或list等数组字段
          articles.value = response.data.items || response.data.list || [response.data];
        }
      } else {
        // 如果没有data字段，尝试直接使用response作为数据
        articles.value = Array.isArray(response) ? response : [response];
      }
    }
    
    console.log('处理后的文章数据:', articles.value);
    console.log('总文章数:', totalArticles.value);
  } catch (err) {
    console.error('获取文章失败:', err);
    error.value = err.response?.data?.error || '获取文章列表失败';
  } finally {
    loading.value = false;
  }
};

// 监听currentPage变化，当页码改变时重新获取数据
watch(() => currentPage.value, () => {
  fetchArticles();
});

// 查看文章详情
const viewArticle = async (id) => {
  loadingArticle.value = true;
  articleError.value = '';
  
  try {
    // 使用adminApi获取文章详情，并正确处理响应格式
    const response = await adminApi.getArticleDetail(id);
    selectedArticle.value = response.data || response;
    showModal.value = true;
  } catch (err) {
    articleError.value = err.response?.data?.error || '获取文章详情失败';
  } finally {
    loadingArticle.value = false;
  }
};

// 关闭模态框
const closeModal = () => {
  showModal.value = false;
  selectedArticle.value = null;
  articleError.value = '';
};

// 删除文章
const deleteArticle = async (id) => {
  if (confirm('确定要删除这篇文章吗？')) {
    try {
      await articleApi.deleteArticle(id);
      // 重新获取文章列表
      fetchArticles();
    } catch (err) {
      alert('删除失败: ' + (err.response?.data?.error || '未知错误'));
    }
  }
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

// 过滤文章 - 在后端分页模式下，过滤主要在后端完成
// 这里只做简单的数据展示，因为我们已经通过API传递了分页参数
const filteredArticles = computed(() => {
  let result = articles.value;
  
  // 注意：在后端分页模式下，这里的过滤只是对当前页数据的补充过滤
  // 如果需要全量过滤，应该在API层面实现（例如通过搜索接口）
  
  // 按标题搜索 - 只过滤当前页数据
  if (searchQuery.value) {
    result = result.filter(article => 
      article.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }
  
  // 按状态过滤 - 只过滤当前页数据
  if (statusFilter.value !== '') {
    result = result.filter(article => article.status === parseInt(statusFilter.value));
  }
  
  // 注意：由于我们使用后端分页，这里不再需要前端slice操作
  // 每页的数据已经由后端返回
  return result;
});

// 搜索按钮点击事件
const handleSearch = () => {
  // 重置到第一页
  currentPage.value = 1;
  // 重新获取数据，让后端处理搜索
  fetchArticles();
};

// 监听搜索和筛选条件变化
watch([searchQuery, statusFilter], () => {
  // 当搜索条件或筛选条件变化时，重置到第一页
  currentPage.value = 1;
  // 重新获取数据
  fetchArticles();
});

// 总文章数（从API获取）
const totalArticles = ref(0);

// 总页数 - 基于API返回的total值
const totalPages = computed(() => {
  // 如果有API返回的total值，优先使用它
  if (totalArticles.value > 0) {
    return Math.ceil(totalArticles.value / pageSize);
  }
  
  // 否则回退到基于当前数据的计算
  let filteredTotal = articles.value;
  
  // 应用相同的过滤条件计算总数
  if (searchQuery.value) {
    filteredTotal = filteredTotal.filter(article => 
      article.title.toLowerCase().includes(searchQuery.value.toLowerCase())
    );
  }
  
  if (statusFilter.value !== '') {
    filteredTotal = filteredTotal.filter(article => 
      article.status === parseInt(statusFilter.value)
    );
  }
  
  return Math.ceil(filteredTotal.length / pageSize);
});
</script>

<style scoped>
/* 现有样式保持不变 */

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  padding: 16px 24px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  color: #333;
}

.close-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.close-btn:hover {
  color: #333;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.article-detail .article-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  color: #666;
  font-size: 14px;
  flex-wrap: wrap;
}

.article-detail .article-meta span {
  padding: 4px 8px;
  background: #f5f5f5;
  border-radius: 4px;
}

.article-detail .article-meta span.published {
  background: #e6f7ff;
  color: #1890ff;
}

.article-detail .article-meta span.draft {
  background: #fff7e6;
  color: #fa8c16;
}

.article-content {
  line-height: 1.8;
  color: #333;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #eee;
  text-align: right;
}

.close-button {
  padding: 8px 16px;
  background: #f5f5f5;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: #333;
}

.close-button:hover {
  background: #e8e8e8;
}
</style>

<style scoped>
.article-management {
  padding: 20px;
  /* 移动端适配 */
  @media (max-width: 768px) {
    padding: 15px;
  }
  @media (max-width: 480px) {
    padding: 10px;
  }
}

/* 现代卡片设计 */
.card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.card:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.card-header {
  padding: 24px 24px 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 15px;
}

.card-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

/* 现代化按钮设计 */
.create-btn {
  background: white;
  color: #667eea;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  text-decoration: none;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.create-btn:hover {
  background: #f8f9ff;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.card-body {
  padding: 24px;
  /* 移动端适配 */
  @media (max-width: 768px) {
    padding: 16px;
  }
}

/* 过滤器区域 */
.filters {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  align-items: stretch;
}

.search-box {
  flex: 1;
  min-width: 200px;
  display: flex;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border-radius: 8px;
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.search-box:focus-within {
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.search-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid #f0f2f5;
  border-right: none;
  border-radius: 8px 0 0 8px;
  outline: none;
  font-size: 14px;
  transition: border-color 0.3s ease;
}

.search-input:focus {
  border-color: #667eea;
}

.search-btn {
  background: #667eea;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 0 8px 8px 0;
  cursor: pointer;
  font-weight: 500;
  transition: background-color 0.3s ease;
  min-width: 80px;
}

.search-btn:hover {
  background: #5a67d8;
}

.status-filter {
  padding: 12px 16px;
  border: 2px solid #f0f2f5;
  border-radius: 8px;
  outline: none;
  background: white;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.status-filter:focus {
  border-color: #667eea;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

/* 表格样式 - 现代化设计 */
.table-container {
  overflow-x: auto;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.article-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
}

.article-table th,
.article-table td {
  padding: 16px;
  text-align: left;
  border-bottom: 1px solid #f0f2f5;
}

.article-table th {
  background: #f8fafc;
  font-weight: 600;
  color: #2d3748;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.article-table tbody tr {
  transition: background-color 0.2s ease;
}

.article-table tbody tr:hover {
  background-color: #f8fafc;
}

.title-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  color: #2d3748;
}

.title-cell:hover {
  color: #667eea;
  text-decoration: underline;
}

/* 状态徽章 */
.status-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.status-badge.published {
  background: #ebf8ff;
  color: #3182ce;
}

.status-badge.draft {
  background: #fffaf0;
  color: #ed8936;
}

/* 操作按钮组 */
.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.view-btn,
.edit-btn,
.delete-btn {
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  min-width: 50px;
  justify-content: center;
}

.view-btn {
  background: #e6fffa;
  color: #319795;
}

.view-btn:hover {
  background: #b2f5ea;
  transform: translateY(-1px);
}

.edit-btn {
  background: #ebf8ff;
  color: #3182ce;
}

.edit-btn:hover {
  background: #bee3f8;
  transform: translateY(-1px);
}

.delete-btn {
  background: #fed7d7;
  color: #e53e3e;
}

.delete-btn:hover {
  background: #feb2b2;
  transform: translateY(-1px);
}

/* 分页组件 */
.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.page-btn {
  padding: 10px 16px;
  border: 2px solid #f0f2f5;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s ease;
}

.page-btn:hover:not(:disabled) {
  border-color: #667eea;
  color: #667eea;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.page-btn:disabled {
  background: #f8f9fa;
  cursor: not-allowed;
  color: #a0aec0;
  border-color: #e2e8f0;
}

.page-info {
  font-size: 14px;
  color: #718096;
  font-weight: 500;
}

/* 加载和错误状态 */
.loading,
.error {
  padding: 60px;
  text-align: center;
  font-size: 16px;
}

.loading {
  color: #667eea;
  font-weight: 500;
}

.error {
  color: #e53e3e;
  background: #fed7d7;
  border-radius: 8px;
  margin: 20px 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .card-header {
    padding: 20px 16px 16px;
  }
  
  .card-header h2 {
    font-size: 20px;
  }
  
  .create-btn {
    padding: 8px 16px;
    font-size: 13px;
  }
  
  .filters {
    flex-direction: column;
    gap: 12px;
  }
  
  .search-box,
  .status-filter {
    width: 100%;
  }
  
  .article-table th,
  .article-table td {
    padding: 12px;
    font-size: 14px;
  }
  
  .title-cell {
    max-width: 150px;
  }
  
  .actions {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .table-container {
    font-size: 12px;
  }
  
  .article-table th,
  .article-table td {
    padding: 8px;
  }
  
  .status-badge {
    padding: 4px 8px;
    font-size: 11px;
  }
  
  .view-btn,
  .edit-btn,
  .delete-btn {
    padding: 6px 10px;
    font-size: 11px;
  }
}
</style>