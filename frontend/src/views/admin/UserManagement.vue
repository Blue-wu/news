<template>
  <div class="user-management">
    <h1>用户管理</h1>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="user-table-container">
      <table class="user-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>昵称</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>状态</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.nickname }}</td>
            <td>{{ user.email }}</td>
            <td>{{ user.role }}</td>
            <td>{{ user.status === 1 ? '活跃' : '禁用' }}</td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>
              <button 
                @click="toggleUserStatus(user.id, user.status)" 
                :disabled="user.role === 'admin'"
                class="action-btn"
              >
                {{ user.status === 1 ? '禁用' : '启用' }}
              </button>
              <button 
                @click="deleteUser(user.id)" 
                :disabled="user.role === 'admin'"
                class="action-btn delete-btn"
              >
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { adminApi } from '@/utils/api.js';
const users = ref([])
const loading = ref(false)
const error = ref('')

onMounted(() => {
  fetchUsers()
})
const fetchUsers = async () => {
  loading.value = true
  error.value = ''
  try {
    // 获取用户列表
    const response = await adminApi.getUsers()
    // 直接访问 response.users，因为响应拦截器已经返回了 response.data
    users.value = response.users
  } catch (err) {
    error.value = err.response?.data?.error || '获取用户列表失败'
  } finally {
    loading.value = false
  }
}

const toggleUserStatus = async (id, currentStatus) => {
  try {
    const newStatus = currentStatus === 1 ? 0 : 1
   // 更新用户状态
await adminApi.updateUserStatus(id, newStatus)
    // 更新本地数据
    const user = users.value.find(u => u.id === id)
    if (user) {
      user.status = newStatus
    }
  } catch (err) {
    error.value = err.response?.data?.error || '更新用户状态失败'
  }
}

const deleteUser = async (id) => {
  if (confirm('确定要删除这个用户吗？')) {
    try {
      // 删除用户
await adminApi.deleteUser(id)
      users.value = users.value.filter(user => user.id !== id)
    } catch (err) {
      error.value = err.response?.data?.error || '删除用户失败'
    }
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString()
}
</script>

<style scoped>
.user-management {
  padding: 20px;
}

.loading,
.error {
  padding: 20px;
  text-align: center;
}

.error {
  color: #ff4757;
}

.user-table-container {
  overflow-x: auto;
}

.user-table {
  width: 100%;
  border-collapse: collapse;
  background: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.user-table th,
.user-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.user-table th {
  background: #f8f8f8;
  font-weight: bold;
}

.action-btn {
  padding: 5px 10px;
  margin-right: 5px;
  cursor: pointer;
  border: none;
  border-radius: 4px;
  background: #646cff;
  color: white;
}

.action-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.delete-btn {
  background: #ff4757;
}
</style>