<script setup>
import { ref, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { articleApi } from '@/utils/api.js';

const route = useRoute();
const router = useRouter();
const articles = ref([]);
const loading = ref(false);
const error = ref('');
const keyword = ref('');
const currentPage = ref(1);
const pageSize = ref(10);
const totalCount = ref(0);
const totalPages = ref(0);

// 获取搜索结果 - 移到watch之前定义
// 修改fetchSearchResults函数
const fetchSearchResults = async () => {
  if (!keyword.value) return;
  
  loading.value = true;
  error.value = '';
  try {
    // 确保调用正确的API方法和参数
    const response = await articleApi.searchArticles(keyword.value, currentPage.value, pageSize.value);
    
    // 处理响应数据
    articles.value = response.data || [];
    totalCount.value = response.total || 0;
    totalPages.value = Math.ceil(totalCount.value / pageSize.value);
  } catch (err) {
    error.value = '搜索失败，请稍后重试';
    console.error('Search error:', err);
  } finally {
    loading.value = false;
  }
};

// 监听路由参数变化
watch(() => route.query.q, (newKeyword) => {
  if (newKeyword) {
    keyword.value = newKeyword;
    currentPage.value = 1;
    fetchSearchResults();
  }
}, { immediate: true });

// 监听页码变化
watch(currentPage, () => {
  fetchSearchResults();
});

// 跳转到文章详情
const goToArticleDetail = (articleId) => {
  router.push(`/articles/${articleId}`);
};

// 分页方法 - 添加const声明
const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page;
  }
};

const goToPrevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
};

const goToNextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

const goToFirstPage = () => {
  currentPage.value = 1;
};

const goToLastPage = () => {
  currentPage.value = totalPages.value;
};
</script>


<template>
	<div class="search-results">
		<div class="container">
			<div class="search-header">
				<h1>搜索结果："{{ keyword }}"</h1>
				<p>找到 {{ totalCount }} 条相关结果</p>
			</div>
			
			<div v-if="loading" class="loading">加载中...</div>
			<div v-else-if="error" class="error">{{ error }}</div>
			<div v-else-if="articles.length === 0" class="empty">
				<p>没有找到相关内容，请尝试其他关键词</p>
			</div>
			<div v-else class="article-list">
				<div v-for="article in articles" :key="article.id" class="article-item">
					<div class="article-header">
						<h3 class="article-title">
							<router-link :to="'/articles/' + article.id">{{ article.title }}</router-link>
						</h3>
					</div>
					<p class="article-summary">{{ article.summary }}</p>
					<div class="article-meta">
						<span>{{ article.user?.username || '未知用户' }}</span>
						<span class="dot">·</span>
						<span>{{ new Date(article.created_at).toLocaleDateString() }}</span>
						<span class="dot">·</span>
						<span>{{ article.view_count }} 阅读</span>
					</div>
				</div>
			</div>
			
			<!-- 分页组件 -->
			<div v-if="totalPages > 1" class="pagination">
				<button @click="goToFirstPage" :disabled="currentPage === 1" class="pagination-btn">首页</button>
				<button @click="goToPrevPage" :disabled="currentPage === 1" class="pagination-btn">上一页</button>
				<span class="pagination-info">第 {{ currentPage }} / {{ totalPages }} 页</span>
				<button @click="goToNextPage" :disabled="currentPage === totalPages" class="pagination-btn">下一页</button>
				<button @click="goToLastPage" :disabled="currentPage === totalPages" class="pagination-btn">末页</button>
			</div>
		</div>
	</div>
</template>

<style scoped>
.search-results {
	padding: 20px 0;
}

.container {
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 20px;
}

.search-header {
	margin-bottom: 30px;
	padding-bottom: 20px;
	border-bottom: 1px solid #eee;
}

.search-header h1 {
	font-size: 24px;
	color: #333;
	margin-bottom: 10px;
}

.search-header p {
	color: #666;
	font-size: 14px;
}

.article-list {
	margin-bottom: 30px;
}

.article-item {
	background: #fff;
	padding: 20px;
	margin-bottom: 15px;
	border-radius: 8px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
	transition: all 0.3s ease;
}

.article-item:hover {
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.article-title {
	font-size: 18px;
	margin-bottom: 10px;
	line-height: 1.4;
}

.article-title a {
	color: #333;
	text-decoration: none;
	transition: color 0.3s ease;
}

.article-title a:hover {
	color: #ff6600;
}

.article-summary {
	color: #666;
	font-size: 14px;
	line-height: 1.6;
	margin-bottom: 10px;
	display: -webkit-box;
	-webkit-line-clamp: 2;
	-webkit-box-orient: vertical;
	overflow: hidden;
}

.article-meta {
	font-size: 13px;
	color: #999;
}

.dot {
	margin: 0 8px;
}

.loading,
.error,
.empty {
	text-align: center;
	padding: 40px;
	color: #666;
}

.pagination {
	display: flex;
	justify-content: center;
	align-items: center;
	gap: 10px;
	margin-top: 30px;
}

.pagination-btn {
	padding: 6px 12px;
	border: 1px solid #ddd;
	background: #fff;
	color: #333;
	border-radius: 4px;
	cursor: pointer;
	transition: all 0.3s ease;
}

.pagination-btn:hover:not(:disabled) {
	background: #f5f5f5;
	border-color: #ccc;
}

.pagination-btn:disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

.pagination-info {
	padding: 0 10px;
	color: #666;
}
</style>