<template>
  <div class="layout">

    <!-- 主要内容区域 -->
    <el-main class="main-content">
      <router-view />
    </el-main>

    <!-- 页脚 -->
    <el-footer class="footer">
      <div class="footer-content">
        <p>&copy; 2024 AI小说平台 All Rights Reserved.</p>
      </div>
    </el-footer>
  </div>
</template>

<script setup>
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { ElMessage } from 'element-plus'
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      userStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
      break
  }
}
</script>

<style scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background-color: #fff;
  position: fixed;
  width: 100%;
  z-index: 100;
  padding: 0;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.logo {
  font-size: 20px;
  font-weight: bold;
}

.logo a {
  color: #409EFF;
  text-decoration: none;
}

.nav-menu {
  flex: 1;
  margin: 0 40px;
}

.user-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  font-size: 14px;
}

.main-content {
  flex: 1;
  padding-top: 80px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 120px);
}

.footer {
  background-color: #fff;
  padding: 20px 0;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  text-align: center;
  color: #666;
}

:deep(.el-menu) {
  border-bottom: none;
}

:deep(.el-menu--horizontal) {
  border-bottom: none;
}

:deep(.el-menu-item) {
  font-size: 15px;
}
</style>