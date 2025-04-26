<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 左侧个人信息卡片 -->
      <el-col :span="6">
        <el-card class="profile-card">
          <div class="avatar-container">
            <el-avatar 
              :size="120" 
              :src="userInfo.avatar"
              @error="() => userInfo.avatar = '/default-avatar.png'"
              :class="{ 'avatar-loading': uploading }"
            />
            <input
              type="file"
              ref="avatarInput"
              style="display: none"
              accept="image/*"
              @change="handleAvatarChange"
            />
            <el-button 
              text 
              class="change-avatar-btn" 
              @click="triggerAvatarUpload"
              :loading="uploading"
            >
              {{ uploading ? '上传中...' : '更换头像' }}
            </el-button>
          </div>
          <h2 class="username">{{ userInfo.username }}</h2>
          <p class="email">{{ userInfo.email }}</p>
          <div class="user-stats">
            <div class="stat-item">
              <h3>{{ userInfo.novelCount || 0 }}</h3>
              <p>作品数</p>
            </div>
            <div class="stat-item">
              <h3>{{ userInfo.wordCount || 0 }}</h3>
              <p>总字数</p>
            </div>
            <div class="stat-item">
              <h3>{{ userInfo.favoriteCount || 0 }}</h3>
              <p>收藏数</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧内容区域 -->
      <el-col :span="18">
        <el-card>
          <el-tabs v-model="activeTab">
            <!-- 基本信息设置 -->
            <el-tab-pane label="基本信息" name="basic">
              <el-form 
                :model="userForm" 
                :rules="rules"
                ref="userFormRef"
                label-width="100px"
              >
                <el-form-item label="用户名" prop="username">
                  <el-input v-model="userForm.username" />
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="userForm.email" />
                </el-form-item>
                <el-form-item label="个人简介" prop="bio">
                  <el-input 
                    v-model="userForm.bio" 
                    type="textarea" 
                    :rows="4"
                    placeholder="写点什么来介绍自己吧..."
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleUpdateProfile">
                    保存修改
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <!-- 修改密码 -->
            <el-tab-pane label="修改密码" name="password">
              <el-form 
                :model="passwordForm"
                :rules="passwordRules"
                ref="passwordFormRef"
                label-width="100px"
              >
                <el-form-item label="当前密码" prop="oldPassword">
                  <el-input 
                    v-model="passwordForm.oldPassword"
                    type="password"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="新密码" prop="newPassword">
                  <el-input 
                    v-model="passwordForm.newPassword"
                    type="password"
                    show-password
                  />
                </el-form-item>
                <el-form-item label="确认新密码" prop="confirmPassword">
                  <el-input 
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    show-password
                  />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="handleUpdatePassword">
                    修改密码
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <!-- 我的作品 -->
            <el-tab-pane label="我的作品" name="novels">
              <div class="novels-header">
                <el-button type="primary" @click="$router.push('/write')">
                  创建新作品
                </el-button>
              </div>
              <el-table :data="myNovels" style="width: 100%">
                <el-table-column prop="title" label="作品名称" />
                <el-table-column prop="wordCount" label="字数" width="120" />
                <el-table-column prop="updateTime" label="最后更新" width="180" />
                <el-table-column prop="status" label="状态" width="120">
                  <template #default="{ row }">
                    <el-tag :type="row.status === '连载中' ? 'success' : 'info'">
                      {{ row.status }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="操作" width="200">
                  <template #default="{ row }">
                    <el-button text @click="$router.push(`/write/${row.id}`)">
                      继续创作
                    </el-button>
                    <el-button text type="danger" @click="handleDeleteNovel(row.id)">
                      删除
                    </el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import { userApi } from '../api/user'

const userStore = useUserStore()
const activeTab = ref('basic')
const avatarInput = ref(null)

// 修改用户信息的初始化
const userInfo = ref({
  username: userStore.user?.username || '',
  email: userStore.user?.email || '',
  avatar: userStore.user?.avatar || '/default-avatar.png',
  bio: '',
  novelCount: 0,
  wordCount: 0,
  favoriteCount: 0
})

// 表单数据
const userForm = ref({
  username: '',
  email: '',
  bio: ''
})

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 我的作品列表
const myNovels = ref([
  {
    id: 1,
    title: '测试小说',
    wordCount: 10000,
    updateTime: '2024-01-20 15:30',
    status: '连载中'
  }
])

// 表单验证规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const uploading = ref(false)

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    // 从后端获取最新的用户信息
    const { data } = await userApi.getUserInfo()
    userInfo.value = {
      ...userInfo.value,
      ...data,
      avatar: data.avatar || '/default-avatar.png'
    }
    // 更新用户表单
    userForm.value = {
      username: data.username,
      email: data.email,
      bio: data.bio
    }
    // 更新 store 中的用户信息
    userStore.updateUser(data)
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  }
}

// 更新个人信息
const handleUpdateProfile = async () => {
  try {
    await userApi.updateProfile(userForm.value)
    ElMessage.success('个人信息更新成功')
    await fetchUserInfo()
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '更新失败')
  }
}

// 更新密码
const handleUpdatePassword = async () => {
  try {
    await userApi.updatePassword({
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    })
    ElMessage.success('密码修改成功')
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '密码修改失败')
  }
}

// 触发头像上传
const triggerAvatarUpload = () => {
  avatarInput.value.click()
}

// 修改头像上传处理函数
const handleAvatarChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return

  // 验证文件类型
  const allowedTypes = ['image/jpeg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    ElMessage.error('只支持 JPG、PNG、GIF 格式的图片')
    return
  }

  // 验证文件大小（例如限制为 2MB）
  const maxSize = 2 * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error('图片大小不能超过 2MB')
    return
  }

  uploading.value = true
  const formData = new FormData()
  formData.append('avatar', file)

  try {
    const { data } = await userApi.uploadAvatar(formData)
    userInfo.value.avatar = data.avatar_url
    userStore.updateUser({
      ...userStore.user,
      avatar: data.avatar_url
    })
    ElMessage.success('头像上传成功')
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '头像上传失败')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

// 删除小说
const handleDeleteNovel = async (novelId) => {
  try {
    await ElMessageBox.confirm('确定要删除这部作品吗？', '提示', {
      type: 'warning'
    })
    // TODO: 调用删除API
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.profile-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.profile-card {
  text-align: center;
}

.avatar-container {
  margin-bottom: 20px;
  position: relative;
  display: inline-block;
}

.change-avatar-btn {
  margin-top: 10px;
}

.username {
  margin: 10px 0;
  font-size: 20px;
}

.email {
  color: #666;
  margin-bottom: 20px;
}

.user-stats {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
  border-top: 1px solid #eee;
}

.stat-item h3 {
  margin: 0;
  font-size: 20px;
  color: #409EFF;
}

.stat-item p {
  margin: 5px 0 0;
  font-size: 14px;
  color: #666;
}

.novels-header {
  margin-bottom: 20px;
}

.avatar-loading {
  opacity: 0.7;
}

.avatar-container {
  position: relative;
}

.change-avatar-btn {
  position: absolute;
  bottom: -30px;
  left: 50%;
  transform: translateX(-50%);
  white-space: nowrap;
}
</style> 