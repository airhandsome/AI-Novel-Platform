<template>
  <div class="header">
    <div class="header-content">
      <div class="logo">
        <router-link to="/">Vision AI小说平台</router-link>
      </div>
      
      <div class="nav-menu">
        <el-menu
          mode="horizontal"
          :router="true"
          :default-active="route.path"
          class="transparent-menu"
        >
          <el-menu-item index="/">
            <el-icon><House /></el-icon>
            首页
          </el-menu-item>
          <el-menu-item index="/library">
            <el-icon><Reading /></el-icon>
            书库
          </el-menu-item>
          <el-menu-item index="/ranking">
            <el-icon><Trophy /></el-icon>
            排行榜
          </el-menu-item>
          <el-menu-item index="/write">
            <el-icon><Edit /></el-icon>
            创作中心
          </el-menu-item>
        </el-menu>
      </div>

      <div class="user-area">
          <template v-if="userStore.token">
            <el-dropdown @command="handleCommand">
              <div class="user-avatar">
                <el-avatar :size="32" :src="userStore.user?.avatar || '/default-avatar.png'" />
                <span class="username">{{ userStore.user?.username }}</span>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="settings">设置</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="router.push('/login')">登录</el-button>
            <el-button @click="router.push('/register')">注册</el-button>
          </template>
        </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { House, Reading, Trophy, Edit } from '@element-plus/icons-vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)
const userName = computed(() => userStore.userName)
const userAvatar = computed(() => userStore.avatar || '/default-avatar.png')

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      handleLogout()
      break
  }
}

const handleLogout = async () => {
  try {
    await userStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  } catch (error) {
    console.error('退出登录失败:', error)
  }
}
</script>

<style scoped>
.header {
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background-color: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(10px);
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
  color: #fff;
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
  color: #fff;
}

:deep(.transparent-menu) {
  background-color: transparent;
  border-bottom: none;
}

:deep(.el-menu--horizontal) {
  border-bottom: none;
}

:deep(.el-menu-item) {
  font-size: 15px;
  color: #fff !important;
}

:deep(.el-menu-item.is-active) {
  color: #409EFF !important;
  background-color: rgba(255, 255, 255, 0.1) !important;
}

:deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

:deep(.el-menu-item .el-icon) {
  color: #fff;
  margin-right: 4px;
}

:deep(.el-button) {
  padding: 8px 16px;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 16px;
}
</style> 