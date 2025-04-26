<template>
  <div class="chapter-reader">
    <!-- 顶部导航栏 -->
    <div class="reader-header" :class="{ 'hidden': isHeaderHidden }">
      <div class="header-left">
        <el-button @click="goBack">
          <el-icon><Back /></el-icon>
          返回目录
        </el-button>
        <h2 class="novel-title">{{ novel.title }}</h2>
      </div>
      <div class="header-right">
        <el-button-group>
          <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
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
          <span>主题</span>
          <div class="theme-selector">
            <div
              v-for="theme in themes"
              :key="theme.name"
              class="theme-item"
              :class="{ active: currentTheme === theme.name }"
              :style="{ background: theme.background, color: theme.color }"
              @click="setTheme(theme)"
            >
              Aa
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

    <!-- 正文内容 -->
    <div
      class="chapter-content"
      :style="{
        fontSize: `${fontSize}px`,
        lineHeight: lineHeight,
        background: themes[currentTheme].background,
        color: themes[currentTheme].color
      }"
      @click="toggleHeader"
    >
      <h1 class="chapter-title">{{ currentChapter.title }}</h1>
      <div class="content-body" v-html="currentChapter.content"></div>
    </div>

    <!-- 底部导航 -->
    <div class="chapter-nav" :class="{ 'hidden': isHeaderHidden }">
      <el-button @click="prevChapter" :disabled="!hasPrev">上一章</el-button>
      <el-button @click="showSettings = true">设置</el-button>
      <el-button @click="nextChapter" :disabled="!hasNext">下一章</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Back } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'

const route = useRoute()
const router = useRouter()

// 基础数据
const novel = ref({})
const currentChapter = ref({})
const chapters = ref([])
const isHeaderHidden = ref(false)
const showSettings = ref(false)
const showChapters = ref(false)

// 阅读设置
const fontSize = ref(16)
const lineHeight = ref(1.6)
const currentTheme = ref('light')
const themes = {
  light: { name: 'light', background: '#ffffff', color: '#333333' },
  sepia: { name: 'sepia', background: '#f4ecd8', color: '#5c4b51' },
  dark: { name: 'dark', background: '#1a1a1a', color: '#dddddd' }
}

// 获取小说信息
const fetchNovelDetail = async () => {
  try {
    const { data } = await novelApi.getNovelDetail(route.params.novelId)
    novel.value = data
  } catch (error) {
    ElMessage.error('获取小说信息失败')
  }
}

// 获取章节内容
const fetchChapter = async (chapterId) => {
  try {
    const { data } = await novelApi.getChapterDetail(chapterId)
    currentChapter.value = data
    // 更新阅读进度
    updateReadProgress()
  } catch (error) {
    ElMessage.error('获取章节内容失败')
  }
}

// 获取章节列表
const fetchChapterList = async () => {
  try {
    const { data } = await novelApi.getNovelChapters({
      novelId: route.params.novelId,
      limit: 1000 // 获取所有章节
    })
    chapters.value = data.chapters
  } catch (error) {
    ElMessage.error('获取章节列表失败')
  }
}

// 更新阅读进度
const updateReadProgress = async () => {
  try {
    await novelApi.updateReadProgress({
      novelId: route.params.novelId,
      chapterId: currentChapter.value.id
    })
  } catch (error) {
    console.error('更新阅读进度失败', error)
  }
}

// 计算属性
const hasPrev = computed(() => {
  const index = chapters.value.findIndex(c => c.id === currentChapter.value.id)
  return index > 0
})

const hasNext = computed(() => {
  const index = chapters.value.findIndex(c => c.id === currentChapter.value.id)
  return index < chapters.value.length - 1
})

// 章节导航
const prevChapter = () => {
  if (hasPrev.value) {
    const index = chapters.value.findIndex(c => c.id === currentChapter.value.id)
    goToChapter(chapters.value[index - 1].id)
  }
}

const nextChapter = () => {
  if (hasNext.value) {
    const index = chapters.value.findIndex(c => c.id === currentChapter.value.id)
    goToChapter(chapters.value[index + 1].id)
  }
}

const goToChapter = (chapterId) => {
  router.push(`/novel/${route.params.novelId}/chapter/${chapterId}`)
  showChapters.value = false
}

// 返回小说页面
const goBack = () => {
  router.push(`/novel/${route.params.novelId}`)
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

const setTheme = (theme) => {
  currentTheme.value = theme.name
  localStorage.setItem('reader-theme', theme.name)
}

// 加载保存的阅读设置
const loadReaderSettings = () => {
  fontSize.value = parseInt(localStorage.getItem('reader-font-size')) || 16
  lineHeight.value = parseFloat(localStorage.getItem('reader-line-height')) || 1.6
  currentTheme.value = localStorage.getItem('reader-theme') || 'light'
}

// 键盘快捷键
const handleKeyPress = (e) => {
  if (e.key === 'ArrowLeft') {
    prevChapter()
  } else if (e.key === 'ArrowRight') {
    nextChapter()
  }
}

onMounted(async () => {
  loadReaderSettings()
  await fetchNovelDetail()
  await fetchChapterList()
  await fetchChapter(route.params.chapterId)
  window.addEventListener('keydown', handleKeyPress)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyPress)
})

// 监听路由变化
watch(
  () => route.params.chapterId,
  (newChapterId) => {
    if (newChapterId) {
      fetchChapter(newChapterId)
    }
  }
)
</script>

<style scoped>
.chapter-reader {
  height: 100vh;
  display: flex;
  flex-direction: column;
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

.reader-header.hidden {
  transform: translateY(-100%);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.novel-title {
  margin: 0;
  font-size: 18px;
}

.chapter-content {
  flex: 1;
  padding: 60px 20px;
  max-width: 800px;
  margin: 0 auto;
  transition: all 0.3s;
}

.chapter-title {
  text-align: center;
  margin-bottom: 30px;
}

.content-body {
  text-align: justify;
  word-break: break-word;
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

.theme-selector {
  display: flex;
  gap: 15px;
  margin-top: 10px;
}

.theme-item {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: 2px solid transparent;
}

.theme-item.active {
  border-color: var(--el-color-primary);
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
</style>