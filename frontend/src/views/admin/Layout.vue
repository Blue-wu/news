<template>
  <div class="admin-layout">
    <!-- 侧边导航栏 - 仅在不隐藏时渲染 -->
    <aside v-if="!sidebarHidden" class="sidebar">
      <div class="sidebar-header">
        <h2>管理后台</h2>
      </div>
    <nav class="sidebar-nav">
      <router-link to="/admin" class="nav-item" exact>
        <i class="icon-home"></i>
        <span>仪表盘</span>
      </router-link>
      <router-link to="/admin/articles" class="nav-item">
        <i class="icon-article"></i>
        <span>文章管理</span>
      </router-link>
      <router-link to="/admin/users" class="nav-item">
        <i class="icon-user"></i>
        <span>用户管理</span>
      </router-link>
      <router-link to="/admin/categories" class="nav-item">
        <i class="icon-category"></i>
        <span>分类管理</span>
      </router-link>
      <router-link to="/admin/settings" class="nav-item">
        <i class="icon-settings"></i>
        <span>系统设置</span>
      </router-link>
    </nav>
  </aside>
  
  <!-- 侧边栏背景遮罩 - 仅在显示侧边栏时显示 -->
  <div v-if="!sidebarHidden" class="sidebar-overlay" @click="toggleSidebar"></div>
  
  <!-- 主内容区域 - 直接根据侧边栏状态设置不同样式 -->
  <main class="main-content" :class="{ 'sidebar-visible': !sidebarHidden }">
    <!-- 顶部工具栏 -->
    <header class="main-header">
      <div class="header-left">
        <button class="toggle-sidebar-btn" @click="toggleSidebar" :title="sidebarHidden ? '显示侧边栏' : '隐藏侧边栏'">
          <i class="icon-menu"></i>
        </button>
        <h1>{{ pageTitle }}</h1>
      </div>
      <div class="header-right">
        <div class="user-info">
          <!-- 添加空值检查 -->
          <span class="username" v-if="user">{{ user.username }}</span>
          <button @click="handleLogout" class="logout-btn">退出</button>
        </div>
      </div>
    </header>
    
    <!-- 页面内容 -->
    <div class="content-wrapper">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </div>
  </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';

const route = useRoute();
const router = useRouter();
const sidebarHidden = ref(false); // 初始默认显示侧边栏
const user = ref(null);

onMounted(() => {
  // 从localStorage获取用户信息
  const userStr = localStorage.getItem('user');
  if (userStr) {
    user.value = JSON.parse(userStr);
  }
});

const pageTitle = computed(() => {
  const titles = {
    '/admin': '仪表盘',
    '/admin/articles': '文章管理',
    '/admin/users': '用户管理',
    '/admin/categories': '分类管理',
    '/admin/settings': '系统设置'
  };
  return titles[route.path] || '管理后台';
});

// 控制侧边栏显示/隐藏
const toggleSidebar = () => {
  sidebarHidden.value = !sidebarHidden.value;
};

const handleLogout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  router.push('/login');
};
</script>

<style scoped>
/* 完全重写布局样式，使用更简单直接的方法 */
.admin-layout {
  display: flex;
  height: 100vh;
  background-color: #f5f7fa;
  position: relative;
  overflow: hidden;
}

/* 侧边栏样式 */
.sidebar {
  width: 250px;
  background-color: #2c3e50;
  color: white;
  height: 100vh;
  position: fixed;
  top: 0;
  left: 0;
  z-index: 1001; /* 确保侧边栏在最上层 */
  transition: transform 0.3s ease;
  transform: translateX(0);
  box-shadow: 2px 0 5px rgba(0, 0, 0, 0.2);
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid #34495e;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 20px;
}

.sidebar-nav {
  padding: 10px 0;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  color: #bdc3c7;
  text-decoration: none;
  transition: all 0.3s ease;
  min-height: 48px;
  box-sizing: border-box;
}

.nav-item:hover,
.nav-item.router-link-active {
  background-color: #34495e;
  color: white;
}

.nav-item i {
  margin-right: 10px;
  font-size: 20px;
  width: 24px;
  text-align: center;
}

/* 侧边栏背景遮罩 */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  cursor: pointer;
  transition: opacity 0.3s ease;
}

/* 主内容区域 - 完全重写，避免选择器优先级问题 */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  margin-left: 0; /* 默认无侧边栏时占据整个宽度 */
  transition: none; /* 不需要过渡效果，直接显示 */
  background-color: #f5f7fa;
  overflow: hidden;
}

/* 当侧边栏可见时，添加左边距 */
.main-content.sidebar-visible {
  margin-left: 250px;
  width: calc(100% - 250px);
}

/* 顶部工具栏 */
.main-header {
  background-color: white;
  padding: 0 20px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 10; /* 确保在内容之上 */
}

.header-left {
  display: flex;
  align-items: center;
}

/* 侧边栏切换按钮 - 提高可见性 */
.toggle-sidebar-btn {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  margin-right: 15px;
  color: #666;
  min-height: 48px;
  padding: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  z-index: 1002; /* 确保按钮在最上层 */
  position: relative;
}

.toggle-sidebar-btn:hover {
  background-color: #f0f0f0;
}

.toggle-sidebar-btn:active {
  background-color: #e0e0e0;
}

.header-left h1 {
  margin: 0;
  font-size: 20px;
  color: #333;
}

.user-info {
  display: flex;
  align-items: center;
}

.username {
  margin-right: 15px;
  color: #666;
}

.logout-btn {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 4px;
  cursor: pointer;
  min-height: 44px;
  font-size: 16px;
  transition: background-color 0.3s ease;
}

.logout-btn:hover {
  background-color: #c0392b;
}

.logout-btn:active {
  background-color: #a93226;
  transform: scale(0.98);
}

/* 内容包装器 - 确保可滚动 */
.content-wrapper {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  position: relative;
  z-index: 1;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 图标样式 */
.icon-home::before { content: "🏠"; }
.icon-article::before { content: "📝"; }
.icon-user::before { content: "👥"; }
.icon-category::before { content: "📁"; }
.icon-settings::before { content: "⚙️"; }
.icon-menu::before { content: "☰"; }

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    width: 240px;
  }
  
  .main-content.sidebar-visible {
    margin-left: 240px;
    width: calc(100% - 240px);
  }
  
  .main-header {
    padding: 0 15px;
    height: 56px;
  }
  
  .header-left h1 {
    font-size: 18px;
  }
  
  .username {
    display: none;
  }
  
  .content-wrapper {
    padding: 15px;
  }
}

@media (max-width: 480px) {
  .sidebar {
    width: 220px;
  }
  
  .main-content.sidebar-visible {
    margin-left: 220px;
    width: calc(100% - 220px);
  }
  
  .nav-item {
    padding: 14px 16px;
  }
  
  .sidebar-header {
    padding: 16px;
  }
  
  .content-wrapper {
    padding: 12px;
  }
}
</style>