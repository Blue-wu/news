<script setup>
// 简化App.vue，只作为路由容器
import { onMounted, onUnmounted } from 'vue';

// 回到顶部功能
const handleScroll = () => {
  const backToTop = document.getElementById('back-to-top');
  if (window.scrollY > 300) {
    backToTop?.classList.add('visible');
  } else {
    backToTop?.classList.remove('visible');
  }
};

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
  <div class="app-container">
    <router-view />
    
    <!-- 底部页面 -->
    <footer class="footer">
      <div class="footer-content">
        <div class="footer-text">
          <p>Copyright © 2024 本站由 优品飞游戏网 强力驱动 | XML地图 | HTML地图 | TXT地图</p>
          <p>免责声明:本站所有信息均来自互联网搜集,与产品相关信息的真实准确性均由发布单位及个人负责,拒绝任何人以任何形式在本站发表与中华人民共和国法律相抵触的言论,请大家仔细辨认！并不代表本站观点,本站对此不承担任何相关法律责任！如果发现本站有任何文章侵犯你的权益,请立刻邮件 pxyzx@proton.me 通知给予删除。</p>
          <p>投诉: pxyzx@proton.me | TG:cnmmhm</p>
          <p>蜀ICP备2024052306号-2 | 运行时长: 0.185秒 | 查询信息: 19次</p>
        </div>
        <div class="footer-image">
          <img src="https://picsum.photos/100/100?random=footer" alt="Footer Image" class="footer-img" />
        </div>
      </div>
    </footer>
    
    <!-- 回到顶部按钮 -->
    <button id="back-to-top" class="back-to-top" @click="scrollToTop">
      ↑
    </button>
  </div>
</template>

<style scoped>
/* 重置默认样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Arial', 'Microsoft YaHei', sans-serif;
  color: #333;
  background-color: #f5f5f5;
  line-height: 1.6;
}

.app-container {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ... 现有样式保持不变 ... */

/* 底部样式 */
.footer {
  background-color: #1a237e;
  color: #fff;
  padding: 20px;
  margin-top: auto;
}

.footer-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  max-width: 1200px;
  margin: 0 auto;
}

.footer-text {
  flex: 1;
  font-size: 14px;
  line-height: 1.8;
}

.footer-text p {
  margin: 5px 0;
}

.footer-image {
  margin-left: 20px;
}

.footer-img {
  width: 100px;
  height: 100px;
  border-radius: 8px;
}

/* 回到顶部按钮 */
.back-to-top {
  position: fixed;
  bottom: 30px;
  right: 30px;
  background-color: #1a237e;
  color: white;
  border: none;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  opacity: 0;
  transition: opacity 0.3s;
  z-index: 1000;
}

.back-to-top.visible {
  opacity: 1;
}

.back-to-top:hover {
  background-color: #283593;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .footer-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }
  
  .footer-image {
    margin-left: 0;
    margin-top: 20px;
  }
  
  .back-to-top {
    bottom: 20px;
    right: 20px;
    width: 40px;
    height: 40px;
    font-size: 16px;
  }
}
</style>