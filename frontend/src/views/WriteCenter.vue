<template>
  <div class="write-center">
    <div class="header glass-effect">
      <div class="header-content">
        <h1 class="display-large">创作中心</h1>
        <div class="header-stats">
          <div class="stat-card">
            <span class="stat-value">{{ formatNumber(totalWordCount) }}</span>
            <span class="stat-label">总字数</span>
          </div>
          <div class="stat-card">
            <span class="stat-value">{{ formatNumber(totalReadCount) }}</span>
            <span class="stat-label">总阅读</span>
          </div>
          <div class="stat-card">
            <span class="stat-value">{{ formatNumber(totalFavorites) }}</span>
            <span class="stat-label">总收藏</span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <el-button type="primary" size="large" class="create-button" @click="createNovel">
          <el-icon><Plus /></el-icon>创建新作品
        </el-button>
        <el-button class="ai-button" @click="openAIDialog">
          <el-icon><Brush /></el-icon>AI助手
        </el-button>
      </div>
    </div>

    <div class="content">
      <!-- 作品列表 -->
      <el-tabs v-model="activeTab" class="novel-tabs">
        <el-tab-pane label="我的作品" name="myNovels">
          <div class="novel-list">
            <div v-for="novel in novels" :key="novel.id" class="novel-card">
              <div class="novel-info">
                <div class="cover-wrapper">
                  <img :src="novel.coverUrl || '/default-cover.jpg'" class="novel-cover" />
                  <div class="cover-overlay"></div>
                </div>
                <div class="novel-details">
                  <h3 class="title-large novel-title">{{ novel.title }}</h3>
                  <div class="novel-meta">
                    <el-dropdown @command="(command) => updateNovelStatus(novel.id, command)">
                      <el-tag :type="getStatusType(novel.status)" class="status-tag" effect="light">
                        {{ getStatusText(novel.status) }}
                        <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                      </el-tag>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item 
                            v-for="option in STATUS_OPTIONS" 
                            :key="option.value"
                            :command="option.value"
                            :disabled="novel.status === option.value"
                          >
                            {{ option.label }}
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                    <span class="update-time">{{ getRelativeTime(novel.updateTime) }}</span>
                  </div>
                  <div class="novel-stats">
                    <div class="stat-item">
                      <el-icon><Document /></el-icon>
                      <span>{{ formatNumber(novel.wordCount) }}字</span>
                    </div>
                    <div class="stat-item">
                      <el-icon><View /></el-icon>
                      <span>{{ formatNumber(novel.readCount) }}阅读</span>
                    </div>
                    <div class="stat-item">
                      <el-icon><Star /></el-icon>
                      <span>{{ formatNumber(novel.favoriteCount) }}收藏</span>
                    </div>
                  </div>
                  <div class="progress-info">
                    <span class="progress-label">最近30天更新进度</span>
                    <el-progress 
                      :percentage="novel.monthProgress || 0"
                      :format="(val) => val + '%'"
                      :stroke-width="8"
                      class="progress-bar"
                    />
                  </div>
                </div>
              </div>
              <div class="novel-actions">                
                <el-button @click="manageChapters(novel.id)">
                  <el-icon><List /></el-icon>章节管理
                </el-button>
                <el-button @click="editNovel(novel.id)">
                  <el-icon><Setting /></el-icon>作品设置
                </el-button>
                <el-dropdown trigger="click" @command="handleCommand($event, novel.id)">
                  <el-button>
                    更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="outline">
                        <el-icon><Connection /></el-icon>大纲管理
                      </el-dropdown-item>
                      <el-dropdown-item command="stats">
                        <el-icon><TrendCharts /></el-icon>数据统计
                      </el-dropdown-item>                    
                      <el-dropdown-item command="delete" divided>
                        <el-icon><Delete /></el-icon>删除作品
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </el-tab-pane>
        <el-tab-pane label="草稿箱" name="drafts">
          <!-- 草稿列表，结构类似上面的作品列表 -->
        </el-tab-pane>
      </el-tabs>

      <!-- 分页器 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalNovels"
          background
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 创建新作品对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建新作品" width="600px" class="custom-dialog">
      <el-form :model="newNovel" :rules="novelRules" ref="novelForm" label-width="100px">
        <el-form-item label="作品名称" prop="title">
          <el-input v-model="newNovel.title" placeholder="请输入作品名称" />
        </el-form-item>
        <el-form-item label="作品分类" prop="category">
          <el-select v-model="newNovel.category" placeholder="请选择分类">
            <el-option
              v-for="item in categories"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="作品简介" prop="description">
          <el-input
            v-model="newNovel.description"
            type="textarea"
            :rows="4"
            placeholder="请输入作品简介"
          />
        </el-form-item>
        <el-form-item label="封面图片">
          <el-upload
            class="cover-uploader"
            action="/api/upload"
            :show-file-list="false"
            :on-success="handleCoverSuccess"
            :before-upload="beforeCoverUpload"
          >
            <img v-if="newNovel.coverUrl" :src="newNovel.coverUrl" class="uploaded-cover" />
            <el-icon v-else class="upload-icon"><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">支持jpg、png格式，建议尺寸400x600px</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="createDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitNovel">确认创建</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- AI助手对话框 -->
    <el-dialog 
      v-model="aiDialogVisible" 
      title="AI创作助手" 
      width="800px"
      class="ai-dialog custom-dialog"
    >
      <div class="ai-chat">
        <div class="chat-messages" ref="chatMessages">
          <div 
            v-for="(message, index) in aiMessages" 
            :key="index"
            :class="['message', message.role]"
          >
            <el-avatar 
              :size="40" 
              :src="message.role === 'assistant' ? '/ai-avatar.png' : '/user-avatar.png'"
            />
            <div class="message-content">
              <div class="message-text" v-html="formatMessage(message.content)"></div>
              <div class="message-time">{{ formatTime(message.timestamp) }}</div>
            </div>
          </div>
        </div>
        <div class="chat-input">
          <el-input
            v-model="aiInput"
            type="textarea"
            :rows="3"
            placeholder="请描述你需要帮助的内容..."
            @keyup.enter.ctrl="sendMessage"
          />
          <div class="input-actions">
            <span class="input-tip">Ctrl + Enter 发送</span>
            <el-button type="primary" @click="sendMessage" :loading="isAiLoading">
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, ArrowDown, Edit, List, Setting, Document, View, Star,
  Delete, Brush, Connection, TrendCharts 
} from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'
import { formatDate, getRelativeTime } from '../utils/date'
import { marked } from 'marked'

const router = useRouter()
const activeTab = ref('myNovels')
const currentPage = ref(1)
const pageSize = ref(10)
const totalNovels = ref(0)
const novels = ref([])
const createDialogVisible = ref(false)

// 新作品表单
const novelForm = ref(null)
const newNovel = reactive({
  title: '',
  category: '',
  description: ''
})

// 表单验证规则
const novelRules = {
  title: [
    { required: true, message: '请输入作品名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择作品分类', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入作品简介', trigger: 'blur' },
    { min: 10, max: 500, message: '长度在 10 到 500 个字符', trigger: 'blur' }
  ]
}

// 分类选项
const categories = [
  { value: 'fantasy', label: '奇幻' },
  { value: 'scifi', label: '科幻' },
  { value: 'romance', label: '言情' },
  { value: 'mystery', label: '悬疑' },
  { value: 'history', label: '历史' },
  { value: 'other', label: '其他' }
]

// 新增状态
const totalWordCount = ref(0)
const totalReadCount = ref(0)
const totalFavorites = ref(0)
const aiDialogVisible = ref(false)
const aiMessages = ref([])
const aiInput = ref('')
const isAiLoading = ref(false)
const chatMessages = ref(null)

// 状态选项
const STATUS_OPTIONS = [
  { value: 0, label: '连载中', type: 'primary' },
  { value: 1, label: '已完结', type: 'success' },
  { value: 2, label: '烂尾', type: 'danger' }
]

// 获取作品列表
const fetchNovels = async () => {
  try {
    const response = await novelApi.getAuthorNovels({
      page: currentPage.value,
      limit: pageSize.value
    })
    
    // 确保response存在且包含novels数组
    if (!response?.novels) {
      console.error('Invalid API response format:', response)
      ElMessage.error('获取作品列表失败：数据格式错误')
      novels.value = []
      totalNovels.value = 0
      return
    }
    
    // Transform the novels data to handle any missing or null values
    novels.value = response.novels.map(novel => ({
      ...novel,
      coverUrl: novel.coverUrl || '/default-cover.jpg',
      wordCount: novel.wordCount || 0,
      readCount: novel.readCount || 0,
      favoriteCount: novel.favoriteCount || 0,
      status: novel.status || 0,
      updateTime: novel.updatedAt, // For compatibility with the template
      author: novel.author || {
        Username: '未知作者',
        Avatar: '/default-avatar.jpg'
      }
    }))
    
    totalNovels.value = response.total || 0
  } catch (error) {
    console.error('获取作品列表失败:', error)
    ElMessage.error('获取作品列表失败')
    novels.value = []
    totalNovels.value = 0
  }
}

// 创建新作品
const createNovel = () => {
  createDialogVisible.value = true
}

// 提交新作品
const submitNovel = async () => {
  if (!novelForm.value) return
  
  await novelForm.value.validate(async (valid) => {
    if (valid) {
      try {
        await novelApi.createNovel(newNovel)
        ElMessage.success('创建成功')
        createDialogVisible.value = false
        fetchNovels()
      } catch (error) {
        ElMessage.error('创建失败')
      }
    }
  })
}

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    0: 'primary',  // 连载中
    1: 'success',  // 已完结
    2: 'danger'    // 烂尾
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    0: '连载中',
    1: '已完结',
    2: '烂尾'
  }
  return texts[status] || '未知'
}

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num
}

// 处理更多操作
const handleCommand = async (command, novelId) => {
  switch (command) {
    case 'outline':
      router.push(`/write/outline/${novelId}`)
      break
    case 'stats':
      router.push(`/write/stats/${novelId}`)
      break
    case 'delete':
      await handleDelete(novelId)
      break
  }
}

// 删除作品
const handleDelete = async (novelId) => {
  try {
    await ElMessageBox.confirm('确定要删除这部作品吗？此操作不可恢复', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await novelApi.deleteNovel(novelId)
    ElMessage.success('删除成功')
    fetchNovels()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 页面跳转方法
const continueWriting = (novelId) => router.push(`/write/editor/${novelId}`)
const manageChapters = (novelId) => router.push(`/write/chapters/${novelId}`)
const editNovel = (novelId) => router.push(`/write/settings/${novelId}`)

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  fetchNovels()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchNovels()
}

// 新增方法
const openAIDialog = () => {
  aiDialogVisible.value = true
  if (aiMessages.value.length === 0) {
    aiMessages.value.push({
      role: 'assistant',
      content: '你好！我是你的AI创作助手。我可以帮你：\n\n- 优化作品情节和人物设定\n- 提供写作技巧和建议\n- 分析作品数据和读者反馈\n- 解答创作过程中的疑问\n\n请告诉我你需要什么帮助？',
      timestamp: Date.now()
    })
  }
}

const sendMessage = async () => {
  if (!aiInput.value.trim() || isAiLoading.value) return
  
  const userMessage = aiInput.value.trim()
  aiMessages.value.push({
    role: 'user',
    content: userMessage,
    timestamp: Date.now()
  })
  aiInput.value = ''
  
  isAiLoading.value = true
  try {
    // 调用DeepSeek API
    const response = await fetch('https://api.deepseek.com/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${import.meta.env.VITE_DEEPSEEK_API_KEY}`
      },
      body: JSON.stringify({
        model: 'deepseek-chat',
        messages: [
          { role: 'system', content: '你是一个专业的小说创作顾问，擅长提供写作建议和优化意见。' },
          { role: 'user', content: userMessage }
        ]
      })
    })
    
    const data = await response.json()
    aiMessages.value.push({
      role: 'assistant',
      content: data.choices[0].message.content,
      timestamp: Date.now()
    })
  } catch (error) {
    ElMessage.error('AI响应失败，请稍后重试')
  } finally {
    isAiLoading.value = false
    await nextTick()
    scrollToBottom()
  }
}

const scrollToBottom = () => {
  if (chatMessages.value) {
    chatMessages.value.scrollTop = chatMessages.value.scrollHeight
  }
}

const formatMessage = (content) => {
  return marked.parse(content)
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleTimeString()
}

// 获取作品统计数据
const fetchNovelStats = async () => {
  try {
    const response = await novelApi.getAuthorStats()
    totalWordCount.value = response.totalWordCount || 0
    totalReadCount.value = response.totalReadCount || 0
    totalFavorites.value = response.totalFavorites || 0
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

// 添加修改状态的方法
const updateNovelStatus = async (novelId, newStatus) => {
  try {
    await novelApi.updateNovelStatus(novelId, newStatus)
    ElMessage.success('状态更新成功')
    await fetchNovels() // 刷新列表
  } catch (error) {
    console.error('更新状态失败:', error)
    ElMessage.error('更新状态失败')
  }
}

onMounted(() => {
  fetchNovels()
  fetchNovelStats()
})
</script>

<style scoped>
.write-center {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  background: var(--system-background);
}

.header {
  margin-bottom: 40px;
  padding: 32px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
  background: var(--system-background);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header h1 {
  margin: 0;
  color: var(--text-primary);
}

.header-stats {
  display: flex;
  gap: 24px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 24px;
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.header-actions {
  display: flex;
  gap: 16px;
}

.create-button {
  padding: 12px 24px;
}

.novel-card {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 24px;
  border: 1px solid var(--border-color);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.novel-card:hover {
  transform: translateY(-2px);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.05),
    0 12px 25px rgba(0, 0, 0, 0.07);
  border-color: var(--apple-blue);
}

.novel-info {
  display: flex;
  gap: 24px;
  margin-bottom: 24px;
}

.cover-wrapper {
  position: relative;
  width: 160px;
  flex-shrink: 0;
  border-radius: 12px;
  overflow: hidden;
}

.novel-cover {
  width: 100%;
  height: 240px;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, transparent 60%, rgba(0, 0, 0, 0.4));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.novel-card:hover .novel-cover {
  transform: scale(1.05);
}

.novel-card:hover .cover-overlay {
  opacity: 1;
}

.novel-details {
  flex: 1;
}

.novel-title {
  margin: 0 0 16px;
  color: var(--text-primary);
}

.novel-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.status-tag {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-tag:hover {
  opacity: 0.8;
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.el-dropdown-menu__item.is-disabled) {
  cursor: not-allowed;
}

.update-time {
  color: var(--text-tertiary);
  font-size: 14px;
}

.novel-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
}

.progress-info {
  margin-top: 20px;
}

.progress-label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-secondary);
  font-size: 14px;
}

.progress-bar {
  width: 100%;
}

.novel-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.novel-actions .el-button {
  display: flex;
  align-items: center;
  gap: 6px;
}

.custom-dialog {
  border-radius: 16px;
}

.cover-uploader {
  width: 200px;
  height: 300px;
  border: 1px dashed var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.upload-icon {
  font-size: 28px;
  color: var(--text-tertiary);
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.uploaded-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-tip {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 8px;
}

.ai-dialog {
  .ai-chat {
    height: 600px;
    display: flex;
    flex-direction: column;
  }

  .chat-messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: var(--system-grouped-secondary);
    border-radius: 12px;
    margin-bottom: 20px;
  }

  .message {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;

    &.assistant {
      .message-content {
        background: var(--system-blue-tint);
      }
    }

    &.user {
      flex-direction: row-reverse;

      .message-content {
        background: var(--system-grouped-primary);
      }
    }
  }

  .message-content {
    flex: 1;
    padding: 16px;
    border-radius: 12px;
    max-width: 80%;
  }

  .message-text {
    color: var(--text-primary);
    line-height: 1.6;

    :deep(p) {
      margin: 0 0 12px;
      &:last-child {
        margin-bottom: 0;
      }
    }

    :deep(code) {
      background: var(--system-grouped-tertiary);
      padding: 2px 6px;
      border-radius: 4px;
    }
  }

  .message-time {
    font-size: 12px;
    color: var(--text-tertiary);
    margin-top: 8px;
  }

  .chat-input {
    background: var(--system-background);
    padding: 16px;
    border-radius: 12px;
  }

  .input-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 12px;
  }

  .input-tip {
    font-size: 12px;
    color: var(--text-tertiary);
  }
}

@media (prefers-color-scheme: dark) {
  .novel-card {
    border-color: rgba(255, 255, 255, 0.1);
  }

  .novel-card:hover {
    border-color: var(--apple-blue);
    box-shadow: 
      0 4px 12px rgba(0, 0, 0, 0.2),
      0 12px 25px rgba(0, 0, 0, 0.3);
  }
}
</style>
