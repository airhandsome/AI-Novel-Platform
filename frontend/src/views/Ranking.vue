<template>
  <div class="ranking">
    <div class="page-header glass-effect">
      <h1 class="display-medium">排行榜</h1>
      <div class="time-filter">
        <el-radio-group v-model="timeRange" size="large">
          <el-radio-button value="day">日榜</el-radio-button>
          <el-radio-button value="week">周榜</el-radio-button>
          <el-radio-button value="month">月榜</el-radio-button>
          <el-radio-button value="all">总榜</el-radio-button>
        </el-radio-group>
      </div>
    </div>

    <el-row :gutter="24">
      <!-- 阅读榜 -->
      <el-col :span="8">
        <div class="rank-card">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><View /></el-icon>
              阅读榜
            </h3>
          </div>
          <div class="rank-list">
            <div
              v-for="(book, index) in readRanking"
              :key="book.id"
              class="rank-item"
              @click="router.push(`/novel/${book.id}`)"
            >
              <div class="rank-number" :class="{ 'top-1': index === 0, 'top-2': index === 1, 'top-3': index === 2 }">
                {{ index + 1 }}
              </div>
              <div class="rank-content">
                <div class="book-info">
                  <h4 class="title-medium">{{ book.title }}</h4>
                  <p class="body-medium">{{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ book.readCount }}次</div>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <!-- 收藏榜 -->
      <el-col :span="8">
        <div class="rank-card">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><Star /></el-icon>
              收藏榜
            </h3>
          </div>
          <div class="rank-list">
            <div
              v-for="(book, index) in favoriteRanking"
              :key="book.id"
              class="rank-item"
              @click="router.push(`/novel/${book.id}`)"
            >
              <div class="rank-number" :class="{ 'top-1': index === 0, 'top-2': index === 1, 'top-3': index === 2 }">
                {{ index + 1 }}
              </div>
              <div class="rank-content">
                <div class="book-info">
                  <h4 class="title-medium">{{ book.title }}</h4>
                  <p class="body-medium">{{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ book.favoriteCount }}收藏</div>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <!-- 推荐榜 -->
      <el-col :span="8">
        <div class="rank-card">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><Trophy /></el-icon>
              推荐榜
            </h3>
          </div>
          <div class="rank-list">
            <div
              v-for="(book, index) in recommendRanking"
              :key="book.id"
              class="rank-item"
              @click="router.push(`/novel/${book.id}`)"
            >
              <div class="rank-number" :class="{ 'top-1': index === 0, 'top-2': index === 1, 'top-3': index === 2 }">
                {{ index + 1 }}
              </div>
              <div class="rank-content">
                <div class="book-info">
                  <h4 class="title-medium">{{ book.title }}</h4>
                  <p class="body-medium">{{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ book.recommendCount }}推荐</div>
              </div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { View, Star, Trophy } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'

const router = useRouter()
const timeRange = ref('week')

// 排行榜数据
const readRanking = ref([])
const favoriteRanking = ref([])
const recommendRanking = ref([])

// 加载排行榜数据
const loadRankings = async () => {
  try {
    const [readData, favoriteData, recommendData] = await Promise.all([
      novelApi.getNovelRanking({ type: 'read', timeRange: timeRange.value }),
      novelApi.getNovelRanking({ type: 'favorite', timeRange: timeRange.value }),
      novelApi.getNovelRanking({ type: 'recommend', timeRange: timeRange.value })
    ])

    readRanking.value = readData.items
    favoriteRanking.value = favoriteData.items
    recommendRanking.value = recommendData.items
  } catch (error) {
    console.error('加载排行榜数据失败:', error)
  }
}

// 监听时间范围变化
watch(timeRange, () => {
  loadRankings()
})

onMounted(() => {
  loadRankings()
})
</script>

<style scoped>
.ranking {
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  color: var(--text-primary);
  margin: 0;
}

:deep(.el-radio-button__inner) {
  border-radius: 20px;
  padding: 8px 20px;
  font-weight: 500;
}

:deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-radius: 20px 0 0 20px;
}

:deep(.el-radio-button:last-child .el-radio-button__inner) {
  border-radius: 0 20px 20px 0;
}

.rank-card {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
  overflow: hidden;
}

.rank-header {
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.rank-header h3 {
  margin: 0;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.rank-list {
  padding: 8px 0;
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.rank-item:hover {
  background-color: var(--system-background);
}

.rank-number {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background-color: var(--apple-gray-4);
  color: var(--text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  margin-right: 16px;
}

.rank-number.top-1 {
  background: linear-gradient(135deg, var(--apple-pink), var(--apple-orange));
  color: white;
}

.rank-number.top-2 {
  background: linear-gradient(135deg, var(--apple-purple), var(--apple-indigo));
  color: white;
}

.rank-number.top-3 {
  background: linear-gradient(135deg, var(--apple-blue), var(--apple-teal));
  color: white;
}

.rank-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.book-info {
  flex: 1;
  padding-right: 16px;
}

.book-info h4 {
  margin: 0;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.book-info p {
  margin: 0;
  color: var(--text-secondary);
}

.rank-value {
  color: var(--text-tertiary);
  min-width: 80px;
  text-align: right;
}

@media (prefers-color-scheme: dark) {
  .rank-card {
    background: var(--system-grouped-secondary);
  }

  .rank-item:hover {
    background-color: var(--system-background);
  }

  .rank-number {
    background-color: var(--apple-gray-1);
  }
}
</style> 