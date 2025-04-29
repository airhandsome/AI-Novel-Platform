<template>
  <div class="home">
    <!-- 轮播图部分 -->
    <div class="banner">
      <el-carousel height="500px">
        <el-carousel-item v-for="item in banners" :key="item.id">
          <div class="banner-content" :style="{ backgroundImage: `url(${item.image})` }">
            <div class="banner-text glass-effect">
              <h2 class="display-large">{{ item.title }}</h2>
              <p class="title-medium">{{ item.description }}</p>
              <el-button type="primary" size="large" class="banner-button" @click="navigateToSection(item.action)">
                {{ item.buttonText }}
              </el-button>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 内容区域 -->
    <div class="content">
      <!-- 推荐小说 -->
      <section class="section card-container" v-loading="loading.recommended">
        <div class="section-header">
          <h2 class="display-small">编辑推荐</h2>
          <el-button text class="title-medium" @click="$router.push('/library')">查看更多</el-button>
        </div>
        <el-empty v-if="recommendedBooks.length === 0 && !loading.recommended" description="暂无推荐小说" />
        <el-row :gutter="24" v-else>
          <el-col :span="4" v-for="book in recommendedBooks" :key="book.id">
            <div class="book-card shadow-md" @click="$router.push(`/novels/${book.id}`)">
              <div class="book-cover">
                <img :src="book.cover || '/default-cover.jpg'" :alt="book.title" />
                <div class="book-overlay">
                  <el-button circle icon="Plus" class="add-button" @click.stop="addToFavorite(book.id)" />
                </div>
              </div>
              <div class="book-info">
                <h3 class="title-medium">{{ book.title }}</h3>
                <p class="body-medium author">✍️ {{ book.author }}</p>
                <p class="caption description">{{ book.description }}</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- 最新更新 -->
      <section class="section card-container" v-loading="loading.latestUpdates">
        <div class="section-header">
          <h2 class="display-small">最新更新</h2>
          <el-button text class="title-medium" @click="$router.push('/library?sort=latest')">查看更多</el-button>
        </div>
        <el-empty v-if="latestUpdates.length === 0 && !loading.latestUpdates" description="暂无最新更新" />
        <div class="latest-updates" v-else>
          <div v-for="update in latestUpdates" :key="update.id" class="update-item">
            <div class="book-meta">
              <router-link :to="`/novels/${update.id}`" class="book-title title-medium">
                {{ update.title }}
              </router-link>
              <span class="author body-medium">✍️ {{ update.author }}</span>
            </div>
            <router-link 
              :to="`/novels/${update.id}/chapter/${update.latestChapterId}`" 
              class="chapter-link body-medium"
            >
              {{ update.latestChapter }}
            </router-link>
            <span class="update-time caption">{{ formatDate(update.updateTime) }}</span>
          </div>
        </div>
      </section>

      <!-- AI创作专区 -->
      <section class="section ai-section card-container">
        <div class="section-header">
          <h2 class="display-small">AI创作专区</h2>
          <el-button type="primary" size="large" @click="$router.push('/write')">开始创作</el-button>
        </div>
        <el-row :gutter="24">
          <el-col :span="8" v-for="feature in aiFeatures" :key="feature.id">
            <div class="feature-card glass-effect">
              <div class="feature-icon">
                <el-icon :size="40">
                  <component :is="feature.icon" />
                </el-icon>
              </div>
              <h3 class="title-large">{{ feature.title }}</h3>
              <p class="body-medium">{{ feature.description }}</p>
            </div>
          </el-col>
        </el-row>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Edit, Brush, Tools } from '@element-plus/icons-vue'
import { novelApi } from '../api/novel'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

// 加载状态
const loading = ref({
  recommended: false,
  latestUpdates: false
})

// 轮播图数据
const banners = ref([
  {
    id: 1,
    image: '/banner1.jpg',
    title: '发现优质小说',
    description: '数千部精选小说，总有一本适合你',
    buttonText: '立即探索',
    action: 'library'
  },
  {
    id: 2,
    image: '/banner2.jpg',
    title: 'AI智能创作',
    description: '让创作更简单，让故事更精彩',
    buttonText: '开始创作',
    action: 'write'
  },
  {
    id: 3,
    image: '/banner3.jpg',
    title: '作者社区',
    description: '和其他作者交流，分享创作经验',
    buttonText: '加入社区',
    action: 'community'
  }
])

// 推荐小说数据
const recommendedBooks = ref([])

// 最新更新数据
const latestUpdates = ref([])

// AI功能特性
const aiFeatures = ref([
  {
    id: 1,
    icon: 'Edit',
    title: '智能写作助手',
    description: '提供情节发展建议，帮助你突破创作瓶颈'
  },
  {
    id: 2,
    icon: 'Brush',
    title: '角色塑造',
    description: '智能分析角色性格，提供人物对话建议'
  },
  {
    id: 3,
    icon: 'Tools',
    title: '场景描写',
    description: '生成丰富的场景描写，让故事更加生动'
  }
])

// 导航到不同区域
const navigateToSection = (action) => {
  switch (action) {
    case 'library':
      router.push('/library')
      break
    case 'write':
      router.push('/write')
      break
    case 'community':
      router.push('/community')
      break
    default:
      router.push('/library')
  }
}

// 加载推荐小说
const loadRecommendedBooks = async () => {
  try {
    loading.value.recommended = true
    const response = await novelApi.getNovelList({
      page: 1,
      pageSize: 6,
      sortBy: 'recommend'
    })
    
    recommendedBooks.value = response.novels.map(novel => ({
      id: novel.id,
      title: novel.title,
      description: novel.description,
      author: novel.author.Username,
      cover: novel.coverUrl || '/default-cover.jpg',
      status: novel.status
    }))
  } catch (error) {
    console.error('加载推荐小说失败:', error)
    ElMessage.error('加载推荐小说失败')
  } finally {
    loading.value.recommended = false
  }
}

// 加载最新更新
const loadLatestUpdates = async () => {
  try {
    loading.value.latestUpdates = true
    const response = await novelApi.getNovelList({
      page: 1,
      pageSize: 5,
      sortBy: 'latest'
    })
    
    latestUpdates.value = response.novels.map(novel => ({
      id: novel.id,
      title: novel.title,
      author: novel.author.Username,
      latestChapter: novel.latestChapter?.title || '暂无章节',
      latestChapterId: novel.latestChapter?.id || 0,
      updateTime: novel.updateTime || novel.createTime
    }))
  } catch (error) {
    console.error('加载最新更新失败:', error)
    ElMessage.error('加载最新更新失败')
  } finally {
    loading.value.latestUpdates = false
  }
}

// 添加到收藏
const addToFavorite = async (novelId) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  try {
    await novelApi.favoriteNovel(novelId)
    ElMessage.success('已添加到收藏')
  } catch (error) {
    console.error('添加收藏失败:', error)
    ElMessage.error('添加收藏失败')
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  
  // 一天内
  if (diff < 24 * 60 * 60 * 1000) {
    const hours = Math.floor(diff / (60 * 60 * 1000))
    if (hours < 1) {
      const minutes = Math.floor(diff / (60 * 1000))
      return minutes <= 0 ? '刚刚' : `${minutes}分钟前`
    }
    return `${hours}小时前`
  }
  
  // 一周内
  if (diff < 7 * 24 * 60 * 60 * 1000) {
    const days = Math.floor(diff / (24 * 60 * 60 * 1000))
    return `${days}天前`
  }
  
  // 一年内
  if (date.getFullYear() === now.getFullYear()) {
    return `${date.getMonth() + 1}月${date.getDate()}日`
  }
  
  // 超过一年
  return `${date.getFullYear()}/${date.getMonth() + 1}/${date.getDate()}`
}

onMounted(() => {
  loadRecommendedBooks()
  loadLatestUpdates()
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

.home {
  min-height: 100vh;
  background-color: #f9f9f9;
}

.banner {
  margin-bottom: 60px;
  margin-top: -60px;
}

.banner-content {
  height: 100%;
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  padding: 0 max(20px, calc((100% - 1200px) / 2));
  position: relative;
}

.banner-content::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to right, rgba(46, 125, 50, 0.8), rgba(27, 94, 32, 0.6));
}

.banner-text {
  position: relative;
  max-width: 600px;
  padding: 40px;
  border-radius: 24px;
  color: white;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.banner-text h2 {
  margin: 0 0 16px;
}

.banner-text p {
  margin: 0 0 24px;
  opacity: 0.9;
}

.banner-button {
  font-weight: 600;
  padding: 12px 32px;
  background: var(--primary);
  border: none;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(27, 94, 32, 0.3);
  transition: transform 0.2s, box-shadow 0.2s;
}

.banner-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(27, 94, 32, 0.4);
  background: var(--dark);
}

.content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section {
  margin-bottom: 60px;
}

.card-container {
  background: white;
  border-radius: 16px;
  padding: 30px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(46, 125, 50, 0.1);
  position: relative;
  overflow: hidden;
}

.card-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--secondary));
  opacity: 0.8;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
}

.section-header h2 {
  color: var(--dark);
  margin: 0;
  font-weight: 600;
}

.section-header .el-button {
  color: var(--primary);
}

.book-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e0e0e0;
  position: relative;
  height: 380px;
  display: flex;
  flex-direction: column;
}

.book-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(46, 125, 50, 0.15);
  border-color: var(--primary);
}

.book-cover {
  position: relative;
  height: 240px;
  flex-shrink: 0;
  overflow: hidden;
}

.book-cover::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    transparent 70%,
    rgba(0, 0, 0, 0.3) 100%
  );
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

.book-overlay {
  position: absolute;
  top: 10px;
  right: 10px;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 2;
}

.book-card:hover .book-overlay {
  opacity: 1;
}

.add-button {
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  color: var(--primary);
  border: none;
}

.add-button:hover {
  background: var(--primary);
  color: white;
  transform: scale(1.1);
}

.book-info {
  padding: 20px;
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  min-height: 140px;
}

.book-info h3 {
  color: var(--text-primary);
  margin: 0 0 8px;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
  font-size: 1.25rem;
}

.author {
  color: var(--text-secondary);
  margin-bottom: 8px;
  flex-shrink: 0;
  font-size: 0.95rem;
}

.description {
  color: var(--text-tertiary);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  flex-grow: 1;
  font-size: 0.9rem;
  line-height: 1.5;
}

.latest-updates {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e0e0e0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.update-item {
  display: flex;
  align-items: center;
  padding: 20px;
  gap: 24px;
  border-bottom: 1px solid rgba(46, 125, 50, 0.1);
  transition: all 0.3s ease;
  position: relative;
  background: white;
}

.update-item:last-child {
  border-bottom: none;
}

.update-item:hover {
  background-color: #f8fffa;
  padding-left: 24px;
}

.update-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--primary);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.update-item:hover::before {
  opacity: 1;
}

.book-meta {
  width: 240px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.book-title {
  color: var(--text-primary);
  text-decoration: none;
  transition: color 0.3s;
  font-weight: 500;
}

.book-title:hover {
  color: var(--primary);
}

.chapter-link {
  flex: 1;
  color: var(--text-primary);
  text-decoration: none;
  transition: color 0.3s;
}

.chapter-link:hover {
  color: var(--primary);
}

.update-time {
  color: var(--text-tertiary);
  min-width: 120px;
  text-align: right;
  background: #f8f9fa;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  border: 1px solid rgba(46, 125, 50, 0.1);
}

.ai-section {
  margin-bottom: 80px;
}

.feature-card {
  padding: 40px;
  border-radius: 16px;
  text-align: center;
  transition: all 0.3s ease;
  border: 1px solid #e0e0e0;
  background: white;
  position: relative;
  overflow: hidden;
  border-top: 3px solid var(--primary);
}

.feature-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 10px 25px rgba(46, 125, 50, 0.15);
  background-color: #f8fffa;
}

.feature-icon {
  margin-bottom: 24px;
  color: var(--primary);
  position: relative;
}

.feature-icon::after {
  content: '';
  position: absolute;
  width: 60px;
  height: 60px;
  background: var(--secondary);
  opacity: 0.1;
  border-radius: 50%;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  transition: all 0.3s ease;
}

.feature-card:hover .feature-icon::after {
  width: 70px;
  height: 70px;
  opacity: 0.15;
}

.feature-card h3 {
  color: var(--dark);
  margin: 0 0 16px;
  font-weight: 500;
}

.feature-card p {
  color: #555;
  margin: 0;
}

/* 移动端适配 */
@media (max-width: 900px) {
  .card-container {
    padding: 20px;
  }
  
  .feature-card {
    padding: 30px;
  }
}

@media (max-width: 600px) {
  .card-container {
    padding: 16px;
  }
  
  .feature-card {
    padding: 24px;
  }
  
  .update-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .book-meta {
    width: 100%;
  }
  
  .update-time {
    width: auto;
  }
}

/* 暗色模式 */
@media (prefers-color-scheme: dark) {
  .home {
    background-color: var(--system-background);
  }

  .banner-content::before {
    background: linear-gradient(to right, rgba(27, 94, 32, 0.9), rgba(20, 70, 24, 0.7));
  }

  .banner-text {
    background: rgba(30, 30, 30, 0.4);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .banner-button {
    background: var(--primary);
    box-shadow: 0 4px 12px rgba(27, 94, 32, 0.2);
  }

  .banner-button:hover {
    box-shadow: 0 6px 16px rgba(27, 94, 32, 0.3);
  }

  .card-container {
    background: var(--system-grouped-secondary);
    border-color: rgba(255, 255, 255, 0.08);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
  }

  .section-header {
    border-color: rgba(255, 255, 255, 0.08);
  }

  .section-header h2 {
    color: #000;
  }

  .section-header .el-button {
    color: var(--secondary);
  }

  .book-card {
    background: var(--system-grouped-secondary);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .book-card:hover {
    border-color: var(--primary);
    box-shadow: 0 10px 25px rgba(46, 125, 50, 0.2);
  }

  .add-button {
    background: rgba(60, 60, 60, 0.9);
    color: var(--secondary);
  }

  .add-button:hover {
    background: var(--primary);
    color: white;
  }

  .latest-updates {
    border-color: rgba(255, 255, 255, 0.1);
    background: var(--system-grouped-secondary);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .update-item {
    background: var(--system-grouped-secondary);
    border-color: rgba(255, 255, 255, 0.08);
  }

  .update-item:hover {
    background-color: rgba(46, 125, 50, 0.1);
  }

  .update-time {
    background: rgba(60, 60, 60, 0.5);
    border-color: rgba(255, 255, 255, 0.08);
  }

  .feature-card {
    border-color: rgba(255, 255, 255, 0.1);
    background: var(--system-grouped-secondary);
  }

  .feature-card:hover {
    border-color: var(--primary);
    background-color: rgba(46, 125, 50, 0.1);
    box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  }

  .feature-card h3 {
    color: #000;
  }

  .feature-card p {
    color: #000;
  }
}
</style> 