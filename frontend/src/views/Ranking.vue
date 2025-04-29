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
        <div class="rank-card" v-loading="loading.read">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><View /></el-icon>
              阅读榜
            </h3>
          </div>
          <div class="rank-list">
            <el-empty v-if="readRanking.length === 0 && !loading.read" description="暂无阅读排行数据" />
            <div
              v-else
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
                  <p class="body-medium">✍️ {{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ formatNumber(book.readCount) }}次</div>
              </div>
            </div>
          </div>
          <div class="rank-footer" v-if="readRanking.length > 0">
            <el-button type="primary" text @click="router.push('/library?sort=most_read')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </el-col>

      <!-- 收藏榜 -->
      <el-col :span="8">
        <div class="rank-card" v-loading="loading.favorite">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><Star /></el-icon>
              收藏榜
            </h3>
          </div>
          <div class="rank-list">
            <el-empty v-if="favoriteRanking.length === 0 && !loading.favorite" description="暂无收藏排行数据" />
            <div
              v-else
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
                  <p class="body-medium">✍️ {{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ formatNumber(book.favoriteCount) }}收藏</div>
              </div>
            </div>
          </div>
          <div class="rank-footer" v-if="favoriteRanking.length > 0">
            <el-button type="primary" text @click="router.push('/library?sort=most_favorite')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </el-col>

      <!-- 推荐榜 -->
      <el-col :span="8">
        <div class="rank-card" v-loading="loading.recommend">
          <div class="rank-header">
            <h3 class="title-large">
              <el-icon><Trophy /></el-icon>
              推荐榜
            </h3>
          </div>
          <div class="rank-list">
            <el-empty v-if="recommendRanking.length === 0 && !loading.recommend" description="暂无推荐排行数据" />
            <div
              v-else
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
                  <p class="body-medium">✍️ {{ book.author }}</p>
                </div>
                <div class="rank-value caption">{{ formatNumber(book.recommendCount) }}推荐</div>
              </div>
            </div>
          </div>
          <div class="rank-footer" v-if="recommendRanking.length > 0">
            <el-button type="primary" text @click="router.push('/library?sort=most_recommend')">
              查看更多 <el-icon><ArrowRight /></el-icon>
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 页面底部 -->
    <div class="ranking-footer">
      <p>排行榜数据每日更新，统计周期根据所选时间范围而定</p>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { View, Star, Trophy, ArrowRight } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'
import { ElMessage } from 'element-plus'

const router = useRouter()
const timeRange = ref('week')

// 加载状态
const loading = ref({
  read: false,
  favorite: false,
  recommend: false
})

// 排行榜数据
const readRanking = ref([])
const favoriteRanking = ref([])
const recommendRanking = ref([])

// 加载排行榜数据
const loadRankings = async () => {
  try {
    // 同时加载三种排行榜
    loading.value.read = true
    loading.value.favorite = true
    loading.value.recommend = true

    const [readData, favoriteData, recommendData] = await Promise.all([
      novelApi.getNovelRanking({ type: 'read', timeRange: timeRange.value }),
      novelApi.getNovelRanking({ type: 'favorite', timeRange: timeRange.value }),
      novelApi.getNovelRanking({ type: 'recommend', timeRange: timeRange.value })
    ])

    // 处理读取的数据
    readRanking.value = readData?.items?.map(novel => ({
      id: novel.id,
      title: novel.title,
      author: novel.author.Username,
      readCount: novel.readCount || 0
    })) || []

    favoriteRanking.value = favoriteData?.items?.map(novel => ({
      id: novel.id,
      title: novel.title,
      author: novel.author.Username,
      favoriteCount: novel.favoriteCount || 0
    })) || []

    recommendRanking.value = recommendData?.items?.map(novel => ({
      id: novel.id,
      title: novel.title,
      author: novel.author.Username,
      recommendCount: novel.recommendCount || 0
    })) || []
  } catch (error) {
    console.error('加载排行榜数据失败:', error)
    ElMessage.error('加载排行榜数据失败，请稍后重试')
  } finally {
    loading.value.read = false
    loading.value.favorite = false
    loading.value.recommend = false
  }
}

// 格式化数字（简化大数显示）
const formatNumber = (num) => {
  if (!num) return 0
  
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  
  return num
}

// 监听时间范围变化，重新加载数据
watch(timeRange, () => {
  loadRankings()
})

onMounted(() => {
  loadRankings()
})
</script>

<style scoped>
/* 使用Library.vue中相同的根变量 */
:root {
  --primary: #2e7d32;
  --secondary: #81c784;
  --light: #f1f8e9;
  --dark: #1b5e20;
  --border: #c8e6c9;
  --code-bg: #e8f5e9;
}

.ranking {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  background-color: #f9f9f9;
}

.page-header {
  position: sticky;
  top: 80px;
  z-index: 10;
  padding: 24px 30px;
  margin: -20px -20px 24px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
  backdrop-filter: blur(10px);
  background-color: rgba(255, 255, 255, 0.9);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  border-radius: 0 0 16px 16px;
}

.page-header h1 {
  color: var(--dark);
  margin: 0;
  font-weight: 600;
  font-size: 1.8rem;
}

:deep(.el-radio-button__inner) {
  border-radius: 20px;
  padding: 8px 20px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.el-radio-button__original) {
  opacity: 0;
}

:deep(.el-radio-button:first-child .el-radio-button__inner) {
  border-radius: 20px 0 0 20px;
}

:deep(.el-radio-button:last-child .el-radio-button__inner) {
  border-radius: 0 20px 20px 0;
}

:deep(.el-radio-button.is-active .el-radio-button__inner) {
  background-color: var(--primary);
  border-color: var(--primary);
  box-shadow: 0 0 10px rgba(46, 125, 50, 0.3);
}

.rank-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  margin-bottom: 24px;
  overflow: hidden;
  border: 1px solid rgba(46, 125, 50, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.rank-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--secondary));
  opacity: 0.8;
  z-index: 1;
}

.rank-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 15px 35px rgba(46, 125, 50, 0.15);
}

.rank-header {
  padding: 20px 24px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  background: var(--light);
  position: relative;
}

.rank-header h3 {
  margin: 0;
  color: var(--dark);
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.rank-header h3 .el-icon {
  color: var(--primary);
}

.rank-list {
  padding: 8px 0;
  flex-grow: 1;
  min-height: 300px;
}

.rank-item {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  margin: 4px 8px;
  border-radius: 12px;
  position: relative;
  overflow: hidden;
}

.rank-item:hover {
  background-color: #f8fffa;
  transform: scale(1.02) translateX(4px);
}

.rank-number {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background-color: var(--border);
  color: var(--dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  margin-right: 16px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
  z-index: 1;
  transition: all 0.3s ease;
}

.rank-item:hover .rank-number {
  transform: scale(1.1);
}

.rank-number.top-1 {
  background: linear-gradient(135deg, #f44336, #ff9800);
  color: white;
  box-shadow: 0 3px 10px rgba(244, 67, 54, 0.3);
}

.rank-number.top-2 {
  background: linear-gradient(135deg, #9c27b0, #673ab7);
  color: white;
  box-shadow: 0 3px 10px rgba(156, 39, 176, 0.3);
}

.rank-number.top-3 {
  background: linear-gradient(135deg, #2196f3, #00bcd4);
  color: white;
  box-shadow: 0 3px 10px rgba(33, 150, 243, 0.3);
}

.rank-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 1;
}

.book-info {
  flex: 1;
  padding-right: 16px;
}

.book-info h4 {
  margin: 0;
  color: #333;
  margin-bottom: 4px;
  font-weight: 500;
  transition: color 0.3s;
}

.rank-item:hover .book-info h4 {
  color: var(--primary);
}

.book-info p {
  margin: 0;
  color: #666;
  font-size: 0.9rem;
}

.rank-value {
  color: #666;
  min-width: 80px;
  text-align: right;
  background: var(--light);
  padding: 4px 10px;
  border-radius: 12px;
  font-weight: 500;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  font-size: 0.85rem;
}

.rank-item:hover .rank-value {
  background: var(--secondary);
  color: var(--dark);
  box-shadow: 0 2px 8px rgba(46, 125, 50, 0.15);
}

.rank-footer {
  padding: 16px;
  text-align: center;
  border-top: 1px solid rgba(46, 125, 50, 0.1);
}

.rank-footer .el-button {
  color: var(--primary);
  transition: all 0.3s ease;
}

.rank-footer .el-button:hover {
  color: var(--dark);
  transform: translateX(4px);
}

.ranking-footer {
  text-align: center;
  color: #666;
  font-size: 0.9rem;
  margin-top: 20px;
  padding: 20px 0;
  border-top: 1px solid rgba(46, 125, 50, 0.1);
}

/* 移动端适配 */
@media (max-width: 900px) {
  .page-header {
    flex-direction: column;
    gap: 16px;
    padding: 20px;
  }
  
  .rank-card {
    margin-bottom: 16px;
  }
}

@media (max-width: 600px) {
  .page-header {
    padding: 16px;
  }
  
  .time-filter {
    width: 100%;
    overflow-x: auto;
  }
  
  :deep(.el-radio-button__inner) {
    padding: 8px 12px;
  }
  
  .rank-item {
    padding: 12px 16px;
    margin: 2px 4px;
  }
}

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .ranking {
    background-color: var(--system-background);
  }

  .page-header {
    background-color: rgba(30, 30, 30, 0.85);
    border-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  }

  .page-header h1 {
    color: #eee;
  }

  .rank-card {
    background: var(--system-grouped-secondary);
    border-color: rgba(255, 255, 255, 0.08);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
  }

  .rank-header {
    background: rgba(46, 125, 50, 0.15);
    border-color: rgba(255, 255, 255, 0.08);
  }
  
  .rank-header h3 {
    color: #eee;
  }

  .rank-item:hover {
    background-color: rgba(46, 125, 50, 0.1);
  }

  .rank-number {
    background-color: var(--dark);
    color: #eee;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
  }
  
  .book-info h4 {
    color: #eee;
  }
  
  .book-info p {
    color: #aaa;
  }
  
  .rank-value {
    background: rgba(46, 125, 50, 0.15);
    color: #aaa;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  }

  .rank-item:hover .rank-value {
    background: var(--secondary);
    color: var(--dark);
  }
  
  .rank-footer {
    border-color: rgba(255, 255, 255, 0.08);
  }
  
  .rank-footer .el-button {
    color: var(--secondary);
  }
  
  .ranking-footer {
    color: #aaa;
    border-color: rgba(255, 255, 255, 0.08);
  }
}
</style>