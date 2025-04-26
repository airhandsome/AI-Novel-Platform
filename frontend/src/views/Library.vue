<template>
  <div class="library">
    <div class="page-header glass-effect">
      <h1 class="display-medium">书库</h1>
      <div class="filters">
        <el-select v-model="filters.category" placeholder="分类" clearable class="apple-select">
          <el-option
            v-for="item in categories"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        <el-select v-model="filters.status" placeholder="状态" clearable class="apple-select">
          <el-option label="连载中" value="ongoing" />
          <el-option label="已完结" value="completed" />
        </el-select>
        <el-select v-model="filters.sortBy" placeholder="排序" class="apple-select">
          <el-option label="最新发布" value="latest" />
          <el-option label="最多阅读" value="most_read" />
          <el-option label="最多收藏" value="most_favorite" />
        </el-select>
      </div>
    </div>

    <div class="book-list">
      <el-row :gutter="24">
        <el-col :span="6" v-for="book in books" :key="book.id">
          <div class="book-card shadow-md" @click="router.push(`/novel/${book.id}`)">
            <div class="book-cover">
              <img :src="book.cover" :alt="book.title" />
              <div class="book-status" :class="book.status">
                {{ book.status === 'ongoing' ? '连载中' : '已完结' }}
              </div>
            </div>
            <div class="book-info">
              <h3 class="title-medium">{{ book.title }}</h3>
              <p class="body-medium author">{{ book.author }}</p>
              <p class="caption description">{{ book.description }}</p>
              <div class="book-meta">
                <span>
                  <el-icon><View /></el-icon>
                  {{ book.readCount }}
                </span>
                <span>
                  <el-icon><Star /></el-icon>
                  {{ book.favoriteCount }}
                </span>
                <span>
                  <el-icon><Collection /></el-icon>
                  {{ book.wordCount }}字
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
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { View, Star, Collection } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'

const router = useRouter()

// 筛选条件
const filters = reactive({
  category: '',
  status: '',
  sortBy: 'latest'
})

// 分页
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 书籍列表
const books = ref([])

// 分类选项
const categories = [
  { label: '玄幻', value: 'fantasy' },
  { label: '科幻', value: 'sci-fi' },
  { label: '武侠', value: 'martial-arts' },
  { label: '都市', value: 'urban' },
  { label: '历史', value: 'history' },
  { label: '推理', value: 'mystery' }
]

// 加载书籍列表
const loadBooks = async () => {
  try {
    const response = await novelApi.getNovelList({
      page: currentPage.value,
      pageSize: pageSize.value,
      ...filters
    })
    books.value = response.items
    total.value = response.total
  } catch (error) {
    console.error('加载书籍列表失败:', error)
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
watch(filters, () => {
  currentPage.value = 1
  loadBooks()
})

onMounted(() => {
  loadBooks()
})
</script>

<style scoped>
.library {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  background: var(--system-background);
}

.page-header {
  position: sticky;
  top: 80px;
  z-index: 10;
  padding: 24px;
  margin: -20px -20px 24px;
  border-bottom: 1px solid var(--border-color);
}

.page-header h1 {
  color: var(--text-primary);
  margin-bottom: 20px;
}

.filters {
  display: flex;
  gap: 16px;
}

.apple-select {
  min-width: 120px;
}

.book-list {
  margin-bottom: 32px;
}

.book-card {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.book-card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.1);
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
}

.book-status {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  color: #fff;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.book-status.ongoing {
  background: rgba(0, 122, 255, 0.8);
}

.book-status.completed {
  background: rgba(52, 199, 89, 0.8);
}

.book-info {
  padding: 20px;
}

.book-info h3 {
  color: var(--text-primary);
  margin-bottom: 8px;
}

.author {
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.description {
  color: var(--text-tertiary);
  margin-bottom: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.book-meta {
  display: flex;
  gap: 20px;
  color: var(--text-tertiary);
}

.book-meta span {
  display: flex;
  align-items: center;
  gap: 6px;
}

.book-meta .el-icon {
  font-size: 16px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
  padding-bottom: 40px;
}

:deep(.el-pagination) {
  --el-pagination-button-bg-color: var(--system-grouped-secondary);
  --el-pagination-hover-color: var(--apple-blue);
}

:deep(.el-select) {
  --el-select-border-color-hover: var(--apple-blue);
  --el-select-input-focus-border-color: var(--apple-blue);
}

:deep(.el-input__wrapper) {
  background-color: var(--system-grouped-secondary);
}

@media (prefers-color-scheme: dark) {
  .book-card {
    background: var(--system-grouped-secondary);
  }
  
  .book-status {
    background: rgba(255, 255, 255, 0.2);
  }
}
</style> 