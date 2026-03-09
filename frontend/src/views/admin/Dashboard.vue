<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon articles-icon">📝</div>
        <div class="stat-content">
          <h3>文章总数</h3>
          <p class="stat-number">{{ articleCount }}</p>
          <p class="stat-change positive">+12% 较上月</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon users-icon">👥</div>
        <div class="stat-content">
          <h3>用户总数</h3>
          <p class="stat-number">{{ userCount }}</p>
          <p class="stat-change positive">+8% 较上月</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon views-icon">👁️</div>
        <div class="stat-content">
          <h3>总浏览量</h3>
          <p class="stat-number">{{ viewCount }}</p>
          <p class="stat-change positive">+23% 较上月</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon drafts-icon">📄</div>
        <div class="stat-content">
          <h3>草稿数</h3>
          <p class="stat-number">{{ draftCount }}</p>
          <p class="stat-change negative">-5% 较上月</p>
        </div>
      </div>
    </div>

    <!-- 图表和最近活动 -->
    <div class="dashboard-grid">
      <!-- 流量统计图表 -->
      <div class="chart-card">
        <div class="card-header">
          <h3>流量统计</h3>
          <div class="chart-controls">
            <button class="chart-btn active">日</button>
            <button class="chart-btn">周</button>
            <button class="chart-btn">月</button>
          </div>
        </div>
        <div class="chart-content">
          <!-- 这里可以集成图表库如Chart.js -->
          <div class="chart-placeholder">
            <div class="chart-bars">
              <div class="chart-bar" :style="{ height: '40%' }"></div>
              <div class="chart-bar" :style="{ height: '65%' }"></div>
              <div class="chart-bar" :style="{ height: '52%' }"></div>
              <div class="chart-bar" :style="{ height: '78%' }"></div>
              <div class="chart-bar" :style="{ height: '85%' }"></div>
              <div class="chart-bar" :style="{ height: '72%' }"></div>
              <div class="chart-bar" :style="{ height: '90%' }"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 最近文章 -->
      <div class="recent-articles-card">
        <div class="card-header">
          <h3>最近文章</h3>
          <router-link to="/admin/articles" class="view-all-link">查看全部</router-link>
        </div>
        <div class="recent-articles-list">
          <div v-for="article in recentArticles" :key="article.id" class="recent-article-item">
            <div class="article-info">
              <h4>{{ article.title }}</h4>
              <p class="article-meta">{{ formatDate(article.created_at) }} · {{ article.status }}</p>
            </div>
            <div class="article-actions">
              <button class="action-btn view-btn">查看</button>
              <button class="action-btn edit-btn">编辑</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { articleApi } from '@/utils/api.js';

const articleCount = ref(28);
const userCount = ref(156);
const viewCount = ref(12543);
const draftCount = ref(5);
const recentArticles = ref([
  {
    id: 1,
    title: "教大家使用阳光巴厘岛到底有没有挂",
    created_at: "2025-10-22T10:30:00Z",
    status: "已发布"
  },
  {
    id: 2,
    title: "科技推荐欢聚斗地主是不是有挂",
    created_at: "2025-10-21T14:20:00Z",
    status: "已发布"
  },
  {
    id: 3,
    title: "分享技术新世界牛牛到底有没有挂",
    created_at: "2025-10-20T09:15:00Z",
    status: "草稿"
  }
]);

onMounted(async () => {
  try {
    // 实际项目中应该从API获取数据
    // const articles = await articleApi.getArticles();
    // articleCount.value = articles.length;
  } catch (error) {
    console.error('获取统计数据失败:', error);
  }
});

const formatDate = (dateString) => {
  const date = new Date(dateString);
  return date.toLocaleDateString();
};
</script>

<style scoped>
.dashboard {
  color: #333;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 15px;
}

.articles-icon {
  background-color: #e8f4fd;
  color: #3498db;
}

.users-icon {
  background-color: #e8f8f5;
  color: #2ecc71;
}

.views-icon {
  background-color: #fef5e7;
  color: #f39c12;
}

.drafts-icon {
  background-color: #fef1f2;
  color: #e74c3c;
}

.stat-content h3 {
  margin: 0 0 5px 0;
  font-size: 14px;
  color: #7f8c8d;
  font-weight: normal;
}

.stat-number {
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: bold;
  color: #2c3e50;
}

.stat-change {
  margin: 0;
  font-size: 12px;
}

.stat-change.positive {
  color: #2ecc71;
}

.stat-change.negative {
  color: #e74c3c;
}

/* 仪表盘网格 */
.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.chart-card,
.recent-articles-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 15px 20px;
  border-bottom: 1px solid #ecf0f1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #2c3e50;
}

/* 图表控制 */
.chart-controls {
  display: flex;
  gap: 5px;
}

.chart-btn {
  padding: 4px 10px;
  border: 1px solid #ddd;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.chart-btn.active {
  background: #3498db;
  color: white;
  border-color: #3498db;
}

.view-all-link {
  font-size: 12px;
  color: #3498db;
  text-decoration: none;
}

/* 图表内容 */
.chart-content {
  padding: 20px;
}

.chart-placeholder {
  height: 250px;
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  width: 100%;
  height: 200px;
}

.chart-bar {
  width: 30px;
  background: linear-gradient(to top, #3498db, #2980b9);
  border-radius: 4px 4px 0 0;
  transition: height 1s ease;
}

/* 最近文章列表 */
.recent-articles-list {
  padding: 0;
}

.recent-article-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #ecf0f1;
}

.recent-article-item:last-child {
  border-bottom: none;
}

.article-info h4 {
  margin: 0 0 5px 0;
  font-size: 14px;
  color: #2c3e50;
}

.article-meta {
  margin: 0;
  font-size: 12px;
  color: #7f8c8d;
}

.article-actions {
  display: flex;
  gap: 5px;
}

.action-btn {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
}

.view-btn {
  background: #ecf0f1;
  color: #2c3e50;
}

.edit-btn {
  background: #3498db;
  color: white;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>