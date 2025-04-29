<template>
  <div class="novel-detail">
    <!-- 小说基本信息 -->
    <div class="novel-header glass-effect">
      <div class="novel-info">
        <div class="cover-wrapper">
          <img :src="novel.coverUrl || '../../public/assets/default-cover.png'" class="novel-cover" />
          <div class="cover-overlay"></div>
        </div>
        <div class="novel-meta">
          <h1 class="display-large novel-title">{{ novel.title }}</h1>
          <div class="novel-author">
            <el-avatar :size="32" :src="novel.authorAvatar" />
            <span class="title-medium">{{ novel.authorName }}</span>
          </div>
          <div class="novel-stats">
            <span class="stat-item">
              <el-icon><Document /></el-icon>
              <span class="body-medium">{{ formatNumber(novel.wordCount) }}字</span>
            </span>
            <span class="stat-item">
              <el-icon><View /></el-icon>
              <span class="body-medium">{{ formatNumber(novel.readCount) }}次阅读</span>
            </span>
            <span class="stat-item">
              <el-icon><Star /></el-icon>
              <span class="body-medium">{{ formatNumber(novel.favoriteCount) }}收藏</span>
            </span>
          </div>
          <div class="novel-tags">
            <el-tag class="custom-tag" effect="plain">{{ novel.category }}</el-tag>
            <el-tag 
              class="custom-tag" 
              effect="plain"
              :type="novel.status === 2 ? 'success' : 'warning'"
            >
              {{ novel.status === 2 ? '已完结' : '连载中' }}
            </el-tag>
          </div>
          <div class="novel-actions">
            <el-button type="primary" size="large" class="action-button" @click="startReading">
              开始阅读
            </el-button>
            <el-button 
              :type="isFavorited ? 'success' : ''" 
              size="large"
              class="action-button"
              @click="toggleFavorite"
            >
              <el-icon><Star /></el-icon>
              {{ isFavorited ? '已收藏' : '收藏' }}
            </el-button>
          </div>
        </div>
      </div>
      <div class="novel-description">
        <h3 class="title-large">作品简介</h3>
        <p class="body-medium">{{ novel.description }}</p>
      </div>
    </div>

    <!-- 章节列表 -->
    <div class="chapter-list glass-effect">
      <div class="chapter-header">
        <h2 class="display-small">章节列表</h2>
        <div class="chapter-filters">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索章节"
            clearable
            class="search-input"
            @input="filterChapters"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-switch
            v-model="reverseOrder"
            active-text="倒序"
            inactive-text="正序"
            class="custom-switch"
            @change="toggleOrder"
          />
        </div>
      </div>

      <div class="chapters-list">
        <div v-for="chapter in displayChapters" :key="chapter.id" class="chapter-item">
          <div class="chapter-info">
            <div class="chapter-title">
              {{ chapter.title }}
            </div>
            <div class="chapter-meta">
              <span>{{ chapter.formattedWordCount }}</span>
              <span>更新于: {{ chapter.updateTime }}</span>
              <span :class="['chapter-status', chapter.statusClass]">
                {{ chapter.statusText }}
              </span>
            </div>
          </div>
          <div class="chapter-actions">
            <el-button 
              type="primary" 
              size="small" 
              @click="handleReadChapter(chapter.id)"
            >
              阅读
            </el-button>
            <el-button 
              v-if="isAuthor" 
              type="warning" 
              size="small" 
              @click="handleEditChapter(chapter.id)"
            >
              编辑
            </el-button>
          </div>
        </div>
        
        <!-- 空状态 -->
        <el-empty 
          v-if="displayChapters.length === 0" 
          description="暂无章节" 
          :image-size="200"
        >
          <template #image>
            <img src="../../public/assets/empty-chapter.png" alt="暂无章节" style="width: 200px;" />
          </template>
          <el-button v-if="isAuthor" type="primary" @click="handleCreateChapter">
            创建新章节
          </el-button>
        </el-empty>
      </div>

      <!-- 分页器 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="totalChapters"
          :layout="'total, sizes, prev, pager, next, jumper'"
          :background="true"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Document, View, Star, Search } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'
import { useUserStore } from '../stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

// 小说信息
const novel = ref({})
const chapters = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const totalChapters = ref(0)
const searchKeyword = ref('')
const reverseOrder = ref(false)
const isFavorited = ref(false)

// 检查当前用户是否是作者
const isAuthor = computed(() => {
  return userStore.user?.id === novel.value.authorId
})

// 获取小说信息
const fetchNovelDetail = async () => {
  try {
    const data = await novelApi.getNovelDetail(route.params.id)
    novel.value = data
    
    // 检查是否已收藏
    try {
      const response = await novelApi.checkFavorite(route.params.id)
      isFavorited.value = response.message === 'true'
    } catch (error) {
      console.error('获取收藏状态失败:', error)
      isFavorited.value = false
    }
  } catch (error) {
    ElMessage.error('获取小说信息失败')
  }
}

// 获取章节列表
const fetchChapters = async () => {
  try {
    const data  = await novelApi.getNovelChapters({
      novelId: route.params.id,
      page: currentPage.value,
      limit: pageSize.value,
      keyword: searchKeyword.value,
      reverse: reverseOrder.value
    })
    
    // 处理章节数据
    chapters.value = data.chapters.map(chapter => ({
      ...chapter,
      // 解码 HTML 实体
      content: decodeHTMLEntities(chapter.content),
      // 格式化时间
      createdAt: formatDate(chapter.createdAt),
      updateTime: formatDate(chapter.updateTime),
      // 格式化字数
      formattedWordCount: formatWordCount(chapter.wordCount),
      // 状态文本和样式
      statusText: getStatusText(chapter.status),
      statusClass: getStatusClass(chapter.status)
    }))
    
    totalChapters.value = data.total
  } catch (error) {
    ElMessage.error('获取章节列表失败')
  }
}

// HTML 实体解码函数
const decodeHTMLEntities = (text) => {
  const textarea = document.createElement('textarea')
  textarea.innerHTML = text
  return textarea.value
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 格式化字数
const formatWordCount = (count) => {
  if (count >= 10000) {
    return (count / 10000).toFixed(1) + '万字'
  } else if (count >= 1000) {
    return (count / 1000).toFixed(1) + '千字'
  }
  return count + '字'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    0: '草稿',
    1: '已发布'
  }
  return statusMap[status] || '未知'
}

// 获取状态样式类名
const getStatusClass = (status) => {
  const classMap = {
    0: 'status-draft',
    1: 'status-published'
  }
  return classMap[status] || 'status-unknown'
}

// 过滤显示的章节
const displayChapters = computed(() => {
  let filtered = chapters.value
  
  // 搜索过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(chapter => 
      chapter.title.toLowerCase().includes(keyword) ||
      chapter.content.toLowerCase().includes(keyword)
    )
  }
  
  // 排序
  if (reverseOrder.value) {
    filtered = [...filtered].reverse()
  }
  
  return filtered
})

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num
}

// 开始阅读
const startReading = () => {
  if (chapters.value.length > 0) {
    goToChapter(chapters.value[0].id)
  }
}

// 跳转到章节
const goToChapter = (chapterId) => {
  router.push(`/novels/${route.params.id}/chapter/${chapterId}`)
}

// 收藏/取消收藏
const toggleFavorite = async () => {
  try {
    if (isFavorited.value) {
      await novelApi.unfavoriteNovel(route.params.id)
      ElMessage.success('已取消收藏')
      novel.value.favoriteCount--
    } else {
      await novelApi.favoriteNovel(route.params.id)
      ElMessage.success('收藏成功')
      novel.value.favoriteCount++
    }
    isFavorited.value = !isFavorited.value
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 搜索章节
const filterChapters = () => {
  currentPage.value = 1
  fetchChapters()
}

// 切换排序
const toggleOrder = () => {
  fetchChapters()
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  fetchChapters()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchChapters()
}

// 创建新章节
const handleCreateChapter = () => {
  router.push(`/novels/${route.params.id}/chapter/create`)
}

// 加载指定章节
const handleReadChapter = (chapterId) => {
  goToChapter(chapterId)
}

// 编辑章节
const handleEditChapter = (chapterId) => {
  router.push(`/novels/${route.params.id}/chapter/${chapterId}/edit`)
}

onMounted(() => {
  fetchNovelDetail()
  fetchChapters()
})
</script>

<style scoped>
.novel-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  min-height: 100vh;
  background: var(--system-background);
}

.novel-header {
  margin-bottom: 40px;
  border-radius: 24px;
  padding: 40px;
  border: 1px solid var(--border-color);
}

.novel-info {
  display: flex;
  gap: 40px;
  margin-bottom: 40px;
}

.cover-wrapper {
  position: relative;
  width: 240px;
  flex-shrink: 0;
}

.novel-cover {
  width: 100%;
  height: 340px;
  object-fit: cover;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transition: transform 0.3s ease;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  border-radius: 16px;
  background: linear-gradient(to bottom, transparent 60%, rgba(0, 0, 0, 0.4));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.cover-wrapper:hover .novel-cover {
  transform: translateY(-8px);
}

.cover-wrapper:hover .cover-overlay {
  opacity: 1;
}

.novel-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.novel-title {
  color: var(--text-primary);
  margin: 0;
}

.novel-author {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-secondary);
}

.novel-stats {
  display: flex;
  gap: 24px;
  margin: 8px 0;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
}

.novel-tags {
  display: flex;
  gap: 12px;
  margin: 8px 0;
}

.custom-tag {
  border-radius: 20px;
  padding: 6px 16px;
  font-weight: 500;
}

.novel-actions {
  display: flex;
  gap: 16px;
  margin-top: auto;
}

.action-button {
  padding: 12px 32px;
  font-weight: 600;
  border-radius: 12px;
}

.novel-description {
  background: var(--system-grouped-secondary);
  padding: 24px;
  border-radius: 16px;
  margin-top: 20px;
}

.novel-description h3 {
  color: var(--text-primary);
  margin: 0 0 16px;
}

.novel-description p {
  color: var(--text-secondary);
  margin: 0;
  line-height: 1.6;
}

.chapter-list {
  background: var(--system-grouped-secondary);
  padding: 32px;
  border-radius: 24px;
  border: 1px solid var(--border-color);
}

.chapter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chapter-header h2 {
  color: var(--text-primary);
  margin: 0;
}

.chapter-filters {
  display: flex;
  gap: 20px;
  align-items: center;
}

.search-input {
  width: 240px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  background: var(--system-background);
}

.chapters-list {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  overflow: hidden;
  border: 1px solid var(--border-color);
}

.chapter-item {
  display: flex;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
}

.chapter-item:last-child {
  border-bottom: none;
}

.chapter-item:hover {
  background: rgba(46, 125, 50, 0.05);
}

.chapter-info {
  flex: 1;
  margin-right: 24px;
}

.chapter-title {
  font-size: 1.1rem;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.chapter-meta {
  display: flex;
  gap: 24px;
  color: var(--text-secondary);
  font-size: 0.9rem;
  align-items: center;
}

.chapter-status {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 500;
}

.status-draft {
  background: #fff3e0;
  color: #e65100;
}

.status-published {
  background: #e8f5e9;
  color: #2e7d32;
}

.status-unknown {
  background: #f5f5f5;
  color: #757575;
}

.chapter-actions {
  display: flex;
  gap: 12px;
}

.chapter-actions .el-button {
  min-width: 80px;
}

.pagination {
  margin-top: 32px;
  text-align: center;
}

:deep(.el-pagination) {
  --el-pagination-button-bg-color: var(--system-background);
  --el-pagination-hover-color: var(--apple-blue);
}

@media (prefers-color-scheme: dark) {
  .chapter-item {
    border-color: rgba(255, 255, 255, 0.1);
  }

  .chapter-item:hover {
    border-color: var(--apple-blue);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }

  .novel-cover {
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  }

  .search-input :deep(.el-input__wrapper) {
    background: var(--system-grouped-primary);
  }
}

@media (max-width: 768px) {
  .chapter-item {
    flex-direction: column;
    align-items: flex-start;
    padding: 16px;
  }

  .chapter-info {
    margin-right: 0;
    margin-bottom: 16px;
    width: 100%;
  }

  .chapter-meta {
    flex-wrap: wrap;
    gap: 12px;
  }

  .chapter-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .chapter-filters {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .search-input {
    width: 100%;
  }
}

/* 添加收藏按钮样式 */
.novel-actions .el-button {
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.novel-actions .el-button:hover {
  transform: translateY(-2px);
}

.novel-actions .el-button:active {
  transform: translateY(0);
}

.novel-actions .el-button[type="success"] {
  background: var(--primary);
  border-color: var(--primary);
}

.novel-actions .el-button[type="success"]:hover {
  background: var(--dark);
  border-color: var(--dark);
}

/* 暗色模式适配 */
@media (prefers-color-scheme: dark) {
  .novel-actions .el-button[type="success"] {
    background: var(--primary);
    border-color: var(--primary);
  }

  .novel-actions .el-button[type="success"]:hover {
    background: var(--dark);
    border-color: var(--dark);
  }
}
</style>