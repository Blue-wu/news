<!-- frontend/src/views/admin/SitemapManagement.vue -->
<template>
  <div class="sitemap-management">
    <h2>站点地图管理</h2>
    
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>
    
    <template v-else>
      <div class="card">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>站点地图地址:</label>
            <a :href="sitemapUrl" target="_blank">{{ sitemapUrl }}</a>
          </div>
          <div class="info-item">
            <label>缓存状态:</label>
            <span :class="['status', settings.cacheStatus.isCached ? 'cached' : 'not-cached']">
              {{ settings.cacheStatus.isCached ? '已缓存' : '未缓存' }}
            </span>
          </div>
          <div class="info-item">
            <label>缓存时间:</label>
            <span>{{ formatDate(settings.cacheStatus.cachedAt) }}</span>
          </div>
          <div class="info-item">
            <label>缓存时长:</label>
            <span>{{ formatDuration(settings.cacheDuration) }}</span>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>访问统计</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>总访问次数:</label>
            <span class="stat-value">{{ statistics.visitCount }}</span>
          </div>
          <div class="info-item">
            <label>最后访问时间:</label>
            <span>{{ formatDate(statistics.lastVisitTime) }}</span>
          </div>
        </div>
      </div>

      <div class="card">
        <h3>缓存设置</h3>
        <form @submit.prevent="handleUpdateSettings">
          <div class="form-group">
            <label for="cacheDuration">缓存时长 (秒):</label>
            <input 
              id="cacheDuration"
              type="number" 
              v-model.number="settingsForm.cacheDuration"
              min="60" 
              step="60"
              class="form-control"
            />
            <small class="form-text">最小60秒，推荐86400秒（24小时）</small>
          </div>
          <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
            {{ isSubmitting ? '保存中...' : '保存设置' }}
          </button>
        </form>
      </div>

      <div class="card actions">
        <h3>操作</h3>
        <div class="action-buttons">
          <button 
            @click="handleRegenerateSitemap" 
            class="btn btn-primary" 
            :disabled="isProcessing"
          >
            {{ isProcessing ? '生成中...' : '重新生成站点地图' }}
          </button>
          <button 
            @click="handleClearCache" 
            class="btn btn-warning"
            :disabled="isProcessing"
          >
            清除缓存
          </button>
          <button 
            @click="handleCopySitemapUrl" 
            class="btn btn-secondary"
          >
            复制站点地图URL
          </button>
          <a :href="baiduSubmissionUrl" target="_blank" class="btn btn-info">
            提交到百度
          </a>
          <a :href="googleSubmissionUrl" target="_blank" class="btn btn-info">
            提交到Google
          </a>
        </div>
      </div>
    </template>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

export default {
  name: 'SitemapManagement',
  metaInfo: {
    title: '站点地图管理 - 后台管理'
  },
  setup() {
    const router = useRouter()
    
    const baseUrl = 'http://pjsam.wa541.cn'
    const sitemapUrl = computed(() => `${baseUrl}/sitemap.xml`)
    const baiduSubmissionUrl = 'https://ziyuan.baidu.com/site/index'
    const googleSubmissionUrl = 'https://search.google.com/search-console'
    
    // 本地状态管理
    const settings = ref({
      cacheEnabled: true,
      cacheDuration: 86400,
      cacheStatus: {
        isCached: false,
        cachedAt: '',
        cacheAge: ''
      }
    })
    
    const statistics = ref({
      visitCount: 0,
      lastVisitTime: ''
    })
    
    const isLoading = ref(false)
    const error = ref(null)
    
    // 表单状态
    const settingsForm = reactive({
      cacheEnabled: true,
      cacheDuration: 86400
    })
    
    const isSubmitting = ref(false)
    const isProcessing = ref(false)
    
    // 获取站点地图设置
    const fetchSettings = async () => {
      isLoading.value = true
      try {
        const response = await axios.get('/api/admin/settings/sitemap', {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        settings.value = response.data
        // 更新表单数据
        settingsForm.cacheEnabled = settings.value.cacheEnabled
        settingsForm.cacheDuration = settings.value.cacheDuration
      } catch (error) {
        console.error('获取站点地图设置失败:', error)
        showError('获取站点地图设置失败')
      } finally {
        isLoading.value = false
      }
    }
    
    // 更新站点地图设置
    const handleUpdateSettings = async () => {
      isSubmitting.value = true
      try {
        await axios.put('/api/admin/settings/sitemap', {
          cacheEnabled: settingsForm.cacheEnabled,
          cacheDuration: settingsForm.cacheDuration
        }, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        showSuccess('设置已成功更新')
        // 重新获取最新设置
        fetchSettings()
      } catch (error) {
        console.error('更新设置失败:', error)
        showError('更新设置失败')
      } finally {
        isSubmitting.value = false
      }
    }
    
    // 重新生成站点地图
    const handleRegenerateSitemap = async () => {
      if (!confirm('确定要重新生成站点地图吗？')) return
      
      isProcessing.value = true
      try {
        await axios.post('/api/admin/sitemap/generate', {}, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        showSuccess('站点地图已成功重新生成')
      } catch (error) {
        console.error('重新生成站点地图失败:', error)
        showError('重新生成站点地图失败')
      } finally {
        isProcessing.value = false
      }
    }
    
    // 清除缓存
    const handleClearCache = async () => {
      if (!confirm('确定要清除站点地图缓存吗？')) return
      
      isProcessing.value = true
      try {
        await axios.post('/api/admin/sitemap/clear-cache', {}, {
          headers: {
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          }
        })
        showSuccess('缓存已成功清除')
        // 重新获取最新设置
        fetchSettings()
      } catch (error) {
        console.error('清除缓存失败:', error)
        showError('清除缓存失败')
      } finally {
        isProcessing.value = false
      }
    }
    
    // 复制站点地图URL
    const handleCopySitemapUrl = async () => {
      try {
        await navigator.clipboard.writeText(sitemapUrl.value)
        showSuccess('URL已成功复制到剪贴板')
      } catch (error) {
        console.error('复制失败:', error)
        showError('复制失败，请手动复制')
      }
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }
    
    // 格式化时长
    const formatDuration = (seconds) => {
      if (seconds < 60) return `${seconds}秒`
      
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      
      if (hours > 0) {
        return `${hours}小时${minutes > 0 ? minutes + '分钟' : ''}`
      }
      return `${minutes}分钟`
    }
    
    // 显示成功提示
    const showSuccess = (message) => {
      // 根据项目的提示组件实现
      // 例如：this.$message.success(message)
      alert(`成功: ${message}`)
    }
    
    // 显示错误提示
    const showError = (message) => {
      // 根据项目的提示组件实现
      // 例如：this.$message.error(message)
      alert(`错误: ${message}`)
    }
    
    // 组件挂载时获取数据
    onMounted(() => {
      fetchSettings()
    })
    
    return {
      sitemapUrl,
      baiduSubmissionUrl,
      googleSubmissionUrl,
      settings,
      statistics,
      isLoading,
      error,
      settingsForm,
      isSubmitting,
      isProcessing,
      handleUpdateSettings,
      handleRegenerateSitemap,
      handleClearCache,
      handleCopySitemapUrl,
      formatDate,
      formatDuration
    }
  }
}
</script>

<style scoped>
.sitemap-management {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.card {
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e8e8e8;
}

h2 {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 24px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e8e8;
}

h3 {
  font-size: 18px;
  font-weight: 500;
  color: #555;
  margin-bottom: 16px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: #666;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #4285f4;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f8f9fa;
  border-radius: 4px;
}

.info-item label {
  font-weight: 500;
  color: #555;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: #4285f4;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-control:focus {
  border-color: #4285f4;
  outline: none;
}

.form-text {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #6c757d;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  transition: background-color 0.2s, opacity 0.2s;
  text-decoration: none;
  display: inline-block;
  text-align: center;
  min-width: 120px;
}

.btn-primary {
  background-color: #4285f4;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #3367d6;
}

.btn-primary:disabled {
  background-color: #ccc;
  cursor: not-allowed;
  opacity: 0.7;
}

.btn-warning {
  background-color: #fbbc05;
  color: white;
}

.btn-warning:hover:not(:disabled) {
  background-color: #e6a600;
}

.btn-secondary {
  background-color: #f1f3f4;
  color: #333;
}

.btn-secondary:hover {
  background-color: #e8eaed;
}

.btn-info {
  background-color: #34a853;
  color: white;
}

.btn-info:hover {
  background-color: #2d8f45;
}

.status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status.cached {
  background-color: #e6f4ea;
  color: #34a853;
}

.status.not-cached {
  background-color: #fce8e6;
  color: #ea4335;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>