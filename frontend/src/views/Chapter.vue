<template>
  <div class="chapter-reader">
    <!-- 顶部导航栏 -->
    <div class="reader-header" :class="{ 'hidden': isHeaderHidden }">
      <div class="header-left">
        <el-button @click="goBack">
          <el-icon><Back /></el-icon>
          返回目录
        </el-button>
        <h2 class="novel-title">{{ novel?.title || '' }}</h2>
      </div>
      <div class="header-right">
        <el-button-group>
          <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
          <el-button @click="showSettings = true">
            <el-icon><Setting /></el-icon>
          </el-button>
          <el-button @click="showChapterList">目录</el-button>
          <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
        </el-button-group>
      </div>
    </div>

    <!-- 阅读设置抽屉 -->
    <el-drawer
      v-model="showSettings"
      title="阅读设置"
      direction="rtl"
      size="300px"
    >
      <div class="reader-settings">
        <div class="setting-item">
          <span>字体大小</span>
          <el-slider
            v-model="fontSize"
            :min="12"
            :max="24"
            :step="1"
            @change="updateFontSize"
          />
        </div>
        <div class="setting-item">
          <span>行间距</span>
          <el-slider
            v-model="lineHeight"
            :min="1.2"
            :max="2"
            :step="0.1"
            @change="updateLineHeight"
          />
        </div>
        <div class="setting-item">
          <span>内容宽度</span>
          <el-slider
            v-model="contentWidth"
            :min="600"
            :max="2400"
            :step="100"
            @change="updateContentWidth"
          />
        </div>
        <div class="setting-item">
          <span>分栏数量</span>
          <el-slider
            v-model="columnCount"
            :min="1"
            :max="3"
            :step="1"
            :marks="{1: '1列', 2: '2列', 3: '3列'}"
            @change="updateColumnCount"
          />
        </div>
        <div class="setting-item" v-if="columnCount > 1">
          <span>分栏间距</span>
          <el-slider
            v-model="columnGap"
            :min="40"
            :max="120"
            :step="20"
            @change="updateColumnGap"
          />
        </div>
        <div class="setting-item">
          <span>主题</span>
          <div class="theme-grid">
            <div
              v-for="theme in themes"
              :key="theme.name"
              class="theme-item"
              :class="{ active: currentTheme === theme.name }"
              :style="{ 
                background: theme.background, 
                color: theme.color,
                borderColor: currentTheme === theme.name ? theme.color : 'transparent'
              }"
              @click="setTheme(theme)"
            >
              <span class="theme-text">Aa</span>
              <span class="theme-name">{{ theme.description }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- 章节列表抽屉 -->
    <el-drawer
      v-model="showChapters"
      title="章节目录"
      direction="rtl"
      size="300px"
    >
      <el-scrollbar height="100%">
        <div
          v-for="chapter in chapters"
          :key="chapter.id"
          class="chapter-item"
          :class="{ active: chapter.id === currentChapter.id }"
          @click="goToChapter(chapter.id)"
        >
          {{ chapter.title }}
        </div>
      </el-scrollbar>
    </el-drawer>

    <!-- 正文内容容器 -->
    <div 
      class="content-container"
      @wheel.prevent="handleWheel"
      @mousedown="handleMouseDown"
    >
      <div
        class="chapter-content"
        :style="{
          fontSize: `${fontSize}px`,
          lineHeight: lineHeight,
          background: themes[currentTheme].background,
          color: themes[currentTheme].color,
          width: `${actualContentWidth}px`,
          minWidth: '300px',
          columnCount: columnCount,
          columnGap: `${columnGap}px`,
          columnRule: columnCount > 1 ? `1px solid ${themes[currentTheme].color}` : 'none',
          columnFill: 'auto',
          height: 'calc(100vh - 160px)',
          margin: 'auto',
          padding: '20px 0'
        }"
        @click="toggleHeader"
      >
        <div 
          class="content-wrapper"
          :style="{
            columnGap: `${columnGap}px`,
            padding: '0 40px',
            margin: '0 auto'
          }"
        >
          <h1 class="chapter-title">{{ currentChapter?.title || '' }}</h1>
          <div 
            class="content-body" 
            v-html="currentChapter?.content || ''"
          ></div>
        </div>
        <!-- Add page numbers -->
        <div 
          class="page-numbers"
          :class="{ 'dark-theme': currentTheme === 'dark' }"
          v-if="totalPages > 1"
        >
          {{ currentPage }} / {{ totalPages }}
        </div>
      </div>

      <!-- Update pagination indicator -->
      <div 
        class="pagination-indicator" 
        v-if="columnCount > 1 || totalPages > 1" 
        :class="{ 'hidden': isHeaderHidden }"
      >
        <el-button @click="handlePageChange('prev')" :disabled="currentPage === 1">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
        <el-button @click="handlePageChange('next')" :disabled="currentPage >= totalPages">
          <el-icon><ArrowRight /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 底部导航 -->
    <div class="chapter-nav" :class="{ 'hidden': isHeaderHidden }">
      <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
      <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Back, ArrowLeft, ArrowRight, Setting } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'

const route = useRoute()
const router = useRouter()

// 基础数据
const novel = ref({
  title: '',
  // Add other novel properties you need with default values
})
const currentChapter = ref({
  id: null,
  title: '',
  content: ''
})
const chapters = ref([])
const isHeaderHidden = ref(false)
const showSettings = ref(false)
const showChapters = ref(false)

// 阅读设置
const fontSize = ref(16)
const lineHeight = ref(1.6)
const contentWidth = ref(800)
const columnCount = ref(1)
const columnGap = ref(60)
const currentPage = ref(1)
const currentTheme = ref('light')
const themes = {
  light: { 
    name: 'light', 
    background: '#ffffff', 
    color: '#333333',
    description: '默认模式'
  },
  sepia: { 
    name: 'sepia', 
    background: '#f4ecd8', 
    color: '#5c4b51',
    description: '护眼黄'
  },
  dark: { 
    name: 'dark', 
    background: '#1a1a1a', 
    color: '#dddddd',
    description: '暗黑'
  },
  green: {
    name: 'green',
    background: '#e6f3ec', 
    color: '#2c4a3c', 
    description: '护眼绿'
  },
  blue: {
    name: 'blue',
    background: '#f0f8ff',
    color: '#1e4e8c',
    description: '清新蓝'
  },
  pink: {
    name: 'pink',
    background: '#fff0f5', 
    color: '#c25975', 
    description: '浪漫粉'
  }
}

// 添加总页数计算
const totalPages = ref(1)

// 计算实际内容宽度
const actualContentWidth = computed(() => {
  // 获取容器宽度（考虑padding）
  const containerWidth = window.innerWidth - 40 // 减去左右各20px的padding
  
  // 单列模式下的宽度计算
  if (columnCount.value === 1) {
    // 使用较小的最大宽度限制，并确保不超过容器宽度
    return Math.min(contentWidth.value, containerWidth, 800)
  }
  
  // 多列模式下的宽度计算
  const availableWidth = containerWidth * 0.95 // 预留一些边距
  const singleColumnWidth = Math.min(700, (availableWidth - (columnCount.value - 1) * columnGap.value) / columnCount.value)
  return singleColumnWidth * columnCount.value + (columnCount.value - 1) * columnGap.value
})

// 监听窗口大小变化
const handleResize = () => {
  calculateTotalPages()
  // 重置滚动位置到当前页
  nextTick(() => {
    scrollToPage(currentPage.value)
  })
}

// 计算总页数
const calculateTotalPages = () => {
  const content = document.querySelector('.chapter-content')
  if (!content) return
  
  // 使用 scrollWidth 和 clientWidth 计算总页数
  totalPages.value = Math.ceil(content.scrollWidth / content.clientWidth)
}

// 处理鼠标滚轮事件
const handleWheel = (e) => {
  // 阻止默认滚动行为
  e.preventDefault()
  
  // 设置滚动阈值和冷却时间
  const scrollThreshold = 50 // 滚动阈值
  const now = Date.now()
  
  // 如果距离上次滚动时间太短，则忽略此次滚动
  if (now - lastWheelTime.value < 200) { // 200ms冷却时间
    return
  }
  
  // 根据滚动方向翻页
  if (Math.abs(e.deltaY) > scrollThreshold) {
    if (e.deltaY > 0) {
      if (currentPage.value < totalPages.value) {
        currentPage.value++
        scrollToPage(currentPage.value)
      } else if (e.deltaY > 100) { // 在最后一页继续滚动时切换到下一章
        nextChapter()
      }
    } else {
      if (currentPage.value > 1) {
        currentPage.value--
        scrollToPage(currentPage.value)
      } else if (e.deltaY < -100) { // 在第一页继续向上滚动时切换到上一章
        prevChapter()
      }
    }
    lastWheelTime.value = now
  }
}

// 添加最后一次滚轮事件时间记录
const lastWheelTime = ref(0)

// 滚动到指定页
const scrollToPage = (pageNumber) => {
  const content = document.querySelector('.chapter-content')
  if (!content) return
  
  const scrollPosition = (pageNumber - 1) * content.clientWidth
  content.scrollTo({
    left: scrollPosition,
    behavior: 'smooth'
  })
}

// 处理鼠标按下事件
const handleMouseDown = (e) => {
  // 检查是否为鼠标中键
  if (e.button === 1) {
    e.preventDefault()
    // 记录初始Y坐标
    const startY = e.clientY
    
    const handleMouseMove = (moveEvent) => {
      const deltaY = moveEvent.clientY - startY
      if (Math.abs(deltaY) > 50) { // 移动超过50px才触发翻页
        if (deltaY > 0) {
          prevChapter()
        } else {
          nextChapter()
        }
        document.removeEventListener('mousemove', handleMouseMove)
        document.removeEventListener('mouseup', handleMouseUp)
      }
    }
    
    const handleMouseUp = () => {
      document.removeEventListener('mousemove', handleMouseMove)
      document.removeEventListener('mouseup', handleMouseUp)
    }
    
    document.addEventListener('mousemove', handleMouseMove)
    document.addEventListener('mouseup', handleMouseUp)
  }
}

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
const fetchChapterList = async () => {
  try {
    const response = await novelApi.getNovelChapters({
      novelId: route.params.id,
      limit: 1000 // 获取所有章节
    })
    
    // 直接使用response，因为data已经在API层被解构
    if (response && Array.isArray(response.chapters)) {
      // 按order字段排序
      chapters.value = response.chapters.sort((a, b) => a.order - b.order)
    } else {
      chapters.value = []
      ElMessage.error('获取章节列表失败：数据格式错误')
    }
  } catch (error) {
    ElMessage.error('获取章节列表失败')
    chapters.value = []
  }
}

// 获取章节内容
const fetchChapter = async (chapterId) => {
  try {
    const response = await novelApi.getChapterDetail(chapterId)
    // 直接使用response，因为data已经在API层被解构
    if (response) {
      currentChapter.value = {
        id: response.id,
        title: response.title,
        content: response.content,
        novelId: response.novelId,
        order: response.order
      }
      // 更新阅读进度
      updateReadProgress()
    }
  } catch (error) {
    ElMessage.error('获取章节内容失败')
  }
}

// 更新阅读进度
const updateReadProgress = async () => {
  try {
    if (!currentChapter.value?.id || !currentChapter.value?.novelId) {
      return; // Skip if no chapter id or novel id available
    }
    await novelApi.updateReadProgress({
      novelId: currentChapter.value.novelId,
      chapterId: currentChapter.value.id
    })
  } catch (error) {
    console.error('更新阅读进度失败', error)
  }
}

// 计算属性
const hasPrev = computed(() => {
  if (!currentChapter.value?.order || !chapters.value.length) return false
  return chapters.value.some(chapter => chapter.order === currentChapter.value.order - 1)
})

const hasNext = computed(() => {
  if (!currentChapter.value?.order || !chapters.value.length) return false
  return chapters.value.some(chapter => chapter.order === currentChapter.value.order + 1)
})

// 章节导航
const prevChapter = async () => {
  if (hasPrev.value && currentChapter.value?.order) {
    const prevChapter = chapters.value.find(c => c.order === currentChapter.value.order - 1)
    if (prevChapter) {
      await router.push(`/novels/${route.params.id}/chapter/${prevChapter.id}`)
    }
  }
}

const nextChapter = async () => {
  if (hasNext.value && currentChapter.value?.order) {
    const nextChapter = chapters.value.find(c => c.order === currentChapter.value.order + 1)
    if (nextChapter) {
      await router.push(`/novels/${route.params.id}/chapter/${nextChapter.id}`)
    }
  }
}

const goToChapter = (chapterId) => {
  router.push(`/novels/${route.params.id}/chapter/${chapterId}`)
  showChapters.value = false
}

// 返回小说页面
const goBack = () => {
  router.push(`/novels/${route.params.id}`)
}

// 显示章节列表
const showChapterList = () => {
  showChapters.value = true
}

// 切换顶部导航栏显示
const toggleHeader = () => {
  isHeaderHidden.value = !isHeaderHidden.value
}

// 更新阅读设置
const updateFontSize = (size) => {
  localStorage.setItem('reader-font-size', size)
}

const updateLineHeight = (height) => {
  localStorage.setItem('reader-line-height', height)
}

const updateContentWidth = (width) => {
  localStorage.setItem('reader-content-width', width)
}

const updateColumnCount = (count) => {
  localStorage.setItem('reader-column-count', count)
}

const updateColumnGap = (gap) => {
  localStorage.setItem('reader-column-gap', gap)
}

const setTheme = (theme) => {
  currentTheme.value = theme.name
  localStorage.setItem('reader-theme', theme.name)
}

// 加载保存的阅读设置
const loadReaderSettings = () => {
  fontSize.value = parseInt(localStorage.getItem('reader-font-size')) || 16
  lineHeight.value = parseFloat(localStorage.getItem('reader-line-height')) || 1.6
  contentWidth.value = parseInt(localStorage.getItem('reader-content-width')) || 800
  columnCount.value = parseInt(localStorage.getItem('reader-column-count')) || 1
  columnGap.value = parseInt(localStorage.getItem('reader-column-gap')) || 60
  currentTheme.value = localStorage.getItem('reader-theme') || 'light'
}

// 翻页处理
const handlePageChange = (direction) => {
  const content = document.querySelector('.chapter-content')
  if (!content) return

  const totalWidth = content.scrollWidth
  const viewWidth = content.clientWidth
  const maxPage = Math.ceil(totalWidth / viewWidth)

  if (direction === 'next' && currentPage.value < maxPage) {
    currentPage.value++
    content.scrollTo({
      left: (currentPage.value - 1) * viewWidth,
      behavior: 'smooth'
    })
  } else if (direction === 'prev' && currentPage.value > 1) {
    currentPage.value--
    content.scrollTo({
      left: (currentPage.value - 1) * viewWidth,
      behavior: 'smooth'
    })
  }
}

// 键盘快捷键
const handleKeyPress = (e) => {
  if (e.key === 'ArrowLeft') {
    if (e.altKey) {
      prevChapter()
    } else {
      handlePageChange('prev')
    }
  } else if (e.key === 'ArrowRight') {
    if (e.altKey) {
      nextChapter()
    } else {
      handlePageChange('next')
    }
  }
}

// 监听内容变化重新计算总页数
const observer = ref(null)

onMounted(async () => {
  loadReaderSettings()
  window.addEventListener('resize', handleResize)
  window.addEventListener('wheel', handleWheel, { passive: false })
  handleResize()
  await fetchNovelDetail()
  await fetchChapterList()
  await fetchChapter(route.params.chapterId)
  window.addEventListener('keydown', handleKeyPress)
  
  // 创建 ResizeObserver 监听内容变化
  observer.value = new ResizeObserver(() => {
    calculateTotalPages()
  })
  
  const content = document.querySelector('.chapter-content')
  if (content) {
    observer.value.observe(content)
  }

  // 自动滚动到正确位置
  window.scrollTo({
    top: 100,
    behavior: 'instant'
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('wheel', handleWheel)
  window.removeEventListener('keydown', handleKeyPress)
  
  // 清理 ResizeObserver
  if (observer.value) {
    observer.value.disconnect()
  }
})

// 监听路由参数变化
watch(
  () => route.params.chapterId,
  async (newChapterId, oldChapterId) => {
    if (newChapterId && newChapterId !== oldChapterId) {
      currentPage.value = 1
      await fetchChapter(newChapterId)
      // 等待内容加载完成后计算总页数
      nextTick(() => {
        calculateTotalPages()
        scrollToPage(1)
      })
    }
  },
  { immediate: true }
)
</script>

<style scoped>
/* Override main-content padding for chapter view */
:deep(.main-content) {
  padding-top: 0;
  min-height: 100vh;
}

.chapter-reader {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  max-width: 100vw;
  position: relative;
}

.reader-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 10px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: transform 0.3s;
}

.header-right .el-button-group {
  display: flex;
  gap: 5px;
}

.content-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  box-sizing: border-box;
  padding: 0px 20px 80px;
  min-height: 100vh;
}

.chapter-content {
  box-sizing: border-box;
  overflow-y: auto;
  overflow-x: hidden;
  column-rule-width: 2px;
  column-rule-style: solid;
  column-rule-color: rgba(0, 0, 0, 0.1);
  max-width: 100%;
  height: calc(100vh - 160px);
  margin: auto;
  scroll-behavior: smooth;
  padding: 20px 0;
  position: relative;
}

.content-wrapper {
  break-inside: avoid;
  max-width: 100%;
  padding: 0 40px;
  margin: 0 auto;
}

.content-body {
  text-align: justify;
  word-break: break-word;
  max-width: 100%;
  overflow-wrap: break-word;
  hyphens: auto;
}

.content-body :deep(p) {
  margin: 1.2em 0;
  text-indent: 2em;
  break-inside: avoid;
  position: relative;
  max-width: 100%;
  line-height: inherit;
}

.chapter-title {
  text-align: center;
  margin-bottom: 30px;
  break-after: avoid;
  break-inside: avoid;
  font-size: 1.5em;
  padding: 0 20px;
}

.pagination-indicator {
  position: fixed;
  bottom: 80px;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.9);
  padding: 5px 10px;
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: opacity 0.3s;
  z-index: 100;
}

.pagination-indicator.hidden {
  opacity: 0;
  pointer-events: none;
}

.page-info {
  min-width: 30px;
  text-align: center;
  font-size: 14px;
}

.chapter-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  padding: 10px;
  display: flex;
  justify-content: center;
  gap: 20px;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s;
  z-index: 100;
}

.chapter-nav.hidden {
  transform: translateY(100%);
}

.reader-settings {
  padding: 20px;
}

.setting-item {
  margin-bottom: 30px;
}

.theme-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  margin-top: 10px;
}

.theme-item {
  height: 80px;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  position: relative;
}

.theme-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.theme-text {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 4px;
}

.theme-name {
  font-size: 12px;
  opacity: 0.8;
}

.theme-item.active {
  border-width: 2px;
  border-style: solid;
}

.chapter-item {
  padding: 12px 20px;
  cursor: pointer;
  transition: all 0.3s;
}

.chapter-item:hover {
  background: #f5f7fa;
}

.chapter-item.active {
  color: var(--el-color-primary);
  background: var(--el-color-primary-light-9);
}

.chapter-content[style*="background: #ffffff"] {
  column-rule-color: rgba(0, 0, 0, 0.1);
}

.chapter-content[style*="background: #1a1a1a"] {
  column-rule-color: rgba(255, 255, 255, 0.1);
}

.chapter-content[style*="background: #f4ecd8"] {
  column-rule-color: rgba(92, 75, 81, 0.1);
}

.chapter-content[style*="background: #e6f3ec"] {
  column-rule-color: rgba(44, 74, 60, 0.1);
}

.chapter-content[style*="background: #f0f8ff"] {
  column-rule-color: rgba(30, 78, 140, 0.1);
}

.chapter-content[style*="background: #fff0f5"] {
  column-rule-color: rgba(194, 89, 117, 0.1);
}

.page-numbers {
  position: absolute;
  bottom: 10px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(255, 255, 255, 0.8);
  padding: 4px 12px;
  border-radius: 15px;
  font-size: 14px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  z-index: 10;
  transition: opacity 0.3s;
}

.dark-theme .page-numbers {
  background: rgba(0, 0, 0, 0.8);
  color: #fff;
}

/* Responsive design */
@media screen and (max-width: 1440px) {
  .chapter-content {
    width: 90%;
    max-width: 1200px;
  }
}

@media screen and (max-width: 1024px) {
  .chapter-content {
    width: 95%;
    max-width: 900px;
  }
  
  .content-wrapper {
    padding: 0 30px;
  }
}

@media screen and (max-width: 768px) {
  .content-container {
    padding: 50px 10px 70px;
  }

  .chapter-content {
    width: 100%;
    height: calc(100vh - 140px);
  }

  .content-wrapper {
    padding: 0 20px;
  }

  .chapter-title {
    font-size: 1.2em;
    margin-bottom: 20px;
  }

  .page-numbers {
    font-size: 12px;
    padding: 3px 10px;
  }
}

@media screen and (max-width: 480px) {
  .content-container {
    padding: 40px 5px 60px;
  }

  .chapter-content {
    height: calc(100vh - 120px);
  }

  .content-wrapper {
    padding: 0 15px;
  }

  .content-body :deep(p) {
    margin: 1em 0;
    text-indent: 1.5em;
  }
}
</style>