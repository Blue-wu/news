<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { authApi } from '../../utils/api.js';

const username = ref('');
const password = ref('');
const error = ref('');
const loading = ref(false);
const router = useRouter();

const handleLogin = async () => {
  if (!username.value || !password.value) {
    error.value = '请输入用户名和密码';
    return;
  }

  try {
    loading.value = true;
    error.value = '';

    const response = await authApi.login({
      username: username.value,
      password: password.value
    });

    // 保存token和用户信息到localStorage
    localStorage.setItem('token', response.token);
    localStorage.setItem('user', JSON.stringify(response.user));

    // 根据用户角色进行不同的跳转
    if (response.user.role === 'admin') {
      // 管理员用户跳转到管理后台
      router.push('/admin');
    } else {
      // 普通用户跳转到首页
      router.push('/');
    }
  } catch (err) {
    error.value = err.response?.data?.error || '登录失败，请重试';
  } finally {
    loading.value = false;
  }
};
</script>
<template>
	<div class="login-container">
		<h1>登录</h1>
		<div v-if="error" class="error-message">{{ error }}</div>
		<form @submit.prevent="handleLogin">
			<div class="form-group">
				<label for="username">用户名</label>
				<input
					id="username"
					v-model="username"
					type="text"
					placeholder="请输入用户名"
					required
				/>
			</div>
			<div class="form-group">
				<label for="password">密码</label>
				<input
					id="password"
					v-model="password"
					type="password"
					placeholder="请输入密码"
					required
				/>
			</div>
			<button type="submit" :disabled="loading" class="login-button">
				{{ loading ? '登录中...' : '登录' }}
			</button>
		</form>
	</div>
</template>

<style scoped>
.login-container {
	max-width: 400px;
	margin: 0 auto;
	padding: 2rem;
	border: 1px solid #ddd;
	border-radius: 8px;
	margin-top: 4rem;
	/* 添加阴影提升视觉效果 */
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	/* 确保在移动设备上内容不会超出屏幕 */
	box-sizing: border-box;
}
	h1 {
		text-align: center;
		margin-bottom: 2rem;
		/* 调整标题大小以适应移动端 */
		font-size: 1.8rem;
	}

.form-group {
	margin-bottom: 1.5rem;
}
	label {
		display: block;
		margin-bottom: 0.5rem;
		font-weight: bold;
		/* 增加标签字体大小以提高可读性 */
		font-size: 0.95rem;
	}
	input {
		width: 100%;
		padding: 0.75rem;
		border: 1px solid #ddd;
		border-radius: 4px;
		font-size: 1rem;
		/* 增加输入框高度以方便触摸操作 */
		min-height: 48px;
		box-sizing: border-box;
		/* 添加焦点状态样式 */
		outline: none;
	}
	input:focus {
		border-color: #646cff;
		box-shadow: 0 0 0 2px rgba(100, 108, 255, 0.1);
	}

.login-button {
	width: 100%;
	padding: 0.75rem;
	background-color: #646cff;
	color: white;
	border: none;
	border-radius: 4px;
	font-size: 1rem;
	cursor: pointer;
	transition: background-color 0.3s;
	/* 增加按钮高度以方便触摸操作 */
	min-height: 48px;
	/* 添加触摸反馈 */
	user-select: none;
}
	.login-button:hover {
		background-color: #535bf2;
	}
	.login-button:active {
		/* 添加点击效果 */
		background-color: #4338ca;
		transform: translateY(1px);
	}
	.login-button:disabled {
		background-color: #ccc;
		cursor: not-allowed;
	}

.error-message {
	background-color: #fee;
	color: #c33;
	padding: 0.75rem;
	border-radius: 4px;
	margin-bottom: 1.5rem;
	/* 增加错误消息的可见性 */
	font-size: 0.9rem;
}

/* 响应式设计 - 针对移动设备 */
@media (max-width: 768px) {
	.login-container {
		max-width: 100%;
		margin-top: 2rem;
		margin-left: 1rem;
		margin-right: 1rem;
		padding: 1.5rem;
		/* 移除移动端的边框，减少视觉干扰 */
		border: none;
		box-shadow: 0 1px 5px rgba(0, 0, 0, 0.1);
	}
		h1 {
			font-size: 1.5rem;
			margin-bottom: 1.5rem;
		}
	
	.form-group {
		margin-bottom: 1.25rem;
	}
		input {
			font-size: 1.05rem; /* 稍大字体提高可读性 */
			padding: 0.85rem;  /* 增加内边距 */
		}
	
	.login-button {
		font-size: 1.05rem;
		font-weight: 500;
	}
	
	.error-message {
		padding: 0.85rem;
		font-size: 0.95rem;
	}
}

/* 针对小屏幕手机的额外优化 */
@media (max-width: 480px) {
	.login-container {
		margin-left: 0.75rem;
		margin-right: 0.75rem;
		padding: 1.25rem;
	}
	
	/* 调整整体页面边距 */
	body {
		margin: 0;
		padding: 0;
	}
}
</style>