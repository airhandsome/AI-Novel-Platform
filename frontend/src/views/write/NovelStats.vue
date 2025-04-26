<template>
  <div class="novel-stats">
    <div class="header">
      <h1>数据统计</h1>
      <el-button type="primary" @click="refreshStats">
        <el-icon><Refresh /></el-icon>刷新数据
      </el-button>
    </div>

    <div class="stats-content">
      <!-- 基础数据卡片 -->
      <div class="stats-cards">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card class="stats-card">
              <template #header>
                <div class="card-header">
                  <span>总字数</span>
                  <el-tooltip content="作品所有章节的总字数">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <div class="card-value">{{ formatNumber(stats.totalWords) }}</div>
              <div class="card-trend" :class="{ 'up': stats.wordsTrend > 0 }">
                {{ stats.wordsTrend > 0 ? '+' : '' }}{{ stats.wordsTrend }}%
                <el-icon><ArrowUp v-if="stats.wordsTrend > 0" /><ArrowDown v-else /></el-icon>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <template #header>
                <div class="card-header">
                  <span>总阅读量</span>
                  <el-tooltip content="作品的总阅读次数">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <div class="card-value">{{ formatNumber(stats.totalReads) }}</div>
              <div class="card-trend" :class="{ 'up': stats.readsTrend > 0 }">
                {{ stats.readsTrend > 0 ? '+' : '' }}{{ stats.readsTrend }}%
                <el-icon><ArrowUp v-if="stats.readsTrend > 0" /><ArrowDown v-else /></el-icon>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <template #header>
                <div class="card-header">
                  <span>总收藏数</span>
                  <el-tooltip content="作品被收藏的总次数">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <div class="card-value">{{ formatNumber(stats.totalFavorites) }}</div>
              <div class="card-trend" :class="{ 'up': stats.favoritesTrend > 0 }">
                {{ stats.favoritesTrend > 0 ? '+' : '' }}{{ stats.favoritesTrend }}%
                <el-icon><ArrowUp v-if="stats.favoritesTrend > 0" /><ArrowDown v-else /></el-icon>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="stats-card">
              <template #header>
                <div class="card-header">
                  <span>评论数</span>
                  <el-tooltip content="作品的总评论数">
                    <el-icon><InfoFilled /></el-icon>
                  </el-tooltip>
                </div>
              </template>
              <div class="card-value">{{ formatNumber(stats.totalComments) }}</div>
              <div class="card-trend" :class="{ 'up': stats.commentsTrend > 0 }">
                {{ stats.commentsTrend > 0 ? '+' : '' }}{{ stats.commentsTrend }}%
                <el-icon><ArrowUp v-if="stats.commentsTrend > 0" /><ArrowDown v-else /></el-icon>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 趋势图表 -->
      <div class="stats-charts">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-card class="chart-card">
              <template #header>
                <div class="card-header">
                  <span>阅读趋势</span>
                  <el-select v-model="readTimeRange" size="small" @change="updateReadChart">
                    <el-option label="最近7天" value="7" />
                    <el-option label="最近30天" value="30" />
                    <el-option label="最近90天" value="90" />
                  </el-select>
                </div>
              </template>
              <div ref="readChartRef" class="chart-container"></div>
            </el-card>
          </el-col>
          <el-col :span="12">
            <el-card class="chart-card">
              <template #header>
                <div class="card-header">
                  <span>写作进度</span>
                  <el-select v-model="writeTimeRange" size="small" @change="updateWriteChart">
                    <el-option label="最近7天" value="7" />
                    <el-option label="最近30天" value="30" />
                    <el-option label="最近90天" value="90" />
                  </el-select>
                </div>
              </template>
              <div ref="writeChartRef" class="chart-container"></div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 章节统计 -->
      <div class="chapter-stats">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>章节统计</span>
              <el-radio-group v-model="chapterSortBy" size="small" @change="sortChapterStats">
                <el-radio-button label="reads">按阅读量</el-radio-button>
                <el-radio-button label="comments">按评论数</el-radio-button>
                <el-radio-button label="words">按字数</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <el-table :data="chapterStats" style="width: 100%">
            <el-table-column prop="title" label="章节标题" min-width="200" />
            <el-table-column prop="words" label="字数" width="120">
              <template #default="{ row }">
                {{ formatNumber(row.words) }}
              </template>
            </el-table-column>
            <el-table-column prop="reads" label="阅读量" width="120">
              <template #default="{ row }">
                {{ formatNumber(row.reads) }}
              </template>
            </el-table-column>
            <el-table-column prop="comments" label="评论数" width="120">
              <template #default="{ row }">
                {{ formatNumber(row.comments) }}
              </template>
            </el-table-column>
            <el-table-column prop="updateTime" label="更新时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.updateTime) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  InfoFilled,
  ArrowUp,
  ArrowDown
} from '@element-plus/icons-vue'
import { novelApi } from '../../api/novel'
import { formatDate } from '../../utils/date'
import * as echarts from 'echarts'

const route = useRoute()
const readChartRef = ref(null)
const writeChartRef = ref(null)
let readChart = null
let writeChart = null

// 统计数据
const stats = ref({
  totalWords: 0,
  totalReads: 0,
  totalFavorites: 0,
  totalComments: 0,
  wordsTrend: 0,
  readsTrend: 0,
  favoritesTrend: 0,
  commentsTrend: 0
})

// 图表相关
const readTimeRange = ref('7')
const writeTimeRange = ref('7')
const chapterSortBy = ref('reads')
const chapterStats = ref([])

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await novelApi.getNovelStats(route.params.id)
    stats.value = response.data
    updateCharts()
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  }
}

// 刷新数据
const refreshStats = () => {
  loadStats()
}

// 格式化数字
const formatNumber = (num) => {
  if (!num) return '0'
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 初始化图表
const initCharts = () => {
  if (readChartRef.value) {
    readChart = echarts.init(readChartRef.value)
  }
  if (writeChartRef.value) {
    writeChart = echarts.init(writeChartRef.value)
  }
}

// 更新图表数据
const updateCharts = async () => {
  try {
    const [readData, writeData] = await Promise.all([
      novelApi.getNovelReadTrend(route.params.id, readTimeRange.value),
      novelApi.getNovelWriteTrend(route.params.id, writeTimeRange.value)
    ])
    
    updateReadChart(readData.data)
    updateWriteChart(writeData.data)
  } catch (error) {
    console.error('更新图表数据失败:', error)
    ElMessage.error('更新图表数据失败')
  }
}

// 更新阅读趋势图表
const updateReadChart = (data) => {
  if (!readChart) return

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.dates
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: data.reads,
      type: 'line',
      smooth: true,
      areaStyle: {}
    }]
  }

  readChart.setOption(option)
}

// 更新写作进度图表
const updateWriteChart = (data) => {
  if (!writeChart) return

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: data.dates
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: data.words,
      type: 'line',
      smooth: true,
      areaStyle: {}
    }]
  }

  writeChart.setOption(option)
}

// 加载章节统计
const loadChapterStats = async () => {
  try {
    const response = await novelApi.getNovelChapterStats(route.params.id)
    chapterStats.value = response.data
    sortChapterStats()
  } catch (error) {
    console.error('加载章节统计失败:', error)
    ElMessage.error('加载章节统计失败')
  }
}

// 排序章节统计
const sortChapterStats = () => {
  chapterStats.value.sort((a, b) => b[chapterSortBy.value] - a[chapterSortBy.value])
}

// 监听窗口大小变化
const handleResize = () => {
  readChart?.resize()
  writeChart?.resize()
}

onMounted(() => {
  loadStats()
  loadChapterStats()
  initCharts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  readChart?.dispose()
  writeChart?.dispose()
})
</script>

<style scoped>
.novel-stats {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.stats-cards {
  margin-bottom: 30px;
}

.stats-card {
  height: 160px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  margin: 10px 0;
}

.card-trend {
  color: #f56c6c;
  display: flex;
  align-items: center;
  gap: 5px;
}

.card-trend.up {
  color: #67c23a;
}

.stats-charts {
  margin-bottom: 30px;
}

.chart-card {
  height: 400px;
}

.chart-container {
  height: 300px;
}

.chapter-stats {
  margin-bottom: 30px;
}
</style> 