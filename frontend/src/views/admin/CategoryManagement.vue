<script setup>
import { ref, onMounted } from 'vue';
import { adminApi } from '@/utils/api'; // 导入封装的adminApi

const categories = ref([]);
const loading = ref(false);
const error = ref('');
const showCreateForm = ref(false);
const showEditForm = ref(false);
const currentCategory = ref({
  id: '',
  name: '',
  slug: '',
  description: '',
  parent_id: 0,
  sort_order: 0
});

// 获取分类列表
const fetchCategories = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await adminApi.getCategories(); // 使用adminApi
    categories.value = response.data || [];
  } catch (err) {
    error.value = '获取分类失败：' + (err.response?.data?.message || err.message);
    console.error('获取分类失败:', err);
  } finally {
    loading.value = false;
  }
};

// 创建分类
const createCategory = async () => {
  try {
    // 确保数据结构与后端API匹配
    const categoryData = {
      name: currentCategory.value.name,
      slug: currentCategory.value.slug,
      description: currentCategory.value.description || '',
      parent_id: currentCategory.value.parent_id,
      sort_order: currentCategory.value.sort_order || 0
    };
    
    await adminApi.createCategory(categoryData); // 使用adminApi
    showCreateForm.value = false;
    resetForm();
    fetchCategories();
  } catch (err) {
    error.value = '创建分类失败：' + (err.response?.data?.message || err.message);
    console.error('创建分类失败:', err);
  }
};

// 更新分类
const updateCategory = async () => {
  try {
    const categoryData = {
      name: currentCategory.value.name,
      slug: currentCategory.value.slug,
      description: currentCategory.value.description || '',
      parent_id: currentCategory.value.parent_id,
      sort_order: currentCategory.value.sort_order || 0
    };
    
    await adminApi.updateCategory(currentCategory.value.id, categoryData); // 使用adminApi
    showEditForm.value = false;
    resetForm();
    fetchCategories();
  } catch (err) {
    error.value = '更新分类失败：' + (err.response?.data?.message || err.message);
    console.error('更新分类失败:', err);
  }
};

// 删除分类
const deleteCategory = async (id) => {
  if (!confirm('确定要删除这个分类吗？')) return;
  
  try {
    await adminApi.deleteCategory(id); // 使用adminApi
    fetchCategories();
  } catch (err) {
    error.value = '删除分类失败：' + (err.response?.data?.message || err.message);
    console.error('删除分类失败:', err);
  }
};

// 打开编辑表单
const openEditForm = (category) => {
  // 确保字段映射正确，可能需要根据后端返回的数据结构进行调整
  currentCategory.value = {
    id: category.id,
    name: category.name,
    slug: category.slug,
    description: category.description || '',
    parent_id: category.parent_id || 0,
    sort_order: category.sort_order || 0
  };
  showEditForm.value = true;
  showCreateForm.value = false;
};

// 打开创建表单
const openCreateForm = () => {
  resetForm();
  showCreateForm.value = true;
  showEditForm.value = false;
};

// 重置表单
const resetForm = () => {
  currentCategory.value = {
    id: '',
    name: '',
    slug: '',
    description: '',
    parent_id: 0,
    sort_order: 0
  };
};

// 自动生成slug
const generateSlug = () => {
  if (currentCategory.value.name) {
    currentCategory.value.slug = currentCategory.value.name
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')
      .replace(/\s+/g, '-');
  }
};

onMounted(() => {
  fetchCategories();
});
</script>

<template>
  <div class="category-management">
    <h2>分类管理</h2>
    
    <button @click="openCreateForm" class="create-btn">添加分类</button>
    
    <!-- 错误提示 -->
    <div v-if="error" class="error">{{ error }}</div>
    
    <!-- 加载状态 -->
    <div v-if="loading" class="loading">加载中...</div>
    
    <!-- 分类列表 -->
    <div v-else class="category-list">
      <table class="category-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>分类名称</th>
            <th>别名</th>
            <th>描述</th>
            <th>父分类</th>
            <th>排序</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="category in categories" :key="category.id">
            <td>{{ category.id }}</td>
            <td>{{ category.name }}</td>
            <td>{{ category.slug }}</td>
            <td class="description-cell">{{ category.description }}</td>
            <td>{{ category.parent_id || '无' }}</td>
            <td>{{ category.sort_order }}</td>
            <td class="actions">
              <button @click="openEditForm(category)" class="edit-btn">编辑</button>
              <button @click="deleteCategory(category.id)" class="delete-btn">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="categories.length === 0" class="empty">暂无分类</div>
    </div>
    
    <!-- 创建分类表单 -->
    <div v-if="showCreateForm" class="form-overlay">
      <div class="form-container">
        <h3>创建分类</h3>
        <form @submit.prevent="createCategory">
          <div class="form-group">
            <label for="name">分类名称 *</label>
            <input 
              id="name" 
              v-model="currentCategory.name" 
              @input="generateSlug"
              required
              placeholder="请输入分类名称"
            >
          </div>
          
          <div class="form-group">
            <label for="slug">别名 *</label>
            <input 
              id="slug" 
              v-model="currentCategory.slug" 
              required
              placeholder="URL友好的别名"
            >
          </div>
          
          <div class="form-group">
            <label for="description">描述</label>
            <textarea 
              id="description" 
              v-model="currentCategory.description"
              placeholder="请输入分类描述"
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="parent_id">父分类</label>
            <select id="parent_id" v-model.number="currentCategory.parent_id">
              <option :value="0">顶级分类</option>
              <option 
                v-for="cat in categories.filter(c => !currentCategory.id || c.id !== currentCategory.id)" 
                :key="cat.id" 
                :value="cat.id"
              >
                {{ cat.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="sort_order">排序</label>
            <input 
              id="sort_order" 
              v-model.number="currentCategory.sort_order" 
              type="number"
              placeholder="排序序号"
            >
          </div>
          
          <div class="form-actions">
            <button type="submit" class="submit-btn">创建</button>
            <button type="button" @click="showCreateForm = false" class="cancel-btn">取消</button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- 编辑分类表单 -->
    <div v-if="showEditForm" class="form-overlay">
      <div class="form-container">
        <h3>编辑分类</h3>
        <form @submit.prevent="updateCategory">
          <div class="form-group">
            <label for="edit-name">分类名称 *</label>
            <input 
              id="edit-name" 
              v-model="currentCategory.name" 
              @input="generateSlug"
              required
              placeholder="请输入分类名称"
            >
          </div>
          
          <div class="form-group">
            <label for="edit-slug">别名 *</label>
            <input 
              id="edit-slug" 
              v-model="currentCategory.slug" 
              required
              placeholder="URL友好的别名"
            >
          </div>
          
          <div class="form-group">
            <label for="edit-description">描述</label>
            <textarea 
              id="edit-description" 
              v-model="currentCategory.description"
              placeholder="请输入分类描述"
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="edit-parent_id">父分类</label>
            <select id="edit-parent_id" v-model.number="currentCategory.parent_id">
              <option :value="0">顶级分类</option>
              <option 
                v-for="cat in categories.filter(c => c.id !== currentCategory.id)" 
                :key="cat.id" 
                :value="cat.id"
              >
                {{ cat.name }}
              </option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="edit-sort_order">排序</label>
            <input 
              id="edit-sort_order" 
              v-model.number="currentCategory.sort_order" 
              type="number"
              placeholder="排序序号"
            >
          </div>
          
          <div class="form-actions">
            <button type="submit" class="submit-btn">保存</button>
            <button type="button" @click="showEditForm = false" class="cancel-btn">取消</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.category-management {
  padding: 20px;
}

h2 {
  margin-bottom: 20px;
  color: #333;
}

.create-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  margin-bottom: 20px;
}

.create-btn:hover {
  background-color: #45a049;
}

.error {
  color: #f44336;
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

.category-table {
  width: 100%;
  border-collapse: collapse;
  background-color: white;
  border-radius: 4px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.category-table th, .category-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.category-table th {
  background-color: #f5f5f5;
  font-weight: 600;
  color: #333;
}

.category-table tr:hover {
  background-color: #f9f9f9;
}

.description-cell {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.actions {
  display: flex;
  gap: 5px;
}

.edit-btn {
  background-color: #2196F3;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.edit-btn:hover {
  background-color: #0b7dda;
}

.delete-btn {
  background-color: #f44336;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 4px;
  cursor: pointer;
}

.delete-btn:hover {
  background-color: #d32f2f;
}

.empty {
  text-align: center;
  padding: 40px;
  color: #666;
  background-color: white;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

/* 表单样式 */
.form-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.form-container {
  background-color: white;
  padding: 30px;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.2);
}

.form-container h3 {
  margin-bottom: 20px;
  color: #333;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: 500;
  color: #333;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group textarea {
  resize: vertical;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: #4CAF50;
}

.form-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.submit-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.submit-btn:hover {
  background-color: #45a049;
}

.cancel-btn {
  background-color: #ccc;
  color: #333;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
}

.cancel-btn:hover {
  background-color: #999;
}
</style>