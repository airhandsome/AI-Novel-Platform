<template>
  <div class="library">
    <div class="header">
      <div class="header-content">
        <h1>📚 探索书库</h1>
        <p>发现属于你的精彩故事，开启文字的奇妙旅程</p>
        <div class="tagline">让我们一起探索文学的无限可能</div>
      </div>
    </div>

    <div class="feature-list">
      <div class="feature-card">
        <div class="feature-icon">📖</div>
        <h3 class="feature-title">海量小说</h3>
        <p>各类题材，应有尽有</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">✨</div>
        <h3 class="feature-title">优质创作</h3>
        <p>精品内容，创意无限</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🎯</div>
        <h3 class="feature-title">个性推荐</h3>
        <p>智能算法，量身定制</p>
      </div>
      <div class="feature-card">
        <div class="feature-icon">🚀</div>
        <h3 class="feature-title">即时更新</h3>
        <p>实时连载，追更无忧</p>
      </div>
    </div>

    <div class="main-content">
      <div class="filters-container">
        <div class="search-group">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索你感兴趣的小说或作者 ✨"
            class="search-input"
            clearable
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button type="primary" :icon="Search" @click="handleSearch" :loading="isLoading" class="search-button">
            搜索
          </el-button>
        </div>
        <div class="filter-group">
          <el-select v-model="filters.category" placeholder="选择分类 🎯" clearable class="apple-select">
            <el-option
              v-for="item in categories"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
          <el-select v-model="filters.status" placeholder="作品状态 📖" clearable class="apple-select">
            <el-option label="全部" value="-1" />
            <el-option label="连载中 ⚡️" value="0" />
            <el-option label="已完结 🎉" value="1" />
            <el-option label="烂尾 💔" value="2" />
          </el-select>
          <el-select v-model="filters.sortBy" placeholder="排序方式 🔍" class="apple-select">
            <el-option label="最新发布 🆕" value="latest" />
            <el-option label="最多阅读 👀" value="most_read" />
            <el-option label="最多收藏 ⭐️" value="most_favorite" />
          </el-select>
        </div>
      </div>

      <div v-if="books.length === 0" class="empty-state">
        <el-empty description="暂时没有找到相关小说 🔍">
          <template #image>
            <img src="../../public/assets/empty-books.png" class="empty-image" alt="没有找到书籍" />
          </template>
          <el-button type="primary" @click="resetFilters">重置筛选 ↺</el-button>
        </el-empty>
      </div>

      <div v-else class="book-list">
        <el-row :gutter="24">
          <el-col :lg="6" :md="8" :sm="12" :xs="24" v-for="book in books" :key="book.id">
            <div class="book-card shadow-md" @click="router.push(`/novels/${book.id}`)">
              <div class="book-cover">
                <img :src="book.cover" :alt="book.title" />
                <div class="book-status-bar">
                  <div class="book-category">{{ getCategoryText(book.category) }}</div>
                  <div class="book-status" :class="getStatusStyle(book.status)">
                    {{ getStatusText(book.status) }}
                  </div>
                </div>
              </div>
              <div class="book-info">
                <h3 class="title-medium">{{ book.title }}</h3>
                <p class="body-medium author">✍️ {{ book.author }}</p>
                <p class="caption description">{{ book.description }}</p>
                <div class="book-meta">
                  <span>
                    <el-icon><View /></el-icon>
                    {{ formatNumber(book.readCount) }}
                  </span>
                  <span>
                    <el-icon><Star /></el-icon>
                    {{ formatNumber(book.favoriteCount) }}
                  </span>
                  <span>
                    <el-icon><Collection /></el-icon>
                    {{ formatNumber(book.wordCount) }}字
                  </span>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { View, Star, Collection, Search } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'
import { ElMessage } from 'element-plus'

const router = useRouter()

// 搜索相关状态
const searchKeyword = ref('')
const isLoading = ref(false)

// 筛选条件
const filters = reactive({
  category: 'all',  // 默认全部
  status: '-1',  // 默认全部
  sortBy: 'latest',
  keyword: ''  // 搜索关键词
})

// 分页
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 书籍列表
const books = ref([])

// 分类选项
const categories = [
  { label: '全部', value: 'all' },
  { label: '玄幻', value: 'fantasy' },
  { label: '科幻', value: 'sci-fi' },
  { label: '武侠', value: 'martial-arts' },
  { label: '都市', value: 'urban' },
  { label: '历史', value: 'history' },
  { label: '推理', value: 'mystery' }
]

// 获取分类文本
const getCategoryText = (categoryValue) => {
  const category = categories.find(c => c.value === categoryValue)
  return category ? category.label : '未知'
}

// 处理搜索
const handleSearch = async () => {
  filters.keyword = searchKeyword.value
  currentPage.value = 1
  await loadBooks()
}

// 加载书籍列表
const loadBooks = async () => {
  try {
    isLoading.value = true
    const response = await novelApi.getNovelList({
      page: currentPage.value,
      pageSize: pageSize.value,
      ...filters
    })
    
    books.value = response.novels.map(novel => ({
      id: novel.id,
      title: novel.title,
      description: novel.description,
      author: novel.author.Username,
      cover: novel.coverUrl || '/default-cover.jpg',
      status: novel.status,
      readCount: novel.readCount,
      favoriteCount: novel.favoriteCount,
      wordCount: novel.wordCount,
      category: novel.category
    }))
    
    total.value = response.total
  } catch (error) {
    console.error('加载书籍列表失败:', error)
    ElMessage.error('加载书籍列表失败')
  } finally {
    isLoading.value = false
  }
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  loadBooks()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadBooks()
}

// 监听筛选条件变化
watch(
  () => [filters.category, filters.status, filters.sortBy],
  () => {
    currentPage.value = 1
    loadBooks()
  }
)

// 修改状态显示逻辑
const getStatusText = (status) => {
  const statusMap = {
    0: '连载中',
    1: '已完结',
    2: '烂尾',
    '-1': '全部'
  }
  return statusMap[status] || '未知'
}

// 获取状态样式
const getStatusStyle = (status) => {
  const styleMap = {
    0: 'status-ongoing',
    1: 'status-completed',
    2: 'status-abandoned'
  }
  return styleMap[status] || ''
}

// 添加重置筛选方法
const resetFilters = () => {
  searchKeyword.value = ''
  filters.category = 'all'
  filters.status = '-1'
  filters.sortBy = 'latest'
  filters.keyword = ''
  loadBooks()
}

// 格式化数字
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

onMounted(() => {
  loadBooks()
})
</script>

<style scoped>
:root {
  --primary: #2e7d32;
  --secondary: #81c784;
  --light: #f1f8e9;
  --dark: #1b5e20;
  --border: #c8e6c9;
  --code-bg: #e8f5e9;
}

.library {
  min-height: 100vh;
  background-color: #f9f9f9;
}

.header {
  background: linear-gradient(135deg, #2E7D32 0%, #1B5E20 100%);
  color: white;
  padding: 3rem 0;
  text-align: center;
  position: relative;
  margin-bottom: 2rem;
}

.header::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 20% 30%, rgba(255,255,255,0.2) 0%, transparent 20%),
    radial-gradient(circle at 80% 70%, rgba(255,255,255,0.2) 0%, transparent 20%);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  position: relative;
  z-index: 1;
}

.header h1 {
  font-size: 2.8rem;
  margin: 0;
  text-shadow: 1px 1px 3px rgba(0,0,0,0.2);
}

.header p {
  font-size: 1.2rem;
  max-width: 800px;
  margin: 1rem auto 0;
  opacity: 0.9;
}

.tagline {
  display: inline-block;
  background: rgba(255,255,255,0.2);
  padding: 0.5rem 1.5rem;
  border-radius: 20px;
  margin-top: 1.5rem;
  font-size: 0.9rem;
  backdrop-filter: blur(5px);
}

.feature-list {
  display: flex;
  justify-content: space-between;
  gap: 1.5rem;
  margin: 3rem auto;
  max-width: 1200px;
  padding: 0 2rem;
}

.feature-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  flex: 1;
  box-shadow: 0 3px 10px rgba(46,125,50,0.1);
  border-top: 3px solid var(--primary);
  transition: all 0.3s ease;
  text-align: center;
  min-width: 0;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(46,125,50,0.15);
  background-color: #f8fffa;
}

.feature-icon {
  font-size: 2.5rem;
  margin-bottom: 1.5rem;
  color: var(--primary);
}

.feature-title {
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: var(--dark);
  font-size: 1.1rem;
}

.feature-card p {
  color: #555;
  font-size: 0.95rem;
  line-height: 1.6;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.filters-container {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  margin-bottom: 2rem;
  border: 1px solid rgba(46, 125, 50, 0.1);
  position: relative;
  overflow: hidden;
}

.filters-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--secondary));
}

.search-group {
  display: flex;
  gap: 16px;
  margin-bottom: 1.5rem;
  position: relative;
}

.search-input {
  flex: 1;
}

.search-input :deep(.el-input__wrapper) {
  background-color: #f8f9fa;
  border: 2px solid transparent;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  padding: 0 16px;
}

.search-input :deep(.el-input__wrapper:hover) {
  background-color: #fff;
  border-color: rgba(46, 125, 50, 0.1);
}

.search-input :deep(.el-input__wrapper.is-focus) {
  background-color: #fff;
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(46, 125, 50, 0.1);
}

.search-button {
  min-width: 120px;
  height: 40px;
  border-radius: 12px;
  background: var(--primary);
  border-color: var(--primary);
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.search-button:hover {
  background: var(--dark);
  border-color: var(--dark);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(46, 125, 50, 0.2);
}

.search-button:active {
  transform: translateY(0);
}

.filter-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 12px;
  border: 1px solid rgba(46, 125, 50, 0.1);
}

.apple-select {
  min-width: 180px;
  flex: 1;
}

.apple-select :deep(.el-input__wrapper) {
  background-color: white;
  border: 2px solid transparent;
  border-radius: 10px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.apple-select :deep(.el-input__wrapper:hover) {
  border-color: rgba(46, 125, 50, 0.1);
}

.apple-select :deep(.el-input__wrapper.is-focus) {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px rgba(46, 125, 50, 0.1);
}

.apple-select :deep(.el-select-dropdown__item) {
  padding: 8px 16px;
}

.apple-select :deep(.el-select-dropdown__item.selected) {
  background-color: rgba(46, 125, 50, 0.1);
  color: var(--primary);
  font-weight: 500;
}

.apple-select :deep(.el-select-dropdown__item:hover) {
  background-color: rgba(46, 125, 50, 0.05);
}

.book-list {
  margin-top: 2rem;
}

.book-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  border: 1px solid #e0e0e0;
}

.book-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(46,125,50,0.15);
}

.book-cover {
  position: relative;
  height: 240px;
  overflow: hidden;
}

.book-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.book-card:hover .book-cover img {
  transform: scale(1.08);
}

.book-info {
  padding: 20px;
}

.book-info h3 {
  font-size: 1.25rem;
  margin: 0 0 8px;
  color: var(--text-primary);
}

.author {
  font-size: 0.95rem;
  color: var(--text-secondary);
  margin: 0 0 12px;
}

.description {
  font-size: 0.9rem;
  color: var(--text-tertiary);
  margin: 0 0 16px;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.book-meta {
  display: flex;
  gap: 16px;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

.book-meta span {
  display: flex;
  align-items: center;
  gap: 6px;
}

.empty-state {
  margin: 80px 0;
  text-align: center;
}

.empty-image {
  width: 240px;
  margin-bottom: 24px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin: 40px 0;
  padding-bottom: 40px;
}

/* 移动端适配 */
@media (max-width: 900px) {
  .filters-container {
    padding: 1.5rem;
  }

  .filter-group {
    gap: 12px;
  }
  
  .apple-select {
    min-width: calc(50% - 8px);
  }
}

@media (max-width: 600px) {
  .filters-container {
    padding: 1rem;
    border-radius: 12px;
  }

  .search-group {
    flex-direction: column;
    gap: 12px;
  }

  .search-button {
    width: 100%;
    height: 44px;
  }

  .filter-group {
    padding: 0.75rem;
    gap: 8px;
  }

  .apple-select {
    width: 100%;
    min-width: 100%;
  }
}

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .book-card {
    background: var(--system-grouped-secondary);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .book-card:hover {
    border-color: var(--apple-blue);
    box-shadow: 
      0 8px 16px rgba(0, 0, 0, 0.2),
      0 16px 32px rgba(0, 0, 0, 0.3);
  }

  .header-content h1 {
    background: linear-gradient(45deg, #0A84FF, #5E5CE6);
    -webkit-background-clip: text;
  }

  .search-group {
    background: var(--system-grouped-primary);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .search-input :deep(.el-input__wrapper) {
    background-color: var(--system-background);
  }

  .apple-select :deep(.el-input__wrapper) {
    background-color: var(--system-grouped-primary);
  }
}
</style> 