<template>
  <div class="chapter-manager">
    <div class="header">
      <h1>章节管理</h1>
      <el-button type="primary" @click="createChapter">
        <el-icon><Plus /></el-icon>新建章节
      </el-button>
    </div>

    <div class="chapter-list">
      <el-table :data="chapters" style="width: 100%">
        <el-table-column prop="order" label="章节号" width="100">
          <template #default="{ row }">
            {{ row.order }}
          </template>
        </el-table-column>
        
        <el-table-column prop="title" label="章节标题" />
        
        <el-table-column prop="wordCount" label="字数" width="120">
          <template #default="{ row }">
            {{ formatNumber(row.wordCount) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="updateTime" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updateTime) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="300">
          <template #default="{ row }">
            <el-button @click="editChapter(row.id)">编辑</el-button>
            <el-button @click="previewChapter(row.id)">预览</el-button>
            <el-button 
              v-if="row.status === 0"
              type="success" 
              @click="publishChapter(row.id)"
            >发布</el-button>
            <el-dropdown trigger="click" @command="(cmd) => handleCommand(cmd, row.id)">
              <el-button>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="moveUp">上移</el-dropdown-item>
                  <el-dropdown-item command="moveDown">下移</el-dropdown-item>
                  <el-dropdown-item command="delete" divided>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 预览对话框 -->
    <el-dialog v-model="previewVisible" title="章节预览" width="60%">
      <h2>{{ previewChapterData.title }}</h2>
      <div class="chapter-content">{{ previewChapterData.content }}</div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, ArrowDown } from '@element-plus/icons-vue'
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
    const response = await novelApi.getNovelChapters(novelId)
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

onMounted(() => {
  fetchChapters()
})
</script>

<style scoped>
.chapter-manager {
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

.chapter-list {
  margin-bottom: 30px;
}

.chapter-content {
  white-space: pre-wrap;
  line-height: 1.6;
  max-height: 500px;
  overflow-y: auto;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 4px;
}
</style> 