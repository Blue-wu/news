<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articleApi } from '@/utils/api.js'; // 使用封装的API

const route = useRoute()
const router = useRouter()
const article = ref(null)
const loading = ref(true)
const error = ref('')
const relatedArticles = ref([]) // 相关文章数据
const randomArticles = ref([]) // 随机文章数据
// 添加上下篇文章变量
const prevArticle = ref(null)
const nextArticle = ref(null)

// 初始加载
onMounted(async () => {
  await fetchData()
})

// 添加路由参数监听
watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      // 当ID变化时重新加载所有数据
      fetchData()
    }
  }
)

// 将所有数据获取方法封装到一个函数中
const fetchData = async () => {
  await Promise.all([
    fetchArticle(),
    fetchRelatedArticles(),
    fetchRandomArticles(),
    fetchPrevNextArticles()
  ])
}

const fetchArticle = async () => {
  const id = route.params.id
  try {
    loading.value = true
    // 使用封装的API方法，确保token正确传递
    article.value = await articleApi.getArticle(id)
     // 增加阅读计数，使用try-catch避免影响主要功能
    try {
      // 移除会话级别限制，每次加载都增加阅读计数
      await articleApi.incrementViewCount(id)
      // 更新本地的阅读计数显示
      article.value.view_count = (article.value.view_count || 0) + 1
    } catch (viewError) {
      console.warn('更新阅读计数失败:', viewError)
      // 即使更新失败也不影响页面显示
    }
  } catch (err) {
    error.value = '获取文章详情失败'
    console.error('Error fetching article:', err)
    // 使用模拟数据，防止页面空白
    article.value = {
      id: id,
      title: '文章标题',
      content: '<p>文章内容加载失败，这里是示例内容。</p>',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      view_count: 0,
      comment_count: 0,
      user: {
        username: '未知用户',
        verified: false
      },
      categories: [
        { name: '未分类', slug: 'uncategorized' }
      ]
    }
  } finally {
    loading.value = false
  }
}

// 获取相关文章
const fetchRelatedArticles = async () => {
  try {
    // 这里简化处理，实际应该根据当前文章的分类或标签获取相关文章
    const allArticles = await articleApi.getArticles()
    relatedArticles.value = (Array.isArray(allArticles) ? allArticles : [])
      .slice(0, 3)
      .map(article => ({
        id: article.id,
        title: article.title,
        view_count: article.view_count || 0
      }))
  } catch (err) {
    console.error('Error fetching related articles:', err)
    // 提供默认数据
    relatedArticles.value = [
      { id: 1, title: '热门文章1', view_count: 1234 },
      { id: 2, title: '热门文章2', view_count: 987 },
      { id: 3, title: '热门文章3', view_count: 567 }
    ]
  }
}
// 添加获取上下篇文章的函数
const fetchPrevNextArticles = async () => {
  const currentId = parseInt(route.params.id) // 确保ID是数字类型
  try {
    // 获取所有文章来确定上下篇关系
    const allArticles = await articleApi.getArticles(1, 100)
    
    // 正确处理API响应格式
    const articles = Array.isArray(allArticles) ? allArticles : (allArticles.data || [])
    
    // 确保只处理有ID的文章并按发布时间排序
    const validArticles = articles.filter(article => article.id && article.created_at)
    const sortedArticles = [...validArticles].sort((a, b) => 
      new Date(b.created_at) - new Date(a.created_at)
    )
    
    // 找到当前文章的索引，使用严格相等
    const currentIndex = sortedArticles.findIndex(article => article.id === currentId)
    
    if (currentIndex !== -1) {
      // 设置上一篇文章（索引+1的文章，因为是降序排列）
      if (currentIndex < sortedArticles.length - 1) {
        prevArticle.value = sortedArticles[currentIndex + 1]
      }
      
      // 设置下一篇文章（索引-1的文章，因为是降序排列）
      if (currentIndex > 0) {
        nextArticle.value = sortedArticles[currentIndex - 1]
      }
    }
  } catch (err) {
    console.error('获取上下篇文章失败:', err)
    // 提供更安全的模拟数据
    if (!prevArticle.value) {
      prevArticle.value = {
        id: currentId - 1,
        title: '上一篇文章标题',
        created_at: new Date().toISOString()
      }
    }
    if (!nextArticle.value) {
      nextArticle.value = {
        id: currentId + 1,
        title: '下一篇文章标题',
        created_at: new Date().toISOString()
      }
    }
  }
}
// 获取随机文章
const fetchRandomArticles = async () => {
  try {
    // 这里简化处理，实际应该获取随机推荐的文章
    const allArticles = await articleApi.getArticles()
    const shuffled = (Array.isArray(allArticles) ? allArticles : [])
      .sort(() => 0.5 - Math.random())
      .slice(0, 4)
    randomArticles.value = shuffled.map(article => ({
      id: article.id,
      title: article.title,
      cover_image: article.cover_image || `https://picsum.photos/300/200?random=${article.id}`
    }))
  } catch (err) {
    console.error('Error fetching random articles:', err)
    // 提供默认数据
    randomArticles.value = [
      { id: 1, title: '文章标题1', cover_image: 'https://picsum.photos/300/200?random=1' },
      { id: 2, title: '文章标题2', cover_image: 'https://picsum.photos/300/200?random=2' },
      { id: 3, title: '文章标题3', cover_image: 'https://picsum.photos/300/200?random=3' },
      { id: 4, title: '文章标题4', cover_image: 'https://picsum.photos/300/200?random=4' }
    ]
  }
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 刷新随机文章
const refreshRandomArticles = async () => {
  await fetchRandomArticles()
}
</script>

<template>
  <div class="article-detail-container">
    <!-- 错误提示 -->
    <div v-if="error" class="error">{{ error }}</div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading">加载中...</div>
    
    <div v-else-if="article" class="main-content">
      
      <!-- 文章详情主内容 -->
      <div class="content-wrapper">
        <article class="article-content">
          <!-- 文章标题 -->
          <h1 class="article-title">{{ article.title }}</h1>
          
          <!-- 作者信息和元数据 -->
          <div class="article-meta">
            <div class="author-info">
              <img 
                src="https://picsum.photos/32/32?random={{ article.user?.id || 1 }}" 
                alt="作者头像" 
                class="author-avatar"
              >
              <div>
                <span class="author-name">{{ article.user?.nickname || article.user?.username || article.user_name || '未知用户' }}</span>
                <span v-if="article.user?.verified" class="verified-badge">V</span>
              </div>
            </div>
            
            <div class="meta-stats">
              <span class="publish-time">{{ formatDate(article.created_at) }}</span>
              <span class="separator">/</span>
              <span class="view-count">{{ article.view_count || 0 }}阅读</span>
              <span class="separator">/</span>
              <span class="comment-count">{{ article.comment_count || 0 }}评论</span>
            </div>
          </div>
          
          <!-- 文章正文 -->
          <div class="article-body" v-html="article.content"></div>
          
          <!-- 文章标签 -->
          <div v-if="article.tags && article.tags.length > 0" class="article-tags">
            <span v-for="tag in article.tags" :key="tag.id" class="tag">
              {{ tag.name }}
            </span>
          </div>
          <!-- 本站声明 -->
          <div class="site-disclaimer">
            <blockquote>
              <p>本站声明：以上都本文内容由互联网用户自发贡献，该文观点仅代表作者本人。本站仅提供信息存储空间服务，不拥有所有权，不承担相关法律责任。如发现本站有涉嫌抄袭侵权/违法违规的内容，一经查实，本站将立刻删除。</p>
            </blockquote>
          </div>
           <!-- 添加上下篇导航 -->
          <div class="prev-next-nav">
  <!-- 上一篇 -->
<router-link 
  :to="`/articles/${prevArticle.id}`" 
  class="prev-article"
  v-if="prevArticle && prevArticle.id"
>
  <span class="nav-label">上一篇</span>
  <span class="nav-title">{{ prevArticle.title }}</span>
  <span class="nav-date">{{ prevArticle.created_at ? formatDate(prevArticle.created_at) : '' }}</span>
</router-link>

<!-- 下一篇 -->
<router-link 
  :to="`/articles/${nextArticle.id}`" 
  class="next-article"
  v-if="nextArticle && nextArticle.id"
>
  <span class="nav-label">下一篇</span>
  <span class="nav-title">{{ nextArticle.title }}</span>
  <span class="nav-date">{{ nextArticle.created_at ? formatDate(nextArticle.created_at) : '' }}</span>
</router-link>
</div>
          <!-- 返回按钮 -->
          <router-link to="/articles" class="back-btn">返回文章列表</router-link>
        </article>
        

      </div>
    </div>
    
    <div v-else class="not-found">文章不存在</div>
  </div>
</template>

<style scoped>
.article-detail-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
/* 上下篇导航样式 */
.prev-next-nav {
  display: flex;
  justify-content: space-between;
  margin: 30px 0;
  gap: 20px;
}
.article-tags {
  margin-bottom: 20px;
}

/* 本站声明样式 */
.site-disclaimer {
  margin: 30px 0;
  padding: 0;
}

.site-disclaimer blockquote {
  margin: 0;
  padding: 15px 20px;
  background-color: #f8f9fa;
  border-left: 4px solid #6c757d;
  border-radius: 0 4px 4px 0;
  font-size: 14px;
  line-height: 1.6;
  color: #666;
}

.site-disclaimer p {
  margin: 0;
  text-align: justify;
}

.prev-article,
.next-article {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
  text-decoration: none;
  color: #333;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.prev-article:hover,
.next-article:hover {
  background-color: #e9ecef;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.prev-article::before,
.next-article::after {
  content: '';
  position: absolute;
  width: 100%;
  height: 4px;
  bottom: 0;
  left: 0;
  background: linear-gradient(90deg, #ff6600, #ff9933);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.prev-article::before {
  transform-origin: left;
}

.next-article::after {
  transform-origin: right;
}

.prev-article:hover::before,
.next-article:hover::after {
  transform: scaleX(1);
}

.nav-label {
  font-size: 14px;
  color: #ff6600;
  margin-bottom: 8px;
  font-weight: 500;
}

.nav-title {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav-date {
  font-size: 12px;
  color: #999;
  margin-top: auto;
}

.next-article {
  text-align: right;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .prev-next-nav {
    flex-direction: column;
    gap: 15px;
  }
  
  .prev-article,
  .next-article {
    padding: 15px;
  }
  
  .next-article {
    text-align: left;
  }
}
.error {
  color: #ff4757;
  background-color: #ffebee;
  padding: 10px;
  border-radius: 4px;
  margin-bottom: 20px;
  text-align: center;
}

.loading {
  text-align: center;
  padding: 40px;
  color: #666;
}

.breadcrumb {
  font-size: 14px;
  color: #666;
  margin-bottom: 20px;
  padding: 10px 0;
}

.breadcrumb a {
  color: #666;
  text-decoration: none;
}

.breadcrumb a:hover {
  color: #ff6600;
}

.breadcrumb span {
  margin: 0 8px;
}

.content-wrapper {
  display: flex;
  gap: 30px;
}

.article-content {
  flex: 1;
  background-color: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.article-title {
  font-size: 28px;
  font-weight: bold;
  color: #333;
  margin-bottom: 20px;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #eee;
  margin-bottom: 30px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.author-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.author-name {
  font-weight: 500;
  color: #333;
}

.verified-badge {
  display: inline-block;
  background-color: #ff6600;
  color: white;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  font-size: 12px;
  text-align: center;
  line-height: 16px;
  margin-left: 4px;
}

.meta-stats {
  font-size: 14px;
  color: #999;
}

.separator {
  margin: 0 8px;
  color: #ccc;
}

.article-body {
  line-height: 1.8;
  color: #333;
  margin-bottom: 30px;
}

.article-body h1,
.article-body h2,
.article-body h3 {
  margin: 20px 0 15px 0;
  color: #333;
}

.article-body p {
  margin-bottom: 15px;
}

.article-body img {
  max-width: 100%;
  height: auto;
  margin: 15px 0;
  border-radius: 4px;
}

.article-tags {
  margin-bottom: 20px;
}

.tag {
  display: inline-block;
  background-color: #f5f5f5;
  color: #666;
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 14px;
  margin-right: 10px;
  margin-bottom: 10px;
}

.back-btn {
  display: inline-block;
  padding: 10px 20px;
  background-color: #f5f5f5;
  color: #333;
  text-decoration: none;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background-color: #e6e6e6;
}

/* 侧边栏样式 */
.sidebar {
  width: 300px;
  flex-shrink: 0;
}

.sidebar-section {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.sidebar-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
  position: relative;
  padding-bottom: 10px;
}

.sidebar-title::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 40px;
  height: 3px;
  background-color: #ff6600;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.refresh-btn {
  background: none;
  border: none;
  color: #999;
  cursor: pointer;
  font-size: 14px;
  padding: 5px;
}

.refresh-btn:hover {
  color: #ff6600;
}

/* 热门文章样式 */
.hot-articles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hot-article-item {
  display: flex;
  align-items: flex-start;
  cursor: pointer;
  transition: all 0.3s ease;
  padding: 5px;
  border-radius: 4px;
}

.hot-article-item:hover {
  background-color: #f9f9f9;
}

.article-rank {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background-color: #ff6600;
  color: white;
  font-size: 12px;
  font-weight: bold;
  border-radius: 50%;
  margin-right: 10px;
  flex-shrink: 0;
}

.hot-article-item:nth-child(1) .article-rank {
  background-color: #ff4444;
}

.hot-article-item:nth-child(2) .article-rank {
  background-color: #ff6600;
}

.hot-article-item:nth-child(3) .article-rank {
  background-color: #ffaa00;
}

.hot-article-item .article-title {
  font-size: 14px;
  color: #333;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

/* 随便看看样式 */
.random-articles {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.random-article-item {
  cursor: pointer;
  transition: transform 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.random-article-item:hover {
  transform: translateY(-3px);
}

.article-thumbnail {
  width: 100%;
  height: 80px;
  object-fit: cover;
  border-radius: 4px;
}

.random-article-item .article-title {
  font-size: 12px;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.3;
}

.not-found {
  text-align: center;
  padding: 60px 20px;
  color: #999;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

/* 响应式设计 */
@media (max-width: 992px) {
  .content-wrapper {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
  }
  
  .random-articles {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .article-detail-container {
    padding: 10px;
  }
  
  .article-content {
    padding: 20px;
  }
  
  .article-title {
    font-size: 24px;
  }
  
  .article-meta {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .random-articles {
    grid-template-columns: repeat(2, 1fr);
  }
}

.article-body {
  line-height: 1.8;
  color: #333;
  margin-bottom: 30px;
  font-size: 16px;
}

/* 增强富文本样式支持 */
.article-body h1,
.article-body h2,
.article-body h3,
.article-body h4,
.article-body h5,
.article-body h6 {
  margin: 20px 0 15px 0;
  color: #333;
  font-weight: 600;
}

.article-body h1 {
  font-size: 28px;
  padding-bottom: 10px;
  border-bottom: 1px solid #eee;
}

.article-body h2 {
  font-size: 24px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

.article-body h3 {
  font-size: 20px;
}

.article-body p {
  margin-bottom: 15px;
  text-align: justify;
}

.article-body img {
  max-width: 100%;
  height: auto;
  margin: 15px auto;
  border-radius: 4px;
  display: block;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 支持表格样式 */
.article-body table {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
}

.article-body th,
.article-body td {
  padding: 10px 12px;
  border: 1px solid #e0e0e0;
  text-align: left;
}

.article-body th {
  background-color: #f8f9fa;
  font-weight: 600;
}

.article-body tr:nth-child(even) {
  background-color: #f9f9f9;
}

/* 支持列表样式 */
.article-body ul,
.article-body ol {
  margin: 15px 0;
  padding-left: 30px;
}

.article-body li {
  margin-bottom: 8px;
}

/* 支持引用样式 */
.article-body blockquote {
  border-left: 4px solid #ff6600;
  padding-left: 15px;
  margin: 15px 0;
  color: #666;
  font-style: italic;
}

/* 支持代码块样式 */
.article-body pre {
  background-color: #f6f8fa;
  padding: 16px;
  border-radius: 4px;
  overflow-x: auto;
  margin: 15px 0;
}

.article-body code {
  font-family: 'Courier New', Courier, monospace;
  background-color: #f6f8fa;
  padding: 2px 4px;
  border-radius: 3px;
  font-size: 14px;
}

.article-body pre code {
  background-color: transparent;
  padding: 0;
}

/* 支持链接样式 */
.article-body a {
  color: #ff6600;
  text-decoration: none;
}

.article-body a:hover {
  text-decoration: underline;
}

/* 支持加粗、斜体等内联样式 */
.article-body strong {
  font-weight: 600;
  color: #333;
}

.article-body em {
  font-style: italic;
}

.article-body u {
  text-decoration: underline;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .article-body {
    font-size: 15px;
  }
  
  .article-body h1 {
    font-size: 24px;
  }
  
  .article-body h2 {
    font-size: 20px;
  }
  
  .article-body h3 {
    font-size: 18px;
  }
}
</style>