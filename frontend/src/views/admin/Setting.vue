<template>
  <div class="settings-container">
    <h2>系统配置</h2>
    
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error-message">{{ error }}</div>
    <div v-else>
      <form @submit.prevent="saveSettings" class="settings-form">
        <div class="settings-section">
          <h3>网站基本信息</h3>
          
          <div class="form-group">
            <label for="siteName">网站名称</label>
            <input 
              type="text" 
              id="siteName" 
              v-model="settings.siteName" 
              placeholder="请输入网站名称"
              required
            >
          </div>
          
          <div class="form-group">
            <label for="siteDescription">网站描述</label>
            <textarea 
              id="siteDescription" 
              v-model="settings.siteDescription" 
              placeholder="请输入网站描述"
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="siteKeywords">网站关键词</label>
            <input 
              type="text" 
              id="siteKeywords" 
              v-model="settings.siteKeywords" 
              placeholder="请输入网站关键词，用逗号分隔"
            >
          </div>
          
          <div class="form-group">
            <label for="footerText">底部文字</label>
            <textarea 
              id="footerText" 
              v-model="settings.footerText" 
              placeholder="请输入底部版权信息"
              rows="2"
            ></textarea>
          </div>
        </div>
        
        <div class="settings-section">
          <h3>文章设置</h3>
          
          <div class="form-group">
            <label for="articlesPerPage">每页显示文章数量</label>
            <input 
              type="number" 
              id="articlesPerPage" 
              v-model.number="settings.articlesPerPage" 
              min="1" 
              max="50"
              required
            >
          </div>
          
          <div class="form-group">
            <label for="allowComments">允许评论</label>
            <label class="switch">
              <input 
                type="checkbox" 
                v-model="settings.allowComments"
              >
              <span class="slider round"></span>
            </label>
          </div>
          
          <div class="form-group">
            <label for="autoPublish">自动发布</label>
            <label class="switch">
              <input 
                type="checkbox" 
                v-model="settings.autoPublish"
              >
              <span class="slider round"></span>
            </label>
          </div>
        </div>
        
        <div class="settings-section">
          <h3>联系方式</h3>
          
          <div class="form-group">
            <label for="contactEmail">联系邮箱</label>
            <input 
              type="email" 
              id="contactEmail" 
              v-model="settings.contactEmail" 
              placeholder="请输入联系邮箱"
            >
          </div>
          
          <div class="form-group">
            <label for="socialLinks">社交媒体链接</label>
            <div v-for="(link, index) in settings.socialLinks" :key="index" class="social-link-item">
              <select v-model="link.platform" class="social-platform">
                <option value="weixin">微信</option>
                <option value="weibo">微博</option>
                <option value="github">GitHub</option>
                <option value="twitter">Twitter</option>
                <option value="facebook">Facebook</option>
              </select>
              <input 
                type="url" 
                v-model="link.url" 
                placeholder="请输入链接地址"
                class="social-url"
              >
              <button type="button" @click="removeSocialLink(index)" class="remove-btn">删除</button>
            </div>
            <button type="button" @click="addSocialLink" class="add-btn">添加社交媒体</button>
          </div>
        </div>
        
        <div class="form-actions">
          <button type="button" @click="cancelEdit" class="btn btn-secondary">取消</button>
          <button type="submit" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存设置' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { adminApi } from '../../utils/api';
import { ElMessage } from 'element-plus';

const loading = ref(true);
const saving = ref(false);
const error = ref('');
const settings = ref({
  siteName: '',
  siteDescription: '',
  siteKeywords: '',
  footerText: '',
  articlesPerPage: 10,
  allowComments: true,
  autoPublish: false,
  contactEmail: '',
  socialLinks: []
});

const originalSettings = ref({});

// 加载系统设置
onMounted(async () => {
  try {
    loading.value = true;
    const response = await adminApi.getSystemSettings();
    settings.value = response.data || {
      siteName: '',
      siteDescription: '',
      siteKeywords: '',
      footerText: '',
      articlesPerPage: 10,
      allowComments: true,
      autoPublish: false,
      contactEmail: '',
      socialLinks: []
    };
    // 深拷贝原始设置，用于取消操作
    originalSettings.value = JSON.parse(JSON.stringify(settings.value));
  } catch (err) {
    error.value = '加载系统设置失败，请重试';
    console.error('加载系统设置失败:', err);
  } finally {
    loading.value = false;
  }
});

// 保存系统设置
const saveSettings = async () => {
  try {
    saving.value = true;
    await adminApi.updateSystemSettings(settings.value);
    ElMessage.success('系统设置已成功保存');
    // 更新原始设置引用
    originalSettings.value = JSON.parse(JSON.stringify(settings.value));
  } catch (err) {
    ElMessage.error('保存失败，请重试');
    console.error('保存系统设置失败:', err);
  } finally {
    saving.value = false;
  }
};

// 取消编辑
const cancelEdit = () => {
  settings.value = JSON.parse(JSON.stringify(originalSettings.value));
  ElMessage.info('已取消更改');
};

// 添加社交媒体链接
const addSocialLink = () => {
  settings.value.socialLinks.push({
    platform: 'weixin',
    url: ''
  });
};

// 删除社交媒体链接
const removeSocialLink = (index) => {
  settings.value.socialLinks.splice(index, 1);
};
</script>

<style scoped>
.settings-container {
  padding: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 24px;
}

.loading, .error-message {
  padding: 40px;
  text-align: center;
}

.error-message {
  color: #f56c6c;
}

.settings-form {
  max-width: 800px;
}

.settings-section {
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eaeaea;
}

.settings-section:last-child {
  border-bottom: none;
}

h3 {
  margin-bottom: 15px;
  color: #606266;
  font-size: 18px;
}

.form-group {
  margin-bottom: 16px;
}

label {
  display: block;
  margin-bottom: 6px;
  color: #606266;
  font-weight: 500;
}

input[type="text"],
input[type="number"],
input[type="email"],
input[type="url"],
textarea,
select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  transition: border-color 0.2s;
}

input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #409eff;
}

textarea {
  resize: vertical;
  min-height: 60px;
}

/* 开关样式 */
.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: 0.4s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: #409eff;
}

input:checked + .slider:before {
  transform: translateX(24px);
}

/* 社交媒体链接样式 */
.social-link-item {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.social-platform {
  width: 120px;
}

.social-url {
  flex: 1;
}

.add-btn,
.remove-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.add-btn {
  background-color: #f0f9eb;
  color: #67c23a;
  margin-top: 8px;
}

.add-btn:hover {
  background-color: #e1f3d8;
}

.remove-btn {
  background-color: #fef0f0;
  color: #f56c6c;
  white-space: nowrap;
}

.remove-btn:hover {
  background-color: #fde2e2;
}

/* 表单操作按钮 */
.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
  min-width: 80px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #409eff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #66b1ff;
}

.btn-secondary {
  background-color: #f5f7fa;
  color: #606266;
}

.btn-secondary:hover {
  background-color: #e6e8eb;
}
</style>