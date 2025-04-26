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
              <el-button type="primary" size="large" class="banner-button">立即探索</el-button>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 内容区域 -->
    <div class="content">
      <!-- 推荐小说 -->
      <section class="section">
        <div class="section-header">
          <h2 class="display-small">编辑推荐</h2>
          <el-button text class="title-medium">查看更多</el-button>
        </div>
        <el-row :gutter="24">
          <el-col :span="4" v-for="book in recommendedBooks" :key="book.id">
            <div class="book-card shadow-md" @click="$router.push(`/novel/${book.id}`)">
              <div class="book-cover">
                <img :src="book.cover" :alt="book.title" />
              </div>
              <div class="book-info">
                <h3 class="title-medium">{{ book.title }}</h3>
                <p class="body-medium author">{{ book.author }}</p>
                <p class="caption description">{{ book.description }}</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- 最新更新 -->
      <section class="section">
        <div class="section-header">
          <h2 class="display-small">最新更新</h2>
          <el-button text class="title-medium">查看更多</el-button>
        </div>
        <div class="latest-updates">
          <div v-for="update in latestUpdates" :key="update.id" class="update-item">
            <div class="book-meta">
              <router-link :to="`/novel/${update.id}`" class="book-title title-medium">
                {{ update.title }}
              </router-link>
              <span class="author body-medium">{{ update.author }}</span>
            </div>
            <router-link 
              :to="`/novel/${update.id}/chapter/${update.latestChapterId}`" 
              class="chapter-link body-medium"
            >
              {{ update.latestChapter }}
            </router-link>
            <span class="update-time caption">{{ update.updateTime }}</span>
          </div>
        </div>
      </section>

      <!-- AI创作专区 -->
      <section class="section ai-section">
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
import { ref } from 'vue'
import { Edit, Brush, Tools } from '@element-plus/icons-vue'

// 模拟数据 - 轮播图
const banners = ref([
  {
    id: 1,
    image: '/banner1.jpg',
    title: '发现优质小说',
    description: '数千部精选小说，总有一本适合你'
  },
  {
    id: 2,
    image: '/banner2.jpg',
    title: 'AI智能创作',
    description: '让创作更简单，让故事更精彩'
  },
  {
    id: 3,
    image: '/banner3.jpg',
    title: '作者社区',
    description: '和其他作者交流，分享创作经验'
  }
])

// 模拟数据 - 推荐小说
const recommendedBooks = ref([
  {
    id: 1,
    title: '星际迷航',
    author: '张三',
    cover: '/book1.jpg',
    description: '一部精彩的科幻小说'
  },
  {
    id: 2,
    title: '魔法世界',
    author: '李四',
    cover: '/book2.jpg',
    description: '奇幻冒险故事'
  },
  // ... 更多书籍
])

// 模拟数据 - 最新更新
const latestUpdates = ref([
  {
    id: 1,
    title: '星际迷航',
    author: '张三',
    latestChapter: '第100章 惊天大战',
    latestChapterId: 100,
    updateTime: '2024-01-20 15:30'
  },
  {
    id: 2,
    title: '魔法世界',
    author: '李四',
    latestChapter: '第80章 魔法觉醒',
    latestChapterId: 80,
    updateTime: '2024-01-20 14:20'
  },
  // ... 更多更新
])

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
</script>

<style scoped>
.home {
  min-height: 100vh;
  background: var(--system-background);
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
  background: linear-gradient(to right, rgba(0, 0, 0, 0.6), transparent);
}

.banner-text {
  position: relative;
  max-width: 600px;
  padding: 40px;
  border-radius: 24px;
  color: white;
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
}

.content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.section {
  margin-bottom: 60px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.section-header h2 {
  color: var(--text-primary);
  margin: 0;
}

.book-card {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
  position: relative;
  height: 380px;
  display: flex;
  flex-direction: column;
}

.book-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 16px;
  background: linear-gradient(
    to bottom,
    transparent 0%,
    transparent 70%,
    rgba(0, 0, 0, 0.05) 100%
  );
  pointer-events: none;
}

.book-card:hover {
  transform: translateY(-4px);
  box-shadow: 
    0 4px 12px rgba(0, 0, 0, 0.05),
    0 12px 25px rgba(0, 0, 0, 0.07);
  border-color: var(--apple-blue);
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
  transition: transform 0.5s ease;
}

.book-card:hover .book-cover img {
  transform: scale(1.05);
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
}

.author {
  color: var(--text-secondary);
  margin-bottom: 8px;
  flex-shrink: 0;
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
  font-size: 13px;
  line-height: 1.5;
}

.latest-updates {
  background: var(--system-grouped-secondary);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.update-item {
  display: flex;
  align-items: center;
  padding: 20px;
  gap: 24px;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.3s ease;
  position: relative;
  background: linear-gradient(
    to right,
    transparent 0%,
    transparent 50%,
    var(--system-background) 100%
  );
  background-size: 200% 100%;
  background-position: 0% 0%;
}

.update-item:last-child {
  border-bottom: none;
}

.update-item:hover {
  background-position: 100% 0%;
  padding-left: 24px;
}

.update-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--apple-blue);
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
}

.book-title:hover {
  color: var(--apple-blue);
}

.author {
  color: var(--text-secondary);
}

.chapter-link {
  flex: 1;
  color: var(--text-primary);
  text-decoration: none;
}

.chapter-link:hover {
  color: var(--apple-blue);
}

.update-time {
  color: var(--text-tertiary);
  width: 120px;
  text-align: right;
}

.ai-section {
  margin-bottom: 80px;
}

.feature-card {
  padding: 40px;
  border-radius: 24px;
  text-align: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid var(--border-color);
  background: linear-gradient(
    145deg,
    var(--system-grouped-secondary) 0%,
    var(--system-background) 100%
  );
  position: relative;
  overflow: hidden;
}

.feature-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent,
    var(--apple-blue),
    transparent
  );
  opacity: 0;
  transition: opacity 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-8px);
  border-color: var(--apple-blue);
  box-shadow: 
    0 8px 16px rgba(0, 0, 0, 0.05),
    0 16px 32px rgba(0, 0, 0, 0.07);
}

.feature-card:hover::before {
  opacity: 1;
}

.feature-icon {
  margin-bottom: 24px;
  color: var(--apple-blue);
  position: relative;
}

.feature-icon::after {
  content: '';
  position: absolute;
  width: 60px;
  height: 60px;
  background: var(--apple-blue);
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
  color: var(--text-primary);
  margin: 0 0 16px;
}

.feature-card p {
  color: var(--text-secondary);
  margin: 0;
}

@media (prefers-color-scheme: dark) {
  .banner-content::before {
    background: linear-gradient(to right, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0.4));
  }

  .book-card {
    border-color: rgba(255, 255, 255, 0.1);
  }

  .book-card:hover {
    border-color: var(--apple-blue);
    box-shadow: 
      0 4px 12px rgba(0, 0, 0, 0.2),
      0 12px 25px rgba(0, 0, 0, 0.3);
  }

  .latest-updates {
    border-color: rgba(255, 255, 255, 0.1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .feature-card {
    border-color: rgba(255, 255, 255, 0.1);
    background: linear-gradient(
      145deg,
      var(--system-grouped-secondary) 0%,
      rgba(255, 255, 255, 0.05) 100%
    );
  }

  .feature-card:hover {
    border-color: var(--apple-blue);
    box-shadow: 
      0 8px 16px rgba(0, 0, 0, 0.2),
      0 16px 32px rgba(0, 0, 0, 0.3);
  }
}
</style> 