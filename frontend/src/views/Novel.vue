<template>
  <div class="novel-detail">
    <!-- 小说基本信息 -->
    <div class="novel-header glass-effect">
      <div class="novel-info">
        <div class="cover-wrapper">
          <img :src="novel.coverUrl || '/default-cover.jpg'" class="novel-cover" />
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
              :type="novel.isFavorited ? 'success' : ''" 
              size="large"
              class="action-button"
              @click="toggleFavorite"
            >
              <el-icon><Star /></el-icon>
              {{ novel.isFavorited ? '已收藏' : '收藏' }}
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

      <el-scrollbar height="500px" class="chapter-scrollbar">
        <div class="chapter-grid">
          <div
            v-for="chapter in displayChapters"
            :key="chapter.id"
            class="chapter-item"
            @click="goToChapter(chapter.id)"
          >
            <span class="chapter-title title-medium" :class="{ 'read': chapter.isRead }">
              {{ chapter.title }}
            </span>
            <span class="chapter-word-count caption">{{ formatNumber(chapter.wordCount) }}字</span>
          </div>
        </div>
      </el-scrollbar>

      <!-- 分页器 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="totalChapters"
          layout="total, sizes, prev, pager, next, jumper"
          background
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

const route = useRoute()
const router = useRouter()

// 小说信息
const novel = ref({})
const chapters = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const totalChapters = ref(0)
const searchKeyword = ref('')
const reverseOrder = ref(false)

// 获取小说信息
const fetchNovelDetail = async () => {
  try {
    const { data } = await novelApi.getNovelDetail(route.params.id)
    novel.value = data
  } catch (error) {
    ElMessage.error('获取小说信息失败')
  }
}

// 获取章节列表
const fetchChapters = async () => {
  try {
    const { data } = await novelApi.getNovelChapters({
      novelId: route.params.id,
      page: currentPage.value,
      limit: pageSize.value,
      keyword: searchKeyword.value,
      reverse: reverseOrder.value
    })
    chapters.value = data.chapters
    totalChapters.value = data.total
  } catch (error) {
    ElMessage.error('获取章节列表失败')
  }
}

// 过滤显示的章节
const displayChapters = computed(() => {
  return chapters.value
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
  router.push(`/novel/${route.params.id}/chapter/${chapterId}`)
}

// 收藏/取消收藏
const toggleFavorite = async () => {
  try {
    if (novel.value.isFavorited) {
      await novelApi.unfavoriteNovel(route.params.id)
      ElMessage.success('已取消收藏')
    } else {
      await novelApi.favoriteNovel(route.params.id)
      ElMessage.success('收藏成功')
    }
    novel.value.isFavorited = !novel.value.isFavorited
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

.chapter-scrollbar :deep(.el-scrollbar__view) {
  padding: 4px;
}

.chapter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.chapter-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--system-background);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
}

.chapter-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: var(--apple-blue);
}

.chapter-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--text-primary);
}

.chapter-title.read {
  color: var(--text-tertiary);
}

.chapter-word-count {
  color: var(--text-tertiary);
  margin-left: 12px;
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
}
</style>