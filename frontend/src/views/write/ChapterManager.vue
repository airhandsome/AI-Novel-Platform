<template>
  <div class="chapter-manager">
    <div class="header">
      <div class="header-left">
        <h1>章节管理</h1>
        <div class="stats">
          <div class="stat-item">
            <span class="label">总章节</span>
            <span class="value">{{ chapters.length }}</span>
          </div>
          <div class="stat-item">
            <span class="label">总字数</span>
            <span class="value">{{ formatNumber(totalWordCount) }}</span>
          </div>
        </div>
      </div>
      <el-button type="primary" class="create-btn" @click="createChapter">
        <el-icon><Plus /></el-icon>新建章节
      </el-button>
    </div>

    <div class="chapter-list">
      <el-table 
        :data="chapters" 
        style="width: 100%"
        :header-cell-style="{
          background: '#f5f7fa',
          color: '#606266',
          fontWeight: 600
        }"
        :row-class-name="tableRowClassName"
      >
        <el-table-column prop="order" label="章节号" width="100" align="center">
          <template #default="{ row }">
            <span class="chapter-order"># {{ row.order }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="title" label="章节标题">
          <template #default="{ row }">
            <div class="chapter-title-cell">
              <span class="title">{{ row.title }}</span>
              <span class="update-time">{{ formatDate(row.updateTime) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="wordCount" label="字数" width="120" align="center">
          <template #default="{ row }">
            <span class="word-count">{{ formatNumber(row.wordCount) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" class="status-tag">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="320" align="center">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button type="primary" plain size="small" @click="editChapter(row.id)">
                <el-icon><Edit /></el-icon>编辑
              </el-button>
              <el-button type="info" plain size="small" @click="previewChapter(row.id)">
                <el-icon><View /></el-icon>预览
              </el-button>
              <el-button 
                v-if="row.status === 0"
                type="success" 
                plain
                size="small"
                @click="publishChapter(row.id)"
              >
                <el-icon><Upload /></el-icon>发布
              </el-button>
              <el-dropdown 
                trigger="click" 
                @command="(cmd) => handleCommand(cmd, row.id)"
                class="more-dropdown"
              >
                <el-button size="small">
                  <el-icon><More /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="moveUp">
                      <el-icon><Top /></el-icon>上移
                    </el-dropdown-item>
                    <el-dropdown-item command="moveDown">
                      <el-icon><Bottom /></el-icon>下移
                    </el-dropdown-item>
                    <el-dropdown-item command="delete" divided>
                      <el-icon><Delete /></el-icon>删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 预览对话框 -->
    <el-dialog 
      v-model="previewVisible" 
      :title="previewChapterData.title"
      width="80%"
      class="preview-dialog"
      :close-on-click-modal="false"
      :destroy-on-close="true"
    >
      <div class="preview-content">
        <div class="chapter-content" :style="previewStyle">
          <div class="content-wrapper">
            <h1 class="chapter-title">{{ previewChapterData.title }}</h1>
            <div class="content-body" v-html="previewChapterData.content"></div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="preview-footer">
          <el-button @click="previewVisible = false">关闭</el-button>
          <el-button type="primary" @click="editChapter(previewChapterData.id)">编辑此章节</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, ArrowDown, Edit, View, Upload, More, 
  Top, Bottom, Delete 
} from '@element-plus/icons-vue'
import { novelApi } from '../../api/novel'
import { formatDate } from '../../utils/date'

const route = useRoute()
const router = useRouter()
const novelId = route.params.id

const chapters = ref([])
const previewVisible = ref(false)
const previewChapterData = ref({})

// 获取章节列表
const fetchChapters = async () => {
  try {
    const response = await novelApi.getNovelChapters({novelId: novelId})
    // 从response中获取chapters数组
    chapters.value = response.chapters || []
  } catch (error) {
    console.error('获取章节列表失败:', error)
    ElMessage.error('获取章节列表失败')
  }
}

// 创建新章节
const createChapter = () => {
  router.push(`/write/editor/${novelId}`)
}

// 编辑章节
const editChapter = (chapterId) => {
  router.push(`/write/editor/${novelId}?chapterId=${chapterId}`)
}

// 预览章节
const previewChapter = async (chapterId) => {
  try {
    const response = await novelApi.getChapter(chapterId)
    previewChapterData.value = response
    previewVisible.value = true
  } catch (error) {
    console.error('获取章节内容失败:', error)
    ElMessage.error('获取章节内容失败')
  }
}

// 处理更多操作
const handleCommand = async (command, chapterId) => {
  switch (command) {
    case 'moveUp':
      await moveChapter(chapterId, 'up')
      break
    case 'moveDown':
      await moveChapter(chapterId, 'down')
      break
    case 'delete':
      await deleteChapter(chapterId)
      break
  }
}

// 移动章节
const moveChapter = async (chapterId, direction) => {
  try {
    await novelApi.moveChapter(chapterId, direction)
    await fetchChapters()
    ElMessage.success('移动成功')
  } catch (error) {
    console.error('移动章节失败:', error)
    ElMessage.error('移动章节失败')
  }
}

// 删除章节
const deleteChapter = async (chapterId) => {
  try {
    await ElMessageBox.confirm('确定要删除这个章节吗？此操作不可恢复', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await novelApi.deleteChapter(chapterId)
    await fetchChapters()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除章节失败:', error)
      ElMessage.error('删除章节失败')
    }
  }
}

// 发布章节
const publishChapter = async (chapterId) => {
  try {
    await ElMessageBox.confirm('确定要发布这个章节吗？发布后将对读者可见', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await novelApi.updateChapterStatus(chapterId, 1) // 1表示已发布状态
    await fetchChapters()
    ElMessage.success('发布成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('发布章节失败:', error)
      ElMessage.error('发布章节失败')
    }
  }
}

// 格式化数字
const formatNumber = (num) => {
  if (!num) return 0
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num
}

// 获取状态类型
const getStatusType = (status) => {
  const types = {
    0: 'info',    // 草稿
    1: 'success', // 已发布
    2: 'warning'  // 待审核
  }
  return types[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    0: '草稿',
    1: '已发布',
    2: '待审核'
  }
  return texts[status] || '未知'
}

// 计算总字数
const totalWordCount = computed(() => {
  return chapters.value.reduce((sum, chapter) => sum + (chapter.wordCount || 0), 0)
})

// 表格行样式
const tableRowClassName = ({ row, rowIndex }) => {
  return row.status === 1 ? 'published-row' : ''
}

// 预览样式
const previewStyle = {
  fontSize: '16px',
  lineHeight: '1.8',
  background: '#ffffff',
  color: '#333333',
  width: '100%',
  maxWidth: '800px',
  margin: '0 auto',
  padding: '20px 40px',
  boxSizing: 'border-box',
  height: '60vh',
  overflowY: 'auto'
}

onMounted(() => {
  fetchChapters()
})
</script>

<style scoped>
:root {
  --primary-green: #42b983;
  --primary-green-light: #4caf50;
  --primary-green-dark: #388e3c;
  --success-green: #67c23a;
  --success-light: #f0f9eb;
  --text-primary: #2c3e50;
  --text-secondary: #606266;
  --text-tertiary: #909399;
  --bg-primary: #f9fafb;
  --bg-secondary: #ecf5ec;
  --border-color: #e4ebe4;
}

.chapter-manager {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: var(--bg-primary);
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(66, 185, 131, 0.1);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 24px;
  background: var(--bg-secondary);
  border-radius: 12px;
  border: 1px solid var(--border-color);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 40px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  color: var(--primary-green-dark);
  font-weight: 600;
}

.stats {
  display: flex;
  gap: 32px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item .label {
  font-size: 14px;
  color: var(--text-tertiary);
}

.stat-item .value {
  font-size: 22px;
  font-weight: 600;
  color: var(--primary-green);
}

.create-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  font-size: 16px;
  background-color: var(--primary-green) !important;
  border-color: var(--primary-green) !important;
}

.create-btn:hover {
  background-color: var(--primary-green-dark) !important;
  border-color: var(--primary-green-dark) !important;
}

.chapter-list {
  margin-bottom: 30px;
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(66, 185, 131, 0.08);
}

.chapter-order {
  font-family: 'Monaco', monospace;
  font-weight: 600;
  color: var(--primary-green);
}

.chapter-title-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.chapter-title-cell .title {
  font-weight: 500;
  color: var(--text-primary);
}

.chapter-title-cell .update-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.word-count {
  font-family: 'Monaco', monospace;
  color: var(--success-green);
  font-weight: 500;
}

.status-tag {
  padding: 4px 16px;
  border-radius: 16px;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.action-buttons :deep(.el-button--primary) {
  --el-button-bg-color: var(--primary-green);
  --el-button-border-color: var(--primary-green);
  --el-button-hover-bg-color: var(--primary-green-dark);
  --el-button-hover-border-color: var(--primary-green-dark);
}

.action-buttons :deep(.el-button--success) {
  --el-button-bg-color: var(--success-green);
  --el-button-border-color: var(--success-green);
}

:deep(.el-table) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.el-table th) {
  background-color: var(--bg-secondary) !important;
  color: var(--primary-green-dark) !important;
  font-weight: 600;
}

:deep(.el-table__row) {
  transition: all 0.3s ease;
}

:deep(.published-row) {
  background-color: var(--success-light);
}

:deep(.el-table__row:hover) {
  background-color: var(--bg-secondary) !important;
  transform: translateY(-2px);
}

/* 预览对话框样式 */
.preview-dialog {
  :deep(.el-dialog__header) {
    background: var(--bg-secondary);
    margin: 0;
    padding: 20px;
    border-bottom: 1px solid var(--border-color);
  }
  
  :deep(.el-dialog__title) {
    color: var(--primary-green-dark);
    font-weight: 600;
  }
}

.preview-content {
  background: var(--bg-primary);
  padding: 24px;
}

.chapter-content {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(66, 185, 131, 0.1);
}

.content-wrapper {
  max-width: 100%;
}

.chapter-title {
  text-align: center;
  margin-bottom: 30px;
  font-size: 24px;
  color: var(--primary-green-dark);
  font-weight: 600;
}

.content-body {
  text-align: justify;
  word-break: break-word;
  line-height: 1.8;
  color: var(--text-primary);
}

.content-body :deep(p) {
  margin: 1.2em 0;
  text-indent: 2em;
}

.preview-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  background: #fff;
  border-top: 1px solid var(--border-color);
}

:deep(.el-dropdown-menu__item) {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
}

:deep(.el-dropdown-menu__item:hover) {
  color: var(--primary-green);
  background-color: var(--bg-secondary);
}

:deep(.el-button .el-icon) {
  margin-right: 4px;
}

:deep(.el-table__row .el-button--primary.is-plain:hover) {
  color: var(--primary-green);
  border-color: var(--primary-green);
  background: var(--bg-secondary);
}

:deep(.el-button--primary.is-plain) {
  --el-button-hover-text-color: var(--primary-green);
  --el-button-hover-border-color: var(--primary-green);
  --el-button-hover-bg-color: var(--bg-secondary);
}
</style> 