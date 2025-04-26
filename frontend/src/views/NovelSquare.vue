<template>
  <div class="novel-square">
    <div class="header glass-effect">
      <h1 class="display-medium">小说广场</h1>
      <div class="filters">
        <el-select 
          v-model="filter.category" 
          placeholder="分类" 
          clearable
          class="apple-select"
        >
          <el-option
            v-for="item in categories"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        <el-select 
          v-model="filter.status" 
          placeholder="状态" 
          clearable
          class="apple-select"
        >
          <el-option
            v-for="item in statusOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        <el-select 
          v-model="filter.sort" 
          placeholder="排序"
          class="apple-select"
        >
          <el-option
            v-for="item in sortOptions"
            :key="item.value"
            :label="item.label"
            :value="item.value"
          />
        </el-select>
        <div class="search-wrapper">
          <el-input
            v-model="filter.keyword"
            placeholder="搜索小说或作者"
            clearable
            class="search-input"
            @keyup.enter="searchNovels"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <div class="content">
      <!-- 小说列表 -->
      <el-row :gutter="24">
        <el-col :xs="12" :sm="8" :md="6" :lg="4" v-for="novel in novels" :key="novel.id">
          <div class="novel-card" @click="goToNovel(novel.id)">
            <div class="cover-wrapper">
              <img :src="novel.coverUrl || '/default-cover.jpg'" class="novel-cover" />
              <div class="cover-overlay"></div>
              <div class="novel-status" :class="{ 'completed': novel.status === 2 }">
                {{ novel.status === 1 ? '连载中' : novel.status === 2 ? '已完结' : '草稿' }}
              </div>
            </div>
            <div class="novel-info">
              <h3 class="title-medium novel-title">{{ novel.title }}</h3>
              <p class="body-medium novel-author">{{ novel.author }}</p>
              <p class="caption novel-description">{{ novel.description }}</p>
              <div class="novel-stats">
                <span class="stat-item">
                  <el-icon><View /></el-icon> {{ formatNumber(novel.readCount) }}
                </span>
                <span class="stat-item">
                  <el-icon><Star /></el-icon> {{ formatNumber(novel.favoriteCount) }}
                </span>
                <span class="stat-item">
                  <el-icon><Document /></el-icon> {{ formatNumber(novel.wordCount) }}字
                </span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 分页器 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[12, 24, 36, 48]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="totalNovels"
          background
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, View, Star, Document } from '@element-plus/icons-vue'
// 假设我们已经创建了这个API
import { novelApi } from '../api/novel'

const router = useRouter()

// 分页数据
const currentPage = ref(1)
const pageSize = ref(12)
const totalNovels = ref(0)

// 过滤器
const filter = reactive({
  category: '',
  status: '',
  sort: 'latest',
  keyword: ''
})

// 分类选项
const categories = ref([
  { value: 'fantasy', label: '奇幻' },
  { value: 'scifi', label: '科幻' },
  { value: 'romance', label: '言情' },
  { value: 'mystery', label: '悬疑' },
  { value: 'history', label: '历史' },
  { value: 'other', label: '其他' }
])

// 状态选项
const statusOptions = ref([
  { value: 1, label: '连载中' },
  { value: 2, label: '已完结' }
])

// 排序选项
const sortOptions = ref([
  { value: 'latest', label: '最新发布' },
  { value: 'popular', label: '人气最高' },
  { value: 'favorites', label: '收藏最多' },
  { value: 'updateTime', label: '最近更新' }
])

// 小说数据
const novels = ref([])

// 获取小说列表
const fetchNovels = async () => {
  try {
    const { data } = await novelApi.getNovels({
      page: currentPage.value,
      limit: pageSize.value,
      category: filter.category,
      status: filter.status, 
      sort: filter.sort,
      keyword: filter.keyword
    })
    
    novels.value = data.novels
    totalNovels.value = data.total
  } catch (error) {
    ElMessage.error('获取小说列表失败')
    
    // 模拟数据，实际应从API获取
    novels.value = Array(12).fill(null).map((_, i) => ({
      id: i + 1,
      title: `示例小说 ${i + 1}`,
      author: `作者 ${i % 5 + 1}`,
      description: '这是一个小说描述，简短介绍小说的内容和特点。',
      coverUrl: `/sample-covers/cover${i % 8 + 1}.jpg`,
      status: i % 3 + 1, // 0: 草稿, 1: 连载中, 2: 已完结
      readCount: Math.floor(Math.random() * 10000),
      favoriteCount: Math.floor(Math.random() * 1000),
      wordCount: Math.floor(Math.random() * 500000)
    }))
    totalNovels.value = 100
  }
}

// 搜索小说
const searchNovels = () => {
  currentPage.value = 1
  fetchNovels()
}

// 处理页码变化
const handleCurrentChange = (page) => {
  currentPage.value = page
  fetchNovels()
}

// 处理每页数量变化
const handleSizeChange = (size) => {
  pageSize.value = size
  fetchNovels()
}

// 跳转到小说详情
const goToNovel = (id) => {
  router.push(`/novel/${id}`)
}

// 格式化数字，例如 1000 -> 1k
const formatNumber = (num) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num
}

// 监听过滤条件变化
watch(
  () => [filter.category, filter.status, filter.sort],
  () => {
    currentPage.value = 1
    fetchNovels()
  }
)

onMounted(() => {
  fetchNovels()
})
</script>

<style scoped>
.novel-square {
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
  position: sticky;
  top: 20px;
  z-index: 10;
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.header h1 {
  margin: 0 0 24px;
  color: var(--text-primary);
}

.filters {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.apple-select {
  width: 160px;
}

.search-wrapper {
  flex: 1;
  min-width: 280px;
}

.search-input {
  width: 100%;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 12px;
  background: var(--system-background);
}

.content {
  margin-top: 32px;
}

.novel-card {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
  height: 100%;
  display: flex;
  flex-direction: column;
}

.novel-card:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.05),
    0 12px 25px rgba(0, 0, 0, 0.07);
  border-color: var(--apple-blue);
}

.cover-wrapper {
  position: relative;
  height: 240px;
  overflow: hidden;
}

.novel-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.5s ease;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    transparent 50%,
    rgba(0, 0, 0, 0.4)
  );
  opacity: 0;
  transition: opacity 0.3s ease;
}

.novel-card:hover .novel-cover {
  transform: scale(1.05);
}

.novel-card:hover .cover-overlay {
  opacity: 1;
}

.novel-status {
  position: absolute;
  top: 12px;
  right: 12px;
  background: rgba(52, 199, 89, 0.9);
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  transition: transform 0.3s ease;
}

.novel-status.completed {
  background: rgba(142, 142, 147, 0.9);
}

.novel-card:hover .novel-status {
  transform: translateY(-2px);
}

.novel-info {
  padding: 20px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.novel-title {
  color: var(--text-primary);
  margin: 0 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
}

.novel-author {
  color: var(--text-secondary);
  margin: 0 0 12px;
}

.novel-description {
  color: var(--text-tertiary);
  margin: 0 0 16px;
  flex-grow: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
}

.novel-stats {
  display: flex;
  justify-content: space-between;
  color: var(--text-tertiary);
  border-top: 1px solid var(--border-color);
  padding-top: 12px;
  margin-top: auto;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.stat-item .el-icon {
  font-size: 16px;
}

.pagination {
  margin-top: 48px;
  margin-bottom: 48px;
  text-align: center;
}

:deep(.el-pagination) {
  --el-pagination-button-bg-color: var(--system-background);
  --el-pagination-hover-color: var(--apple-blue);
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

  .novel-status {
    background: rgba(52, 199, 89, 0.7);
  }

  .novel-status.completed {
    background: rgba(142, 142, 147, 0.7);
  }
}
</style> 