<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articleApi, adminApi } from '../utils/api.js'
import Quill from 'quill'
import 'quill/dist/quill.snow.css'
const route = useRoute()
const router = useRouter()
const articleId = route.params.id
const isEditMode = !!articleId

// 更新字段与后端模型对应
const title = ref('')
const content = ref('')
const summary = ref('')
const coverImage = ref('')
const status = ref(1) // 1: 已发布, 0: 草稿
const isTop = ref(0) // 是否置顶
const selectedCategoryIds = ref([])
const selectedTagIds = ref([])
const categories = ref([])
const tags = ref([])
// 添加封面图片上传相关变量
const uploadingCover = ref(false)
const coverImageInput = ref(null)

const loading = ref(isEditMode)
const error = ref('')
const successMsg = ref('')

onMounted(async () => {
  // 获取分类和标签数据
  await fetchCategoriesAndTags()
  
  // 如果是编辑模式，获取文章详情
  if (isEditMode) {
    await fetchArticle()
  }
  // 初始化Quill编辑器
  initQuillEditor()
})

// 添加封面图片上传处理函数
const handleCoverImageUpload = async (event) => {
  const file = event.target.files[0]
  if (!file) return
  
  uploadingCover.value = true
  const formData = new FormData()
  formData.append('image', file)
  
  try {
    // 获取token
    const token = localStorage.getItem('token')
    if (!token) {
      error.value = '请先登录'
      return
    }
    
    // 上传图片
    const response = await fetch('/api/upload/image', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    })
    
    if (!response.ok) {
      throw new Error('图片上传失败')
    }
    
    const data = await response.json()
    coverImage.value = data.url // 设置封面图片URL
  } catch (err) {
    console.error('封面图片上传失败:', err)
    error.value = '封面图片上传失败，请重试'
  } finally {
    uploadingCover.value = false
    // 清空文件输入，允许重新选择相同文件
    event.target.value = ''
  }
}

// 添加移除封面图片函数
const removeCoverImage = () => {
  coverImage.value = ''
}
// 找到initQuillEditor函数，修改为以下内容
const initQuillEditor = () => {
  const toolbarOptions = [
    [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
    ['bold', 'italic', 'underline', 'strike'],
    [{ 'color': [] }, { 'background': [] }],
    [{ 'list': 'ordered'}, { 'list': 'bullet' }],
    [{ 'indent': '-1'}, { 'indent': '+1' }],
    [{ 'align': [] }],
    ['link', 'image'],
    ['clean']
  ];
  
  const quill = new Quill('#editor', {
    theme: 'snow',
    modules: {
      toolbar: {
        container: toolbarOptions,
        handlers: {
          // 自定义图片上传处理
          'image': function(value) {
            if (value) {
              // 创建文件输入元素
              const input = document.createElement('input');
              input.setAttribute('type', 'file');
              input.setAttribute('accept', 'image/*');
              input.click();
              
              // 监听文件选择事件
              input.onchange = async function() {
                if (this.files && this.files[0]) {
                  const formData = new FormData();
                  formData.append('image', this.files[0]);
                  
                  try {
                    // 上传图片到服务器
                    const token = localStorage.getItem('token');
                    const response = await fetch('/api/upload/image', {
                      method: 'POST',
                      headers: {
                        'Authorization': `Bearer ${token}`
                      },
                      body: formData
                    });
                    
                    if (!response.ok) {
                      throw new Error('图片上传失败');
                    }
                    
                    const data = await response.json();
                    
                    // 将图片URL插入到编辑器中
                    const range = quill.getSelection(true);
                    quill.insertEmbed(range.index, 'image', data.url);
                    quill.setSelection(range.index + 1);
                  } catch (error) {
                    console.error('图片上传错误:', error);
                    alert('图片上传失败，请重试');
                  }
                }
              };
            } else {
              // 默认为false时，执行默认行为
              quill.format('image', false);
            }
          }
        }
      }
    },
    placeholder: '开始编辑文章内容...'
  });
  
  // 如果有现有内容，设置到编辑器中
  if (content.value) {
    quill.root.innerHTML = content.value;
  }
  
  // 监听内容变化，更新content变量
  quill.on('text-change', () => {
    content.value = quill.root.innerHTML;
  });
}
const fetchCategoriesAndTags = async () => {
  try {
    // 获取所有分类，正确提取data数组
    const response = await adminApi.getCategories()
    categories.value = response.data || []
    
    // 假设API也有获取标签的接口，如果没有则需要添加
    // tags.value = await adminApi.getTags()
  } catch (err) {
    console.error('获取分类和标签失败:', err)
  }
}

const fetchArticle = async () => {
  try {
    // 检查是否已登录并有有效token
    const token = localStorage.getItem('token');
    if (!token) {
      error.value = '请先登录';
      router.push('/login');
      return;
    }
    
    // 由于这是管理员操作，使用adminApi而不是articleApi
    const response = await adminApi.getArticleDetail(articleId);
    const article = response.data || response;
    
    title.value = article.title
    content.value = article.content
    summary.value = article.summary || ''
    coverImage.value = article.cover_image || ''
    status.value = article.status || 1
    isTop.value = article.is_top || 0
    
    // 设置已选分类和标签
    if (article.categories && article.categories.length) {
      selectedCategoryIds.value = article.categories.map(c => c.id)
    }
    if (article.tags && article.tags.length) {
      selectedTagIds.value = article.tags.map(t => t.id)
    }
  } catch (err) {
    error.value = '获取文章详情失败，请检查权限或重新登录';
    console.error('Error fetching article:', err)
    // 认证失败时跳转到登录页
    if (err.response && err.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      router.push('/login');
    }
  } finally {
    loading.value = false
  }
}

const saveArticle = async () => {
  if (!title.value.trim() || !content.value.trim()) {
    error.value = '标题和内容不能为空'
    return
  }

    try {
    // 获取当前登录用户信息
    const userStr = localStorage.getItem('user');
    const user = userStr ? JSON.parse(userStr) : null;
    
    if (!user || !user.id) {
      error.value = '用户未登录或登录信息无效';
      return;
    }
    
    // 组装文章数据
    const articleData = {
      title: title.value,
      content: content.value,
      summary: summary.value,
      cover_image: coverImage.value,
      status: status.value,
      is_top: isTop.value,
      user_id: user.id, // 添加用户ID
      category_ids: selectedCategoryIds.value,
      tag_ids: selectedTagIds.value
    }

    if (isEditMode) {
      await articleApi.updateArticle(articleId, articleData)
      successMsg.value = '文章更新成功'
      
      // 对于编辑模式，如果文章状态为已发布，可以考虑提交更新后的URL
      // 注意：这个操作是异步的，不阻塞主要流程
      if (status.value === 1) {
        const siteDomain = window.location.origin;
        // 异步调用，不要等待结果
        submitToBaidu(siteDomain, articleId).catch(err => {
          console.error('百度URL提交失败，但不影响文章更新:', err);
        });
      }
    } else {
      // 创建文章并获取返回的文章数据，包括ID
      const response = await articleApi.createArticle(articleData)
      successMsg.value = '文章创建成功'
      
      // 重置表单
      title.value = ''
      content.value = ''
      summary.value = ''
      coverImage.value = ''
      status.value = 1
      isTop.value = 0
      selectedCategoryIds.value = []
      selectedTagIds.value = []
      
      // 只有在创建新文章且状态为已发布时才提交给百度收录
      // 注意：这个操作是异步的，不阻塞主要流程
      if (status.value === 1) {
        const siteDomain = window.location.origin;
        const newArticleId = response.id || (response.data && response.data.id);
        // 异步调用，不要等待结果
        submitToBaidu(siteDomain, newArticleId).catch(err => {
          console.error('百度URL提交失败，但不影响文章创建:', err);
        });
      }
    }
    
    // 成功后返回管理员文章管理页面
    setTimeout(() => {
      router.push('/admin/articles')
    }, 1500)
  } catch (err) {
    error.value = isEditMode ? '更新文章失败' : '创建文章失败'
    console.error('Error saving article:', err)
  }
}

// 提交URL到百度收录
const submitToBaidu = async (siteUrl, articleId) => {
  // 仅在开发环境打印信息，生产环境可根据需要调整
  console.log('准备提交URL到百度收录');
  
  // 检查必要参数
  if (!siteUrl) {
    console.error('提交百度收录失败: 缺少网站URL');
    return;
  }
  
  try {
    // 百度站长平台的token，需要替换为实际的值
    // 注意：实际使用前必须替换为真实token
    const baiduToken = '3066aIne0xQw5+jyRZBtmooVoLFEcik1Nr/TDrbjhcblAh+GEJ1iPlgqpYvmmjLOK4t+UTZWqCpQPYZadhtyHSfDkEABt6QfbFblyy9erbaC7gwzT1DpuoDCpKhzlHi9WApUdxGOktPvJrnY35dbul7WDLPnYtoblR8Peogvw3OBZUo=';
    
    // 开发环境可以跳过实际提交，仅打印信息
    // 这样可以避免因为token无效或网络问题影响开发
    if (process.env.NODE_ENV === 'development') {
      console.log('开发环境：跳过实际百度URL提交');
      console.log('将要提交的URL:', articleId ? `${siteUrl}/articles/${articleId}` : `${siteUrl}/articles`);
      return;
    }
    
    // 百度URL提交API
    const baiduApi = `https://data.zz.baidu.com/urls?site=${encodeURIComponent(siteUrl)}&token=${baiduToken}`;
    
    // 构建需要提交的URL
    let submitUrl;
    if (articleId) {
      // 如果有文章ID，提交文章详情页
      submitUrl = `${siteUrl}/articles/${articleId}`;
    } else {
      // 如果没有文章ID，提交文章列表页作为后备
      submitUrl = `${siteUrl}/articles`;
    }
    
    // 发送POST请求到百度API
    const response = await fetch(baiduApi, {
      method: 'POST',
      headers: {
        'Content-Type': 'text/plain'
      },
      body: submitUrl,
      // 添加超时控制，避免长时间等待
      signal: AbortSignal.timeout(5000) // 5秒超时
    });
    
    // 检查响应状态
    if (!response.ok) {
      throw new Error(`HTTP错误，状态码: ${response.status}`);
    }
    
    const data = await response.json();
    console.log('百度收录提交结果:', data);
    
  } catch (error) {
    console.error('提交到百度收录失败:', error);
    // 这里不抛出错误，避免影响文章创建的流程
  }
}
</script>

<template>
  <div class="article-edit-container">
    <h2 class="page-title">{{ isEditMode ? '编辑文章' : '创建新文章' }}</h2>
    
    <div v-if="error" class="error-message">{{ error }}</div>
    <div v-if="successMsg" class="success-message">{{ successMsg }}</div>
    
    <div v-if="loading" class="loading">加载中...</div>
    
    <form v-else @submit.prevent="saveArticle" class="article-form">
      <div class="form-group">
        <label for="title">标题 <span class="required">*</span></label>
        <input type="text" id="title" v-model="title" placeholder="请输入文章标题" required />
      </div>
      
      <div class="form-group">
        <label for="summary">摘要</label>
        <textarea id="summary" v-model="summary" placeholder="请输入文章摘要" rows="3"></textarea>
      </div>
      
      <div class="form-group">
        <label for="coverImageUpload">封面图片</label>
        <input 
          type="file" 
          id="coverImageUpload" 
          accept="image/*" 
          @change="handleCoverImageUpload"
          style="display: none"
          ref="coverImageInput"
        />
        <div class="upload-container">
          <!-- 上传按钮 -->
          <button 
            type="button" 
            class="upload-btn"
            @click="$refs.coverImageInput.click()"
          >
            {{ coverImage ? '更换图片' : '选择图片' }}
          </button>
          <!-- 预览区域 -->
          <div v-if="coverImage" class="image-preview">
            <img :src="coverImage" alt="封面预览" />
            <button 
              type="button" 
              class="remove-btn"
              @click="removeCoverImage"
            >
              移除
            </button>
          </div>
          <!-- 上传中状态 -->
          <div v-if="uploadingCover" class="uploading">上传中...</div>
        </div>
      </div>
      
      <div class="form-row">
        <div class="form-group">
          <label for="status">状态</label>
          <select id="status" v-model="status">
            <option value="1">已发布</option>
            <option value="0">草稿</option>
          </select>
        </div>
        
        <div class="form-group">
          <label for="isTop">置顶</label>
          <select id="isTop" v-model="isTop">
            <option value="0">否</option>
            <option value="1">是</option>
          </select>
        </div>
      </div>
      
      <div class="form-group">
        <label>分类</label>
        <div class="category-select">
          <label v-for="category in categories" :key="category.id" class="checkbox-label">
            <input 
              type="checkbox" 
              :value="category.id" 
              v-model="selectedCategoryIds" 
            />
            {{ category.name }}
          </label>
        </div>
      </div>
      
  <div class="form-group">
    <label for="content">内容 <span class="required">*</span></label>
    <!-- 替换textarea为div -->
    <div id="editor"></div>
    <input type="hidden" v-model="content" required />
  </div>
      
      <div class="form-actions">
        <button type="submit" class="save-btn">{{ isEditMode ? '更新文章' : '创建文章' }}</button>
        <router-link to="/admin/articles" class="cancel-btn">取消</router-link>
      </div>
    </form>
  </div>
</template>

<style scoped>
.article-edit-container {
  padding: 20px;
}

.page-title {
  font-size: 24px;
  margin-bottom: 20px;
  color: #2c3e50;
  border-bottom: 1px solid #e9ecef;
  padding-bottom: 10px;
}
:deep(.ql-container) {
  min-height: 400px;
  font-size: 16px;
  border-radius: 0 0 4px 4px;
}

:deep(.ql-toolbar) {
  border-radius: 4px 4px 0 0;
}
.article-form {
  max-width: 100%;
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 20px;
}

.form-row {
  display: flex;
  gap: 20px;
}

.form-row .form-group {
  flex: 1;
}

label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #495057;
}

.required {
  color: #dc3545;
}

input, textarea, select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  transition: border-color 0.15s ease-in-out;
}

input:focus, textarea:focus, select:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.25);
}

textarea {
  resize: vertical;
  min-height: 200px;
}

.category-select {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin-right: 6px;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.save-btn {
  padding: 10px 24px;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
  transition: background-color 0.2s;
}

.save-btn:hover {
  background-color: #2980b9;
}

.cancel-btn {
  padding: 10px 24px;
  background-color: #f8f9fa;
  color: #495057;
  text-decoration: none;
  border-radius: 4px;
  font-size: 16px;
  border: 1px solid #dee2e6;
  transition: background-color 0.2s;
}

.cancel-btn:hover {
  background-color: #e9ecef;
}

.error-message {
  color: #dc3545;
  background-color: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 15px;
}

.success-message {
  color: #155724;
  background-color: #d4edda;
  border: 1px solid #c3e6cb;
  border-radius: 4px;
  padding: 10px;
  margin-bottom: 15px;
}

.loading {
  padding: 20px;
  text-align: center;
  color: #6c757d;
}

/* 响应式布局 */
@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
  }
  
  .article-edit-container {
    padding: 10px;
  }
  
  .article-form {
    padding: 15px;
  }
}
</style>